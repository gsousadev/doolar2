package presentation

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func UploadAudio(w http.ResponseWriter, r *http.Request) {

	defer func() {
		if rec := recover(); rec != nil {
			fmt.Printf("Recovered from panic: %v\n", rec)
		}
	}()

	// IMPORTANTE: Processa o arquivo ANTES de configurar headers de streaming
	transcription, err := sendAudioFileToWhisper(r, w)
	if err != nil {
		fmt.Println("Erro transcrevendo:", err)
		http.Error(w, "Falha ao transcrever o áudio", http.StatusInternalServerError)
		return
	}

	fmt.Println("Transcrição recebida:", transcription)

	// Agora configura streaming para a resposta do Ollama
	w.Header().Set("Content-Type", "application/x-ndjson")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming não suportado", http.StatusInternalServerError)
		return
	}

	streamTextFromOllama(w, flusher, transcription)
}

func streamTextFromOllama(w http.ResponseWriter, flusher http.Flusher, text string) {

	prompt := fmt.Sprintf(`Você é um assistente que extrai informações de tarefas de áudio transcrito.

Analise o texto abaixo e extraia as informações de uma tarefa em formato JSON estruturado:

Texto transcrito: "%s"

Retorne APENAS um JSON válido com esta estrutura exata (sem markdown, sem explicações):
{
  "title": "título curto da tarefa (máximo 100 caracteres)",
  "description": "descrição detalhada da tarefa",
  "status": "pending",
  "created_at": "data e hora atual no formato ISO 8601"
}

Regras:
- Se o texto não mencionar uma tarefa clara, use o conteúdo como descrição e crie um título resumido
- Status sempre "pending" por padrão
- Data atual: %s
- Não adicione comentários ou texto extra, apenas o JSON`, text, time.Now().Format(time.RFC3339))

	fmt.Println("🔵 Enviando prompt para Ollama (stream)...")

	reqBody, err := json.Marshal(OllamaRequest{
		Model:  "deepseek-r1",
		Prompt: prompt,
		Stream: true,
	})
	if err != nil {
		fmt.Printf("❌ Erro ao criar request body: %v\n", err)
		return
	}

	req, err := http.NewRequest("POST", "http://ollama:11434/api/generate", bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Printf("❌ Erro ao criar request: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 0, // Sem timeout para streaming
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ Erro ao conectar com Ollama: %v\n", err)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Falha ao conectar com o serviço de IA",
		})
		flusher.Flush()
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("❌ Ollama retornou status %d\n", resp.StatusCode)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Serviço de IA retornou erro: %d", resp.StatusCode),
		})
		flusher.Flush()
		return
	}

	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024) // Buffer maior para chunks grandes

	fmt.Println("📥 Iniciando leitura do stream do Ollama...")
	chunkCount := 0

	for scanner.Scan() {
		chunkCount++
		line := scanner.Bytes()

		fmt.Printf("📦 Chunk #%d recebido (%d bytes)\n", chunkCount, len(line))

		if len(line) == 0 {
			fmt.Println("⚠️ Chunk vazio, ignorando...")
			continue
		}

		var chunk OllamaResponse
		if err := json.Unmarshal(line, &chunk); err != nil {
			fmt.Printf("⚠️ Erro ao decodificar chunk: %v - linha: %s\n", err, string(line))
			continue
		}

		fmt.Printf("✉️ Chunk decodificado - response: '%s' (len=%d), done: %v\n",
			chunk.Response, len(chunk.Response), chunk.Done)

		// Envia o chunk completo (não apenas o texto)
		responseData := map[string]interface{}{
			"response": chunk.Response,
			"done":     chunk.Done,
		}

		if chunk.Response != "" || chunk.Done {
			fmt.Printf("📤 Enviando chunk #%d para cliente...\n", chunkCount)

			// Tenta enviar, mas não quebra se houver erro de timeout
			encoder := json.NewEncoder(w)
			if err := encoder.Encode(responseData); err != nil {
				fmt.Printf("⚠️ Cliente desconectou ou timeout: %v\n", err)
				return // Cliente desconectou, para o streaming
			}

			// Flush imediatamente após cada chunk
			flusher.Flush()
			fmt.Printf("✅ Chunk #%d enviado e flushed\n", chunkCount)
		} else {
			fmt.Printf("⏭️ Chunk #%d sem conteúdo, pulando envio\n", chunkCount)
		}

		if chunk.Done {
			fmt.Printf("🏁 Stream finalizado após %d chunks\n", chunkCount)
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("❌ Erro ao ler stream: %v\n", err)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Erro ao processar resposta da IA",
			"done":  "true",
		})
		flusher.Flush()
	}
}

func sendAudioFileToWhisper(r *http.Request, w http.ResponseWriter) (string, error) {

	// Pega arquivo do form (isso faz o parse internamente)
	file, handler, err := r.FormFile("audio")
	if err != nil {
		return "", fmt.Errorf("erro ao obter arquivo: %w", err)
	}
	defer file.Close()

	fmt.Printf("Recebido: %s (%d bytes)\n", handler.Filename, handler.Size)

	// Detecta extensão do arquivo
	ext := filepath.Ext(handler.Filename)
	if ext == "" {
		ext = ".webm" // Default
	}

	// Gerar nome do arquivo local
	filename := "audio-" + time.Now().Format("20060102-150405") + ext

	filePath := "/app/internal/tasks/uploads/" + filename

	// Cria o arquivo local
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("erro ao criar arquivo local: %w", err)
	}
	defer dst.Close()

	// Copia o conteúdo enviado
	_, err = io.Copy(dst, file)
	if err != nil {
		return "", fmt.Errorf("erro ao salvar arquivo: %w", err)
	}

	fmt.Println("Arquivo salvo:", filePath)

	// ------------------------------------------
	// 🔥 ENVIA PARA A API WHISPER PYTHON
	// ------------------------------------------

	savedFile, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("erro ao reabrir arquivo: %w", err)
	}
	defer savedFile.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("audio", filepath.Base(filePath))
	if err != nil {
		return "", fmt.Errorf("erro ao criar form file: %w", err)
	}

	_, err = io.Copy(part, savedFile)
	if err != nil {
		return "", fmt.Errorf("erro ao copiar arquivo: %w", err)
	}

	writer.Close()

	req, err := http.NewRequest("POST", "http://whisper-asr:8000/transcribe", body)
	if err != nil {
		return "", fmt.Errorf("erro ao criar request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("erro ao conectar com Whisper: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("whisper retornou status %d", resp.StatusCode)
	}

	whisperBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("erro ao ler resposta do Whisper: %w", err)
	}

	return string(whisperBody), nil
}

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

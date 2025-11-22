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

func streamTextFromOllama(w http.ResponseWriter, flusher http.Flusher, text string) {

	prompt := fmt.Sprintf(`repita o texto. Apensar formate ele e não retorne nenhum dado a mais: %s`, text)

	fmt.Println("🔵 Enviando prompt para Ollama (stream)...")

	reqBody, _ := json.Marshal(OllamaRequest{
		Model:  "deepseek-r1",
		Prompt: prompt,
		Stream: true, // STREAM ATIVADO! 🔥
	})

	req, _ := http.NewRequest("POST", "http://ollama:11434/api/generate", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 0, // NÃO TER TIMEOUT PARA STREAM
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)

	for {
		line, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}

		println(string(line))

		var chunk OllamaResponse
		if err := json.Unmarshal(line, &chunk); err != nil {
			continue
		}

		if chunk.Response != "" {
			// ENVIA PARA O CLIENTE EM TEMPO REAL
			fmt.Println(chunk.Response)

			json.NewEncoder(w).Encode(chunk.Response)

			flusher.Flush()
		}
	}
}

func UploadAudio(w http.ResponseWriter, r *http.Request) {

	defer func() {
		if rec := recover(); rec != nil {
			fmt.Printf("Recovered from panic: %v\n", rec)
		}
	}()

	w.Header().Set("Content-Type", "application/x-ndjson")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming não suportado", http.StatusInternalServerError)
		return
	}

	transcription, err := sendAudioFileToWhisper(r, w)
	if err != nil {
		fmt.Println("Erro transcrevendo:", err)
		RespondError(w, http.StatusInternalServerError, "Falha ao transcrever o áudio")
		return
	}

	fmt.Println("Transcrição recebida:", transcription)

	streamTextFromOllama(w, flusher, transcription)
}

func sendAudioFileToWhisper(r *http.Request, w http.ResponseWriter) (string, error) {

	err := r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Erro ao processar o upload")
		return "", err
	}

	// Pega arquivo do form
	file, handler, err := r.FormFile("audio")
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Arquivo não encontrado no campo 'audio'")
		return "", err
	}
	defer file.Close()

	fmt.Printf("Recebido: %s (%d bytes)\n", handler.Filename, handler.Size)

	// Gerar nome do arquivo local
	filename := "audio-" + time.Now().Format("20060102-150405") + ".webm"

	filePath := "./uploads/" + filename

	// Cria o arquivo local
	dst, err := os.Create(filePath)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "Erro ao criar arquivo local")
		return "", err
	}
	defer dst.Close()

	// Copia o conteúdo enviado
	_, err = io.Copy(dst, file)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "Erro ao salvar arquivo no servidor")
		return "", err
	}

	fmt.Println("Arquivo salvo:", filePath)

	// ------------------------------------------
	// 🔥 ENVIA PARA A API WHISPER PYTHON
	// ------------------------------------------

	file, err = os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("audio", filepath.Base(filePath))
	if err != nil {
		return "", err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return "", err
	}

	writer.Close()

	req, err := http.NewRequest("POST", "http://whisper-asr:8000/transcribe", body)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	whisperBody, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
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

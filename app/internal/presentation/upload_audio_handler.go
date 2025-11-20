package presentation

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func UploadAudio(w http.ResponseWriter, r *http.Request) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic:", r)
		}
	}()

	// Limita o tamanho do upload (ex: 10MB)
	r.ParseMultipartForm(10 << 20)

	// "audio" é o nome do campo no form
	file, handler, err := r.FormFile("audio")
	if err != nil {
		fmt.Println("Erro ao recuperar o arquivo do formulário:", err)
		RespondError(w, http.StatusInternalServerError, "Erro ao salvar o arquivo")
		return
	}

	defer file.Close()

	fmt.Printf("Recebido: %s (%d bytes)\n", handler.Filename, handler.Size)

	filename := "audio-" + time.Now().Format(time.RFC3339) + ".webm"

	// Criar arquivo local
	dst, err := os.Create("./uploads/" + filename)

	if err != nil {
		RespondError(w, http.StatusInternalServerError, "Erro ao criar o arquivo local")
		return
	}

	fmt.Printf("Arquivo salvo em: %s\n", "./uploads/"+filename)

	defer dst.Close()

	// Copiar o conteúdo
	_, err = io.Copy(dst, file)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "Erro ao copiar o arquivo")
		return
	}

	fmt.Printf("Upload feito com sucesso: %s\n", "./uploads/"+filename)

	RespondSuccess(w, http.StatusCreated, "Upload feito com sucesso", nil)
}

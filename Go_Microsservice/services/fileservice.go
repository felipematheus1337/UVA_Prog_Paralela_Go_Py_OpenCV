package services

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path"
	"strings"
	"sync"
)

var tamanhoMaximoDoArquivo = 10 * 1024 * 1024
var pyServerBaseURL = "http://localhost:8087/upload_image"

type UploadedFile struct {
	Name     string
	Contents *bytes.Reader
}

func verificarTamanhoPermitido(tamanhoArquivo int64) bool {
	if tamanhoArquivo > int64(tamanhoMaximoDoArquivo) {
		return false
	} else {
		return true
	}
}

func verificarExtensaoPermitida(filename string) bool {
	filename = strings.TrimPrefix(filename, ".")
	if filename != "jpg" && filename != "png" && filename != "jpeg" {
		return false
	} else {
		return true
	}
}

func EnviarImagem(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		var waitGroup sync.WaitGroup

		waitGroup.Add(2)

		tamanho := make(chan bool)
		extensao := make(chan bool)

		go func() {
			tamanho <- verificarTamanhoPermitido(r.ContentLength)
			waitGroup.Done()

		}()

		arquivo, cabecalho, err := r.FormFile("file")

		if err != nil {
			fmt.Println("Erro arqv: ", err)
			return
		}

		defer arquivo.Close()

		extensaoArquivo := strings.ToLower(path.Ext(cabecalho.Filename))

		fmt.Println(extensaoArquivo)

		go func() {
			extensao <- verificarExtensaoPermitida(extensaoArquivo)
			waitGroup.Done()

		}()

		verifyTamanho := <-tamanho
		if !verifyTamanho {
			w.Write([]byte("Tamanho invalido, imagens até 10MB!"))
			return
		}

		verifyExtensao := <-extensao
		if !verifyExtensao {
			w.Write([]byte("Extensao invalida! Apenas png, jpg ou jpeg"))
			return
		}

		waitGroup.Wait()

		nomeDoArquivo := cabecalho.Filename

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)

		part, err := writer.CreateFormFile("image", nomeDoArquivo)

		if err != nil {
			fmt.Println("Error")
			return
		}

		_, err = io.Copy(part, arquivo)

		req, err := http.NewRequest("POST", "http://localhost:8777/upload_image", body)
		if err != nil {
			return
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())

		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			return
		}
		defer resp.Body.Close()

		// Se tudo ocorrer bem, a resposta deve ser "200 OK"
		if resp.StatusCode != http.StatusOK {
			return
		}

		fmt.Fprintln(w, "Arquivo  enviado com sucesso!", nomeDoArquivo)
	} else {
		http.Error(w, "Método não permitiodo", http.StatusMethodNotAllowed)
	}
}

/*
func EnviarImagens(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		var waitGroup sync.WaitGroup

		err := r.ParseMultipartForm(int64(tamanhoMaximoDoArquivo))

		if err != nil {
			fmt.Println(err)
			return
		}

		formulario := r.MultipartForm

		arquivos, ok := formulario.File["file"]
		if !ok {
			http.Error(w, "Nenhuma imagem encontrada", http.StatusBadRequest)
			return
		}

		numeroDeArquivos := len(arquivos)

		waitGroup.Add(numeroDeArquivos * 2)

		for _, arquivo := range arquivos {

			tamanho := make(chan bool)
			extensao := make(chan bool)

			go func(arquivo *multipart.FileHeader) {
				tamanho <- verificarTamanhoPermitido(arquivo.Size)
				waitGroup.Done()
			}(arquivo)

			extensaoArquivo := strings.ToLower(path.Ext(arquivo.Filename))

			go func(extensaoArquivo string) {
				extensao <- verificarExtensaoPermitida(extensaoArquivo)
				waitGroup.Done()
			}(extensaoArquivo)

			verifyTamanho := <-tamanho
			if !verifyTamanho {
				w.Write([]byte("Tamanho inválido, imagens até 10MB!"))
				return
			}

			verifyExtensao := <-extensao
			if !verifyExtensao {
				w.Write([]byte("Extensão inválida! Apenas png, jpg ou jpeg"))
				return
			}
		}

		waitGroup.Wait()

		for _, arquivo := range arquivos {
			nomeDoArquivo := arquivo.Filename
			fmt.Fprintf(w, "Arquivo %s enviado com sucesso!\n", nomeDoArquivo)
		}

		req, err := http.NewRequest("POST", "http://localhost:8777/upload_image", arquivos)

		if err != nil {
			return
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())

	} else {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}
*/

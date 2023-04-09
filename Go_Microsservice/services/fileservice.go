package services

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"path"
	"path/filepath"
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

		w.Write([]byte("Arquivo enviado com sucesso! " + nomeDoArquivo))
		fmt.Fprintln(w, "Arquivo  enviado com sucesso!", nomeDoArquivo)
	} else {
		http.Error(w, "Método não permitiodo", http.StatusMethodNotAllowed)
	}
}

func EnviarImagens(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(int64(tamanhoMaximoDoArquivo) * 2)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var waitGroup sync.WaitGroup

	files := r.MultipartForm.File["files"]

	waitGroup.Add(len(files))

	var imagens [][]byte

	for _, file := range files {

		if !verificarTamanhoPermitido(file.Size) {
			waitGroup.Done()
			http.Error(w, "O tamanho máximo permitido para cada arquivo é de 10 MB", http.StatusBadRequest)
			return
		}

		ext := filepath.Ext(file.Filename)

		if !(verificarExtensaoPermitida(ext)) {
			waitGroup.Done()
			http.Error(w, "A extensão do arquivo não é permitida", http.StatusBadRequest)
			return
		}

		fileContent, err := file.Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer fileContent.Close()

		fileBytes, err := ioutil.ReadAll(fileContent)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		imagens = append(imagens, fileBytes)
	}

	waitGroup.Wait()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for i, file := range files {
		part, err := writer.CreateFormFile("file", file.Filename)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		part.Write(imagens[i])
	}

	writer.Close()

	fmt.Println(body)

	req, err := http.NewRequest("POST", "http://localhost:8777/upload_images", body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Arquivos enviados com sucesso!"))
}

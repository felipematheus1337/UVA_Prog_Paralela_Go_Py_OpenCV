package services

import (
	"fmt"
	"net/http"
	"path/filepath"
)

var tamanhoMaximoDoArquivo = 10 * 1024 * 1024

func EnviarImagens(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		if r.ContentLength > int64(tamanhoMaximoDoArquivo) {
			w.Write([]byte("Não é possível fazer upload de arquivo com tamanho maior de 10mb"))
			return
		}

		arquivo, cabecalho, err := r.FormFile("arquivo")

		if err != nil {
			fmt.Println(err)
			return
		}

		extensao := filepath.Ext(cabecalho.Filename)

		if extensao != ".jpg" && extensao != ".png" && extensao != ".jpeg" {
			w.Write([]byte("Apenas arquivos .jpg, .png e .jpeg são permitidos"))
			return
		}

		defer arquivo.Close()

		nomeDoArquivo := cabecalho.Filename

		fmt.Println(arquivo)

		fmt.Fprintf(w, "Arquivo %s enviado com sucesso!", nomeDoArquivo)
	} else {
		http.Error(w, "Método não permitiodo", http.StatusMethodNotAllowed)
	}

}

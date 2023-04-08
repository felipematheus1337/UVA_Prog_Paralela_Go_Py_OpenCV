package services

import (
	"fmt"
	"net/http"
)

func EnviarImagens(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Chegou")
	w.Write([]byte("Chegou .."))
	w.WriteHeader(http.StatusOK)
}

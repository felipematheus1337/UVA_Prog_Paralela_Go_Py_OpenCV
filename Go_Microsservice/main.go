package main

import (
	"fmt"
	"log"
	"microsservice/services"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/upload/single", services.EnviarImagem).Methods(http.MethodPost)
	router.HandleFunc("/upload/multi", services.EnviarImagens).Methods(http.MethodPost)

	c := cors.Default()
	handler := c.Handler(router)

	fmt.Println("Escutando na porta 8081")
	log.Fatal(http.ListenAndServe(":8081", handler))
}

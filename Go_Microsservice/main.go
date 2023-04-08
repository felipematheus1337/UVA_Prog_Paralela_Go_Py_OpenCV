package main

import (
	"fmt"
	"log"
	"microsservice/services"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/upload/single", services.EnviarImagens).Methods(http.MethodPost)
	http.Handle("/", router)

	fmt.Println("Escutando na porta 8081")
	log.Fatal(http.ListenAndServe(":8081", router))

}

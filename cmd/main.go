package main

import (
	"log"
	"net/http"

	"github.com/lucasfarolfi/hire.me/infrastructure/webserver/handlers"
)

func main() {
	log.Println("Application starting...")

	shortenerHandler := handlers.NewShortenerHandler()

	http.HandleFunc("POST /shortener", shortenerHandler.Create)
	http.HandleFunc("GET /shortener/{alias}", shortenerHandler.Retrieve)

	log.Println("Server is running at port 8080")
	http.ListenAndServe(":8080", nil)
}

package main

import (
	"log"
	"net/http"

	"github.com/lucasfarolfi/hire.me/infrastructure/db"
	"github.com/lucasfarolfi/hire.me/infrastructure/webserver/handlers"
)

func main() {
	log.Println("Application starting...")

	db := db.InitializeDatabase()
	handler := handlers.NewURLShortenerHandler(db)

	http.HandleFunc("POST /shortener", handler.Create)
	http.HandleFunc("GET /shortener/{alias}", handler.Retrieve)

	log.Println("Server is running at port 8080")
	http.ListenAndServe(":8080", nil)
}

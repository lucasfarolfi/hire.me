package main

import (
	"log"
	"net/http"

	"github.com/lucasfarolfi/hire.me/infrastructure/db"
	"github.com/lucasfarolfi/hire.me/infrastructure/repository"
	"github.com/lucasfarolfi/hire.me/infrastructure/webserver/handlers"
	"github.com/lucasfarolfi/hire.me/internal/service"
)

func main() {
	log.Println("Application starting...")

	db := db.InitializeDatabase()
	repository := repository.NewShortenedURLRepository(db)
	service := service.NewURLShortenerService(repository)
	handler := handlers.NewURLShortenerHandler(service)

	http.HandleFunc("POST /", handler.Create)
	http.HandleFunc("GET /u/{alias}", handler.RetrieveByAlias)

	log.Println("Server is running at port 8080")
	http.ListenAndServe(":8080", nil)
}

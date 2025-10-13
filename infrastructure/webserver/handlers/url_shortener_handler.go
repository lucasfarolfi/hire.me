package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/lucasfarolfi/hire.me/internal/dto"
	"github.com/lucasfarolfi/hire.me/internal/entity"
	"gorm.io/gorm"
)

type URLShortenerHandler struct {
	DB *gorm.DB
}

func NewURLShortenerHandler(db *gorm.DB) *URLShortenerHandler {
	return &URLShortenerHandler{db}
}

func (h *URLShortenerHandler) Create(w http.ResponseWriter, r *http.Request) {
	s := &entity.ShortenedURL{Alias: "abcde4", Url: "http://www.bemobi.com.br"}

	res := dto.NewCreatedShortenedURLDTO(s.Alias, s.Url, "15ms")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, "failed to encode shortener response", http.StatusInternalServerError)
	}
}

func (h *URLShortenerHandler) Retrieve(w http.ResponseWriter, r *http.Request) {
	alias := r.PathValue("alias")
	log.Println(alias)
	s := &entity.ShortenedURL{Alias: "abcde", Url: "http://www.bemobi.com.br"}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(s)
	if err != nil {
		http.Error(w, "failed to encode shortener response", http.StatusInternalServerError)
	}
}

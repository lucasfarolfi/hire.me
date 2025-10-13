package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/lucasfarolfi/hire.me/internal/entity"
)

type ShortenerHandler struct {
}

func NewShortenerHandler() *ShortenerHandler {
	return &ShortenerHandler{}
}

func (h *ShortenerHandler) Create(w http.ResponseWriter, r *http.Request) {
	s := &entity.Shortener{Alias: "abcde", Url: "http://www.bemobi.com.br"}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err := json.NewEncoder(w).Encode(s)
	if err != nil {
		http.Error(w, "failed to encode shortener response", http.StatusInternalServerError)
	}
}

func (h *ShortenerHandler) Retrieve(w http.ResponseWriter, r *http.Request) {
	alias := r.PathValue("alias")
	log.Println(alias)
	s := &entity.Shortener{Alias: "abcde", Url: "http://www.bemobi.com.br"}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(s)
	if err != nil {
		http.Error(w, "failed to encode shortener response", http.StatusInternalServerError)
	}
}

package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/lucasfarolfi/hire.me/internal/dto"
	"github.com/lucasfarolfi/hire.me/internal/entity"
	"github.com/lucasfarolfi/hire.me/internal/service"
)

type URLShortenerHandler struct {
	service *service.URLShortenerService
}

func NewURLShortenerHandler(service *service.URLShortenerService) *URLShortenerHandler {
	return &URLShortenerHandler{service}
}

func (h *URLShortenerHandler) Create(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	url := r.URL.Query().Get("url")
	alias := r.URL.Query().Get("alias")
	if url == "" {
		http.Error(w, "url is required", http.StatusBadRequest)
		return
	}

	if alias == "" {
		alias = h.service.GenerateRandomAlias()
	} else if h.service.ExistsByAlias(alias) {
		retrieveErrorResponseBody(w, http.StatusBadRequest, "001", "CUSTOM ALIAS ALREADY EXISTS", alias)
		return
	}

	created, err := h.service.Create(alias, url)
	if err != nil {
		http.Error(w, "failed to create shortened URL", http.StatusInternalServerError)
		return
	}
	durationStr := fmt.Sprintf("%.3fms", float64(time.Since(startTime).Nanoseconds())/1e6)
	res := dto.NewCreatedShortenedURLDTO(created.Alias, created.Url, durationStr)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "failed to encode shortener response", http.StatusInternalServerError)
	}
}

func (h *URLShortenerHandler) RetrieveByAlias(w http.ResponseWriter, r *http.Request) {
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

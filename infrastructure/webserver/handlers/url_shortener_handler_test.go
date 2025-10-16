package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/lucasfarolfi/hire.me/infrastructure/repository"
	"github.com/lucasfarolfi/hire.me/internal/dto"
	"github.com/lucasfarolfi/hire.me/internal/entity"
	"github.com/lucasfarolfi/hire.me/internal/service"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestShortenerHandlerIntegration_Create(t *testing.T) {
	t.Run("Given a valid URL, when the API receives the request, then it should create a shortened URL", func(t *testing.T) {
		db := loadDB(t)
		service := service.NewURLShortenerService(repository.NewShortenedURLRepository(db))
		handler := NewURLShortenerHandler(service)

		server := httptest.NewServer(http.HandlerFunc(handler.Create))
		defer server.Close()

		params := url.Values{}
		params.Add("url", "http://www.bemobi.com.br")
		fullUrl := server.URL + "?" + params.Encode()

		resp, err := http.Post(fullUrl, "application/json", nil)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response dto.CreatedShortenedURLDTO
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)

		assert.Regexp(t, "^[a-zA-Z0-9]{6}$", response.Alias, "Alias should be a 6-character alphanumeric string")
		assert.Regexp(t, `^[0-9]*\.[0-9]+ms$`, response.Statistics.TimeTaken, "TimeTaken should be a positive duration in milliseconds")
		assert.Equal(t, params.Get("url"), response.URL, "The returned URL should match the input URL")
	})

	t.Run("Given a valid URL and an new optional custom alias, hen the API receives the request, then it should create a shortened URL using the custom alias", func(t *testing.T) {
		db := loadDB(t)
		service := service.NewURLShortenerService(repository.NewShortenedURLRepository(db))
		handler := NewURLShortenerHandler(service)

		server := httptest.NewServer(http.HandlerFunc(handler.Create))
		defer server.Close()

		params := url.Values{}
		params.Add("url", "http://www.bemobi.com.br")
		params.Add("alias", "XYhakR")
		fullUrl := server.URL + "?" + params.Encode()

		resp, err := http.Post(fullUrl, "application/json", nil)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response dto.CreatedShortenedURLDTO
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)

		assert.Equal(t, params.Get("alias"), response.Alias, "The returned Alias should match the input Alias")
		assert.Equal(t, params.Get("url"), response.URL, "The returned URL should match the input URL")
		assert.Regexp(t, `^[0-9]*\.[0-9]+ms$`, response.Statistics.TimeTaken, "TimeTaken should be a positive duration in milliseconds")
	})

	t.Run("Given a valid URL and an existent custom alias, hen the API receives the request, then it should return a custom error response", func(t *testing.T) {
		db := loadDB(t)
		service := service.NewURLShortenerService(repository.NewShortenedURLRepository(db))
		handler := NewURLShortenerHandler(service)

		service.Create("XYhakR", "http://www.abcde.com.br")

		server := httptest.NewServer(http.HandlerFunc(handler.Create))
		defer server.Close()

		params := url.Values{}
		params.Add("url", "http://www.bemobi.com.br")
		params.Add("alias", "XYhakR")
		fullUrl := server.URL + "?" + params.Encode()

		resp, err := http.Post(fullUrl, "application/json", nil)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var response HttpResponseErrorBody
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)

		assert.Equal(t, params.Get("alias"), response.Alias, "The returned Alias should match the input Alias")
		assert.Equal(t, "001", response.ErrCode, "ErrCode should be '001'")
		assert.Equal(t, "CUSTOM ALIAS ALREADY EXISTS", response.Description, "Description should indicate the alias already exists")
	})
}

func TestShortenerHandlerIntegration_RetrieveByAlias(t *testing.T) {
	t.Run("Given a valid alias, when the API receives the GET request, then it should retrieve the shortened URL", func(t *testing.T) {
		db := loadDB(t)
		service := service.NewURLShortenerService(repository.NewShortenedURLRepository(db))
		handler := NewURLShortenerHandler(service)

		mux := http.NewServeMux()
		mux.HandleFunc("GET /shortener/{alias}", handler.RetrieveByAlias)
		server := httptest.NewServer(mux)
		defer server.Close()

		alias := "abc123"
		url := "http://www.bemobi.com.br"
		db.Create(&entity.ShortenedURL{Alias: alias, Url: url})

		resp, err := http.Get(server.URL + "/shortener/" + alias)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response dto.ShortenedUrlRetrieveDTO
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)

		assert.Equal(t, url, response.URL, "The returned URL should match the stored URL")
	})

	t.Run("Given an non-existing alias, when the API receives the GET request, then it should return a custom error response", func(t *testing.T) {
		db := loadDB(t)
		service := service.NewURLShortenerService(repository.NewShortenedURLRepository(db))
		handler := NewURLShortenerHandler(service)

		mux := http.NewServeMux()
		mux.HandleFunc("GET /shortener/{alias}", handler.RetrieveByAlias)
		server := httptest.NewServer(mux)
		defer server.Close()

		alias := "non-existing"

		resp, err := http.Get(server.URL + "/shortener/" + alias)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var response HttpResponseErrorBody
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)

		assert.Equal(t, "002", response.ErrCode, "ErrCode should be '002'")
		assert.Equal(t, "SHORTENED URL NOT FOUND", response.Description, "Description should indicate the shortened URL was not found")
	})
}

func loadDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	err = db.AutoMigrate(&entity.ShortenedURL{})
	assert.NoError(t, err)
	return db
}

package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lucasfarolfi/hire.me/internal/dto"
	"github.com/lucasfarolfi/hire.me/internal/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestShortenerHandlerIntegration_Create(t *testing.T) {
	t.Run("Given a valid URL, when the API receives the request, then it should create a shortened UR", func(t *testing.T) {
		db := loadDB(t)
		handler := NewURLShortenerHandler(db)

		server := httptest.NewServer(http.HandlerFunc(handler.Create))
		defer server.Close()

		reqBody := &dto.URLShortenerCreateDTO{
			URL: "http://www.bemobi.com.br",
		}
		encodedBody, err := json.Marshal(reqBody)
		assert.NoError(t, err)

		resp, err := http.Post(server.URL, "application/json", bytes.NewReader(encodedBody))
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response dto.CreatedShortenedURLDTO
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)

		assert.Regexp(t, "^[a-zA-Z0-9]{6}$", response.Alias, "Alias should be a 6-character alphanumeric string")
		assert.Regexp(t, `^[1-9][0-9]*ms$`, response.Statistics.TimeTaken, "TimeTaken should be a positive duration in milliseconds")
		assert.Equal(t, reqBody.URL, response.URL, "The returned URL should match the input URL")
	})

	t.Run("Given a valid URL and an new optional custom alias, hen the API receives the request, then it should create a shortened URL using the custom alias", func(t *testing.T) {
		db := loadDB(t)
		handler := NewURLShortenerHandler(db)

		server := httptest.NewServer(http.HandlerFunc(handler.Create))
		defer server.Close()

		reqBody := &dto.URLShortenerCreateDTO{
			URL:   "http://www.bemobi.com.br",
			Alias: "XYhakR",
		}
		encodedBody, err := json.Marshal(reqBody)
		assert.NoError(t, err)

		resp, err := http.Post(server.URL, "application/json", bytes.NewReader(encodedBody))
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response dto.CreatedShortenedURLDTO
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)

		assert.Equal(t, reqBody.Alias, response.Alias, "The returned Alias should match the input Alias")
		assert.Equal(t, reqBody.URL, response.URL, "The returned URL should match the input URL")
		assert.Regexp(t, `^[1-9][0-9]*ms$`, response.Statistics.TimeTaken, "TimeTaken should be a positive duration in milliseconds")
	})

	t.Run("Given a valid URL and an existent custom alias, hen the API receives the request, then it should return a custom error response", func(t *testing.T) {
		db := loadDB(t)
		handler := NewURLShortenerHandler(db)

		server := httptest.NewServer(http.HandlerFunc(handler.Create))
		defer server.Close()

		reqBody := &dto.URLShortenerCreateDTO{
			URL:   "http://www.bemobi.com.br",
			Alias: "XYhakR",
		}
		encodedBody, err := json.Marshal(reqBody)
		assert.NoError(t, err)

		resp, err := http.Post(server.URL, "application/json", bytes.NewReader(encodedBody))
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var response HttpResponseErrorBody
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)

		assert.Equal(t, reqBody.Alias, response.Alias, "The returned Alias should match the input Alias")
		assert.Equal(t, "001", response.ErrCode, "ErrCode should be '001'")
		assert.Equal(t, "CUSTOM ALIAS ALREADY EXISTS", response.Description, "Description should indicate the alias already exists")
	})
}

func loadDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file:memory:"), &gorm.Config{})
	assert.NoError(t, err)
	err = db.AutoMigrate(&entity.ShortenedURL{})
	assert.NoError(t, err)
	return db
}

package handlers

import (
	"encoding/json"
	"fmt"
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

		assert.Regexp(t, "^[a-zA-Z0-9]{11}$", response.Alias, "Alias should be a 6-character alphanumeric string")
		assert.Regexp(t, `^[0-9]*\.[0-9]+ms$`, response.Statistics.TimeTaken, "TimeTaken should be a positive duration in milliseconds")
		assert.Equal(t, fmt.Sprintf("%s/u/%s", server.URL, response.Alias), response.URL, "The returned URL should match the input URL")
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
		assert.Equal(t, fmt.Sprintf("%s/u/%s", server.URL, params.Get("alias")), response.URL, "The returned URL should match the input URL")
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
		mux.HandleFunc("GET /u/{alias}", handler.RetrieveByAlias)
		server := httptest.NewServer(mux)
		defer server.Close()

		alias := "abc123"
		url := "http://www.bemobi.com.br"
		db.Create(&entity.ShortenedURL{Alias: alias, Url: url})

		resp, err := http.Get(server.URL + "/u/" + alias)
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
		mux.HandleFunc("GET /u/{alias}", handler.RetrieveByAlias)
		server := httptest.NewServer(mux)
		defer server.Close()

		alias := "non-existing"

		resp, err := http.Get(server.URL + "/u/" + alias)
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

func TestShortenerHandlerIntegration_CreatexRetrieve(t *testing.T) {
	t.Run("Given a valid alias and a url, when create is called followed by retrieve endpoint, then it should receive the shorten URL and redirect to the full URL", func(t *testing.T) {
		db := loadDB(t)
		service := service.NewURLShortenerService(repository.NewShortenedURLRepository(db))
		handler := NewURLShortenerHandler(service)

		mux := http.NewServeMux()
		mux.HandleFunc("POST /", handler.Create)
		mux.HandleFunc("GET /u/{alias}", handler.RetrieveByAlias)
		server := httptest.NewServer(mux)
		defer server.Close()

		alias := "abc123"
		urlToShort := "http://www.bemobi.com.br"

		params := url.Values{}
		params.Add("url", urlToShort)
		params.Add("alias", alias)
		fullUrl := server.URL + "?" + params.Encode()

		resp, err := http.Post(fullUrl, "application/json", nil)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var createResBody dto.CreatedShortenedURLDTO
		err = json.NewDecoder(resp.Body).Decode(&createResBody)
		assert.NoError(t, err)

		// Now get the full URL using the shorten url retrieved by create
		resp, err = http.Get(createResBody.URL)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var retrieveResBody dto.ShortenedUrlRetrieveDTO
		err = json.NewDecoder(resp.Body).Decode(&retrieveResBody)
		assert.NoError(t, err)

		assert.Equal(t, urlToShort, retrieveResBody.URL, "The returned URL should match the stored URL")
	})
}

func TestShortenerHandlerIntegration_GetMostAcessedUrls(t *testing.T) {
	t.Run("When GetMostAcessedUrls is called, then it should retrieve the 10 most accessed URLs ordering by access times DESC", func(t *testing.T) {
		db := loadDB(t)
		service := service.NewURLShortenerService(repository.NewShortenedURLRepository(db))
		handler := NewURLShortenerHandler(service)

		alias2 := "ABcdeF"
		service.Create(alias2, "http://www.bemobi.com.br")
		for i := 0; i < 3; i++ {
			_, err := service.RetrieveByAlias(alias2)
			assert.NoError(t, err)
		}

		alias3 := "123abc"
		service.Create(alias3, "http://www.example.com")
		for i := 0; i < 2; i++ {
			_, err := service.RetrieveByAlias(alias2)
			assert.NoError(t, err)
		}

		alias1 := "XYhakR"
		service.Create(alias1, "http://www.abcde.com.br")
		for i := 0; i < 5; i++ {
			_, err := service.RetrieveByAlias(alias1)
			assert.NoError(t, err)
		}

		mux := http.NewServeMux()
		mux.HandleFunc("GET /most_acessed", handler.RetrieveByAlias)
		server := httptest.NewServer(mux)
		defer server.Close()

		fullUrl := server.URL + "/most_acessed"
		resp, err := http.Get(fullUrl)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var resBody []dto.MostAcessedUrlDTO
		err = json.NewDecoder(resp.Body).Decode(resBody)
		assert.NoError(t, err)

		assert.Len(t, resBody, 3, "There should be 3 most accessed URLs")
		assert.Equal(t, "http://www.abcde.com.br", resBody[0].URL, "The most accessed URL should be first")
		assert.Equal(t, 5, resBody[0].AccessTimes, "The access times for the most accessed URL should be 5")
		assert.Equal(t, "http://www.bemobi.com.br", resBody[1].URL, "The second most accessed URL should be second")
		assert.Equal(t, 3, resBody[1].AccessTimes, "The access times for the second most accessed URL should be 3")
		assert.Equal(t, "http://www.example.com", resBody[2].URL, "The third most accessed URL should be third")
		assert.Equal(t, 2, resBody[2].AccessTimes, "The access times for the third most accessed URL should be 2")
	})
}

func loadDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	err = db.AutoMigrate(&entity.ShortenedURL{})
	assert.NoError(t, err)
	return db
}

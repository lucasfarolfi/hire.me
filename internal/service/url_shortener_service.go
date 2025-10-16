package service

import (
	"fmt"

	"github.com/lucasfarolfi/hire.me/internal/entity"
)

var ErrAliasAlreadyExists = fmt.Errorf("alias already exists")

type URLShortenerService struct {
	Repository ShortenedURLRepository
}

type ShortenedURLRepository interface {
	Create(shortUrl *entity.ShortenedURL) error
	FindByAlias(alias string) (*entity.ShortenedURL, error)
	ExistsByAlias(alias string) bool
	IncrementAccessTimesByID(id int) error
}

func NewURLShortenerService(repository ShortenedURLRepository) *URLShortenerService {
	return &URLShortenerService{repository}
}

func (s *URLShortenerService) GenerateRandomAlias() string {
	for {
		alias := "abcdef"
		if !s.Repository.ExistsByAlias("abcdef") {
			return alias
		}
	}
}

func (s *URLShortenerService) Create(alias, url string) (*entity.ShortenedURL, error) {
	shortenedUrl := entity.NewShortenedURL(alias, url)
	err := s.Repository.Create(shortenedUrl)
	if err != nil {
		return nil, err
	}
	return shortenedUrl, nil
}

func (s *URLShortenerService) RetrieveByAlias(alias string) (*entity.ShortenedURL, error) {
	return s.Repository.FindByAlias(alias)
}

func (s *URLShortenerService) ExistsByAlias(alias string) bool {
	if s.Repository.ExistsByAlias(alias) {
		return true
	}
	return false
}

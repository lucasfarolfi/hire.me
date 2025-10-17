package service

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"time"

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
		timestamp := uint64(time.Now().UnixMilli())
		randomPart := uint64(randomUint32())
		combined := (timestamp << 20) | randomPart
		alias := encodeBase62(combined)
		if !s.Repository.ExistsByAlias(alias) {
			return alias
		}
	}
}

func randomUint32() uint32 {
	var b [4]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(err)
	}
	return binary.BigEndian.Uint32(b[:])
}

func encodeBase62(num uint64) string {
	if num == 0 {
		return "0"
	}

	result := ""
	for num > 0 {
		remainder := num % 62
		result += string("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"[remainder])
		num /= 62
	}
	return result
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

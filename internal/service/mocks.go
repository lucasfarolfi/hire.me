package service

import (
	"github.com/lucasfarolfi/hire.me/internal/entity"
	"github.com/stretchr/testify/mock"
)

type MockShortenedURLRepository struct {
	mock.Mock
}

func (m *MockShortenedURLRepository) Create(shortUrl *entity.ShortenedURL) error {
	args := m.Called(shortUrl)
	return args.Error(0)
}

func (m *MockShortenedURLRepository) FindByAlias(alias string) (*entity.ShortenedURL, error) {
	args := m.Called(alias)
	if args.Get(0) != nil {
		return args.Get(0).(*entity.ShortenedURL), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockShortenedURLRepository) ExistsByAlias(alias string) bool {
	args := m.Called(alias)
	return args.Bool(0)
}

func (m *MockShortenedURLRepository) IncrementAccessTimesByID(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockShortenedURLRepository) Get10MostAcessedUrls() ([]entity.ShortenedURL, error) {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).([]entity.ShortenedURL), args.Error(1)
	}
	return nil, args.Error(1)
}

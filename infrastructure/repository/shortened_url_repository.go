package repository

import (
	"github.com/lucasfarolfi/hire.me/internal/entity"
	"gorm.io/gorm"
)

type ShortenedURLRepository struct {
	DB *gorm.DB
}

func NewShortenedURLRepository(db *gorm.DB) *ShortenedURLRepository {
	return &ShortenedURLRepository{DB: db}
}

func (ur *ShortenedURLRepository) Create(shortUrl *entity.ShortenedURL) error {
	return ur.DB.Create(shortUrl).Error
}

func (ur *ShortenedURLRepository) FindByAlias(alias string) (*entity.ShortenedURL, error) {
	var shortUrl entity.ShortenedURL
	err := ur.DB.Where("alias = ?", alias).First(&shortUrl).Error
	if err != nil {
		return nil, err
	}
	return &shortUrl, nil
}

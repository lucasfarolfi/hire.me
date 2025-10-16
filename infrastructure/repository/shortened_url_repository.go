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

func (ur *ShortenedURLRepository) ExistsByAlias(alias string) bool {
	var count int64
	err := ur.DB.Model(&entity.ShortenedURL{}).Where("alias = ?", alias).Count(&count).Error
	if err != nil || count == 0 {
		return false
	}
	return true
}

func (ur *ShortenedURLRepository) IncrementAccessTimesByID(id int) error {
	return ur.DB.Model(&entity.ShortenedURL{}).Where("id = ?", id).
		UpdateColumn("access_times", gorm.Expr("access_times + ?", 1)).Error
}

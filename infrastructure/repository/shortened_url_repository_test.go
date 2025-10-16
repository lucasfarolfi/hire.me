package repository

import (
	"testing"

	"github.com/lucasfarolfi/hire.me/internal/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestShortenerUrlRepositoryIntegration_Create(t *testing.T) {
	t.Run("Given a valid shortened url, when the create method is called, then it should create a shortened URL in database", func(t *testing.T) {
		db := loadDB(t)

		repository := NewShortenedURLRepository(db)

		shortUrl := &entity.ShortenedURL{
			Alias:       "abc123",
			Url:         "http://www.bemobi.com.br",
			AccessTimes: 0,
		}

		err := repository.Create(shortUrl)
		assert.NoError(t, err)

		shortUrlCreated := &entity.ShortenedURL{}
		err = db.First(shortUrlCreated, "alias = ?", shortUrl.Alias).Error
		assert.NoError(t, err)

		expected := shortUrl
		assert.Equal(t, expected, shortUrlCreated, "The created short URL should match the input short URL")
	})

	t.Run("Given a duplicate alias, when the Create method is called, then it should return an error", func(t *testing.T) {
		db := loadDB(t)
		repository := NewShortenedURLRepository(db)

		shortUrl := &entity.ShortenedURL{
			Alias:       "abc123",
			Url:         "http://www.bemobi.com.br",
			AccessTimes: 0,
		}
		err := repository.Create(shortUrl)
		assert.NoError(t, err)

		err = repository.Create(shortUrl)

		assert.Error(t, err, "An error should be returned when trying to create a duplicate alias")
	})
}

func TestShortenerUrlRepositoryIntegration_FindByAlias(t *testing.T) {
	t.Run("Given a valid alias that was already stored in database, when the FindByAlias method is called, then it should retrieve the stored shortened URL", func(t *testing.T) {
		db := loadDB(t)
		repository := NewShortenedURLRepository(db)

		shortUrl := &entity.ShortenedURL{
			Alias:       "abc123",
			Url:         "http://www.bemobi.com.br",
			AccessTimes: 0,
		}

		err := db.Create(shortUrl).Error
		assert.NoError(t, err)

		retrievedShortUrl, err := repository.FindByAlias(shortUrl.Alias)

		assert.NoError(t, err)
		expected := shortUrl
		assert.Equal(t, expected, retrievedShortUrl, "The retrieved short URL should match the input short URL")
	})

	t.Run("Given an alias that does not exist in the database, when the FindByAlias method is called, then it should return an error", func(t *testing.T) {
		db := loadDB(t)
		repository := NewShortenedURLRepository(db)

		nonExistentAlias := "nonexistent"

		retrievedShortUrl, err := repository.FindByAlias(nonExistentAlias)

		assert.Error(t, err, "An error should be returned when the alias does not exist in the database")
		assert.Nil(t, retrievedShortUrl, "The retrieved short URL should be nil when the alias does not exist")
	})
}

func loadDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	err = db.AutoMigrate(&entity.ShortenedURL{})
	assert.NoError(t, err)
	return db
}

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

func TestShortenedURLRepository_ExistsByAlias(t *testing.T) {
	t.Run("Given an alias that exists in the database, when ExistsByAlias is called, then it should return true", func(t *testing.T) {
		db := loadDB(t)
		repository := NewShortenedURLRepository(db)

		shortUrl := &entity.ShortenedURL{
			Alias:       "abc123",
			Url:         "http://www.bemobi.com.br",
			AccessTimes: 0,
		}
		err := db.Create(shortUrl).Error
		assert.NoError(t, err)

		exists := repository.ExistsByAlias("abc123")
		assert.True(t, exists, "ExistsByAlias should return true for an existing alias")
	})

	t.Run("Given an alias that does not exist in the database, when ExistsByAlias is called, then it should return false", func(t *testing.T) {
		db := loadDB(t)
		repository := NewShortenedURLRepository(db)

		exists := repository.ExistsByAlias("nonexistent")
		assert.False(t, exists, "ExistsByAlias should return false for a non-existing alias")
	})
}

func TestShortenedURLRepository_IncrementAccessTimes(t *testing.T) {
	t.Run("Given an alias that exists in the database, when IncrementAccessTimes is called, then it should increment the access times", func(t *testing.T) {
		db := loadDB(t)
		repository := NewShortenedURLRepository(db)

		shortUrl := &entity.ShortenedURL{
			Alias:       "abc123",
			Url:         "http://www.bemobi.com.br",
			AccessTimes: 0,
		}
		err := db.Create(shortUrl).Error
		assert.NoError(t, err)

		err = db.Where("alias = ?", "abc123").First(&shortUrl).Error
		assert.NoError(t, err)

		err = repository.IncrementAccessTimesByID(shortUrl.ID)
		assert.NoError(t, err)

		var updatedShortUrl entity.ShortenedURL
		err = db.First(&updatedShortUrl, "alias = ?", "abc123").Error
		assert.NoError(t, err)
		assert.Equal(t, int32(1), updatedShortUrl.AccessTimes, "AccessTimes should be incremented to 1")
	})
}

func TestShortenedURLRepository_Get10MostAcessedUrls(t *testing.T) {
	t.Run("When Get10MostAcessedUrls is called, then it should return the top 10 most accessed URLs", func(t *testing.T) {
		db := loadDB(t)
		repository := NewShortenedURLRepository(db)

		// Create 15 shortened URLs with varying access times
		for i := 1; i <= 15; i++ {
			shortUrl := &entity.ShortenedURL{
				Alias:       "alias" + string(i),
				Url:         "http://www.example" + string(i) + ".com",
				AccessTimes: int32(i * 10), // Access times: 10, 20, ..., 150
			}
			err := db.Create(shortUrl).Error
			assert.NoError(t, err)
		}

		mostAccessedUrls, err := repository.Get10MostAcessedUrls()
		assert.NoError(t, err)
		assert.Len(t, mostAccessedUrls, 10, "Should return exactly 10 most accessed URLs")

		// Verify that the URLs are ordered by AccessTimes descending
		for i := 0; i < len(mostAccessedUrls)-1; i++ {
			assert.GreaterOrEqual(t, mostAccessedUrls[i].AccessTimes, mostAccessedUrls[i+1].AccessTimes,
				"URLs should be ordered by AccessTimes in descending order")
		}
	})
}

func loadDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	err = db.AutoMigrate(&entity.ShortenedURL{})
	assert.NoError(t, err)
	return db
}

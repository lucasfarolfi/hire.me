package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestShortenerServiceUnit_GenerateRandomAlias(t *testing.T) {
	t.Run("Should keep generating a valid random alias until a non-existing alias is found in the database", func(t *testing.T) {
		repo := &MockShortenedURLRepository{}
		repo.On("ExistsByAlias", mock.AnythingOfType("string")).Return(true).Once()
		repo.On("ExistsByAlias", mock.AnythingOfType("string")).Return(true).Once()
		repo.On("ExistsByAlias", mock.AnythingOfType("string")).Return(false).Once()

		service := NewURLShortenerService(repo)

		resultedAlias := service.GenerateRandomAlias()

		assert.Regexp(t, "^[a-zA-Z0-9]{11}$", resultedAlias, "Alias should be a 6-character alphanumeric string")
		repo.AssertNumberOfCalls(t, "ExistsByAlias", 3)
	})
}

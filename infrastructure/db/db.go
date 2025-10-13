package db

import (
	"github.com/lucasfarolfi/hire.me/internal/entity"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitializeDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&entity.ShortenedURL{})
	if err != nil {
		panic(err)
	}
	return db
}

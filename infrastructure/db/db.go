package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/lucasfarolfi/hire.me/internal/entity"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitializeDatabase() *gorm.DB {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	if dbHost == "" || dbPort == "" || dbUser == "" || dbPassword == "" || dbName == "" {
		log.Fatal("One or more DB environment variable are not available!")
	}
	connString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbPort, dbName)

	var db *gorm.DB
	var err error
	for i := 0; i < 10; i++ {
		log.Println("Trying to connect to database... Attempt", i+1)
		db, err = gorm.Open(mysql.Open(connString), &gorm.Config{})
		if err == nil {
			sqlDB, err := db.DB()
			if err != nil {
				log.Println("Error while getting DB instance:", err)
				break
			}
			err = sqlDB.Ping()
			if err == nil {
				log.Println("Database connection established successfully!")
				break
			}
		}

		log.Println("Failed to connect, retrying...")
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&entity.ShortenedURL{})
	if err != nil {
		panic(err)
	}
	return db
}

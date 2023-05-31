package database

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/joho/godotenv"
)

var DBConn *gorm.DB

func InitDB() error {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	dsn := os.Getenv("POSTGRES_URL")
	DBConn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	return nil
}

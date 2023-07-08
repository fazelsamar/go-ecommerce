package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbConn *gorm.DB

func getDSN() string {
	host := os.Getenv("DATABASE_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	db := os.Getenv("POSTGRES_DB")
	port := os.Getenv("DATABASE_PORT")
	ssl := os.Getenv("POSTGRES_SSL")
	timezone := os.Getenv("POSTGRES_TIMEZONE")
	conStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", host, user, password, db, port, ssl, timezone)
	return conStr
}

func InitDB() {
	var err error
	dbConn, err = gorm.Open(postgres.Open(getDSN()), &gorm.Config{})
	if err != nil {
		panic(err)
	}
}

func GetDatabaseConnection() *gorm.DB {
	if dbConn == nil {
		InitDB()
	}
	return dbConn
}

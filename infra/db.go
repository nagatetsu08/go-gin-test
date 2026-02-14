package infra

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupDB() *gorm.DB {

	// テスト時はSQLiteを通常時はpostgresを使用する
	env := os.Getenv("ENV")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tokyo",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	var (
		db  *gorm.DB
		err error
	)

	if env == "prod" {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	} else {
		db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	}

	if err != nil {
		panic("Failed to connect database")
	}
	return db
}

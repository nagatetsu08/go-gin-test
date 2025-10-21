package main

import (
	"gin-freemarket/infra"
	"gin-freemarket/models"
)

func main() {
	infra.Initialize()

	db := infra.SetupDB()

	if err := db.AutoMigrate(&models.Item{}, &models.User{}); err != nil {
		panic("Failed to migrate database")
	}
}

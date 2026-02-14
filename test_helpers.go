package main

import (
	"gin-freemarket/infra"
	"gin-freemarket/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setUpTestData(db *gorm.DB) {
	items := []models.Item{
		{Name: "テストアイテム1", Price: 1000, Description: "テスト１", SoldOut: false, UserId: 1},
		{Name: "テストアイテム2", Price: 2000, Description: "テスト２", SoldOut: true, UserId: 1},
		{Name: "テストアイテム3", Price: 3000, Description: "テスト３", SoldOut: false, UserId: 2},
	}

	users := []models.User{
		{Email: "test1@example.com", Password: "test1Pass"},
		{Email: "test2@example.com", Password: "test2Pass"},
	}

	for _, user := range users {
		db.Create(&user)
	}

	for _, item := range items {
		db.Create(&item)
	}
}

func setUp() *gin.Engine {
	db := infra.SetupDB()

	// sqlite上のDBにテーブルを作成
	db.AutoMigrate(&models.Item{}, &models.User{})

	//テストデータ作成
	setUpTestData(db)

	// ルータ作成
	router := setUpRouter(db)

	return router
}

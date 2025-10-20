package main

import (
	"gin-freemarket/controllers"
	"gin-freemarket/models"
	"gin-freemarket/repositories"
	"gin-freemarket/services"

	"github.com/gin-gonic/gin"
)

func main() {
	items := []models.Item{
		{ID: 1, Name: "商品1", Price: 1000, Description: "説明1", SoldOut: false},
		{ID: 2, Name: "商品2", Price: 2000, Description: "説明2", SoldOut: true},
		{ID: 3, Name: "商品3", Price: 3000, Description: "説明3", SoldOut: false},
	}

	itemRepository := repositories.NewItemMemoryRepository(items) //DBインスタンスそのもの（リポジトリ）
	itemService := services.NewItemService(itemRepository)
	itemController := controllers.NewItemController(itemService)

	// エンドポイント設定
	router := gin.Default()
	router.GET("/items", itemController.FindAll)
	router.Run("localhost:8080") // デフォルトで0.0.0.0:8080で待機します

}

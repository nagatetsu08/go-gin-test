package main

import (
	"gin-freemarket/controllers"
	"gin-freemarket/infra"
	"gin-freemarket/repositories"
	"gin-freemarket/services"

	"github.com/gin-gonic/gin"
)

func main() {

	infra.Initialize()
	db := infra.SetupDB()

	// items := []models.Item{
	// 	{ID: 1, Name: "商品1", Price: 1000, Description: "説明1", SoldOut: false},
	// 	{ID: 2, Name: "商品2", Price: 2000, Description: "説明2", SoldOut: true},
	// 	{ID: 3, Name: "商品3", Price: 3000, Description: "説明3", SoldOut: false},
	// }

	// リポジトリ形式にしているので、切り替えが簡単（大元を変えればいいだけ。）
	// 実用的な例で言うと、モックで作っていた部分を本番ように差し替えたりするときに使える。

	// itemRepository := repositories.NewItemMemoryRepository(items) //サーバーのメモリをDB代わりにしたリポジトリ
	itemRepository := repositories.NewItemRepository(db) // DBを利用したリポジトリ
	itemService := services.NewItemService(itemRepository)
	itemController := controllers.NewItemController(itemService)

	authRepository := repositories.NewAuthRepository(db)
	authService := services.NewAuthService(authRepository)
	authController := controllers.NewAuthController(authService)

	// エンドポイント設定
	router := gin.Default()

	// ルーティングをグルーピング化する
	itemRouter := router.Group("/items")
	authRouter := router.Group("/auth")

	itemRouter.GET("/", itemController.FindAll)
	itemRouter.GET("/:id", itemController.FindById)
	itemRouter.POST("/", itemController.Create)
	itemRouter.PUT("/:id", itemController.Update)
	itemRouter.DELETE("/:id", itemController.Delete)

	authRouter.POST("/signup", authController.Signup)
	authRouter.POST("/login", authController.Login)

	router.Run("localhost:8080") // デフォルトで0.0.0.0:8080で待機します

}

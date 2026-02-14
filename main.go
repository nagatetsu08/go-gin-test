package main

import (
	"gin-freemarket/controllers"
	"gin-freemarket/infra"
	"gin-freemarket/middlewares"
	"gin-freemarket/repositories"
	"gin-freemarket/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setUpRouter(db *gorm.DB) *gin.Engine {
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
	router.Use(cors.Default()) // 全てのサイトからのアクセスを許可する（本番ではあまり良くない）

	// 以下本番設定
	// config := cors.DefaultConfig()
	// config.AllowOrigins = []string{"http://google.com"}
	// router.Use(cors.New(config))

	// ルーティングをグルーピング化する
	itemRouter := router.Group("/items")
	itemRoouteWithAuth := router.Group("/items", middlewares.AuthMiddleware(authService))
	authRouter := router.Group("/auth")

	itemRouter.GET("/", itemController.FindAll)
	itemRoouteWithAuth.GET("/:id", itemController.FindById)
	itemRoouteWithAuth.POST("/", itemController.Create)
	itemRoouteWithAuth.PUT("/:id", itemController.Update)
	itemRoouteWithAuth.DELETE("/:id", itemController.Delete)

	authRouter.POST("/signup", authController.Signup)
	authRouter.POST("/login", authController.Login)

	return router
}

func main() {

	infra.Initialize()
	db := infra.SetupDB()

	router := setUpRouter(db)

	router.Run("localhost:8080") // デフォルトで0.0.0.0:8080で待機します

}

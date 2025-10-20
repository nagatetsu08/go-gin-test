package controllers

// コントローラもリポジトリクラスと同じ構造。
//

import (
	"gin-freemarket/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Itemコントローラのインタフェース
// 何を実装すべきかをメソッド単位で書く
type IItemController interface {
	FindAll(ctx *gin.Context)
}

// コントローラクラスの実態（classに相当。goにはクラスの概念がない。。。）
// Laravelとかだとこの中にメソッドを書いて処理を書くイメージとなるが、goはfunctionとして外に切り出す。
// プロパティとして、services.IItemService型を投入できる「service」という名のプロパティを定義
type ItemController struct {
	service services.IItemService
}

// コンストラクタ
func NewItemController(service services.IItemService) IItemController {
	return &ItemController{service: service}
}

// コントローラメソッド
func (c *ItemController) FindAll(ctx *gin.Context) {
	items, err := c.service.FindAll()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": items})
}

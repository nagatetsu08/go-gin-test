package controllers

// コントローラもリポジトリクラスと同じ構造。
//

import (
	"gin-freemarket/dto"
	"gin-freemarket/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Itemコントローラのインタフェース
// 何を実装すべきかをメソッド単位で書く
type IItemController interface {
	FindAll(ctx *gin.Context)
	FindById(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
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

func (c *ItemController) FindById(ctx *gin.Context) {
	// パスパラメータで受け取ったものは全てstring型になるのでuintに変更してやる
	itemId, err := strconv.ParseUint(ctx.Param("id"), 10, 64)

	// パラメータチェック
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	//サービスクラスメソッド実行
	item, err := c.service.FindById(uint(itemId))

	if err != nil {
		if err.Error() == "Item is not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": item})
}

func (c *ItemController) Create(ctx *gin.Context) {

	// ユーザーからのパラメータ受取用の箱を準備
	var input dto.CreateItemInput

	// ctxの中に持っているリクエストデータをinput（dto.CreateItemInput）にバインドする（当てこむ）
	// dto.CreateItemInputに定義された対応するプロパティにマッピングされ、かつ同時にバリデーションも実施される
	// なお、以下のif文の書き方はGoらしい書き方の一つ。汎用的すぎる名前のスコープを極力狭めるために工夫らしい
	// if err := 式; err != nil {...}という書き方は覚えておいた方がいい。
	// ただし、実際に処理に使う値について、参照に値を上書きするような処理出ないと意味がなく、returnをしてしまうメソッドの場合は
	// 大人しく一旦変数で受け取ってから、if判定

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newItem, err := c.service.Create(input)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": newItem})
}

func (c *ItemController) Update(ctx *gin.Context) {

	// uintは環境依存（32ビット or 64ビット）となる。これはサーバ環境がどちらでも動くようにするため
	// 一方でuint64、uint32という型も存在し、上記のuintとは全く異なる型。
	// goはString型を直接uintにはキャストできず、ParseUintを噛ませる必要がある。ただ、これを使うと
	// uint32 or uint64のどちらかに変換されてしまう。uint64とuintは異なる型なのでそのまま引数で渡すと怒られる。
	// なので、uint64をuintにキャストしている。（メソッド側をuint64戻り値にしても解消する）

	itemId, err := strconv.ParseUint(ctx.Param("id"), 10, 64)

	// パラメータチェック
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	// ユーザーからのパラメータ受取用の箱を準備
	var input dto.UpdateItemInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateedItem, err := c.service.Update(uint(itemId), input)

	if err != nil {
		if err.Error() == "Item is not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": updateedItem})
}

func (c *ItemController) Delete(ctx *gin.Context) {
	itemId, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	err = c.service.Delete(uint(itemId))

	if err != nil {
		if err.Error() == "Item is not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error"})
		return
	}
	ctx.Status(http.StatusOK) // ステータスコードのみを返す
}

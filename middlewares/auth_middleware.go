package middlewares

import (
	"gin-freemarket/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authService services.IAuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// リクエストヘッダーのチェック
		header := ctx.GetHeader("Authorization")
		if header == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized) //AbortWithStatusを使用すると後続のミドルウェア及び処理もSTOPできる
			return
		}

		// AutorizationHeaderがBearerトークンかどうかをチェックする。
		if !strings.HasPrefix(header, "Bearer ") {
			ctx.AbortWithStatus(http.StatusUnauthorized) //AbortWithStatusを使用すると後続のミドルウェア及び処理もSTOPできる
			return
		}

		tokenString := strings.TrimPrefix(header, "Bearer ") //Bearトークンの取得
		user, err := authService.GetUserFromToken(tokenString)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized) //AbortWithStatusを使用すると後続のミドルウェア及び処理もSTOPできる
			return
		}

		// ユーザー情報をリクエストコンテキストにセット
		ctx.Set("user", user)

		// 次のミドルウェア なければメイン処理に移行する
		ctx.Next()
	}
}

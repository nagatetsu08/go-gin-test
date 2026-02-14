package main

import (
	"encoding/json"
	"gin-freemarket/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// goではTestで始まる関数がテスト関数と認識される
func TestFindAll(t *testing.T) {
	// テストのセットアップ
	router := setUp()

	// テストの結果を記録するレコーダを作成
	w := httptest.NewRecorder()

	// リクエストの生成（第一引数：メソッド、第二引数：path、第三引数：body(post時に必要)）
	req, _ := http.NewRequest("GET", "/items", nil)

	// APIリクエスト実行
	router.ServeHTTP(w, req)

	// API実行結果取得変数
	var res map[string][]models.Item

	// レスポンスレコーダからBodyを抽出してresに格納
	json.Unmarshal(w.Body.Bytes(), &res)

	// アサーション
	// 第一引数：testingオブジェクト、第二引数：予想値、第3引数：実値
	assert.Equal(t, http.StatusOK, w.Code) //ステータスコード
	assert.Equal(t, 3, len(res["data"]))   //件数
}

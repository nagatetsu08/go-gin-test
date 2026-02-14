package main

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

// TestMain関数はテストランナーによって全てのテストが実行される前に呼び出されるやつ
func TestMain(m *testing.M) {

	// テスト用の環境変数を読み込む
	if err := godotenv.Load(".env.test"); err != nil {
		log.Fatalln("Error loading .envtest")
	}

	code := m.Run()

	// テスト終了
	os.Exit(code)
}

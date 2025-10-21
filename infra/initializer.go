package infra

import (
	"log"

	"github.com/joho/godotenv"
)

func Initialize() {
	// .dnvを読み込む。引数でファイル名を指定できる。省略すると.envが読み込まれる
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading.env file")

	}
}

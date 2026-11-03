package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/laineus/digimon/api/controllers"
	"os"
)

func main() {
	// .env読込
	godotenv.Load()
	// ルーター準備
	router := gin.Default()

	// ルーティング
	router.GET("/digimon", controller.GetDigimon)
	router.POST("/digimon", controller.PostDigimon)

	// リクエスト受付開始
	portNumber := os.Getenv("PORT")
	router.Run(":" + portNumber)
}

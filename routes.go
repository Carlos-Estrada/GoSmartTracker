package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	serverPort := os.Getenv("PORT")
	if serverPort == "" {
		log.Fatal("$PORT must be set")
	}

	apiRouter := gin.Default()

	apiRouter.GET("/status", respondServerStatus)
	apiRouter.POST("/create", createTrackingItem)
	apiRouter.GET("/items", listTrackingItems)

	apiRouter.Run(":" + serverPort)
}

func respondServerStatus(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "up"})
}

func createTrackingItem(ctx *gin.Context) {
	ctx.JSON(http.StatusCreated, gin.H{"message": "Tracking item created"})
}

func listTrackingItems(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"items": []string{"item1", "item2"}})
}
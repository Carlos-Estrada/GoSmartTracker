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

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.Default()

	router.GET("/status", getStatus)
	router.POST("/create", createItem)
	router.GET("/items", getItems)

	router.Run(":" + port)
}

func getStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "up"})
}

func createItem(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "Item created"})
}

func getItems(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"items": []string{"item1", "item2"}})
}
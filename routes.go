package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.Default()

	// Define routes
	router.GET("/status", getStatus)
	router.POST("/create", createItem)
	router.GET("/items", getItems)

	// Start server
	router.Run(":" + port)
}

// getStatus checks the API status
func getStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "up"})
}

// createItem is a placeholder for item creation endpoint
func createItem(c *gin.Context) {
	// Mock implementation
	c.JSON(http.StatusCreated, gin.H{"message": "Item created"})
}

// getItems is a placeholder for fetching items endpoint
func getItems(c *gin.Context) {
	// Mock implementation
	c.JSON(http.StatusOK, gin.H{"items": []string{"item1", "item2"}})
}
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var usersCollection *mongo.Collection
var ctx = context.TODO()

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mongoUri := os.Getenv("MONGO_URI")

	clientOptions := options.Client().ApplyURI(mongoUri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	usersCollection = client.Database("testdb").Collection("users")

	logMessage("MongoDB connection established")
}

func logMessage(message string) {
	// Customize this function as needed for more advanced logging features
	log.Println(message)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func register(c *gin.Context) {
	var newUser User
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	hashedPassword, err := hashPassword(newUser.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": "Could not hash password"})
		return
	}

	newUser.Password = hashedPassword

	_, err = usersCollection.InsertOne(ctx, newUser)
	if err != nil {
		c.JSON(500, gin.H{"error": "Could not create user"})
		return
	}

	logMessage(fmt.Sprintf("User registered: %s", newUser.Username))

	c.JSON(201, gin.H{"message": "User created"})
}

func login(c *gin.Context) {
	var loginUser, foundUser User
	if err := c.BindJSON(&loginUser); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	err := usersCollection.FindOne(ctx, bson.M{"username": loginUser.Username}).Decode(&foundUser)
	if err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	if !checkPasswordHash(loginUser.Password, foundUser.Password) {
		c.JSON(401, gin.H{"error": "Invalid password"})
		return
	}

	logMessage(fmt.Sprintf("User logged in: %s", loginUser.Username))

	c.JSON(200, gin.H{"message": "Login successful"})
}

func main() {
	router := gin.Default()

	router.POST("/register", register)
	router.POST("/login", login)

	err := router.Run(":8080")
	if err != nil {
		fmt.Println("Failed to start server: ", err)
	}
}
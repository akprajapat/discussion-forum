package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"discussion-forum/handlers"

	"github.com/joho/godotenv"
)

var DB *mongo.Database

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

}

func ConnectDB() *mongo.Client {
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI"))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Ensure the context is cancelled when the function exits

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Ping the database to verify the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")
	return client
}

func main() {
	client := ConnectDB()
	DB = client.Database(os.Getenv("DATABASE"))

	// Ensure unique indexes for username and email
	userCol := DB.Collection("users")
	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "username", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}
	if _, err := userCol.Indexes().CreateMany(context.Background(), indexModels); err != nil {
		log.Fatal("Failed to create indexes:", err)
	}

	r := gin.Default()
	r.Use(cors.Default())

	handlers.Init(DB)
	// Initialize routes
	setupRoutes(r)

	// Start server
	r.Run(":8080")
}

func setupRoutes(r *gin.Engine) {
	api := r.Group("/api")

	// Auth routes
	api.POST("/register", handlers.Register)
	api.POST("/login", handlers.Login)

	// Question routes
	api.GET("/questions", handlers.GetQuestions)
	api.POST("/questions", handlers.AuthMiddleware(), handlers.CreateQuestion)
	api.GET("/questions/:id", handlers.GetQuestion)
	api.PUT("/questions/:id/vote", handlers.AuthMiddleware(), handlers.VoteQuestion)

	// Answer routes
	api.POST("/questions/:id/answers", handlers.AuthMiddleware(), handlers.CreateAnswer)
	api.PUT("/answers/:id/vote", handlers.AuthMiddleware(), handlers.VoteAnswer)

	// Comment routes
	api.POST("/answers/:id/comments", handlers.AuthMiddleware(), handlers.CreateComment)
	api.GET("/answers/:id/comments", handlers.GetComments)
}

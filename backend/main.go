package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"discussion-forum/handlers"
)

var DB *mongo.Database

func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.Connect(ctx); err != nil {
		log.Fatal(err)
	}
	DB = client.Database("discussion_forum")

	r := gin.Default()

	// Enable CORS
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
}

package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"tmdt-backend/bootstrap"
	"tmdt-backend/route"
)

func main() {
	// Load .env file from local directory or parent directory if it exists
	if err := godotenv.Load(); err != nil {
		if err := godotenv.Load("../.env"); err != nil {
			log.Println("Warning: No .env file found, using system environment variables")
		}
	}

	ctx := context.Background()

	// Connect PostgreSQL
	dbPool, err := bootstrap.ConnectDB(ctx)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer dbPool.Close()

	// Connect Redis
	redisClient, err := bootstrap.ConnectRedis(ctx)
	if err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	}
	defer redisClient.Close()

	// Setup Gin
	r := gin.New()

	// Setup Routes
	route.Setup(dbPool, redisClient, r)

	// Determine port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server is running on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

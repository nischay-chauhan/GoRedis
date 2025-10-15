package main

import (
	"context"
	"log"
	"net/http"

	"go-redis/internal/config"
	"go-redis/internal/routes"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func main() {
	cfg := config.Load()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: "",
		DB:       0,
	})

	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Connected to Redis")

	router := routes.SetupRoutes(redisClient)

	serverAddr := ":" + cfg.Port
	log.Printf("Server starting on port %s\n", cfg.Port)
	log.Printf("Available endpoints:")
	log.Printf("  POST http://localhost:%s/score - Submit a score (requires JSON body: {\"player\": \"name\", \"score\": 100})", cfg.Port)
	log.Printf("  GET  http://localhost:%s/health - Health check\n", cfg.Port)

	if err := http.ListenAndServe(serverAddr, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

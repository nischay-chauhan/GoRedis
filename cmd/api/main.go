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

	rgb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: "", 
		DB:       0,  
	})

	_, err := rgb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Connected to Redis")

	router := routes.SetupRoutes()

	serverAddr := ":" + cfg.Port
	log.Printf("Server starting on port %s\n", cfg.Port)
	log.Printf("Available endpoints:")
	log.Printf("  POST http://localhost:%s/score", cfg.Port)
	log.Printf("  GET  http://localhost:%s/health\n", cfg.Port)

	if err := http.ListenAndServe(serverAddr, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func main() {
	rgb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := rgb.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Failed to connect to Redis", err)
		return
	}
	fmt.Println("Connected to Redis")

	http.HandleFunc("/health", healthCheck)
	fmt.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

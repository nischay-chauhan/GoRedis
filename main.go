package main

import (
    "fmt"
	"context"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()	

func main() {
	rgb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	})

	_,err := rgb.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Failed to connect to Redis", err)
		return
	}

	fmt.Println("Connected to Redis")
	fmt.Println("Game leaderboard API")
}
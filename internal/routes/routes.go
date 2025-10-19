package routes

import (
    "net/http"
    "time"
    "go-redis/internal/handlers"
    "go-redis/internal/middleware"
    "github.com/redis/go-redis/v9"
)

func SetupRoutes(redisClient *redis.Client) *http.ServeMux {
    mux := http.NewServeMux()
    scoreHandler := handlers.NewScoreHandler(redisClient)
    leaderboardHandler := handlers.NewLeaderboardHandler(redisClient)

    rateLimiter := middleware.NewRateLimiter(60, 10, 1*time.Minute)

    mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    })

    mux.Handle("POST /score", rateLimiter.Limit(http.HandlerFunc(scoreHandler.SubmitScore)))
    mux.Handle("GET /score", rateLimiter.Limit(http.HandlerFunc(scoreHandler.GetScore)))

    mux.Handle("GET /leaderboard/top", rateLimiter.Limit(http.HandlerFunc(leaderboardHandler.Top)))
    mux.Handle("GET /leaderboard/player", rateLimiter.Limit(http.HandlerFunc(leaderboardHandler.Player)))
    mux.Handle("GET /leaderboard/around/{player}", rateLimiter.Limit(http.HandlerFunc(leaderboardHandler.Around)))

    return mux
}

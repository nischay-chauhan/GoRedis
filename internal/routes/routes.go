package routes

import (
    "net/http"
    "time"
    "go-redis/internal/handlers"
    "go-redis/internal/middleware"
    "github.com/redis/go-redis/v9"
)

const (
    defaultTimeout = 30 * time.Second
)

func SetupRoutes(redisClient *redis.Client) *http.ServeMux {
    mux := http.NewServeMux()
    scoreHandler := handlers.NewScoreHandler(redisClient)
    leaderboardHandler := handlers.NewLeaderboardHandler(redisClient)

    rateLimiter := middleware.NewRateLimiter(60, 10, 1*time.Minute)
    
    cors := middleware.NewCors(&middleware.CorsConfig{
        AllowedOrigins: []string{"*"},
        AllowedMethods: []string{"GET", "POST", "OPTIONS"},
        AllowedHeaders: []string{"Content-Type", "Authorization"},
    })
    
    timeout := middleware.NewTimeout(middleware.TimeoutConfig{
        DefaultTimeout: defaultTimeout,
    })

    mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    })

    mux.Handle("POST /score", cors(timeout(rateLimiter.Limit(http.HandlerFunc(scoreHandler.SubmitScore)))))
    mux.Handle("GET /score", cors(timeout(rateLimiter.Limit(http.HandlerFunc(scoreHandler.GetScore)))))

    mux.Handle("GET /leaderboard/top", cors(timeout(rateLimiter.Limit(http.HandlerFunc(leaderboardHandler.Top)))))
    mux.Handle("GET /leaderboard/player", cors(timeout(rateLimiter.Limit(http.HandlerFunc(leaderboardHandler.Player)))))
    mux.Handle("GET /leaderboard/around/{player}", cors(timeout(rateLimiter.Limit(http.HandlerFunc(leaderboardHandler.Around)))))

    return mux
}

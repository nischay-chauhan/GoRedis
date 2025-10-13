package routes

import (
	"net/http"
	"go-redis/internal/handlers"
	"github.com/redis/go-redis/v9"
)

func SetupRoutes(redisClient *redis.Client) *http.ServeMux {
	mux := http.NewServeMux()
	scoreHandler := handlers.NewScoreHandler(redisClient)

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("POST /score", scoreHandler.SubmitScore)
	mux.HandleFunc("GET /score", scoreHandler.GetScore)

	return mux
}

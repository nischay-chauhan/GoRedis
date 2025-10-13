package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/redis/go-redis/v9"
)

const (
	scoreSet = "scores"
)

type ScoreRequest struct {
	Player string `json:"player"`
	Score  int    `json:"score"`
}

type ScoreHandler struct {
	redisClient *redis.Client
}

func NewScoreHandler(redisClient *redis.Client) *ScoreHandler {
	return &ScoreHandler{
		redisClient: redisClient,
	}
}

func (h *ScoreHandler) SubmitScore(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ScoreRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Player == "" {
		http.Error(w, "Player name is required", http.StatusBadRequest)
		return
	}

	if req.Score <= 0 {
		http.Error(w, "Score must be a positive number", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	score, err := h.redisClient.ZIncrBy(ctx, scoreSet, float64(req.Score), req.Player).Result()
	if err != nil {
		log.Printf("Failed to update score in Redis: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Score updated successfully",
		"player":  req.Player,
		"score":   score,
	})
}

func (h *ScoreHandler) GetScore(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	player := r.URL.Query().Get("player")
	if player == "" {
		http.Error(w, "Player name is required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	score, err := h.redisClient.ZScore(ctx, scoreSet, player).Result()
	if err != nil {
		if err == redis.Nil {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status":  "success",
				"message": "Player not found",
				"player":  player,
				"score":   0,
			})
			return
		}
		log.Printf("Failed to get score from Redis: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Score retrieved successfully",
		"player":  player,
		"score":   score,
	})
}

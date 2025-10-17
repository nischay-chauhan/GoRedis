package handlers

import (
	"encoding/json"
	"go-redis/internal/models"
	"net/http"
	"strconv"

	"github.com/redis/go-redis/v9"
)

const leaderboardSet = "scores"

type LeaderboardHandler struct {
	redisClient *redis.Client
}

func NewLeaderboardHandler(redisClient *redis.Client) *LeaderboardHandler {
	return &LeaderboardHandler{redisClient: redisClient}
}

// Top handles GET /leaderboard/top?limit=10
func (h *LeaderboardHandler) Top(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	limit := 10
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			limit = n
		}
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	ctx := r.Context()
	stop := int64(limit - 1)
	zs, err := h.redisClient.ZRevRangeWithScores(ctx, leaderboardSet, 0, stop).Result()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	entries := make([]models.LeaderboardEntry, 0, len(zs))
	for i, z := range zs {
		memberStr, _ := z.Member.(string)
		entries = append(entries, models.LeaderboardEntry{
			Rank:   i + 1,
			Player: memberStr,
			Score:  z.Score,
		})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Top players retrieved successfully",
		"limit":   limit,
		"data":    entries,
	})
}

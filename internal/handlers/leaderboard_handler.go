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

// Player handles GET /leaderboard/player?player=:id
// Returns rank (1-based), score, and percentile (0-100 where higher is better)
func (h *LeaderboardHandler) Player(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    player := r.URL.Query().Get("player")
    if player == "" {
        http.Error(w, "Player is required", http.StatusBadRequest)
        return
    }

    ctx := r.Context()

    rank0, err := h.redisClient.ZRevRank(ctx, leaderboardSet, player).Result()
    if err != nil {
        if err == redis.Nil {
            w.WriteHeader(http.StatusOK)
            json.NewEncoder(w).Encode(map[string]interface{}{
                "status":  "success",
                "message": "Player not found",
                "player":  player,
                "rank":    nil,
                "score":   0,
                "percentile": 0,
            })
            return
        }
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    pipe := h.redisClient.Pipeline()
    scoreCmd := pipe.ZScore(ctx, leaderboardSet, player)
    totalCmd := pipe.ZCard(ctx, leaderboardSet)
    if _, err := pipe.Exec(ctx); err != nil && err != redis.Nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    score, _ := scoreCmd.Result()
    total, _ := totalCmd.Result()

    rank := int(rank0) + 1
    var percentile float64
    if total > 0 {
        percentile = (1 - (float64(rank) / float64(total))) * 100.0
        if percentile < 0 {
            percentile = 0
        }
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "status":     "success",
        "message":    "Player rank retrieved successfully",
        "player":     player,
        "rank":       rank,
        "score":      score,
        "total":      total,
        "percentile": percentile,
    })
}

// Around handles GET /leaderboard/around/{player}?radius=2
// Returns entries around the player's current rank.
func (h *LeaderboardHandler) Around(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    player := r.PathValue("player")
    if player == "" {
        http.Error(w, "Player is required", http.StatusBadRequest)
        return
    }

    radius := 2
    if v := r.URL.Query().Get("radius"); v != "" {
        if n, err := strconv.Atoi(v); err == nil {
            radius = n
        }
    }
    if radius < 1 {
        radius = 1
    }
    if radius > 10 {
        radius = 10
    }

    ctx := r.Context()
    rank0, err := h.redisClient.ZRevRank(ctx, leaderboardSet, player).Result()
    if err != nil {
        if err == redis.Nil {
            w.WriteHeader(http.StatusOK)
            json.NewEncoder(w).Encode(map[string]interface{}{
                "status":  "success",
                "message": "Player not found",
                "player":  player,
                "data":    []models.LeaderboardEntry{},
            })
            return
        }
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    start := int64(0)
    if int64(rank0)-int64(radius) > 0 {
        start = int64(rank0) - int64(radius)
    }
    end := int64(rank0) + int64(radius)

    zs, err := h.redisClient.ZRevRangeWithScores(ctx, leaderboardSet, start, end).Result()
    if err != nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    entries := make([]models.LeaderboardEntry, 0, len(zs))
    for i, z := range zs {
        memberStr, _ := z.Member.(string)
        entries = append(entries, models.LeaderboardEntry{
            Rank:   int(start) + i + 1,
            Player: memberStr,
            Score:  z.Score,
        })
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "status":  "success",
        "message": "Around player window retrieved successfully",
        "player":  player,
        "radius":  radius,
        "data":    entries,
    })
}

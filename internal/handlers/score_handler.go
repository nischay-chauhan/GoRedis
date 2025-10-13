package handlers

import (
	"encoding/json"
	"net/http"
)

type ScoreRequest struct {
	Player string `json:"player"`
	Score  int    `json:"score"`
}

func SubmitScore(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ScoreRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Player == "" {
		http.Error(w, "Player name is required", http.StatusBadRequest)
		return
	}

	if req.Score < 0 {
		http.Error(w, "Score must be a positive number", http.StatusBadRequest)
		return
	}

	// TODO: Add to Redis

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Score submitted successfully",
		"player":  req.Player,
		"score":   string(rune(req.Score)),
	})
}

package models

type ScoreRequest struct {
	Player string `json:"player"`
	Score  int    `json:"score"`
}

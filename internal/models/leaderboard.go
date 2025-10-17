package models

type ScoreRequest struct {
    Player string `json:"player"`
    Score  int    `json:"score"`
}

type LeaderboardEntry struct {
    Rank   int     `json:"rank"`
    Player string  `json:"player"`
    Score  float64 `json:"score"`
}

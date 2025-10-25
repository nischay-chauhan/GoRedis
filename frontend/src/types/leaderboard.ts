export interface Score {
  player: string;
  score: number;
  timestamp?: string;
}

export interface PlayerScore {
  player: string;
  score: number;
  rank: number;
}

export interface LeaderboardResponse {
  data: PlayerScore[];
  total: number;
}

export interface SubmitScoreRequest {
  player: string;
  score: number;
}

export interface PlayerScoresResponse {
  player: string;
  scores: {
    score: number;
    timestamp: string;
  }[];
  rank: number;
  total: number;
}

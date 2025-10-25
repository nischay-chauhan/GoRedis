import axios from 'axios';
import { SubmitScoreRequest, LeaderboardResponse, PlayerScoresResponse } from '@/types/leaderboard';

const API_BASE_URL = '/api';

export const api = {
  // Submit a new score
  async submitScore(data: SubmitScoreRequest): Promise<void> {
    await axios.post(`${API_BASE_URL}/score`, data);
  },

  // Get top players
  async getTopPlayers(limit: number = 10): Promise<LeaderboardResponse> {
    const response = await axios.get(`${API_BASE_URL}/leaderboard/top?limit=${limit}`);
    return response.data;
  },

  // Get player's scores
  async getPlayerScores(player: string): Promise<PlayerScoresResponse> {
    const response = await axios.get(`${API_BASE_URL}/leaderboard/player?player=${encodeURIComponent(player)}`);
    return response.data;
  },

  // Get players around a specific player
  async getPlayersAround(player: string, limit: number = 5): Promise<LeaderboardResponse> {
    const response = await axios.get(
      `${API_BASE_URL}/leaderboard/around/${encodeURIComponent(player)}?limit=${limit}`
    );
    return response.data;
  },

  // Health check
  async healthCheck(): Promise<{ status: string }> {
    const response = await axios.get(`${API_BASE_URL}/health`);
    return response.data;
  },
};

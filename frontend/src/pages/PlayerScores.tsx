import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { useSearchParams } from 'react-router-dom';
import { api } from '@/services/api';

interface PlayerScore {
  score: number;
  timestamp: string;
  // Add other properties that might be in your score objects
}

interface PlayerData {
  player: string;
  rank?: number;
  score?: number;
  scores?: PlayerScore[];
  total?: number;
  // Add other properties that might be in your player data
}

export function PlayerScores() {
  const [searchParams, setSearchParams] = useSearchParams();
  const [playerName, setPlayerName] = useState(searchParams.get('player') || '');
  const [inputValue, setInputValue] = useState(searchParams.get('player') || '');

  const { data, isLoading, isError } = useQuery<PlayerData | null>({
    queryKey: ['playerScores', playerName],
    queryFn: () => api.getPlayerScores(playerName),
    enabled: !!playerName,
    retry: false,
  });

  // Safely get scores array or default to empty array
  const scores = data?.scores || [];

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const trimmedName = inputValue.trim();
    if (trimmedName) {
      setPlayerName(trimmedName);
      setSearchParams({ player: trimmedName });
    }
  };

  return (
    <div>
      <div className="mb-8">
        <h1 className="text-2xl font-bold text-gray-900 mb-2">Player Stats</h1>
        <p className="text-gray-600">
          View your scores and ranking on the leaderboard.
        </p>
      </div>

      <div className="bg-white shadow rounded-lg p-6 mb-8">
        <form onSubmit={handleSubmit} className="flex flex-col sm:flex-row gap-4">
          <div className="flex-grow">
            <label htmlFor="player" className="sr-only">
              Player Name
            </label>
            <input
              type="text"
              id="player"
              value={inputValue}
              onChange={(e) => setInputValue(e.target.value)}
              className="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
              placeholder="Enter player name"
              required
            />
          </div>
          <button
            type="submit"
            className="inline-flex items-center justify-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
          >
            View Stats
          </button>
        </form>
      </div>

      {isLoading && playerName && (
        <div className="flex justify-center items-center h-64">
          <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-primary-600"></div>
        </div>
      )}

      {isError && (
        <div className="bg-red-50 border-l-4 border-red-500 p-4 mb-6">
          <div className="flex">
            <div className="flex-shrink-0">
              <svg className="h-5 w-5 text-red-500" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
              </svg>
            </div>
            <div className="ml-3">
              <p className="text-sm text-red-700">
                Player not found or no scores available.
              </p>
            </div>
          </div>
        </div>
      )}

      {data && data.player && (
        <div className="space-y-8">
          <div className="bg-white overflow-hidden shadow rounded-lg">
            <div className="px-4 py-5 sm:p-6">
              <div className="flex items-center">
                <div className="flex-shrink-0 bg-primary-500 rounded-md p-3">
                  <svg className="h-6 w-6 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                  </svg>
                </div>
                <div className="ml-5 w-0 flex-1">
                  <dl>
                    <dt className="text-sm font-medium text-gray-500 truncate">Player</dt>
                    <dd className="flex items-baseline">
                      <div className="text-2xl font-semibold text-gray-900">{data.player}</div>
                    </dd>
                  </dl>
                </div>
                <div className="ml-5 w-0 flex-1">
                  <dl>
                    <dt className="text-sm font-medium text-gray-500 truncate">Rank</dt>
                    <dd className="flex items-baseline">
                      <div className="text-2xl font-semibold text-gray-900">#{data.rank}</div>
                    </dd>
                  </dl>
                </div>
                <div className="ml-5 w-0 flex-1">
                  <dl>
                    <dt className="text-sm font-medium text-gray-500 truncate">Total Players</dt>
                    <dd className="flex items-baseline">
                      <div className="text-2xl font-semibold text-gray-900">{data.total}</div>
                    </dd>
                  </dl>
                </div>
              </div>
            </div>
          </div>

          <div className="bg-white shadow overflow-hidden sm:rounded-lg p-6">
            <div className="mb-4">
              <h3 className="text-lg leading-6 font-medium text-gray-900">Score History</h3>
              <p className="mt-1 max-w-2xl text-sm text-gray-500">
                {data?.player ? `Recent scores submitted by ${data.player}` : 'No player found'}
              </p>
            </div>
            <div className="bg-white overflow-hidden">
              <ul className="divide-y divide-gray-200">
                {scores.length > 0 ? (
                  scores.map((score: PlayerScore, index: number) => (
                    <li key={index} className="px-4 py-4 sm:px-6">
                      <div className="flex items-center justify-between">
                        <div className="text-sm font-medium text-gray-900">
                          {score?.timestamp ? new Date(score.timestamp).toLocaleString() : 'No timestamp found'}
                        </div>
                        <span className="px-2 inline-flex text-sm leading-5 font-semibold rounded-full bg-green-100 text-green-800">
                          {score?.score?.toLocaleString() || 'N/A'}
                        </span>
                      </div>
                    </li>
                  ))
                ) : (
                  <li className="px-4 py-4 sm:px-6 text-center text-gray-500">
                    No scores found for this player.
                  </li>
                )}
              </ul>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

import { useQuery } from '@tanstack/react-query';
import { api } from '@/services/api';
import { useEffect } from 'react';
import { LeaderboardResponse } from '@/types/leaderboard';

interface Player {
  rank: number;
  player: string;
  score: number;
}

export function Leaderboard() {
  console.log(' Leaderboard component rendered');

  const { data, isLoading, error, isFetching } = useQuery<LeaderboardResponse>({
    queryKey: ['leaderboard'],
    queryFn: async () => {
      console.log(' Fetching leaderboard data...');
      const response = await api.getTopPlayers(10);
      console.log(' API Response:', response);
      return response;
    },
  });

  const players = data?.data || [];

  useEffect(() => {
    console.log('🔄 Component mounted or data updated', {
      isLoading,
      isFetching,
      error,
      hasData: !!data,
      playerCount: players.length,
    });
  }, [isLoading, isFetching, error, data, players.length]);

  if (isLoading) {
    console.log('⏳ Rendering loading state');
    return (
      <div className="flex justify-center items-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-primary-600"></div>
      </div>
    );
  }

  if (error) {
    console.error(' Rendering error state:', error);
    return <div className="text-red-600 p-4">Error loading leaderboard</div>;
  }

  console.log(' Rendering leaderboard with', players.length, 'players');

  return (
    <div className="max-w-4xl mx-auto px-4 py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900">Leaderboard</h1>
        <p className="mt-2 text-sm text-gray-600">Top players by score</p>
      </div>

      {players.length === 0 ? (
        <div className="text-center py-12">
          <svg
            className="mx-auto h-12 w-12 text-gray-400"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={1}
              d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
            />
          </svg>
          <h3 className="mt-2 text-sm font-medium text-gray-900">No players yet</h3>
          <p className="mt-1 text-sm text-gray-500">Be the first to submit a score!</p>
        </div>
      ) : (
        <div className="bg-white shadow overflow-hidden sm:rounded-md">
          <ul className="divide-y divide-gray-200">
            {players.map((player: Player, index: number) => (
              <li key={`${player.player}-${index}`}>
                <div className="px-4 py-4 sm:px-6">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center">
                      <div className="flex-shrink-0 h-10 w-10 flex items-center justify-center rounded-full bg-blue-100 text-blue-700 font-bold">
                        {index + 1}
                      </div>
                      <div className="ml-4">
                        <div className="text-sm font-medium text-gray-900">{player.player}</div>
                        <div className="text-sm text-gray-500">Score: {player.score.toLocaleString()}</div>
                      </div>
                    </div>
                    <div className="ml-2 flex-shrink-0">
                      <div className="px-2 inline-flex text-sm leading-5 font-semibold rounded-full bg-green-100 text-green-800">
                        {player.score.toLocaleString()}
                      </div>
                    </div>
                  </div>
                </div>
              </li>
            ))}
          </ul>
        </div>
      )}

      <div className="mt-4 text-sm text-gray-500 text-right">
        Showing {players.length} players
      </div>
    </div>
  );
}

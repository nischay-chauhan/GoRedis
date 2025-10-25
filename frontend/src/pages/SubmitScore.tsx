import { useState } from 'react';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '@/services/api';

export function SubmitScore() {
  const [formData, setFormData] = useState({
    player: '',
    score: '',
  });
  const [isSubmitted, setIsSubmitted] = useState(false);
  const queryClient = useQueryClient();

  const submitScoreMutation = useMutation({
    mutationFn: () => 
      api.submitScore({
        player: formData.player.trim(),
        score: Number(formData.score),
      }),
    onSuccess: () => {
      setIsSubmitted(true);
      setFormData({ player: '', score: '' });
      // Invalidate and refetch leaderboard data
      queryClient.invalidateQueries({ queryKey: ['leaderboard'] });
      
      // Reset the success message after 5 seconds
      setTimeout(() => setIsSubmitted(false), 5000);
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (formData.player.trim() && formData.score) {
      submitScoreMutation.mutate();
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: name === 'score' ? value.replace(/\D/g, '') : value,
    }));
  };

  return (
    <div className="max-w-md mx-auto">
      <div className="text-center mb-8">
        <h1 className="text-2xl font-bold text-gray-900">Submit Your Score</h1>
        <p className="mt-2 text-sm text-gray-600">
          Enter your name and score to be added to the leaderboard.
        </p>
      </div>

      {isSubmitted && (
        <div className="mb-6 p-4 bg-green-50 rounded-md">
          <div className="flex">
            <div className="flex-shrink-0">
              <svg className="h-5 w-5 text-green-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
              </svg>
            </div>
            <div className="ml-3">
              <p className="text-sm font-medium text-green-800">
                Score submitted successfully!
              </p>
            </div>
          </div>
        </div>
      )}

      {submitScoreMutation.isError && (
        <div className="mb-6 p-4 bg-red-50 rounded-md">
          <div className="flex">
            <div className="flex-shrink-0">
              <svg className="h-5 w-5 text-red-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
              </svg>
            </div>
            <div className="ml-3">
              <p className="text-sm font-medium text-red-800">
                Failed to submit score. Please try again.
              </p>
            </div>
          </div>
        </div>
      )}

      <div className="bg-white py-8 px-6 shadow rounded-lg">
        <form onSubmit={handleSubmit} className="space-y-6">
          <div>
            <label htmlFor="player" className="block text-sm font-medium text-gray-700">
              Player Name
            </label>
            <div className="mt-1">
              <input
                type="text"
                name="player"
                id="player"
                required
                value={formData.player}
                onChange={handleChange}
                maxLength={50}
                className="input"
                placeholder="Enter your name"
              />
            </div>
          </div>

          <div>
            <label htmlFor="score" className="block text-sm font-medium text-gray-700">
              Score
            </label>
            <div className="mt-1">
              <input
                type="text"
                name="score"
                id="score"
                required
                value={formData.score}
                onChange={handleChange}
                pattern="\d*"
                className="input"
                placeholder="Enter your score"
              />
            </div>
          </div>

          <div>
            <button
              type="submit"
              disabled={submitScoreMutation.isPending || !formData.player.trim() || !formData.score}
              className="w-full flex justify-center py-3 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {submitScoreMutation.isPending ? (
                <>
                  <svg className="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                    <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  Submitting...
                </>
              ) : (
                'Submit Score'
              )}
            </button>
          </div>
        </form>
      </div>

      <div className="mt-8 text-center">
        <p className="text-sm text-gray-500">
          Your score will be added to the leaderboard after submission.
        </p>
      </div>
    </div>
  );
}

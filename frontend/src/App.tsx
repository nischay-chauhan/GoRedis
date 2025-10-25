import { Routes, Route, Link } from 'react-router-dom';
import { Home } from '@/pages/Home';
import { Leaderboard } from '@/pages/Leaderboard';
import { SubmitScore } from '@/pages/SubmitScore';
import { PlayerScores } from '@/pages/PlayerScores';

export function App() {
  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
          <div className="flex justify-between items-center">
            <Link to="/" className="text-2xl font-bold text-primary-600">
              Redis Leaderboard
            </Link>
            <nav className="flex space-x-4">
              <Link to="/" className="text-gray-700 hover:text-primary-600">
                Home
              </Link>
              <Link to="/submit" className="text-gray-700 hover:text-primary-600">
                Submit Score
              </Link>
              <Link to="/leaderboard" className="text-gray-700 hover:text-primary-600">
                Leaderboard
              </Link>
            </nav>
          </div>
        </div>
      </header>

      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/submit" element={<SubmitScore />} />
          <Route path="/leaderboard" element={<Leaderboard />} />
          <Route path="/player" element={<PlayerScores />} />
        </Routes>
      </main>

      <footer className="bg-white border-t border-gray-200 mt-12">
        <div className="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
          <p className="text-center text-gray-500 text-sm">
            &copy; {new Date().getFullYear()} Redis Leaderboard Demo
          </p>
        </div>
      </footer>
    </div>
  );
}

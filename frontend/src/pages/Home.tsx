import { Link } from 'react-router-dom';

export function Home() {
  return (
    <div className="text-center">
      <h1 className="text-4xl font-bold text-gray-900 mb-6">Redis Leaderboard Demo</h1>
      <p className="text-xl text-gray-600 mb-8">
        A high-performance leaderboard system built with Go, Redis, and React
      </p>
      
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mt-12">
        <div className="card hover:shadow-lg transition-shadow">
          <div className="text-primary-600 mb-4">
            <svg xmlns="http://www.w3.org/2000/svg" className="h-12 w-12 mx-auto" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
            </svg>
          </div>
          <h3 className="text-xl font-semibold mb-2">View Leaderboard</h3>
          <p className="text-gray-600 mb-4">See the top players and their scores in real-time.</p>
          <Link to="/leaderboard" className="btn btn-primary inline-block">
            View Leaderboard
          </Link>
        </div>

        <div className="card hover:shadow-lg transition-shadow">
          <div className="text-primary-600 mb-4">
            <svg xmlns="http://www.w3.org/2000/svg" className="h-12 w-12 mx-auto" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
            </svg>
          </div>
          <h3 className="text-xl font-semibold mb-2">Submit Score</h3>
          <p className="text-gray-600 mb-4">Add your score to the leaderboard and see where you rank.</p>
          <Link to="/submit" className="btn btn-primary inline-block">
            Submit Your Score
          </Link>
        </div>

        <div className="card hover:shadow-lg transition-shadow">
          <div className="text-primary-600 mb-4">
            <svg xmlns="http://www.w3.org/2000/svg" className="h-12 w-12 mx-auto" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
            </svg>
          </div>
          <h3 className="text-xl font-semibold mb-2">Your Stats</h3>
          <p className="text-gray-600 mb-4">Check your ranking and score history.</p>
          <Link to="/player" className="btn btn-primary inline-block">
            View Your Stats
          </Link>
        </div>
      </div>

      <div className="mt-16 bg-white p-6 rounded-lg shadow">
        <h2 className="text-2xl font-bold mb-4">About This Demo</h2>
        <p className="text-gray-700 mb-4">
          This application showcases a high-performance leaderboard system built with:
        </p>
        <ul className="grid grid-cols-1 md:grid-cols-3 gap-4 text-left">
          <li className="flex items-center">
            <span className="bg-blue-100 text-blue-800 text-xs font-medium mr-2 px-2.5 py-0.5 rounded">Backend</span>
            <span>Go with Redis for high-performance data storage</span>
          </li>
          <li className="flex items-center">
            <span className="bg-green-100 text-green-800 text-xs font-medium mr-2 px-2.5 py-0.5 rounded">Frontend</span>
            <span>React with TypeScript and Tailwind CSS</span>
          </li>
          <li className="flex items-center">
            <span className="bg-yellow-100 text-yellow-800 text-xs font-medium mr-2 px-2.5 py-0.5 rounded">Features</span>
            <span>Real-time updates, caching, and rate limiting</span>
          </li>
        </ul>
      </div>
    </div>
  );
}

import React, { useEffect, useState } from 'react';
import axios from 'axios';

interface LeaderboardEntry {
  rank: number;
  user: {
    id: string;
    username: string;
    display_name: string;
    avatar_url: string;
  };
  total_co2_saved_kg: number;
  sustainability_score: number;
  level: number;
}

interface LeaderboardResponse {
  leaderboard: LeaderboardEntry[];
}

export const LeaderboardSidebar: React.FC = () => {
  const [leaderboard, setLeaderboard] = useState<LeaderboardEntry[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchLeaderboard = async () => {
      try {
        const response = await axios.get<LeaderboardResponse>(
          'http://localhost:8080/v1/sustainability/leaderboard'
        );
        setLeaderboard(response.data.leaderboard);
        setLoading(false);
      } catch (err) {
        console.error('Failed to fetch leaderboard:', err);
        setError('リーダーボードの読み込みに失敗しました');
        setLoading(false);
      }
    };

    fetchLeaderboard();
  }, []);

  if (loading) {
    return (
      <div className="w-80 bg-white/80 backdrop-blur-sm rounded-lg shadow-lg p-6">
        <div className="animate-pulse">
          <div className="h-8 bg-gray-200 rounded w-3/4 mb-4"></div>
          {[1, 2, 3, 4, 5].map((i) => (
            <div key={i} className="flex items-center gap-3 mb-3">
              <div className="w-12 h-12 bg-gray-200 rounded-full"></div>
              <div className="flex-1">
                <div className="h-4 bg-gray-200 rounded w-3/4 mb-2"></div>
                <div className="h-3 bg-gray-200 rounded w-1/2"></div>
              </div>
            </div>
          ))}
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="w-80 bg-white/80 backdrop-blur-sm rounded-lg shadow-lg p-6">
        <p className="text-red-600 text-center">{error}</p>
      </div>
    );
  }

  const getMedalEmoji = (rank: number) => {
    switch (rank) {
      case 1:
        return '🥇';
      case 2:
        return '🥈';
      case 3:
        return '🥉';
      default:
        return `${rank}位`;
    }
  };

  return (
    <div className="w-80 bg-white/80 backdrop-blur-sm rounded-lg shadow-lg p-6 sticky top-4">
      <div className="flex items-center gap-2 mb-6">
        <span className="text-3xl">🏆</span>
        <h2 className="text-2xl font-bold text-gray-900">エコランキング</h2>
      </div>

      <div className="space-y-3">
        {leaderboard.slice(0, 10).map((entry) => (
          <div
            key={`${entry.user.id}-${entry.rank}`}
            className={`flex items-center gap-3 p-3 rounded-lg transition-all hover:shadow-md ${
              entry.rank <= 3
                ? 'bg-gradient-to-r from-yellow-50 to-green-50 border border-yellow-200'
                : 'bg-gray-50 hover:bg-gray-100'
            }`}
          >
            <div className="flex-shrink-0 text-2xl font-bold w-12 text-center">
              {getMedalEmoji(entry.rank)}
            </div>

            <img
              src={entry.user.avatar_url}
              alt={entry.user.username}
              className="w-12 h-12 rounded-full object-cover border-2 border-green-200"
            />

            <div className="flex-1 min-w-0">
              <p className="font-semibold text-gray-900 truncate">
                {entry.user.display_name || entry.user.username}
              </p>
              <div className="flex items-center gap-2 text-sm">
                <span className="text-green-600 font-bold">
                  {entry.total_co2_saved_kg.toFixed(1)}kg
                </span>
                <span className="text-gray-500">CO₂削減</span>
              </div>
            </div>

            <div className="flex flex-col items-end">
              <span className="text-xs font-medium text-white bg-green-600 px-2 py-1 rounded-full">
                Lv.{entry.level}
              </span>
            </div>
          </div>
        ))}
      </div>

      <div className="mt-6 p-4 bg-green-50 rounded-lg border border-green-200">
        <div className="flex items-center gap-2 mb-2">
          <span className="text-xl">🌱</span>
          <p className="text-sm font-semibold text-gray-900">
            中古品購入でCO₂削減に貢献しよう！
          </p>
        </div>
        <p className="text-xs text-gray-600">
          新品購入を避けることで地球環境を守り、ランキング上位を目指そう。
        </p>
      </div>
    </div>
  );
};

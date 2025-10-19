import React, { useState, useEffect } from 'react';
import axios from 'axios';

interface AIAgentStats {
  total_ai_generations: number;
  listings_created: number;
  negotiations_handled: number;
  shipments_prepared: number;
  average_confidence: number;
  time_saved_minutes: number;
  acceptance_rate: number;
}

const AIAgentDashboard: React.FC = () => {
  const [stats, setStats] = useState<AIAgentStats | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchStats();
  }, []);

  const fetchStats = async () => {
    try {
      const response = await axios.get('/api/v1/ai-agent/stats', {
        headers: { Authorization: `Bearer ${localStorage.getItem('token')}` },
      });
      setStats(response.data);
    } catch (error) {
      console.error('Failed to fetch stats:', error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-4 border-blue-600"></div>
      </div>
    );
  }

  const timeSavedHours = stats ? (stats.time_saved_minutes / 60).toFixed(1) : '0';

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-7xl mx-auto px-4">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900 mb-2">
            🤖 AIエージェント ダッシュボード
          </h1>
          <p className="text-gray-600">
            AIがあなたの代わりに処理したタスクの統計情報
          </p>
        </div>

        {/* Key Metrics */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
          {/* Time Saved */}
          <div className="bg-gradient-to-br from-blue-500 to-blue-600 rounded-lg shadow-lg p-6 text-white">
            <div className="flex items-center justify-between mb-2">
              <h3 className="text-sm font-medium opacity-90">節約時間</h3>
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <div className="text-4xl font-bold mb-1">{timeSavedHours}</div>
            <div className="text-sm opacity-90">時間</div>
          </div>

          {/* Listings Created */}
          <div className="bg-gradient-to-br from-green-500 to-green-600 rounded-lg shadow-lg p-6 text-white">
            <div className="flex items-center justify-between mb-2">
              <h3 className="text-sm font-medium opacity-90">AI出品数</h3>
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
            </div>
            <div className="text-4xl font-bold mb-1">{stats?.listings_created || 0}</div>
            <div className="text-sm opacity-90">件</div>
          </div>

          {/* Negotiations */}
          <div className="bg-gradient-to-br from-purple-500 to-purple-600 rounded-lg shadow-lg p-6 text-white">
            <div className="flex items-center justify-between mb-2">
              <h3 className="text-sm font-medium opacity-90">AI交渉処理</h3>
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 8h2a2 2 0 012 2v6a2 2 0 01-2 2h-2v4l-4-4H9a1.994 1.994 0 01-1.414-.586m0 0L11 14h4a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2v4l.586-.586z" />
              </svg>
            </div>
            <div className="text-4xl font-bold mb-1">{stats?.negotiations_handled || 0}</div>
            <div className="text-sm opacity-90">件 ({stats?.acceptance_rate.toFixed(1)}% 承認)</div>
          </div>

          {/* Shipments */}
          <div className="bg-gradient-to-br from-orange-500 to-orange-600 rounded-lg shadow-lg p-6 text-white">
            <div className="flex items-center justify-between mb-2">
              <h3 className="text-sm font-medium opacity-90">AI配送準備</h3>
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4" />
              </svg>
            </div>
            <div className="text-4xl font-bold mb-1">{stats?.shipments_prepared || 0}</div>
            <div className="text-sm opacity-90">件</div>
          </div>
        </div>

        {/* Detailed Stats */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
          {/* AI Performance */}
          <div className="bg-white rounded-lg shadow-lg p-6">
            <h2 className="text-xl font-bold text-gray-900 mb-4">AI性能</h2>
            <div className="space-y-4">
              <div>
                <div className="flex justify-between mb-2">
                  <span className="text-sm font-medium text-gray-700">平均確信度</span>
                  <span className="text-sm font-bold text-blue-600">{stats?.average_confidence.toFixed(1)}%</span>
                </div>
                <div className="w-full bg-gray-200 rounded-full h-2">
                  <div
                    className="bg-blue-600 h-2 rounded-full"
                    style={{ width: `${stats?.average_confidence || 0}%` }}
                  ></div>
                </div>
              </div>

              <div>
                <div className="flex justify-between mb-2">
                  <span className="text-sm font-medium text-gray-700">交渉承認率</span>
                  <span className="text-sm font-bold text-green-600">{stats?.acceptance_rate.toFixed(1)}%</span>
                </div>
                <div className="w-full bg-gray-200 rounded-full h-2">
                  <div
                    className="bg-green-600 h-2 rounded-full"
                    style={{ width: `${stats?.acceptance_rate || 0}%` }}
                  ></div>
                </div>
              </div>
            </div>
          </div>

          {/* Time Breakdown */}
          <div className="bg-white rounded-lg shadow-lg p-6">
            <h2 className="text-xl font-bold text-gray-900 mb-4">節約時間内訳</h2>
            <div className="space-y-3">
              <div className="flex items-center justify-between p-3 bg-blue-50 rounded-lg">
                <div>
                  <div className="font-medium text-gray-900">AI出品</div>
                  <div className="text-sm text-gray-600">{stats?.listings_created || 0}件 × 15分</div>
                </div>
                <div className="text-xl font-bold text-blue-600">
                  {((stats?.listings_created || 0) * 15 / 60).toFixed(1)}h
                </div>
              </div>

              <div className="flex items-center justify-between p-3 bg-purple-50 rounded-lg">
                <div>
                  <div className="font-medium text-gray-900">AI交渉</div>
                  <div className="text-sm text-gray-600">{stats?.negotiations_handled || 0}件 × 5分</div>
                </div>
                <div className="text-xl font-bold text-purple-600">
                  {((stats?.negotiations_handled || 0) * 5 / 60).toFixed(1)}h
                </div>
              </div>

              <div className="flex items-center justify-between p-3 bg-orange-50 rounded-lg">
                <div>
                  <div className="font-medium text-gray-900">AI配送準備</div>
                  <div className="text-sm text-gray-600">{stats?.shipments_prepared || 0}件 × 10分</div>
                </div>
                <div className="text-xl font-bold text-orange-600">
                  {((stats?.shipments_prepared || 0) * 10 / 60).toFixed(1)}h
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Benefits */}
        <div className="bg-gradient-to-r from-blue-500 to-purple-600 rounded-lg shadow-lg p-8 text-white">
          <h2 className="text-2xl font-bold mb-4">AIエージェントがもたらす効果</h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <div className="text-center">
              <div className="text-4xl font-bold mb-2">97%</div>
              <div className="text-sm opacity-90">出品時間削減</div>
            </div>
            <div className="text-center">
              <div className="text-4xl font-bold mb-2">100%</div>
              <div className="text-sm opacity-90">交渉時間削減</div>
            </div>
            <div className="text-center">
              <div className="text-4xl font-bold mb-2">90%</div>
              <div className="text-sm opacity-90">配送準備時間削減</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default AIAgentDashboard;

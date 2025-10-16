import { useState, useEffect } from 'react'
import { sustainabilityService, LeaderboardEntry } from '../services/sustainability'
import { Card } from '../components/common/Card'
import { Header } from '../components/layout/Header'

export const Leaderboard = () => {
  const [leaderboard, setLeaderboard] = useState<LeaderboardEntry[]>([])
  const [period, setPeriod] = useState('all')
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadLeaderboard()
  }, [period])

  const loadLeaderboard = async () => {
    try {
      setLoading(true)
      const data = await sustainabilityService.getLeaderboard(10, period)
      setLeaderboard(data)
    } catch (error) {
      console.error('Failed to load leaderboard:', error)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen">
      <Header />
      <div className="container mx-auto px-4 py-8 bg-white/70 backdrop-blur-sm">
        <div className="max-w-4xl mx-auto">
          <h1 className="text-3xl font-bold mb-6">ランキング</h1>

        <div className="mb-6">
          <select
            value={period}
            onChange={(e) => setPeriod(e.target.value)}
            className="px-4 py-2 border border-gray-300 rounded-lg"
          >
            <option value="all">All Time</option>
            <option value="month">This Month</option>
            <option value="week">This Week</option>
          </select>
        </div>

        {loading ? (
          <div>Loading...</div>
        ) : (
          <div className="space-y-4">
            {leaderboard.map((entry) => (
              <Card key={entry.user.id} padding="md">
                <div className="flex items-center justify-between">
                  <div className="flex items-center space-x-4">
                    <div className={`text-2xl font-bold ${
                      entry.rank === 1 ? 'text-yellow-500' :
                      entry.rank === 2 ? 'text-gray-400' :
                      entry.rank === 3 ? 'text-orange-600' :
                      'text-gray-600'
                    }`}>
                      #{entry.rank}
                    </div>
                    <div>
                      {entry.user.avatar_url ? (
                        <img
                          src={entry.user.avatar_url}
                          alt={entry.user.display_name}
                          className="w-12 h-12 rounded-full"
                        />
                      ) : (
                        <div className="w-12 h-12 rounded-full bg-gray-200 flex items-center justify-center">
                          {entry.user.display_name.charAt(0)}
                        </div>
                      )}
                    </div>
                    <div>
                      <div className="font-semibold">{entry.user.display_name}</div>
                      <div className="text-sm text-gray-600">Level {entry.level}</div>
                    </div>
                  </div>
                  <div className="text-right">
                    <div className="text-xl font-bold text-green-600">
                      {entry.total_co2_saved_kg.toFixed(2)} kg
                    </div>
                    <div className="text-sm text-gray-600">
                      CO2 Saved
                    </div>
                  </div>
                </div>
              </Card>
            ))}
          </div>
        )}
        </div>
      </div>
    </div>
  )
}

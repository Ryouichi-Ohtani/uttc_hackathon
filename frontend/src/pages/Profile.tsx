import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuthStore } from '@/store/authStore'
import { Button } from '@/components/common/Button'
import { Card } from '@/components/common/Card'
import { Header } from '@/components/layout/Header'
import api from '@/services/api'

interface DashboardData {
  total_co2_saved_kg: number
  level: number
  sustainability_score: number
  next_level_threshold: number
  achievements: any[]
  recent_logs: any[]
  monthly_stats: {
    current_month_co2_saved: number
    transactions: number
  }
  comparisons: {
    equivalent_trees: number
    car_km_avoided: number
  }
}

export const Profile = () => {
  const navigate = useNavigate()
  const { user, logout } = useAuthStore()
  const [dashboard, setDashboard] = useState<DashboardData | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadDashboard()
  }, [])

  const loadDashboard = async () => {
    try {
      const response = await api.get('/sustainability/dashboard')
      setDashboard(response.data)
    } catch (error) {
      console.error('Failed to load dashboard:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  if (!user) return null

  return (
    <div className="min-h-screen">
      <Header />

      <div className="max-w-7xl mx-auto px-4 py-8 bg-white/70 backdrop-blur-sm">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* User Info */}
          <Card className="lg:col-span-1">
            <div className="text-center">
              <div className="w-24 h-24 rounded-full bg-primary-500 flex items-center justify-center text-white text-3xl font-bold mx-auto mb-4">
                {user.level}
              </div>
              <h2 className="text-2xl font-bold text-gray-900">{user.display_name}</h2>
              <p className="text-gray-600">@{user.username}</p>
              <p className="text-sm text-gray-500 mt-2">{user.email}</p>

              <div className="mt-6 space-y-4">
                <div className="p-4 bg-green-50 rounded-lg">
                  <div className="text-3xl font-bold text-green-600">
                    {user.total_co2_saved_kg.toFixed(1)}kg
                  </div>
                  <div className="text-sm text-gray-600">Total CO2 Saved</div>
                </div>

                <div className="p-4 bg-blue-50 rounded-lg">
                  <div className="text-3xl font-bold text-blue-600">
                    Level {user.level}
                  </div>
                  <div className="text-sm text-gray-600">Sustainability Level</div>
                </div>

                <Button variant="outline" className="w-full" onClick={handleLogout}>
                  Logout
                </Button>
              </div>
            </div>
          </Card>

          {/* Dashboard */}
          <div className="lg:col-span-2 space-y-6">
            {loading ? (
              <div className="flex justify-center items-center h-64">
                <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-500" />
              </div>
            ) : dashboard ? (
              <>
                <Card>
                  <h3 className="text-lg font-semibold mb-4">This Month</h3>
                  <div className="grid grid-cols-2 gap-4">
                    <div className="p-4 bg-gray-50 rounded-lg">
                      <div className="text-2xl font-bold text-gray-900">
                        {dashboard.monthly_stats.current_month_co2_saved.toFixed(1)}kg
                      </div>
                      <div className="text-sm text-gray-600">CO2 Saved</div>
                    </div>
                    <div className="p-4 bg-gray-50 rounded-lg">
                      <div className="text-2xl font-bold text-gray-900">
                        {dashboard.monthly_stats.transactions}
                      </div>
                      <div className="text-sm text-gray-600">Transactions</div>
                    </div>
                  </div>
                </Card>

                <Card>
                  <h3 className="text-lg font-semibold mb-4">Environmental Impact</h3>
                  <div className="space-y-3">
                    <div className="flex items-center justify-between p-3 bg-green-50 rounded-lg">
                      <span className="text-gray-700">🌳 Equivalent Trees Planted</span>
                      <span className="font-bold text-green-600">
                        {dashboard.comparisons.equivalent_trees.toFixed(1)}
                      </span>
                    </div>
                    <div className="flex items-center justify-between p-3 bg-blue-50 rounded-lg">
                      <span className="text-gray-700">🚗 Car KM Avoided</span>
                      <span className="font-bold text-blue-600">
                        {dashboard.comparisons.car_km_avoided.toFixed(1)} km
                      </span>
                    </div>
                  </div>
                </Card>

                <Card>
                  <h3 className="text-lg font-semibold mb-4">Achievements</h3>
                  {dashboard.achievements.length > 0 ? (
                    <div className="grid grid-cols-2 gap-3">
                      {dashboard.achievements.map((achievement: any) => (
                        <div
                          key={achievement.id}
                          className="p-3 border border-gray-200 rounded-lg text-center"
                        >
                          <div className="text-2xl mb-1">🏆</div>
                          <div className="font-semibold text-sm">
                            {achievement.achievement?.name}
                          </div>
                        </div>
                      ))}
                    </div>
                  ) : (
                    <p className="text-gray-500 text-center py-4">
                      No achievements yet. Keep going!
                    </p>
                  )}
                </Card>

                <Card>
                  <h3 className="text-lg font-semibold mb-4">Recent Activity</h3>
                  {dashboard.recent_logs.length > 0 ? (
                    <div className="space-y-2">
                      {dashboard.recent_logs.map((log: any) => (
                        <div
                          key={log.id}
                          className="p-3 border-l-4 border-primary-500 bg-gray-50"
                        >
                          <div className="flex items-center justify-between">
                            <span className="text-sm text-gray-700">{log.description}</span>
                            <span className="text-sm font-semibold text-green-600">
                              +{log.co2_saved_kg.toFixed(2)}kg
                            </span>
                          </div>
                          <div className="text-xs text-gray-500 mt-1">
                            {new Date(log.created_at).toLocaleDateString()}
                          </div>
                        </div>
                      ))}
                    </div>
                  ) : (
                    <p className="text-gray-500 text-center py-4">No recent activity</p>
                  )}
                </Card>
              </>
            ) : null}
          </div>
        </div>
      </div>
    </div>
  )
}

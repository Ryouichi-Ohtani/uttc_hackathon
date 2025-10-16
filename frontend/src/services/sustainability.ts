import { api } from './api'
import { User, UserAchievement, SustainabilityLog, Product } from '../types'

export interface DashboardData {
  total_co2_saved_kg: number
  level: number
  sustainability_score: number
  next_level_threshold: number
  achievements: UserAchievement[]
  recent_logs: SustainabilityLog[]
  monthly_stats: {
    current_month_co2_saved: number
    transactions: number
  }
  comparisons: {
    equivalent_trees: number
    car_km_avoided: number
  }
}

export interface LeaderboardEntry {
  rank: number
  user: User
  total_co2_saved_kg: number
  sustainability_score: number
  level: number
}

export const sustainabilityService = {
  async getDashboard(): Promise<DashboardData> {
    const response = await api.get<DashboardData>('/sustainability/dashboard')
    return response.data
  },

  async getLeaderboard(limit: number = 10, period: string = 'all'): Promise<LeaderboardEntry[]> {
    const response = await api.get<{ leaderboard: LeaderboardEntry[] }>('/sustainability/leaderboard', {
      params: { limit, period }
    })
    return response.data.leaderboard
  },

  async getFavorites(): Promise<Product[]> {
    const response = await api.get<{ favorites: Product[] }>('/sustainability/favorites')
    return response.data.favorites
  }
}

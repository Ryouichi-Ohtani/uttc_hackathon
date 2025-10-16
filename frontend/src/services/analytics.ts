import { api } from './api'
import { Product } from '../types'

export interface UserBehaviorSummary {
  total_views: number
  total_searches: number
  total_likes: number
  total_offers: number
  favorite_categories: CategoryCount[]
  recently_viewed: Product[]
}

export interface CategoryCount {
  category: string
  count: number
}

export interface ProductAnalytics {
  product_id: string
  product: Product
  view_count: number
  like_count: number
  offer_count: number
}

export interface SearchKeyword {
  keyword: string
  count: number
}

export const analyticsService = {
  async getUserBehavior(): Promise<UserBehaviorSummary> {
    const response = await api.get<UserBehaviorSummary>('/analytics/user/behavior')
    return response.data
  },

  async getPopularProducts(days: number = 7): Promise<ProductAnalytics[]> {
    const response = await api.get<{ products: ProductAnalytics[] }>('/analytics/popular-products', {
      params: { days }
    })
    return response.data.products
  },

  async getSearchTrends(days: number = 7): Promise<SearchKeyword[]> {
    const response = await api.get<{ keywords: SearchKeyword[] }>('/analytics/search-trends', {
      params: { days }
    })
    return response.data.keywords
  }
}

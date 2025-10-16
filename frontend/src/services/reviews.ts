import { api } from './api'
import { Review } from '../types'

export interface CreateReviewRequest {
  product_id: string
  purchase_id: string
  rating: number
  comment: string
}

export interface ReviewsResponse {
  reviews: Review[]
  average_rating: number
}

export const reviewService = {
  async create(data: CreateReviewRequest): Promise<Review> {
    const response = await api.post<Review>('/reviews', data)
    return response.data
  },

  async getProductReviews(productId: string): Promise<ReviewsResponse> {
    const response = await api.get<ReviewsResponse>(`/reviews/products/${productId}`)
    return response.data
  }
}

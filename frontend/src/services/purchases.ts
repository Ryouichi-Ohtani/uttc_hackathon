import { api } from './api'

export interface Purchase {
  id: string
  product_id: string
  buyer_id: string
  seller_id: string
  price: number
  co2_saved_kg: number
  status: 'pending' | 'completed' | 'cancelled'
  payment_method: string
  shipping_address: string
  completed_at?: string
  created_at: string
  product?: any
  buyer?: any
  seller?: any
}

export interface CreatePurchaseRequest {
  product_id: string
  shipping_address: string
  payment_method: string
}

export interface PurchaseListResponse {
  purchases: Purchase[]
  pagination: {
    page: number
    limit: number
    total: number
    total_pages: number
  }
}

class PurchaseService {
  async create(data: CreatePurchaseRequest): Promise<Purchase> {
    const response = await api.post('/purchases', data)
    return response.data
  }

  async getById(id: string): Promise<Purchase> {
    const response = await api.get(`/purchases/${id}`)
    return response.data
  }

  async list(role?: 'buyer' | 'seller', page = 1, limit = 20): Promise<PurchaseListResponse> {
    const params = new URLSearchParams({
      page: page.toString(),
      limit: limit.toString(),
    })
    if (role) {
      params.append('role', role)
    }
    const response = await api.get(`/purchases?${params}`)
    return response.data
  }

  async complete(id: string): Promise<void> {
    await api.patch(`/purchases/${id}/complete`)
  }
}

export const purchaseService = new PurchaseService()

import { api } from './api'
import { Purchase, ShippingLabel } from '@/types'

export interface CreatePurchaseRequest {
  product_id: string
  shipping_address: string
  payment_method: string
  // Delivery information
  delivery_date?: string
  delivery_time_slot?: 'morning' | 'afternoon' | 'evening' | 'anytime'
  use_registered_address?: boolean
  recipient_name?: string
  recipient_phone_number?: string
  recipient_postal_code?: string
  recipient_prefecture?: string
  recipient_city?: string
  recipient_address_line1?: string
  recipient_address_line2?: string
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

  async getShippingLabel(purchaseId: string): Promise<ShippingLabel> {
    const response = await api.get(`/purchases/${purchaseId}/shipping-label`)
    return response.data
  }

  async generateShippingLabel(purchaseId: string): Promise<ShippingLabel> {
    const response = await api.post(`/purchases/${purchaseId}/shipping-label`)
    return response.data
  }
}

export const purchaseService = new PurchaseService()

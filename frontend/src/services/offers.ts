import { api } from './api'

export interface Offer {
  id: string
  product_id: string
  buyer_id: string
  offer_price: number
  message: string
  status: 'pending' | 'accepted' | 'rejected' | 'cancelled'
  response_message: string
  created_at: string
  updated_at: string
  responded_at?: string
  product?: any
  buyer?: any
}

export interface CreateOfferRequest {
  product_id: string
  offer_price: number
  message: string
}

export interface RespondOfferRequest {
  accept: boolean
  message: string
}

export const offerService = {
  async create(data: CreateOfferRequest): Promise<Offer> {
    const response = await api.post<Offer>('/offers', data)
    return response.data
  },

  async getMyOffers(role: 'buyer' | 'seller' = 'buyer'): Promise<Offer[]> {
    const response = await api.get<{ offers: Offer[] }>('/offers/my', {
      params: { role }
    })
    return response.data.offers
  },

  async getProductOffers(productId: string): Promise<Offer[]> {
    const response = await api.get<{ offers: Offer[] }>(`/offers/products/${productId}`)
    return response.data.offers
  },

  async respond(offerId: string, data: RespondOfferRequest): Promise<Offer> {
    const response = await api.patch<Offer>(`/offers/${offerId}/respond`, data)
    return response.data
  }
}

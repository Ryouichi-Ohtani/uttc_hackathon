import api from '@/lib/api'
import {
  PaymentAuthorizationRequest,
  PaymentAuthorizationResponse,
  AutoPurchaseWatch,
  CreateAutoPurchaseWatchRequest,
} from '@/types'

class AutoPurchaseService {
  /**
   * Authorize payment for auto-purchase
   */
  async authorizePayment(
    request: PaymentAuthorizationRequest
  ): Promise<PaymentAuthorizationResponse> {
    const response = await api.post('/auto-purchases/authorize-payment', request)
    return response.data
  }

  /**
   * Create a new auto-purchase watch
   */
  async createWatch(
    request: CreateAutoPurchaseWatchRequest
  ): Promise<AutoPurchaseWatch> {
    const response = await api.post('/auto-purchases', request)
    return response.data
  }

  /**
   * Get all watches for the authenticated user
   */
  async getUserWatches(): Promise<AutoPurchaseWatch[]> {
    const response = await api.get('/auto-purchases')
    return response.data.watches
  }

  /**
   * Get a specific watch by ID
   */
  async getWatch(watchId: string): Promise<AutoPurchaseWatch> {
    const response = await api.get(`/auto-purchases/${watchId}`)
    return response.data
  }

  /**
   * Cancel an auto-purchase watch
   */
  async cancelWatch(watchId: string): Promise<void> {
    await api.delete(`/auto-purchases/${watchId}`)
  }
}

export const autoPurchaseService = new AutoPurchaseService()

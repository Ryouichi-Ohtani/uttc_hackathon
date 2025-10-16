import api from './api'
import { Product, ProductFilters, PaginationResponse, CO2Comparison } from '@/types'

export const productService = {
  async list(filters?: ProductFilters) {
    const response = await api.get<{
      products: Product[]
      pagination: PaginationResponse
    }>('/products', { params: filters })
    return response.data
  },

  async getById(id: string) {
    const response = await api.get<Product>(`/products/${id}`)
    return response.data
  },

  async create(data: any) {
    const response = await api.post<Product>('/products', data)
    return response.data
  },

  async update(id: string, data: any) {
    const response = await api.put<Product>(`/products/${id}`, data)
    return response.data
  },

  async delete(id: string) {
    const response = await api.delete(`/products/${id}`)
    return response.data
  },

  async addFavorite(id: string) {
    const response = await api.post(`/products/${id}/favorite`)
    return response.data
  },

  async removeFavorite(id: string) {
    const response = await api.delete(`/products/${id}/favorite`)
    return response.data
  },

  getCO2Comparison(product: Product): CO2Comparison {
    const savedKg = product.co2_impact_kg
    const buyingNewKg = savedKg / 0.33
    const buyingUsedKg = buyingNewKg - savedKg

    return {
      buying_new_kg: buyingNewKg,
      buying_used_kg: buyingUsedKg,
      saved_kg: savedKg,
      equivalent_trees: savedKg / 20,
    }
  },
}

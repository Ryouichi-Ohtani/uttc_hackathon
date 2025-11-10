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
    const savedKg = product.co2_impact_kg ?? 0
    const newProductCO2 = savedKg / 0.33
    const productCO2 = newProductCO2 - savedKg

    return {
      product_co2: productCO2,
      new_product_co2: newProductCO2,
      saved_co2: savedKg,
      saved_percentage: savedKg > 0 ? (savedKg / newProductCO2) * 100 : 0,
    }
  },
}

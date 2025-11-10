import api from '@/lib/api'
import { User, Product, Purchase, PaginationResponse } from '@/types'

interface AdminPaginationResponse {
  page: number
  limit: number
  total: number
  total_pages: number
}

interface UsersListResponse {
  users: User[]
  pagination: AdminPaginationResponse
}

interface ProductsListResponse {
  products: Product[]
  pagination: PaginationResponse
}

interface PurchasesListResponse {
  purchases: Purchase[]
  pagination: PaginationResponse
}

interface UpdateUserRequest {
  display_name?: string
  role?: string
  bio?: string
  postal_code?: string
  prefecture?: string
  city?: string
  address_line1?: string
  address_line2?: string
  phone_number?: string
}

interface UpdateProductRequest {
  title?: string
  description?: string
  price?: number
  category?: string
  condition?: 'new' | 'like_new' | 'good' | 'fair' | 'poor'
  status?: 'draft' | 'active' | 'sold' | 'reserved' | 'deleted'
}

interface UpdatePurchaseRequest {
  status: 'pending' | 'completed' | 'cancelled'
}

class AdminService {
  // User Management
  async getUsers(page = 1, limit = 100): Promise<UsersListResponse> {
    const response = await api.get('/admin/users', {
      params: { page, limit },
    })
    return response.data
  }

  async updateUser(userId: string, updates: UpdateUserRequest): Promise<User> {
    const response = await api.put(`/admin/users/${userId}`, updates)
    return response.data
  }

  async deleteUser(userId: string): Promise<void> {
    await api.delete(`/admin/users/${userId}`)
  }

  // Product Management
  async getProducts(page = 1, limit = 100): Promise<ProductsListResponse> {
    const response = await api.get('/admin/products', {
      params: { page, limit },
    })
    return response.data
  }

  async updateProduct(
    productId: string,
    updates: UpdateProductRequest
  ): Promise<Product> {
    const response = await api.put(`/admin/products/${productId}`, updates)
    return response.data
  }

  async deleteProduct(productId: string): Promise<void> {
    await api.delete(`/admin/products/${productId}`)
  }

  // Purchase Management
  async getPurchases(page = 1, limit = 100): Promise<PurchasesListResponse> {
    const response = await api.get('/admin/purchases', {
      params: { page, limit },
    })
    return response.data
  }

  async updatePurchase(
    purchaseId: string,
    updates: UpdatePurchaseRequest
  ): Promise<Purchase> {
    const response = await api.put(`/admin/purchases/${purchaseId}`, updates)
    return response.data
  }
}

export const adminService = new AdminService()

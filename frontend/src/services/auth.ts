import api from './api'
import { AuthResponse } from '@/types'

export const authService = {
  async register(data: {
    email: string
    username: string
    password: string
    display_name: string
  }): Promise<AuthResponse> {
    const response = await api.post<AuthResponse>('/auth/register', data)
    return response.data
  },

  async login(data: { email: string; password: string }): Promise<AuthResponse> {
    const response = await api.post<AuthResponse>('/auth/login', data)
    return response.data
  },

  async getMe() {
    const response = await api.get('/auth/me')
    return response.data
  },
}

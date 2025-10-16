import { api } from './api'

export interface Notification {
  id: string
  user_id: string
  type: 'message' | 'purchase' | 'favorite' | 'review'
  title: string
  message: string
  link: string
  is_read: boolean
  created_at: string
}

export const notificationService = {
  async getNotifications(): Promise<Notification[]> {
    const response = await api.get<Notification[]>('/notifications')
    return response.data
  },

  async getUnreadCount(): Promise<number> {
    const response = await api.get<{ count: number }>('/notifications/unread-count')
    return response.data.count
  },

  async markAsRead(id: string): Promise<void> {
    await api.patch(`/notifications/${id}/read`)
  },

  async markAllAsRead(): Promise<void> {
    await api.post('/notifications/read-all')
  }
}

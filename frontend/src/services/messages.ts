import { api } from './api'

export interface Conversation {
  id: string
  product_id?: string
  product?: any
  participants: ConversationParticipant[]
  messages?: Message[]
  last_message_at: string
  created_at: string
}

export interface ConversationParticipant {
  id: string
  conversation_id: string
  user_id: string
  user?: any
  last_read_at: string
  created_at: string
}

export interface Message {
  id: string
  conversation_id: string
  sender_id: string
  sender?: any
  content: string
  is_read: boolean
  created_at: string
}

export interface MessageListResponse {
  messages: Message[]
  pagination: {
    page: number
    limit: number
    total: number
    total_pages: number
  }
}

class MessageService {
  async getOrCreateConversation(productId: string, sellerId: string): Promise<Conversation> {
    const response = await api.get(`/conversations/product/${productId}/seller/${sellerId}`)
    return response.data
  }

  async listConversations(): Promise<Conversation[]> {
    const response = await api.get('/conversations')
    return response.data.conversations
  }

  async getMessages(conversationId: string, page = 1, limit = 50): Promise<MessageListResponse> {
    const response = await api.get(`/conversations/${conversationId}/messages`, {
      params: { page, limit },
    })
    return response.data
  }

  async sendMessage(conversationId: string, content: string): Promise<Message> {
    const response = await api.post(`/conversations/${conversationId}/messages`, { content })
    return response.data
  }

  createWebSocket(conversationId: string, token: string): WebSocket {
    const wsUrl = import.meta.env.VITE_WS_URL || 'ws://localhost:8080'
    const ws = new WebSocket(`${wsUrl}/v1/ws/conversations/${conversationId}`)

    ws.onopen = () => {
      // Send auth message
      ws.send(JSON.stringify({
        type: 'auth',
        token: token,
      }))
    }

    return ws
  }
}

export const messageService = new MessageService()

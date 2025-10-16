import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { messageService, Conversation } from '@/services/messages'
import { Button } from '@/components/common/Button'
import { Card } from '@/components/common/Card'
import { Header } from '@/components/layout/Header'
import { useAuthStore } from '@/store/authStore'
import toast from 'react-hot-toast'

export const Messages = () => {
  const navigate = useNavigate()
  const { user } = useAuthStore()
  const [conversations, setConversations] = useState<Conversation[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadConversations()
  }, [])

  const loadConversations = async () => {
    try {
      setLoading(true)
      const data = await messageService.listConversations()
      setConversations(data)
    } catch (error) {
      toast.error('Failed to load conversations')
    } finally {
      setLoading(false)
    }
  }

  const getOtherParticipant = (conv: Conversation) => {
    return conv.participants.find((p) => p.user_id !== user?.id)?.user
  }

  const getLastMessage = (conv: Conversation) => {
    return conv.messages?.[0]
  }

  return (
    <div className="min-h-screen">
      <Header />

      <div className="max-w-4xl mx-auto px-4 py-8 bg-white/70 backdrop-blur-sm">
        {loading ? (
          <div className="flex justify-center items-center h-64">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-500" />
          </div>
        ) : conversations.length === 0 ? (
          <Card className="text-center py-12">
            <div className="text-6xl mb-4">💬</div>
            <h3 className="text-xl font-semibold text-gray-900 mb-2">No messages yet</h3>
            <p className="text-gray-600">Start chatting with sellers to see messages here</p>
            <Button className="mt-4" onClick={() => navigate('/')}>
              Browse Products
            </Button>
          </Card>
        ) : (
          <div className="space-y-3">
            {conversations.map((conv) => {
              const otherUser = getOtherParticipant(conv)
              const lastMsg = getLastMessage(conv)

              return (
                <Card
                  key={conv.id}
                  className="cursor-pointer hover:shadow-lg transition"
                  onClick={() => navigate(`/chat/${conv.id}`)}
                >
                  <div className="flex items-start gap-4">
                    <div className="w-12 h-12 rounded-full bg-primary-100 flex items-center justify-center text-lg font-bold flex-shrink-0">
                      {otherUser?.username?.[0]?.toUpperCase() || '?'}
                    </div>
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center justify-between mb-1">
                        <h3 className="font-semibold text-gray-900 truncate">
                          {otherUser?.display_name || 'Unknown User'}
                        </h3>
                        <span className="text-xs text-gray-500 flex-shrink-0">
                          {new Date(conv.last_message_at).toLocaleDateString()}
                        </span>
                      </div>
                      {conv.product && (
                        <p className="text-sm text-gray-600 mb-1">
                          📦 {conv.product.title}
                        </p>
                      )}
                      {lastMsg && (
                        <p className="text-sm text-gray-500 truncate">
                          {lastMsg.sender_id === user?.id ? 'You: ' : ''}
                          {lastMsg.content}
                        </p>
                      )}
                    </div>
                  </div>
                </Card>
              )
            })}
          </div>
        )}
      </div>
    </div>
  )
}

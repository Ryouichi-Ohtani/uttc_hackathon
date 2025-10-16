import { useState, useEffect, useRef } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { messageService, Message } from '@/services/messages'
import { Button } from '@/components/common/Button'
import { Input } from '@/components/common/Input'
import { Header } from '@/components/layout/Header'
import { useAuthStore } from '@/store/authStore'
import toast from 'react-hot-toast'

export const Chat = () => {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const { user, token } = useAuthStore()
  const [messages, setMessages] = useState<Message[]>([])
  const [newMessage, setNewMessage] = useState('')
  const [loading, setLoading] = useState(true)
  const [ws, setWs] = useState<WebSocket | null>(null)
  const messagesEndRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    if (id) {
      loadMessages()
      setupWebSocket()
    }

    return () => {
      if (ws) {
        ws.close()
      }
    }
  }, [id])

  useEffect(() => {
    scrollToBottom()
  }, [messages])

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }

  const loadMessages = async () => {
    try {
      setLoading(true)
      const data = await messageService.getMessages(id!)
      setMessages(data.messages.reverse())
    } catch (error) {
      toast.error('Failed to load messages')
    } finally {
      setLoading(false)
    }
  }

  const setupWebSocket = () => {
    if (!id || !token) return

    const websocket = messageService.createWebSocket(id, token)

    websocket.onmessage = (event) => {
      const data = JSON.parse(event.data)

      if (data.type === 'message' && data.data) {
        setMessages((prev) => [...prev, data.data])
      }
    }

    websocket.onerror = (error) => {
      console.error('WebSocket error:', error)
      toast.error('Real-time messaging connection failed')
    }

    setWs(websocket)
  }

  const handleSend = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!newMessage.trim() || !id) return

    try {
      // Send via WebSocket if available
      if (ws && ws.readyState === WebSocket.OPEN) {
        ws.send(JSON.stringify({
          type: 'send_message',
          content: newMessage,
        }))
      } else {
        // Fallback to HTTP
        const message = await messageService.sendMessage(id, newMessage)
        setMessages((prev) => [...prev, message])
      }

      setNewMessage('')
    } catch (error) {
      toast.error('Failed to send message')
    }
  }

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-500" />
      </div>
    )
  }

  return (
    <div className="h-screen flex flex-col bg-gray-50">
      <Header />

      {/* Messages */}
      <div className="flex-1 overflow-y-auto px-4 py-6">
        <div className="max-w-3xl mx-auto space-y-4">
          {messages.map((msg) => {
            const isOwnMessage = msg.sender_id === user?.id

            return (
              <div
                key={msg.id}
                className={`flex ${isOwnMessage ? 'justify-end' : 'justify-start'}`}
              >
                <div
                  className={`max-w-md px-4 py-2 rounded-2xl ${
                    isOwnMessage
                      ? 'bg-primary-500 text-white'
                      : 'bg-white border border-gray-200'
                  }`}
                >
                  <p className="text-sm">{msg.content}</p>
                  <p
                    className={`text-xs mt-1 ${
                      isOwnMessage ? 'text-primary-100' : 'text-gray-400'
                    }`}
                  >
                    {new Date(msg.created_at).toLocaleTimeString()}
                  </p>
                </div>
              </div>
            )
          })}
          <div ref={messagesEndRef} />
        </div>
      </div>

      {/* Input */}
      <div className="bg-white border-t border-gray-200 px-4 py-4">
        <form onSubmit={handleSend} className="max-w-3xl mx-auto flex gap-2">
          <Input
            type="text"
            value={newMessage}
            onChange={(e) => setNewMessage(e.target.value)}
            placeholder="Type a message..."
            className="flex-1"
          />
          <Button type="submit" disabled={!newMessage.trim()}>
            Send
          </Button>
        </form>
      </div>
    </div>
  )
}

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
  const [aiSuggestion, setAiSuggestion] = useState<string | null>(null)
  const [aiLoading, setAiLoading] = useState(false)
  const messagesEndRef = useRef<HTMLDivElement>(null)
  const wsRef = useRef<WebSocket | null>(null)

  useEffect(() => {
    if (id) {
      loadMessages()
      setupWebSocket()
    }

    return () => {
      if (wsRef.current) {
        wsRef.current.close()
        wsRef.current = null
      }
    }
  }, [id, token])

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
      setMessages(data.messages ? data.messages.reverse() : [])
    } catch (error: any) {
      console.error('Failed to load messages:', error)
      toast.error(error.response?.data?.error || 'Failed to load messages')
      setMessages([])
    } finally {
      setLoading(false)
    }
  }

  const setupWebSocket = () => {
    if (!id || !token) return

    // Close existing WebSocket if any
    if (wsRef.current) {
      wsRef.current.close()
    }

    const websocket = messageService.createWebSocket(id, token)

    websocket.onmessage = (event) => {
      const data = JSON.parse(event.data)

      if (data.type === 'message' && data.data) {
        // Prevent duplicate messages by checking if message already exists
        setMessages((prev) => {
          const exists = prev.some((msg) => msg.id === data.data.id)
          if (exists) {
            return prev
          }
          return [...prev, data.data]
        })
      }
    }

    websocket.onerror = (error) => {
      console.error('WebSocket error:', error)
      // Don't show error toast - will fall back to HTTP
    }

    websocket.onclose = (event) => {
      console.log('WebSocket closed:', event.code, event.reason)
      // Connection closed, will use HTTP fallback for sending messages
    }

    wsRef.current = websocket
    setWs(websocket)
  }

  const handleSend = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!newMessage.trim() || !id) return

    const messageContent = newMessage
    setNewMessage('') // Clear input immediately for better UX

    try {
      // Send via WebSocket if available
      if (ws && ws.readyState === WebSocket.OPEN) {
        ws.send(JSON.stringify({
          type: 'send_message',
          content: messageContent,
        }))
        // Message will be added via WebSocket onmessage handler
      } else {
        // Fallback to HTTP
        const message = await messageService.sendMessage(id, messageContent)
        setMessages((prev) => [...prev, message])
      }
      setAiSuggestion(null)
    } catch (error) {
      toast.error('Failed to send message')
      setNewMessage(messageContent) // Restore message on error
    }
  }

  const handleAISuggestion = async () => {
    if (!id || aiLoading) return

    try {
      setAiLoading(true)
      const suggestion = await messageService.requestSuggestion(id)
      setAiSuggestion(suggestion)
      toast.success('AI提案を取得しました')
    } catch (error) {
      console.error('AI suggestion error:', error)
      toast.error('AI提案の取得に失敗しました')
    } finally {
      setAiLoading(false)
    }
  }

  const applyAISuggestion = () => {
    if (aiSuggestion) {
      setNewMessage(aiSuggestion)
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
          {messages.length === 0 && !loading && (
            <div className="text-center py-12 text-gray-500">
              <p className="text-lg mb-2">💬</p>
              <p>No messages yet. Start the conversation!</p>
            </div>
          )}
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
        <div className="max-w-3xl mx-auto space-y-2">
          <div className="flex flex-wrap items-center gap-2">
            <Button
              size="sm"
              variant="outline"
              onClick={handleAISuggestion}
              disabled={aiLoading}
              className="bg-gray-50 text-gray-800"
            >
              {aiLoading ? 'AI提案を取得中…' : 'AIにメッセージ案を作ってもらう'}
            </Button>
            {aiSuggestion && (
              <button
                type="button"
                className="text-xs font-semibold text-primary-700 underline"
                onClick={applyAISuggestion}
              >
                提案を入力
              </button>
            )}
          </div>

          {aiSuggestion && (
            <div className="rounded-lg border border-yellow-200 bg-yellow-50 p-3 text-sm text-gray-800">
              <p className="whitespace-pre-wrap">{aiSuggestion}</p>
              <p className="mt-2 text-xs text-gray-500">AIからの提案</p>
            </div>
          )}

          <form onSubmit={handleSend} className="flex gap-2">
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
    </div>
  )
}

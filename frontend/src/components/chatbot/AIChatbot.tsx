import { useState, useRef, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { api } from '@/services/api'
import { Button } from '@/components/common/Button'
import toast from 'react-hot-toast'

interface Message {
  role: 'user' | 'assistant'
  content: string
}

interface AIChatbotProps {
  isOpen: boolean
  onClose: () => void
}

interface ProductCard {
  id: string
  name: string
  imageUrl: string
}

export const AIChatbot = ({ isOpen, onClose }: AIChatbotProps) => {
  const navigate = useNavigate()
  const [messages, setMessages] = useState<Message[]>([
    {
      role: 'assistant',
      content: 'こんにちは！EcoMateのAIアシスタントです🌱\n\n商品探しやエコに関する質問など、何でもお手伝いします！'
    }
  ])
  const [input, setInput] = useState('')
  const [loading, setLoading] = useState(false)
  const messagesEndRef = useRef<HTMLDivElement>(null)

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }

  useEffect(() => {
    scrollToBottom()
  }, [messages])

  const handleSend = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!input.trim() || loading) return

    const userMessage: Message = {
      role: 'user',
      content: input
    }

    setMessages(prev => [...prev, userMessage])
    setInput('')
    setLoading(true)

    try {
      const response = await api.post('/chatbot/chat', {
        messages: [...messages, userMessage],
        context: ''
      })

      const assistantMessage: Message = {
        role: 'assistant',
        content: response.data.message
      }

      setMessages(prev => [...prev, assistantMessage])
    } catch (error: any) {
      toast.error('メッセージの送信に失敗しました')
      console.error('Chatbot error:', error)
    } finally {
      setLoading(false)
    }
  }

  const parseProductCards = (content: string): { text: string; products: ProductCard[] } => {
    const productRegex = /\[PRODUCT:([^:]+):([^:]+):([^\]]+)\]/g
    const products: ProductCard[] = []
    let match

    while ((match = productRegex.exec(content)) !== null) {
      products.push({
        id: match[1],
        name: match[2],
        imageUrl: match[3]
      })
    }

    const text = content.replace(productRegex, '').trim()
    return { text, products }
  }

  const handleProductClick = (productId: string) => {
    navigate(`/products/${productId}`)
    onClose()
  }

  if (!isOpen) return null

  return (
    <div className="fixed right-0 top-0 h-screen w-96 bg-white shadow-2xl border-l border-gray-200 flex flex-col z-50">
      {/* Header */}
      <div className="bg-gradient-to-r from-green-500 to-emerald-600 text-white px-4 py-4 flex items-center justify-between">
        <div className="flex items-center gap-2">
          <span className="text-2xl">💬</span>
          <div>
            <h3 className="font-semibold">AIアシスタント</h3>
            <p className="text-xs text-green-100">powered by Gemini</p>
          </div>
        </div>
        <button
          onClick={onClose}
          className="text-white hover:bg-white/20 rounded-full p-2 transition"
        >
          ✕
        </button>
      </div>

      {/* Messages */}
      <div className="flex-1 overflow-y-auto p-4 space-y-4">
        {messages.map((msg, index) => {
          const { text, products } = msg.role === 'assistant' ? parseProductCards(msg.content) : { text: msg.content, products: [] }

          return (
            <div
              key={index}
              className={`flex ${msg.role === 'user' ? 'justify-end' : 'justify-start'}`}
            >
              <div className={`max-w-[85%] ${msg.role === 'user' ? '' : 'w-full'}`}>
                <div
                  className={`rounded-2xl px-4 py-3 ${
                    msg.role === 'user'
                      ? 'bg-primary-500 text-white'
                      : 'bg-gray-100 text-gray-900 border border-gray-200'
                  }`}
                >
                  {msg.role === 'assistant' && (
                    <div className="flex items-center gap-2 mb-2">
                      <span className="text-lg">🌱</span>
                      <span className="text-xs font-semibold text-primary-600">EcoMate AI</span>
                    </div>
                  )}
                  <p className="text-sm whitespace-pre-wrap">{text}</p>
                </div>

                {/* Product Cards */}
                {products.length > 0 && (
                  <div className="mt-3 space-y-2">
                    {products.map((product, idx) => (
                      <button
                        key={idx}
                        onClick={() => handleProductClick(product.id)}
                        className="w-full bg-white border border-gray-200 rounded-lg p-3 hover:shadow-md transition-shadow flex items-center gap-3 text-left"
                      >
                        <img
                          src={product.imageUrl || 'https://via.placeholder.com/80x80/10B981/FFFFFF?text=EcoMate'}
                          alt={product.name}
                          className="w-16 h-16 object-cover rounded"
                          onError={(e) => {
                            const target = e.target as HTMLImageElement
                            target.src = 'https://via.placeholder.com/80x80/10B981/FFFFFF?text=EcoMate'
                          }}
                        />
                        <div className="flex-1">
                          <p className="text-sm font-semibold text-gray-900">{product.name}</p>
                          <p className="text-xs text-primary-600 mt-1">→ 商品詳細を見る</p>
                        </div>
                      </button>
                    ))}
                  </div>
                )}
              </div>
            </div>
          )
        })}
        {loading && (
          <div className="flex justify-start">
            <div className="bg-gray-100 border border-gray-200 rounded-2xl px-4 py-3">
              <div className="flex items-center gap-2">
                <div className="flex space-x-1">
                  <div className="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style={{ animationDelay: '0ms' }}></div>
                  <div className="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style={{ animationDelay: '150ms' }}></div>
                  <div className="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style={{ animationDelay: '300ms' }}></div>
                </div>
                <span className="text-xs text-gray-500">考え中...</span>
              </div>
            </div>
          </div>
        )}
        <div ref={messagesEndRef} />
      </div>

      {/* Quick Actions */}
      <div className="px-4 py-2 border-t border-gray-100">
        <div className="text-xs text-gray-500 mb-2">クイック質問</div>
        <div className="flex flex-wrap gap-2">
          {['おすすめの商品は？', 'スニーカー探してる', 'CO2削減について'].map((question) => (
            <button
              key={question}
              onClick={() => setInput(question)}
              className="text-xs px-3 py-1 bg-gray-100 hover:bg-gray-200 rounded-full transition"
              disabled={loading}
            >
              {question}
            </button>
          ))}
        </div>
      </div>

      {/* Input */}
      <form onSubmit={handleSend} className="p-4 border-t border-gray-200 bg-gray-50">
        <div className="flex gap-2">
          <input
            type="text"
            value={input}
            onChange={(e) => setInput(e.target.value)}
            placeholder="メッセージを入力..."
            className="flex-1 px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 text-sm"
            disabled={loading}
          />
          <Button type="submit" disabled={loading || !input.trim()}>
            {loading ? '...' : '送信'}
          </Button>
        </div>
      </form>
    </div>
  )
}

import { useState } from 'react'
import { Button } from '../common/Button'
import { Input } from '../common/Input'
import api from '@/services/api'
import toast from 'react-hot-toast'

interface OfferDialogProps {
  productId: string
  currentPrice: number
  onClose: () => void
  onSuccess: () => void
}

export const OfferDialog = ({ productId, currentPrice, onClose, onSuccess }: OfferDialogProps) => {
  const [offerPrice, setOfferPrice] = useState('')
  const [message, setMessage] = useState('')
  const [loading, setLoading] = useState(false)
  const [aiSuggestion, setAiSuggestion] = useState<any>(null)
  const [loadingSuggestion, setLoadingSuggestion] = useState(false)

  const handleGetSuggestion = async () => {
    try {
      setLoadingSuggestion(true)
      const response = await api.get(`/v1/offers/products/${productId}/ai-suggestion`)
      setAiSuggestion(response.data)

      if (response.data.suggested_price) {
        setOfferPrice(response.data.suggested_price.toString())
      }

      toast.success('AI価格提案を取得しました！')
    } catch (error) {
      toast.error('AI提案の取得に失敗しました')
    } finally {
      setLoadingSuggestion(false)
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    const price = parseInt(offerPrice)
    if (isNaN(price) || price <= 0) {
      toast.error('有効な金額を入力してください')
      return
    }

    if (price >= currentPrice) {
      toast.error('提示価格は現在の価格より低くしてください')
      return
    }

    setLoading(true)
    try {
      await api.post('/v1/offers', {
        product_id: productId,
        offered_price: price,
        message: message || undefined,
      })

      toast.success('価格交渉を送信しました！')
      onSuccess()
      onClose()
    } catch (error: any) {
      toast.error(error.response?.data?.error || '送信に失敗しました')
    } finally {
      setLoading(false)
    }
  }

  const discount = offerPrice ? Math.round(((currentPrice - parseInt(offerPrice)) / currentPrice) * 100) : 0

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-lg max-w-md w-full p-6 max-h-[90vh] overflow-y-auto">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-xl font-bold">価格交渉</h2>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-gray-600 text-2xl"
          >
            ×
          </button>
        </div>

        <div className="mb-6 p-4 bg-gray-50 rounded-lg">
          <div className="text-sm text-gray-600 mb-1">現在の価格</div>
          <div className="text-2xl font-bold">¥{currentPrice.toLocaleString()}</div>
        </div>

        {!aiSuggestion ? (
          <Button
            variant="outline"
            className="w-full mb-4"
            onClick={handleGetSuggestion}
            disabled={loadingSuggestion}
          >
            {loadingSuggestion ? '🤖 AI分析中...' : '🤖 AI価格提案を受ける'}
          </Button>
        ) : (
          <div className="mb-4 p-4 bg-primary-50 rounded-lg border border-primary-200">
            <div className="flex items-start gap-2 mb-2">
              <span className="text-2xl">🤖</span>
              <div className="flex-1">
                <div className="font-semibold text-primary-900 mb-1">AI価格提案</div>
                <div className="text-xl font-bold text-primary-600 mb-2">
                  ¥{aiSuggestion.suggested_price?.toLocaleString()}
                </div>
                <div className="text-sm text-gray-700 whitespace-pre-wrap">
                  {aiSuggestion.reasoning}
                </div>
              </div>
            </div>
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              提示価格 *
            </label>
            <Input
              type="number"
              required
              value={offerPrice}
              onChange={(e) => setOfferPrice(e.target.value)}
              placeholder="例: 5000"
              min="1"
              max={currentPrice - 1}
            />
            {offerPrice && discount > 0 && (
              <div className="mt-2 text-sm">
                <span className="text-green-600 font-semibold">
                  {discount}% OFF
                </span>
                <span className="text-gray-600 ml-2">
                  (¥{(currentPrice - parseInt(offerPrice)).toLocaleString()} 安く)
                </span>
              </div>
            )}
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              メッセージ (任意)
            </label>
            <textarea
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent"
              rows={3}
              value={message}
              onChange={(e) => setMessage(e.target.value)}
              placeholder="出品者へのメッセージ..."
            />
          </div>

          <div className="flex gap-2">
            <Button
              type="button"
              variant="outline"
              className="flex-1"
              onClick={onClose}
            >
              キャンセル
            </Button>
            <Button type="submit" className="flex-1" disabled={loading}>
              {loading ? '送信中...' : '交渉を送信'}
            </Button>
          </div>
        </form>

        <div className="mt-4 p-3 bg-blue-50 rounded text-sm text-gray-700">
          💡 ヒント: 出品者が承認すると、新しい価格で購入できます
        </div>
      </div>
    </div>
  )
}

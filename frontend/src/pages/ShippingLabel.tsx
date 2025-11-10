import { useState, useEffect } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { purchaseService } from '@/services/purchases'
import { ShippingLabelView } from '@/components/shipping/ShippingLabelView'
import { Button } from '@/components/common/Button'
import { Card } from '@/components/common/Card'
import toast from 'react-hot-toast'
import { ShippingLabel as ShippingLabelType } from '@/types'

export const ShippingLabel = () => {
  const { purchaseId } = useParams<{ purchaseId: string }>()
  const navigate = useNavigate()
  const [label, setLabel] = useState<ShippingLabelType | null>(null)
  const [loading, setLoading] = useState(true)
  const [generating, setGenerating] = useState(false)

  useEffect(() => {
    if (purchaseId) {
      loadShippingLabel()
    }
  }, [purchaseId])

  const loadShippingLabel = async () => {
    if (!purchaseId) return

    try {
      setLoading(true)
      const data = await purchaseService.getShippingLabel(purchaseId)
      setLabel(data)
    } catch (error: any) {
      // If label doesn't exist, that's okay - we'll show a generate button
      if (error.response?.status !== 404) {
        toast.error('Failed to load shipping label')
      }
    } finally {
      setLoading(false)
    }
  }

  const handleGenerate = async () => {
    if (!purchaseId) return

    try {
      setGenerating(true)
      const data = await purchaseService.generateShippingLabel(purchaseId)
      setLabel(data)
      toast.success('Shipping label generated successfully!')
    } catch (error: any) {
      toast.error(error.response?.data?.error || 'Failed to generate shipping label')
    } finally {
      setGenerating(false)
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
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-4xl mx-auto px-4">
        <Button
          variant="outline"
          onClick={() => navigate('/purchases')}
          className="mb-6"
        >
          ← 購入履歴に戻る
        </Button>

        {label ? (
          <ShippingLabelView label={label} />
        ) : (
          <Card className="text-center py-12">
            <h2 className="text-2xl font-bold text-gray-900 mb-4">
              配送伝票を生成
            </h2>
            <p className="text-gray-600 mb-6">
              この購入の配送伝票がまだ生成されていません。<br />
              下のボタンをクリックして、配送伝票を生成してください。
            </p>
            <Button
              onClick={handleGenerate}
              disabled={generating}
              className="inline-flex items-center"
            >
              {generating ? (
                <>
                  <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2" />
                  生成中...
                </>
              ) : (
                '📦 配送伝票を生成する'
              )}
            </Button>
          </Card>
        )}
      </div>
    </div>
  )
}

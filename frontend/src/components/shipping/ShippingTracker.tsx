import { useState, useEffect } from 'react'
import { Card } from '@/components/common/Card'
import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

interface ShippingTracking {
  id: string
  purchase_id: string
  tracking_number: string
  carrier: string
  status: string
  shipped_at?: string
  delivered_at?: string
  estimated_arrival?: string
  shipping_method: string
  co2_saved: number
}

interface ShippingTrackerProps {
  purchaseId: string
}

export const ShippingTracker = ({ purchaseId }: ShippingTrackerProps) => {
  const [tracking, setTracking] = useState<ShippingTracking | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadTracking()
  }, [purchaseId])

  const loadTracking = async () => {
    try {
      setLoading(true)
      const token = localStorage.getItem('token')
      const response = await axios.get(
        `${API_BASE_URL}/v1/shipping/purchase/${purchaseId}`,
        {
          headers: { Authorization: `Bearer ${token}` },
        }
      )
      setTracking(response.data)
    } catch (error) {
      console.error('Failed to load tracking:', error)
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return (
      <Card>
        <div className="flex items-center justify-center py-4">
          <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-primary-500" />
        </div>
      </Card>
    )
  }

  if (!tracking) {
    return (
      <Card>
        <p className="text-gray-500 text-sm">配送情報はまだ登録されていません</p>
      </Card>
    )
  }

  const statusSteps = [
    { key: 'pending', label: '準備中', icon: '○' },
    { key: 'shipped', label: '発送済み', icon: '○' },
    { key: 'in_transit', label: '配送中', icon: '○' },
    { key: 'delivered', label: '配達完了', icon: '✓' },
  ]

  const currentStepIndex = statusSteps.findIndex((step) => step.key === tracking.status)

  const carrierNames: Record<string, string> = {
    yamato: 'ヤマト運輸',
    sagawa: '佐川急便',
    yupack: 'ゆうパック',
  }

  return (
    <Card>
      <div className="space-y-4">
        <div className="flex items-center justify-between">
          <h3 className="font-semibold">配送状況</h3>
          {tracking.shipping_method === 'eco' && (
            <span className="px-2 py-1 bg-green-100 text-green-800 text-xs rounded-full">
              エコ配送 -{tracking.co2_saved.toFixed(2)}kg CO2
            </span>
          )}
        </div>

        <div className="grid grid-cols-2 gap-3 text-sm">
          <div>
            <p className="text-gray-500 text-xs">配送業者</p>
            <p className="font-medium">{carrierNames[tracking.carrier] || tracking.carrier}</p>
          </div>
          <div>
            <p className="text-gray-500 text-xs">追跡番号</p>
            <p className="font-medium font-mono text-xs">{tracking.tracking_number}</p>
          </div>
        </div>

        {/* Status Timeline */}
        <div className="relative pt-2">
          <div className="flex justify-between">
            {statusSteps.map((step, index) => (
              <div key={step.key} className="flex flex-col items-center flex-1">
                <div
                  className={`w-8 h-8 rounded-full flex items-center justify-center text-lg transition-colors ${
                    index <= currentStepIndex
                      ? 'bg-primary-500 text-white'
                      : 'bg-gray-200 text-gray-400'
                  }`}
                >
                  {step.icon}
                </div>
                <p
                  className={`text-xs mt-2 text-center ${
                    index <= currentStepIndex ? 'text-primary-600 font-medium' : 'text-gray-400'
                  }`}
                >
                  {step.label}
                </p>
              </div>
            ))}
          </div>

          {/* Progress Line */}
          <div className="absolute top-6 left-0 right-0 h-1 bg-gray-200" style={{ zIndex: -1 }}>
            <div
              className="h-full bg-primary-500 transition-all duration-500"
              style={{
                width: `${(currentStepIndex / (statusSteps.length - 1)) * 100}%`,
              }}
            />
          </div>
        </div>

        {/* Dates */}
        <div className="grid grid-cols-2 gap-3 pt-4 border-t text-sm">
          {tracking.shipped_at && (
            <div>
              <p className="text-gray-500 text-xs">発送日時</p>
              <p className="font-medium">
                {new Date(tracking.shipped_at).toLocaleString('ja-JP', {
                  month: 'short',
                  day: 'numeric',
                  hour: '2-digit',
                  minute: '2-digit',
                })}
              </p>
            </div>
          )}

          {tracking.estimated_arrival && !tracking.delivered_at && (
            <div>
              <p className="text-gray-500 text-xs">お届け予定</p>
              <p className="font-medium text-primary-600">
                {new Date(tracking.estimated_arrival).toLocaleString('ja-JP', {
                  month: 'short',
                  day: 'numeric',
                })}
              </p>
            </div>
          )}

          {tracking.delivered_at && (
            <div>
              <p className="text-gray-500 text-xs">配達完了日時</p>
              <p className="font-medium text-green-600">
                {new Date(tracking.delivered_at).toLocaleString('ja-JP', {
                  month: 'short',
                  day: 'numeric',
                  hour: '2-digit',
                  minute: '2-digit',
                })}
              </p>
            </div>
          )}
        </div>
      </div>
    </Card>
  )
}

import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { Header } from '@/components/layout/Header'
import { Card } from '@/components/common/Card'
import { Button } from '@/components/common/Button'
import { autoPurchaseService } from '@/services/autoPurchase'
import { AutoPurchaseWatch } from '@/types'
import toast from 'react-hot-toast'

export const AutoPurchaseWatches = () => {
  const navigate = useNavigate()
  const [watches, setWatches] = useState<AutoPurchaseWatch[]>([])
  const [loading, setLoading] = useState(true)
  const [filter, setFilter] = useState<'all' | 'active' | 'executed' | 'cancelled'>('all')

  useEffect(() => {
    loadWatches()
  }, [])

  const loadWatches = async () => {
    try {
      setLoading(true)
      const data = await autoPurchaseService.getUserWatches()
      setWatches(data)
    } catch (error) {
      toast.error('監視リストの読み込みに失敗しました')
    } finally {
      setLoading(false)
    }
  }

  const handleCancelWatch = async (watchId: string) => {
    if (!confirm('この自動購入監視をキャンセルしますか？')) {
      return
    }

    try {
      await autoPurchaseService.cancelWatch(watchId)
      toast.success('監視をキャンセルしました')
      loadWatches()
    } catch (error: any) {
      toast.error(error.response?.data?.error || 'キャンセルに失敗しました')
    }
  }

  const getStatusBadge = (status: string) => {
    const colors = {
      active: 'bg-blue-100 text-blue-800',
      executed: 'bg-green-100 text-green-800',
      cancelled: 'bg-gray-100 text-gray-800',
      expired: 'bg-red-100 text-red-800',
    }
    const labels = {
      active: '監視中',
      executed: '購入済み',
      cancelled: 'キャンセル',
      expired: '期限切れ',
    }
    return (
      <span className={`px-2 py-1 rounded-full text-xs font-medium ${colors[status as keyof typeof colors]}`}>
        {labels[status as keyof typeof labels]}
      </span>
    )
  }

  const getTimeSlotLabel = (slot?: string) => {
    switch (slot) {
      case 'morning':
        return '午前 (8:00-12:00)'
      case 'afternoon':
        return '午後 (12:00-18:00)'
      case 'evening':
        return '夜間 (18:00-21:00)'
      default:
        return '指定なし'
    }
  }

  const filteredWatches = watches.filter((watch) => {
    if (filter === 'all') return true
    return watch.status === filter
  })

  return (
    <div className="min-h-screen">
      <Header />

      <div className="max-w-7xl mx-auto px-4 py-8 bg-white/70 backdrop-blur-sm">
        <div className="mb-6">
          <h1 className="text-3xl font-bold text-gray-900 mb-2">
            🤖 自動購入監視リスト
          </h1>
          <p className="text-gray-600">
            設定した価格以下になったら自動で購入されます
          </p>
        </div>

        {/* Filter */}
        <div className="mb-6 flex gap-2">
          <Button
            variant={filter === 'all' ? 'primary' : 'outline'}
            onClick={() => setFilter('all')}
          >
            すべて
          </Button>
          <Button
            variant={filter === 'active' ? 'primary' : 'outline'}
            onClick={() => setFilter('active')}
          >
            監視中
          </Button>
          <Button
            variant={filter === 'executed' ? 'primary' : 'outline'}
            onClick={() => setFilter('executed')}
          >
            購入済み
          </Button>
          <Button
            variant={filter === 'cancelled' ? 'primary' : 'outline'}
            onClick={() => setFilter('cancelled')}
          >
            キャンセル
          </Button>
        </div>

        {loading ? (
          <div className="flex justify-center items-center h-64">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-500" />
          </div>
        ) : filteredWatches.length === 0 ? (
          <Card className="text-center py-12">
            <h3 className="text-xl font-semibold text-gray-900 mb-2">
              {filter === 'all' ? '監視中の商品はありません' : `${filter === 'active' ? '監視中' : filter === 'executed' ? '購入済み' : 'キャンセル'}の商品はありません`}
            </h3>
            <p className="text-gray-600 mb-4">
              商品詳細ページから自動購入を設定できます
            </p>
            <Button onClick={() => navigate('/')}>
              商品を探す
            </Button>
          </Card>
        ) : (
          <div className="space-y-4">
            {filteredWatches.map((watch) => (
              <Card key={watch.id} className={watch.status === 'active' ? 'border-blue-300 bg-blue-50/30' : ''}>
                <div className="flex items-start gap-4">
                  {/* Product Image */}
                  {watch.product && watch.product.images.length > 0 && (
                    <img
                      src={watch.product.images[0].image_url}
                      alt={watch.product.title}
                      className="w-24 h-24 object-cover rounded-md"
                    />
                  )}

                  {/* Details */}
                  <div className="flex-1">
                    <div className="flex items-center gap-3 mb-2">
                      <h3
                        className="text-lg font-semibold cursor-pointer hover:text-primary-600"
                        onClick={() => watch.product && navigate(`/products/${watch.product.id}`)}
                      >
                        {watch.product?.title || '商品が削除されました'}
                      </h3>
                      {getStatusBadge(watch.status)}
                    </div>

                    <div className="grid grid-cols-2 gap-2 text-sm mb-2">
                      <div>
                        <span className="text-gray-600">現在価格:</span>{' '}
                        <span className="font-bold">¥{watch.product?.price.toLocaleString() || '-'}</span>
                      </div>
                      <div>
                        <span className="text-gray-600">最大購入価格:</span>{' '}
                        <span className="font-bold text-purple-600">¥{watch.max_price.toLocaleString()}</span>
                      </div>
                      <div>
                        <span className="text-gray-600">配送希望日:</span>{' '}
                        {watch.delivery_date
                          ? new Date(watch.delivery_date).toLocaleDateString('ja-JP')
                          : '指定なし'}
                      </div>
                      <div>
                        <span className="text-gray-600">配送時間帯:</span>{' '}
                        {getTimeSlotLabel(watch.delivery_time_slot)}
                      </div>
                    </div>

                    <div className="text-sm text-gray-500 mb-2">
                      <p>配送先: {watch.shipping_address}</p>
                    </div>

                    {watch.status === 'executed' && watch.purchase_id && (
                      <Button
                        size="sm"
                        onClick={() => navigate('/purchases')}
                        className="mt-2"
                      >
                        購入履歴を見る
                      </Button>
                    )}

                    {watch.status === 'active' && (
                      <div className="flex gap-2 mt-2">
                        <Button
                          size="sm"
                          variant="outline"
                          onClick={() => handleCancelWatch(watch.id)}
                          className="text-red-600 border-red-300 hover:bg-red-50"
                        >
                          監視をキャンセル
                        </Button>
                      </div>
                    )}

                    <p className="text-xs text-gray-400 mt-2">
                      作成日時: {new Date(watch.created_at).toLocaleString('ja-JP')}
                      {watch.status === 'active' && (
                        <>
                          {' | '}
                          有効期限: {new Date(watch.expires_at).toLocaleDateString('ja-JP')}
                        </>
                      )}
                    </p>
                  </div>
                </div>
              </Card>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}

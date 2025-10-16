import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { purchaseService, Purchase } from '@/services/purchases'
import { Button } from '@/components/common/Button'
import { Card } from '@/components/common/Card'
import { Header } from '@/components/layout/Header'
import toast from 'react-hot-toast'

export const Purchases = () => {
  const navigate = useNavigate()
  const [purchases, setPurchases] = useState<Purchase[]>([])
  const [loading, setLoading] = useState(true)
  const [filter, setFilter] = useState<'all' | 'buyer' | 'seller'>('all')

  useEffect(() => {
    loadPurchases()
  }, [filter])

  const loadPurchases = async () => {
    try {
      setLoading(true)
      const role = filter === 'all' ? undefined : filter
      const data = await purchaseService.list(role)
      setPurchases(data.purchases)
    } catch (error) {
      toast.error('Failed to load purchases')
    } finally {
      setLoading(false)
    }
  }

  const handleComplete = async (id: string) => {
    try {
      await purchaseService.complete(id)
      toast.success('Purchase completed!')
      loadPurchases()
    } catch (error: any) {
      toast.error(error.response?.data?.error || 'Failed to complete purchase')
    }
  }

  const getStatusBadge = (status: string) => {
    const colors = {
      pending: 'bg-yellow-100 text-yellow-800',
      completed: 'bg-green-100 text-green-800',
      cancelled: 'bg-red-100 text-red-800',
    }
    return (
      <span className={`px-2 py-1 rounded-full text-xs font-medium ${colors[status as keyof typeof colors]}`}>
        {status}
      </span>
    )
  }

  return (
    <div className="min-h-screen">
      <Header />

      <div className="max-w-7xl mx-auto px-4 py-8 bg-white/70 backdrop-blur-sm">
        <div className="mb-6 flex gap-2">
          <Button
            variant={filter === 'all' ? 'primary' : 'outline'}
            onClick={() => setFilter('all')}
          >
            All
          </Button>
          <Button
            variant={filter === 'buyer' ? 'primary' : 'outline'}
            onClick={() => setFilter('buyer')}
          >
            As Buyer
          </Button>
          <Button
            variant={filter === 'seller' ? 'primary' : 'outline'}
            onClick={() => setFilter('seller')}
          >
            As Seller
          </Button>
        </div>

        {loading ? (
          <div className="flex justify-center items-center h-64">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-500" />
          </div>
        ) : purchases.length === 0 ? (
          <Card className="text-center py-12">
            <div className="text-6xl mb-4">📦</div>
            <h3 className="text-xl font-semibold text-gray-900 mb-2">No purchases yet</h3>
            <p className="text-gray-600">Start shopping to see your purchases here</p>
            <Button className="mt-4" onClick={() => navigate('/')}>
              Browse Products
            </Button>
          </Card>
        ) : (
          <div className="space-y-4">
            {purchases.map((purchase) => (
              <Card key={purchase.id}>
                <div className="flex items-start justify-between">
                  <div className="flex-1">
                    <div className="flex items-center gap-3 mb-2">
                      <h3 className="text-lg font-semibold">
                        {purchase.product?.title || 'Product'}
                      </h3>
                      {getStatusBadge(purchase.status)}
                    </div>
                    <p className="text-gray-600 mb-2">
                      Price: ¥{purchase.price.toLocaleString()}
                    </p>
                    <p className="text-sm text-gray-500">
                      CO2 Saved: {purchase.co2_saved_kg.toFixed(2)}kg
                    </p>
                    <p className="text-sm text-gray-500">
                      Payment: {purchase.payment_method}
                    </p>
                    <p className="text-sm text-gray-500">
                      Shipping: {purchase.shipping_address}
                    </p>
                    <p className="text-xs text-gray-400 mt-2">
                      {new Date(purchase.created_at).toLocaleDateString()}
                    </p>
                  </div>
                  <div className="ml-4">
                    {purchase.status === 'pending' && filter === 'seller' && (
                      <Button
                        size="sm"
                        onClick={() => handleComplete(purchase.id)}
                      >
                        Complete
                      </Button>
                    )}
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

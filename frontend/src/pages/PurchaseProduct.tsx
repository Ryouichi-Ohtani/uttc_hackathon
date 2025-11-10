import { useState, useEffect } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { purchaseService } from '@/services/purchases'
import { productService } from '@/services/products'
import { Button } from '@/components/common/Button'
import { Card } from '@/components/common/Card'
import toast from 'react-hot-toast'
import { Product } from '@/types'

export const PurchaseProduct = () => {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const [loading, setLoading] = useState(false)
  const [product, setProduct] = useState<Product | null>(null)
  const [formData, setFormData] = useState({
    shipping_address: '',
    payment_method: 'credit_card',
  })

  useEffect(() => {
    if (id) {
      loadProduct(id)
    }
  }, [id])

  const loadProduct = async (productId: string) => {
    try {
      const data = await productService.getById(productId)
      setProduct(data)
    } catch (error) {
      toast.error('Failed to load product')
      navigate('/')
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!id) return

    setLoading(true)
    try {
      await purchaseService.create({
        product_id: id,
        shipping_address: formData.shipping_address,
        payment_method: formData.payment_method,
      })
      toast.success('Purchase completed successfully!')
      navigate('/purchases')
    } catch (error: any) {
      toast.error(error.response?.data?.error || 'Failed to complete purchase')
    } finally {
      setLoading(false)
    }
  }

  if (!product) {
    return <div className="flex justify-center items-center h-screen">Loading...</div>
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-3xl mx-auto px-4">
        <Button
          variant="outline"
          onClick={() => navigate(`/products/${id}`)}
          className="mb-6"
        >
          ← Back to Product
        </Button>

        <Card>
          <h1 className="text-2xl font-bold mb-6">Complete Purchase</h1>

          <div className="mb-6 p-4 bg-gray-50 rounded-lg">
            <h2 className="font-semibold mb-2">{product.title}</h2>
            <p className="text-2xl font-bold text-primary-600">¥{product.price.toLocaleString()}</p>
            {product.co2_impact_kg && (
              <p className="text-sm text-gray-600 mt-2">
                🌱 CO2 Saved: {product.co2_impact_kg.toFixed(2)}kg
              </p>
            )}
          </div>

          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Shipping Address *
              </label>
              <textarea
                required
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                rows={3}
                value={formData.shipping_address}
                onChange={(e) =>
                  setFormData({ ...formData, shipping_address: e.target.value })
                }
                placeholder="Enter your shipping address"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Payment Method *
              </label>
              <select
                required
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                value={formData.payment_method}
                onChange={(e) =>
                  setFormData({ ...formData, payment_method: e.target.value })
                }
              >
                <option value="credit_card">Credit Card</option>
                <option value="bank_transfer">Bank Transfer</option>
                <option value="cash_on_delivery">Cash on Delivery</option>
              </select>
            </div>

            <div className="pt-4 border-t">
              <Button type="submit" className="w-full" disabled={loading}>
                {loading ? 'Processing...' : 'Complete Purchase'}
              </Button>
            </div>
          </form>
        </Card>
      </div>
    </div>
  )
}

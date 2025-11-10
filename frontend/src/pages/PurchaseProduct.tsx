import { useState, useEffect } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { purchaseService } from '@/services/purchases'
import { productService } from '@/services/products'
import { Button } from '@/components/common/Button'
import { Card } from '@/components/common/Card'
import { AddressConfirmationDialog } from '@/components/shipping/AddressConfirmationDialog'
import { useAuthStore } from '@/store/authStore'
import toast from 'react-hot-toast'
import { Product } from '@/types'

export const PurchaseProduct = () => {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const { user } = useAuthStore()
  const [loading, setLoading] = useState(false)
  const [product, setProduct] = useState<Product | null>(null)
  const [showAddressDialog, setShowAddressDialog] = useState(false)
  const [useRegisteredAddress, setUseRegisteredAddress] = useState(false)
  const [formData, setFormData] = useState({
    shipping_address: '',
    payment_method: 'credit_card',
    delivery_date: '',
    delivery_time_slot: 'anytime' as 'morning' | 'afternoon' | 'evening' | 'anytime',
    recipient_name: '',
    recipient_phone_number: '',
    recipient_postal_code: '',
    recipient_prefecture: '',
    recipient_city: '',
    recipient_address_line1: '',
    recipient_address_line2: '',
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

  const handleAddressConfirm = () => {
    if (user) {
      setUseRegisteredAddress(true)
      setFormData({
        ...formData,
        recipient_name: user.display_name || user.username,
        recipient_phone_number: user.phone_number || '',
        recipient_postal_code: user.postal_code || '',
        recipient_prefecture: user.prefecture || '',
        recipient_city: user.city || '',
        recipient_address_line1: user.address_line1 || '',
        recipient_address_line2: user.address_line2 || '',
      })
      setShowAddressDialog(false)
    }
  }

  const handleAddressCancel = () => {
    setUseRegisteredAddress(false)
    setShowAddressDialog(false)
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!id) return

    // Show address confirmation dialog if not already confirmed
    if (!useRegisteredAddress && !formData.recipient_name && user) {
      setShowAddressDialog(true)
      return
    }

    setLoading(true)
    try {
      await purchaseService.create({
        product_id: id,
        shipping_address: formData.shipping_address,
        payment_method: formData.payment_method,
        delivery_date: formData.delivery_date || undefined,
        delivery_time_slot: formData.delivery_time_slot,
        use_registered_address: useRegisteredAddress,
        recipient_name: formData.recipient_name,
        recipient_phone_number: formData.recipient_phone_number,
        recipient_postal_code: formData.recipient_postal_code,
        recipient_prefecture: formData.recipient_prefecture,
        recipient_city: formData.recipient_city,
        recipient_address_line1: formData.recipient_address_line1,
        recipient_address_line2: formData.recipient_address_line2,
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
            {/* Delivery Date & Time */}
            <div className="bg-blue-50 rounded-lg p-4">
              <h3 className="font-semibold text-blue-900 mb-3">配送情報</h3>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    希望受取日
                  </label>
                  <input
                    type="date"
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                    value={formData.delivery_date}
                    min={new Date().toISOString().split('T')[0]}
                    onChange={(e) =>
                      setFormData({ ...formData, delivery_date: e.target.value })
                    }
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    希望受取時間帯
                  </label>
                  <select
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                    value={formData.delivery_time_slot}
                    onChange={(e) =>
                      setFormData({ ...formData, delivery_time_slot: e.target.value as any })
                    }
                  >
                    <option value="anytime">指定なし</option>
                    <option value="morning">午前 (8:00-12:00)</option>
                    <option value="afternoon">午後 (12:00-18:00)</option>
                    <option value="evening">夜間 (18:00-21:00)</option>
                  </select>
                </div>
              </div>
            </div>

            {/* Recipient Address */}
            {!useRegisteredAddress && (
              <div className="bg-gray-50 rounded-lg p-4">
                <h3 className="font-semibold text-gray-900 mb-3">配送先情報</h3>

                <div className="space-y-3">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      宛名 *
                    </label>
                    <input
                      type="text"
                      required
                      className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                      value={formData.recipient_name}
                      onChange={(e) =>
                        setFormData({ ...formData, recipient_name: e.target.value })
                      }
                      placeholder="山田 太郎"
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      電話番号 *
                    </label>
                    <input
                      type="tel"
                      required
                      className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                      value={formData.recipient_phone_number}
                      onChange={(e) =>
                        setFormData({ ...formData, recipient_phone_number: e.target.value })
                      }
                      placeholder="090-1234-5678"
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      郵便番号 *
                    </label>
                    <input
                      type="text"
                      required
                      className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                      value={formData.recipient_postal_code}
                      onChange={(e) =>
                        setFormData({ ...formData, recipient_postal_code: e.target.value })
                      }
                      placeholder="123-4567"
                    />
                  </div>

                  <div className="grid grid-cols-2 gap-3">
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">
                        都道府県 *
                      </label>
                      <input
                        type="text"
                        required
                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                        value={formData.recipient_prefecture}
                        onChange={(e) =>
                          setFormData({ ...formData, recipient_prefecture: e.target.value })
                        }
                        placeholder="東京都"
                      />
                    </div>

                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">
                        市区町村 *
                      </label>
                      <input
                        type="text"
                        required
                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                        value={formData.recipient_city}
                        onChange={(e) =>
                          setFormData({ ...formData, recipient_city: e.target.value })
                        }
                        placeholder="渋谷区"
                      />
                    </div>
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      番地 *
                    </label>
                    <input
                      type="text"
                      required
                      className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                      value={formData.recipient_address_line1}
                      onChange={(e) =>
                        setFormData({ ...formData, recipient_address_line1: e.target.value })
                      }
                      placeholder="渋谷1-2-3"
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      建物名・部屋番号
                    </label>
                    <input
                      type="text"
                      className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                      value={formData.recipient_address_line2}
                      onChange={(e) =>
                        setFormData({ ...formData, recipient_address_line2: e.target.value })
                      }
                      placeholder="渋谷マンション 101号室"
                    />
                  </div>
                </div>
              </div>
            )}

            {useRegisteredAddress && (
              <div className="bg-green-50 rounded-lg p-4">
                <div className="flex items-start justify-between">
                  <div className="flex-1">
                    <h3 className="font-semibold text-green-900 mb-2">登録住所を使用</h3>
                    <div className="text-sm text-green-800">
                      <p>{formData.recipient_name}</p>
                      <p>〒{formData.recipient_postal_code}</p>
                      <p>{formData.recipient_prefecture}{formData.recipient_city}</p>
                      <p>{formData.recipient_address_line1} {formData.recipient_address_line2}</p>
                      <p>{formData.recipient_phone_number}</p>
                    </div>
                  </div>
                  <Button
                    type="button"
                    variant="outline"
                    size="sm"
                    onClick={() => setUseRegisteredAddress(false)}
                  >
                    変更
                  </Button>
                </div>
              </div>
            )}

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Shipping Address (備考)
              </label>
              <textarea
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                rows={2}
                value={formData.shipping_address}
                onChange={(e) =>
                  setFormData({ ...formData, shipping_address: e.target.value })
                }
                placeholder="その他の配送指示があれば入力してください"
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
                <option value="credit_card">クレジットカード</option>
                <option value="bank_transfer">銀行振込</option>
                <option value="cash_on_delivery">代金引換</option>
              </select>
            </div>

            <div className="pt-4 border-t">
              <Button type="submit" className="w-full" disabled={loading}>
                {loading ? '処理中...' : '購入を確定する'}
              </Button>
            </div>
          </form>

          {user && (
            <AddressConfirmationDialog
              user={user}
              onConfirm={handleAddressConfirm}
              onCancel={handleAddressCancel}
              isOpen={showAddressDialog}
            />
          )}
        </Card>
      </div>
    </div>
  )
}

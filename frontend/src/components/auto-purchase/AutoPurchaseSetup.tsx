import { useState } from 'react'
import { Card } from '@/components/common/Card'
import { Button } from '@/components/common/Button'
import { PaymentAuthorizationDialog } from './PaymentAuthorizationDialog'
import { AddressConfirmationDialog } from '@/components/shipping/AddressConfirmationDialog'
import { autoPurchaseService } from '@/services/autoPurchase'
import { Product, User, PaymentAuthorizationRequest } from '@/types'
import toast from 'react-hot-toast'

interface AutoPurchaseSetupProps {
  product: Product
  user: User
  onSuccess?: () => void
}

export const AutoPurchaseSetup = ({ product, user, onSuccess }: AutoPurchaseSetupProps) => {
  const [isExpanded, setIsExpanded] = useState(false)
  const [showPaymentDialog, setShowPaymentDialog] = useState(false)
  const [showAddressDialog, setShowAddressDialog] = useState(false)
  const [loading, setLoading] = useState(false)

  const [maxPrice, setMaxPrice] = useState(product.price)
  const [deliveryDate, setDeliveryDate] = useState('')
  const [deliveryTimeSlot, setDeliveryTimeSlot] = useState<'morning' | 'afternoon' | 'evening' | 'anytime'>('anytime')

  const [addressData, setAddressData] = useState({
    recipient_name: user.display_name || '',
    recipient_phone_number: user.phone_number || '',
    recipient_postal_code: user.postal_code || '',
    recipient_prefecture: user.prefecture || '',
    recipient_city: user.city || '',
    recipient_address_line1: user.address_line1 || '',
    recipient_address_line2: user.address_line2 || '',
  })

  const [paymentAuth, setPaymentAuth] = useState<{
    payment_method_id: string
    payment_auth_token: string
    authorized_amount: number
  } | null>(null)

  const hasAddress = Boolean(
    user.postal_code &&
    user.prefecture &&
    user.city &&
    user.address_line1
  )

  const handleAddressConfirm = () => {
    setAddressData({
      recipient_name: user.display_name || '',
      recipient_phone_number: user.phone_number || '',
      recipient_postal_code: user.postal_code || '',
      recipient_prefecture: user.prefecture || '',
      recipient_city: user.city || '',
      recipient_address_line1: user.address_line1 || '',
      recipient_address_line2: user.address_line2 || '',
    })
    setShowAddressDialog(false)
    toast.success('登録住所を設定しました')
  }

  const handlePaymentAuthorization = async (request: PaymentAuthorizationRequest) => {
    try {
      const response = await autoPurchaseService.authorizePayment(request)
      setPaymentAuth({
        payment_method_id: response.payment_method_id,
        payment_auth_token: response.payment_auth_token,
        authorized_amount: response.authorized_amount,
      })
      setShowPaymentDialog(false)
      toast.success('決済認証が完了しました')
    } catch (error: any) {
      throw error
    }
  }

  const handleCreateWatch = async () => {
    // Validation
    if (!paymentAuth) {
      toast.error('決済情報を認証してください')
      return
    }

    if (!addressData.recipient_postal_code || !addressData.recipient_prefecture) {
      toast.error('配送先住所を設定してください')
      return
    }

    if (maxPrice > product.price) {
      toast.error('最大購入価格は現在の価格以下に設定してください')
      return
    }

    try {
      setLoading(true)
      const shippingAddress = `${addressData.recipient_postal_code} ${addressData.recipient_prefecture}${addressData.recipient_city}${addressData.recipient_address_line1}${addressData.recipient_address_line2 ? ' ' + addressData.recipient_address_line2 : ''}`

      await autoPurchaseService.createWatch({
        product_id: product.id,
        max_price: maxPrice,
        payment_method_id: paymentAuth.payment_method_id,
        payment_auth_token: paymentAuth.payment_auth_token,
        authorized_amount: paymentAuth.authorized_amount,
        delivery_date: deliveryDate || undefined,
        delivery_time_slot: deliveryTimeSlot,
        shipping_address: shippingAddress,
        recipient_name: addressData.recipient_name,
        recipient_phone_number: addressData.recipient_phone_number,
        recipient_postal_code: addressData.recipient_postal_code,
        recipient_prefecture: addressData.recipient_prefecture,
        recipient_city: addressData.recipient_city,
        recipient_address_line1: addressData.recipient_address_line1,
        recipient_address_line2: addressData.recipient_address_line2,
      })

      toast.success('自動購入監視を開始しました！')
      setIsExpanded(false)
      onSuccess?.()
    } catch (error: any) {
      toast.error(error.response?.data?.error || '自動購入の設定に失敗しました')
    } finally {
      setLoading(false)
    }
  }

  const getTomorrowDate = () => {
    const tomorrow = new Date()
    tomorrow.setDate(tomorrow.getDate() + 1)
    return tomorrow.toISOString().split('T')[0]
  }

  if (!isExpanded) {
    return (
      <Card className="bg-gradient-to-r from-purple-50 to-blue-50 border-purple-200">
        <div className="flex items-center justify-between">
          <div className="flex-1">
            <h3 className="font-bold text-lg text-purple-900 mb-1">
              🤖 AI自動購入
            </h3>
            <p className="text-sm text-purple-700">
              価格が設定した金額以下になったら自動で購入します
            </p>
          </div>
          <Button
            onClick={() => setIsExpanded(true)}
            className="ml-4 bg-purple-600 hover:bg-purple-700"
          >
            設定する
          </Button>
        </div>
      </Card>
    )
  }

  return (
    <>
      <Card className="bg-gradient-to-r from-purple-50 to-blue-50 border-purple-200">
        <div className="space-y-6">
          <div>
            <h3 className="font-bold text-xl text-purple-900 mb-2">
              🤖 AI自動購入の設定
            </h3>
            <p className="text-sm text-purple-700">
              以下の条件を満たしたときに自動で購入されます
            </p>
          </div>

          {/* Max Price */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              最大購入価格 (現在: ¥{product.price.toLocaleString()})
            </label>
            <div className="flex items-center gap-2">
              <span className="text-lg">¥</span>
              <input
                type="number"
                value={maxPrice}
                onChange={(e) => setMaxPrice(parseInt(e.target.value))}
                max={product.price}
                min={1}
                className="flex-1 px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-purple-500 focus:border-transparent text-lg font-bold"
              />
            </div>
            <p className="text-xs text-gray-500 mt-1">
              この価格以下になったら自動購入されます
            </p>
          </div>

          {/* Delivery Date */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              配送希望日 (任意)
            </label>
            <input
              type="date"
              value={deliveryDate}
              onChange={(e) => setDeliveryDate(e.target.value)}
              min={getTomorrowDate()}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-purple-500 focus:border-transparent"
            />
          </div>

          {/* Delivery Time Slot */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              配送時間帯
            </label>
            <select
              value={deliveryTimeSlot}
              onChange={(e) => setDeliveryTimeSlot(e.target.value as any)}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-purple-500 focus:border-transparent"
            >
              <option value="anytime">指定なし</option>
              <option value="morning">午前 (8:00-12:00)</option>
              <option value="afternoon">午後 (12:00-18:00)</option>
              <option value="evening">夜間 (18:00-21:00)</option>
            </select>
          </div>

          {/* Address */}
          <div>
            <div className="flex items-center justify-between mb-2">
              <label className="block text-sm font-medium text-gray-700">
                配送先住所
              </label>
              {hasAddress && (
                <Button
                  size="sm"
                  variant="outline"
                  onClick={() => setShowAddressDialog(true)}
                  className="text-xs"
                >
                  登録住所を使用
                </Button>
              )}
            </div>
            {addressData.recipient_postal_code ? (
              <div className="bg-white p-3 rounded-md border border-gray-300 text-sm">
                <p className="font-medium">{addressData.recipient_name}</p>
                <p className="text-gray-600">
                  〒{addressData.recipient_postal_code}
                </p>
                <p className="text-gray-600">
                  {addressData.recipient_prefecture}
                  {addressData.recipient_city}
                  {addressData.recipient_address_line1}
                  {addressData.recipient_address_line2 && ` ${addressData.recipient_address_line2}`}
                </p>
                <p className="text-gray-600">{addressData.recipient_phone_number}</p>
              </div>
            ) : (
              <p className="text-sm text-gray-500 italic">
                住所が設定されていません
              </p>
            )}
          </div>

          {/* Payment Authorization */}
          <div>
            <div className="flex items-center justify-between mb-2">
              <label className="block text-sm font-medium text-gray-700">
                決済情報
              </label>
              {!paymentAuth && (
                <Button
                  size="sm"
                  onClick={() => setShowPaymentDialog(true)}
                  className="bg-green-600 hover:bg-green-700"
                >
                  決済を認証
                </Button>
              )}
            </div>
            {paymentAuth ? (
              <div className="bg-green-50 border border-green-200 p-3 rounded-md">
                <p className="text-sm text-green-800">
                  ✓ 決済認証完了 (最大金額: ¥{paymentAuth.authorized_amount.toLocaleString()})
                </p>
              </div>
            ) : (
              <div className="bg-yellow-50 border border-yellow-200 p-3 rounded-md">
                <p className="text-sm text-yellow-800">
                  ⚠️ 決済情報の認証が必要です
                </p>
              </div>
            )}
          </div>

          {/* Action Buttons */}
          <div className="flex gap-3 pt-2">
            <Button
              variant="outline"
              onClick={() => setIsExpanded(false)}
              disabled={loading}
              className="flex-1"
            >
              キャンセル
            </Button>
            <Button
              onClick={handleCreateWatch}
              disabled={loading || !paymentAuth}
              className="flex-1 bg-purple-600 hover:bg-purple-700"
            >
              {loading ? (
                <>
                  <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2" />
                  設定中...
                </>
              ) : (
                '自動購入を開始'
              )}
            </Button>
          </div>
        </div>
      </Card>

      {/* Payment Authorization Dialog */}
      <PaymentAuthorizationDialog
        isOpen={showPaymentDialog}
        amount={maxPrice}
        onConfirm={handlePaymentAuthorization}
        onCancel={() => setShowPaymentDialog(false)}
      />

      {/* Address Confirmation Dialog */}
      <AddressConfirmationDialog
        isOpen={showAddressDialog}
        user={user}
        onConfirm={handleAddressConfirm}
        onCancel={() => setShowAddressDialog(false)}
      />
    </>
  )
}

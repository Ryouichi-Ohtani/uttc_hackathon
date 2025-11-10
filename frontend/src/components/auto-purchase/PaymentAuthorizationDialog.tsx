import { useState } from 'react'
import { Card } from '@/components/common/Card'
import { Button } from '@/components/common/Button'
import { PaymentAuthorizationRequest } from '@/types'

interface PaymentAuthorizationDialogProps {
  isOpen: boolean
  amount: number
  onConfirm: (request: PaymentAuthorizationRequest) => Promise<void>
  onCancel: () => void
}

export const PaymentAuthorizationDialog = ({
  isOpen,
  amount,
  onConfirm,
  onCancel,
}: PaymentAuthorizationDialogProps) => {
  const [formData, setFormData] = useState({
    card_number: '',
    expiry_month: '',
    expiry_year: '',
    cvv: '',
    cardholder_name: '',
  })
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  if (!isOpen) return null

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target

    // Format card number with spaces
    if (name === 'card_number') {
      const cleaned = value.replace(/\s/g, '')
      const formatted = cleaned.replace(/(\d{4})/g, '$1 ').trim()
      setFormData((prev) => ({ ...prev, [name]: formatted }))
      return
    }

    setFormData((prev) => ({ ...prev, [name]: value }))
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')

    // Validation
    const cardNumber = formData.card_number.replace(/\s/g, '')
    if (cardNumber.length !== 16) {
      setError('カード番号は16桁で入力してください')
      return
    }

    const expiryMonth = parseInt(formData.expiry_month)
    const expiryYear = parseInt(formData.expiry_year)

    if (!expiryMonth || expiryMonth < 1 || expiryMonth > 12) {
      setError('有効期限(月)が正しくありません')
      return
    }

    if (!expiryYear || expiryYear < new Date().getFullYear() % 100) {
      setError('有効期限(年)が正しくありません')
      return
    }

    if (formData.cvv.length < 3 || formData.cvv.length > 4) {
      setError('セキュリティコードは3-4桁で入力してください')
      return
    }

    if (!formData.cardholder_name.trim()) {
      setError('カード名義人を入力してください')
      return
    }

    try {
      setLoading(true)
      await onConfirm({
        card_number: cardNumber,
        expiry_month: expiryMonth,
        expiry_year: 2000 + expiryYear,
        cvv: formData.cvv,
        cardholder_name: formData.cardholder_name,
        amount,
      })
    } catch (err: any) {
      setError(err.response?.data?.error || '決済認証に失敗しました')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50">
      <Card className="w-full max-w-md mx-4">
        <div className="mb-6">
          <h2 className="text-2xl font-bold text-gray-900 mb-2">
            決済情報の認証
          </h2>
          <p className="text-gray-600 text-sm">
            自動購入を有効にするため、決済情報を事前認証します。<br />
            最大購入金額: <span className="font-bold text-lg">¥{amount.toLocaleString()}</span>
          </p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          {/* Card Number */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              カード番号
            </label>
            <input
              type="text"
              name="card_number"
              value={formData.card_number}
              onChange={handleChange}
              placeholder="1234 5678 9012 3456"
              maxLength={19}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-primary-500 focus:border-transparent font-mono"
              required
            />
          </div>

          {/* Expiry Date */}
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                有効期限(月)
              </label>
              <input
                type="number"
                name="expiry_month"
                value={formData.expiry_month}
                onChange={handleChange}
                placeholder="MM"
                min="1"
                max="12"
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                有効期限(年)
              </label>
              <input
                type="number"
                name="expiry_year"
                value={formData.expiry_year}
                onChange={handleChange}
                placeholder="YY"
                min={new Date().getFullYear() % 100}
                max={99}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                required
              />
            </div>
          </div>

          {/* CVV */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              セキュリティコード(CVV)
            </label>
            <input
              type="password"
              name="cvv"
              value={formData.cvv}
              onChange={handleChange}
              placeholder="123"
              maxLength={4}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-primary-500 focus:border-transparent font-mono"
              required
            />
          </div>

          {/* Cardholder Name */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              カード名義人
            </label>
            <input
              type="text"
              name="cardholder_name"
              value={formData.cardholder_name}
              onChange={handleChange}
              placeholder="TARO YAMADA"
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-primary-500 focus:border-transparent uppercase"
              required
            />
          </div>

          {/* Error Message */}
          {error && (
            <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-md text-sm">
              {error}
            </div>
          )}

          {/* Notice */}
          <div className="bg-blue-50 border border-blue-200 text-blue-800 px-4 py-3 rounded-md text-xs">
            <p className="font-semibold mb-1">ℹ️ ご注意</p>
            <ul className="list-disc list-inside space-y-1">
              <li>即座に請求されることはありません</li>
              <li>条件に一致した場合のみ自動購入されます</li>
              <li>認証は30日間有効です</li>
            </ul>
          </div>

          {/* Buttons */}
          <div className="flex gap-3 pt-2">
            <Button
              type="button"
              variant="outline"
              onClick={onCancel}
              disabled={loading}
              className="flex-1"
            >
              キャンセル
            </Button>
            <Button
              type="submit"
              disabled={loading}
              className="flex-1"
            >
              {loading ? (
                <>
                  <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2" />
                  認証中...
                </>
              ) : (
                '認証する'
              )}
            </Button>
          </div>
        </form>
      </Card>
    </div>
  )
}

import { ShippingLabel } from '@/types'
import { Card } from '@/components/common/Card'

interface ShippingLabelViewProps {
  label: ShippingLabel
}

export const ShippingLabelView = ({ label }: ShippingLabelViewProps) => {
  const formatDate = (dateString?: string) => {
    if (!dateString) return '指定なし'
    const date = new Date(dateString)
    return date.toLocaleDateString('ja-JP', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      weekday: 'short'
    })
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

  const getCarrierName = (carrier: string) => {
    switch (carrier) {
      case 'yamato':
        return 'ヤマト運輸'
      case 'sagawa':
        return '佐川急便'
      case 'japan_post':
        return '日本郵便'
      default:
        return carrier
    }
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="bg-gradient-to-r from-blue-600 to-blue-700 text-white p-6 rounded-lg">
        <h2 className="text-2xl font-bold mb-2">配送伝票</h2>
        <p className="text-blue-100">出荷時にご使用ください</p>
      </div>

      {/* Label Content - Mimics actual shipping label */}
      <Card className="print:shadow-none">
        <div className="border-4 border-black p-6">
          {/* Carrier Info */}
          <div className="flex items-center justify-between mb-6 pb-4 border-b-2 border-black">
            <div className="text-2xl font-bold">{getCarrierName(label.carrier)}</div>
            <div className="text-right">
              <div className="text-xs text-gray-600">伝票番号</div>
              <div className="text-lg font-mono font-bold">
                {label.tracking_number || '未発行'}
              </div>
            </div>
          </div>

          {/* Delivery Info */}
          <div className="grid grid-cols-2 gap-6 mb-6">
            {/* Delivery Date/Time */}
            <div className="bg-yellow-50 border-2 border-yellow-400 p-4">
              <div className="font-bold text-sm mb-2 text-yellow-900">配送希望日時</div>
              <div className="space-y-1">
                <div className="text-lg font-bold">{formatDate(label.delivery_date)}</div>
                <div className="text-sm">{getTimeSlotLabel(label.delivery_time_slot)}</div>
              </div>
            </div>

            {/* Package Info */}
            <div className="bg-blue-50 border-2 border-blue-400 p-4">
              <div className="font-bold text-sm mb-2 text-blue-900">荷物情報</div>
              <div className="space-y-1">
                <div>サイズ: <span className="font-bold text-lg">{label.package_size}サイズ</span></div>
                <div>重量: <span className="font-bold">{label.weight}kg</span></div>
              </div>
            </div>
          </div>

          {/* Recipient (TO) */}
          <div className="mb-6 border-2 border-black p-4">
            <div className="bg-black text-white px-2 py-1 inline-block mb-3 font-bold">
              お届け先
            </div>
            <div className="grid grid-cols-3 gap-4">
              <div className="col-span-1">
                <div className="text-xs text-gray-600 mb-1">郵便番号</div>
                <div className="text-2xl font-mono font-bold">
                  〒{label.recipient_postal_code}
                </div>
              </div>
              <div className="col-span-2">
                <div className="text-xs text-gray-600 mb-1">住所</div>
                <div className="font-bold text-lg leading-tight">
                  {label.recipient_prefecture}{label.recipient_city}
                </div>
                <div className="font-bold text-lg leading-tight">
                  {label.recipient_address_line1}
                </div>
                {label.recipient_address_line2 && (
                  <div className="font-bold text-lg leading-tight">
                    {label.recipient_address_line2}
                  </div>
                )}
              </div>
            </div>
            <div className="mt-3 pt-3 border-t border-gray-300">
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <div className="text-xs text-gray-600 mb-1">お名前</div>
                  <div className="text-xl font-bold">{label.recipient_name} 様</div>
                </div>
                <div>
                  <div className="text-xs text-gray-600 mb-1">電話番号</div>
                  <div className="text-lg font-mono font-bold">{label.recipient_phone_number}</div>
                </div>
              </div>
            </div>
          </div>

          {/* Sender (FROM) */}
          <div className="border-2 border-gray-400 p-4 bg-gray-50">
            <div className="bg-gray-600 text-white px-2 py-1 inline-block mb-3 font-bold text-sm">
              ご依頼主
            </div>
            <div className="grid grid-cols-3 gap-4 text-sm">
              <div className="col-span-1">
                <div className="text-xs text-gray-600 mb-1">郵便番号</div>
                <div className="font-mono font-bold">
                  〒{label.sender_postal_code}
                </div>
              </div>
              <div className="col-span-2">
                <div className="text-xs text-gray-600 mb-1">住所</div>
                <div className="font-bold">
                  {label.sender_prefecture}{label.sender_city}
                  {label.sender_address_line1}
                  {label.sender_address_line2 && ` ${label.sender_address_line2}`}
                </div>
              </div>
            </div>
            <div className="mt-2 pt-2 border-t border-gray-300">
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <div className="text-xs text-gray-600 mb-1">お名前</div>
                  <div className="font-bold">{label.sender_name}</div>
                </div>
                <div>
                  <div className="text-xs text-gray-600 mb-1">電話番号</div>
                  <div className="font-mono font-bold">{label.sender_phone_number}</div>
                </div>
              </div>
            </div>
          </div>

          {/* Product Name */}
          <div className="mt-4 p-3 bg-gray-100 border border-gray-300">
            <div className="text-xs text-gray-600 mb-1">品名</div>
            <div className="font-bold">{label.product_name}</div>
          </div>
        </div>
      </Card>

      {/* Print Button */}
      <div className="flex justify-center print:hidden">
        <button
          onClick={() => window.print()}
          className="bg-blue-600 hover:bg-blue-700 text-white font-bold py-3 px-8 rounded-lg shadow-lg transition"
        >
          🖨️ この伝票を印刷する
        </button>
      </div>

      {/* Instructions */}
      <Card className="bg-yellow-50 border-yellow-300 print:hidden">
        <h3 className="font-bold text-yellow-900 mb-2">📦 発送の手順</h3>
        <ol className="list-decimal list-inside space-y-1 text-sm text-yellow-900">
          <li>この伝票を印刷して、荷物に貼り付けてください</li>
          <li>商品を丁寧に梱包してください</li>
          <li>{getCarrierName(label.carrier)}の営業所またはコンビニから発送してください</li>
          <li>発送が完了したら、購入者に発送完了の連絡をしてください</li>
        </ol>
      </Card>
    </div>
  )
}

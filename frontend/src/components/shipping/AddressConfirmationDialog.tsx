import { User } from '@/types'
import { Button } from '@/components/common/Button'
import { Card } from '@/components/common/Card'

interface AddressConfirmationDialogProps {
  user: User
  onConfirm: () => void
  onCancel: () => void
  isOpen: boolean
}

export const AddressConfirmationDialog = ({
  user,
  onConfirm,
  onCancel,
  isOpen
}: AddressConfirmationDialogProps) => {
  if (!isOpen) return null

  const hasAddress = user.postal_code && user.prefecture && user.city && user.address_line1

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50">
      <Card className="max-w-md w-full mx-4">
        <h2 className="text-xl font-bold mb-4">配送先住所の確認</h2>

        {hasAddress ? (
          <>
            <p className="text-gray-600 mb-4">
              以下の登録住所を配送先として使用しますか？
            </p>

            <div className="bg-gray-50 rounded-lg p-4 mb-6">
              <div className="space-y-2">
                <div>
                  <span className="text-sm text-gray-500">宛名：</span>
                  <p className="font-medium">{user.display_name || user.username}</p>
                </div>
                <div>
                  <span className="text-sm text-gray-500">郵便番号：</span>
                  <p className="font-medium">〒{user.postal_code}</p>
                </div>
                <div>
                  <span className="text-sm text-gray-500">住所：</span>
                  <p className="font-medium">
                    {user.prefecture}{user.city}
                    {user.address_line1}
                    {user.address_line2 && ` ${user.address_line2}`}
                  </p>
                </div>
                {user.phone_number && (
                  <div>
                    <span className="text-sm text-gray-500">電話番号：</span>
                    <p className="font-medium">{user.phone_number}</p>
                  </div>
                )}
              </div>
            </div>

            <div className="flex gap-3">
              <Button
                variant="primary"
                className="flex-1"
                onClick={onConfirm}
              >
                この住所を使用
              </Button>
              <Button
                variant="outline"
                className="flex-1"
                onClick={onCancel}
              >
                別の住所を入力
              </Button>
            </div>
          </>
        ) : (
          <>
            <p className="text-gray-600 mb-6">
              登録されている住所情報がありません。配送先住所を入力してください。
            </p>
            <Button
              variant="primary"
              className="w-full"
              onClick={onCancel}
            >
              住所を入力する
            </Button>
          </>
        )}
      </Card>
    </div>
  )
}

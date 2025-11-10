import { useState, useEffect } from 'react'
import { Header } from '@/components/layout/Header'
import { Card } from '@/components/common/Card'
import { Button } from '@/components/common/Button'
import { adminService } from '@/services/admin'
import { User, Product, Purchase } from '@/types'
import toast from 'react-hot-toast'

type TabType = 'users' | 'products' | 'purchases'

export const AdminDashboard = () => {
  const [activeTab, setActiveTab] = useState<TabType>('users')
  const [loading, setLoading] = useState(true)

  // Users state
  const [users, setUsers] = useState<User[]>([])
  const [usersTotal, setUsersTotal] = useState(0)

  // Products state
  const [products, setProducts] = useState<Product[]>([])
  const [productsTotal, setProductsTotal] = useState(0)

  // Purchases state
  const [purchases, setPurchases] = useState<Purchase[]>([])
  const [purchasesTotal, setPurchasesTotal] = useState(0)

  // Edit modal state
  const [editingUser, setEditingUser] = useState<User | null>(null)
  const [editingProduct, setEditingProduct] = useState<Product | null>(null)
  const [editingPurchase, setEditingPurchase] = useState<Purchase | null>(null)

  useEffect(() => {
    loadData()
  }, [activeTab])

  const loadData = async () => {
    setLoading(true)
    try {
      if (activeTab === 'users') {
        await loadUsers()
      } else if (activeTab === 'products') {
        await loadProducts()
      } else if (activeTab === 'purchases') {
        await loadPurchases()
      }
    } finally {
      setLoading(false)
    }
  }

  const loadUsers = async () => {
    try {
      const data = await adminService.getUsers()
      setUsers(data.users)
      setUsersTotal(data.pagination.total)
    } catch (error) {
      toast.error('ユーザーの読み込みに失敗しました')
    }
  }

  const loadProducts = async () => {
    try {
      const data = await adminService.getProducts()
      setProducts(data.products)
      setProductsTotal(data.pagination.total)
    } catch (error) {
      toast.error('商品の読み込みに失敗しました')
    }
  }

  const loadPurchases = async () => {
    try {
      const data = await adminService.getPurchases()
      setPurchases(data.purchases)
      setPurchasesTotal(data.pagination.total)
    } catch (error) {
      toast.error('購入履歴の読み込みに失敗しました')
    }
  }

  const handleDeleteUser = async (userId: string) => {
    if (!confirm('このユーザーを削除しますか？')) return

    try {
      await adminService.deleteUser(userId)
      toast.success('ユーザーを削除しました')
      loadUsers()
    } catch (error: any) {
      toast.error(error.response?.data?.error?.message || 'ユーザーの削除に失敗しました')
    }
  }

  const handleDeleteProduct = async (productId: string) => {
    if (!confirm('この商品を削除しますか？')) return

    try {
      await adminService.deleteProduct(productId)
      toast.success('商品を削除しました')
      loadProducts()
    } catch (error: any) {
      toast.error(error.response?.data?.error?.message || '商品の削除に失敗しました')
    }
  }

  const handleUpdateUser = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!editingUser) return

    try {
      await adminService.updateUser(editingUser.id, {
        display_name: editingUser.display_name,
        role: editingUser.role,
        bio: editingUser.bio,
      })
      toast.success('ユーザー情報を更新しました')
      setEditingUser(null)
      loadUsers()
    } catch (error: any) {
      toast.error(error.response?.data?.error?.message || '更新に失敗しました')
    }
  }

  const handleUpdateProduct = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!editingProduct) return

    try {
      await adminService.updateProduct(editingProduct.id, {
        title: editingProduct.title,
        description: editingProduct.description,
        price: editingProduct.price,
        status: editingProduct.status,
        condition: editingProduct.condition,
      })
      toast.success('商品情報を更新しました')
      setEditingProduct(null)
      loadProducts()
    } catch (error: any) {
      toast.error(error.response?.data?.error?.message || '更新に失敗しました')
    }
  }

  const handleUpdatePurchase = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!editingPurchase) return

    try {
      await adminService.updatePurchase(editingPurchase.id, {
        status: editingPurchase.status,
      })
      toast.success('購入情報を更新しました')
      setEditingPurchase(null)
      loadPurchases()
    } catch (error: any) {
      toast.error(error.response?.data?.error?.message || '更新に失敗しました')
    }
  }

  const getStatusBadge = (status: string, type: 'product' | 'purchase') => {
    if (type === 'product') {
      const colors = {
        active: 'bg-green-100 text-green-800',
        sold: 'bg-blue-100 text-blue-800',
        draft: 'bg-gray-100 text-gray-800',
        reserved: 'bg-yellow-100 text-yellow-800',
        deleted: 'bg-red-100 text-red-800',
      }
      const labels = {
        active: '販売中',
        sold: '売却済み',
        draft: '下書き',
        reserved: '予約済み',
        deleted: '削除済み',
      }
      return (
        <span className={`px-2 py-1 rounded-full text-xs font-medium ${colors[status as keyof typeof colors]}`}>
          {labels[status as keyof typeof labels] || status}
        </span>
      )
    } else {
      const colors = {
        pending: 'bg-yellow-100 text-yellow-800',
        completed: 'bg-green-100 text-green-800',
        cancelled: 'bg-red-100 text-red-800',
      }
      const labels = {
        pending: '処理中',
        completed: '完了',
        cancelled: 'キャンセル',
      }
      return (
        <span className={`px-2 py-1 rounded-full text-xs font-medium ${colors[status as keyof typeof colors]}`}>
          {labels[status as keyof typeof labels] || status}
        </span>
      )
    }
  }

  const getRoleBadge = (role: string) => {
    const colors = {
      admin: 'bg-purple-100 text-purple-800',
      moderator: 'bg-blue-100 text-blue-800',
      user: 'bg-gray-100 text-gray-800',
    }
    const labels = {
      admin: '管理者',
      moderator: 'モデレーター',
      user: 'ユーザー',
    }
    return (
      <span className={`px-2 py-1 rounded-full text-xs font-medium ${colors[role as keyof typeof colors]}`}>
        {labels[role as keyof typeof labels] || role}
      </span>
    )
  }

  return (
    <div className="min-h-screen">
      <Header />

      <div className="max-w-7xl mx-auto px-4 py-8 bg-white/70 backdrop-blur-sm">
        <div className="mb-6">
          <h1 className="text-3xl font-bold text-gray-900 mb-2">
            管理者ダッシュボード
          </h1>
          <p className="text-gray-600">
            ユーザー、商品、購入履歴を管理します
          </p>
        </div>

        {/* Tabs */}
        <div className="mb-6 flex gap-2 border-b border-gray-200">
          <button
            className={`px-4 py-2 font-medium border-b-2 transition-colors ${
              activeTab === 'users'
                ? 'border-primary-500 text-primary-600'
                : 'border-transparent text-gray-600 hover:text-gray-900'
            }`}
            onClick={() => setActiveTab('users')}
          >
            ユーザー ({usersTotal})
          </button>
          <button
            className={`px-4 py-2 font-medium border-b-2 transition-colors ${
              activeTab === 'products'
                ? 'border-primary-500 text-primary-600'
                : 'border-transparent text-gray-600 hover:text-gray-900'
            }`}
            onClick={() => setActiveTab('products')}
          >
            商品 ({productsTotal})
          </button>
          <button
            className={`px-4 py-2 font-medium border-b-2 transition-colors ${
              activeTab === 'purchases'
                ? 'border-primary-500 text-primary-600'
                : 'border-transparent text-gray-600 hover:text-gray-900'
            }`}
            onClick={() => setActiveTab('purchases')}
          >
            購入履歴 ({purchasesTotal})
          </button>
        </div>

        {loading ? (
          <div className="flex justify-center items-center h-64">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-500" />
          </div>
        ) : (
          <>
            {/* Users Tab */}
            {activeTab === 'users' && (
              <Card>
                <div className="overflow-x-auto">
                  <table className="min-w-full divide-y divide-gray-200">
                    <thead className="bg-gray-50">
                      <tr>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                          ユーザー
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                          メール
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                          ロール
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                          登録日
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                          操作
                        </th>
                      </tr>
                    </thead>
                    <tbody className="bg-white divide-y divide-gray-200">
                      {users.map((user) => (
                        <tr key={user.id} className="hover:bg-gray-50">
                          <td className="px-6 py-4 whitespace-nowrap">
                            <div className="flex items-center">
                              {user.avatar_url && (
                                <img
                                  src={user.avatar_url}
                                  alt={user.display_name}
                                  className="w-8 h-8 rounded-full mr-3"
                                />
                              )}
                              <div>
                                <div className="text-sm font-medium text-gray-900">
                                  {user.display_name}
                                </div>
                                <div className="text-sm text-gray-500">
                                  @{user.username}
                                </div>
                              </div>
                            </div>
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                            {user.email}
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap">
                            {getRoleBadge(user.role)}
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                            {new Date(user.created_at).toLocaleDateString('ja-JP')}
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                            <Button
                              size="sm"
                              variant="outline"
                              onClick={() => setEditingUser(user)}
                              className="mr-2"
                            >
                              編集
                            </Button>
                            <Button
                              size="sm"
                              variant="outline"
                              onClick={() => handleDeleteUser(user.id)}
                              className="text-red-600 border-red-300 hover:bg-red-50"
                            >
                              削除
                            </Button>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </Card>
            )}

            {/* Products Tab */}
            {activeTab === 'products' && (
              <Card>
                <div className="overflow-x-auto">
                  <table className="min-w-full divide-y divide-gray-200">
                    <thead className="bg-gray-50">
                      <tr>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                          商品
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                          販売者
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                          価格
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                          ステータス
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                          操作
                        </th>
                      </tr>
                    </thead>
                    <tbody className="bg-white divide-y divide-gray-200">
                      {products.map((product) => (
                        <tr key={product.id} className="hover:bg-gray-50">
                          <td className="px-6 py-4 whitespace-nowrap">
                            <div className="flex items-center">
                              {product.images.length > 0 && (
                                <img
                                  src={product.images[0].image_url}
                                  alt={product.title}
                                  className="w-12 h-12 object-cover rounded mr-3"
                                />
                              )}
                              <div className="text-sm font-medium text-gray-900">
                                {product.title}
                              </div>
                            </div>
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                            {product.seller?.display_name || '-'}
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                            ¥{product.price.toLocaleString()}
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap">
                            {getStatusBadge(product.status, 'product')}
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                            <Button
                              size="sm"
                              variant="outline"
                              onClick={() => setEditingProduct(product)}
                              className="mr-2"
                            >
                              編集
                            </Button>
                            <Button
                              size="sm"
                              variant="outline"
                              onClick={() => handleDeleteProduct(product.id)}
                              className="text-red-600 border-red-300 hover:bg-red-50"
                            >
                              削除
                            </Button>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </Card>
            )}

            {/* Purchases Tab */}
            {activeTab === 'purchases' && (
              <Card>
                <div className="overflow-x-auto">
                  <table className="min-w-full divide-y divide-gray-200">
                    <thead className="bg-gray-50">
                      <tr>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                          商品
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                          購入者
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                          販売者
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                          金額
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                          ステータス
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                          操作
                        </th>
                      </tr>
                    </thead>
                    <tbody className="bg-white divide-y divide-gray-200">
                      {purchases.map((purchase) => (
                        <tr key={purchase.id} className="hover:bg-gray-50">
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                            {purchase.product?.title || '-'}
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                            {purchase.buyer?.display_name || '-'}
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                            {purchase.seller?.display_name || '-'}
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                            ¥{purchase.price.toLocaleString()}
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap">
                            {getStatusBadge(purchase.status, 'purchase')}
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                            <Button
                              size="sm"
                              variant="outline"
                              onClick={() => setEditingPurchase(purchase)}
                            >
                              編集
                            </Button>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </Card>
            )}
          </>
        )}
      </div>

      {/* Edit User Modal */}
      {editingUser && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50">
          <Card className="w-full max-w-md mx-4">
            <h2 className="text-2xl font-bold text-gray-900 mb-4">
              ユーザー編集
            </h2>
            <form onSubmit={handleUpdateUser} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  表示名
                </label>
                <input
                  type="text"
                  value={editingUser.display_name}
                  onChange={(e) =>
                    setEditingUser({ ...editingUser, display_name: e.target.value })
                  }
                  className="w-full px-3 py-2 border border-gray-300 rounded-md"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  ロール
                </label>
                <select
                  value={editingUser.role}
                  onChange={(e) =>
                    setEditingUser({ ...editingUser, role: e.target.value })
                  }
                  className="w-full px-3 py-2 border border-gray-300 rounded-md"
                >
                  <option value="user">ユーザー</option>
                  <option value="moderator">モデレーター</option>
                  <option value="admin">管理者</option>
                </select>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  自己紹介
                </label>
                <textarea
                  value={editingUser.bio || ''}
                  onChange={(e) =>
                    setEditingUser({ ...editingUser, bio: e.target.value })
                  }
                  rows={3}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md"
                />
              </div>

              <div className="flex gap-3 pt-2">
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => setEditingUser(null)}
                  className="flex-1"
                >
                  キャンセル
                </Button>
                <Button type="submit" className="flex-1">
                  更新
                </Button>
              </div>
            </form>
          </Card>
        </div>
      )}

      {/* Edit Product Modal */}
      {editingProduct && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50">
          <Card className="w-full max-w-md mx-4">
            <h2 className="text-2xl font-bold text-gray-900 mb-4">
              商品編集
            </h2>
            <form onSubmit={handleUpdateProduct} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  商品名
                </label>
                <input
                  type="text"
                  value={editingProduct.title}
                  onChange={(e) =>
                    setEditingProduct({ ...editingProduct, title: e.target.value })
                  }
                  className="w-full px-3 py-2 border border-gray-300 rounded-md"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  価格
                </label>
                <input
                  type="number"
                  value={editingProduct.price}
                  onChange={(e) =>
                    setEditingProduct({
                      ...editingProduct,
                      price: parseInt(e.target.value),
                    })
                  }
                  className="w-full px-3 py-2 border border-gray-300 rounded-md"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  ステータス
                </label>
                <select
                  value={editingProduct.status}
                  onChange={(e) =>
                    setEditingProduct({
                      ...editingProduct,
                      status: e.target.value as any,
                    })
                  }
                  className="w-full px-3 py-2 border border-gray-300 rounded-md"
                >
                  <option value="draft">下書き</option>
                  <option value="active">販売中</option>
                  <option value="sold">売却済み</option>
                  <option value="reserved">予約済み</option>
                  <option value="deleted">削除済み</option>
                </select>
              </div>

              <div className="flex gap-3 pt-2">
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => setEditingProduct(null)}
                  className="flex-1"
                >
                  キャンセル
                </Button>
                <Button type="submit" className="flex-1">
                  更新
                </Button>
              </div>
            </form>
          </Card>
        </div>
      )}

      {/* Edit Purchase Modal */}
      {editingPurchase && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50">
          <Card className="w-full max-w-md mx-4">
            <h2 className="text-2xl font-bold text-gray-900 mb-4">
              購入編集
            </h2>
            <form onSubmit={handleUpdatePurchase} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  ステータス
                </label>
                <select
                  value={editingPurchase.status}
                  onChange={(e) =>
                    setEditingPurchase({
                      ...editingPurchase,
                      status: e.target.value as any,
                    })
                  }
                  className="w-full px-3 py-2 border border-gray-300 rounded-md"
                >
                  <option value="pending">処理中</option>
                  <option value="completed">完了</option>
                  <option value="cancelled">キャンセル</option>
                </select>
              </div>

              <div className="flex gap-3 pt-2">
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => setEditingPurchase(null)}
                  className="flex-1"
                >
                  キャンセル
                </Button>
                <Button type="submit" className="flex-1">
                  更新
                </Button>
              </div>
            </form>
          </Card>
        </div>
      )}
    </div>
  )
}

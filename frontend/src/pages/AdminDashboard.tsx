import { useState, useEffect } from 'react'
import { Header } from '@/components/layout/Header'
import { Button } from '@/components/common/Button'
import { LoadingSpinner } from '@/components/common/LoadingSpinner'
import { User, Product, Purchase } from '@/types'
import {
  useAdminUsers,
  useAdminProducts,
  useAdminPurchases,
  useUpdateUser,
  useDeleteUser,
  useAdminUpdateProduct,
  useAdminDeleteProduct,
  useAdminUpdatePurchase,
} from '@/hooks/useAdmin'
import {
  UsersIcon,
  ShoppingBagIcon,
  CurrencyYenIcon,
  ChartBarIcon,
  ArrowTrendingUpIcon,
  ArrowTrendingDownIcon,
  XMarkIcon,
  PencilIcon,
  TrashIcon,
  FunnelIcon,
  MagnifyingGlassIcon
} from '@heroicons/react/24/outline'
import { CheckCircleIcon, XCircleIcon, ClockIcon } from '@heroicons/react/24/solid'

type TabType = 'users' | 'products' | 'purchases'

interface Stats {
  totalUsers: number
  totalProducts: number
  totalPurchases: number
  totalRevenue: number
  userGrowth: number
  productGrowth: number
  purchaseGrowth: number
  revenueGrowth: number
}

export const AdminDashboard = () => {
  const [activeTab, setActiveTab] = useState<TabType>('users')
  const [searchTerm, setSearchTerm] = useState('')
  const [stats, setStats] = useState<Stats>({
    totalUsers: 0,
    totalProducts: 0,
    totalPurchases: 0,
    totalRevenue: 0,
    userGrowth: 12.5,
    productGrowth: 8.3,
    purchaseGrowth: 15.7,
    revenueGrowth: 23.4
  })

  // Edit modal state
  const [editingUser, setEditingUser] = useState<User | null>(null)
  const [editingProduct, setEditingProduct] = useState<Product | null>(null)
  const [editingPurchase, setEditingPurchase] = useState<Purchase | null>(null)

  // React Query hooks
  const { data: usersData, isLoading: usersLoading } = useAdminUsers()
  const { data: productsData, isLoading: productsLoading } = useAdminProducts()
  const { data: purchasesData, isLoading: purchasesLoading } = useAdminPurchases()

  const updateUserMutation = useUpdateUser()
  const deleteUserMutation = useDeleteUser()
  const updateProductMutation = useAdminUpdateProduct()
  const deleteProductMutation = useAdminDeleteProduct()
  const updatePurchaseMutation = useAdminUpdatePurchase()

  // Extract data from queries
  const users = usersData?.users || []
  const usersTotal = usersData?.pagination?.total || 0
  const products = productsData?.products || []
  const productsTotal = productsData?.pagination?.total || 0
  const purchases = purchasesData?.purchases || []
  const purchasesTotal = purchasesData?.pagination?.total || 0

  // Calculate revenue
  const totalRevenue = purchases.reduce((sum, p) => sum + p.price, 0)

  // Determine loading state based on active tab
  const loading =
    (activeTab === 'users' && usersLoading) ||
    (activeTab === 'products' && productsLoading) ||
    (activeTab === 'purchases' && purchasesLoading)

  // Update stats when data changes
  useEffect(() => {
    setStats(prev => ({
      ...prev,
      totalUsers: usersTotal,
      totalProducts: productsTotal,
      totalPurchases: purchasesTotal,
      totalRevenue,
    }))
  }, [usersTotal, productsTotal, purchasesTotal, totalRevenue])

  const handleDeleteUser = (userId: string) => {
    if (!confirm('このユーザーを削除しますか？この操作は取り消せません。')) return
    deleteUserMutation.mutate(userId)
  }

  const handleDeleteProduct = (productId: string) => {
    if (!confirm('この商品を削除しますか？この操作は取り消せません。')) return
    deleteProductMutation.mutate(productId)
  }

  const handleUpdateUser = (e: React.FormEvent) => {
    e.preventDefault()
    if (!editingUser) return

    updateUserMutation.mutate(
      {
        userId: editingUser.id,
        data: {
          display_name: editingUser.display_name,
          role: editingUser.role,
          bio: editingUser.bio,
        },
      },
      {
        onSuccess: () => setEditingUser(null),
      }
    )
  }

  const handleUpdateProduct = (e: React.FormEvent) => {
    e.preventDefault()
    if (!editingProduct) return

    updateProductMutation.mutate(
      {
        productId: editingProduct.id,
        data: {
          title: editingProduct.title,
          description: editingProduct.description,
          price: editingProduct.price,
          status: editingProduct.status,
          condition: editingProduct.condition,
        },
      },
      {
        onSuccess: () => setEditingProduct(null),
      }
    )
  }

  const handleUpdatePurchase = (e: React.FormEvent) => {
    e.preventDefault()
    if (!editingPurchase) return

    updatePurchaseMutation.mutate(
      {
        purchaseId: editingPurchase.id,
        data: {
          status: editingPurchase.status,
        },
      },
      {
        onSuccess: () => setEditingPurchase(null),
      }
    )
  }

  const getStatusBadge = (status: string, type: 'product' | 'purchase') => {
    if (type === 'product') {
      const configs = {
        active: { icon: CheckCircleIcon, class: 'bg-emerald-100 dark:bg-emerald-900/30 text-emerald-700 dark:text-emerald-400', label: '販売中' },
        sold: { icon: CheckCircleIcon, class: 'bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-400', label: '売却済み' },
        draft: { icon: ClockIcon, class: 'bg-slate-100 dark:bg-slate-800 text-slate-700 dark:text-slate-400', label: '下書き' },
        reserved: { icon: ClockIcon, class: 'bg-amber-100 dark:bg-amber-900/30 text-amber-700 dark:text-amber-400', label: '予約済み' },
        deleted: { icon: XCircleIcon, class: 'bg-red-100 dark:bg-red-900/30 text-red-700 dark:text-red-400', label: '削除済み' },
      }
      const config = configs[status as keyof typeof configs] || configs.draft
      const Icon = config.icon
      return (
        <span className={`inline-flex items-center gap-1.5 px-3 py-1.5 rounded-full text-xs font-bold ${config.class}`}>
          <Icon className="w-3.5 h-3.5" />
          {config.label}
        </span>
      )
    } else {
      const configs = {
        pending: { icon: ClockIcon, class: 'bg-amber-100 dark:bg-amber-900/30 text-amber-700 dark:text-amber-400', label: '処理中' },
        completed: { icon: CheckCircleIcon, class: 'bg-emerald-100 dark:bg-emerald-900/30 text-emerald-700 dark:text-emerald-400', label: '完了' },
        cancelled: { icon: XCircleIcon, class: 'bg-red-100 dark:bg-red-900/30 text-red-700 dark:text-red-400', label: 'キャンセル' },
      }
      const config = configs[status as keyof typeof configs] || configs.pending
      const Icon = config.icon
      return (
        <span className={`inline-flex items-center gap-1.5 px-3 py-1.5 rounded-full text-xs font-bold ${config.class}`}>
          <Icon className="w-3.5 h-3.5" />
          {config.label}
        </span>
      )
    }
  }

  const getRoleBadge = (role: string) => {
    const configs = {
      admin: { class: 'bg-gradient-to-r from-purple-500 to-pink-500 text-white shadow-md', label: '管理者' },
      moderator: { class: 'bg-gradient-to-r from-blue-500 to-cyan-500 text-white shadow-md', label: 'モデレーター' },
      user: { class: 'bg-slate-200 dark:bg-slate-700 text-slate-700 dark:text-slate-300', label: 'ユーザー' },
    }
    const config = configs[role as keyof typeof configs] || configs.user
    return (
      <span className={`inline-flex items-center px-3 py-1.5 rounded-full text-xs font-bold ${config.class}`}>
        {config.label}
      </span>
    )
  }

  const filteredUsers = users.filter(user =>
    user.display_name.toLowerCase().includes(searchTerm.toLowerCase()) ||
    user.email.toLowerCase().includes(searchTerm.toLowerCase()) ||
    user.username.toLowerCase().includes(searchTerm.toLowerCase())
  )

  const filteredProducts = products.filter(product =>
    product.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
    product.category.toLowerCase().includes(searchTerm.toLowerCase())
  )

  const filteredPurchases = purchases.filter(purchase =>
    purchase.product?.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
    purchase.buyer?.display_name.toLowerCase().includes(searchTerm.toLowerCase())
  )

  const StatCard = ({ icon: Icon, title, value, change, color }: any) => (
    <div className="card p-6 bg-gradient-to-br from-white to-slate-50/50 dark:from-dark-card dark:to-slate-900/50 hover:shadow-2xl transition-all duration-300 hover:-translate-y-1">
      <div className="flex items-start justify-between mb-4">
        <div className={`p-3 rounded-xl bg-gradient-to-br ${color} shadow-lg`}>
          <Icon className="w-6 h-6 text-white" />
        </div>
        <div className={`flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-bold ${
          change >= 0
            ? 'bg-emerald-100 dark:bg-emerald-900/30 text-emerald-700 dark:text-emerald-400'
            : 'bg-red-100 dark:bg-red-900/30 text-red-700 dark:text-red-400'
        }`}>
          {change >= 0 ? <ArrowTrendingUpIcon className="w-3.5 h-3.5" /> : <ArrowTrendingDownIcon className="w-3.5 h-3.5" />}
          {Math.abs(change)}%
        </div>
      </div>
      <h3 className="text-slate-600 dark:text-slate-400 text-sm font-medium mb-1">{title}</h3>
      <p className="text-3xl font-bold text-slate-900 dark:text-white">{value}</p>
    </div>
  )

  return (
    <div className="min-h-screen bg-slate-50 dark:bg-dark">
      <div className="mesh-gradient" />
      <Header />

      <div className="container-max relative z-10 py-8">
        {/* Page Header */}
        <div className="mb-8 animate-fade-up">
          <div className="flex items-center gap-3 mb-3">
            <div className="p-3 bg-gradient-to-br from-purple-500 to-pink-500 rounded-xl shadow-lg">
              <ChartBarIcon className="w-8 h-8 text-white" />
            </div>
            <div>
              <h1 className="text-4xl font-bold gradient-text">管理者ダッシュボード</h1>
              <p className="text-slate-600 dark:text-slate-400 mt-1">
                システム全体の管理と監視
              </p>
            </div>
          </div>
        </div>

        {/* Stats Cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8 animate-fade-up" style={{ animationDelay: '100ms' }}>
          <StatCard
            icon={UsersIcon}
            title="総ユーザー数"
            value={stats.totalUsers.toLocaleString()}
            change={stats.userGrowth}
            color="from-blue-500 to-cyan-600"
          />
          <StatCard
            icon={ShoppingBagIcon}
            title="総商品数"
            value={stats.totalProducts.toLocaleString()}
            change={stats.productGrowth}
            color="from-purple-500 to-pink-600"
          />
          <StatCard
            icon={ChartBarIcon}
            title="総取引数"
            value={stats.totalPurchases.toLocaleString()}
            change={stats.purchaseGrowth}
            color="from-emerald-500 to-green-600"
          />
          <StatCard
            icon={CurrencyYenIcon}
            title="総売上"
            value={`¥${stats.totalRevenue.toLocaleString()}`}
            change={stats.revenueGrowth}
            color="from-amber-500 to-orange-600"
          />
        </div>

        {/* Tabs */}
        <div className="mb-6 animate-fade-up" style={{ animationDelay: '200ms' }}>
          <div className="card p-1 bg-white dark:bg-dark-card inline-flex gap-1">
            {[
              { id: 'users', label: 'ユーザー', count: usersTotal, icon: UsersIcon },
              { id: 'products', label: '商品', count: productsTotal, icon: ShoppingBagIcon },
              { id: 'purchases', label: '購入履歴', count: purchasesTotal, icon: ChartBarIcon }
            ].map(tab => (
              <button
                key={tab.id}
                className={`flex items-center gap-2 px-6 py-3 font-bold rounded-lg transition-all duration-200 ${
                  activeTab === tab.id
                    ? 'bg-gradient-to-r from-primary-500 to-accent-500 text-white shadow-lg'
                    : 'text-slate-600 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800'
                }`}
                onClick={() => setActiveTab(tab.id as TabType)}
              >
                <tab.icon className="w-5 h-5" />
                {tab.label}
                <span className={`px-2 py-0.5 rounded-full text-xs font-bold ${
                  activeTab === tab.id
                    ? 'bg-white/20'
                    : 'bg-slate-200 dark:bg-slate-700'
                }`}>
                  {tab.count}
                </span>
              </button>
            ))}
          </div>
        </div>

        {/* Search & Filters */}
        <div className="card p-4 mb-6 animate-fade-up" style={{ animationDelay: '300ms' }}>
          <div className="flex flex-col sm:flex-row gap-4">
            <div className="flex-1 relative">
              <MagnifyingGlassIcon className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-slate-400" />
              <input
                type="text"
                placeholder={`${activeTab === 'users' ? 'ユーザー名、メールアドレス' : activeTab === 'products' ? '商品名、カテゴリー' : '商品名、購入者'}で検索...`}
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="w-full pl-12 pr-4 py-3 bg-slate-50 dark:bg-slate-800 border-2 border-transparent rounded-xl focus:border-primary-500 focus:bg-white dark:focus:bg-slate-900 transition-colors font-medium"
              />
            </div>
            <Button variant="outline" className="flex items-center gap-2">
              <FunnelIcon className="w-5 h-5" />
              フィルター
            </Button>
          </div>
        </div>

        {/* Content */}
        {loading ? (
          <div className="card p-12 flex flex-col items-center justify-center">
            <LoadingSpinner
              type="spinner"
              size="xl"
              text="データを読み込んでいます..."
            />
          </div>
        ) : (
          <div className="animate-fade-up" style={{ animationDelay: '400ms' }}>
            {/* Users Tab */}
            {activeTab === 'users' && (
              <div className="card overflow-hidden">
                <div className="overflow-x-auto">
                  <table className="min-w-full">
                    <thead>
                      <tr className="bg-slate-50 dark:bg-slate-800/50 border-b border-slate-200 dark:border-slate-700">
                        <th className="px-6 py-4 text-left text-xs font-bold text-slate-700 dark:text-slate-300 uppercase tracking-wider">
                          ユーザー
                        </th>
                        <th className="px-6 py-4 text-left text-xs font-bold text-slate-700 dark:text-slate-300 uppercase tracking-wider">
                          メールアドレス
                        </th>
                        <th className="px-6 py-4 text-left text-xs font-bold text-slate-700 dark:text-slate-300 uppercase tracking-wider">
                          ロール
                        </th>
                        <th className="px-6 py-4 text-left text-xs font-bold text-slate-700 dark:text-slate-300 uppercase tracking-wider">
                          登録日
                        </th>
                        <th className="px-6 py-4 text-right text-xs font-bold text-slate-700 dark:text-slate-300 uppercase tracking-wider">
                          操作
                        </th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-slate-200 dark:divide-slate-700">
                      {filteredUsers.map((user, index) => (
                        <tr
                          key={user.id}
                          className="hover:bg-slate-50 dark:hover:bg-slate-800/50 transition-colors animate-fade-up"
                          style={{ animationDelay: `${index * 30}ms` }}
                        >
                          <td className="px-6 py-4">
                            <div className="flex items-center gap-3">
                              {user.avatar_url ? (
                                <img
                                  src={user.avatar_url}
                                  alt={user.display_name}
                                  className="w-10 h-10 rounded-full object-cover ring-2 ring-slate-200 dark:ring-slate-700"
                                />
                              ) : (
                                <div className="w-10 h-10 rounded-full bg-gradient-to-br from-primary-500 to-accent-500 flex items-center justify-center text-white font-bold ring-2 ring-slate-200 dark:ring-slate-700">
                                  {user.display_name[0].toUpperCase()}
                                </div>
                              )}
                              <div>
                                <div className="font-bold text-slate-900 dark:text-white">
                                  {user.display_name}
                                </div>
                                <div className="text-sm text-slate-500 dark:text-slate-400">
                                  @{user.username}
                                </div>
                              </div>
                            </div>
                          </td>
                          <td className="px-6 py-4 text-slate-700 dark:text-slate-300 font-medium">
                            {user.email}
                          </td>
                          <td className="px-6 py-4">
                            {getRoleBadge(user.role)}
                          </td>
                          <td className="px-6 py-4 text-slate-600 dark:text-slate-400 font-medium">
                            {new Date(user.created_at).toLocaleDateString('ja-JP')}
                          </td>
                          <td className="px-6 py-4">
                            <div className="flex items-center justify-end gap-2">
                              <button
                                onClick={() => setEditingUser(user)}
                                className="p-2 hover:bg-blue-50 dark:hover:bg-blue-900/20 text-blue-600 dark:text-blue-400 rounded-lg transition-colors"
                                title="編集"
                              >
                                <PencilIcon className="w-5 h-5" />
                              </button>
                              <button
                                onClick={() => handleDeleteUser(user.id)}
                                className="p-2 hover:bg-red-50 dark:hover:bg-red-900/20 text-red-600 dark:text-red-400 rounded-lg transition-colors"
                                title="削除"
                              >
                                <TrashIcon className="w-5 h-5" />
                              </button>
                            </div>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>
            )}

            {/* Products Tab */}
            {activeTab === 'products' && (
              <div className="card overflow-hidden">
                <div className="overflow-x-auto">
                  <table className="min-w-full">
                    <thead>
                      <tr className="bg-slate-50 dark:bg-slate-800/50 border-b border-slate-200 dark:border-slate-700">
                        <th className="px-6 py-4 text-left text-xs font-bold text-slate-700 dark:text-slate-300 uppercase tracking-wider">
                          商品
                        </th>
                        <th className="px-6 py-4 text-left text-xs font-bold text-slate-700 dark:text-slate-300 uppercase tracking-wider">
                          販売者
                        </th>
                        <th className="px-6 py-4 text-left text-xs font-bold text-slate-700 dark:text-slate-300 uppercase tracking-wider">
                          価格
                        </th>
                        <th className="px-6 py-4 text-left text-xs font-bold text-slate-700 dark:text-slate-300 uppercase tracking-wider">
                          ステータス
                        </th>
                        <th className="px-6 py-4 text-right text-xs font-bold text-slate-700 dark:text-slate-300 uppercase tracking-wider">
                          操作
                        </th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-slate-200 dark:divide-slate-700">
                      {filteredProducts.map((product, index) => (
                        <tr
                          key={product.id}
                          className="hover:bg-slate-50 dark:hover:bg-slate-800/50 transition-colors animate-fade-up"
                          style={{ animationDelay: `${index * 30}ms` }}
                        >
                          <td className="px-6 py-4">
                            <div className="flex items-center gap-3">
                              {product.images.length > 0 && (
                                <img
                                  src={product.images[0].image_url}
                                  alt={product.title}
                                  className="w-14 h-14 object-cover rounded-lg ring-2 ring-slate-200 dark:ring-slate-700"
                                />
                              )}
                              <div className="font-bold text-slate-900 dark:text-white max-w-xs truncate">
                                {product.title}
                              </div>
                            </div>
                          </td>
                          <td className="px-6 py-4 text-slate-700 dark:text-slate-300 font-medium">
                            {product.seller?.display_name || '-'}
                          </td>
                          <td className="px-6 py-4">
                            <span className="text-lg font-bold text-primary-600 dark:text-primary-400">
                              ¥{product.price.toLocaleString()}
                            </span>
                          </td>
                          <td className="px-6 py-4">
                            {getStatusBadge(product.status, 'product')}
                          </td>
                          <td className="px-6 py-4">
                            <div className="flex items-center justify-end gap-2">
                              <button
                                onClick={() => setEditingProduct(product)}
                                className="p-2 hover:bg-blue-50 dark:hover:bg-blue-900/20 text-blue-600 dark:text-blue-400 rounded-lg transition-colors"
                                title="編集"
                              >
                                <PencilIcon className="w-5 h-5" />
                              </button>
                              <button
                                onClick={() => handleDeleteProduct(product.id)}
                                className="p-2 hover:bg-red-50 dark:hover:bg-red-900/20 text-red-600 dark:text-red-400 rounded-lg transition-colors"
                                title="削除"
                              >
                                <TrashIcon className="w-5 h-5" />
                              </button>
                            </div>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>
            )}

            {/* Purchases Tab */}
            {activeTab === 'purchases' && (
              <div className="card overflow-hidden">
                <div className="overflow-x-auto">
                  <table className="min-w-full">
                    <thead>
                      <tr className="bg-slate-50 dark:bg-slate-800/50 border-b border-slate-200 dark:border-slate-700">
                        <th className="px-6 py-4 text-left text-xs font-bold text-slate-700 dark:text-slate-300 uppercase tracking-wider">
                          商品
                        </th>
                        <th className="px-6 py-4 text-left text-xs font-bold text-slate-700 dark:text-slate-300 uppercase tracking-wider">
                          購入者
                        </th>
                        <th className="px-6 py-4 text-left text-xs font-bold text-slate-700 dark:text-slate-300 uppercase tracking-wider">
                          販売者
                        </th>
                        <th className="px-6 py-4 text-left text-xs font-bold text-slate-700 dark:text-slate-300 uppercase tracking-wider">
                          金額
                        </th>
                        <th className="px-6 py-4 text-left text-xs font-bold text-slate-700 dark:text-slate-300 uppercase tracking-wider">
                          ステータス
                        </th>
                        <th className="px-6 py-4 text-right text-xs font-bold text-slate-700 dark:text-slate-300 uppercase tracking-wider">
                          操作
                        </th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-slate-200 dark:divide-slate-700">
                      {filteredPurchases.map((purchase, index) => (
                        <tr
                          key={purchase.id}
                          className="hover:bg-slate-50 dark:hover:bg-slate-800/50 transition-colors animate-fade-up"
                          style={{ animationDelay: `${index * 30}ms` }}
                        >
                          <td className="px-6 py-4 font-bold text-slate-900 dark:text-white">
                            {purchase.product?.title || '-'}
                          </td>
                          <td className="px-6 py-4 text-slate-700 dark:text-slate-300 font-medium">
                            {purchase.buyer?.display_name || '-'}
                          </td>
                          <td className="px-6 py-4 text-slate-700 dark:text-slate-300 font-medium">
                            {purchase.seller?.display_name || '-'}
                          </td>
                          <td className="px-6 py-4">
                            <span className="text-lg font-bold text-emerald-600 dark:text-emerald-400">
                              ¥{purchase.price.toLocaleString()}
                            </span>
                          </td>
                          <td className="px-6 py-4">
                            {getStatusBadge(purchase.status, 'purchase')}
                          </td>
                          <td className="px-6 py-4">
                            <div className="flex items-center justify-end gap-2">
                              <button
                                onClick={() => setEditingPurchase(purchase)}
                                className="p-2 hover:bg-blue-50 dark:hover:bg-blue-900/20 text-blue-600 dark:text-blue-400 rounded-lg transition-colors"
                                title="編集"
                              >
                                <PencilIcon className="w-5 h-5" />
                              </button>
                            </div>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>
            )}
          </div>
        )}
      </div>

      {/* Edit User Modal */}
      {editingUser && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm animate-fade-up">
          <div className="card w-full max-w-lg p-6 shadow-2xl animate-scale">
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-2xl font-bold gradient-text">ユーザー編集</h2>
              <button
                onClick={() => setEditingUser(null)}
                className="p-2 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-lg transition-colors"
              >
                <XMarkIcon className="w-6 h-6" />
              </button>
            </div>
            <form onSubmit={handleUpdateUser} className="space-y-5">
              <div>
                <label className="label">表示名</label>
                <input
                  type="text"
                  value={editingUser.display_name}
                  onChange={(e) =>
                    setEditingUser({ ...editingUser, display_name: e.target.value })
                  }
                  className="input"
                />
              </div>

              <div>
                <label className="label">ロール</label>
                <select
                  value={editingUser.role}
                  onChange={(e) =>
                    setEditingUser({ ...editingUser, role: e.target.value })
                  }
                  className="input"
                >
                  <option value="user">ユーザー</option>
                  <option value="moderator">モデレーター</option>
                  <option value="admin">管理者</option>
                </select>
              </div>

              <div>
                <label className="label">自己紹介</label>
                <textarea
                  value={editingUser.bio || ''}
                  onChange={(e) =>
                    setEditingUser({ ...editingUser, bio: e.target.value })
                  }
                  rows={3}
                  className="input"
                />
              </div>

              <div className="flex gap-3 pt-4">
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => setEditingUser(null)}
                  className="flex-1 btn-ripple"
                >
                  キャンセル
                </Button>
                <Button type="submit" className="flex-1 btn-gradient btn-ripple">
                  更新
                </Button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Edit Product Modal */}
      {editingProduct && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm animate-fade-up">
          <div className="card w-full max-w-lg p-6 shadow-2xl animate-scale max-h-[90vh] overflow-y-auto">
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-2xl font-bold gradient-text">商品編集</h2>
              <button
                onClick={() => setEditingProduct(null)}
                className="p-2 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-lg transition-colors"
              >
                <XMarkIcon className="w-6 h-6" />
              </button>
            </div>
            <form onSubmit={handleUpdateProduct} className="space-y-5">
              <div>
                <label className="label">商品名</label>
                <input
                  type="text"
                  value={editingProduct.title}
                  onChange={(e) =>
                    setEditingProduct({ ...editingProduct, title: e.target.value })
                  }
                  className="input"
                />
              </div>

              <div>
                <label className="label">説明</label>
                <textarea
                  value={editingProduct.description}
                  onChange={(e) =>
                    setEditingProduct({ ...editingProduct, description: e.target.value })
                  }
                  rows={3}
                  className="input"
                />
              </div>

              <div>
                <label className="label">価格</label>
                <input
                  type="number"
                  value={editingProduct.price}
                  onChange={(e) =>
                    setEditingProduct({
                      ...editingProduct,
                      price: parseInt(e.target.value),
                    })
                  }
                  className="input"
                />
              </div>

              <div>
                <label className="label">状態</label>
                <select
                  value={editingProduct.condition}
                  onChange={(e) =>
                    setEditingProduct({
                      ...editingProduct,
                      condition: e.target.value as any,
                    })
                  }
                  className="input"
                >
                  <option value="new">新品</option>
                  <option value="like_new">未使用に近い</option>
                  <option value="good">良好</option>
                  <option value="fair">使用感あり</option>
                </select>
              </div>

              <div>
                <label className="label">ステータス</label>
                <select
                  value={editingProduct.status}
                  onChange={(e) =>
                    setEditingProduct({
                      ...editingProduct,
                      status: e.target.value as any,
                    })
                  }
                  className="input"
                >
                  <option value="draft">下書き</option>
                  <option value="active">販売中</option>
                  <option value="sold">売却済み</option>
                  <option value="reserved">予約済み</option>
                  <option value="deleted">削除済み</option>
                </select>
              </div>

              <div className="flex gap-3 pt-4">
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => setEditingProduct(null)}
                  className="flex-1 btn-ripple"
                >
                  キャンセル
                </Button>
                <Button type="submit" className="flex-1 btn-gradient btn-ripple">
                  更新
                </Button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Edit Purchase Modal */}
      {editingPurchase && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm animate-fade-up">
          <div className="card w-full max-w-md p-6 shadow-2xl animate-scale">
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-2xl font-bold gradient-text">購入編集</h2>
              <button
                onClick={() => setEditingPurchase(null)}
                className="p-2 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-lg transition-colors"
              >
                <XMarkIcon className="w-6 h-6" />
              </button>
            </div>
            <form onSubmit={handleUpdatePurchase} className="space-y-5">
              <div>
                <label className="label">ステータス</label>
                <select
                  value={editingPurchase.status}
                  onChange={(e) =>
                    setEditingPurchase({
                      ...editingPurchase,
                      status: e.target.value as any,
                    })
                  }
                  className="input"
                >
                  <option value="pending">処理中</option>
                  <option value="completed">完了</option>
                  <option value="cancelled">キャンセル</option>
                </select>
              </div>

              <div className="flex gap-3 pt-4">
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => setEditingPurchase(null)}
                  className="flex-1 btn-ripple"
                >
                  キャンセル
                </Button>
                <Button type="submit" className="flex-1 btn-gradient btn-ripple">
                  更新
                </Button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}

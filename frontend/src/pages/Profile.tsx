import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuthStore } from '@/store/authStore'
import { Button } from '@/components/common/Button'
import { Card } from '@/components/common/Card'
import { Header } from '@/components/layout/Header'
import api from '@/services/api'
import { Product } from '@/types'
import { offerService, Offer, MarketPriceAnalysis } from '@/services/offers'
import toast from 'react-hot-toast'
import { PROFILE_PLACEHOLDER } from '@/utils/placeholderImages'

export const Profile = () => {
  const navigate = useNavigate()
  const { user, logout } = useAuthStore()
  const [myProducts, setMyProducts] = useState<Product[]>([])
  const [loading, setLoading] = useState(true)
  const [filter, setFilter] = useState<'all' | 'active' | 'draft' | 'sold'>('all')
  const [offers, setOffers] = useState<Offer[]>([])
  const [activeTab, setActiveTab] = useState<'products' | 'offers'>('products')
  const [marketAnalysis, setMarketAnalysis] = useState<{ [offerId: string]: MarketPriceAnalysis }>({})
  const [loadingAnalysis, setLoadingAnalysis] = useState<{ [offerId: string]: boolean }>({})

  useEffect(() => {
    loadMyProducts()
    if (activeTab === 'offers') {
      loadMyOffers()
    }
  }, [filter, activeTab])

  const loadMyProducts = async () => {
    try {
      setLoading(true)
      const response = await api.get('/products', {
        params: {
          seller_id: user?.id,
          status: filter === 'all' ? undefined : filter
        }
      })
      setMyProducts(response.data.products || [])
    } catch (error) {
      console.error('Failed to load products:', error)
      toast.error('商品の読み込みに失敗しました')
    } finally {
      setLoading(false)
    }
  }

  const loadMyOffers = async () => {
    try {
      setLoading(true)
      const data = await offerService.getMyOffers('seller')
      setOffers(data)
    } catch (error) {
      console.error('Failed to load offers:', error)
      toast.error('価格交渉の読み込みに失敗しました')
    } finally {
      setLoading(false)
    }
  }

  const handleDeleteProduct = async (id: string) => {
    if (!confirm('この商品を削除してもよろしいですか？')) return

    try {
      await api.delete(`/products/${id}`)
      toast.success('商品を削除しました')
      loadMyProducts()
    } catch (error: any) {
      toast.error(error.response?.data?.error || '削除に失敗しました')
    }
  }

  const handleStartAINegotiation = async (offerId: string) => {
    try {
      console.log('Starting AI negotiation for offer:', offerId)
      await offerService.startAINegotiation(offerId)
      toast.success('AI交渉を開始しました')
      await loadMyOffers()
    } catch (error: any) {
      console.error('AI negotiation error:', error)
      const errorMessage = error.response?.data?.error || error.message || 'AI交渉の開始に失敗しました'
      toast.error(errorMessage)
    }
  }

  const handleRetryAINegotiation = async (offerId: string, customPrompt: string) => {
    try {
      console.log('Retrying AI negotiation with custom prompt for offer:', offerId, 'Prompt:', customPrompt)
      await offerService.retryAINegotiationWithPrompt(offerId, customPrompt)
      toast.success('カスタムプロンプトで再交渉を開始しました')
      await loadMyOffers()
    } catch (error: any) {
      console.error('AI re-negotiation error:', error)
      const errorMessage = error.response?.data?.error || error.message || 'AI再交渉の開始に失敗しました'
      toast.error(errorMessage)
    }
  }

  const handleGetMarketAnalysis = async (offerId: string) => {
    try {
      setLoadingAnalysis(prev => ({ ...prev, [offerId]: true }))
      const analysis = await offerService.getMarketPriceAnalysis(offerId)
      setMarketAnalysis(prev => ({ ...prev, [offerId]: analysis }))
      toast.success('市場価格分析が完了しました')
    } catch (error: any) {
      console.error('Market analysis error:', error)
      const errorMessage = error.response?.data?.error || error.message || '市場価格分析に失敗しました'
      toast.error(errorMessage)
    } finally {
      setLoadingAnalysis(prev => ({ ...prev, [offerId]: false }))
    }
  }

  const handleRespondToOffer = async (offerId: string, accept: boolean) => {
    try {
      console.log('Responding to offer:', offerId, 'Accept:', accept)
      const result = await offerService.respond(offerId, {
        accept,
        message: accept ? '価格交渉を承認しました！' : '申し訳ございませんが、この価格では承認できません。'
      })
      console.log('Offer response result:', result)
      toast.success(accept ? '価格交渉を承認しました！' : '価格交渉を拒否しました')
      await loadMyOffers()
    } catch (error: any) {
      console.error('Offer response error:', error)
      console.error('Error details:', error.response?.data)
      const errorMessage = error.response?.data?.error || error.response?.data?.message || error.message || '価格交渉の返答に失敗しました'
      toast.error(errorMessage)
    }
  }

  const handleEditProduct = (id: string) => {
    navigate(`/products/${id}/edit`)
  }

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  if (!user) return null

  const getStatusBadge = (status: string) => {
    const badges = {
      draft: { color: 'bg-yellow-100 text-yellow-800', text: '下書き（未公開）' },
      active: { color: 'bg-green-100 text-green-800', text: '出品中' },
      sold: { color: 'bg-blue-100 text-blue-800', text: '売却済み' },
      reserved: { color: 'bg-purple-100 text-purple-800', text: '予約済み' },
    }
    const badge = badges[status as keyof typeof badges] || { color: 'bg-gray-100 text-gray-800', text: status }
    return (
      <span className={`px-3 py-1 rounded-full text-xs font-medium ${badge.color}`}>
        {badge.text}
      </span>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <Header />

      <div className="max-w-7xl mx-auto px-4 py-8">
        <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
          {/* User Info Sidebar */}
          <Card className="lg:col-span-1 h-fit">
            <div className="text-center">
              {user.avatar_url ? (
                <img
                  src={user.avatar_url}
                  alt={user.username}
                  className="w-24 h-24 rounded-full mx-auto mb-4 object-cover"
                />
              ) : (
                <div className="w-24 h-24 rounded-full bg-primary-500 flex items-center justify-center text-white text-3xl font-bold mx-auto mb-4">
                  {user.username[0].toUpperCase()}
                </div>
              )}
              <h2 className="text-2xl font-bold text-gray-900">
                {user.display_name || user.username}
              </h2>
              <p className="text-gray-600">@{user.username}</p>

              {user.bio && (
                <p className="mt-4 text-sm text-gray-600">{user.bio}</p>
              )}

              <div className="mt-6 space-y-3">
                <Button variant="primary" className="w-full" onClick={() => navigate('/create')}>
                  新しく出品
                </Button>
                <Button variant="outline" className="w-full" onClick={() => navigate('/ai/create')}>
                  AI自動出品
                </Button>
                <Button variant="outline" className="w-full" onClick={handleLogout}>
                  ログアウト
                </Button>
              </div>
            </div>
          </Card>

          {/* My Products List */}
          <div className="lg:col-span-3">
            <div className="mb-6">
              <h1 className="text-3xl font-bold text-gray-900 mb-4">マイページ</h1>

              {/* Main Tabs */}
              <div className="flex gap-4 mb-6 border-b">
                <button
                  onClick={() => setActiveTab('products')}
                  className={`px-4 py-2 font-medium transition border-b-2 ${
                    activeTab === 'products'
                      ? 'border-primary-500 text-primary-600'
                      : 'border-transparent text-gray-600 hover:text-gray-900'
                  }`}
                >
                  出品した商品
                </button>
                <button
                  onClick={() => setActiveTab('offers')}
                  className={`px-4 py-2 font-medium transition border-b-2 ${
                    activeTab === 'offers'
                      ? 'border-primary-500 text-primary-600'
                      : 'border-transparent text-gray-600 hover:text-gray-900'
                  }`}
                >
                  価格交渉リクエスト
                  {offers.filter(o => o.status === 'pending').length > 0 && (
                    <span className="ml-2 px-2 py-1 bg-red-500 text-white text-xs rounded-full">
                      {offers.filter(o => o.status === 'pending').length}
                    </span>
                  )}
                </button>
              </div>

              {/* Filter Tabs (for products) */}
              {activeTab === 'products' && (
                <div className="flex gap-2 flex-wrap">
                  <button
                    onClick={() => setFilter('all')}
                    className={`px-4 py-2 rounded-lg font-medium transition ${
                      filter === 'all'
                        ? 'bg-primary-500 text-white'
                        : 'bg-white text-gray-700 hover:bg-gray-100'
                    }`}
                  >
                    すべて
                  </button>
                  <button
                    onClick={() => setFilter('draft')}
                    className={`px-4 py-2 rounded-lg font-medium transition ${
                      filter === 'draft'
                        ? 'bg-yellow-500 text-white'
                        : 'bg-white text-gray-700 hover:bg-gray-100'
                    }`}
                  >
                    下書き
                  </button>
                  <button
                    onClick={() => setFilter('active')}
                    className={`px-4 py-2 rounded-lg font-medium transition ${
                      filter === 'active'
                        ? 'bg-green-500 text-white'
                        : 'bg-white text-gray-700 hover:bg-gray-100'
                    }`}
                  >
                    出品中
                  </button>
                  <button
                    onClick={() => setFilter('sold')}
                    className={`px-4 py-2 rounded-lg font-medium transition ${
                      filter === 'sold'
                        ? 'bg-blue-500 text-white'
                        : 'bg-white text-gray-700 hover:bg-gray-100'
                    }`}
                  >
                    売却済み
                  </button>
                </div>
              )}
            </div>

            {/* Products List */}
            {activeTab === 'products' && (
              <>
                {loading ? (
                  <div className="flex justify-center items-center h-64">
                    <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-500" />
                  </div>
                ) : myProducts.length === 0 ? (
                  <Card className="text-center py-12">
                    <h3 className="text-xl font-semibold text-gray-900 mb-2">
                      出品商品がありません
                    </h3>
                    <p className="text-gray-600 mb-4">
                      最初の商品を出品してみましょう
                    </p>
                    <div className="flex gap-3 justify-center">
                      <Button onClick={() => navigate('/create')}>
                        通常出品
                      </Button>
                      <Button variant="outline" onClick={() => navigate('/ai/create')}>
                        AI自動出品
                      </Button>
                    </div>
                  </Card>
                ) : (
                  <div className="space-y-4">
                    {myProducts.map((product) => (
                      <Card key={product.id} className="hover:shadow-lg transition-shadow">
                        <div className="flex gap-4">
                          {/* Product Image */}
                          <div className="flex-shrink-0">
                            <img
                              src={
                                product.images?.[0]?.cdn_url ||
                                product.images?.[0]?.image_url ||
                                PROFILE_PLACEHOLDER
                              }
                              alt={product.title}
                              className="w-32 h-32 object-cover rounded-lg"
                              onError={(e) => {
                                const target = e.target as HTMLImageElement
                                target.src = PROFILE_PLACEHOLDER
                              }}
                            />
                          </div>

                          {/* Product Info */}
                          <div className="flex-1 min-w-0">
                            <div className="flex items-start justify-between gap-4">
                              <div className="flex-1">
                                <h3 className="text-lg font-bold text-gray-900 mb-1">
                                  {product.title}
                                </h3>
                                <p className="text-gray-600 text-sm line-clamp-2 mb-2">
                                  {product.description}
                                </p>
                                <div className="flex items-center gap-2 mb-2">
                                  {getStatusBadge(product.status)}
                                  <span className="text-xs text-gray-500 bg-gray-100 px-2 py-1 rounded">
                                    {product.category}
                                  </span>
                                </div>
                                <div className="flex items-baseline gap-2">
                                  <span className="text-2xl font-bold text-primary-600">
                                    ¥{product.price.toLocaleString()}
                                  </span>
                                  <span className="text-sm text-gray-500">
                                    閲覧数: {product.view_count || 0}
                                  </span>
                                </div>
                              </div>

                              {/* Action Buttons */}
                              <div className="flex flex-col gap-2">
                                {product.status === 'draft' && (
                                  <div className="text-xs text-yellow-600 font-medium mb-1">
                                    未公開
                                  </div>
                                )}
                                <Button
                                  size="sm"
                                  variant="outline"
                                  onClick={() => navigate(`/products/${product.id}`)}
                                >
                                  詳細
                                </Button>
                                {(product.status === 'draft' || product.status === 'active') && (
                                  <>
                                    <Button
                                      size="sm"
                                      variant="outline"
                                      onClick={() => handleEditProduct(product.id)}
                                    >
                                      編集
                                    </Button>
                                    <Button
                                      size="sm"
                                      variant="outline"
                                      className="text-red-600 hover:bg-red-50"
                                      onClick={() => handleDeleteProduct(product.id)}
                                    >
                                      削除
                                    </Button>
                                  </>
                                )}
                              </div>
                            </div>
                          </div>
                        </div>
                      </Card>
                    ))}
                  </div>
                )}
              </>
            )}

            {/* Offers List */}
            {activeTab === 'offers' && (
              <>
                {loading ? (
                  <div className="flex justify-center items-center h-64">
                    <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-500" />
                  </div>
                ) : offers.length === 0 ? (
                  <Card className="text-center py-12">
                    <h3 className="text-xl font-semibold text-gray-900 mb-2">
                      価格交渉リクエストはありません
                    </h3>
                    <p className="text-gray-600">
                      購入者からの価格交渉リクエストがここに表示されます
                    </p>
                  </Card>
                ) : (
                  <div className="space-y-4">
                    {offers.map((offer) => {
                      const getOfferStatusBadge = (status: string) => {
                        const badges = {
                          pending: { color: 'bg-yellow-100 text-yellow-800', text: '保留中' },
                          accepted: { color: 'bg-green-100 text-green-800', text: '承認済み' },
                          rejected: { color: 'bg-red-100 text-red-800', text: '拒否済み' },
                          cancelled: { color: 'bg-gray-100 text-gray-800', text: 'キャンセル済み' },
                        }
                        const badge = badges[status as keyof typeof badges] || { color: 'bg-gray-100 text-gray-800', text: status }
                        return (
                          <span className={`px-3 py-1 rounded-full text-xs font-medium ${badge.color}`}>
                            {badge.text}
                          </span>
                        )
                      }

                      return (
                        <Card key={offer.id} className="hover:shadow-lg transition-shadow">
                          <div className="flex gap-4">
                            {/* Product Image */}
                            <div className="flex-shrink-0">
                              <img
                                src={
                                  offer.product?.images?.[0]?.cdn_url ||
                                  offer.product?.images?.[0]?.image_url ||
                                  PROFILE_PLACEHOLDER
                                }
                                alt={offer.product?.title}
                                className="w-24 h-24 object-cover rounded-lg"
                                onError={(e) => {
                                  const target = e.target as HTMLImageElement
                                  target.src = PROFILE_PLACEHOLDER
                                }}
                              />
                            </div>

                            {/* Offer Info */}
                            <div className="flex-1 min-w-0">
                              <div className="flex items-start justify-between gap-4 mb-3">
                                <div>
                                  <h3 className="text-lg font-bold text-gray-900 mb-1">
                                    {offer.product?.title}
                                  </h3>
                                  {getOfferStatusBadge(offer.status)}
                                </div>
                                <div className="text-right">
                                  <div className="text-xs text-gray-500 mb-1">
                                    {new Date(offer.created_at).toLocaleDateString()}
                                  </div>
                                </div>
                              </div>

                              {/* Buyer Info */}
                              <div className="flex items-center gap-2 mb-3">
                                <div className="w-8 h-8 rounded-full bg-blue-100 flex items-center justify-center text-sm font-bold">
                                  {offer.buyer?.username?.[0]?.toUpperCase() || 'B'}
                                </div>
                                <div className="text-sm">
                                  <span className="text-gray-600">購入希望者: </span>
                                  <span className="font-medium">{offer.buyer?.display_name || offer.buyer?.username}</span>
                                </div>
                              </div>

                              {/* Price Comparison */}
                              <div className="flex items-baseline gap-4 mb-3">
                                <div>
                                  <div className="text-xs text-gray-500">販売価格</div>
                                  <span className="text-lg text-gray-500 line-through">
                                    ¥{offer.product?.price?.toLocaleString()}
                                  </span>
                                </div>
                                <div className="text-2xl font-bold text-gray-400">→</div>
                                <div>
                                  <div className="text-xs text-gray-500">希望価格</div>
                                  <span className="text-2xl font-bold text-primary-600">
                                    ¥{offer.offer_price?.toLocaleString()}
                                  </span>
                                </div>
                              </div>

                              {/* Message */}
                              {offer.message && (
                                <div className="bg-gray-50 rounded-lg p-3 mb-3">
                                  <div className="text-xs text-gray-500 mb-1">メッセージ:</div>
                                  <p className="text-sm text-gray-700">{offer.message}</p>
                                </div>
                              )}

                              {/* Response Message */}
                              {offer.response_message && (
                                <div className="bg-blue-50 rounded-lg p-3 mb-3">
                                  <div className="text-xs text-blue-600 mb-1">あなたの返答:</div>
                                  <p className="text-sm text-gray-700">{offer.response_message}</p>
                                </div>
                              )}

                              {/* Market Price Analysis */}
                              {offer.status === 'pending' && (
                                <div className="bg-gradient-to-r from-green-50 to-teal-50 rounded-lg p-4 mb-3">
                                  <div className="flex items-center justify-between mb-3">
                                    <div className="text-sm font-bold text-green-700">AI市場価格分析</div>
                                    <Button
                                      size="sm"
                                      variant="outline"
                                      onClick={() => handleGetMarketAnalysis(offer.id)}
                                      disabled={loadingAnalysis[offer.id]}
                                      className="bg-green-500 text-white hover:bg-green-600 text-xs"
                                    >
                                      {loadingAnalysis[offer.id] ? '分析中...' : marketAnalysis[offer.id] ? '再分析' : '価格分析'}
                                    </Button>
                                  </div>

                                  {marketAnalysis[offer.id] && (
                                    <div className="space-y-3">
                                      <div className="grid grid-cols-3 gap-2 text-xs">
                                        <div className="bg-white rounded p-2">
                                          <div className="text-gray-500">推奨価格</div>
                                          <div className="text-lg font-bold text-green-600">
                                            ¥{marketAnalysis[offer.id].recommended_price.toLocaleString()}
                                          </div>
                                        </div>
                                        <div className="bg-white rounded p-2">
                                          <div className="text-gray-500">最低価格</div>
                                          <div className="text-sm font-semibold">
                                            ¥{marketAnalysis[offer.id].min_price.toLocaleString()}
                                          </div>
                                        </div>
                                        <div className="bg-white rounded p-2">
                                          <div className="text-gray-500">最高価格</div>
                                          <div className="text-sm font-semibold">
                                            ¥{marketAnalysis[offer.id].max_price.toLocaleString()}
                                          </div>
                                        </div>
                                      </div>

                                      <div className="bg-white rounded p-3">
                                        <div className="text-xs font-semibold text-gray-700 mb-2">市場データ</div>
                                        <div className="grid grid-cols-2 gap-2">
                                          {marketAnalysis[offer.id].market_data_sources.map((data, idx) => (
                                            <div key={idx} className="text-xs border-l-2 border-green-400 pl-2">
                                              <div className="font-semibold">{data.platform}</div>
                                              <div className="text-gray-600">
                                                ¥{data.price.toLocaleString()} ({data.condition})
                                              </div>
                                            </div>
                                          ))}
                                        </div>
                                      </div>

                                      <div className="bg-white rounded p-3">
                                        <div className="text-xs font-semibold text-gray-700 mb-1">分析結果</div>
                                        <div className="text-xs text-gray-600">{marketAnalysis[offer.id].analysis}</div>
                                        <div className="mt-2 flex items-center gap-2">
                                          <span className="text-xs text-gray-500">信頼度:</span>
                                          <span className={`px-2 py-0.5 rounded text-xs font-semibold ${
                                            marketAnalysis[offer.id].confidence_level === 'high' ? 'bg-green-100 text-green-700' :
                                            marketAnalysis[offer.id].confidence_level === 'medium' ? 'bg-yellow-100 text-yellow-700' :
                                            'bg-red-100 text-red-700'
                                          }`}>
                                            {marketAnalysis[offer.id].confidence_level === 'high' ? '高' :
                                             marketAnalysis[offer.id].confidence_level === 'medium' ? '中' : '低'}
                                          </span>
                                        </div>
                                      </div>
                                    </div>
                                  )}
                                </div>
                              )}

                              {/* AI Negotiation Logs */}
                              {offer.ai_negotiation_logs && offer.ai_negotiation_logs.length > 0 && (
                                <div className="bg-gradient-to-r from-blue-50 to-purple-50 rounded-lg p-4 mb-3">
                                  <div className="flex items-center gap-2 mb-3">
                                    <div className="text-sm font-bold text-blue-700">AI自動交渉の記録</div>
                                    {offer.final_ai_price && (
                                      <div className="ml-auto px-3 py-1 bg-blue-600 text-white text-xs font-bold rounded-full">
                                        AI推奨価格: ¥{offer.final_ai_price.toLocaleString()}
                                      </div>
                                    )}
                                  </div>
                                  <div className="space-y-2 max-h-60 overflow-y-auto">
                                    {offer.ai_negotiation_logs.map((log) => (
                                      <div
                                        key={log.id}
                                        className={`text-xs p-2 rounded ${
                                          log.role === 'buyer_ai'
                                            ? 'bg-blue-100 text-blue-900'
                                            : log.role === 'seller_ai'
                                            ? 'bg-green-100 text-green-900'
                                            : 'bg-gray-100 text-gray-900 font-bold'
                                        }`}
                                      >
                                        <div className="font-semibold mb-1">
                                          {log.role === 'buyer_ai' ? '購入者AI' : log.role === 'seller_ai' ? '出品者AI' : 'システム'}
                                          {log.price && ` (¥${log.price.toLocaleString()})`}
                                        </div>
                                        <div>{log.message}</div>
                                      </div>
                                    ))}
                                  </div>

                                  {/* Custom Prompt Re-negotiation */}
                                  {offer.status === 'pending' && (
                                    <div className="mt-3 pt-3 border-t border-purple-200">
                                      <div className="text-xs font-semibold text-purple-700 mb-2">カスタムプロンプトで再交渉</div>
                                      <div className="flex gap-2">
                                        <input
                                          type="text"
                                          placeholder="例: もっと値下げ交渉を強気にして"
                                          className="flex-1 text-xs px-3 py-2 border border-purple-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500"
                                          id={`custom-prompt-${offer.id}`}
                                        />
                                        <Button
                                          size="sm"
                                          variant="outline"
                                          onClick={() => {
                                            const input = document.getElementById(`custom-prompt-${offer.id}`) as HTMLInputElement
                                            const prompt = input?.value?.trim()
                                            if (prompt) {
                                              handleRetryAINegotiation(offer.id, prompt)
                                              input.value = ''
                                            }
                                          }}
                                          className="bg-purple-500 text-white hover:bg-purple-600 text-xs whitespace-nowrap"
                                        >
                                          再交渉
                                        </Button>
                                      </div>
                                    </div>
                                  )}
                                </div>
                              )}

                              {/* Action Buttons */}
                              {offer.status === 'pending' && (
                                <div className="space-y-2">
                                  {/* AI Negotiation Button */}
                                  {(!offer.ai_negotiation_logs || offer.ai_negotiation_logs.length === 0) && (
                                    <Button
                                      size="sm"
                                      variant="outline"
                                      onClick={() => handleStartAINegotiation(offer.id)}
                                      className="w-full bg-gradient-to-r from-blue-500 to-purple-500 text-white hover:from-blue-600 hover:to-purple-600"
                                    >
                                      AIに交渉させる
                                    </Button>
                                  )}

                                  <div className="flex gap-2">
                                    <Button
                                      size="sm"
                                      variant="primary"
                                      onClick={() => handleRespondToOffer(offer.id, true)}
                                      className="flex-1"
                                    >
                                      承認する
                                    </Button>
                                    <Button
                                      size="sm"
                                      variant="outline"
                                      onClick={() => handleRespondToOffer(offer.id, false)}
                                      className="flex-1 text-red-600 hover:bg-red-50"
                                    >
                                      拒否する
                                    </Button>
                                  </div>
                                </div>
                              )}
                            </div>
                          </div>
                        </Card>
                      )
                    })}
                  </div>
                )}
              </>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}

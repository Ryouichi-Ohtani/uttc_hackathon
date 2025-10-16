import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { productService } from '@/services/products'
import { messageService } from '@/services/messages'
import { Product } from '@/types'
import { Button } from '@/components/common/Button'
import { Card } from '@/components/common/Card'
import { Header } from '@/components/layout/Header'
import { Product3DViewer } from '@/components/products/Product3DViewer'
import { ARTryOn } from '@/components/ar/ARTryOn'
import { PricePrediction } from '@/components/analytics/PricePrediction'
import { LiveAuction } from '@/components/auction/LiveAuction'
import { OfferDialog } from '@/components/offers/OfferDialog'
import { useTranslation } from '@/i18n/useTranslation'
import toast from 'react-hot-toast'

export const ProductDetail = () => {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const { t: _t } = useTranslation()
  const [product, setProduct] = useState<Product | null>(null)
  const [loading, setLoading] = useState(true)
  const [selectedImage, setSelectedImage] = useState(0)
  const [activeTab, setActiveTab] = useState<'2d' | '3d' | 'ar'>('2d')
  const [showOfferDialog, setShowOfferDialog] = useState(false)

  useEffect(() => {
    if (id) loadProduct()
  }, [id])

  const loadProduct = async () => {
    try {
      setLoading(true)
      const data = await productService.getById(id!)
      setProduct(data)
    } catch (error) {
      toast.error('Failed to load product')
      navigate('/')
    } finally {
      setLoading(false)
    }
  }

  const handleMessageSeller = async () => {
    if (!product) return
    try {
      const conversation = await messageService.getOrCreateConversation(product.id, product.seller_id)
      navigate(`/chat/${conversation.id}`)
    } catch (error) {
      toast.error('Failed to start conversation')
    }
  }

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-500" />
      </div>
    )
  }

  if (!product) return null

  const co2Comparison = productService.getCO2Comparison(product)
  const currentImage = product.images[selectedImage]

  return (
    <div className="min-h-screen bg-gray-50">
      <Header />
      <div className="max-w-7xl mx-auto px-4 py-8">

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          {/* Images with 3D/AR Toggle */}
          <div>
            {/* View Mode Tabs */}
            <div className="flex gap-2 mb-4">
              <button
                onClick={() => setActiveTab('2d')}
                className={`px-4 py-2 rounded-lg transition ${
                  activeTab === '2d' ? 'bg-primary-500 text-white' : 'bg-gray-200'
                }`}
              >
                📷 2D
              </button>
              <button
                onClick={() => setActiveTab('3d')}
                className={`px-4 py-2 rounded-lg transition ${
                  activeTab === '3d' ? 'bg-primary-500 text-white' : 'bg-gray-200'
                }`}
              >
                🎮 3D
              </button>
              <button
                onClick={() => setActiveTab('ar')}
                className={`px-4 py-2 rounded-lg transition ${
                  activeTab === 'ar' ? 'bg-primary-500 text-white' : 'bg-gray-200'
                }`}
              >
                📸 AR
              </button>
            </div>

            {/* Display based on active tab */}
            {activeTab === '2d' && (
              <>
                <Card padding="none" className="overflow-hidden mb-4">
                  <img
                    src={(currentImage?.cdn_url && currentImage.cdn_url.trim() !== '') ? currentImage.cdn_url : (currentImage?.image_url || 'https://via.placeholder.com/800x600/10B981/FFFFFF?text=EcoMate')}
                    alt={product.title}
                    className="w-full h-96 object-cover"
                    onError={(e) => {
                      const target = e.target as HTMLImageElement
                      target.src = 'https://via.placeholder.com/800x600/10B981/FFFFFF?text=EcoMate'
                    }}
                  />
                </Card>
                <div className="grid grid-cols-4 gap-2">
                  {product.images.map((img, idx) => (
                    <button
                      key={img.id}
                      onClick={() => setSelectedImage(idx)}
                      className={`border-2 rounded-lg overflow-hidden ${
                        idx === selectedImage ? 'border-primary-500' : 'border-gray-200'
                      }`}
                    >
                      <img
                        src={(img.cdn_url && img.cdn_url.trim() !== '') ? img.cdn_url : (img.image_url || 'https://via.placeholder.com/200x200/10B981/FFFFFF?text=EcoMate')}
                        alt=""
                        className="w-full h-20 object-cover"
                        onError={(e) => {
                          const target = e.target as HTMLImageElement
                          target.src = 'https://via.placeholder.com/200x200/10B981/FFFFFF?text=EcoMate'
                        }}
                      />
                    </button>
                  ))}
                </div>
              </>
            )}

            {activeTab === '3d' && (
              <Product3DViewer
                modelUrl={product.model_url}
                fallbackImage={currentImage?.cdn_url || currentImage?.image_url}
                productName={product.title}
              />
            )}

            {activeTab === 'ar' && (
              <ARTryOn
                productImage={(currentImage?.cdn_url && currentImage.cdn_url.trim() !== '') ? currentImage.cdn_url : (currentImage?.image_url || 'https://via.placeholder.com/800x600/10B981/FFFFFF?text=EcoMate')}
                productName={product.title}
                category={product.category}
              />
            )}
          </div>

          {/* Details */}
          <div className="space-y-6">
            <div>
              <h1 className="text-3xl font-bold text-gray-900 mb-2">
                {product.title}
              </h1>
              <div className="flex items-center gap-2 text-sm text-gray-600">
                <span className="bg-gray-100 px-2 py-1 rounded">{product.category}</span>
                <span className="bg-blue-100 text-blue-800 px-2 py-1 rounded">
                  {product.condition}
                </span>
              </div>
            </div>

            <div className="flex items-baseline gap-4">
              <span className="text-4xl font-bold text-primary-600">
                ¥{product.price.toLocaleString()}
              </span>
              {product.ai_suggested_price && (
                <span className="text-lg text-gray-500 line-through">
                  ¥{product.ai_suggested_price.toLocaleString()}
                </span>
              )}
            </div>

            {/* CO2 Impact */}
            <Card className="bg-gradient-to-r from-green-50 to-primary-50">
              <div className="flex items-center gap-3 mb-3">
                <span className="text-3xl">🌱</span>
                <div>
                  <h3 className="font-semibold text-lg text-gray-900">
                    Environmental Impact
                  </h3>
                  <p className="text-sm text-gray-600">
                    Buying this saves the planet!
                  </p>
                </div>
              </div>

              <div className="space-y-2">
                <div className="flex justify-between text-sm">
                  <span className="text-gray-600">Buying new:</span>
                  <span className="font-medium">
                    {co2Comparison.buying_new_kg.toFixed(1)}kg CO2
                  </span>
                </div>
                <div className="flex justify-between text-sm">
                  <span className="text-gray-600">Buying used:</span>
                  <span className="font-medium">
                    {co2Comparison.buying_used_kg.toFixed(1)}kg CO2
                  </span>
                </div>
                <div className="flex justify-between text-lg font-bold text-primary-600 pt-2 border-t border-primary-200">
                  <span>You save:</span>
                  <span>{co2Comparison.saved_kg.toFixed(1)}kg CO2</span>
                </div>
              </div>

              <div className="mt-4 p-3 bg-white/50 rounded-lg text-sm">
                <p className="text-gray-700">
                  That's like planting{' '}
                  <strong>{co2Comparison.equivalent_trees.toFixed(2)} trees</strong> 🌳
                </p>
              </div>
            </Card>

            {/* Description */}
            <Card>
              <h3 className="font-semibold text-lg mb-2">Description</h3>
              <p className="text-gray-700 whitespace-pre-wrap">{product.description}</p>
              {product.ai_generated_description && (
                <div className="mt-2 text-xs text-gray-500 flex items-center gap-1">
                  <span>✨</span>
                  <span>AI-generated description</span>
                </div>
              )}
            </Card>

            {/* Seller Info */}
            {product.seller && (
              <Card>
                <h3 className="font-semibold text-lg mb-3">Seller</h3>
                <div className="flex items-center gap-3">
                  <div className="w-12 h-12 rounded-full bg-primary-100 flex items-center justify-center text-lg font-bold">
                    {product.seller.username[0].toUpperCase()}
                  </div>
                  <div>
                    <div className="font-medium text-gray-900">
                      {product.seller.display_name}
                    </div>
                    <div className="text-sm text-gray-600">
                      @{product.seller.username}
                    </div>
                    <div className="flex items-center gap-2 mt-1">
                      <span className="text-xs bg-primary-100 text-primary-700 px-2 py-0.5 rounded">
                        Level {product.seller.level}
                      </span>
                      <span className="text-xs text-gray-600">
                        {product.seller.total_co2_saved_kg.toFixed(0)}kg CO2 saved
                      </span>
                    </div>
                  </div>
                </div>
              </Card>
            )}

            {/* Actions */}
            <div className="space-y-3">
              <div className="flex gap-3">
                <Button className="flex-1" size="lg" onClick={handleMessageSeller}>
                  💬 メッセージ
                </Button>
                <Button
                  variant="primary"
                  className="flex-1"
                  size="lg"
                  onClick={() => navigate(`/purchase/${product.id}`)}
                >
                  🛒 購入する
                </Button>
              </div>

              <Button
                variant="outline"
                className="w-full"
                size="lg"
                onClick={() => setShowOfferDialog(true)}
              >
                💰 価格交渉をする
              </Button>
            </div>
          </div>
        </div>

        {/* Offer Dialog */}
        {showOfferDialog && (
          <OfferDialog
            productId={product.id}
            currentPrice={product.price}
            onClose={() => setShowOfferDialog(false)}
            onSuccess={() => {
              toast.success('価格交渉を送信しました！出品者の返答をお待ちください。')
              setShowOfferDialog(false)
            }}
          />
        )}

        {/* AI Price Prediction */}
        <div className="mt-8">
          <PricePrediction productId={product.id} currentPrice={product.price} />
        </div>

        {/* Live Auction (if product supports auction) */}
        {(product as any).auction_enabled && (
          <div className="mt-8">
            <LiveAuction
              productId={product.id}
              startingPrice={product.price}
              currentPrice={product.price}
              endTime={new Date(Date.now() + 24 * 60 * 60 * 1000).toISOString()}
            />
          </div>
        )}
      </div>
    </div>
  )
}

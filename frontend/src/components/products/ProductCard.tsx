import { Link } from 'react-router-dom'
import { Product } from '@/types'
import { Card } from '@/components/common/Card'

interface ProductCardProps {
  product: Product
}

export const ProductCard = ({ product }: ProductCardProps) => {
  const primaryImage = product.images.find((img) => img.is_primary) || product.images[0]
  const getImageUrl = () => {
    if (primaryImage?.cdn_url && primaryImage.cdn_url.trim() !== '') {
      return primaryImage.cdn_url
    }
    if (primaryImage?.image_url && primaryImage.image_url.trim() !== '') {
      return primaryImage.image_url
    }
    return 'https://via.placeholder.com/800x600/10B981/FFFFFF?text=EcoMate'
  }
  const imageUrl = getImageUrl()

  return (
    <Link to={`/products/${product.id}`}>
      <Card padding="none" hover className="overflow-hidden !bg-white/80 backdrop-blur-sm">
        <div className="relative">
          <img
            src={imageUrl}
            alt={product.title}
            className="w-full h-48 object-cover"
            onError={(e) => {
              const target = e.target as HTMLImageElement
              target.src = 'https://via.placeholder.com/800x600/10B981/FFFFFF?text=EcoMate'
            }}
          />
          <div className="absolute top-2 right-2 bg-white/90 backdrop-blur-sm px-2 py-1 rounded-full text-xs font-medium">
            {product.condition}
          </div>
        </div>

        <div className="p-4">
          <h3 className="font-semibold text-lg text-gray-900 line-clamp-2 mb-2">
            {product.title}
          </h3>

          <div className="flex items-center justify-between mb-3">
            <span className="text-2xl font-bold text-primary-600">
              ¥{product.price.toLocaleString()}
            </span>
            {product.ai_suggested_price && (
              <span className="text-xs text-gray-500 line-through">
                ¥{product.ai_suggested_price.toLocaleString()}
              </span>
            )}
          </div>

          <div className="flex items-center gap-2 text-sm text-primary-600 bg-primary-50 rounded-lg p-2">
            <span className="text-lg">🌱</span>
            <span className="font-medium">
              {product.co2_impact_kg.toFixed(1)}kg CO2 saved
            </span>
          </div>

          <div className="flex items-center justify-between mt-3 text-sm text-gray-500">
            <div className="flex items-center gap-1">
              <span>♡</span>
              <span>{product.favorite_count}</span>
            </div>
            <div className="flex items-center gap-1">
              <span>👁</span>
              <span>{product.view_count}</span>
            </div>
            <div className="text-xs bg-gray-100 px-2 py-1 rounded">
              {product.category}
            </div>
          </div>

          {product.seller && (
            <div className="mt-3 pt-3 border-t border-gray-100 flex items-center gap-2">
              <div className="w-6 h-6 rounded-full bg-primary-100 flex items-center justify-center text-xs">
                {product.seller.username[0].toUpperCase()}
              </div>
              <span className="text-sm text-gray-600">{product.seller.username}</span>
              <span className="text-xs text-primary-600">
                Lv.{product.seller.level}
              </span>
            </div>
          )}
        </div>
      </Card>
    </Link>
  )
}

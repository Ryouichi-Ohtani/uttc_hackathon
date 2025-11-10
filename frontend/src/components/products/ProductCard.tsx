import { Link } from 'react-router-dom'
import { Product } from '@/types'
import {
  HeartIcon,
  EyeIcon,
  SparklesIcon,
  ComputerDesktopIcon,
  ShoppingBagIcon,
  HomeModernIcon,
  BookOpenIcon,
  PuzzlePieceIcon,
  TrophyIcon,
  CubeIcon
} from '@heroicons/react/24/outline'
import { HeartIcon as HeartIconSolid, StarIcon as StarIconSolid } from '@heroicons/react/24/solid'
import { useState } from 'react'
import { PRODUCT_PLACEHOLDER } from '@/utils/placeholderImages'

interface ProductCardProps {
  product: Product
}

export const ProductCard = ({ product }: ProductCardProps) => {
  const [isFavorited, setIsFavorited] = useState(false)
  const [isImageLoading, setIsImageLoading] = useState(true)

  const primaryImage = product.images?.find((img) => img.is_primary) || product.images?.[0]
  const getImageUrl = () => {
    if (primaryImage?.cdn_url && primaryImage.cdn_url.trim() !== '') {
      return primaryImage.cdn_url
    }
    if (primaryImage?.image_url && primaryImage.image_url.trim() !== '') {
      return primaryImage.image_url
    }
    return PRODUCT_PLACEHOLDER
  }
  const imageUrl = getImageUrl()

  const getConditionBadge = (condition: string) => {
    const conditions: Record<string, { label: string; color: string }> = {
      new: { label: '新品', color: 'bg-gradient-to-r from-emerald-500 to-green-600' },
      like_new: { label: '未使用に近い', color: 'bg-gradient-to-r from-blue-500 to-cyan-600' },
      good: { label: '良好', color: 'bg-gradient-to-r from-indigo-500 to-blue-600' },
      fair: { label: '使用感あり', color: 'bg-gradient-to-r from-amber-500 to-orange-600' },
    }
    return conditions[condition.toLowerCase()] || { label: condition, color: 'bg-slate-600' }
  }

  const getCategoryIcon = (category: string) => {
    const icons: Record<string, any> = {
      electronics: ComputerDesktopIcon,
      clothing: ShoppingBagIcon,
      furniture: HomeModernIcon,
      books: BookOpenIcon,
      toys: PuzzlePieceIcon,
      sports: TrophyIcon,
    }
    return icons[category.toLowerCase()] || CubeIcon
  }

  const conditionBadge = getConditionBadge(product.condition)

  const handleFavorite = (e: React.MouseEvent) => {
    e.preventDefault()
    e.stopPropagation()
    setIsFavorited(!isFavorited)
    // TODO: Add API call to toggle favorite
  }

  return (
    <Link to={`/products/${product.id}`} className="block group">
      <div className="card card-hover overflow-hidden transition-all duration-300 animate-scale bg-white">
        {/* Image Container */}
        <div className="relative aspect-square overflow-hidden bg-gray-100">
          {/* Loading Skeleton */}
          {isImageLoading && (
            <div className="absolute inset-0 skeleton" />
          )}

          <img
            src={imageUrl}
            alt={product.title}
            className={`w-full h-full object-cover transition-all duration-700 group-hover:scale-105 ${
              isImageLoading ? 'opacity-0' : 'opacity-100'
            }`}
            onLoad={() => setIsImageLoading(false)}
            onError={(e) => {
              const target = e.target as HTMLImageElement
              target.src = PRODUCT_PLACEHOLDER
              setIsImageLoading(false)
            }}
          />

          {/* Top Badges */}
          <div className="absolute top-3 left-3 right-3 flex items-start justify-between gap-2">
            {/* AI Badge - More prominent */}
            {product.ai_generated_description && (
              <div className="flex items-center gap-1.5 px-3 py-1.5 bg-gradient-to-r from-primary-500 to-accent-500 text-white text-sm font-bold rounded-full shadow-mercari-hover animate-pulse">
                <SparklesIcon className="w-4 h-4" />
                <span>AI生成</span>
              </div>
            )}

            {/* Condition Badge */}
            <div className={`${conditionBadge.color} text-white text-xs font-bold px-3 py-1 rounded-full shadow-mercari ml-auto`}>
              {conditionBadge.label}
            </div>
          </div>

          {/* Favorite Button */}
          <button
            onClick={handleFavorite}
            className="absolute bottom-3 right-3 w-11 h-11 bg-white/95 backdrop-blur-sm rounded-full flex items-center justify-center shadow-mercari hover:scale-110 transition-transform duration-200"
            aria-label="Add to favorites"
          >
            {isFavorited ? (
              <HeartIconSolid className="w-6 h-6 text-red-500" />
            ) : (
              <HeartIcon className="w-6 h-6 text-gray-600" />
            )}
          </button>

          {/* Hover Overlay */}
          <div className="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300" />
        </div>

        {/* Content */}
        <div className="p-4">
          {/* Category & Title */}
          <div className="mb-3">
            <div className="flex items-center gap-2 mb-1.5">
              {(() => {
                const Icon = getCategoryIcon(product.category)
                return <Icon className="w-4 h-4 text-gray-500" />
              })()}
              <span className="text-xs font-medium text-gray-500 uppercase tracking-wider">
                {product.category}
              </span>
            </div>
            <h3 className="font-semibold text-gray-900 line-clamp-2 group-hover:text-primary-600 transition-colors">
              {product.title}
            </h3>
          </div>

          {/* Price */}
          <div className="mb-3">
            <div className="flex items-baseline gap-2">
              <span className="text-2xl font-bold text-primary-600">
                ¥{product.price.toLocaleString()}
              </span>
              {product.ai_suggested_price && product.ai_suggested_price !== product.price && (
                <span className="text-sm text-gray-400 line-through">
                  ¥{product.ai_suggested_price.toLocaleString()}
                </span>
              )}
            </div>
            {product.ai_suggested_price && product.price < product.ai_suggested_price && (
              <div className="mt-1 inline-flex items-center gap-1 px-2 py-0.5 bg-green-100 text-green-700 rounded-full">
                <SparklesIcon className="w-3 h-3" />
                <span className="text-xs font-bold">
                  {Math.round(((product.ai_suggested_price - product.price) / product.ai_suggested_price) * 100)}% お得
                </span>
              </div>
            )}
          </div>

          {/* Stats */}
          <div className="flex items-center justify-between text-sm text-gray-600 mb-3">
            <div className="flex items-center gap-3">
              <div className="flex items-center gap-1">
                {product.favorite_count > 0 ? (
                  <HeartIconSolid className="w-4 h-4 text-red-500" />
                ) : (
                  <HeartIcon className="w-4 h-4" />
                )}
                <span className="font-medium">{product.favorite_count}</span>
              </div>
              <div className="flex items-center gap-1">
                <EyeIcon className="w-4 h-4" />
                <span className="font-medium">{product.view_count}</span>
              </div>
            </div>
            <div className="px-2 py-1 bg-green-100 text-green-700 text-xs font-bold rounded">
              ECO
            </div>
          </div>

          {/* Seller Info */}
          {product.seller && (
            <div className="pt-3 border-t border-gray-200 flex items-center justify-between">
              <div className="flex items-center gap-2">
                {product.seller.avatar_url ? (
                  <img
                    src={product.seller.avatar_url}
                    alt={product.seller.username}
                    className="w-7 h-7 rounded-full object-cover ring-2 ring-gray-100"
                  />
                ) : (
                  <div className="w-7 h-7 rounded-full bg-gradient-to-br from-primary-500 to-accent-500 flex items-center justify-center text-white text-xs font-bold ring-2 ring-gray-100">
                    {product.seller.username[0].toUpperCase()}
                  </div>
                )}
                <span className="text-sm font-medium text-gray-700">
                  {product.seller.username}
                </span>
              </div>

              {/* Rating */}
              <div className="flex items-center gap-0.5">
                {[...Array(5)].map((_, i) => (
                  <StarIconSolid
                    key={i}
                    className={`w-3.5 h-3.5 ${
                      i < 4 ? 'text-yellow-400' : 'text-gray-300'
                    }`}
                  />
                ))}
                <span className="text-xs font-bold text-gray-700 ml-1">
                  4.8
                </span>
              </div>
            </div>
          )}
        </div>
      </div>
    </Link>
  )
}
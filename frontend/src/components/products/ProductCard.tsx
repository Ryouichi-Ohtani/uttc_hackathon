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
      <div className="card overflow-hidden transition-all duration-500 bg-white dark:bg-dark-card hover:shadow-2xl hover:-translate-y-2 border border-slate-200 dark:border-slate-700 hover:border-primary-300 dark:hover:border-primary-600">
        {/* Image Container */}
        <div className="relative aspect-square overflow-hidden bg-gradient-to-br from-slate-100 to-slate-200 dark:from-slate-800 dark:to-slate-900">
          {/* Loading Skeleton */}
          {isImageLoading && (
            <div className="absolute inset-0 skeleton" />
          )}

          <img
            src={imageUrl}
            alt={product.title}
            className={`w-full h-full object-cover transition-all duration-700 group-hover:scale-110 group-hover:rotate-1 ${
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
          <div className="absolute top-3 left-3 right-3 flex items-start justify-between gap-2 z-10">
            {/* AI Badge - More prominent */}
            {product.ai_generated_description && (
              <div className="flex items-center gap-1.5 px-3 py-2 bg-gradient-to-r from-primary-500 via-accent-500 to-primary-500 bg-size-200 animate-gradient text-white text-sm font-bold rounded-full shadow-lg transform transition-transform duration-300 group-hover:scale-105">
                <SparklesIcon className="w-4 h-4 animate-pulse" />
                <span>AI生成</span>
              </div>
            )}

            {/* Condition Badge */}
            <div className={`${conditionBadge.color} text-white text-xs font-bold px-3 py-1.5 rounded-full shadow-lg ml-auto transform transition-transform duration-300 group-hover:scale-105`}>
              {conditionBadge.label}
            </div>
          </div>

          {/* Favorite Button */}
          <button
            onClick={handleFavorite}
            className="absolute bottom-3 right-3 w-12 h-12 bg-white/95 dark:bg-dark-card/95 backdrop-blur-md rounded-full flex items-center justify-center shadow-lg hover:shadow-xl hover:scale-125 transition-all duration-300 z-10 border-2 border-slate-200 dark:border-slate-700 hover:border-red-300 dark:hover:border-red-500"
            aria-label="Add to favorites"
          >
            {isFavorited ? (
              <HeartIconSolid className="w-6 h-6 text-red-500 animate-pulse" />
            ) : (
              <HeartIcon className="w-6 h-6 text-slate-600 dark:text-slate-400 group-hover:text-red-500 transition-colors" />
            )}
          </button>

          {/* Hover Overlay */}
          <div className="absolute inset-0 bg-gradient-to-t from-black/70 via-black/20 to-transparent opacity-0 group-hover:opacity-100 transition-all duration-500" />

          {/* Quick View Badge - Shows on Hover */}
          <div className="absolute bottom-3 left-3 opacity-0 group-hover:opacity-100 transform translate-y-2 group-hover:translate-y-0 transition-all duration-300 z-10">
            <div className="px-3 py-1.5 bg-white/95 dark:bg-dark-card/95 backdrop-blur-md rounded-lg text-xs font-bold text-slate-900 dark:text-white shadow-lg">
              クリックして詳細
            </div>
          </div>
        </div>

        {/* Content */}
        <div className="p-5 space-y-4">
          {/* Category & Title */}
          <div>
            <div className="flex items-center gap-2 mb-2">
              {(() => {
                const Icon = getCategoryIcon(product.category)
                return <Icon className="w-4 h-4 text-slate-500 dark:text-slate-400" />
              })()}
              <span className="text-xs font-bold text-slate-500 dark:text-slate-400 uppercase tracking-wider">
                {product.category}
              </span>
            </div>
            <h3 className="font-bold text-lg text-slate-900 dark:text-white line-clamp-2 group-hover:text-primary-600 dark:group-hover:text-primary-400 transition-colors leading-snug">
              {product.title}
            </h3>
          </div>

          {/* Price */}
          <div>
            <div className="flex items-baseline gap-2.5">
              <span className="text-3xl font-bold text-primary-600 dark:text-primary-400">
                ¥{product.price.toLocaleString()}
              </span>
              {product.ai_suggested_price && product.ai_suggested_price !== product.price && (
                <span className="text-sm text-slate-400 dark:text-slate-500 line-through font-medium">
                  ¥{product.ai_suggested_price.toLocaleString()}
                </span>
              )}
            </div>
            {product.ai_suggested_price && product.price < product.ai_suggested_price && (
              <div className="mt-2 inline-flex items-center gap-1.5 px-3 py-1.5 bg-gradient-to-r from-emerald-100 to-green-100 dark:from-emerald-900/30 dark:to-green-900/30 text-emerald-700 dark:text-emerald-400 rounded-full shadow-sm">
                <SparklesIcon className="w-4 h-4" />
                <span className="text-xs font-bold">
                  {Math.round(((product.ai_suggested_price - product.price) / product.ai_suggested_price) * 100)}% お得
                </span>
              </div>
            )}
          </div>

          {/* Stats */}
          <div className="flex items-center justify-between text-sm text-slate-600 dark:text-slate-400">
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
            <div className="pt-4 border-t-2 border-slate-200 dark:border-slate-700 flex items-center justify-between">
              <div className="flex items-center gap-3">
                {product.seller.avatar_url ? (
                  <img
                    src={product.seller.avatar_url}
                    alt={product.seller.username}
                    className="w-9 h-9 rounded-full object-cover ring-2 ring-slate-200 dark:ring-slate-700 group-hover:ring-primary-300 dark:group-hover:ring-primary-600 transition-all"
                  />
                ) : (
                  <div className="w-9 h-9 rounded-full bg-gradient-to-br from-primary-500 to-accent-500 flex items-center justify-center text-white text-sm font-bold ring-2 ring-slate-200 dark:ring-slate-700 group-hover:ring-primary-300 dark:group-hover:ring-primary-600 transition-all shadow-md">
                    {product.seller.username[0].toUpperCase()}
                  </div>
                )}
                <span className="text-sm font-bold text-slate-700 dark:text-slate-300 group-hover:text-primary-600 dark:group-hover:text-primary-400 transition-colors">
                  {product.seller.username}
                </span>
              </div>

              {/* Rating */}
              <div className="flex items-center gap-1">
                {[...Array(5)].map((_, i) => (
                  <StarIconSolid
                    key={i}
                    className={`w-4 h-4 transition-transform group-hover:scale-110 ${
                      i < 4 ? 'text-amber-400' : 'text-slate-300 dark:text-slate-600'
                    }`}
                    style={{ transitionDelay: `${i * 30}ms` }}
                  />
                ))}
                <span className="text-xs font-bold text-slate-700 dark:text-slate-300 ml-1">
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
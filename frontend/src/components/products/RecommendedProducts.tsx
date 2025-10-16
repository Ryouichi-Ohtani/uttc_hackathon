import { useState, useEffect } from 'react'
import { ProductCard } from './ProductCard'
import { Product } from '@/types'
import api from '@/services/api'

export const RecommendedProducts = ({ currentProductId }: { currentProductId?: string }) => {
  const [products, setProducts] = useState<Product[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadRecommendations()
  }, [currentProductId])

  const loadRecommendations = async () => {
    try {
      const response = await api.get('/recommendations', {
        params: { limit: 4 }
      })
      setProducts(response.data.recommendations || [])
    } catch (error) {
      // Fallback to latest products
      try {
        const response = await api.get('/products', {
          params: { limit: 4, sort: 'created_desc' },
        })
        setProducts(response.data.products || [])
      } catch (err) {
        console.error('Failed to load recommendations:', err)
      }
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return <div className="text-center py-4">Loading recommendations...</div>
  }

  if (products.length === 0) return null

  return (
    <div>
      <h2 className="text-2xl font-bold mb-4">Recommended for You</h2>
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {products.map((product) => (
          <ProductCard key={product.id} product={product} />
        ))}
      </div>
    </div>
  )
}

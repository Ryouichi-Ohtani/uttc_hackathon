import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { sustainabilityService } from '../services/sustainability'
import { Product } from '../types'
import { Card } from '../components/common/Card'
import { Header } from '../components/layout/Header'

export const Favorites = () => {
  const [favorites, setFavorites] = useState<Product[]>([])
  const [loading, setLoading] = useState(true)
  const navigate = useNavigate()

  useEffect(() => {
    loadFavorites()
  }, [])

  const loadFavorites = async () => {
    try {
      setLoading(true)
      const data = await sustainabilityService.getFavorites()
      setFavorites(data)
    } catch (error) {
      console.error('Failed to load favorites:', error)
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return (
      <div className="min-h-screen">
        <Header />
        <div className="container mx-auto px-4 py-8">
          <div className="flex justify-center items-center h-64">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-500" />
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen">
      <Header />
      <div className="container mx-auto px-4 py-8 bg-white/70 backdrop-blur-sm">
        <h1 className="text-3xl font-bold mb-6">お気に入り</h1>

      {favorites.length === 0 ? (
        <Card padding="lg">
          <div className="text-center py-8 text-gray-500">
            <p>You haven't favorited any products yet.</p>
          </div>
        </Card>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
          {favorites.map((product) => (
            <Card
              key={product.id}
              padding="none"
              hover
              onClick={() => navigate(`/products/${product.id}`)}
            >
              {product.images && product.images.length > 0 ? (
                <img
                  src={product.images[0].cdn_url || product.images[0].image_url}
                  alt={product.title}
                  className="w-full h-48 object-cover rounded-t-lg"
                />
              ) : (
                <div className="w-full h-48 bg-gray-200 flex items-center justify-center rounded-t-lg">
                  <span className="text-gray-400">No image</span>
                </div>
              )}
              <div className="p-4">
                <h3 className="font-semibold text-lg mb-2 truncate">{product.title}</h3>
                <p className="text-2xl font-bold text-green-600 mb-2">¥{product.price.toLocaleString()}</p>
                <div className="flex items-center justify-between text-sm text-gray-600">
                  <span className="capitalize">{product.condition}</span>
                  <span>❤️ {product.favorite_count}</span>
                </div>
                {product.co2_impact_kg && product.co2_impact_kg > 0 && (
                  <div className="mt-2 text-sm text-green-600">
                    🌱 {product.co2_impact_kg.toFixed(2)} kg CO2 saved
                  </div>
                )}
              </div>
            </Card>
          ))}
        </div>
      )}
      </div>
    </div>
  )
}

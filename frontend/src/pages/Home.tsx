import { useState, useEffect } from 'react'
import { useAuthStore } from '@/store/authStore'
import { productService } from '@/services/products'
import { aiService } from '@/services/ai'
import { Product, ProductFilters } from '@/types'
import { ProductCard } from '@/components/products/ProductCard'
import { Input } from '@/components/common/Input'
import { Button } from '@/components/common/Button'
import { Card } from '@/components/common/Card'
import { VoiceSearch } from '@/components/search/VoiceSearch'
import { AIChatbot } from '@/components/chatbot/AIChatbot'
import { ChatbotButton } from '@/components/chatbot/ChatbotButton'
import { LeaderboardSidebar } from '@/components/sustainability/LeaderboardSidebar'
import { Header } from '@/components/layout/Header'
import { useTranslation } from '@/i18n/useTranslation'
import toast from 'react-hot-toast'

export const Home = () => {
  const { user } = useAuthStore()
  const { t, language, setLanguage } = useTranslation()
  const [products, setProducts] = useState<Product[]>([])
  const [loading, setLoading] = useState(true)
  const [filters, setFilters] = useState<ProductFilters>({
    page: 1,
    limit: 20,
    sort: 'created_desc',
  })
  const [searchTerm, setSearchTerm] = useState('')
  const [isChatbotOpen, setIsChatbotOpen] = useState(false)
  const [isTranslating, setIsTranslating] = useState(false)
  const [translationInfo, setTranslationInfo] = useState<string>('')

  useEffect(() => {
    loadProducts()
  }, [filters])

  const loadProducts = async () => {
    try {
      setLoading(true)
      const data = await productService.list(filters)
      setProducts(data.products)
    } catch (error: any) {
      toast.error('Failed to load products')
    } finally {
      setLoading(false)
    }
  }

  const handleSearch = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!searchTerm.trim()) {
      setFilters({ ...filters, search: '', page: 1 })
      setTranslationInfo('')
      return
    }

    try {
      setIsTranslating(true)
      setTranslationInfo('Translating your search...')

      // Get multilingual translation
      const translation = await aiService.translateSearch(searchTerm)

      // Build comprehensive search query with all translations and keywords
      const multilingualQuery = [
        searchTerm,
        translation.japanese,
        translation.english,
        translation.romanized,
        ...translation.keywords
      ].filter(Boolean).join(' ')

      // Show translation info to user
      const langName = translation.detected_language === 'ja' ? '日本語' :
                      translation.detected_language === 'en' ? 'English' :
                      translation.detected_language
      setTranslationInfo(`Detected: ${langName} | Searching: ${translation.search_intent}`)

      // Search with expanded query
      setFilters({ ...filters, search: multilingualQuery, page: 1 })
    } catch (error) {
      console.error('Translation error:', error)
      // Fallback to original search
      setFilters({ ...filters, search: searchTerm, page: 1 })
      setTranslationInfo('Translation unavailable, searching with original query')
    } finally {
      setIsTranslating(false)
    }
  }

  const handleFilterChange = (key: keyof ProductFilters, value: any) => {
    setFilters({ ...filters, [key]: value, page: 1 })
  }

  return (
    <div className="min-h-screen">
      <Header />

      {/* Search and Language Section */}
      <div className="bg-white border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-4 py-4">
          {/* Language Selector */}
          <div className="mb-4 flex items-center gap-2">
            <span className="text-sm text-gray-600">Language:</span>
            <select
              value={language}
              onChange={(e) => setLanguage(e.target.value as any)}
              className="px-3 py-1 border border-gray-300 rounded-lg text-sm"
            >
              <option value="ja">日本語</option>
              <option value="en">English</option>
              <option value="zh">中文</option>
            </select>
          </div>

          {/* Search Bar with Voice */}
          <form onSubmit={handleSearch}>
            <div className="flex gap-2">
              <Input
                type="text"
                placeholder={t('common.search')}
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="flex-1"
                disabled={isTranslating}
              />
              <VoiceSearch
                onTranscript={(text) => {
                  setSearchTerm(text)
                  setFilters({ ...filters, search: text, page: 1 })
                }}
              />
              <Button type="submit" isLoading={isTranslating}>
                {t('common.search')}
              </Button>
            </div>
            {translationInfo && (
              <div className="mt-2 text-sm text-gray-600 bg-primary-50 px-3 py-2 rounded-lg">
                {translationInfo}
              </div>
            )}
          </form>
        </div>
      </div>

      <div className="px-4 py-8 min-h-screen">
        <div className="flex gap-6 justify-start">
          {/* Leaderboard Sidebar - Left */}
          <aside className="hidden lg:block flex-shrink-0">
            <LeaderboardSidebar />
          </aside>

          {/* Main Content */}
          <div className="flex-1 min-w-0">
            <div className="grid grid-cols-12 gap-6">
              {/* Filters Sidebar */}
              <aside className="col-span-12 md:col-span-3">
            <Card className="bg-white/80 backdrop-blur-sm">
              <h2 className="font-semibold text-lg mb-4">Filters</h2>

              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Sort By
                  </label>
                  <select
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg"
                    value={filters.sort}
                    onChange={(e) => handleFilterChange('sort', e.target.value)}
                  >
                    <option value="created_desc">Newest First</option>
                    <option value="price_asc">Price: Low to High</option>
                    <option value="price_desc">Price: High to Low</option>
                    <option value="eco_impact_desc">Most Eco-Friendly</option>
                  </select>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Category
                  </label>
                  <select
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg"
                    value={filters.category || ''}
                    onChange={(e) => handleFilterChange('category', e.target.value || undefined)}
                  >
                    <option value="">All Categories</option>
                    <option value="clothing">Clothing</option>
                    <option value="electronics">Electronics</option>
                    <option value="furniture">Furniture</option>
                    <option value="books">Books</option>
                    <option value="toys">Toys</option>
                    <option value="sports">Sports</option>
                  </select>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Condition
                  </label>
                  <select
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg"
                    value={filters.condition || ''}
                    onChange={(e) => handleFilterChange('condition', e.target.value || undefined)}
                  >
                    <option value="">All Conditions</option>
                    <option value="new">New</option>
                    <option value="like_new">Like New</option>
                    <option value="good">Good</option>
                    <option value="fair">Fair</option>
                    <option value="poor">Poor</option>
                  </select>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Price Range
                  </label>
                  <div className="flex gap-2">
                    <Input
                      type="number"
                      placeholder="Min"
                      onChange={(e) =>
                        handleFilterChange('min_price', parseInt(e.target.value) || undefined)
                      }
                    />
                    <Input
                      type="number"
                      placeholder="Max"
                      onChange={(e) =>
                        handleFilterChange('max_price', parseInt(e.target.value) || undefined)
                      }
                    />
                  </div>
                </div>

                <Button
                  variant="outline"
                  className="w-full"
                  onClick={() => {
                    setFilters({ page: 1, limit: 20, sort: 'created_desc' })
                    setSearchTerm('')
                  }}
                >
                  Clear Filters
                </Button>
              </div>
            </Card>
          </aside>

          {/* Products Grid */}
          <main className="col-span-12 md:col-span-9">
            {loading ? (
              <div className="flex items-center justify-center h-64">
                <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-500" />
              </div>
            ) : products.length === 0 ? (
              <Card className="text-center py-12">
                <h3 className="text-xl font-semibold text-gray-900 mb-2">
                  No products found
                </h3>
                <p className="text-gray-600">Try adjusting your search or filters</p>
              </Card>
            ) : (
              <>
                <div className="mb-4 text-sm text-gray-600">
                  Showing {products.length} products
                </div>
                <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
                  {products.map((product) => (
                    <ProductCard key={product.id} product={product} />
                  ))}
                </div>
              </>
            )}
          </main>
            </div>
          </div>
        </div>
      </div>

      {/* AI Chatbot */}
      <AIChatbot isOpen={isChatbotOpen} onClose={() => setIsChatbotOpen(false)} />
      <ChatbotButton onClick={() => setIsChatbotOpen(true)} isOpen={isChatbotOpen} />
    </div>
  )
}

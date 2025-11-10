import { useState, useEffect } from 'react'
import { useAuthStore } from '@/store/authStore'
import { productService } from '@/services/products'
import { aiService } from '@/services/ai'
import { Product, ProductFilters } from '@/types'
import { ProductCard } from '@/components/products/ProductCard'
import { VoiceSearch } from '@/components/search/VoiceSearch'
import { FloatingAIAssistant } from '@/components/ai/FloatingAIAssistant'
import { LeaderboardSidebar } from '@/components/sustainability/LeaderboardSidebar'
import { Header } from '@/components/layout/Header'
import { useTranslation } from '@/i18n/useTranslation'
import {
  MagnifyingGlassIcon,
  SparklesIcon,
  GlobeAltIcon,
  AdjustmentsHorizontalIcon,
  XMarkIcon,
  ChevronDownIcon,
  ArrowUpIcon,
  ArrowDownIcon,
  ComputerDesktopIcon,
  ShoppingBagIcon,
  HomeModernIcon,
  BookOpenIcon,
  PuzzlePieceIcon,
  TrophyIcon
} from '@heroicons/react/24/outline'
import { StarIcon } from '@heroicons/react/24/solid'
import toast from 'react-hot-toast'

export const Home = () => {
  const { t, language, setLanguage } = useTranslation()
  const [products, setProducts] = useState<Product[]>([])
  const [loading, setLoading] = useState(true)
  const [filters, setFilters] = useState<ProductFilters>({
    page: 1,
    limit: 20,
    sort: 'created_desc',
  })
  const [searchTerm, setSearchTerm] = useState('')
  const [isTranslating, setIsTranslating] = useState(false)
  const [translationInfo, setTranslationInfo] = useState<string>('')
  const [showFilters, setShowFilters] = useState(false)

  useEffect(() => {
    loadProducts()
  }, [filters])

  const loadProducts = async () => {
    try {
      setLoading(true)
      const data = await productService.list(filters)
      setProducts(data.products)
    } catch (error: any) {
      toast.error('商品の読み込みに失敗しました')
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
      setTranslationInfo('検索中...')

      const translation = await aiService.translateSearch(searchTerm)
      const multilingualQuery = [
        searchTerm,
        translation.japanese,
        translation.english,
        translation.romanized,
        ...translation.keywords
      ].filter(Boolean).join(' ')

      const langName = translation.detected_language === 'ja' ? '日本語' :
                      translation.detected_language === 'en' ? 'English' :
                      translation.detected_language
      setTranslationInfo(`検出: ${langName} | ${translation.search_intent}`)

      setFilters({ ...filters, search: multilingualQuery, page: 1 })
    } catch (error) {
      setFilters({ ...filters, search: searchTerm, page: 1 })
      setTranslationInfo('元のクエリで検索中')
    } finally {
      setIsTranslating(false)
    }
  }

  const handleFilterChange = (key: keyof ProductFilters, value: any) => {
    setFilters({ ...filters, [key]: value, page: 1 })
  }

  const categories = [
    { value: 'all', label: 'すべて', Icon: SparklesIcon, gradient: 'from-primary-500 to-accent-500' },
    { value: 'electronics', label: '電子機器', Icon: ComputerDesktopIcon, gradient: 'from-secondary-500 to-secondary-600' },
    { value: 'clothing', label: 'ファッション', Icon: ShoppingBagIcon, gradient: 'from-accent-500 to-accent-600' },
    { value: 'furniture', label: '家具', Icon: HomeModernIcon, gradient: 'from-amber-500 to-orange-500' },
    { value: 'books', label: '書籍', Icon: BookOpenIcon, gradient: 'from-purple-500 to-indigo-500' },
    { value: 'toys', label: 'おもちゃ', Icon: PuzzlePieceIcon, gradient: 'from-green-500 to-emerald-500' },
    { value: 'sports', label: 'スポーツ', Icon: TrophyIcon, gradient: 'from-primary-600 to-accent-600' },
  ]

  const sortOptions = [
    { value: 'created_desc', label: '新着順', icon: <ArrowDownIcon className="w-4 h-4" /> },
    { value: 'created_asc', label: '古い順', icon: <ArrowUpIcon className="w-4 h-4" /> },
    { value: 'price_asc', label: '価格の安い順', icon: <ArrowUpIcon className="w-4 h-4" /> },
    { value: 'price_desc', label: '価格の高い順', icon: <ArrowDownIcon className="w-4 h-4" /> },
    { value: 'popular', label: '人気順', icon: <StarIcon className="w-4 h-4" /> },
  ]

  const conditionOptions = [
    { value: 'all', label: 'すべて' },
    { value: 'new', label: '新品' },
    { value: 'like_new', label: '未使用に近い' },
    { value: 'good', label: '良好' },
    { value: 'fair', label: '使用感あり' },
  ]

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Subtle mesh gradient background */}
      <div className="mesh-gradient" />

      <Header />

      <div className="container-max relative z-10">
        {/* Hero Section with Search */}
        <div className="py-8 mb-8">
          <div className="text-center mb-6 animate-fade-up">
            <h1 className="text-4xl md:text-5xl font-bold mb-3">
              <span className="gradient-text">サステナブルな未来へ</span>
            </h1>
            <p className="text-lg text-slate-600 dark:text-slate-400">
              AI が最適な価格と出品をサポート
            </p>
          </div>

          {/* Language Selector */}
          <div className="flex items-center justify-center gap-2 mb-6">
            <GlobeAltIcon className="w-5 h-5 text-slate-500 dark:text-slate-400" />
            <div className="flex gap-1">
              {['ja', 'en', 'zh'].map((lang) => (
                <button
                  key={lang}
                  onClick={() => setLanguage(lang as any)}
                  className={`px-3 py-1.5 rounded-lg text-sm font-medium transition-all ${
                    language === lang
                      ? 'bg-primary-600 text-white shadow-md'
                      : 'bg-white dark:bg-dark-card text-slate-600 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800'
                  }`}
                >
                  {lang === 'ja' ? '日本語' : lang === 'en' ? 'English' : '中文'}
                </button>
              ))}
            </div>
          </div>

          {/* Search Bar */}
          <form onSubmit={handleSearch} className="max-w-3xl mx-auto">
            <div className="relative group">
              <div className="absolute inset-0 bg-gradient-to-r from-primary-500 to-accent-500 rounded-2xl blur-xl opacity-10 group-hover:opacity-20 transition-opacity" />
              <div className="relative flex gap-2 bg-white rounded-2xl shadow-mercari-hover p-2">
                <div className="flex-1 flex items-center">
                  <MagnifyingGlassIcon className="w-5 h-5 text-gray-400 ml-4 mr-3" />
                  <input
                    type="text"
                    placeholder={t('search.placeholder')}
                    value={searchTerm}
                    onChange={(e) => setSearchTerm(e.target.value)}
                    disabled={isTranslating}
                    className="flex-1 py-3 bg-transparent text-gray-900 placeholder-gray-400 focus:outline-none font-medium"
                  />
                </div>
                <VoiceSearch
                  onTranscript={(text) => {
                    setSearchTerm(text)
                    setTimeout(() => {
                      const form = document.querySelector('form')
                      form?.dispatchEvent(new Event('submit', { bubbles: true }))
                    }, 100)
                  }}
                />
                <button
                  type="submit"
                  disabled={isTranslating}
                  className="btn-gradient px-6 rounded-xl flex items-center gap-2"
                >
                  {isTranslating ? (
                    <>
                      <SparklesIcon className="w-5 h-5 animate-spin" />
                      <span className="hidden sm:inline">翻訳中...</span>
                    </>
                  ) : (
                    <>
                      <SparklesIcon className="w-5 h-5" />
                      <span className="hidden sm:inline">AI検索</span>
                    </>
                  )}
                </button>
              </div>
            </div>
            {translationInfo && (
              <div className="mt-2 text-center text-sm text-slate-600 dark:text-slate-400">
                {translationInfo}
              </div>
            )}
          </form>
        </div>

        {/* Category Pills */}
        <div className="mb-6 overflow-x-auto no-scrollbar">
          <div className="flex gap-3 pb-2">
            {categories.map((cat) => {
              const Icon = cat.Icon
              return (
                <button
                  key={cat.value}
                  onClick={() => handleFilterChange('category', cat.value === 'all' ? undefined : cat.value)}
                  className={`
                    flex items-center gap-2 px-5 py-3 rounded-xl whitespace-nowrap
                    transition-all duration-200 hover-lift font-semibold
                    ${filters.category === cat.value || (cat.value === 'all' && !filters.category)
                      ? `bg-gradient-to-r ${cat.gradient} text-white shadow-mercari-hover`
                      : 'bg-white text-gray-700 hover:shadow-mercari border border-gray-200'
                    }
                  `}
                >
                  <Icon className="w-5 h-5" />
                  <span>{cat.label}</span>
                </button>
              )
            })}
          </div>
        </div>

        <div className="flex gap-6">
          {/* Sidebar - Leaderboard */}
          <aside className="hidden lg:block w-80">
            <div className="sticky top-24">
              <LeaderboardSidebar />
            </div>
          </aside>

          {/* Main Content */}
          <div className="flex-1">
            {/* Toolbar */}
            <div className="card mb-6 p-4">
              <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
                <div className="flex items-center gap-3">
                  <span className="text-sm font-medium text-slate-600 dark:text-slate-400">
                    {products.length} 件の商品
                  </span>
                  {filters.search && (
                    <button
                      onClick={() => {
                        setSearchTerm('')
                        setFilters({ ...filters, search: '', page: 1 })
                        setTranslationInfo('')
                      }}
                      className="flex items-center gap-1 px-2 py-1 bg-slate-100 dark:bg-slate-800 rounded-lg text-xs text-slate-600 dark:text-slate-400 hover:bg-slate-200 dark:hover:bg-slate-700"
                    >
                      <XMarkIcon className="w-3 h-3" />
                      検索をクリア
                    </button>
                  )}
                </div>

                <div className="flex items-center gap-3">
                  {/* Sort Dropdown */}
                  <div className="relative">
                    <select
                      value={filters.sort}
                      onChange={(e) => handleFilterChange('sort', e.target.value)}
                      className="appearance-none bg-white dark:bg-dark-card border border-slate-200 dark:border-dark-border rounded-lg px-4 py-2 pr-10 text-sm font-medium text-slate-700 dark:text-slate-300 focus:outline-none focus:ring-2 focus:ring-primary-500 cursor-pointer"
                    >
                      {sortOptions.map((option) => (
                        <option key={option.value} value={option.value}>
                          {option.label}
                        </option>
                      ))}
                    </select>
                    <ChevronDownIcon className="absolute right-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400 pointer-events-none" />
                  </div>

                  {/* Filter Button */}
                  <button
                    onClick={() => setShowFilters(!showFilters)}
                    className={`btn-secondary flex items-center gap-2 ${showFilters ? 'bg-primary-50 dark:bg-primary-900/20 text-primary-600 dark:text-primary-400' : ''}`}
                  >
                    <AdjustmentsHorizontalIcon className="w-4 h-4" />
                    <span>フィルター</span>
                  </button>
                </div>
              </div>

              {/* Advanced Filters */}
              {showFilters && (
                <div className="mt-4 pt-4 border-t border-slate-200 dark:border-slate-700 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 animate-fade-up">
                  {/* Condition Filter */}
                  <div>
                    <label className="label">状態</label>
                    <select
                      value={filters.condition || 'all'}
                      onChange={(e) => handleFilterChange('condition', e.target.value === 'all' ? undefined : e.target.value)}
                      className="input"
                    >
                      {conditionOptions.map((option) => (
                        <option key={option.value} value={option.value}>
                          {option.label}
                        </option>
                      ))}
                    </select>
                  </div>

                  {/* Price Range */}
                  <div>
                    <label className="label">最低価格</label>
                    <input
                      type="number"
                      placeholder="¥0"
                      value={filters.min_price || ''}
                      onChange={(e) => handleFilterChange('min_price', e.target.value ? Number(e.target.value) : undefined)}
                      className="input"
                    />
                  </div>
                  <div>
                    <label className="label">最高価格</label>
                    <input
                      type="number"
                      placeholder="¥999,999"
                      value={filters.max_price || ''}
                      onChange={(e) => handleFilterChange('max_price', e.target.value ? Number(e.target.value) : undefined)}
                      className="input"
                    />
                  </div>

                  {/* AI Generated Only */}
                  <div>
                    <label className="label">AI生成のみ</label>
                    <button
                      onClick={() => handleFilterChange('ai_generated', !filters.ai_generated)}
                      className={`w-full py-2.5 rounded-lg border-2 transition-colors ${
                        filters.ai_generated
                          ? 'border-primary-500 bg-primary-50 dark:bg-primary-900/20 text-primary-600 dark:text-primary-400'
                          : 'border-slate-300 dark:border-slate-600 hover:border-slate-400 dark:hover:border-slate-500'
                      }`}
                    >
                      <span className="flex items-center justify-center gap-2">
                        <SparklesIcon className="w-4 h-4" />
                        {filters.ai_generated ? 'オン' : 'オフ'}
                      </span>
                    </button>
                  </div>
                </div>
              )}
            </div>

            {/* Products Grid */}
            {loading ? (
              <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
                {[...Array(6)].map((_, i) => (
                  <div key={i} className="card animate-pulse">
                    <div className="aspect-square bg-slate-200 dark:bg-slate-700 rounded-t-xl" />
                    <div className="p-4">
                      <div className="h-4 bg-slate-200 dark:bg-slate-700 rounded mb-3" />
                      <div className="h-6 bg-slate-200 dark:bg-slate-700 rounded mb-3 w-2/3" />
                      <div className="h-4 bg-slate-200 dark:bg-slate-700 rounded w-1/2" />
                    </div>
                  </div>
                ))}
              </div>
            ) : products.length > 0 ? (
              <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
                {products.map((product, index) => (
                  <div
                    key={product.id}
                    className="animate-fade-up"
                    style={{ animationDelay: `${index * 50}ms` }}
                  >
                    <ProductCard product={product} />
                  </div>
                ))}
              </div>
            ) : (
              <div className="card p-12 text-center">
                <div className="mb-4 flex justify-center">
                  <MagnifyingGlassIcon className="w-16 h-16 text-gray-300" />
                </div>
                <h3 className="text-xl font-semibold text-gray-900 mb-2">
                  商品が見つかりませんでした
                </h3>
                <p className="text-gray-600">
                  検索条件を変更してもう一度お試しください
                </p>
              </div>
            )}

            {/* Load More */}
            {products.length >= (filters.limit ?? 20) && (
              <div className="mt-8 text-center">
                <button
                  onClick={() => handleFilterChange('page', (filters.page ?? 1) + 1)}
                  className="btn-primary px-8 py-3"
                >
                  もっと見る
                </button>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Floating AI Assistant */}
      <FloatingAIAssistant />
    </div>
  )
}
import { useNavigate, useLocation } from 'react-router-dom'
import { useAuthStore } from '@/store/authStore'
import { NotificationBadge } from '@/components/notifications/NotificationBadge'
import { LanguageSwitcher } from '@/components/common/LanguageSwitcher'
import { useTranslation } from '@/i18n'
import { useState, useEffect } from 'react'
import {
  HomeIcon,
  PlusCircleIcon,
  ShoppingBagIcon,
  ChatBubbleLeftRightIcon,
  HeartIcon,
  UserCircleIcon,
  SparklesIcon,
  SunIcon,
  MoonIcon,
  Bars3Icon,
  XMarkIcon,
  ShieldCheckIcon
} from '@heroicons/react/24/outline'

export const Header = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const { user } = useAuthStore()
  const { t } = useTranslation()
  const [isDarkMode, setIsDarkMode] = useState(false)
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false)
  const [scrolled, setScrolled] = useState(false)

  useEffect(() => {
    const savedTheme = localStorage.getItem('theme')
    const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
    const shouldBeDark = savedTheme === 'dark' || (!savedTheme && prefersDark)

    setIsDarkMode(shouldBeDark)
    if (shouldBeDark) {
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
    }
  }, [])

  useEffect(() => {
    const handleScroll = () => {
      setScrolled(window.scrollY > 10)
    }
    window.addEventListener('scroll', handleScroll)
    return () => window.removeEventListener('scroll', handleScroll)
  }, [])

  const toggleTheme = () => {
    const newTheme = !isDarkMode
    setIsDarkMode(newTheme)

    if (newTheme) {
      document.documentElement.classList.add('dark')
      localStorage.setItem('theme', 'dark')
    } else {
      document.documentElement.classList.remove('dark')
      localStorage.setItem('theme', 'light')
    }
  }

  const baseMenuItems = [
    { path: '/', label: t('nav.home'), icon: HomeIcon },
    { path: '/create', label: t('nav.create'), icon: PlusCircleIcon },
    { path: '/ai/create', label: t('nav.aiCreate'), icon: SparklesIcon, aiFeature: true },
    { path: '/purchases', label: t('nav.purchases'), icon: ShoppingBagIcon },
    { path: '/messages', label: t('nav.messages'), icon: ChatBubbleLeftRightIcon },
    { path: '/favorites', label: t('nav.favorites'), icon: HeartIcon },
    { path: '/profile', label: t('nav.profile'), icon: UserCircleIcon },
  ]

  // Add admin link for admin users
  const menuItems = user?.role === 'admin'
    ? [
        ...baseMenuItems,
        { path: '/admin', label: t('nav.admin'), icon: ShieldCheckIcon, adminFeature: true }
      ]
    : baseMenuItems

  const isActive = (path: string) => location.pathname === path

  return (
    <>
      <header
        className={`sticky top-0 z-sticky transition-all duration-300 ${
          scrolled
            ? 'bg-white/95 backdrop-blur-lg shadow-mercari'
            : 'bg-white/80 backdrop-blur-md'
        } border-b border-gray-200`}
      >
        <div className="container-max">
          <div className="flex items-center justify-between h-16">
            {/* Logo */}
            <button
              onClick={() => navigate('/')}
              className="flex items-center gap-3 group"
            >
              <div className="relative">
                <div className="w-10 h-10 bg-gradient-to-br from-primary-500 to-accent-500 rounded-xl flex items-center justify-center transform group-hover:scale-110 transition-transform duration-200 shadow-lg">
                  <SparklesIcon className="w-6 h-6 text-white" />
                </div>
                <div className="absolute inset-0 rounded-xl bg-gradient-to-br from-primary-500 to-accent-500 blur-xl opacity-30 group-hover:opacity-50 transition-opacity"></div>
              </div>
              <div className="hidden sm:block">
                <h1 className="text-xl font-bold bg-gradient-to-r from-primary-600 to-accent-600 bg-clip-text text-transparent">
                  Automate
                </h1>
                <p className="text-xs text-slate-500 dark:text-slate-400 -mt-0.5">
                  AI-Powered Marketplace
                </p>
              </div>
            </button>

            {/* Desktop Navigation */}
            {user && (
              <nav className="hidden lg:flex items-center gap-1">
                {menuItems.map((item) => {
                  const Icon = item.icon
                  const active = isActive(item.path)

                  return (
                    <button
                      key={item.path}
                      onClick={() => navigate(item.path)}
                      className={`
                        relative px-4 py-2 rounded-lg font-semibold text-sm
                        transition-all duration-200
                        ${active
                          ? 'text-primary-600 bg-primary-50'
                          : item.aiFeature
                          ? 'text-white bg-gradient-to-r from-primary-500 to-accent-500 hover:shadow-mercari'
                          : (item as any).adminFeature
                          ? 'text-white bg-gradient-to-r from-purple-500 to-pink-500 hover:shadow-mercari'
                          : 'text-gray-600 hover:text-gray-900 hover:bg-gray-100'
                        }
                      `}
                    >
                      <span className="flex items-center gap-2">
                        <Icon className="w-4 h-4" />
                        {item.label}
                      </span>
                      {active && !item.aiFeature && !(item as any).adminFeature && (
                        <div className="absolute bottom-0 left-1/2 -translate-x-1/2 w-8 h-0.5 bg-gradient-to-r from-primary-500 to-accent-500 rounded-full" />
                      )}
                    </button>
                  )
                })}
              </nav>
            )}

            {/* Right Section */}
            <div className="flex items-center gap-2">
              {/* Language Switcher */}
              <LanguageSwitcher />

              {/* Theme Toggle */}
              <button
                onClick={toggleTheme}
                className="p-2 rounded-lg text-slate-600 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors"
                aria-label="Toggle theme"
              >
                {isDarkMode ? (
                  <SunIcon className="w-5 h-5" />
                ) : (
                  <MoonIcon className="w-5 h-5" />
                )}
              </button>

              {user && (
                <>
                  {/* Notifications */}
                  <div className="relative">
                    <NotificationBadge />
                  </div>

                  {/* User Profile */}
                  <button
                    onClick={() => navigate('/profile')}
                    className="flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors"
                  >
                    <div className="hidden sm:block text-right">
                      <div className="text-sm font-medium text-slate-900 dark:text-slate-100">
                        {user.display_name || user.username}
                      </div>
                      <div className="text-xs text-slate-500 dark:text-slate-400">
                        @{user.username}
                      </div>
                    </div>
                    {user.avatar_url ? (
                      <img
                        src={user.avatar_url}
                        alt={user.username}
                        className="w-9 h-9 rounded-lg object-cover border-2 border-slate-200 dark:border-slate-700"
                      />
                    ) : (
                      <div className="w-9 h-9 rounded-lg bg-gradient-to-br from-primary-500 to-accent-500 flex items-center justify-center text-white font-semibold shadow-md">
                        {user.username[0].toUpperCase()}
                      </div>
                    )}
                  </button>

                  {/* Mobile Menu Toggle */}
                  <button
                    onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
                    className="lg:hidden p-2 rounded-lg text-slate-600 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors"
                  >
                    {isMobileMenuOpen ? (
                      <XMarkIcon className="w-6 h-6" />
                    ) : (
                      <Bars3Icon className="w-6 h-6" />
                    )}
                  </button>
                </>
              )}
            </div>
          </div>

          {/* Mobile Navigation */}
          {user && isMobileMenuOpen && (
            <nav className="lg:hidden py-4 border-t border-slate-200 dark:border-slate-800 animate-fade-up">
              <div className="grid grid-cols-2 gap-2">
                {menuItems.map((item) => {
                  const Icon = item.icon
                  const active = isActive(item.path)

                  return (
                    <button
                      key={item.path}
                      onClick={() => {
                        navigate(item.path)
                        setIsMobileMenuOpen(false)
                      }}
                      className={`
                        flex items-center gap-2 px-4 py-3 rounded-lg font-semibold text-sm
                        transition-all duration-200
                        ${active
                          ? 'text-primary-600 bg-primary-50'
                          : item.aiFeature
                          ? 'text-white bg-gradient-to-r from-primary-500 to-accent-500 hover:shadow-mercari'
                          : (item as any).adminFeature
                          ? 'text-white bg-gradient-to-r from-purple-500 to-pink-500 hover:shadow-mercari'
                          : 'text-gray-600 hover:text-gray-900 hover:bg-gray-100'
                        }
                      `}
                    >
                      <Icon className="w-5 h-5" />
                      <span>{item.label}</span>
                    </button>
                  )
                })}
              </div>
            </nav>
          )}
        </div>
      </header>

      {/* Premium accent line */}
      <div className="h-0.5 bg-gradient-to-r from-primary-500 via-accent-500 to-primary-500 opacity-30" />
    </>
  )
}
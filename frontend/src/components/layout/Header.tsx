import { useNavigate, useLocation } from 'react-router-dom'
import { useAuthStore } from '@/store/authStore'
import { Card } from '@/components/common/Card'
import { NotificationBadge } from '@/components/notifications/NotificationBadge'

export const Header = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const { user } = useAuthStore()

  const menuItems = [
    { path: '/', label: 'ホーム', icon: '' },
    { path: '/create', label: '出品する', icon: '' },
    { path: '/purchases', label: '購入履歴', icon: '' },
    { path: '/messages', label: 'メッセージ', icon: '' },
    { path: '/favorites', label: 'お気に入り', icon: '' },
    { path: '/leaderboard', label: 'ランキング', icon: '' },
    { path: '/profile', label: 'マイページ', icon: '' },
  ]

  const isActive = (path: string) => {
    return location.pathname === path
  }

  return (
    <header className="bg-gradient-to-r from-green-50/95 via-emerald-50/95 to-teal-50/95 backdrop-blur-sm shadow-sm border-b border-gray-200 sticky top-0 z-50">
      <div className="max-w-7xl mx-auto px-4 py-4">
        <div className="flex items-center justify-between mb-4">
          <button
            onClick={() => navigate('/')}
            className="flex items-center gap-2 hover:opacity-80 transition-opacity"
          >
            <h1 className="text-2xl font-bold text-gray-900">Automate</h1>
          </button>

          {user && (
            <div className="flex items-center gap-4">
              <NotificationBadge />
              <Card className="flex items-center gap-3 py-2 cursor-pointer hover:shadow-md transition-shadow" onClick={() => navigate('/profile')}>
                <div className="text-sm">
                  <div className="font-medium text-gray-900">{user.display_name || user.username}</div>
                  <div className="text-gray-500 text-xs">
                    @{user.username}
                  </div>
                </div>
                {user.avatar_url ? (
                  <img src={user.avatar_url} alt={user.username} className="w-10 h-10 rounded-full" />
                ) : (
                  <div className="w-10 h-10 rounded-full bg-primary-500 flex items-center justify-center text-white font-bold">
                    {user.username[0].toUpperCase()}
                  </div>
                )}
              </Card>
            </div>
          )}
        </div>

        {/* Navigation Menu */}
        {user && (
          <nav className="flex items-center gap-1 overflow-x-auto pb-2 scrollbar-hide">
            {menuItems.map((item) => (
              <button
                key={item.path}
                onClick={() => navigate(item.path)}
                className={`flex items-center gap-2 px-4 py-2 rounded-lg whitespace-nowrap transition-all ${
                  isActive(item.path)
                    ? 'bg-primary-500 text-white shadow-md'
                    : 'bg-white/80 text-gray-700 hover:bg-white hover:shadow-sm'
                }`}
              >
                {item.icon && <span>{item.icon}</span>}
                <span className="text-sm font-medium">{item.label}</span>
              </button>
            ))}
          </nav>
        )}
      </div>
    </header>
  )
}

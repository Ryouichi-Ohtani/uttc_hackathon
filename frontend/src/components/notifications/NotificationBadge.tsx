import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import api from '@/services/api'

export const NotificationBadge = () => {
  const [unreadCount, setUnreadCount] = useState(0)

  useEffect(() => {
    loadUnreadCount()

    // Poll for new notifications every 30 seconds
    const interval = setInterval(loadUnreadCount, 30000)
    return () => clearInterval(interval)
  }, [])

  const loadUnreadCount = async () => {
    try {
      const response = await api.get('/notifications/unread-count')
      setUnreadCount(response.data.count || 0)
    } catch (error) {
      // Silent fail
    }
  }

  if (unreadCount === 0) {
    return (
      <Link to="/notifications" className="relative">
        <span className="text-2xl">🔔</span>
      </Link>
    )
  }

  return (
    <Link to="/notifications" className="relative">
      <span className="text-2xl">🔔</span>
      <span className="absolute -top-1 -right-1 bg-red-500 text-white text-xs rounded-full w-5 h-5 flex items-center justify-center font-bold">
        {unreadCount > 99 ? '99+' : unreadCount}
      </span>
    </Link>
  )
}

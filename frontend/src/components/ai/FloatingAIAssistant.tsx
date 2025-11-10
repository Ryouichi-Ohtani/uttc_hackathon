import { useState } from 'react'
import {
  SparklesIcon,
  ChatBubbleLeftRightIcon,
  MagnifyingGlassIcon,
  CameraIcon,
  XMarkIcon
} from '@heroicons/react/24/outline'
import { AIChatbot } from '@/components/chatbot/AIChatbot'

interface QuickAction {
  icon: any
  label: string
  onClick: () => void
  gradient: string
}

export const FloatingAIAssistant = () => {
  const [isChatbotOpen, setIsChatbotOpen] = useState(false)
  const [showActions, setShowActions] = useState(false)

  const quickActions: QuickAction[] = [
    {
      icon: ChatBubbleLeftRightIcon,
      label: 'AIチャット',
      onClick: () => {
        setIsChatbotOpen(true)
        setShowActions(false)
      },
      gradient: 'from-primary-500 to-accent-500'
    },
    {
      icon: MagnifyingGlassIcon,
      label: 'AI検索',
      onClick: () => {
        window.scrollTo({ top: 0, behavior: 'smooth' })
        setShowActions(false)
        // Focus on search input
        setTimeout(() => {
          const searchInput = document.querySelector('input[type="text"]') as HTMLInputElement
          searchInput?.focus()
        }, 500)
      },
      gradient: 'from-secondary-500 to-secondary-600'
    },
    {
      icon: CameraIcon,
      label: 'AI画像検索',
      onClick: () => {
        // Navigate to AI product creation
        window.location.href = '/ai/create'
        setShowActions(false)
      },
      gradient: 'from-accent-500 to-accent-600'
    }
  ]

  return (
    <>
      {/* Floating Button */}
      <div className="fixed bottom-6 right-6 z-50">
        {/* Quick Actions Menu */}
        {showActions && (
          <div className="absolute bottom-20 right-0 mb-2 space-y-2 animate-fade-up">
            {quickActions.map((action, index) => {
              const Icon = action.icon
              return (
                <button
                  key={index}
                  onClick={action.onClick}
                  className={`
                    flex items-center gap-3 px-4 py-3 rounded-full
                    bg-gradient-to-r ${action.gradient} text-white
                    shadow-mercari-hover hover:scale-105
                    transition-all duration-200
                    whitespace-nowrap font-semibold
                  `}
                  style={{ animationDelay: `${index * 50}ms` }}
                >
                  <Icon className="w-5 h-5" />
                  <span>{action.label}</span>
                </button>
              )
            })}
          </div>
        )}

        {/* Main AI Button */}
        <button
          onClick={() => setShowActions(!showActions)}
          className={`
            w-16 h-16 rounded-full
            bg-gradient-to-r from-primary-500 to-accent-500
            text-white shadow-mercari-hover
            hover:scale-110 active:scale-95
            transition-all duration-200
            flex items-center justify-center
            ${showActions ? 'rotate-90' : ''}
          `}
        >
          {showActions ? (
            <XMarkIcon className="w-7 h-7" />
          ) : (
            <SparklesIcon className="w-7 h-7 animate-pulse" />
          )}
        </button>

        {/* AI Badge */}
        {!showActions && (
          <div className="absolute -top-2 -right-2 bg-white rounded-full px-2 py-0.5 shadow-mercari">
            <span className="text-xs font-bold text-primary-600">AI</span>
          </div>
        )}
      </div>

      {/* Chatbot Modal */}
      <AIChatbot
        isOpen={isChatbotOpen}
        onClose={() => setIsChatbotOpen(false)}
      />
    </>
  )
}

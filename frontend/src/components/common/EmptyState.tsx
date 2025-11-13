import { FC, ReactNode } from 'react'
import { SparklesIcon, MagnifyingGlassIcon, ShoppingBagIcon, InboxIcon } from '@heroicons/react/24/outline'

interface EmptyStateProps {
  type?: 'search' | 'products' | 'inbox' | 'custom'
  title: string
  description?: string
  icon?: ReactNode
  action?: {
    label: string
    onClick: () => void
  }
}

export const EmptyState: FC<EmptyStateProps> = ({
  type = 'custom',
  title,
  description,
  icon,
  action
}) => {
  const getDefaultIcon = () => {
    if (icon) return icon

    switch (type) {
      case 'search':
        return <MagnifyingGlassIcon className="w-20 h-20" />
      case 'products':
        return <ShoppingBagIcon className="w-20 h-20" />
      case 'inbox':
        return <InboxIcon className="w-20 h-20" />
      default:
        return <SparklesIcon className="w-20 h-20" />
    }
  }

  return (
    <div className="card p-16 text-center animate-fade-up">
      {/* Icon */}
      <div className="mb-6 flex justify-center">
        <div className="relative">
          <div className="absolute inset-0 bg-gradient-to-r from-primary-500 to-accent-500 opacity-20 blur-3xl rounded-full animate-pulse" />
          <div className="relative text-slate-300 dark:text-slate-600 animate-float">
            {getDefaultIcon()}
          </div>
        </div>
      </div>

      {/* Title */}
      <h3 className="text-2xl font-bold text-slate-900 dark:text-white mb-3 animate-fade-up stagger-1">
        {title}
      </h3>

      {/* Description */}
      {description && (
        <p className="text-slate-600 dark:text-slate-400 max-w-md mx-auto mb-6 animate-fade-up stagger-2">
          {description}
        </p>
      )}

      {/* Action Button */}
      {action && (
        <div className="animate-fade-up stagger-3">
          <button
            onClick={action.onClick}
            className="btn-gradient btn-ripple px-8 py-3"
          >
            {action.label}
          </button>
        </div>
      )}

      {/* Decorative elements */}
      <div className="mt-8 flex justify-center gap-2">
        {[...Array(3)].map((_, i) => (
          <div
            key={i}
            className="w-2 h-2 rounded-full bg-slate-300 dark:bg-slate-600 animate-bounce"
            style={{ animationDelay: `${i * 0.15}s` }}
          />
        ))}
      </div>
    </div>
  )
}

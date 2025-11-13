import { FC } from 'react'

interface LoadingSpinnerProps {
  size?: 'sm' | 'md' | 'lg' | 'xl'
  type?: 'spinner' | 'dots' | 'pulse' | 'bars'
  className?: string
  text?: string
}

export const LoadingSpinner: FC<LoadingSpinnerProps> = ({
  size = 'md',
  type = 'spinner',
  className = '',
  text
}) => {
  const sizeClasses = {
    sm: 'w-4 h-4',
    md: 'w-8 h-8',
    lg: 'w-12 h-12',
    xl: 'w-16 h-16'
  }

  const dotSizes = {
    sm: 'w-2 h-2',
    md: 'w-3 h-3',
    lg: 'w-4 h-4',
    xl: 'w-5 h-5'
  }

  const barSizes = {
    sm: 'w-1 h-4',
    md: 'w-1.5 h-6',
    lg: 'w-2 h-8',
    xl: 'w-3 h-10'
  }

  if (type === 'dots') {
    return (
      <div className={`flex flex-col items-center gap-4 ${className}`}>
        <div className="flex gap-2">
          {[...Array(3)].map((_, i) => (
            <div
              key={i}
              className={`${dotSizes[size]} rounded-full bg-gradient-to-r from-primary-500 to-accent-500 animate-bounce`}
              style={{ animationDelay: `${i * 0.15}s` }}
            />
          ))}
        </div>
        {text && (
          <p className="text-sm font-medium text-slate-600 dark:text-slate-400 animate-pulse">
            {text}
          </p>
        )}
      </div>
    )
  }

  if (type === 'pulse') {
    return (
      <div className={`flex flex-col items-center gap-4 ${className}`}>
        <div className={`${sizeClasses[size]} relative`}>
          <div className="absolute inset-0 rounded-full bg-gradient-to-r from-primary-500 to-accent-500 animate-ping" />
          <div className="absolute inset-0 rounded-full bg-gradient-to-r from-primary-500 to-accent-500 opacity-75" />
        </div>
        {text && (
          <p className="text-sm font-medium text-slate-600 dark:text-slate-400 animate-pulse">
            {text}
          </p>
        )}
      </div>
    )
  }

  if (type === 'bars') {
    return (
      <div className={`flex flex-col items-center gap-4 ${className}`}>
        <div className="flex gap-1.5 items-end">
          {[...Array(5)].map((_, i) => (
            <div
              key={i}
              className={`${barSizes[size]} bg-gradient-to-t from-primary-500 to-accent-500 rounded-full animate-pulse`}
              style={{
                animationDelay: `${i * 0.1}s`,
                animationDuration: '1s'
              }}
            />
          ))}
        </div>
        {text && (
          <p className="text-sm font-medium text-slate-600 dark:text-slate-400 animate-pulse">
            {text}
          </p>
        )}
      </div>
    )
  }

  // Default spinner
  return (
    <div className={`flex flex-col items-center gap-4 ${className}`}>
      <div className={`${sizeClasses[size]} relative`}>
        <div className="absolute inset-0 rounded-full border-4 border-slate-200 dark:border-slate-700" />
        <div className="absolute inset-0 rounded-full border-4 border-transparent border-t-primary-500 border-r-accent-500 animate-spin" />
      </div>
      {text && (
        <p className="text-sm font-medium text-slate-600 dark:text-slate-400 animate-pulse">
          {text}
        </p>
      )}
    </div>
  )
}

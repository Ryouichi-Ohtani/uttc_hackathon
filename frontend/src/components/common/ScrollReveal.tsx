import { FC, ReactNode } from 'react'
import { useScrollReveal } from '@/hooks/useScrollReveal'

interface ScrollRevealProps {
  children: ReactNode
  className?: string
  delay?: number
  threshold?: number
  animation?: 'fade-up' | 'fade-down' | 'fade-left' | 'fade-right' | 'scale' | 'bounce'
}

export const ScrollReveal: FC<ScrollRevealProps> = ({
  children,
  className = '',
  delay = 0,
  threshold = 0.1,
  animation = 'fade-up'
}) => {
  const { ref, isRevealed } = useScrollReveal({ threshold })

  const animationClasses = {
    'fade-up': 'animate-fade-up',
    'fade-down': 'animate-fade-down',
    'fade-left': 'animate-fade-left',
    'fade-right': 'animate-fade-right',
    'scale': 'animate-scale',
    'bounce': 'animate-bounce-in'
  }

  return (
    <div
      ref={ref as any}
      className={`${className} ${isRevealed ? animationClasses[animation] : 'opacity-0'}`}
      style={{ animationDelay: `${delay}ms` }}
    >
      {children}
    </div>
  )
}

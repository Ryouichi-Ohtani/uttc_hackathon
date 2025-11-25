import { useState, useEffect, useRef, ImgHTMLAttributes } from 'react'
import clsx from 'clsx'

interface LazyImageProps extends Omit<ImgHTMLAttributes<HTMLImageElement>, 'src'> {
  src: string
  alt: string
  placeholder?: string
  className?: string
  wrapperClassName?: string
  threshold?: number
  fadeIn?: boolean
}

export const LazyImage = ({
  src,
  alt,
  placeholder = 'data:image/svg+xml,%3Csvg xmlns=\'http://www.w3.org/2000/svg\' viewBox=\'0 0 400 300\'%3E%3Crect fill=\'%23f0f0f0\' width=\'400\' height=\'300\'/%3E%3C/svg%3E',
  className,
  wrapperClassName,
  threshold = 0.1,
  fadeIn = true,
  ...props
}: LazyImageProps) => {
  const [isLoaded, setIsLoaded] = useState(false)
  const [isInView, setIsInView] = useState(false)
  const [hasError, setHasError] = useState(false)
  const imgRef = useRef<HTMLImageElement>(null)
  const wrapperRef = useRef<HTMLDivElement>(null)

  // Intersection Observer for lazy loading
  useEffect(() => {
    if (!wrapperRef.current) return

    const observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            setIsInView(true)
            observer.disconnect()
          }
        })
      },
      {
        threshold,
        rootMargin: '50px', // Start loading 50px before entering viewport
      }
    )

    observer.observe(wrapperRef.current)

    return () => {
      observer.disconnect()
    }
  }, [threshold])

  // Handle image load
  const handleLoad = () => {
    setIsLoaded(true)
  }

  // Handle image error
  const handleError = () => {
    setHasError(true)
    setIsLoaded(true)
  }

  return (
    <div ref={wrapperRef} className={clsx('relative overflow-hidden', wrapperClassName)}>
      {/* Placeholder/Skeleton */}
      {!isLoaded && (
        <div className="absolute inset-0 bg-gradient-to-br from-slate-200 to-slate-300 dark:from-slate-700 dark:to-slate-800 animate-pulse" />
      )}

      {/* Error State */}
      {hasError && (
        <div className="absolute inset-0 flex items-center justify-center bg-slate-100 dark:bg-slate-800">
          <div className="text-center">
            <svg
              className="w-12 h-12 mx-auto mb-2 text-slate-400 dark:text-slate-600"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"
              />
            </svg>
            <p className="text-xs text-slate-500 dark:text-slate-400">画像を読み込めません</p>
          </div>
        </div>
      )}

      {/* Actual Image */}
      {isInView && (
        <img
          ref={imgRef}
          src={hasError ? placeholder : src}
          alt={alt}
          onLoad={handleLoad}
          onError={handleError}
          className={clsx(
            className,
            fadeIn && 'transition-opacity duration-500',
            isLoaded ? 'opacity-100' : 'opacity-0'
          )}
          loading="lazy" // Native lazy loading as fallback
          {...props}
        />
      )}
    </div>
  )
}

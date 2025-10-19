import { useEffect, useRef, useState } from 'react'

interface Product3DViewerProps {
  modelUrl?: string
  productTitle: string
}

export const Product3DViewer = ({ modelUrl, productTitle }: Product3DViewerProps) => {
  const containerRef = useRef<HTMLDivElement>(null)
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    if (!modelUrl) {
      setIsLoading(false)
      return
    }

    // In production, use Three.js to load and render 3D model
    // For now, show placeholder
    const timer = setTimeout(() => {
      setIsLoading(false)
    }, 1000)

    return () => clearTimeout(timer)
  }, [modelUrl])

  if (!modelUrl) {
    return (
      <div className="bg-gray-100 rounded-lg p-8 text-center">
        <div className="text-gray-400 mb-2">
          <svg className="w-16 h-16 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
          </svg>
        </div>
        <p className="text-gray-500">3D model not available</p>
      </div>
    )
  }

  if (isLoading) {
    return (
      <div className="bg-gray-100 rounded-lg p-8 text-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-green-600 mx-auto"></div>
        <p className="mt-4 text-gray-600">Loading 3D model...</p>
      </div>
    )
  }

  return (
    <div className="relative w-full h-96 bg-gradient-to-br from-gray-900 to-gray-700 rounded-lg overflow-hidden">
      <div ref={containerRef} className="w-full h-full">
        {/* In production: Three.js WebGL canvas would render here */}
        <div className="flex items-center justify-center h-full">
          <div className="text-center">
            <div className="text-white mb-4">
              <svg className="w-24 h-24 mx-auto animate-pulse" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
              </svg>
            </div>
            <p className="text-white text-lg font-semibold">{productTitle}</p>
            <p className="text-gray-300 text-sm mt-2">Interactive 3D View</p>
            <div className="mt-4 space-x-2">
              <button className="px-4 py-2 bg-white/20 hover:bg-white/30 text-white rounded-lg text-sm transition-colors">
                Rotate
              </button>
              <button className="px-4 py-2 bg-white/20 hover:bg-white/30 text-white rounded-lg text-sm transition-colors">
                Zoom
              </button>
            </div>
          </div>
        </div>
      </div>

      {/* Controls overlay */}
      <div className="absolute bottom-4 left-4 right-4 bg-black/50 backdrop-blur-sm rounded-lg p-3">
        <div className="flex items-center justify-between text-white text-sm">
          <span>Drag to rotate - Scroll to zoom</span>
          <button className="px-3 py-1 bg-white/20 hover:bg-white/30 rounded transition-colors">
            Fullscreen
          </button>
        </div>
      </div>
    </div>
  )
}

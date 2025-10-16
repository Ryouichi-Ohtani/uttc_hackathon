import { useEffect, useRef, useState } from 'react'
import { Card } from '@/components/common/Card'
import { Button } from '@/components/common/Button'

interface Product3DViewerProps {
  productId: string
  productName: string
  imageUrl?: string
}

export const Product3DViewer = ({ productId, productName, imageUrl }: Product3DViewerProps) => {
  const containerRef = useRef<HTMLDivElement>(null)
  const [isARMode, setIsARMode] = useState(false)
  const [isLoading, setIsLoading] = useState(false)
  const [rotation, setRotation] = useState(0)

  useEffect(() => {
    // Simple 3D simulation with CSS transforms
    if (containerRef.current && !isARMode) {
      const interval = setInterval(() => {
        setRotation((prev) => (prev + 1) % 360)
      }, 50)

      return () => clearInterval(interval)
    }
  }, [isARMode])

  const handleARView = () => {
    setIsLoading(true)
    // Simulate AR loading
    setTimeout(() => {
      setIsARMode(true)
      setIsLoading(false)
    }, 1500)
  }

  const handleClose = () => {
    setIsARMode(false)
  }

  if (isARMode) {
    return (
      <div className="fixed inset-0 bg-black bg-opacity-90 z-50 flex items-center justify-center">
        <div className="relative w-full h-full">
          {/* AR View Simulation */}
          <div className="absolute inset-0 flex items-center justify-center">
            <div className="relative">
              {/* Camera feed simulation */}
              <div className="w-full h-screen bg-gradient-to-b from-gray-800 to-gray-900">
                {/* AR Grid */}
                <div className="absolute inset-0 opacity-20">
                  {Array.from({ length: 10 }).map((_, i) => (
                    <div
                      key={i}
                      className="absolute border-b border-green-400"
                      style={{
                        top: `${i * 10}%`,
                        width: '100%',
                        height: '1px',
                      }}
                    />
                  ))}
                  {Array.from({ length: 10 }).map((_, i) => (
                    <div
                      key={i}
                      className="absolute border-r border-green-400"
                      style={{
                        left: `${i * 10}%`,
                        height: '100%',
                        width: '1px',
                      }}
                    />
                  ))}
                </div>

                {/* Product in AR */}
                <div className="absolute inset-0 flex items-center justify-center">
                  <div
                    className="transform-gpu transition-transform"
                    style={{
                      transform: `perspective(1000px) rotateY(${rotation}deg) rotateX(15deg)`,
                    }}
                  >
                    <div className="relative">
                      {imageUrl ? (
                        <img
                          src={imageUrl}
                          alt={productName}
                          className="w-64 h-64 object-contain rounded-lg shadow-2xl"
                          style={{
                            filter: 'drop-shadow(0 0 30px rgba(34, 197, 94, 0.5))',
                          }}
                        />
                      ) : (
                        <div className="w-64 h-64 bg-gradient-to-br from-primary-400 to-primary-600 rounded-lg shadow-2xl flex items-center justify-center">
                          <span className="text-6xl">📦</span>
                        </div>
                      )}
                      {/* AR Markers */}
                      <div className="absolute -top-2 -left-2 w-4 h-4 border-t-2 border-l-2 border-green-400" />
                      <div className="absolute -top-2 -right-2 w-4 h-4 border-t-2 border-r-2 border-green-400" />
                      <div className="absolute -bottom-2 -left-2 w-4 h-4 border-b-2 border-l-2 border-green-400" />
                      <div className="absolute -bottom-2 -right-2 w-4 h-4 border-b-2 border-r-2 border-green-400" />
                    </div>
                  </div>
                </div>

                {/* AR Info Overlay */}
                <div className="absolute top-4 left-4 bg-black bg-opacity-70 text-white px-4 py-2 rounded-lg">
                  <p className="text-sm font-semibold">{productName}</p>
                  <p className="text-xs text-green-400">AR Mode Active</p>
                </div>

                {/* AR Controls */}
                <div className="absolute bottom-8 left-1/2 transform -translate-x-1/2 flex gap-4">
                  <button className="bg-white bg-opacity-20 backdrop-blur-sm text-white px-6 py-3 rounded-full hover:bg-opacity-30 transition">
                    📸 写真撮影
                  </button>
                  <button className="bg-white bg-opacity-20 backdrop-blur-sm text-white px-6 py-3 rounded-full hover:bg-opacity-30 transition">
                    🔄 回転
                  </button>
                  <button className="bg-white bg-opacity-20 backdrop-blur-sm text-white px-6 py-3 rounded-full hover:bg-opacity-30 transition">
                    📏 サイズ調整
                  </button>
                </div>
              </div>
            </div>
          </div>

          {/* Close button */}
          <button
            onClick={handleClose}
            className="absolute top-4 right-4 bg-red-500 text-white px-6 py-3 rounded-full hover:bg-red-600 transition z-10"
          >
            ✕ 閉じる
          </button>
        </div>
      </div>
    )
  }

  return (
    <Card>
      <div className="space-y-4">
        <h3 className="font-semibold text-lg">🥽 3D / AR プレビュー</h3>

        {/* 3D Preview */}
        <div
          ref={containerRef}
          className="relative bg-gradient-to-br from-gray-100 to-gray-200 rounded-lg overflow-hidden"
          style={{ height: '300px' }}
        >
          <div className="absolute inset-0 flex items-center justify-center">
            <div
              className="transform-gpu transition-transform"
              style={{
                transform: `perspective(1000px) rotateY(${rotation}deg) rotateX(10deg)`,
              }}
            >
              {imageUrl ? (
                <img
                  src={imageUrl}
                  alt={productName}
                  className="w-48 h-48 object-contain rounded-lg shadow-xl"
                />
              ) : (
                <div className="w-48 h-48 bg-gradient-to-br from-primary-300 to-primary-500 rounded-lg shadow-xl flex items-center justify-center">
                  <span className="text-6xl">📦</span>
                </div>
              )}
            </div>
          </div>

          {/* Rotation indicator */}
          <div className="absolute bottom-4 right-4 bg-black bg-opacity-50 text-white px-3 py-1 rounded-full text-sm">
            🔄 回転中
          </div>
        </div>

        {/* AR Button */}
        <Button
          onClick={handleARView}
          isLoading={isLoading}
          className="w-full bg-gradient-to-r from-green-500 to-emerald-600 hover:from-green-600 hover:to-emerald-700"
        >
          🥽 AR で見る
        </Button>

        <div className="text-xs text-gray-500 text-center">
          ARモードで実際の空間に商品を配置してサイズ感を確認できます
        </div>
      </div>
    </Card>
  )
}

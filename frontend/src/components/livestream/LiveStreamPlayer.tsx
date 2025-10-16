import { useState, useEffect, useRef } from 'react'
import { Card } from '@/components/common/Card'
import { Button } from '@/components/common/Button'

interface LiveStreamPlayerProps {
  streamId: string
  title: string
  sellerName: string
  viewerCount: number
  onClose?: () => void
}

interface Comment {
  id: string
  user_name: string
  comment: string
  created_at: string
}

export const LiveStreamPlayer = ({
  streamId,
  title,
  sellerName,
  viewerCount,
  onClose,
}: LiveStreamPlayerProps) => {
  const [comments, setComments] = useState<Comment[]>([])
  const [newComment, setNewComment] = useState('')
  const [isLive] = useState(true)
  const commentsEndRef = useRef<HTMLDivElement>(null)

  // Simulate live comments
  useEffect(() => {
    const interval = setInterval(() => {
      const sampleComments = [
        'すごい商品ですね！',
        '値段はいくらですか？',
        '購入したいです！',
        '状態はどうですか？',
        'カッコいい！',
        'この色好きです',
      ]

      const randomComment = sampleComments[Math.floor(Math.random() * sampleComments.length)]

      setComments((prev) => [
        ...prev.slice(-50), // Keep last 50 comments
        {
          id: Date.now().toString(),
          user_name: `ユーザー${Math.floor(Math.random() * 100)}`,
          comment: randomComment,
          created_at: new Date().toISOString(),
        },
      ])
    }, 3000)

    return () => clearInterval(interval)
  }, [])

  useEffect(() => {
    commentsEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [comments])

  const handleSendComment = (e: React.FormEvent) => {
    e.preventDefault()
    if (!newComment.trim()) return

    setComments((prev) => [
      ...prev,
      {
        id: Date.now().toString(),
        user_name: 'あなた',
        comment: newComment,
        created_at: new Date().toISOString(),
      },
    ])
    setNewComment('')
  }

  return (
    <div className="fixed inset-0 bg-black z-50 flex flex-col">
      {/* Header */}
      <div className="bg-gradient-to-r from-red-600 to-pink-600 text-white px-4 py-3 flex items-center justify-between">
        <div className="flex items-center gap-4">
          <div className="flex items-center gap-2">
            <div className="w-3 h-3 bg-red-500 rounded-full animate-pulse" />
            <span className="font-bold">LIVE</span>
          </div>
          <div className="text-sm">
            <p className="font-semibold">{title}</p>
            <p className="text-xs opacity-90">{sellerName}</p>
          </div>
        </div>

        <div className="flex items-center gap-4">
          <div className="flex items-center gap-2 bg-black bg-opacity-30 px-3 py-1 rounded-full">
            <span className="text-sm">👁️</span>
            <span className="font-semibold">{viewerCount}</span>
          </div>
          {onClose && (
            <button
              onClick={onClose}
              className="bg-white bg-opacity-20 hover:bg-opacity-30 px-4 py-1 rounded-full transition"
            >
              ✕
            </button>
          )}
        </div>
      </div>

      {/* Main content */}
      <div className="flex-1 flex overflow-hidden">
        {/* Video Player */}
        <div className="flex-1 bg-gray-900 flex items-center justify-center relative">
          {/* Simulated video stream */}
          <div className="absolute inset-0 bg-gradient-to-br from-purple-900 via-blue-900 to-pink-900 opacity-50" />

          {/* Animated background */}
          <div className="absolute inset-0">
            {Array.from({ length: 20 }).map((_, i) => (
              <div
                key={i}
                className="absolute bg-white rounded-full opacity-10"
                style={{
                  width: `${Math.random() * 100 + 50}px`,
                  height: `${Math.random() * 100 + 50}px`,
                  left: `${Math.random() * 100}%`,
                  top: `${Math.random() * 100}%`,
                  animation: `float ${Math.random() * 10 + 10}s infinite ease-in-out`,
                }}
              />
            ))}
          </div>

          {/* Mock seller video feed */}
          <div className="relative z-10 text-center">
            <div className="w-64 h-64 bg-gradient-to-br from-primary-400 to-primary-600 rounded-full flex items-center justify-center mb-4 shadow-2xl">
              <span className="text-9xl">📦</span>
            </div>
            <p className="text-white text-2xl font-bold">{sellerName}のライブ配信</p>
            <p className="text-gray-300 mt-2">商品を実際に見ながらご質問いただけます</p>
          </div>

          {/* Live indicators */}
          <div className="absolute top-4 left-4 flex gap-2">
            <div className="bg-red-500 text-white px-3 py-1 rounded-full text-sm font-bold flex items-center gap-2">
              <div className="w-2 h-2 bg-white rounded-full animate-pulse" />
              配信中
            </div>
          </div>
        </div>

        {/* Chat sidebar */}
        <div className="w-96 bg-gray-100 flex flex-col">
          <div className="bg-white border-b px-4 py-3">
            <h3 className="font-semibold">💬 チャット</h3>
          </div>

          {/* Comments */}
          <div className="flex-1 overflow-y-auto p-4 space-y-3">
            {comments.map((comment) => (
              <div
                key={comment.id}
                className="bg-white rounded-lg p-3 shadow-sm hover:shadow-md transition"
              >
                <div className="flex items-center gap-2 mb-1">
                  <span className="font-semibold text-sm text-primary-600">
                    {comment.user_name}
                  </span>
                  <span className="text-xs text-gray-400">
                    {new Date(comment.created_at).toLocaleTimeString('ja-JP', {
                      hour: '2-digit',
                      minute: '2-digit',
                    })}
                  </span>
                </div>
                <p className="text-sm text-gray-800">{comment.comment}</p>
              </div>
            ))}
            <div ref={commentsEndRef} />
          </div>

          {/* Comment input */}
          <form onSubmit={handleSendComment} className="p-4 bg-white border-t">
            <div className="flex gap-2">
              <input
                type="text"
                value={newComment}
                onChange={(e) => setNewComment(e.target.value)}
                placeholder="コメントを入力..."
                className="flex-1 px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
              />
              <Button type="submit">送信</Button>
            </div>
          </form>
        </div>
      </div>

      <style>{`
        @keyframes float {
          0%, 100% {
            transform: translateY(0) translateX(0);
          }
          50% {
            transform: translateY(-20px) translateX(10px);
          }
        }
      `}</style>
    </div>
  )
}

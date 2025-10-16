import { useState, useEffect } from 'react'
import { Card } from '@/components/common/Card'
import { Button } from '@/components/common/Button'
import axios from 'axios'
import toast from 'react-hot-toast'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

interface CO2Goal {
  id: string
  user_id: string
  target_kg: number
  current_kg: number
  target_date: string
  start_date: string
  status: string
}

export const CO2GoalCard = () => {
  const [goal, setGoal] = useState<CO2Goal | null>(null)
  const [progress, setProgress] = useState(0)
  const [loading, setLoading] = useState(true)
  const [showCreateForm, setShowCreateForm] = useState(false)
  const [targetKG, setTargetKG] = useState(10)
  const [targetDate, setTargetDate] = useState('')
  const [submitting, setSubmitting] = useState(false)

  useEffect(() => {
    loadGoal()
  }, [])

  const loadGoal = async () => {
    try {
      setLoading(true)
      const token = localStorage.getItem('token')
      const response = await axios.get(`${API_BASE_URL}/v1/co2-goals`, {
        headers: { Authorization: `Bearer ${token}` },
      })

      setGoal(response.data.goal)
      setProgress(response.data.progress || 0)
    } catch (error: any) {
      if (error.response?.status !== 401) {
        console.error('Failed to load goal:', error)
      }
    } finally {
      setLoading(false)
    }
  }

  const handleCreateGoal = async (e: React.FormEvent) => {
    e.preventDefault()

    try {
      setSubmitting(true)
      const token = localStorage.getItem('token')

      await axios.post(
        `${API_BASE_URL}/v1/co2-goals`,
        {
          target_kg: targetKG,
          target_date: targetDate,
        },
        {
          headers: { Authorization: `Bearer ${token}` },
        }
      )

      toast.success('CO2削減目標を設定しました！')
      setShowCreateForm(false)
      loadGoal()
    } catch (error: any) {
      toast.error(error.response?.data?.error || '目標の設定に失敗しました')
    } finally {
      setSubmitting(false)
    }
  }

  if (loading) {
    return (
      <Card>
        <div className="flex items-center justify-center py-8">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-500" />
        </div>
      </Card>
    )
  }

  if (!goal && !showCreateForm) {
    return (
      <Card>
        <div className="text-center py-6">
          <div className="text-4xl mb-3">🎯</div>
          <h3 className="text-lg font-semibold mb-2">CO2削減目標を設定しましょう！</h3>
          <p className="text-gray-600 text-sm mb-4">
            目標を設定して、エコな買い物を習慣化しましょう
          </p>
          <Button onClick={() => setShowCreateForm(true)}>目標を設定する</Button>
        </div>
      </Card>
    )
  }

  if (showCreateForm) {
    return (
      <Card>
        <h3 className="text-lg font-semibold mb-4">CO2削減目標を設定</h3>
        <form onSubmit={handleCreateGoal} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              目標削減量 (kg)
            </label>
            <input
              type="number"
              value={targetKG}
              onChange={(e) => setTargetKG(Number(e.target.value))}
              min="1"
              step="0.1"
              required
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
            />
            <p className="text-xs text-gray-500 mt-1">
              例: 10kg = 約50〜100回の買い物
            </p>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              達成目標日
            </label>
            <input
              type="date"
              value={targetDate}
              onChange={(e) => setTargetDate(e.target.value)}
              min={new Date().toISOString().split('T')[0]}
              required
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
            />
          </div>

          <div className="flex gap-2">
            <Button type="submit" isLoading={submitting}>
              設定する
            </Button>
            <Button
              type="button"
              variant="outline"
              onClick={() => setShowCreateForm(false)}
            >
              キャンセル
            </Button>
          </div>
        </form>
      </Card>
    )
  }

  const daysLeft = Math.ceil(
    (new Date(goal!.target_date).getTime() - new Date().getTime()) / (1000 * 60 * 60 * 24)
  )

  const statusColor =
    goal!.status === 'completed'
      ? 'bg-green-500'
      : goal!.status === 'expired'
      ? 'bg-red-500'
      : 'bg-primary-500'

  const statusText =
    goal!.status === 'completed'
      ? '達成！🎉'
      : goal!.status === 'expired'
      ? '期限切れ'
      : `あと${daysLeft}日`

  return (
    <Card>
      <div className="space-y-4">
        <div className="flex items-center justify-between">
          <h3 className="text-lg font-semibold">🎯 CO2削減目標</h3>
          <span className={`px-3 py-1 rounded-full text-white text-sm ${statusColor}`}>
            {statusText}
          </span>
        </div>

        <div className="space-y-2">
          <div className="flex items-center justify-between text-sm">
            <span className="text-gray-600">進捗</span>
            <span className="font-semibold">
              {goal!.current_kg.toFixed(1)} / {goal!.target_kg} kg
            </span>
          </div>

          <div className="w-full bg-gray-200 rounded-full h-4">
            <div
              className={`h-4 rounded-full transition-all duration-500 ${
                progress >= 100 ? 'bg-green-500' : statusColor
              }`}
              style={{ width: `${Math.min(progress, 100)}%` }}
            />
          </div>

          <p className="text-xs text-gray-500 text-center">{progress.toFixed(1)}% 達成</p>
        </div>

        <div className="grid grid-cols-2 gap-4 pt-4 border-t">
          <div>
            <p className="text-xs text-gray-500">開始日</p>
            <p className="text-sm font-medium">
              {new Date(goal!.start_date).toLocaleDateString('ja-JP')}
            </p>
          </div>
          <div>
            <p className="text-xs text-gray-500">目標日</p>
            <p className="text-sm font-medium">
              {new Date(goal!.target_date).toLocaleDateString('ja-JP')}
            </p>
          </div>
        </div>

        {goal!.status === 'completed' && (
          <div className="bg-green-50 border border-green-200 rounded-lg p-3 text-center">
            <p className="text-green-800 font-semibold">🎊 目標達成おめでとうございます！</p>
            <p className="text-sm text-green-600 mt-1">
              地球環境のために素晴らしい貢献をしました
            </p>
          </div>
        )}

        {goal!.status === 'active' && progress < 100 && (
          <div className="bg-primary-50 border border-primary-200 rounded-lg p-3">
            <p className="text-sm text-primary-800">
              💪 あと{(goal!.target_kg - goal!.current_kg).toFixed(1)}kg削減で目標達成です！
            </p>
          </div>
        )}
      </div>
    </Card>
  )
}

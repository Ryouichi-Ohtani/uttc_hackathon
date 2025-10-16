import { useState } from 'react'
import { reviewService } from '../../services/reviews'
import { Button } from '../common/Button'
import toast from 'react-hot-toast'

interface CreateReviewProps {
  productId: string
  purchaseId: string
  onReviewCreated?: () => void
}

export const CreateReview = ({ productId, purchaseId, onReviewCreated }: CreateReviewProps) => {
  const [rating, setRating] = useState(5)
  const [comment, setComment] = useState('')
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    try {
      setLoading(true)
      await reviewService.create({
        product_id: productId,
        purchase_id: purchaseId,
        rating,
        comment
      })
      toast.success('Review submitted successfully!')
      setComment('')
      setRating(5)
      onReviewCreated?.()
    } catch (error: any) {
      toast.error(error.response?.data?.error || 'Failed to submit review')
    } finally {
      setLoading(false)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div>
        <label className="block text-sm font-medium mb-2">Rating</label>
        <div className="flex space-x-2">
          {[1, 2, 3, 4, 5].map((star) => (
            <button
              key={star}
              type="button"
              onClick={() => setRating(star)}
              className={`text-3xl ${
                star <= rating ? 'text-yellow-400' : 'text-gray-300'
              } hover:text-yellow-400 transition-colors`}
            >
              ★
            </button>
          ))}
        </div>
      </div>

      <div>
        <label className="block text-sm font-medium mb-2">Comment</label>
        <textarea
          value={comment}
          onChange={(e) => setComment(e.target.value)}
          rows={4}
          className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-green-500 focus:border-transparent"
          placeholder="Share your experience with this product..."
        />
      </div>

      <Button type="submit" disabled={loading} className="w-full">
        {loading ? 'Submitting...' : 'Submit Review'}
      </Button>
    </form>
  )
}

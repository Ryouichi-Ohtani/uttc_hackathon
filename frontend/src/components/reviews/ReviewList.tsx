import { Review } from '../../types'
import { Card } from '../common/Card'

interface ReviewListProps {
  reviews: Review[]
  averageRating: number
}

export const ReviewList = ({ reviews, averageRating }: ReviewListProps) => {
  if (reviews.length === 0) {
    return (
      <Card padding="lg">
        <p className="text-center text-gray-500">No reviews yet. Be the first to review!</p>
      </Card>
    )
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center space-x-4 mb-6">
        <div className="text-4xl font-bold">{averageRating.toFixed(1)}</div>
        <div>
          <div className="flex text-yellow-400">
            {[1, 2, 3, 4, 5].map((star) => (
              <span key={star} className="text-xl">
                {star <= Math.round(averageRating) ? '★' : '☆'}
              </span>
            ))}
          </div>
          <p className="text-sm text-gray-600">{reviews.length} reviews</p>
        </div>
      </div>

      {reviews.map((review) => (
        <Card key={review.id} padding="md">
          <div className="flex items-start justify-between mb-2">
            <div className="flex items-center space-x-3">
              {review.reviewer?.avatar_url ? (
                <img
                  src={review.reviewer.avatar_url}
                  alt={review.reviewer.display_name}
                  className="w-10 h-10 rounded-full"
                />
              ) : (
                <div className="w-10 h-10 rounded-full bg-gray-200 flex items-center justify-center">
                  {review.reviewer?.display_name?.charAt(0) || '?'}
                </div>
              )}
              <div>
                <p className="font-semibold">{review.reviewer?.display_name || 'Anonymous'}</p>
                <div className="flex text-yellow-400 text-sm">
                  {[1, 2, 3, 4, 5].map((star) => (
                    <span key={star}>{star <= review.rating ? '★' : '☆'}</span>
                  ))}
                </div>
              </div>
            </div>
            <p className="text-xs text-gray-400">
              {new Date(review.created_at).toLocaleDateString()}
            </p>
          </div>
          <p className="text-gray-700">{review.comment}</p>
        </Card>
      ))}
    </div>
  )
}

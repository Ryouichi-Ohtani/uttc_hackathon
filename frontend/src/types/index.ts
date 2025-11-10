export interface User {
  id: string
  email: string
  username: string
  display_name: string
  avatar_url?: string
  bio?: string
  role: string
  // Address information for shipping
  postal_code?: string
  prefecture?: string
  city?: string
  address_line1?: string
  address_line2?: string
  phone_number?: string
  created_at: string
  updated_at: string
}

export interface Product {
  id: string
  seller_id: string
  seller?: User
  title: string
  description: string
  price: number
  category: string
  condition: 'new' | 'like_new' | 'good' | 'fair' | 'poor'
  status: 'draft' | 'active' | 'sold' | 'reserved' | 'deleted'
  weight_kg?: number
  manufacturer_country?: string
  estimated_manufacturing_year?: number
  ai_generated_description?: string
  ai_suggested_price?: number
  co2_impact_kg?: number
  view_count: number
  favorite_count: number
  has_3d_model: boolean
  model_url?: string
  images: ProductImage[]
  created_at: string
  updated_at: string
  is_favorited?: boolean
}

export interface ProductImage {
  id: string
  product_id: string
  image_url: string
  cdn_url?: string
  display_order: number
  is_primary: boolean
  width?: number
  height?: number
  created_at: string
}

export interface Purchase {
  id: string
  product_id: string
  product?: Product
  buyer_id: string
  buyer?: User
  seller_id: string
  seller?: User
  price: number
  status: 'pending' | 'completed' | 'cancelled'
  payment_method?: string
  shipping_address?: string
  // Delivery information
  delivery_date?: string
  delivery_time_slot?: 'morning' | 'afternoon' | 'evening' | 'anytime'
  use_registered_address?: boolean
  recipient_name?: string
  recipient_phone_number?: string
  recipient_postal_code?: string
  recipient_prefecture?: string
  recipient_city?: string
  recipient_address_line1?: string
  recipient_address_line2?: string
  completed_at?: string
  created_at: string
}

export interface ShippingLabel {
  id: string
  purchase_id: string
  purchase?: Purchase
  // Sender (Seller) information
  sender_name: string
  sender_postal_code: string
  sender_prefecture: string
  sender_city: string
  sender_address_line1: string
  sender_address_line2?: string
  sender_phone_number: string
  // Recipient (Buyer) information
  recipient_name: string
  recipient_postal_code: string
  recipient_prefecture: string
  recipient_city: string
  recipient_address_line1: string
  recipient_address_line2?: string
  recipient_phone_number: string
  // Delivery details
  delivery_date?: string
  delivery_time_slot?: string
  product_name: string
  package_size: string
  weight: number
  tracking_number?: string
  carrier: string
  created_at: string
  updated_at: string
}

export interface Conversation {
  id: string
  product_id?: string
  product?: Product
  participants: ConversationParticipant[]
  last_message?: Message
  last_message_at: string
  unread_count?: number
  created_at: string
}

export interface ConversationParticipant {
  id: string
  conversation_id: string
  user_id: string
  user?: User
  last_read_at: string
  created_at: string
}

export interface Message {
  id: string
  conversation_id: string
  sender_id: string
  sender?: User
  content: string
  is_read: boolean
  created_at: string
}

export interface Achievement {
  id: string
  name: string
  description: string
  badge_icon_url?: string
  requirement_type: string
  requirement_value: number
  created_at: string
}

export interface UserAchievement {
  id: string
  user_id: string
  achievement_id: string
  achievement?: Achievement
  earned_at: string
}

// Deprecated - kept for compatibility
export interface SustainabilityLog {
  id: string
  user_id: string
  purchase_id?: string
  action_type: 'purchase' | 'sale'
  description: string
  created_at: string
}

export interface DashboardData {
  achievements: UserAchievement[]
  recent_logs: SustainabilityLog[]
  monthly_stats: {
    transactions: number
  }
}

export interface LeaderboardEntry {
  rank: number
  user: User
}

export interface PaginationResponse {
  page: number
  limit: number
  total: number
  total_pages: number
}

export interface ProductFilters {
  category?: string
  min_price?: number
  max_price?: number
  condition?: string
  search?: string
  sort?: string
  page?: number
  limit?: number
  ai_generated?: boolean
}

export interface CO2Comparison {
  product_co2: number
  new_product_co2: number
  saved_co2: number
  saved_percentage: number
}

export interface AuthResponse {
  user: User
  token: string
}

export interface WSMessage {
  type: 'auth' | 'message' | 'typing' | 'read' | 'send_message' | 'mark_read'
  data?: any
  token?: string
  content?: string
  message_id?: string
}

export interface Review {
  id: string
  product_id: string
  purchase_id: string
  reviewer_id: string
  reviewer?: User
  rating: number
  comment: string
  created_at: string
  updated_at: string
}

// Auto-Purchase types
export interface PaymentAuthorizationRequest {
  card_number: string
  expiry_month: number
  expiry_year: number
  cvv: string
  cardholder_name: string
  amount: number
}

export interface PaymentAuthorizationResponse {
  payment_method_id: string
  payment_auth_token: string
  authorized_amount: number
  expires_at: string
  last4: string
}

export interface AutoPurchaseWatch {
  id: string
  user_id: string
  product_id: string
  product?: Product
  max_price: number
  status: 'active' | 'executed' | 'cancelled' | 'expired'
  payment_authorized: boolean
  payment_method_id: string
  payment_auth_token: string
  authorized_amount: number
  delivery_date?: string
  delivery_time_slot?: 'morning' | 'afternoon' | 'evening' | 'anytime'
  shipping_address: string
  recipient_name: string
  recipient_phone_number: string
  recipient_postal_code: string
  recipient_prefecture: string
  recipient_city: string
  recipient_address_line1: string
  recipient_address_line2?: string
  expires_at: string
  executed_at?: string
  purchase_id?: string
  created_at: string
  updated_at: string
}

export interface CreateAutoPurchaseWatchRequest {
  product_id: string
  max_price: number
  payment_method_id: string
  payment_auth_token: string
  authorized_amount: number
  delivery_date?: string
  delivery_time_slot?: 'morning' | 'afternoon' | 'evening' | 'anytime'
  shipping_address: string
  recipient_name: string
  recipient_phone_number: string
  recipient_postal_code: string
  recipient_prefecture: string
  recipient_city: string
  recipient_address_line1: string
  recipient_address_line2?: string
}

# API Specification - EcoMate

## Base URLs

- **Go Backend (REST)**: `https://api.ecomate.example.com/v1`
- **Go Backend (WebSocket)**: `wss://api.ecomate.example.com/v1/ws`
- **Python AI Service (gRPC)**: Internal only, called by Go backend

## Authentication

All authenticated endpoints require JWT token in header:
```
Authorization: Bearer <jwt_token>
```

---

## REST API Endpoints

### Authentication

#### POST /auth/register
Register a new user.

**Request:**
```json
{
  "email": "user@example.com",
  "username": "johndoe",
  "password": "securepassword123",
  "display_name": "John Doe"
}
```

**Response (201):**
```json
{
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "username": "johndoe",
    "display_name": "John Doe",
    "sustainability_score": 0,
    "total_co2_saved_kg": 0.0,
    "level": 1
  },
  "token": "jwt_token_here"
}
```

#### POST /auth/login
Login existing user.

**Request:**
```json
{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Response (200):**
```json
{
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "username": "johndoe",
    "display_name": "John Doe",
    "avatar_url": "https://cdn.example.com/avatars/123.jpg",
    "sustainability_score": 150,
    "total_co2_saved_kg": 25.5,
    "level": 3
  },
  "token": "jwt_token_here"
}
```

#### GET /auth/me
Get current user profile.

**Response (200):**
```json
{
  "id": "uuid",
  "email": "user@example.com",
  "username": "johndoe",
  "display_name": "John Doe",
  "avatar_url": "https://cdn.example.com/avatars/123.jpg",
  "bio": "Eco-conscious buyer",
  "sustainability_score": 150,
  "total_co2_saved_kg": 25.5,
  "level": 3,
  "achievements": [
    {
      "id": "uuid",
      "name": "Eco Warrior",
      "description": "Save 10kg of CO2",
      "badge_icon_url": "https://cdn.example.com/badges/eco-warrior.svg",
      "earned_at": "2024-01-15T10:30:00Z"
    }
  ]
}
```

---

### Products

#### POST /products
Create a new product listing (with AI assistance).

**Request (multipart/form-data):**
```
title: "Vintage Uniqlo Fleece Jacket"
description: "" (optional, AI will generate if empty)
price: 800 (in yen/cents)
category: "clothing"
condition: "good"
weight_kg: 0.5 (optional, AI will estimate)
images: [File, File, ...] (up to 10)
use_ai_assistance: true
```

**Response (201):**
```json
{
  "id": "uuid",
  "seller_id": "uuid",
  "title": "Vintage Uniqlo Fleece Jacket",
  "description": "Cozy fleece jacket from Uniqlo in excellent condition. Perfect for autumn weather...",
  "ai_generated_description": true,
  "price": 800,
  "ai_suggested_price": 850,
  "category": "clothing",
  "condition": "good",
  "status": "active",
  "weight_kg": 0.5,
  "manufacturer_country": "China",
  "estimated_manufacturing_year": 2020,
  "co2_impact_kg": 5.2,
  "images": [
    {
      "id": "uuid",
      "image_url": "https://storage.example.com/products/abc123.jpg",
      "cdn_url": "https://cdn.example.com/products/abc123.jpg",
      "is_primary": true
    }
  ],
  "created_at": "2024-01-20T15:00:00Z"
}
```

#### GET /products
List products with filters and pagination.

**Query Parameters:**
- `category` (optional): Filter by category
- `min_price`, `max_price` (optional): Price range
- `condition` (optional): Filter by condition
- `search` (optional): Full-text search
- `sort` (optional): `price_asc`, `price_desc`, `created_desc` (default), `eco_impact_desc`
- `page` (default: 1)
- `limit` (default: 20, max: 100)

**Response (200):**
```json
{
  "products": [
    {
      "id": "uuid",
      "seller": {
        "id": "uuid",
        "username": "johndoe",
        "display_name": "John Doe",
        "avatar_url": "https://cdn.example.com/avatars/123.jpg"
      },
      "title": "Vintage Uniqlo Fleece Jacket",
      "price": 800,
      "category": "clothing",
      "condition": "good",
      "status": "active",
      "co2_impact_kg": 5.2,
      "primary_image": {
        "cdn_url": "https://cdn.example.com/products/abc123.jpg"
      },
      "favorite_count": 12,
      "created_at": "2024-01-20T15:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 150,
    "total_pages": 8
  }
}
```

#### GET /products/:id
Get product details.

**Response (200):**
```json
{
  "id": "uuid",
  "seller": {
    "id": "uuid",
    "username": "johndoe",
    "display_name": "John Doe",
    "avatar_url": "https://cdn.example.com/avatars/123.jpg",
    "sustainability_score": 150
  },
  "title": "Vintage Uniqlo Fleece Jacket",
  "description": "Cozy fleece jacket...",
  "price": 800,
  "category": "clothing",
  "condition": "good",
  "status": "active",
  "weight_kg": 0.5,
  "manufacturer_country": "China",
  "co2_impact_kg": 5.2,
  "view_count": 45,
  "favorite_count": 12,
  "is_favorited": false,
  "has_3d_model": false,
  "images": [
    {
      "id": "uuid",
      "cdn_url": "https://cdn.example.com/products/abc123.jpg",
      "is_primary": true
    }
  ],
  "co2_comparison": {
    "buying_new_kg": 15.8,
    "buying_used_kg": 10.6,
    "saved_kg": 5.2,
    "equivalent_trees": 0.26
  },
  "created_at": "2024-01-20T15:00:00Z"
}
```

#### POST /products/:id/favorite
Add product to favorites.

**Response (200):**
```json
{
  "favorited": true,
  "favorite_count": 13
}
```

#### DELETE /products/:id/favorite
Remove product from favorites.

**Response (200):**
```json
{
  "favorited": false,
  "favorite_count": 12
}
```

---

### Purchases

#### POST /purchases
Create a purchase (buy a product).

**Request:**
```json
{
  "product_id": "uuid",
  "shipping_address": "123 Main St, Tokyo, Japan",
  "payment_method": "credit_card"
}
```

**Response (201):**
```json
{
  "id": "uuid",
  "product": {
    "id": "uuid",
    "title": "Vintage Uniqlo Fleece Jacket",
    "price": 800
  },
  "buyer_id": "uuid",
  "seller_id": "uuid",
  "price": 800,
  "co2_saved_kg": 5.2,
  "status": "pending",
  "created_at": "2024-01-21T10:00:00Z"
}
```

#### GET /purchases
Get user's purchase history.

**Query Parameters:**
- `role` (optional): `buyer` or `seller`
- `page`, `limit`

**Response (200):**
```json
{
  "purchases": [
    {
      "id": "uuid",
      "product": {
        "id": "uuid",
        "title": "Vintage Uniqlo Fleece Jacket",
        "primary_image_url": "https://cdn.example.com/products/abc123.jpg"
      },
      "buyer": { "id": "uuid", "username": "janedoe" },
      "seller": { "id": "uuid", "username": "johndoe" },
      "price": 800,
      "co2_saved_kg": 5.2,
      "status": "completed",
      "created_at": "2024-01-21T10:00:00Z",
      "completed_at": "2024-01-25T14:30:00Z"
    }
  ],
  "pagination": { "page": 1, "limit": 20, "total": 5 }
}
```

#### PATCH /purchases/:id/complete
Mark purchase as completed (seller only).

**Response (200):**
```json
{
  "id": "uuid",
  "status": "completed",
  "completed_at": "2024-01-25T14:30:00Z"
}
```

---

### Messaging (DM)

#### GET /conversations
Get user's conversations.

**Response (200):**
```json
{
  "conversations": [
    {
      "id": "uuid",
      "product": {
        "id": "uuid",
        "title": "Vintage Uniqlo Fleece Jacket",
        "primary_image_url": "https://cdn.example.com/products/abc123.jpg"
      },
      "participants": [
        {
          "id": "uuid",
          "username": "johndoe",
          "avatar_url": "https://cdn.example.com/avatars/123.jpg"
        }
      ],
      "last_message": {
        "content": "Is this still available?",
        "created_at": "2024-01-21T09:00:00Z"
      },
      "unread_count": 2
    }
  ]
}
```

#### POST /conversations
Create or get conversation.

**Request:**
```json
{
  "product_id": "uuid",
  "participant_id": "uuid"
}
```

**Response (201 or 200):**
```json
{
  "id": "uuid",
  "product_id": "uuid",
  "participants": [...]
}
```

#### GET /conversations/:id/messages
Get messages in a conversation.

**Query Parameters:**
- `page`, `limit`
- `before` (timestamp): Get messages before this time

**Response (200):**
```json
{
  "messages": [
    {
      "id": "uuid",
      "sender": {
        "id": "uuid",
        "username": "johndoe",
        "avatar_url": "https://cdn.example.com/avatars/123.jpg"
      },
      "content": "Yes, it's still available!",
      "is_read": true,
      "created_at": "2024-01-21T09:05:00Z"
    }
  ],
  "pagination": { "page": 1, "limit": 50, "total": 8 }
}
```

#### POST /conversations/:id/messages
Send a message (also available via WebSocket).

**Request:**
```json
{
  "content": "Is this still available?"
}
```

**Response (201):**
```json
{
  "id": "uuid",
  "sender_id": "uuid",
  "content": "Is this still available?",
  "created_at": "2024-01-21T09:00:00Z"
}
```

---

### Sustainability

#### GET /sustainability/dashboard
Get user's sustainability dashboard.

**Response (200):**
```json
{
  "total_co2_saved_kg": 25.5,
  "level": 3,
  "sustainability_score": 150,
  "next_level_threshold": 200,
  "achievements": [
    {
      "id": "uuid",
      "name": "Eco Warrior",
      "badge_icon_url": "https://cdn.example.com/badges/eco-warrior.svg",
      "earned_at": "2024-01-15T10:30:00Z"
    }
  ],
  "recent_logs": [
    {
      "id": "uuid",
      "action_type": "purchase",
      "co2_saved_kg": 5.2,
      "description": "Purchased Vintage Uniqlo Fleece Jacket",
      "created_at": "2024-01-21T10:00:00Z"
    }
  ],
  "monthly_stats": {
    "current_month_co2_saved": 12.3,
    "transactions": 4
  },
  "comparisons": {
    "equivalent_trees": 1.3,
    "car_km_avoided": 85.5
  }
}
```

#### GET /sustainability/leaderboard
Get top users by sustainability score.

**Query Parameters:**
- `period` (optional): `week`, `month`, `all_time` (default)
- `limit` (default: 50)

**Response (200):**
```json
{
  "leaderboard": [
    {
      "rank": 1,
      "user": {
        "id": "uuid",
        "username": "eco_master",
        "display_name": "Eco Master",
        "avatar_url": "https://cdn.example.com/avatars/456.jpg"
      },
      "total_co2_saved_kg": 250.5,
      "sustainability_score": 1500,
      "level": 15
    }
  ],
  "current_user_rank": 42
}
```

---

## WebSocket API

### WS /ws/conversations/:conversation_id

Connect to real-time messaging.

**Authentication:**
Send JWT token as first message:
```json
{
  "type": "auth",
  "token": "jwt_token_here"
}
```

**Incoming Message Types:**

1. Message received:
```json
{
  "type": "message",
  "data": {
    "id": "uuid",
    "sender": {
      "id": "uuid",
      "username": "johndoe"
    },
    "content": "Hello!",
    "created_at": "2024-01-21T10:00:00Z"
  }
}
```

2. User typing:
```json
{
  "type": "typing",
  "data": {
    "user_id": "uuid",
    "username": "johndoe"
  }
}
```

3. Message read:
```json
{
  "type": "read",
  "data": {
    "message_id": "uuid",
    "user_id": "uuid"
  }
}
```

**Outgoing Message Types:**

1. Send message:
```json
{
  "type": "send_message",
  "content": "Hello!"
}
```

2. Typing indicator:
```json
{
  "type": "typing"
}
```

3. Mark as read:
```json
{
  "type": "mark_read",
  "message_id": "uuid"
}
```

---

## gRPC Service (Python AI ↔ Go Backend)

### ProductAnalysisService

#### AnalyzeProduct
Multimodal analysis of product images and text.

**Request:**
```protobuf
message AnalyzeProductRequest {
  repeated bytes images = 1;
  string title = 2;
  string user_provided_description = 3;
  string category = 4;
}
```

**Response:**
```protobuf
message AnalyzeProductResponse {
  string generated_description = 1;
  int32 suggested_price = 2;
  float estimated_weight_kg = 3;
  string manufacturer_country = 4;
  int32 estimated_manufacturing_year = 5;
  float co2_impact_kg = 6;
  bool is_inappropriate = 7;
  string inappropriate_reason = 8;
  repeated string detected_objects = 9;
}
```

#### CalculateCO2Impact
Calculate CO2 savings.

**Request:**
```protobuf
message CalculateCO2Request {
  string category = 1;
  float weight_kg = 2;
  string manufacturer_country = 3;
  int32 manufacturing_year = 4;
}
```

**Response:**
```protobuf
message CalculateCO2Response {
  float buying_new_kg = 1;
  float buying_used_kg = 2;
  float saved_kg = 3;
}
```

---

## Error Responses

All errors follow this format:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid email format",
    "details": {
      "field": "email"
    }
  }
}
```

**Common Error Codes:**
- `VALIDATION_ERROR` (400)
- `UNAUTHORIZED` (401)
- `FORBIDDEN` (403)
- `NOT_FOUND` (404)
- `CONFLICT` (409) - e.g., email already exists
- `INTERNAL_ERROR` (500)

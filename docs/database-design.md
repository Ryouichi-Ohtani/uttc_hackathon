# Database Schema Design - EcoMate

## ER Diagram (Text Format)

```
users ||--o{ products : "sells"
users ||--o{ purchases : "buys"
users ||--o{ messages : "sends"
users ||--o{ user_achievements : "earns"
users ||--o{ sustainability_logs : "logs"

products ||--o{ purchases : "purchased_in"
products ||--o{ product_images : "has"
products ||--o{ favorites : "favorited_in"

conversations ||--o{ messages : "contains"
users ||--o{ conversation_participants : "participates_in"
conversations ||--o{ conversation_participants : "has"

achievements ||--o{ user_achievements : "awarded_to"
```

## Tables

### users
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    display_name VARCHAR(100),
    avatar_url TEXT,
    bio TEXT,
    sustainability_score INTEGER DEFAULT 0,
    total_co2_saved_kg DECIMAL(10, 2) DEFAULT 0.00,
    level INTEGER DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_sustainability_score ON users(sustainability_score DESC);
```

### products
```sql
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    seller_id UUID NOT NULL REFERENCES users(id),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    price INTEGER NOT NULL, -- in cents/yen
    category VARCHAR(100) NOT NULL,
    condition VARCHAR(50) NOT NULL, -- new, like_new, good, fair, poor
    status VARCHAR(50) NOT NULL DEFAULT 'active', -- active, sold, reserved, deleted
    weight_kg DECIMAL(8, 2),
    manufacturer_country VARCHAR(100),
    estimated_manufacturing_year INTEGER,
    ai_generated_description TEXT,
    ai_suggested_price INTEGER,
    co2_impact_kg DECIMAL(10, 2), -- CO2 saved by buying used
    view_count INTEGER DEFAULT 0,
    favorite_count INTEGER DEFAULT 0,
    has_3d_model BOOLEAN DEFAULT FALSE,
    model_url TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    sold_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_products_seller_id ON products(seller_id);
CREATE INDEX idx_products_status ON products(status);
CREATE INDEX idx_products_category ON products(category);
CREATE INDEX idx_products_price ON products(price);
CREATE INDEX idx_products_created_at ON products(created_at DESC);
CREATE INDEX idx_products_search ON products USING GIN(to_tsvector('english', title || ' ' || COALESCE(description, '')));
```

### product_images
```sql
CREATE TABLE product_images (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    image_url TEXT NOT NULL,
    cdn_url TEXT,
    display_order INTEGER DEFAULT 0,
    is_primary BOOLEAN DEFAULT FALSE,
    width INTEGER,
    height INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_product_images_product_id ON product_images(product_id);
CREATE INDEX idx_product_images_primary ON product_images(product_id, is_primary);
```

### purchases
```sql
CREATE TABLE purchases (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id),
    buyer_id UUID NOT NULL REFERENCES users(id),
    seller_id UUID NOT NULL REFERENCES users(id),
    price INTEGER NOT NULL,
    co2_saved_kg DECIMAL(10, 2) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending', -- pending, completed, cancelled
    payment_method VARCHAR(50), -- mock only
    shipping_address TEXT,
    completed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_purchases_buyer_id ON purchases(buyer_id);
CREATE INDEX idx_purchases_seller_id ON purchases(seller_id);
CREATE INDEX idx_purchases_product_id ON purchases(product_id);
CREATE INDEX idx_purchases_created_at ON purchases(created_at DESC);
```

### conversations
```sql
CREATE TABLE conversations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID REFERENCES products(id),
    last_message_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_conversations_last_message ON conversations(last_message_at DESC);
```

### conversation_participants
```sql
CREATE TABLE conversation_participants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id),
    last_read_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(conversation_id, user_id)
);

CREATE INDEX idx_conversation_participants_user ON conversation_participants(user_id);
CREATE INDEX idx_conversation_participants_conversation ON conversation_participants(conversation_id);
```

### messages
```sql
CREATE TABLE messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    sender_id UUID NOT NULL REFERENCES users(id),
    content TEXT NOT NULL,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_messages_conversation_id ON messages(conversation_id, created_at DESC);
CREATE INDEX idx_messages_sender_id ON messages(sender_id);
```

### favorites
```sql
CREATE TABLE favorites (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, product_id)
);

CREATE INDEX idx_favorites_user_id ON favorites(user_id);
CREATE INDEX idx_favorites_product_id ON favorites(product_id);
```

### achievements
```sql
CREATE TABLE achievements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    badge_icon_url TEXT,
    requirement_type VARCHAR(50) NOT NULL, -- co2_saved, transaction_count, level
    requirement_value INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Seed data
INSERT INTO achievements (name, description, requirement_type, requirement_value) VALUES
('First Step', 'Complete your first transaction', 'transaction_count', 1),
('Eco Warrior', 'Save 10kg of CO2', 'co2_saved', 10),
('Planet Hero', 'Save 50kg of CO2', 'co2_saved', 50),
('Climate Champion', 'Save 100kg of CO2', 'co2_saved', 100),
('Master Trader', 'Complete 50 transactions', 'transaction_count', 50);
```

### user_achievements
```sql
CREATE TABLE user_achievements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    achievement_id UUID NOT NULL REFERENCES achievements(id),
    earned_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, achievement_id)
);

CREATE INDEX idx_user_achievements_user_id ON user_achievements(user_id);
```

### sustainability_logs
```sql
CREATE TABLE sustainability_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    purchase_id UUID REFERENCES purchases(id),
    action_type VARCHAR(50) NOT NULL, -- purchase, sale
    co2_saved_kg DECIMAL(10, 2) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_sustainability_logs_user_id ON sustainability_logs(user_id);
CREATE INDEX idx_sustainability_logs_created_at ON sustainability_logs(created_at DESC);
```

## Key Design Decisions

1. **UUID Primary Keys**: Better for distributed systems and security (no sequential ID guessing)
2. **Soft Deletes**: `deleted_at` for users and products to maintain data integrity
3. **GIN Index for Full-Text Search**: Fast product search on title + description
4. **Denormalization**: `total_co2_saved_kg` in users table for fast leaderboard queries
5. **Conversation Model**: Flexible DM system supporting future group chats
6. **CDN URLs**: Separate column for CDN-optimized image URLs
7. **Decimal for CO2/Money**: Precise calculations without floating-point errors

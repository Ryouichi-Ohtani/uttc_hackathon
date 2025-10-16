-- Demo Users
INSERT INTO users (id, email, password_hash, username, avatar_url, bio, total_co2_saved_kg, created_at, updated_at)
VALUES
  ('550e8400-e29b-41d4-a716-446655440001', 'alice@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMye1J8eqXkPqz9fBfXKXJz1bD0h8kC8.6K', 'alice_green', 'https://i.pravatar.cc/150?img=1', 'エコ活動が大好きです！', 125.5, NOW(), NOW()),
  ('550e8400-e29b-41d4-a716-446655440002', 'bob@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMye1J8eqXkPqz9fBfXKXJz1bD0h8kC8.6K', 'bob_eco', 'https://i.pravatar.cc/150?img=12', 'サステナブルな生活を実践中', 89.3, NOW(), NOW()),
  ('550e8400-e29b-41d4-a716-446655440003', 'carol@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMye1J8eqXkPqz9fBfXKXJz1bD0h8kC8.6K', 'carol_reuse', 'https://i.pravatar.cc/150?img=5', 'リユース・リサイクル推進派', 156.8, NOW(), NOW())
ON CONFLICT (email) DO NOTHING;

-- Demo Products
INSERT INTO products (id, seller_id, title, description, category, price, condition, weight_kg, manufacturer_country, estimated_manufacturing_year, co2_impact_kg, status, created_at, updated_at)
VALUES
  -- Electronics
  ('650e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440001', 'iPhone 12 Pro 128GB', 'バッテリー状態良好、画面に小傷あり。動作確認済み。付属品：充電ケーブル、箱付き。', 'electronics', 45000, 'good', 0.189, 'China', 2020, 15.2, 'active', NOW(), NOW()),
  ('650e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440002', 'MacBook Air M1 2020', '使用頻度少なめ、傷なし美品です。バッテリーサイクル数：50回程度。付属品完備。', 'electronics', 78000, 'excellent', 1.29, 'China', 2020, 42.5, 'active', NOW(), NOW()),
  ('650e8400-e29b-41d4-a716-446655440003', '550e8400-e29b-41d4-a716-446655440003', 'Sony WH-1000XM4 ヘッドホン', 'ノイズキャンセリング機能完璧、付属品完備。使用期間：約1年。', 'electronics', 22000, 'excellent', 0.254, 'Malaysia', 2021, 8.7, 'active', NOW(), NOW()),
  ('650e8400-e29b-41d4-a716-446655440011', '550e8400-e29b-41d4-a716-446655440001', 'iPad Air 第4世代 64GB', 'Wi-Fiモデル、ほぼ未使用。保護フィルム・ケース付き。', 'electronics', 52000, 'excellent', 0.458, 'China', 2021, 18.5, 'active', NOW(), NOW()),

  -- Clothing
  ('650e8400-e29b-41d4-a716-446655440004', '550e8400-e29b-41d4-a716-446655440001', 'パタゴニア フリースジャケット', 'サイズM、ほとんど着用していません。クリーニング済み。', 'clothing', 8900, 'excellent', 0.45, 'Vietnam', 2022, 6.8, 'active', NOW(), NOW()),
  ('650e8400-e29b-41d4-a716-446655440005', '550e8400-e29b-41d4-a716-446655440002', 'ノースフェイス ダウンジャケット', 'サイズL、若干の使用感あり。クリーニング済み、機能性抜群。', 'clothing', 12000, 'good', 0.68, 'China', 2021, 10.2, 'active', NOW(), NOW()),
  ('650e8400-e29b-41d4-a716-446655440012', '550e8400-e29b-41d4-a716-446655440003', 'ユニクロ カシミヤセーター', 'サイズM、ネイビー。2回のみ着用。', 'clothing', 3500, 'excellent', 0.25, 'China', 2023, 3.8, 'active', NOW(), NOW()),

  -- Furniture
  ('650e8400-e29b-41d4-a716-446655440006', '550e8400-e29b-41d4-a716-446655440003', 'IKEAオフィスチェア MARKUS', '座り心地良好、目立つ傷なし。引き取り希望（東京都内）。', 'furniture', 15000, 'good', 15.2, 'Poland', 2020, 35.6, 'active', NOW(), NOW()),
  ('650e8400-e29b-41d4-a716-446655440007', '550e8400-e29b-41d4-a716-446655440001', '無印良品 木製デスク', 'シンプルなデザイン、多少の使用感あり。サイズ：120x60cm。', 'furniture', 18000, 'good', 22.5, 'Japan', 2019, 48.3, 'active', NOW(), NOW()),

  -- Books
  ('650e8400-e29b-41d4-a716-446655440008', '550e8400-e29b-41d4-a716-446655440002', 'プログラミング技術書セット 10冊', 'JavaScript、Python、Go関連の技術書。書き込みなし、美品。', 'books', 5500, 'good', 3.2, 'Japan', 2021, 2.1, 'active', NOW(), NOW()),
  ('650e8400-e29b-41d4-a716-446655440013', '550e8400-e29b-41d4-a716-446655440002', 'ビジネス書セット 15冊', 'マーケティング、経営戦略など。カバー付き。', 'books', 4800, 'good', 4.5, 'Japan', 2022, 2.8, 'active', NOW(), NOW()),

  -- Sports
  ('650e8400-e29b-41d4-a716-446655440009', '550e8400-e29b-41d4-a716-446655440003', 'ヨガマット Manduka', '数回使用のみ、ほぼ新品。厚さ6mm、収納ケース付き。', 'sports', 6800, 'excellent', 1.8, 'USA', 2022, 3.4, 'active', NOW(), NOW()),
  ('650e8400-e29b-41d4-a716-446655440010', '550e8400-e29b-41d4-a716-446655440001', 'ロードバイク GIANT TCR', '定期メンテナンス済み、走行距離3000km。引き取り希望。', 'sports', 95000, 'good', 8.5, 'Taiwan', 2020, 78.5, 'active', NOW(), NOW()),
  ('650e8400-e29b-41d4-a716-446655440014', '550e8400-e29b-41d4-a716-446655440003', 'ダンベルセット 20kg', '2.5kg x 4個、5kg x 2個のプレート。シャフト付き。', 'sports', 8500, 'good', 20.0, 'China', 2021, 12.5, 'active', NOW(), NOW()),

  -- Home & Kitchen
  ('650e8400-e29b-41d4-a716-446655440015', '550e8400-e29b-41d4-a716-446655440002', 'バルミューダ トースター', '使用頻度少、動作良好。説明書付き。', 'home', 15000, 'excellent', 4.3, 'China', 2022, 8.9, 'active', NOW(), NOW()),
  ('650e8400-e29b-41d4-a716-446655440016', '550e8400-e29b-41d4-a716-446655440001', 'ル・クルーゼ鍋 22cm', 'オレンジ色、数回使用。IH対応。', 'home', 12000, 'excellent', 3.8, 'France', 2021, 6.2, 'active', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Product Images
INSERT INTO product_images (product_id, image_url, cdn_url, display_order, is_primary)
VALUES
  ('650e8400-e29b-41d4-a716-446655440001', 'https://images.unsplash.com/photo-1592286927505-2ff62f2b5a3d?w=800', 'https://images.unsplash.com/photo-1592286927505-2ff62f2b5a3d?w=800', 0, true),
  ('650e8400-e29b-41d4-a716-446655440002', 'https://images.unsplash.com/photo-1517336714731-489689fd1ca8?w=800', 'https://images.unsplash.com/photo-1517336714731-489689fd1ca8?w=800', 0, true),
  ('650e8400-e29b-41d4-a716-446655440003', 'https://images.unsplash.com/photo-1546435770-a3e426bf472b?w=800', 'https://images.unsplash.com/photo-1546435770-a3e426bf472b?w=800', 0, true),
  ('650e8400-e29b-41d4-a716-446655440004', 'https://images.unsplash.com/photo-1591047139829-d91aecb6caea?w=800', 'https://images.unsplash.com/photo-1591047139829-d91aecb6caea?w=800', 0, true),
  ('650e8400-e29b-41d4-a716-446655440005', 'https://images.unsplash.com/photo-1551028719-00167b16eac5?w=800', 'https://images.unsplash.com/photo-1551028719-00167b16eac5?w=800', 0, true),
  ('650e8400-e29b-41d4-a716-446655440006', 'https://images.unsplash.com/photo-1580480055273-228ff5388ef8?w=800', 'https://images.unsplash.com/photo-1580480055273-228ff5388ef8?w=800', 0, true),
  ('650e8400-e29b-41d4-a716-446655440007', 'https://images.unsplash.com/photo-1518455027359-f3f8164ba6bd?w=800', 'https://images.unsplash.com/photo-1518455027359-f3f8164ba6bd?w=800', 0, true),
  ('650e8400-e29b-41d4-a716-446655440008', 'https://images.unsplash.com/photo-1495446815901-a7297e633e8d?w=800', 'https://images.unsplash.com/photo-1495446815901-a7297e633e8d?w=800', 0, true),
  ('650e8400-e29b-41d4-a716-446655440009', 'https://images.unsplash.com/photo-1601925260368-ae2f83cf8b7f?w=800', 'https://images.unsplash.com/photo-1601925260368-ae2f83cf8b7f?w=800', 0, true),
  ('650e8400-e29b-41d4-a716-446655440010', 'https://images.unsplash.com/photo-1485965120184-e220f721d03e?w=800', 'https://images.unsplash.com/photo-1485965120184-e220f721d03e?w=800', 0, true),
  ('650e8400-e29b-41d4-a716-446655440011', 'https://images.unsplash.com/photo-1544244015-0df4b3ffc6b0?w=800', 'https://images.unsplash.com/photo-1544244015-0df4b3ffc6b0?w=800', 0, true),
  ('650e8400-e29b-41d4-a716-446655440012', 'https://images.unsplash.com/photo-1620799140408-edc6dcb6d633?w=800', 'https://images.unsplash.com/photo-1620799140408-edc6dcb6d633?w=800', 0, true),
  ('650e8400-e29b-41d4-a716-446655440013', 'https://images.unsplash.com/photo-1524995997946-a1c2e315a42f?w=800', 'https://images.unsplash.com/photo-1524995997946-a1c2e315a42f?w=800', 0, true),
  ('650e8400-e29b-41d4-a716-446655440014', 'https://images.unsplash.com/photo-1517836357463-d25dfeac3438?w=800', 'https://images.unsplash.com/photo-1517836357463-d25dfeac3438?w=800', 0, true),
  ('650e8400-e29b-41d4-a716-446655440015', 'https://images.unsplash.com/photo-1585659722983-3a675dabf23d?w=800', 'https://images.unsplash.com/photo-1585659722983-3a675dabf23d?w=800', 0, true),
  ('650e8400-e29b-41d4-a716-446655440016', 'https://images.unsplash.com/photo-1584990347449-c4a5e3497faa?w=800', 'https://images.unsplash.com/photo-1584990347449-c4a5e3497faa?w=800', 0, true)
ON CONFLICT DO NOTHING;

-- Sustainability Logs (購入ベースではなく、アクション履歴として)
INSERT INTO sustainability_logs (user_id, action_type, co2_saved_kg, description, created_at)
VALUES
  ('550e8400-e29b-41d4-a716-446655440001', 'listing', 15.2, 'iPhone 12 Pro を出品しました', NOW() - INTERVAL '5 days'),
  ('550e8400-e29b-41d4-a716-446655440001', 'listing', 6.8, 'パタゴニア フリースジャケットを出品しました', NOW() - INTERVAL '3 days'),
  ('550e8400-e29b-41d4-a716-446655440002', 'listing', 42.5, 'MacBook Air M1を出品しました', NOW() - INTERVAL '7 days'),
  ('550e8400-e29b-41d4-a716-446655440003', 'listing', 35.6, 'IKEAオフィスチェアを出品しました', NOW() - INTERVAL '2 days'),
  ('550e8400-e29b-41d4-a716-446655440001', 'listing', 18.5, 'iPad Airを出品しました', NOW() - INTERVAL '1 days'),
  ('550e8400-e29b-41d4-a716-446655440002', 'listing', 8.9, 'バルミューダトースターを出品しました', NOW() - INTERVAL '4 days')
ON CONFLICT DO NOTHING;

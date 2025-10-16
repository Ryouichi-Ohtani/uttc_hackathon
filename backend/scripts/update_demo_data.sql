-- 新しい商品を追加
INSERT INTO products (id, seller_id, title, description, category, price, condition, weight_kg, manufacturer_country, estimated_manufacturing_year, co2_impact_kg, status, created_at, updated_at)
VALUES
  -- さらに多様な電子機器
  ('650e8400-e29b-41d4-a716-446655440020', '550e8400-e29b-41d4-a716-446655440002', 'Nintendo Switch 本体', 'Joy-Con付き、動作確認済み。ケース・ソフト2本付き。', 'electronics', 25000, 'good', 0.297, 'China', 2020, 11.5, 'active', NOW(), NOW()),
  ('650e8400-e29b-41d4-a716-446655440021', '550e8400-e29b-41d4-a716-446655440003', 'Canon EOS Kiss X10', 'レンズキット、初心者向けミラーレス一眼。撮影枚数少。', 'electronics', 68000, 'excellent', 0.449, 'Japan', 2021, 28.3, 'active', NOW(), NOW()),
  ('650e8400-e29b-41d4-a716-446655440022', '550e8400-e29b-41d4-a716-446655440001', 'iPad Pro 11インチ 256GB', 'Apple Pencil第2世代付き。Magic Keyboard対応。', 'electronics', 95000, 'excellent', 0.468, 'China', 2022, 35.2, 'active', NOW(), NOW()),
  ('650e8400-e29b-41d4-a716-446655440023', '550e8400-e29b-41d4-a716-446655440002', 'Kindle Paperwhite', '第11世代、広告なし。ほぼ未使用。', 'electronics', 12000, 'excellent', 0.205, 'China', 2023, 4.8, 'active', NOW(), NOW()),

  -- 衣類をさらに追加
  ('650e8400-e29b-41d4-a716-446655440024', '550e8400-e29b-41d4-a716-446655440003', 'ユニクロ ウルトラライトダウン', 'サイズS、レッド。コンパクトに収納可能。', 'clothing', 3800, 'good', 0.32, 'Vietnam', 2023, 4.5, 'active', NOW(), NOW()),
  ('650e8400-e29b-41d4-a716-446655440025', '550e8400-e29b-41d4-a716-446655440001', 'ZARA レザージャケット', 'サイズM、ブラック。本革、美品。', 'clothing', 15000, 'excellent', 1.2, 'Turkey', 2022, 18.0, 'active', NOW(), NOW()),
  ('650e8400-e29b-41d4-a716-446655440026', '550e8400-e29b-41d4-a716-446655440002', 'ナイキ エアマックス95', 'サイズ27cm、ほぼ新品。箱付き。', 'clothing', 12000, 'excellent', 0.65, 'Vietnam', 2023, 9.8, 'active', NOW(), NOW()),

  -- 家具・インテリア
  ('650e8400-e29b-41d4-a716-446655440027', '550e8400-e29b-41d4-a716-446655440003', 'ニトリ 3段カラーボックス 2個セット', 'ホワイト、組み立て済み。引き取り希望。', 'furniture', 3000, 'good', 18.0, 'Japan', 2022, 25.2, 'active', NOW(), NOW()),
  ('650e8400-e29b-41d4-a716-446655440028', '550e8400-e29b-41d4-a716-446655440001', 'KEYUCA デスクライト LED', 'タッチ調光機能付き、USB給電。', 'furniture', 4500, 'excellent', 0.8, 'China', 2022, 3.2, 'active', NOW(), NOW()),

  -- スポーツ・アウトドア用品
  ('650e8400-e29b-41d4-a716-446655440029', '550e8400-e29b-41d4-a716-446655440002', 'コールマン テント 4人用', '使用回数3回のみ。ペグ・収納袋完備。', 'sports', 18000, 'excellent', 6.5, 'China', 2021, 15.8, 'active', NOW(), NOW()),
  ('650e8400-e29b-41d4-a716-446655440030', '550e8400-e29b-41d4-a716-446655440003', 'スノーボード板＋ビンディング', '155cm、BURTON製。ブーツは別売。', 'sports', 28000, 'good', 4.8, 'Austria', 2020, 22.5, 'active', NOW(), NOW()),
  ('650e8400-e29b-41d4-a716-446655440031', '550e8400-e29b-41d4-a716-446655440001', 'フィットネスバイク', '折りたたみ式、静音設計。引き取り希望。', 'sports', 15000, 'good', 18.5, 'China', 2022, 32.5, 'active', NOW(), NOW()),

  -- 家電・生活用品
  ('650e8400-e29b-41d4-a716-446655440032', '550e8400-e29b-41d4-a716-446655440002', 'ダイソン掃除機 V8', 'コードレス、バッテリー良好。付属品完備。', 'home', 28000, 'good', 2.6, 'Malaysia', 2021, 15.6, 'active', NOW(), NOW()),
  ('650e8400-e29b-41d4-a716-446655440033', '550e8400-e29b-41d4-a716-446655440003', 'ネスプレッソ コーヒーメーカー', 'カプセル式、ほぼ未使用。説明書付き。', 'home', 12000, 'excellent', 3.2, 'Switzerland', 2022, 8.9, 'active', NOW(), NOW()),
  ('650e8400-e29b-41d4-a716-446655440034', '550e8400-e29b-41d4-a716-446655440001', '象印 炊飯器 5.5合', 'IH式、内釜美品。2年使用。', 'home', 15000, 'good', 4.5, 'Japan', 2022, 12.5, 'active', NOW(), NOW()),

  -- 楽器・趣味
  ('650e8400-e29b-41d4-a716-446655440035', '550e8400-e29b-41d4-a716-446655440002', 'ヤマハ 電子ピアノ P-125', '88鍵盤、スタンド・ペダル付き。引き取り希望。', 'hobbies', 45000, 'excellent', 11.8, 'Indonesia', 2021, 45.8, 'active', NOW(), NOW()),
  ('650e8400-e29b-41d4-a716-446655440036', '550e8400-e29b-41d4-a716-446655440003', 'レゴ クリエイターシリーズ 5セット', '組み立て説明書付き、箱あり。コレクター放出。', 'hobbies', 18000, 'excellent', 5.2, 'Denmark', 2020, 8.5, 'active', NOW(), NOW()),

  -- ベビー・キッズ用品
  ('650e8400-e29b-41d4-a716-446655440037', '550e8400-e29b-41d4-a716-446655440001', 'コンビ ベビーカー', 'A型、使用期間短め。クリーニング済み。', 'kids', 22000, 'good', 5.8, 'China', 2022, 18.5, 'active', NOW(), NOW()),
  ('650e8400-e29b-41d4-a716-446655440038', '550e8400-e29b-41d4-a716-446655440002', '子供服まとめ売り 100-110cm', '20着セット、ブランド服含む。', 'kids', 8000, 'good', 2.5, 'China', 2022, 6.5, 'active', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- 新しい商品の画像を追加
INSERT INTO product_images (product_id, image_url, cdn_url, display_order, is_primary)
VALUES
  ('650e8400-e29b-41d4-a716-446655440020', 'https://images.unsplash.com/photo-1578303512597-81e6cc155b3e?w=800', 'https://images.unsplash.com/photo-1578303512597-81e6cc155b3e?w=800', 0, true),
  ('650e8400-e29b-41d4-a716-446655440021', 'https://images.unsplash.com/photo-1606983340126-99ab4feaa64a?w=800', 'https://images.unsplash.com/photo-1606983340126-99ab4feaa64a?w=800', 0, true),
  ('650e8400-e29b-41d4-a716-446655440022', 'https://images.unsplash.com/photo-1585790050230-5dd28404f549?w=800', 'https://images.unsplash.com/photo-1585790050230-5dd28404f549?w=800', 0, true),
  ('650e8400-e29b-41d4-a716-446655440023', 'https://images.unsplash.com/photo-1589998059171-988d887df646?w=800', 'https://images.unsplash.com/photo-1589998059171-988d887df646?w=800', 0, true),
  ('650e8400-e29b-41d4-a716-446655440024', 'https://images.unsplash.com/photo-1539533018447-63fcce2678e3?w=800', 'https://images.unsplash.com/photo-1539533018447-63fcce2678e3?w=800', 0, true),
  ('650e8400-e29b-41d4-a716-446655440025', 'https://images.unsplash.com/photo-1551028719-00167b16eac5?w=800', 'https://images.unsplash.com/photo-1551028719-00167b16eac5?w=800', 0, true),
  ('650e8400-e29b-41d4-a716-446655440026', 'https://images.unsplash.com/photo-1542291026-7eec264c27ff?w=800', 'https://images.unsplash.com/photo-1542291026-7eec264c27ff?w=800', 0, true),
  ('650e8400-e29b-41d4-a716-446655440027', 'https://images.unsplash.com/photo-1558618666-fcd25c85cd64?w=800', 'https://images.unsplash.com/photo-1558618666-fcd25c85cd64?w=800', 0, true),
  ('650e8400-e29b-41d4-a716-446655440028', 'https://images.unsplash.com/photo-1507473885765-e6ed057f782c?w=800', 'https://images.unsplash.com/photo-1507473885765-e6ed057f782c?w=800', 0, true),
  ('650e8400-e29b-41d4-a716-446655440029', 'https://images.unsplash.com/photo-1478131143081-80f7f84ca84d?w=800', 'https://images.unsplash.com/photo-1478131143081-80f7f84ca84d?w=800', 0, true),
  ('650e8400-e29b-41d4-a716-446655440030', 'https://images.unsplash.com/photo-1418985991508-e47386d96a71?w=800', 'https://images.unsplash.com/photo-1418985991508-e47386d96a71?w=800', 0, true),
  ('650e8400-e29b-41d4-a716-446655440031', 'https://images.unsplash.com/photo-1517836357463-d25dfeac3438?w=800', 'https://images.unsplash.com/photo-1517836357463-d25dfeac3438?w=800', 0, true),
  ('650e8400-e29b-41d4-a716-446655440032', 'https://images.unsplash.com/photo-1558317374-067fb5f30001?w=800', 'https://images.unsplash.com/photo-1558317374-067fb5f30001?w=800', 0, true),
  ('650e8400-e29b-41d4-a716-446655440033', 'https://images.unsplash.com/photo-1495474472287-4d71bcdd2085?w=800', 'https://images.unsplash.com/photo-1495474472287-4d71bcdd2085?w=800', 0, true),
  ('650e8400-e29b-41d4-a716-446655440034', 'https://images.unsplash.com/photo-1556911220-bff31c812dba?w=800', 'https://images.unsplash.com/photo-1556911220-bff31c812dba?w=800', 0, true),
  ('650e8400-e29b-41d4-a716-446655440035', 'https://images.unsplash.com/photo-1520523839897-bd0b52f945a0?w=800', 'https://images.unsplash.com/photo-1520523839897-bd0b52f945a0?w=800', 0, true),
  ('650e8400-e29b-41d4-a716-446655440036', 'https://images.unsplash.com/photo-1587654780291-39c9404d746b?w=800', 'https://images.unsplash.com/photo-1587654780291-39c9404d746b?w=800', 0, true),
  ('650e8400-e29b-41d4-a716-446655440037', 'https://images.unsplash.com/photo-1522771930-78848d9293e8?w=800', 'https://images.unsplash.com/photo-1522771930-78848d9293e8?w=800', 0, true),
  ('650e8400-e29b-41d4-a716-446655440038', 'https://images.unsplash.com/photo-1519457431-44ccd64a579b?w=800', 'https://images.unsplash.com/photo-1519457431-44ccd64a579b?w=800', 0, true)
ON CONFLICT DO NOTHING;

-- 既存ユーザーのCO2削減量を更新
UPDATE users
SET total_co2_saved_kg = total_co2_saved_kg + 50.0,
    updated_at = NOW()
WHERE email IN ('alice@example.com', 'bob@example.com', 'carol@example.com');

-- サステナビリティログを追加
INSERT INTO sustainability_logs (user_id, action_type, co2_saved_kg, description, created_at)
VALUES
  ('550e8400-e29b-41d4-a716-446655440002', 'listing', 11.5, 'Nintendo Switchを出品しました', NOW() - INTERVAL '1 hours'),
  ('550e8400-e29b-41d4-a716-446655440003', 'listing', 28.3, 'Canon EOS Kiss X10を出品しました', NOW() - INTERVAL '2 hours'),
  ('550e8400-e29b-41d4-a716-446655440001', 'listing', 35.2, 'iPad Proを出品しました', NOW() - INTERVAL '3 hours'),
  ('550e8400-e29b-41d4-a716-446655440002', 'listing', 15.8, 'テントを出品しました', NOW() - INTERVAL '5 hours'),
  ('550e8400-e29b-41d4-a716-446655440003', 'listing', 45.8, 'ヤマハ電子ピアノを出品しました', NOW() - INTERVAL '6 hours')
ON CONFLICT DO NOTHING;

-- 商品の閲覧数・お気に入り数を更新（人気商品を作る）
UPDATE products SET view_count = 125, favorite_count = 23 WHERE id = '650e8400-e29b-41d4-a716-446655440001'; -- iPhone
UPDATE products SET view_count = 98, favorite_count = 18 WHERE id = '650e8400-e29b-41d4-a716-446655440002'; -- MacBook
UPDATE products SET view_count = 156, favorite_count = 31 WHERE id = '650e8400-e29b-41d4-a716-446655440010'; -- ロードバイク
UPDATE products SET view_count = 87, favorite_count = 15 WHERE id = '650e8400-e29b-41d4-a716-446655440022'; -- iPad Pro
UPDATE products SET view_count = 65, favorite_count = 12 WHERE id = '650e8400-e29b-41d4-a716-446655440020'; -- Switch
UPDATE products SET view_count = 112, favorite_count = 25 WHERE id = '650e8400-e29b-41d4-a716-446655440032'; -- ダイソン
UPDATE products SET view_count = 45, favorite_count = 8 WHERE id = '650e8400-e29b-41d4-a716-446655440004'; -- パタゴニア
UPDATE products SET view_count = 78, favorite_count = 14 WHERE id = '650e8400-e29b-41d4-a716-446655440029'; -- テント

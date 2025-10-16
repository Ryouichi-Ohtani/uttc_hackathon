-- リアルな3Dモデル付き商品を追加

-- 1. ビンテージカメラ (電子機器)
INSERT INTO products (
  id, seller_id, title, description, category, price, condition,
  ai_generated_description, ai_suggested_price, model_url,
  created_at, updated_at
) VALUES (
  gen_random_uuid(),
  (SELECT id FROM users WHERE username = 'alice_green' LIMIT 1),
  'ビンテージカメラ - クラシックデザイン',
  'レトロなデザインのビンテージカメラです。コレクション用としても実用としても使えます。動作確認済み。',
  'electronics',
  15000,
  'good',
  true,
  18000,
  'https://raw.githubusercontent.com/KhronosGroup/glTF-Sample-Models/master/2.0/AntiqueCamera/glTF-Binary/AntiqueCamera.glb',
  NOW(),
  NOW()
);

-- 2. ヘルメット (スポーツ用品)
INSERT INTO products (
  id, seller_id, title, description, category, price, condition,
  ai_generated_description, ai_suggested_price, model_url,
  created_at, updated_at
) VALUES (
  gen_random_uuid(),
  (SELECT id FROM users WHERE username = 'bob_eco' LIMIT 1),
  'フライトヘルメット - ダメージ加工',
  'ビンテージ風のフライトヘルメット。ディスプレイ用に最適です。軽量で丈夫な作り。',
  'sports',
  8500,
  'like_new',
  true,
  9500,
  'https://raw.githubusercontent.com/KhronosGroup/glTF-Sample-Models/master/2.0/DamagedHelmet/glTF-Binary/DamagedHelmet.glb',
  NOW(),
  NOW()
);

-- 3. バギー車 (おもちゃ)
INSERT INTO products (
  id, seller_id, title, description, category, price, condition,
  ai_generated_description, ai_suggested_price, model_url,
  created_at, updated_at
) VALUES (
  gen_random_uuid(),
  (SELECT id FROM users WHERE username = 'alice_green' LIMIT 1),
  'レーシングバギー - リモコンカー',
  '高速走行可能なラジコンバギー。オフロード走行も楽しめます。',
  'toys',
  12000,
  'good',
  true,
  14000,
  'https://raw.githubusercontent.com/KhronosGroup/glTF-Sample-Models/master/2.0/Buggy/glTF-Binary/Buggy.glb',
  NOW(),
  NOW()
);

-- 4. ランタン (家具・インテリア)
INSERT INTO products (
  id, seller_id, title, description, category, price, condition,
  ai_generated_description, ai_suggested_price, model_url,
  created_at, updated_at
) VALUES (
  gen_random_uuid(),
  (SELECT id FROM users WHERE username = 'carol_reuse' LIMIT 1),
  'アンティークランタン - 照明器具',
  'レトロなデザインのランタン。インテリアとして部屋のアクセントに最適。',
  'furniture',
  6800,
  'good',
  true,
  7500,
  'https://raw.githubusercontent.com/KhronosGroup/glTF-Sample-Models/master/2.0/Lantern/glTF-Binary/Lantern.glb',
  NOW(),
  NOW()
);

-- 5. ウォーターボトル (スポーツ用品)
INSERT INTO products (
  id, seller_id, title, description, category, price, condition,
  ai_generated_description, ai_suggested_price, model_url,
  created_at, updated_at
) VALUES (
  gen_random_uuid(),
  (SELECT id FROM users WHERE username = 'bob_eco' LIMIT 1),
  'ステンレスウォーターボトル',
  '保温・保冷機能付きのステンレスボトル。アウトドアやスポーツに最適。',
  'sports',
  3500,
  'like_new',
  true,
  4000,
  'https://raw.githubusercontent.com/KhronosGroup/glTF-Sample-Models/master/2.0/WaterBottle/glTF-Binary/WaterBottle.glb',
  NOW(),
  NOW()
);

-- 6. 飛行機模型 (おもちゃ)
INSERT INTO products (
  id, seller_id, title, description, category, price, condition,
  ai_generated_description, ai_suggested_price, model_url,
  created_at, updated_at
) VALUES (
  gen_random_uuid(),
  (SELECT id FROM users WHERE username = 'alice_green' LIMIT 1),
  'セスナ機 - 飛行機模型',
  'リアルなセスナ機の模型。コレクターアイテムとしても人気。',
  'toys',
  18000,
  'like_new',
  true,
  20000,
  'https://raw.githubusercontent.com/KhronosGroup/glTF-Sample-Models/master/2.0/Cesium/glTF-Binary/Cesium.glb',
  NOW(),
  NOW()
);

-- 7. バーチャルシティ模型 (おもちゃ)
INSERT INTO products (
  id, seller_id, title, description, category, price, condition,
  ai_generated_description, ai_suggested_price, model_url,
  created_at, updated_at
) VALUES (
  gen_random_uuid(),
  (SELECT id FROM users WHERE username = 'carol_reuse' LIMIT 1),
  'ミニチュアシティ - 建築模型',
  '詳細に作り込まれた都市の建築模型。インテリアやジオラマに。',
  'toys',
  25000,
  'new',
  true,
  28000,
  'https://raw.githubusercontent.com/KhronosGroup/glTF-Sample-Models/master/2.0/VirtualCity/glTF-Binary/VirtualCity.glb',
  NOW(),
  NOW()
);

-- 8. スニーカー (衣類)
INSERT INTO products (
  id, seller_id, title, description, category, price, condition,
  ai_generated_description, ai_suggested_price, model_url,
  created_at, updated_at
) VALUES (
  gen_random_uuid(),
  (SELECT id FROM users WHERE username = 'bob_eco' LIMIT 1),
  'ハイトップスニーカー - 27cm',
  'クラシックなハイトップスタイルのスニーカー。カジュアルコーデに最適。',
  'clothing',
  7800,
  'good',
  true,
  9000,
  'https://raw.githubusercontent.com/KhronosGroup/glTF-Sample-Models/master/2.0/MaterialsVariantsShoe/glTF-Binary/MaterialsVariantsShoe.glb',
  NOW(),
  NOW()
);

-- 各商品に画像を追加（ダミーURL）
INSERT INTO product_images (product_id, image_url, display_order)
SELECT id, 'https://via.placeholder.com/800x600/4CAF50/FFFFFF?text=3D+Model', 1
FROM products
WHERE model_url LIKE '%glTF-Sample-Models%'
AND NOT EXISTS (SELECT 1 FROM product_images WHERE product_images.product_id = products.id);

-- 追加された商品を確認
SELECT
  title,
  category,
  price,
  condition,
  CASE WHEN model_url IS NOT NULL THEN '✅ 3D' ELSE '❌' END as has_3d
FROM products
WHERE model_url LIKE '%glTF-Sample-Models%'
ORDER BY created_at DESC;

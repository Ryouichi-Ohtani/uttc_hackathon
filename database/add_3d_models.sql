-- 3Dモデルのカラムを追加
ALTER TABLE products ADD COLUMN IF NOT EXISTS model_url TEXT;

-- 無料の3DモデルURLを各カテゴリの商品に追加
-- これらは実際に動作するGLB形式の3Dモデルです

-- 衣類カテゴリに3Dモデル追加
UPDATE products
SET model_url = 'https://raw.githubusercontent.com/KhronosGroup/glTF-Sample-Models/master/2.0/BoxAnimated/glTF-Binary/BoxAnimated.glb'
WHERE id IN (
  SELECT id FROM products
  WHERE category = 'clothing' AND model_url IS NULL
  LIMIT 5
);

-- 家具カテゴリに3Dモデル追加
UPDATE products
SET model_url = 'https://raw.githubusercontent.com/KhronosGroup/glTF-Sample-Models/master/2.0/Chair/glTF-Binary/Chair.glb'
WHERE id IN (
  SELECT id FROM products
  WHERE category = 'furniture' AND model_url IS NULL
  LIMIT 5
);

-- 電子機器カテゴリに3Dモデル追加
UPDATE products
SET model_url = 'https://raw.githubusercontent.com/KhronosGroup/glTF-Sample-Models/master/2.0/Box/glTF-Binary/Box.glb'
WHERE id IN (
  SELECT id FROM products
  WHERE category = 'electronics' AND model_url IS NULL
  LIMIT 5
);

-- スポーツ用品カテゴリに3Dモデル追加
UPDATE products
SET model_url = 'https://raw.githubusercontent.com/KhronosGroup/glTF-Sample-Models/master/2.0/Sphere/glTF-Binary/Sphere.glb'
WHERE id IN (
  SELECT id FROM products
  WHERE category = 'sports' AND model_url IS NULL
  LIMIT 5
);

-- おもちゃカテゴリに3Dモデル追加
UPDATE products
SET model_url = 'https://raw.githubusercontent.com/KhronosGroup/glTF-Sample-Models/master/2.0/Duck/glTF-Binary/Duck.glb'
WHERE id IN (
  SELECT id FROM products
  WHERE category = 'toys' AND model_url IS NULL
  LIMIT 5
);

-- 更新された商品数を確認
SELECT category, COUNT(*) as count_with_3d
FROM products
WHERE model_url IS NOT NULL
GROUP BY category
ORDER BY category;

#!/bin/bash

echo "🚀 EcoMate フルテスト開始"
echo "========================"
echo ""

# サーバーが起動しているか確認
if curl -s http://localhost:3000 > /dev/null 2>&1; then
    echo "✅ サーバーは既に起動しています"
else
    echo "🔄 サーバーを起動中..."
    docker-compose up -d

    echo "⏳ サーバーの起動を待機中..."
    sleep 10

    # ヘルスチェック
    max_attempts=30
    attempt=0
    while [ $attempt -lt $max_attempts ]; do
        if curl -s http://localhost:3000 > /dev/null 2>&1; then
            echo "✅ フロントエンド起動完了"
            break
        fi
        echo "   待機中... ($((attempt + 1))/$max_attempts)"
        sleep 2
        attempt=$((attempt + 1))
    done

    if [ $attempt -eq $max_attempts ]; then
        echo "❌ サーバーの起動に失敗しました"
        exit 1
    fi

    # バックエンドのヘルスチェック
    attempt=0
    while [ $attempt -lt $max_attempts ]; do
        if curl -s http://localhost:8080/v1/products > /dev/null 2>&1; then
            echo "✅ バックエンド起動完了"
            break
        fi
        echo "   バックエンド待機中... ($((attempt + 1))/$max_attempts)"
        sleep 2
        attempt=$((attempt + 1))
    done
fi

echo ""
echo "🌐 ブラウザテスト実行中..."
echo "========================"
echo ""

# ブラウザテスト実行
node test-browser.js

echo ""
echo "📊 テスト結果サマリー"
echo "========================"

# スクリーンショットの数を確認
screenshot_count=$(ls -1 screenshots/*.png 2>/dev/null | wc -l)
echo "📸 スクリーンショット: ${screenshot_count}枚"

if [ -f "screenshots/01_home.png" ]; then
    echo "✅ ホーム画面: OK"
fi

if [ -f "screenshots/02_product_detail.png" ]; then
    echo "✅ 商品詳細: OK"
fi

if [ -f "screenshots/03_3d_view.png" ]; then
    echo "✅ 3Dビュー: OK"
fi

if [ -f "screenshots/04_ar_view.png" ]; then
    echo "✅ ARビュー: OK"
fi

echo ""
echo "🎉 テスト完了！"
echo "📁 スクリーンショットは ./screenshots/ を確認してください"

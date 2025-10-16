const { chromium } = require('playwright');
const fs = require('fs');

(async () => {
  console.log('🚀 ブラウザテスト開始...\n');

  // ブラウザを起動（headless: false で実際のブラウザが開く）
  const browser = await chromium.launch({
    headless: false,
    slowMo: 500, // 動作を見やすくするため0.5秒遅延
  });

  const context = await browser.newContext({
    viewport: { width: 1920, height: 1080 },
  });

  const page = await context.newPage();

  try {
    // スクリーンショット保存先ディレクトリ作成
    const screenshotDir = './screenshots';
    if (!fs.existsSync(screenshotDir)) {
      fs.mkdirSync(screenshotDir);
    }

    console.log('📱 http://localhost:3000 にアクセス中...');
    await page.goto('http://localhost:3000', { waitUntil: 'networkidle' });
    await page.waitForTimeout(2000);

    // 1. ホーム画面のスクリーンショット
    console.log('📸 スクリーンショット1: ホーム画面');
    await page.screenshot({
      path: `${screenshotDir}/01_home.png`,
      fullPage: true
    });

    // 2. 商品数を確認
    const productCards = await page.locator('[class*="ProductCard"], .product-card, article').count();
    console.log(`✅ 商品カード検出: ${productCards}件\n`);

    if (productCards > 0) {
      console.log('🖱️  最初の商品をクリック...');
      // 最初の商品をクリック
      await page.locator('[class*="ProductCard"], .product-card, article').first().click();
      await page.waitForTimeout(2000);

      // 3. 商品詳細ページのスクリーンショット
      console.log('📸 スクリーンショット2: 商品詳細ページ');
      await page.screenshot({
        path: `${screenshotDir}/02_product_detail.png`,
        fullPage: true
      });

      // 商品タイトルを取得
      const title = await page.locator('h1').first().textContent();
      console.log(`📦 商品タイトル: ${title}\n`);

      // 4. 3Dタブがあるか確認してクリック
      const threeDButton = page.locator('button:has-text("3D"), button:has-text("🎮")');
      const threeDCount = await threeDButton.count();

      if (threeDCount > 0) {
        console.log('🖱️  3Dタブをクリック...');
        await threeDButton.first().click();
        await page.waitForTimeout(3000); // 3Dモデル読み込み待機

        // 5. 3Dビューのスクリーンショット
        console.log('📸 スクリーンショット3: 3Dビュー');
        await page.screenshot({
          path: `${screenshotDir}/03_3d_view.png`,
          fullPage: true
        });
        console.log('✅ 3Dモデルが表示されました\n');
      } else {
        console.log('⚠️  3Dタブが見つかりませんでした\n');
      }

      // 6. ARタブがあるか確認
      const arButton = page.locator('button:has-text("AR"), button:has-text("📸")');
      const arCount = await arButton.count();

      if (arCount > 0) {
        console.log('🖱️  ARタブをクリック...');
        await arButton.first().click();
        await page.waitForTimeout(2000);

        // 7. ARビューのスクリーンショット
        console.log('📸 スクリーンショット4: ARビュー');
        await page.screenshot({
          path: `${screenshotDir}/04_ar_view.png`,
          fullPage: true
        });
        console.log('✅ ARビューが表示されました\n');
      }

      // 8. 戻るボタンをクリック
      const backButton = page.locator('button:has-text("Back"), button:has-text("←")');
      const backCount = await backButton.count();

      if (backCount > 0) {
        console.log('🖱️  戻るボタンをクリック...');
        await backButton.first().click();
        await page.waitForTimeout(2000);

        console.log('📸 スクリーンショット5: ホームに戻った');
        await page.screenshot({
          path: `${screenshotDir}/05_back_to_home.png`,
          fullPage: true
        });
      }
    } else {
      console.log('⚠️  商品が見つかりませんでした');
    }

    console.log('\n✅ テスト完了！');
    console.log(`📁 スクリーンショットは ${screenshotDir}/ に保存されました`);

    // 5秒待ってからブラウザを閉じる
    console.log('\n5秒後にブラウザを閉じます...');
    await page.waitForTimeout(5000);

  } catch (error) {
    console.error('❌ エラー:', error.message);
    await page.screenshot({ path: './screenshots/error.png' });
  } finally {
    await browser.close();
    console.log('🏁 ブラウザを閉じました');
  }
})();

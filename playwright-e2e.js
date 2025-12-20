/**
 * End-to-end smoke test for listing, negotiation, purchase, and shipping flows.
 * Requires frontend dev server on http://127.0.0.1:3000 and backend on http://127.0.0.1:8080.
 */
const { chromium, request } = require('playwright')

const FRONTEND_URL = process.env.FRONTEND_URL || 'http://127.0.0.1:3000'
const API_URL = process.env.API_URL || 'http://127.0.0.1:8080/v1/'
const HEADLESS = process.env.HEADLESS !== '0'

const ts = Date.now()
const seller = {
  email: `seller+${ts}@auto.test`,
  username: `seller${ts}`,
  password: 'password123',
  display_name: 'Playwright Seller',
}

const buyer = {
  email: `buyer+${ts}@auto.test`,
  username: `buyer${ts}`,
  password: 'password123',
  display_name: 'Playwright Buyer',
}

async function ensureUser(api, user) {
  const payload = {
    email: user.email,
    username: user.username,
    password: user.password,
    display_name: user.display_name,
  }

  let res = await api.post('auth/register', { data: payload })
  if (!res.ok()) {
    // likely already exists, fallback to login
    res = await api.post('auth/login', {
      data: { email: user.email, password: user.password },
    })
  }
  const data = await res.json()
  return { ...user, token: data.token, id: data.user.id }
}

async function createProduct(api, token) {
  const productPayload = {
    title: `Auto Test Product ${ts}`,
    description: 'Playwright generated product for E2E checks',
    category: 'electronics',
    price: 5000,
    condition: 'good',
    images: [
      'https://images.unsplash.com/photo-1512436991641-6745cdb1723f',
    ],
  }
  const res = await api.post('products', {
    data: productPayload,
    headers: { Authorization: `Bearer ${token}` },
  })
  if (!res.ok()) {
    throw new Error(`Failed to create product: ${res.status()} ${await res.text()}`)
  }
  return await res.json()
}

async function main() {
  console.log('🚀 Playwright E2E starting...')
  const api = await request.newContext({ baseURL: API_URL })

  // 1) Setup users and product via API
  const sellerUser = await ensureUser(api, seller)
  const buyerUser = await ensureUser(api, buyer)
  console.log('👤 Users ready:', { seller: sellerUser.email, buyer: buyerUser.email })

  const product = await createProduct(api, sellerUser.token)
  console.log('📦 Product created:', product.id, product.title)

  // 2) Buyer flow: login, AI price suggestion, submit offer, purchase
  const browser = await chromium.launch({ headless: HEADLESS })
  const page = await browser.newPage({ viewport: { width: 1280, height: 800 } })

  const loginBuyer = async () => {
    await page.goto(`${FRONTEND_URL}/login`, { waitUntil: 'networkidle' })
    await page.fill('input#email', buyerUser.email)
    await page.fill('input#password', buyerUser.password)
    await page.click('button:has-text("ログイン")')
    await page.waitForURL('**/', { timeout: 10000 })
  }

  console.log('🔐 Logging in as buyer...')
  await loginBuyer()

  console.log('🛍️ Navigating to product detail...')
  await page.goto(`${FRONTEND_URL}/products/${product.id}`, { waitUntil: 'networkidle' })
  await page.waitForSelector('text=価格交渉をする', { timeout: 10000 })

  console.log('🤝 Opening offer dialog and fetching AI suggestion...')
  await page.click('button:has-text("価格交渉をする")')
  await page.click('button:has-text("AI価格提案を受ける")')
  await page.waitForSelector('text=AI価格提案', { timeout: 15000 })

  const aiPriceText = await page.textContent('div:has-text("AI価格提案") >> text=¥')
  const offerInputVal = await page.inputValue('input[type="number"]')
  console.log('💡 AI suggested price:', aiPriceText || offerInputVal)

  console.log('📨 Submitting price negotiation offer...')
  await page.click('button:has-text("交渉を送信")')
  await page.waitForSelector('text=価格交渉を送信しました', { timeout: 10000 })

  console.log('🛒 Proceeding to purchase...')
  await page.click('button:has-text("購入する")')
  await page.waitForSelector('text=Complete Purchase', { timeout: 10000 })
  await page.fill('textarea', 'Test Address, Tokyo, JP')
  await page.selectOption('select', 'credit_card')
  await page.click('button:has-text("Complete Purchase")')
  await page.waitForSelector('text=Purchase completed successfully', { timeout: 15000 })
  console.log('✅ Purchase completed (buyer view)')

  // 3) Seller flow: login, open purchases, generate & approve shipping
  const sellerPage = await browser.newPage({ viewport: { width: 1280, height: 800 } })
  const loginSeller = async () => {
    await sellerPage.goto(`${FRONTEND_URL}/login`, { waitUntil: 'networkidle' })
    await sellerPage.fill('input#email', sellerUser.email)
    await sellerPage.fill('input#password', sellerUser.password)
    await sellerPage.click('button:has-text("ログイン")')
    await sellerPage.waitForURL('**/', { timeout: 10000 })
  }

  console.log('🔐 Logging in as seller...')
  await loginSeller()

  console.log('🚚 Opening purchases as seller to run AI shipping...')
  await sellerPage.goto(`${FRONTEND_URL}/purchases`, { waitUntil: 'networkidle' })
  await sellerPage.click('button:has-text("As Seller")')

  const purchaseCard = sellerPage.locator(`div:has(h3:has-text("${product.title}"))`).first()
  await purchaseCard.waitFor({ state: 'visible', timeout: 15000 })

  const aiShipButton = purchaseCard.getByRole('button', { name: /AIに配送情報を準備させる/ })
  await aiShipButton.click()
  await sellerPage.waitForSelector('text=AI配送準備完了', { timeout: 15000 })
  console.log('📦 AI shipping proposal generated')

  const approveButton = sellerPage.getByRole('button', { name: /この内容で承認/ })
  await approveButton.click()
  await sellerPage.waitForSelector('text=配送情報を承認しました', { timeout: 15000 })
  console.log('✅ Shipping approved (seller view)')

  await browser.close()
  console.log('🏁 E2E completed successfully')
}

main().catch((err) => {
  console.error('❌ E2E failed:', err)
  process.exit(1)
})

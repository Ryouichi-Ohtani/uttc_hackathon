export type Language = 'ja' | 'en'

export interface Translations {
  common: {
    welcome: string
    login: string
    logout: string
    register: string
    search: string
    save: string
    cancel: string
    delete: string
    edit: string
    update: string
    create: string
    back: string
    next: string
    loading: string
    error: string
    success: string
    confirm: string
    yes: string
    no: string
  }
  nav: {
    home: string
    create: string
    aiCreate: string
    purchases: string
    messages: string
    favorites: string
    profile: string
    admin: string
    leaderboard: string
    notifications: string
  }
  auth: {
    loginTitle: string
    registerTitle: string
    email: string
    password: string
    username: string
    displayName: string
    confirmPassword: string
    forgotPassword: string
    noAccount: string
    haveAccount: string
    signUp: string
    signIn: string
  }
  product: {
    title: string
    description: string
    price: string
    category: string
    condition: string
    status: string
    seller: string
    buyer: string
    images: string
    listProduct: string
    editProduct: string
    deleteProduct: string
    viewDetails: string
    addToFavorites: string
    removeFromFavorites: string
    categories: {
      electronics: string
      fashion: string
      home: string
      books: string
      sports: string
      toys: string
      other: string
    }
    conditions: {
      new: string
      likeNew: string
      good: string
      fair: string
      poor: string
    }
    statuses: {
      draft: string
      active: string
      sold: string
      reserved: string
      deleted: string
    }
  }
  purchase: {
    buyNow: string
    purchaseHistory: string
    orderDetails: string
    shippingAddress: string
    paymentMethod: string
    totalAmount: string
    orderStatus: string
    trackShipment: string
    statuses: {
      pending: string
      completed: string
      cancelled: string
    }
  }
  message: {
    sendMessage: string
    typeMessage: string
    conversations: string
    noMessages: string
    startConversation: string
  }
  ai: {
    aiListing: string
    aiNegotiation: string
    aiShipping: string
    generating: string
    analyzing: string
    suggestion: string
    autoMode: string
    assistMode: string
    manualMode: string
  }
  admin: {
    dashboard: string
    users: string
    products: string
    purchases: string
    manage: string
    edit: string
    delete: string
    role: string
    roles: {
      admin: string
      moderator: string
      user: string
    }
  }
  sustainability: {
    leaderboard: string
    co2Saved: string
    ranking: string
    myImpact: string
  }
}

const ja: Translations = {
  common: {
    welcome: 'ようこそ',
    login: 'ログイン',
    logout: 'ログアウト',
    register: '新規登録',
    search: '検索',
    save: '保存',
    cancel: 'キャンセル',
    delete: '削除',
    edit: '編集',
    update: '更新',
    create: '作成',
    back: '戻る',
    next: '次へ',
    loading: '読み込み中...',
    error: 'エラーが発生しました',
    success: '成功しました',
    confirm: '確認',
    yes: 'はい',
    no: 'いいえ',
  },
  nav: {
    home: 'ホーム',
    create: '出品',
    aiCreate: 'AI出品',
    purchases: '購入履歴',
    messages: 'メッセージ',
    favorites: 'お気に入り',
    profile: 'マイページ',
    admin: '管理者',
    leaderboard: 'ランキング',
    notifications: '通知',
  },
  auth: {
    loginTitle: 'ログイン',
    registerTitle: '新規登録',
    email: 'メールアドレス',
    password: 'パスワード',
    username: 'ユーザー名',
    displayName: '表示名',
    confirmPassword: 'パスワード（確認）',
    forgotPassword: 'パスワードを忘れた方',
    noAccount: 'アカウントをお持ちでない方',
    haveAccount: 'すでにアカウントをお持ちの方',
    signUp: '新規登録',
    signIn: 'ログイン',
  },
  product: {
    title: '商品名',
    description: '商品説明',
    price: '価格',
    category: 'カテゴリー',
    condition: '状態',
    status: 'ステータス',
    seller: '出品者',
    buyer: '購入者',
    images: '画像',
    listProduct: '商品を出品',
    editProduct: '商品を編集',
    deleteProduct: '商品を削除',
    viewDetails: '詳細を見る',
    addToFavorites: 'お気に入りに追加',
    removeFromFavorites: 'お気に入りから削除',
    categories: {
      electronics: '電子機器',
      fashion: 'ファッション',
      home: '家庭用品',
      books: '本・雑誌',
      sports: 'スポーツ',
      toys: 'おもちゃ',
      other: 'その他',
    },
    conditions: {
      new: '新品',
      likeNew: '未使用に近い',
      good: '良い',
      fair: '普通',
      poor: '悪い',
    },
    statuses: {
      draft: '下書き',
      active: '販売中',
      sold: '売却済み',
      reserved: '予約済み',
      deleted: '削除済み',
    },
  },
  purchase: {
    buyNow: '今すぐ購入',
    purchaseHistory: '購入履歴',
    orderDetails: '注文詳細',
    shippingAddress: '配送先住所',
    paymentMethod: '支払い方法',
    totalAmount: '合計金額',
    orderStatus: '注文ステータス',
    trackShipment: '配送状況を確認',
    statuses: {
      pending: '処理中',
      completed: '完了',
      cancelled: 'キャンセル',
    },
  },
  message: {
    sendMessage: 'メッセージを送信',
    typeMessage: 'メッセージを入力',
    conversations: '会話一覧',
    noMessages: 'メッセージはありません',
    startConversation: '会話を開始',
  },
  ai: {
    aiListing: 'AI出品',
    aiNegotiation: 'AI交渉',
    aiShipping: 'AI配送',
    generating: '生成中...',
    analyzing: '分析中...',
    suggestion: '提案',
    autoMode: '自動モード',
    assistMode: 'アシストモード',
    manualMode: '手動モード',
  },
  admin: {
    dashboard: '管理者ダッシュボード',
    users: 'ユーザー',
    products: '商品',
    purchases: '購入履歴',
    manage: '管理',
    edit: '編集',
    delete: '削除',
    role: 'ロール',
    roles: {
      admin: '管理者',
      moderator: 'モデレーター',
      user: 'ユーザー',
    },
  },
  sustainability: {
    leaderboard: 'サステナビリティランキング',
    co2Saved: 'CO2削減量',
    ranking: 'ランキング',
    myImpact: '私の貢献',
  },
}

const en: Translations = {
  common: {
    welcome: 'Welcome',
    login: 'Login',
    logout: 'Logout',
    register: 'Sign Up',
    search: 'Search',
    save: 'Save',
    cancel: 'Cancel',
    delete: 'Delete',
    edit: 'Edit',
    update: 'Update',
    create: 'Create',
    back: 'Back',
    next: 'Next',
    loading: 'Loading...',
    error: 'An error occurred',
    success: 'Success',
    confirm: 'Confirm',
    yes: 'Yes',
    no: 'No',
  },
  nav: {
    home: 'Home',
    create: 'List Item',
    aiCreate: 'AI Listing',
    purchases: 'Purchases',
    messages: 'Messages',
    favorites: 'Favorites',
    profile: 'Profile',
    admin: 'Admin',
    leaderboard: 'Leaderboard',
    notifications: 'Notifications',
  },
  auth: {
    loginTitle: 'Login',
    registerTitle: 'Sign Up',
    email: 'Email',
    password: 'Password',
    username: 'Username',
    displayName: 'Display Name',
    confirmPassword: 'Confirm Password',
    forgotPassword: 'Forgot Password?',
    noAccount: "Don't have an account?",
    haveAccount: 'Already have an account?',
    signUp: 'Sign Up',
    signIn: 'Sign In',
  },
  product: {
    title: 'Product Name',
    description: 'Description',
    price: 'Price',
    category: 'Category',
    condition: 'Condition',
    status: 'Status',
    seller: 'Seller',
    buyer: 'Buyer',
    images: 'Images',
    listProduct: 'List Product',
    editProduct: 'Edit Product',
    deleteProduct: 'Delete Product',
    viewDetails: 'View Details',
    addToFavorites: 'Add to Favorites',
    removeFromFavorites: 'Remove from Favorites',
    categories: {
      electronics: 'Electronics',
      fashion: 'Fashion',
      home: 'Home & Garden',
      books: 'Books & Magazines',
      sports: 'Sports',
      toys: 'Toys',
      other: 'Other',
    },
    conditions: {
      new: 'New',
      likeNew: 'Like New',
      good: 'Good',
      fair: 'Fair',
      poor: 'Poor',
    },
    statuses: {
      draft: 'Draft',
      active: 'Active',
      sold: 'Sold',
      reserved: 'Reserved',
      deleted: 'Deleted',
    },
  },
  purchase: {
    buyNow: 'Buy Now',
    purchaseHistory: 'Purchase History',
    orderDetails: 'Order Details',
    shippingAddress: 'Shipping Address',
    paymentMethod: 'Payment Method',
    totalAmount: 'Total Amount',
    orderStatus: 'Order Status',
    trackShipment: 'Track Shipment',
    statuses: {
      pending: 'Pending',
      completed: 'Completed',
      cancelled: 'Cancelled',
    },
  },
  message: {
    sendMessage: 'Send Message',
    typeMessage: 'Type a message',
    conversations: 'Conversations',
    noMessages: 'No messages',
    startConversation: 'Start Conversation',
  },
  ai: {
    aiListing: 'AI Listing',
    aiNegotiation: 'AI Negotiation',
    aiShipping: 'AI Shipping',
    generating: 'Generating...',
    analyzing: 'Analyzing...',
    suggestion: 'Suggestion',
    autoMode: 'Auto Mode',
    assistMode: 'Assist Mode',
    manualMode: 'Manual Mode',
  },
  admin: {
    dashboard: 'Admin Dashboard',
    users: 'Users',
    products: 'Products',
    purchases: 'Purchases',
    manage: 'Manage',
    edit: 'Edit',
    delete: 'Delete',
    role: 'Role',
    roles: {
      admin: 'Administrator',
      moderator: 'Moderator',
      user: 'User',
    },
  },
  sustainability: {
    leaderboard: 'Sustainability Leaderboard',
    co2Saved: 'CO2 Saved',
    ranking: 'Ranking',
    myImpact: 'My Impact',
  },
}

export const translations: Record<Language, Translations> = {
  ja,
  en,
}

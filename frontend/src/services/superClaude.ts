import axios from 'axios'

export interface SuperClaudeMessage {
  id: string
  role: 'user' | 'assistant' | 'system'
  content: string
  timestamp: Date
  attachments?: Attachment[]
  codeBlocks?: CodeBlock[]
  metadata?: MessageMetadata
}

export interface Attachment {
  id: string
  name: string
  type: 'image' | 'document' | 'code' | 'data'
  url?: string
  content?: string
  mimeType: string
  size: number
}

export interface CodeBlock {
  id: string
  language: string
  code: string
  filename?: string
  executable?: boolean
  output?: string
}

export interface MessageMetadata {
  model?: 'claude-3-opus' | 'claude-3-sonnet' | 'claude-3-haiku'
  temperature?: number
  maxTokens?: number
  processingTime?: number
  tokenCount?: {
    input: number
    output: number
  }
}

export interface SuperClaudeCapability {
  id: string
  name: string
  description: string
  icon: string
  category: 'analysis' | 'creation' | 'coding' | 'research' | 'translation' | 'education'
  examples: string[]
}

export interface ConversationContext {
  id: string
  title: string
  messages: SuperClaudeMessage[]
  createdAt: Date
  updatedAt: Date
  settings: ConversationSettings
  sharedWith?: string[]
}

export interface ConversationSettings {
  model: 'claude-3-opus' | 'claude-3-sonnet' | 'claude-3-haiku'
  temperature: number
  maxTokens: number
  systemPrompt?: string
  capabilities: string[]
}

class SuperClaudeService {
  private baseUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080/api'
  private conversations: Map<string, ConversationContext> = new Map()
  private currentConversationId: string | null = null

  // Available capabilities
  readonly capabilities: SuperClaudeCapability[] = [
    {
      id: 'code-generation',
      name: 'コード生成',
      description: 'あらゆるプログラミング言語でコードを生成',
      icon: '💻',
      category: 'coding',
      examples: [
        'Reactコンポーネントを作成して',
        'このアルゴリズムをPythonで実装して',
        'SQLクエリを最適化して'
      ]
    },
    {
      id: 'data-analysis',
      name: 'データ分析',
      description: 'データの分析と可視化の提案',
      icon: '📊',
      category: 'analysis',
      examples: [
        'このCSVファイルを分析して',
        '売上データから傾向を見つけて',
        'A/Bテストの結果を解釈して'
      ]
    },
    {
      id: 'creative-writing',
      name: 'クリエイティブライティング',
      description: 'ストーリー、詩、マーケティングコピーの作成',
      icon: '✍️',
      category: 'creation',
      examples: [
        '商品説明を魅力的に書き直して',
        'ブログ記事のアウトラインを作成',
        'キャッチコピーを10個提案して'
      ]
    },
    {
      id: 'research',
      name: 'リサーチ・調査',
      description: '深い調査と包括的な情報収集',
      icon: '🔍',
      category: 'research',
      examples: [
        '最新のAI技術トレンドをまとめて',
        '競合他社の分析をして',
        'この技術の長所短所を比較して'
      ]
    },
    {
      id: 'translation',
      name: '高度な翻訳',
      description: '文脈を理解した自然な翻訳',
      icon: '🌐',
      category: 'translation',
      examples: [
        '技術文書を日本語に翻訳',
        'このメールを丁寧な英語に',
        'マニュアルを多言語化して'
      ]
    },
    {
      id: 'tutoring',
      name: '学習サポート',
      description: 'あらゆる分野の学習を支援',
      icon: '🎓',
      category: 'education',
      examples: [
        'この概念をわかりやすく説明',
        '練習問題を作成して',
        'ステップバイステップで教えて'
      ]
    },
    {
      id: 'debugging',
      name: 'デバッグ支援',
      description: 'コードのエラー発見と修正提案',
      icon: '🐛',
      category: 'coding',
      examples: [
        'このエラーの原因を特定',
        'パフォーマンスを改善して',
        'セキュリティの脆弱性をチェック'
      ]
    },
    {
      id: 'design-system',
      name: 'デザインシステム',
      description: 'UIデザインとUX改善の提案',
      icon: '🎨',
      category: 'creation',
      examples: [
        'カラーパレットを提案',
        'UIコンポーネントを設計',
        'アクセシビリティを改善'
      ]
    }
  ]

  // Initialize a new conversation
  async createConversation(title?: string): Promise<ConversationContext> {
    const conversationId = this.generateId()
    const context: ConversationContext = {
      id: conversationId,
      title: title || `新しい会話 ${new Date().toLocaleString('ja-JP')}`,
      messages: [],
      createdAt: new Date(),
      updatedAt: new Date(),
      settings: {
        model: 'claude-3-opus',
        temperature: 0.7,
        maxTokens: 4096,
        capabilities: this.capabilities.map(c => c.id)
      }
    }

    this.conversations.set(conversationId, context)
    this.currentConversationId = conversationId
    this.saveToLocalStorage()

    return context
  }

  // Send a message to Super Claude
  async sendMessage(
    content: string,
    attachments?: Attachment[],
    settings?: Partial<ConversationSettings>
  ): Promise<SuperClaudeMessage> {
    if (!this.currentConversationId) {
      await this.createConversation()
    }

    const conversation = this.conversations.get(this.currentConversationId!)
    if (!conversation) throw new Error('Conversation not found')

    // Create user message
    const userMessage: SuperClaudeMessage = {
      id: this.generateId(),
      role: 'user',
      content,
      timestamp: new Date(),
      attachments
    }

    conversation.messages.push(userMessage)

    try {
      // Call the enhanced API endpoint
      const response = await axios.post(`${this.baseUrl}/v1/super-claude/chat`, {
        conversation_id: conversation.id,
        message: content,
        attachments: attachments?.map(a => ({
          type: a.type,
          content: a.content,
          name: a.name,
          mime_type: a.mimeType
        })),
        settings: {
          ...conversation.settings,
          ...settings
        },
        context: conversation.messages.slice(-10) // Send last 10 messages for context
      }, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('auth_token')}`,
          'Content-Type': 'application/json'
        }
      })

      // Create assistant message
      const assistantMessage: SuperClaudeMessage = {
        id: this.generateId(),
        role: 'assistant',
        content: response.data.response,
        timestamp: new Date(),
        codeBlocks: this.extractCodeBlocks(response.data.response),
        metadata: {
          model: conversation.settings.model,
          temperature: conversation.settings.temperature,
          processingTime: response.data.processing_time,
          tokenCount: response.data.token_count
        }
      }

      conversation.messages.push(assistantMessage)
      conversation.updatedAt = new Date()
      this.saveToLocalStorage()

      return assistantMessage
    } catch (error: any) {
      console.error('Super Claude error:', error)

      // Create a mock response for demo purposes
      const mockResponse = this.generateMockResponse(content, conversation.settings)
      const assistantMessage: SuperClaudeMessage = {
        id: this.generateId(),
        role: 'assistant',
        content: mockResponse,
        timestamp: new Date(),
        codeBlocks: this.extractCodeBlocks(mockResponse),
        metadata: {
          model: conversation.settings.model,
          temperature: conversation.settings.temperature,
          processingTime: Math.random() * 2000 + 500,
          tokenCount: {
            input: content.length,
            output: mockResponse.length
          }
        }
      }

      conversation.messages.push(assistantMessage)
      conversation.updatedAt = new Date()
      this.saveToLocalStorage()

      return assistantMessage
    }
  }

  // Analyze uploaded file
  async analyzeFile(file: File): Promise<Attachment> {
    const attachment: Attachment = {
      id: this.generateId(),
      name: file.name,
      type: this.getFileType(file),
      mimeType: file.type,
      size: file.size,
      content: await this.readFileContent(file)
    }

    return attachment
  }

  // Execute code in a sandbox (mock)
  async executeCode(codeBlock: CodeBlock): Promise<string> {
    // This would normally call a secure sandbox API
    // For demo, we'll return a mock output
    return `// Execution result for ${codeBlock.language}:\n// Code executed successfully\n// [Mock output]`
  }

  // Export conversation
  exportConversation(format: 'markdown' | 'json' | 'pdf' = 'markdown'): string {
    if (!this.currentConversationId) return ''

    const conversation = this.conversations.get(this.currentConversationId)
    if (!conversation) return ''

    switch (format) {
      case 'markdown':
        return this.exportAsMarkdown(conversation)
      case 'json':
        return JSON.stringify(conversation, null, 2)
      default:
        return ''
    }
  }

  // Share conversation
  async shareConversation(emails: string[]): Promise<string> {
    if (!this.currentConversationId) throw new Error('No active conversation')

    const conversation = this.conversations.get(this.currentConversationId)
    if (!conversation) throw new Error('Conversation not found')

    // Generate share link (mock)
    const shareId = this.generateId()
    conversation.sharedWith = emails

    return `https://ecomate.ai/shared/${shareId}`
  }

  // Get conversation history
  getConversationHistory(): ConversationContext[] {
    return Array.from(this.conversations.values())
      .sort((a, b) => b.updatedAt.getTime() - a.updatedAt.getTime())
  }

  // Switch conversation
  switchConversation(conversationId: string): ConversationContext | null {
    const conversation = this.conversations.get(conversationId)
    if (conversation) {
      this.currentConversationId = conversationId
      return conversation
    }
    return null
  }

  // Delete conversation
  deleteConversation(conversationId: string): void {
    this.conversations.delete(conversationId)
    if (this.currentConversationId === conversationId) {
      this.currentConversationId = null
    }
    this.saveToLocalStorage()
  }

  // Get current conversation
  getCurrentConversation(): ConversationContext | null {
    if (!this.currentConversationId) return null
    return this.conversations.get(this.currentConversationId) || null
  }

  // Private helper methods
  private generateId(): string {
    return `${Date.now()}-${Math.random().toString(36).substr(2, 9)}`
  }

  private extractCodeBlocks(content: string): CodeBlock[] {
    const codeBlocks: CodeBlock[] = []
    const regex = /```(\w+)?\n([\s\S]*?)```/g
    let match

    while ((match = regex.exec(content)) !== null) {
      codeBlocks.push({
        id: this.generateId(),
        language: match[1] || 'plaintext',
        code: match[2].trim(),
        executable: ['python', 'javascript', 'typescript', 'sql'].includes(match[1]?.toLowerCase() || '')
      })
    }

    return codeBlocks
  }

  private getFileType(file: File): 'image' | 'document' | 'code' | 'data' {
    if (file.type.startsWith('image/')) return 'image'
    if (file.type.includes('pdf') || file.type.includes('document')) return 'document'
    if (file.name.match(/\.(js|ts|py|java|cpp|c|go|rs|swift|kt|rb|php|sh|sql)$/)) return 'code'
    if (file.name.match(/\.(json|csv|xml|yaml|yml)$/)) return 'data'
    return 'document'
  }

  private async readFileContent(file: File): Promise<string> {
    return new Promise((resolve, reject) => {
      const reader = new FileReader()
      reader.onload = (e) => resolve(e.target?.result as string)
      reader.onerror = reject

      if (file.type.startsWith('image/')) {
        reader.readAsDataURL(file)
      } else {
        reader.readAsText(file)
      }
    })
  }

  private exportAsMarkdown(conversation: ConversationContext): string {
    let markdown = `# ${conversation.title}\n\n`
    markdown += `Created: ${conversation.createdAt.toLocaleString()}\n\n`
    markdown += `---\n\n`

    conversation.messages.forEach(msg => {
      markdown += `## ${msg.role === 'user' ? '👤 User' : '🤖 Assistant'}\n\n`
      markdown += `${msg.content}\n\n`

      if (msg.codeBlocks?.length) {
        msg.codeBlocks.forEach(block => {
          markdown += `\`\`\`${block.language}\n${block.code}\n\`\`\`\n\n`
        })
      }
    })

    return markdown
  }

  private saveToLocalStorage(): void {
    const data = {
      conversations: Array.from(this.conversations.entries()),
      currentConversationId: this.currentConversationId
    }
    localStorage.setItem('super_claude_data', JSON.stringify(data))
  }

  private loadFromLocalStorage(): void {
    const stored = localStorage.getItem('super_claude_data')
    if (stored) {
      try {
        const data = JSON.parse(stored)
        this.conversations = new Map(data.conversations)
        this.currentConversationId = data.currentConversationId
      } catch (error) {
        console.error('Failed to load conversation history:', error)
      }
    }
  }

  // Generate mock response for demo
  private generateMockResponse(input: string, settings: ConversationSettings): string {
    const lowerInput = input.toLowerCase()

    // Code generation
    if (lowerInput.includes('コード') || lowerInput.includes('code') || lowerInput.includes('実装')) {
      return `了解しました！リクエストされたコードを生成します。

\`\`\`typescript
// AI-Powered Price Negotiation System
export class PriceNegotiationAI {
  private model: string = '${settings.model}'

  async negotiate(initialPrice: number, targetPrice: number): Promise<number> {
    // Advanced negotiation algorithm
    const factors = this.analyzeMarketConditions()
    const optimalPrice = this.calculateOptimalPrice(initialPrice, targetPrice, factors)

    return optimalPrice
  }

  private analyzeMarketConditions(): MarketFactors {
    // Real-time market analysis
    return {
      demand: 0.8,
      supply: 0.6,
      seasonality: 0.9,
      competition: 0.7
    }
  }

  private calculateOptimalPrice(
    initial: number,
    target: number,
    factors: MarketFactors
  ): number {
    const weight = Object.values(factors).reduce((a, b) => a + b) / 4
    return Math.round(target + (initial - target) * weight)
  }
}
\`\`\`

このコードは以下の特徴があります：

1. **マーケット分析**: リアルタイムの市場状況を考慮
2. **最適価格計算**: 複数の要因を重み付けして計算
3. **柔軟な設定**: モデルやパラメータを調整可能

実装についてさらに詳しく説明が必要でしょうか？`
    }

    // Data analysis
    if (lowerInput.includes('分析') || lowerInput.includes('analyze') || lowerInput.includes('データ')) {
      return `データ分析を実行しました。以下が主要な洞察です：

## 📊 分析結果サマリー

### 1. トレンド分析
- **成長率**: 過去3ヶ月で23.5%の成長
- **ピーク時間**: 午後2時〜4時が最もアクティブ
- **人気カテゴリ**: Electronics (35%), Fashion (28%), Books (22%)

### 2. ユーザー行動パターン
\`\`\`python
import pandas as pd
import matplotlib.pyplot as plt

# ユーザーセグメント分析
segments = {
    'Active Buyers': 42,
    'Window Shoppers': 31,
    'First Time Users': 18,
    'Returning Customers': 9
}

# 視覚化コード
fig, ax = plt.subplots()
ax.pie(segments.values(), labels=segments.keys(), autopct='%1.1f%%')
ax.set_title('User Segmentation Analysis')
\`\`\`

### 3. 推奨アクション
1. ピーク時間にプロモーションを実施
2. Electronicsカテゴリの在庫を増やす
3. 初回ユーザー向けのオンボーディング改善

詳細な分析レポートが必要ですか？`
    }

    // Creative writing
    if (lowerInput.includes('説明') || lowerInput.includes('write') || lowerInput.includes('作成')) {
      return `素晴らしい商品説明を作成しました：

## ✨ AI最適化済み商品説明

**キャッチコピー**: 「未来を今、手の中に」

この革新的な商品は、最新のAI技術と伝統的な職人技を融合させた、まさに新時代のマスターピースです。

### 主な特徴：
- 🌟 **スマートAI機能**: 使用パターンを学習し、あなたに最適化
- 🌍 **エコフレンドリー**: 100%リサイクル可能な素材を使用
- ⚡ **高速パフォーマンス**: 従来モデルの3倍の処理速度
- 🛡️ **安心保証**: 2年間の完全保証付き

### なぜ選ぶべきか：
毎日の生活をより豊かに、より効率的に。このアイテムは単なる製品ではなく、あなたのライフスタイルを進化させるパートナーです。

**限定特典**: 今なら早期購入者限定で20%OFF！

他のバリエーションも作成しましょうか？`
    }

    // Default response
    return `ご質問ありがとうございます！Super Claudeの高度なAI機能を使用して回答します。

## 🤖 AI分析結果

入力内容を分析した結果、以下の情報を提供いたします：

### 理解した内容
"${input}" についてのご質問ですね。

### 提案
1. **詳細分析**: より具体的な情報があれば、さらに精密な分析が可能です
2. **実装例**: 実際のコード例や実装方法をお見せできます
3. **最適化**: 現在のアプローチを改善する方法を提案できます

### 利用可能な機能
- 🔍 深層分析
- 💻 コード生成
- 📊 データ可視化
- ✍️ コンテンツ作成
- 🌐 多言語対応
- 🎓 学習支援

どの機能を使用してさらにサポートいたしましょうか？

**Model**: ${settings.model} | **Temperature**: ${settings.temperature}`
  }

  // Initialize on creation
  constructor() {
    this.loadFromLocalStorage()
  }
}

export const superClaudeService = new SuperClaudeService()
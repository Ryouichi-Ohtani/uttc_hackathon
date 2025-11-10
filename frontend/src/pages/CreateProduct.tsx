import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Button } from '@/components/common/Button'
import { Card } from '@/components/common/Card'
import { Header } from '@/components/layout/Header'
import { ImageUpload } from '@/components/common/ImageUpload'
import api from '@/services/api'
import toast from 'react-hot-toast'

interface AIGeneratedData {
  product_id: string
  listing_data: {
    generated_title: string
    generated_description: string
    generated_category: string
    generated_condition: string
    generated_price: number
    ai_confidence_score: number
  }
  suggested_product: {
    title: string
    description: string
    category: string
    condition: string
    price: number
    weight_kg: number
    detected_brand: string
    detected_model: string
    key_features: string[]
    pricing_rationale: string
  }
  confidence_breakdown: {
    title: number
    description: number
    category: number
    price: number
  }
}

type Step = 'upload' | 'generating' | 'approval'

export const CreateProduct = () => {
  const navigate = useNavigate()
  const [step, setStep] = useState<Step>('upload')
  const [loading, setLoading] = useState(false)
  const [imageUrls, setImageUrls] = useState<string[]>([])
  const [userHints, setUserHints] = useState('')
  const [aiData, setAiData] = useState<AIGeneratedData | null>(null)

  // Editable fields for approval step
  const [editableData, setEditableData] = useState({
    title: '',
    description: '',
    price: 0,
    category: '',
    condition: '',
  })

  const handleGenerateWithAI = async () => {
    if (imageUrls.length === 0) {
      toast.error('少なくとも1枚の画像をアップロードしてください')
      return
    }

    setLoading(true)
    setStep('generating')

    try {
      // Call AI Agent Listing Generation API
      const response = await api.post('/ai-agent/listing/generate', {
        image_urls: imageUrls,
        user_hints: userHints,
        auto_publish: false, // Don't auto-publish, require user approval
      })

      const data: AIGeneratedData = response.data
      setAiData(data)

      // Set editable data from AI suggestions
      setEditableData({
        title: data.suggested_product.title,
        description: data.suggested_product.description,
        price: data.suggested_product.price,
        category: data.suggested_product.category,
        condition: data.suggested_product.condition,
      })

      setStep('approval')
      toast.success('🤖 AIが商品情報を生成しました！')
    } catch (error: any) {
      console.error('AI generation failed:', error)
      toast.error(error.response?.data?.error || 'AI生成に失敗しました')
      setStep('upload')
    } finally {
      setLoading(false)
    }
  }

  const handleApproveAndPublish = async () => {
    if (!aiData) return

    setLoading(true)

    try {
      // Approve and modify the AI-generated listing
      await api.post(`/ai-agent/listing/${aiData.product_id}/approve`, {
        approved: true,
        modifications: {
          title: editableData.title,
          description: editableData.description,
          price: editableData.price,
          category: editableData.category,
          condition: editableData.condition,
        },
      })

      toast.success('✅ 商品を出品しました！')
      navigate(`/products/${aiData.product_id}`)
    } catch (error: any) {
      toast.error(error.response?.data?.error || '出品に失敗しました')
    } finally {
      setLoading(false)
    }
  }

  const handleReject = () => {
    setStep('upload')
    setAiData(null)
    setUserHints('')
    toast('再度画像をアップロードしてください')
  }

  return (
    <div className="min-h-screen">
      <Header />
      <div className="py-8">
        <div className="max-w-4xl mx-auto px-4">
          {/* Progress Indicator */}
          <div className="mb-8">
            <div className="flex items-center justify-center gap-4">
              <div className={`flex items-center gap-2 ${step === 'upload' ? 'text-primary-600 font-bold' : 'text-gray-400'}`}>
                <div className={`w-8 h-8 rounded-full flex items-center justify-center ${step === 'upload' ? 'bg-primary-600 text-white' : 'bg-gray-300'}`}>1</div>
                <span>画像アップロード</span>
              </div>
              <div className="w-12 h-1 bg-gray-300"></div>
              <div className={`flex items-center gap-2 ${step === 'generating' ? 'text-primary-600 font-bold' : 'text-gray-400'}`}>
                <div className={`w-8 h-8 rounded-full flex items-center justify-center ${step === 'generating' ? 'bg-primary-600 text-white' : 'bg-gray-300'}`}>2</div>
                <span>AI生成中</span>
              </div>
              <div className="w-12 h-1 bg-gray-300"></div>
              <div className={`flex items-center gap-2 ${step === 'approval' ? 'text-primary-600 font-bold' : 'text-gray-400'}`}>
                <div className={`w-8 h-8 rounded-full flex items-center justify-center ${step === 'approval' ? 'bg-primary-600 text-white' : 'bg-gray-300'}`}>3</div>
                <span>確認・修正</span>
              </div>
            </div>
          </div>

          {/* Step 1: Upload */}
          {step === 'upload' && (
            <Card className="bg-white/80 backdrop-blur-sm">
              <div className="text-center mb-6">
                <h1 className="text-3xl font-bold mb-2">AI自動出品</h1>
                <p className="text-gray-600">画像をアップロードするだけで、AIが全ての設定を自動で行います</p>
              </div>

              <div className="space-y-6">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    商品画像をアップロード *
                  </label>
                  <ImageUpload
                    onUpload={setImageUrls}
                    maxImages={10}
                    existingImages={imageUrls}
                  />
                  <p className="text-xs text-gray-500 mt-2">
                    複数の角度から撮影した写真をアップロードすると、AIの精度が向上します
                  </p>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    ヒント（任意）
                  </label>
                  <textarea
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                    rows={3}
                    value={userHints}
                    onChange={(e) => setUserHints(e.target.value)}
                    placeholder="例: ブランド名、購入時期、特徴など（入力すると精度が上がります）"
                  />
                </div>

                <div className="bg-gradient-to-r from-blue-50 to-purple-50 p-4 rounded-lg border border-blue-200">
                  <div className="flex items-start gap-3">
                    <div>
                      <h3 className="font-bold text-gray-900 mb-1">AI自律エージェントが以下を自動で行います：</h3>
                      <ul className="text-sm text-gray-600 space-y-1">
                        <li>商品タイトルの生成</li>
                        <li>詳細説明の作成</li>
                        <li>カテゴリの自動選択</li>
                        <li>商品状態の判定</li>
                        <li>適正価格の算出</li>
                        <li>ブランド・モデルの検出</li>
                      </ul>
                    </div>
                  </div>
                </div>

                <Button
                  onClick={handleGenerateWithAI}
                  className="w-full"
                  size="lg"
                  disabled={imageUrls.length === 0}
                >
                  AIに任せて自動生成
                </Button>
              </div>
            </Card>
          )}

          {/* Step 2: Generating */}
          {step === 'generating' && (
            <Card className="bg-white/80 backdrop-blur-sm">
              <div className="text-center py-12">
                <div className="animate-spin rounded-full h-16 w-16 border-b-4 border-primary-500 mx-auto mb-6"></div>
                <h2 className="text-2xl font-bold mb-2">AI分析中...</h2>
                <p className="text-gray-600 mb-4">画像を解析して商品情報を生成しています</p>
                <div className="flex items-center justify-center gap-2 text-sm text-gray-500">
                  <div className="animate-pulse">画像認識</div>
                  <span>→</span>
                  <div className="animate-pulse delay-100">テキスト生成</div>
                  <span>→</span>
                  <div className="animate-pulse delay-200">価格算出</div>
                </div>
              </div>
            </Card>
          )}

          {/* Step 3: Approval */}
          {step === 'approval' && aiData && (
            <div className="space-y-6">
              <Card className="bg-white/80 backdrop-blur-sm">
                <div className="flex items-center justify-between mb-6">
                  <h2 className="text-2xl font-bold">AI生成結果</h2>
                  <div className="flex items-center gap-2">
                    <span className="text-sm text-gray-600">信頼度:</span>
                    <span className="text-lg font-bold text-primary-600">
                      {aiData.listing_data.ai_confidence_score.toFixed(0)}%
                    </span>
                  </div>
                </div>

                {/* Detected Brand/Model */}
                {(aiData.suggested_product.detected_brand || aiData.suggested_product.detected_model) && (
                  <div className="mb-6 p-4 bg-blue-50 rounded-lg border border-blue-200">
                    <div className="flex items-center gap-2 mb-2">
                      <h3 className="font-bold text-gray-900">検出情報</h3>
                    </div>
                    <div className="text-sm text-gray-700">
                      {aiData.suggested_product.detected_brand && (
                        <div>ブランド: <strong>{aiData.suggested_product.detected_brand}</strong></div>
                      )}
                      {aiData.suggested_product.detected_model && (
                        <div>モデル: <strong>{aiData.suggested_product.detected_model}</strong></div>
                      )}
                    </div>
                  </div>
                )}

                {/* Editable Fields */}
                <div className="space-y-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                      商品タイトル
                      <span className="ml-2 text-xs text-green-600">信頼度: {aiData.confidence_breakdown.title.toFixed(0)}%</span>
                    </label>
                    <input
                      type="text"
                      className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500"
                      value={editableData.title}
                      onChange={(e) => setEditableData({ ...editableData, title: e.target.value })}
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                      説明
                      <span className="ml-2 text-xs text-green-600">信頼度: {aiData.confidence_breakdown.description.toFixed(0)}%</span>
                    </label>
                    <textarea
                      className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500"
                      rows={6}
                      value={editableData.description}
                      onChange={(e) => setEditableData({ ...editableData, description: e.target.value })}
                    />
                  </div>

                  <div className="grid grid-cols-2 gap-4">
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">
                        価格 (円)
                        <span className="ml-2 text-xs text-green-600">信頼度: {aiData.confidence_breakdown.price.toFixed(0)}%</span>
                      </label>
                      <input
                        type="number"
                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500"
                        value={editableData.price}
                        onChange={(e) => setEditableData({ ...editableData, price: parseInt(e.target.value) })}
                      />
                      {aiData.suggested_product.pricing_rationale && (
                        <p className="text-xs text-gray-500 mt-1">{aiData.suggested_product.pricing_rationale}</p>
                      )}
                    </div>

                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">
                        カテゴリ
                        <span className="ml-2 text-xs text-green-600">信頼度: {aiData.confidence_breakdown.category.toFixed(0)}%</span>
                      </label>
                      <select
                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500"
                        value={editableData.category}
                        onChange={(e) => setEditableData({ ...editableData, category: e.target.value })}
                      >
                        <option value="clothing">Clothing</option>
                        <option value="electronics">Electronics</option>
                        <option value="furniture">Furniture</option>
                        <option value="books">Books</option>
                        <option value="toys">Toys</option>
                        <option value="sports">Sports</option>
                      </select>
                    </div>
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                      状態
                    </label>
                    <select
                      className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500"
                      value={editableData.condition}
                      onChange={(e) => setEditableData({ ...editableData, condition: e.target.value })}
                    >
                      <option value="new">New</option>
                      <option value="like_new">Like New</option>
                      <option value="good">Good</option>
                      <option value="fair">Fair</option>
                    </select>
                  </div>
                </div>

                {/* Key Features */}
                {aiData.suggested_product.key_features?.length > 0 && (
                  <div className="mt-6 p-4 bg-green-50 rounded-lg border border-green-200">
                    <h3 className="font-bold text-gray-900 mb-2">検出された特徴</h3>
                    <div className="flex flex-wrap gap-2">
                      {aiData.suggested_product.key_features.map((feature, index) => (
                        <span key={index} className="px-3 py-1 bg-white border border-green-300 rounded-full text-sm">
                          {feature}
                        </span>
                      ))}
                    </div>
                  </div>
                )}

                {/* Action Buttons */}
                <div className="flex gap-4 mt-6">
                  <Button
                    variant="outline"
                    onClick={handleReject}
                    className="flex-1"
                  >
                    やり直す
                  </Button>
                  <Button
                    onClick={handleApproveAndPublish}
                    isLoading={loading}
                    className="flex-1"
                  >
                    この内容で出品する
                  </Button>
                </div>
              </Card>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

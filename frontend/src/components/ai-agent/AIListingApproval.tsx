import React, { useState, useEffect } from 'react';

interface AIListingData {
  product_id: string;
  listing_data: {
    ai_confidence_score: number;
    generated_title: string;
    generated_description: string;
    generated_category: string;
    generated_condition: string;
    generated_price: number;
  };
  suggested_product: {
    title: string;
    description: string;
    category: string;
    condition: string;
    price: number;
    weight_kg: number;
    detected_brand?: string;
    detected_model?: string;
    key_features?: string[];
    pricing_rationale?: string;
    category_rationale?: string;
  };
  confidence_breakdown: {
    [key: string]: number;
  };
  requires_approval: boolean;
}

interface AIListingApprovalProps {
  listingData: AIListingData;
  onApprove: (modifications: any) => void;
  onCancel: () => void;
}

const AIListingApproval: React.FC<AIListingApprovalProps> = ({
  listingData,
  onApprove,
  onCancel,
}) => {
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [category, setCategory] = useState('');
  const [condition, setCondition] = useState('');
  const [price, setPrice] = useState(0);
  const [weightKg, setWeightKg] = useState(0);

  useEffect(() => {
    const product = listingData.suggested_product;
    setTitle(product.title);
    setDescription(product.description);
    setCategory(product.category);
    setCondition(product.condition);
    setPrice(product.price);
    setWeightKg(product.weight_kg);
  }, [listingData]);

  const handleApprove = () => {
    const modifications = {
      title,
      description,
      category,
      condition,
      price,
      weight_kg: weightKg,
    };
    onApprove(modifications);
  };

  const confidenceColor = (score: number) => {
    if (score >= 80) return 'text-green-600';
    if (score >= 60) return 'text-yellow-600';
    return 'text-red-600';
  };

  return (
    <div className="max-w-4xl mx-auto p-6 bg-white rounded-lg shadow-lg">
      {/* Header */}
      <div className="mb-6">
        <h2 className="text-2xl font-bold text-gray-900 mb-2">
          AI出品エージェント - 承認画面
        </h2>
        <p className="text-gray-600">
          AIが自動生成した商品情報を確認して、必要に応じて修正してください
        </p>
      </div>

      {/* AI Confidence Score */}
      <div className="mb-6 p-4 bg-blue-50 rounded-lg">
        <div className="flex items-center justify-between mb-2">
          <span className="font-semibold text-gray-700">AI確信度</span>
          <span className={`text-2xl font-bold ${confidenceColor(listingData.listing_data.ai_confidence_score)}`}>
            {listingData.listing_data.ai_confidence_score.toFixed(1)}%
          </span>
        </div>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-2 text-sm">
          {Object.entries(listingData.confidence_breakdown).map(([field, score]) => (
            <div key={field} className="flex justify-between">
              <span className="text-gray-600">{field}:</span>
              <span className={confidenceColor(score)}>{score.toFixed(0)}%</span>
            </div>
          ))}
        </div>
      </div>

      {/* Brand & Model Detection */}
      {(listingData.suggested_product.detected_brand || listingData.suggested_product.detected_model) && (
        <div className="mb-6 p-4 bg-purple-50 rounded-lg">
          <h3 className="font-semibold text-gray-700 mb-2">検出情報</h3>
          {listingData.suggested_product.detected_brand && (
            <p className="text-sm text-gray-600">ブランド: <span className="font-medium">{listingData.suggested_product.detected_brand}</span></p>
          )}
          {listingData.suggested_product.detected_model && (
            <p className="text-sm text-gray-600">モデル: <span className="font-medium">{listingData.suggested_product.detected_model}</span></p>
          )}
        </div>
      )}

      {/* Form Fields */}
      <div className="space-y-4 mb-6">
        {/* Title */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            商品名
          </label>
          <input
            type="text"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
        </div>

        {/* Description */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            説明文
          </label>
          <textarea
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            rows={6}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
        </div>

        {/* Category & Condition */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              カテゴリー
            </label>
            <select
              value={category}
              onChange={(e) => setCategory(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            >
              <option value="clothing">衣類</option>
              <option value="electronics">電子機器</option>
              <option value="furniture">家具</option>
              <option value="books">本</option>
              <option value="toys">おもちゃ</option>
              <option value="sports">スポーツ</option>
            </select>
            {listingData.suggested_product.category_rationale && (
              <p className="text-xs text-gray-500 mt-1">
                {listingData.suggested_product.category_rationale}
              </p>
            )}
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              状態
            </label>
            <select
              value={condition}
              onChange={(e) => setCondition(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            >
              <option value="new">新品</option>
              <option value="like_new">未使用に近い</option>
              <option value="good">良好</option>
              <option value="fair">可</option>
              <option value="poor">悪い</option>
            </select>
          </div>
        </div>

        {/* Price & Weight */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              価格 (円)
            </label>
            <input
              type="number"
              value={price}
              onChange={(e) => setPrice(parseInt(e.target.value))}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
            {listingData.suggested_product.pricing_rationale && (
              <p className="text-xs text-gray-500 mt-1">
                {listingData.suggested_product.pricing_rationale}
              </p>
            )}
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              重量 (kg)
            </label>
            <input
              type="number"
              step="0.1"
              value={weightKg}
              onChange={(e) => setWeightKg(parseFloat(e.target.value))}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
          </div>
        </div>
      </div>

      {/* Key Features */}
      {listingData.suggested_product.key_features && listingData.suggested_product.key_features.length > 0 && (
        <div className="mb-6 p-4 bg-green-50 rounded-lg">
          <h3 className="font-semibold text-gray-700 mb-2">検出された特徴</h3>
          <ul className="list-disc list-inside text-sm text-gray-600 space-y-1">
            {listingData.suggested_product.key_features.map((feature, index) => (
              <li key={index}>{feature}</li>
            ))}
          </ul>
        </div>
      )}

      {/* Action Buttons */}
      <div className="flex gap-4">
        <button
          onClick={handleApprove}
          className="flex-1 bg-blue-600 text-white py-3 px-6 rounded-lg font-semibold hover:bg-blue-700 transition-colors"
        >
          承認して出品
        </button>
        <button
          onClick={onCancel}
          className="flex-1 bg-gray-200 text-gray-700 py-3 px-6 rounded-lg font-semibold hover:bg-gray-300 transition-colors"
        >
          キャンセル
        </button>
      </div>
    </div>
  );
};

export default AIListingApproval;

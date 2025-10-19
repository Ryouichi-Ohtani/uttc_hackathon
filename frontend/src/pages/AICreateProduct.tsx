import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import AIListingApproval from '../components/ai-agent/AIListingApproval';

const AICreateProduct: React.FC = () => {
  const navigate = useNavigate();
  const [step, setStep] = useState<'upload' | 'generating' | 'approval'>('upload');
  const [images, setImages] = useState<File[]>([]);
  const [imagePreviewUrls, setImagePreviewUrls] = useState<string[]>([]);
  const [userHints, setUserHints] = useState('');
  const [listingData, setListingData] = useState<any>(null);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleImageSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files) {
      const files = Array.from(e.target.files);
      setImages(files);

      // プレビューURL作成
      const urls = files.map((file) => URL.createObjectURL(file));
      setImagePreviewUrls(urls);
    }
  };

  const handleRemoveImage = (index: number) => {
    const newImages = images.filter((_, i) => i !== index);
    const newUrls = imagePreviewUrls.filter((_, i) => i !== index);
    setImages(newImages);
    setImagePreviewUrls(newUrls);

    // URL解放
    URL.revokeObjectURL(imagePreviewUrls[index]);
  };

  const handleGenerateListing = async () => {
    if (images.length === 0) {
      setError('少なくとも1枚の画像をアップロードしてください');
      return;
    }

    setLoading(true);
    setStep('generating');
    setError('');

    try {
      // Step 1: 画像をアップロード
      const formData = new FormData();
      images.forEach((image) => {
        formData.append('images', image);
      });

      const uploadResponse = await axios.post('/api/v1/upload/images', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
          Authorization: `Bearer ${localStorage.getItem('token')}`,
        },
      });

      const imageUrls = uploadResponse.data.urls;

      // Step 2: AI出品生成
      const generateResponse = await axios.post(
        '/api/v1/ai-agent/listing/generate',
        {
          image_urls: imageUrls,
          user_hints: userHints,
          auto_publish: false,
        },
        {
          headers: { Authorization: `Bearer ${localStorage.getItem('token')}` },
        }
      );

      setListingData(generateResponse.data);
      setStep('approval');
    } catch (err: any) {
      console.error('AI listing generation failed:', err);
      setError(
        err.response?.data?.error || 'AI生成に失敗しました。再度お試しください。'
      );
      setStep('upload');
    } finally {
      setLoading(false);
    }
  };

  const handleApprove = async (modifications: any) => {
    setLoading(true);
    try {
      await axios.post(
        `/api/v1/ai-agent/listing/${listingData.product_id}/approve`,
        modifications,
        {
          headers: { Authorization: `Bearer ${localStorage.getItem('token')}` },
        }
      );

      alert('商品が正常に出品されました！');
      navigate(`/products/${listingData.product_id}`);
    } catch (err) {
      console.error('Approval failed:', err);
      alert('承認に失敗しました');
    } finally {
      setLoading(false);
    }
  };

  const handleCancel = () => {
    if (window.confirm('生成された内容を破棄して最初からやり直しますか？')) {
      setStep('upload');
      setListingData(null);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-4xl mx-auto px-4">
        {/* Progress Indicator */}
        <div className="mb-8">
          <div className="flex items-center justify-center space-x-4">
            <div className={`flex items-center ${step === 'upload' ? 'text-blue-600' : 'text-gray-400'}`}>
              <div className={`w-10 h-10 rounded-full flex items-center justify-center ${step === 'upload' ? 'bg-blue-600 text-white' : 'bg-gray-300'}`}>
                1
              </div>
              <span className="ml-2 font-medium">画像アップロード</span>
            </div>
            <div className="flex-1 h-1 bg-gray-300"></div>
            <div className={`flex items-center ${step === 'generating' ? 'text-blue-600' : 'text-gray-400'}`}>
              <div className={`w-10 h-10 rounded-full flex items-center justify-center ${step === 'generating' ? 'bg-blue-600 text-white animate-pulse' : 'bg-gray-300'}`}>
                2
              </div>
              <span className="ml-2 font-medium">AI生成中</span>
            </div>
            <div className="flex-1 h-1 bg-gray-300"></div>
            <div className={`flex items-center ${step === 'approval' ? 'text-blue-600' : 'text-gray-400'}`}>
              <div className={`w-10 h-10 rounded-full flex items-center justify-center ${step === 'approval' ? 'bg-blue-600 text-white' : 'bg-gray-300'}`}>
                3
              </div>
              <span className="ml-2 font-medium">承認</span>
            </div>
          </div>
        </div>

        {/* Error Message */}
        {error && (
          <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg text-red-700">
            {error}
          </div>
        )}

        {/* Upload Step */}
        {step === 'upload' && (
          <div className="bg-white rounded-lg shadow-lg p-8">
            <h1 className="text-3xl font-bold text-gray-900 mb-2">
              AI出品エージェント
            </h1>
            <p className="text-gray-600 mb-8">
              商品の写真をアップロードするだけで、AIがすべての情報を自動生成します
            </p>

            {/* Image Upload */}
            <div className="mb-6">
              <label className="block text-sm font-medium text-gray-700 mb-2">
                商品画像 *
              </label>
              <div className="border-2 border-dashed border-gray-300 rounded-lg p-8 text-center hover:border-blue-500 transition-colors cursor-pointer">
                <input
                  type="file"
                  accept="image/*"
                  multiple
                  onChange={handleImageSelect}
                  className="hidden"
                  id="image-upload"
                />
                <label htmlFor="image-upload" className="cursor-pointer">
                  <div className="text-gray-600">
                    <svg
                      className="mx-auto h-12 w-12 mb-4"
                      stroke="currentColor"
                      fill="none"
                      viewBox="0 0 48 48"
                    >
                      <path
                        d="M28 8H12a4 4 0 00-4 4v20m32-12v8m0 0v8a4 4 0 01-4 4H12a4 4 0 01-4-4v-4m32-4l-3.172-3.172a4 4 0 00-5.656 0L28 28M8 32l9.172-9.172a4 4 0 015.656 0L28 28m0 0l4 4m4-24h8m-4-4v8m-12 4h.02"
                        strokeWidth={2}
                        strokeLinecap="round"
                        strokeLinejoin="round"
                      />
                    </svg>
                    <p className="text-lg font-medium">クリックして画像を選択</p>
                    <p className="text-sm text-gray-500">または、ドラッグ&ドロップ</p>
                  </div>
                </label>
              </div>

              {/* Image Previews */}
              {imagePreviewUrls.length > 0 && (
                <div className="grid grid-cols-3 gap-4 mt-4">
                  {imagePreviewUrls.map((url, index) => (
                    <div key={index} className="relative">
                      <img
                        src={url}
                        alt={`Preview ${index + 1}`}
                        className="w-full h-32 object-cover rounded-lg"
                      />
                      <button
                        onClick={() => handleRemoveImage(index)}
                        className="absolute top-2 right-2 bg-red-500 text-white rounded-full p-1 hover:bg-red-600"
                      >
                        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                        </svg>
                      </button>
                    </div>
                  ))}
                </div>
              )}
            </div>

            {/* User Hints */}
            <div className="mb-8">
              <label className="block text-sm font-medium text-gray-700 mb-2">
                ヒント（任意）
              </label>
              <input
                type="text"
                value={userHints}
                onChange={(e) => setUserHints(e.target.value)}
                placeholder="例: iPhone 13 Pro 128GB、ほぼ未使用"
                className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
              <p className="text-sm text-gray-500 mt-1">
                商品名やブランド、特徴などを入力するとAIの精度が向上します
              </p>
            </div>

            {/* Generate Button */}
            <button
              onClick={handleGenerateListing}
              disabled={images.length === 0 || loading}
              className="w-full bg-blue-600 text-white py-4 rounded-lg font-bold text-lg hover:bg-blue-700 disabled:bg-gray-300 disabled:cursor-not-allowed transition-colors"
            >
              {loading ? 'AI生成中...' : 'AIで自動生成'}
            </button>

            <p className="text-center text-sm text-gray-500 mt-4">
              所要時間: 約5-10秒
            </p>
          </div>
        )}

        {/* Generating Step */}
        {step === 'generating' && (
          <div className="bg-white rounded-lg shadow-lg p-12 text-center">
            <div className="animate-spin rounded-full h-16 w-16 border-b-4 border-blue-600 mx-auto mb-6"></div>
            <h2 className="text-2xl font-bold text-gray-900 mb-2">
              AIが商品情報を生成中...
            </h2>
            <p className="text-gray-600">
              画像を分析して、タイトル・説明・価格・カテゴリーを自動生成しています
            </p>
          </div>
        )}

        {/* Approval Step */}
        {step === 'approval' && listingData && (
          <AIListingApproval
            listingData={listingData}
            onApprove={handleApprove}
            onCancel={handleCancel}
          />
        )}
      </div>
    </div>
  );
};

export default AICreateProduct;

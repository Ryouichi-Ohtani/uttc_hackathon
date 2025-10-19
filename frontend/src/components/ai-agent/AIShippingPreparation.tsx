import React, { useState, useEffect } from 'react';
import api from '@/services/api';
import { Card } from '@/components/common/Card';
import { Button } from '@/components/common/Button';
import toast from 'react-hot-toast';

interface AIShippingPreparation {
  id: string;
  purchase_id: string;
  is_ai_prepared: boolean;
  suggested_carrier: string;
  suggested_package_size: string;
  estimated_weight: number;
  estimated_cost: number;
  shipping_instructions: string;
  user_approved: boolean;
  approved_at?: string;
}

interface AIShippingPreparationProps {
  purchaseId: string;
  onApprove?: () => void;
}

const AIShippingPreparation: React.FC<AIShippingPreparationProps> = ({
  purchaseId,
  onApprove,
}) => {
  const [preparation, setPreparation] = useState<AIShippingPreparation | null>(null);
  const [loading, setLoading] = useState(false);
  const [generating, setGenerating] = useState(false);
  const [showEdit, setShowEdit] = useState(false);

  // 編集可能フィールド
  const [carrier, setCarrier] = useState('');
  const [packageSize, setPackageSize] = useState('');

  useEffect(() => {
    loadPreparation();
  }, [purchaseId]);

  const loadPreparation = async () => {
    try {
      setLoading(true);
      const response = await api.get(`/ai-agent/shipping/${purchaseId}`);
      setPreparation(response.data);
      setCarrier(response.data.suggested_carrier);
      setPackageSize(response.data.suggested_package_size);
    } catch (error: any) {
      if (error.response?.status === 404) {
        // 準備データがまだない場合
        setPreparation(null);
      } else {
        console.error('Failed to load shipping preparation:', error);
      }
    } finally {
      setLoading(false);
    }
  };

  const handleGenerate = async () => {
    setGenerating(true);
    try {
      const response = await api.post('/ai-agent/shipping/prepare', {
        purchase_id: purchaseId
      });
      setPreparation(response.data);
      setCarrier(response.data.suggested_carrier);
      setPackageSize(response.data.suggested_package_size);
      toast.success('AIが配送情報を準備しました！');
    } catch (error: any) {
      toast.error(error.response?.data?.error || 'AI準備に失敗しました');
    } finally {
      setGenerating(false);
    }
  };

  const handleApprove = async () => {
    if (!preparation) return;

    setLoading(true);
    try {
      await api.post(`/ai-agent/shipping/${purchaseId}/approve`, {
        approved: true,
        carrier: carrier !== preparation.suggested_carrier ? carrier : '',
        package_size: packageSize !== preparation.suggested_package_size ? packageSize : '',
        modifications: JSON.stringify({
          carrier_changed: carrier !== preparation.suggested_carrier,
          package_size_changed: packageSize !== preparation.suggested_package_size,
        }),
      });
      toast.success('配送情報を承認しました！');
      await loadPreparation();
      if (onApprove) onApprove();
    } catch (error: any) {
      toast.error(error.response?.data?.error || '承認に失敗しました');
    } finally {
      setLoading(false);
    }
  };

  if (loading && !preparation) {
    return (
      <Card>
        <div className="flex items-center justify-center py-8">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-500" />
        </div>
      </Card>
    );
  }

  // 配送準備データがない場合
  if (!preparation) {
    return (
      <Card className="bg-gradient-to-r from-blue-50 to-purple-50 border-2 border-blue-200">
        <div className="flex items-start gap-4">
          <div className="flex-1">
            <h3 className="text-lg font-bold text-gray-900 mb-2">
              AI配送準備エージェント
            </h3>
            <p className="text-sm text-gray-600 mb-4">
              AIが商品情報から最適な配送方法を自動で提案します。
              配送業者、パッケージサイズ、配送料、梱包指示まで一括準備！
            </p>
            <Button
              variant="primary"
              onClick={handleGenerate}
              disabled={generating}
              className="w-full sm:w-auto"
            >
              {generating ? (
                <>
                  <span className="animate-spin mr-2">...</span>
                  AI準備中...
                </>
              ) : (
                <>
                  AIに配送情報を準備させる
                </>
              )}
            </Button>
          </div>
        </div>
      </Card>
    );
  }

  // 承認済みの場合
  if (preparation.user_approved) {
    return (
      <Card className="bg-green-50 border-2 border-green-200">
        <div className="flex items-start gap-3">
          <div className="flex-1">
            <h3 className="text-lg font-bold text-green-900 mb-1">
              配送準備完了
            </h3>
            <p className="text-sm text-green-700 mb-3">
              承認日時: {new Date(preparation.approved_at!).toLocaleString('ja-JP')}
            </p>
            <div className="grid grid-cols-2 gap-3 text-sm">
              <div>
                <p className="text-gray-600 text-xs mb-1">配送業者</p>
                <p className="font-medium text-gray-900">{carrier}</p>
              </div>
              <div>
                <p className="text-gray-600 text-xs mb-1">パッケージサイズ</p>
                <p className="font-medium text-gray-900">{packageSize}</p>
              </div>
            </div>
          </div>
        </div>
      </Card>
    );
  }

  // 承認待ちの場合（編集・承認フロー）
  return (
    <Card className="bg-white border-2 border-blue-300">
      <div className="space-y-4">
        {/* ヘッダー */}
        <div className="flex items-start justify-between">
          <div className="flex items-start gap-3">
            <div>
              <h3 className="text-lg font-bold text-gray-900">
                AI配送準備完了
              </h3>
              <p className="text-sm text-gray-600">
                内容を確認して承認してください
              </p>
            </div>
          </div>
          <span className="px-3 py-1 bg-blue-100 text-blue-800 text-xs font-semibold rounded-full">
            承認待ち
          </span>
        </div>

        {/* AI推奨情報 */}
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 p-4 bg-blue-50 rounded-lg">
          <div>
            <label className="block text-xs text-gray-600 mb-1">
              配送業者 {carrier !== preparation.suggested_carrier && <span className="text-orange-600">(変更済み)</span>}
            </label>
            {showEdit ? (
              <input
                type="text"
                value={carrier}
                onChange={(e) => setCarrier(e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded text-sm"
                placeholder="例: ヤマト運輸（宅急便）"
              />
            ) : (
              <p className="font-semibold text-gray-900">{carrier}</p>
            )}
          </div>

          <div>
            <label className="block text-xs text-gray-600 mb-1">
              パッケージサイズ {packageSize !== preparation.suggested_package_size && <span className="text-orange-600">(変更済み)</span>}
            </label>
            {showEdit ? (
              <select
                value={packageSize}
                onChange={(e) => setPackageSize(e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded text-sm"
              >
                <option value="60サイズ">60サイズ</option>
                <option value="80サイズ">80サイズ</option>
                <option value="100サイズ">100サイズ</option>
                <option value="120サイズ">120サイズ</option>
                <option value="140サイズ">140サイズ</option>
                <option value="160サイズ">160サイズ</option>
              </select>
            ) : (
              <p className="font-semibold text-gray-900">{packageSize}</p>
            )}
          </div>

          <div>
            <label className="block text-xs text-gray-600 mb-1">推定重量</label>
            <p className="font-semibold text-gray-900">{preparation.estimated_weight.toFixed(2)} kg</p>
          </div>

          <div>
            <label className="block text-xs text-gray-600 mb-1">推定配送料</label>
            <p className="font-semibold text-gray-900">¥{preparation.estimated_cost.toLocaleString()}</p>
          </div>
        </div>

        {/* 梱包指示 */}
        <div className="p-4 bg-yellow-50 border border-yellow-200 rounded-lg">
          <div className="flex items-start gap-2 mb-2">
            <h4 className="font-semibold text-gray-900 text-sm">梱包指示</h4>
          </div>
          <p className="text-sm text-gray-700 leading-relaxed">
            {preparation.shipping_instructions}
          </p>
        </div>

        {/* アクションボタン */}
        <div className="flex gap-3">
          <Button
            variant="outline"
            onClick={() => setShowEdit(!showEdit)}
            className="flex-1"
            disabled={loading}
          >
            {showEdit ? '編集を終了' : '修正する'}
          </Button>
          <Button
            variant="primary"
            onClick={handleApprove}
            className="flex-1"
            disabled={loading}
          >
            {loading ? (
              <>
                <span className="animate-spin mr-2">...</span>
                承認中...
              </>
            ) : (
              <>
                この内容で承認
              </>
            )}
          </Button>
        </div>

        {/* 時間節約の表示 */}
        <div className="flex items-center justify-center gap-2 text-xs text-gray-500 pt-2 border-t">
          <span>AIが約5分の作業時間を節約しました</span>
        </div>
      </div>
    </Card>
  );
};

export default AIShippingPreparation;

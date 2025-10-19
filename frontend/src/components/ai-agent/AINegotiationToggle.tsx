import React, { useState, useEffect } from 'react';
import axios from 'axios';

interface NegotiationSettings {
  mode: 'ai' | 'manual' | 'hybrid';
  is_enabled: boolean;
  min_acceptable_price: number;
  auto_accept_threshold: number;
  auto_reject_threshold: number;
  negotiation_strategy: 'aggressive' | 'moderate' | 'conservative';
  total_offers_processed: number;
  ai_accepted_count: number;
  ai_rejected_count: number;
}

interface AINegotiationToggleProps {
  productId: string;
  productPrice: number;
  onUpdate?: () => void;
}

const AINegotiationToggle: React.FC<AINegotiationToggleProps> = ({
  productId,
  productPrice,
  onUpdate,
}) => {
  const [settings, setSettings] = useState<NegotiationSettings | null>(null);
  const [mode, setMode] = useState<'ai' | 'manual' | 'hybrid'>('manual');
  const [minPrice, setMinPrice] = useState(Math.floor(productPrice * 0.7));
  const [autoAccept, setAutoAccept] = useState(Math.floor(productPrice * 0.95));
  const [autoReject, setAutoReject] = useState(Math.floor(productPrice * 0.6));
  const [strategy, setStrategy] = useState<'aggressive' | 'moderate' | 'conservative'>('moderate');
  const [loading, setLoading] = useState(false);
  const [showConfig, setShowConfig] = useState(false);

  useEffect(() => {
    fetchSettings();
  }, [productId]);

  const fetchSettings = async () => {
    try {
      const response = await axios.get(`/api/v1/ai-agent/negotiation/${productId}`, {
        headers: { Authorization: `Bearer ${localStorage.getItem('token')}` },
      });
      setSettings(response.data);
      setMode(response.data.mode);
      setMinPrice(response.data.min_acceptable_price);
      setAutoAccept(response.data.auto_accept_threshold);
      setAutoReject(response.data.auto_reject_threshold);
      setStrategy(response.data.negotiation_strategy);
    } catch (error) {
      // Settings don't exist yet
      setSettings(null);
    }
  };

  const handleToggle = async (newMode: 'ai' | 'manual' | 'hybrid') => {
    setLoading(true);
    try {
      await axios.post(
        '/api/v1/ai-agent/negotiation/enable',
        {
          product_id: productId,
          mode: newMode,
          min_acceptable_price: minPrice,
          auto_accept_threshold: autoAccept,
          auto_reject_threshold: autoReject,
          strategy: strategy,
        },
        {
          headers: { Authorization: `Bearer ${localStorage.getItem('token')}` },
        }
      );
      setMode(newMode);
      await fetchSettings();
      if (onUpdate) onUpdate();
    } catch (error) {
      console.error('Failed to update negotiation settings:', error);
      alert('設定の更新に失敗しました');
    } finally {
      setLoading(false);
    }
  };

  const handleDisable = async () => {
    setLoading(true);
    try {
      await axios.delete(`/api/v1/ai-agent/negotiation/${productId}`, {
        headers: { Authorization: `Bearer ${localStorage.getItem('token')}` },
      });
      setMode('manual');
      setSettings(null);
      if (onUpdate) onUpdate();
    } catch (error) {
      console.error('Failed to disable negotiation:', error);
      alert('無効化に失敗しました');
    } finally {
      setLoading(false);
    }
  };

  const acceptanceRate = settings && settings.total_offers_processed > 0
    ? (settings.ai_accepted_count / settings.total_offers_processed * 100).toFixed(1)
    : '0';

  return (
    <div className="bg-white rounded-lg shadow p-6">
      <div className="flex items-center justify-between mb-4">
        <div>
          <h3 className="text-lg font-bold text-gray-900">AI交渉エージェント</h3>
          <p className="text-sm text-gray-600">価格交渉をAIに任せることができます</p>
        </div>
        {settings && (
          <span className={`px-3 py-1 rounded-full text-sm font-semibold ${
            mode === 'ai' ? 'bg-blue-100 text-blue-800' : 'bg-gray-100 text-gray-800'
          }`}>
            {mode === 'ai' ? 'AI自動' : mode === 'hybrid' ? 'AIアシスト' : '手動'}
          </span>
        )}
      </div>

      {/* Mode Selection */}
      <div className="grid grid-cols-3 gap-2 mb-4">
        <button
          onClick={() => handleToggle('ai')}
          disabled={loading}
          className={`p-3 rounded-lg border-2 transition-all ${
            mode === 'ai'
              ? 'border-blue-500 bg-blue-50'
              : 'border-gray-200 hover:border-gray-300'
          }`}
        >
          <div className="font-semibold text-sm mb-1">AI自動</div>
          <div className="text-xs text-gray-600">完全自動交渉</div>
        </button>

        <button
          onClick={() => handleToggle('hybrid')}
          disabled={loading}
          className={`p-3 rounded-lg border-2 transition-all ${
            mode === 'hybrid'
              ? 'border-blue-500 bg-blue-50'
              : 'border-gray-200 hover:border-gray-300'
          }`}
        >
          <div className="font-semibold text-sm mb-1">AIアシスト</div>
          <div className="text-xs text-gray-600">提案+手動承認</div>
        </button>

        <button
          onClick={() => handleToggle('manual')}
          disabled={loading}
          className={`p-3 rounded-lg border-2 transition-all ${
            mode === 'manual'
              ? 'border-blue-500 bg-blue-50'
              : 'border-gray-200 hover:border-gray-300'
          }`}
        >
          <div className="font-semibold text-sm mb-1">手動</div>
          <div className="text-xs text-gray-600">すべて自分で</div>
        </button>
      </div>

      {/* Configuration */}
      {(mode === 'ai' || mode === 'hybrid') && (
        <>
          <button
            onClick={() => setShowConfig(!showConfig)}
            className="w-full text-left text-sm text-blue-600 hover:text-blue-700 font-medium mb-2"
          >
            {showConfig ? '▼' : '▶'} 詳細設定
          </button>

          {showConfig && (
            <div className="space-y-4 p-4 bg-gray-50 rounded-lg mb-4">
              {/* Strategy */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  交渉戦略
                </label>
                <div className="grid grid-cols-3 gap-2">
                  {['aggressive', 'moderate', 'conservative'].map((s) => (
                    <button
                      key={s}
                      onClick={() => setStrategy(s as any)}
                      className={`px-3 py-2 rounded text-sm ${
                        strategy === s
                          ? 'bg-blue-600 text-white'
                          : 'bg-white text-gray-700 border border-gray-300'
                      }`}
                    >
                      {s === 'aggressive' ? '攻撃的' : s === 'moderate' ? '中立' : '保守的'}
                    </button>
                  ))}
                </div>
              </div>

              {/* Price Thresholds */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  自動承認価格 (この金額以上は即承認)
                </label>
                <input
                  type="number"
                  value={autoAccept}
                  onChange={(e) => setAutoAccept(parseInt(e.target.value))}
                  className="w-full px-3 py-2 border border-gray-300 rounded"
                />
                <p className="text-xs text-gray-500 mt-1">
                  定価の {((autoAccept / productPrice) * 100).toFixed(0)}%
                </p>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  最低価格 (これ以下は自動拒否)
                </label>
                <input
                  type="number"
                  value={autoReject}
                  onChange={(e) => setAutoReject(parseInt(e.target.value))}
                  className="w-full px-3 py-2 border border-gray-300 rounded"
                />
                <p className="text-xs text-gray-500 mt-1">
                  定価の {((autoReject / productPrice) * 100).toFixed(0)}%
                </p>
              </div>

              <button
                onClick={() => handleToggle(mode)}
                disabled={loading}
                className="w-full bg-blue-600 text-white py-2 rounded font-semibold hover:bg-blue-700 transition-colors"
              >
                設定を保存
              </button>
            </div>
          )}
        </>
      )}

      {/* Statistics */}
      {settings && settings.total_offers_processed > 0 && (
        <div className="border-t pt-4">
          <h4 className="text-sm font-semibold text-gray-700 mb-2">AI交渉統計</h4>
          <div className="grid grid-cols-3 gap-4 text-center">
            <div>
              <div className="text-2xl font-bold text-gray-900">{settings.total_offers_processed}</div>
              <div className="text-xs text-gray-600">処理済み</div>
            </div>
            <div>
              <div className="text-2xl font-bold text-green-600">{settings.ai_accepted_count}</div>
              <div className="text-xs text-gray-600">承認</div>
            </div>
            <div>
              <div className="text-2xl font-bold text-blue-600">{acceptanceRate}%</div>
              <div className="text-xs text-gray-600">承認率</div>
            </div>
          </div>
        </div>
      )}

      {/* Disable Button */}
      {settings && mode !== 'manual' && (
        <button
          onClick={handleDisable}
          disabled={loading}
          className="w-full mt-4 text-sm text-red-600 hover:text-red-700 font-medium"
        >
          AI交渉を無効化
        </button>
      )}
    </div>
  );
};

export default AINegotiationToggle;

import React, { useEffect, useState } from 'react';
import { Line } from 'react-chartjs-2';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler,
} from 'chart.js';

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler
);

interface PricePredictionData {
  product_id: string;
  predicted_price: number;
  confidence_score: number;
  optimal_list_price: number;
  demand_score: number;
  estimated_days_to_sell: number;
  similar_products: Array<{
    product_id: string;
    sold_price: number;
    days_to_sell: number;
    similarity_score: number;
  }>;
  trend_analysis: {
    category_trend: 'rising' | 'stable' | 'declining';
    seasonal_boost: number;
    market_demand: number;
  };
}

interface PricePredictionProps {
  productId: string;
  currentPrice: number;
}

export const PricePrediction: React.FC<PricePredictionProps> = ({
  productId,
  currentPrice,
}) => {
  const [prediction, setPrediction] = useState<PricePredictionData | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Fetch price prediction from API
    fetchPrediction();
  }, [productId]);

  const fetchPrediction = async () => {
    try {
      // Mock data for demonstration
      const mockPrediction: PricePredictionData = {
        product_id: productId,
        predicted_price: Math.floor(currentPrice * (0.9 + Math.random() * 0.2)),
        confidence_score: 0.75 + Math.random() * 0.2,
        optimal_list_price: Math.floor(currentPrice * 1.05),
        demand_score: 0.6 + Math.random() * 0.3,
        estimated_days_to_sell: Math.floor(5 + Math.random() * 20),
        similar_products: [
          {
            product_id: 'similar1',
            sold_price: currentPrice - 5000,
            days_to_sell: 7,
            similarity_score: 0.92,
          },
          {
            product_id: 'similar2',
            sold_price: currentPrice + 3000,
            days_to_sell: 12,
            similarity_score: 0.85,
          },
        ],
        trend_analysis: {
          category_trend: Math.random() > 0.5 ? 'rising' : 'stable',
          seasonal_boost: Math.random() * 0.15,
          market_demand: 0.6 + Math.random() * 0.3,
        },
      };

      setPrediction(mockPrediction);
      setLoading(false);
    } catch (error) {
      console.error('Error fetching prediction:', error);
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="bg-white rounded-lg shadow-lg p-6">
        <div className="animate-pulse">
          <div className="h-8 bg-gray-200 rounded w-1/3 mb-4"></div>
          <div className="h-32 bg-gray-200 rounded"></div>
        </div>
      </div>
    );
  }

  if (!prediction) {
    return null;
  }

  const getTrendColor = (trend: string) => {
    switch (trend) {
      case 'rising':
        return 'text-green-600';
      case 'declining':
        return 'text-red-600';
      default:
        return 'text-gray-600';
    }
  };

  const getTrendIcon = (trend: string) => {
    switch (trend) {
      case 'rising':
        return '📈';
      case 'declining':
        return '📉';
      default:
        return '➡️';
    }
  };

  const chartData = {
    labels: ['現在の価格', '予測価格', '最適価格'],
    datasets: [
      {
        label: '価格 (¥)',
        data: [currentPrice, prediction.predicted_price, prediction.optimal_list_price],
        borderColor: 'rgb(34, 197, 94)',
        backgroundColor: 'rgba(34, 197, 94, 0.1)',
        fill: true,
      },
    ],
  };

  const chartOptions = {
    responsive: true,
    plugins: {
      legend: {
        display: false,
      },
      title: {
        display: false,
      },
    },
    scales: {
      y: {
        beginAtZero: false,
      },
    },
  };

  return (
    <div className="bg-white rounded-lg shadow-lg p-6">
      <h3 className="text-2xl font-bold mb-6">🤖 AI価格予測</h3>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
        <div className="bg-blue-50 rounded-lg p-4">
          <div className="text-sm text-gray-600 mb-1">予測価格</div>
          <div className="text-2xl font-bold text-blue-600">
            ¥{prediction.predicted_price.toLocaleString()}
          </div>
          <div className="text-xs text-gray-500 mt-1">
            信頼度: {(prediction.confidence_score * 100).toFixed(0)}%
          </div>
        </div>

        <div className="bg-green-50 rounded-lg p-4">
          <div className="text-sm text-gray-600 mb-1">最適出品価格</div>
          <div className="text-2xl font-bold text-green-600">
            ¥{prediction.optimal_list_price.toLocaleString()}
          </div>
          <div className="text-xs text-gray-500 mt-1">
            推奨価格帯
          </div>
        </div>

        <div className="bg-purple-50 rounded-lg p-4">
          <div className="text-sm text-gray-600 mb-1">予想販売日数</div>
          <div className="text-2xl font-bold text-purple-600">
            {prediction.estimated_days_to_sell}日
          </div>
          <div className="text-xs text-gray-500 mt-1">
            需要スコア: {(prediction.demand_score * 100).toFixed(0)}%
          </div>
        </div>
      </div>

      <div className="mb-6">
        <Line data={chartData} options={chartOptions} />
      </div>

      <div className="border-t pt-4">
        <h4 className="font-semibold mb-3">市場トレンド分析</h4>
        <div className="space-y-2">
          <div className="flex justify-between items-center">
            <span className="text-gray-600">カテゴリートレンド</span>
            <span className={`font-semibold ${getTrendColor(prediction.trend_analysis.category_trend)}`}>
              {getTrendIcon(prediction.trend_analysis.category_trend)}{' '}
              {prediction.trend_analysis.category_trend === 'rising'
                ? '上昇中'
                : prediction.trend_analysis.category_trend === 'declining'
                ? '下降中'
                : '安定'}
            </span>
          </div>
          <div className="flex justify-between items-center">
            <span className="text-gray-600">季節的要因</span>
            <span className="font-semibold">
              {prediction.trend_analysis.seasonal_boost > 0 ? '+' : ''}
              {(prediction.trend_analysis.seasonal_boost * 100).toFixed(1)}%
            </span>
          </div>
          <div className="flex justify-between items-center">
            <span className="text-gray-600">市場需要</span>
            <div className="w-32 bg-gray-200 rounded-full h-2">
              <div
                className="bg-green-600 h-2 rounded-full"
                style={{ width: `${prediction.trend_analysis.market_demand * 100}%` }}
              ></div>
            </div>
          </div>
        </div>
      </div>

      <div className="border-t pt-4 mt-4">
        <h4 className="font-semibold mb-3">類似商品の販売実績</h4>
        <div className="space-y-2">
          {prediction.similar_products.map((product, index) => (
            <div key={index} className="flex justify-between items-center bg-gray-50 p-3 rounded">
              <div>
                <div className="text-sm text-gray-600">
                  類似度: {(product.similarity_score * 100).toFixed(0)}%
                </div>
                <div className="text-xs text-gray-500">
                  販売日数: {product.days_to_sell}日
                </div>
              </div>
              <div className="text-lg font-bold text-gray-800">
                ¥{product.sold_price.toLocaleString()}
              </div>
            </div>
          ))}
        </div>
      </div>

      <div className="mt-6 bg-yellow-50 border border-yellow-200 rounded-lg p-4">
        <div className="flex items-start gap-2">
          <span className="text-yellow-600">💡</span>
          <div className="text-sm text-yellow-800">
            <p className="font-semibold mb-1">価格設定のヒント</p>
            <p>
              現在の市場トレンドと類似商品の実績から、
              ¥{prediction.optimal_list_price.toLocaleString()}での出品を推奨します。
              この価格なら約{prediction.estimated_days_to_sell}日以内に売却できる可能性が高いです。
            </p>
          </div>
        </div>
      </div>
    </div>
  );
};

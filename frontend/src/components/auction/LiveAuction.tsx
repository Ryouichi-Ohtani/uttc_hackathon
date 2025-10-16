import React, { useState, useEffect, useRef } from 'react';

interface Bid {
  id: string;
  user_id: string;
  username: string;
  amount: number;
  timestamp: string;
}

interface LiveAuctionProps {
  productId: string;
  startingPrice: number;
  currentPrice: number;
  endTime: string;
}

export const LiveAuction: React.FC<LiveAuctionProps> = ({
  productId,
  startingPrice: _startingPrice,
  currentPrice: initialPrice,
  endTime,
}) => {
  const [currentPrice, setCurrentPrice] = useState(initialPrice);
  const [bidAmount, setBidAmount] = useState(initialPrice + 1000);
  const [bids, setBids] = useState<Bid[]>([]);
  const [timeRemaining, setTimeRemaining] = useState('');
  const [isEnded, setIsEnded] = useState(false);
  const [highestBidder, setHighestBidder] = useState<string | null>(null);

  // WebSocket for real-time bids
  const wsRef = useRef<WebSocket | null>(null);

  useEffect(() => {
    const ws = new WebSocket(`ws://localhost:8080/v1/ws/auction/${productId}`);

    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);

      if (data.type === 'bid') {
        const newBid: Bid = data.bid;
        setBids((prev) => [newBid, ...prev]);
        setCurrentPrice(newBid.amount);
        setHighestBidder(newBid.username);
        setBidAmount(newBid.amount + 1000);
      } else if (data.type === 'auction_end') {
        setIsEnded(true);
      }
    };

    wsRef.current = ws;

    return () => {
      ws.close();
    };
  }, [productId]);

  useEffect(() => {
    const timer = setInterval(() => {
      const end = new Date(endTime).getTime();
      const now = new Date().getTime();
      const diff = end - now;

      if (diff <= 0) {
        setTimeRemaining('終了');
        setIsEnded(true);
        clearInterval(timer);
      } else {
        const hours = Math.floor(diff / (1000 * 60 * 60));
        const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));
        const seconds = Math.floor((diff % (1000 * 60)) / 1000);
        setTimeRemaining(`${hours}時間 ${minutes}分 ${seconds}秒`);
      }
    }, 1000);

    return () => clearInterval(timer);
  }, [endTime]);

  const placeBid = () => {
    if (wsRef.current && bidAmount > currentPrice) {
      wsRef.current.send(
        JSON.stringify({
          type: 'place_bid',
          amount: bidAmount,
        })
      );
    }
  };

  const quickBid = (increment: number) => {
    setBidAmount(currentPrice + increment);
  };

  return (
    <div className="bg-white rounded-lg shadow-lg p-6">
      <div className="mb-6">
        <h2 className="text-2xl font-bold mb-2">🔴 ライブオークション</h2>
        {!isEnded ? (
          <div className="flex items-center gap-2 text-red-600 font-semibold">
            <span className="w-3 h-3 bg-red-600 rounded-full animate-pulse"></span>
            <span>残り時間: {timeRemaining}</span>
          </div>
        ) : (
          <div className="text-gray-600 font-semibold">
            オークション終了
          </div>
        )}
      </div>

      <div className="bg-green-50 border-2 border-green-500 rounded-lg p-4 mb-6">
        <div className="text-sm text-gray-600 mb-1">現在の価格</div>
        <div className="text-4xl font-bold text-green-700">
          ¥{currentPrice.toLocaleString()}
        </div>
        {highestBidder && (
          <div className="text-sm text-gray-600 mt-2">
            最高入札者: {highestBidder}
          </div>
        )}
      </div>

      {!isEnded && (
        <div className="mb-6">
          <label className="block text-sm font-medium text-gray-700 mb-2">
            入札金額
          </label>
          <div className="flex gap-2 mb-3">
            <input
              type="number"
              value={bidAmount}
              onChange={(e) => setBidAmount(Number(e.target.value))}
              className="flex-1 px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-green-500"
              min={currentPrice + 100}
              step={1000}
            />
            <button
              onClick={placeBid}
              disabled={bidAmount <= currentPrice}
              className="px-6 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition disabled:bg-gray-400 disabled:cursor-not-allowed"
            >
              入札する
            </button>
          </div>

          <div className="flex gap-2">
            <button
              onClick={() => quickBid(1000)}
              className="flex-1 px-3 py-1 bg-gray-100 text-gray-700 rounded hover:bg-gray-200 transition text-sm"
            >
              +¥1,000
            </button>
            <button
              onClick={() => quickBid(5000)}
              className="flex-1 px-3 py-1 bg-gray-100 text-gray-700 rounded hover:bg-gray-200 transition text-sm"
            >
              +¥5,000
            </button>
            <button
              onClick={() => quickBid(10000)}
              className="flex-1 px-3 py-1 bg-gray-100 text-gray-700 rounded hover:bg-gray-200 transition text-sm"
            >
              +¥10,000
            </button>
          </div>
        </div>
      )}

      <div className="border-t pt-4">
        <h3 className="font-semibold mb-3">入札履歴</h3>
        <div className="space-y-2 max-h-60 overflow-y-auto">
          {bids.length === 0 ? (
            <p className="text-gray-500 text-sm">まだ入札がありません</p>
          ) : (
            bids.map((bid) => (
              <div
                key={bid.id}
                className="flex justify-between items-center p-2 bg-gray-50 rounded"
              >
                <div>
                  <div className="font-medium">{bid.username}</div>
                  <div className="text-sm text-gray-500">
                    {new Date(bid.timestamp).toLocaleTimeString()}
                  </div>
                </div>
                <div className="text-lg font-bold text-green-600">
                  ¥{bid.amount.toLocaleString()}
                </div>
              </div>
            ))
          )}
        </div>
      </div>

      {isEnded && highestBidder && (
        <div className="mt-4 bg-yellow-50 border border-yellow-200 rounded-lg p-4">
          <div className="text-center">
            <div className="text-xl font-bold text-yellow-800 mb-2">
              🎉 落札されました！
            </div>
            <div className="text-gray-700">
              落札者: <span className="font-semibold">{highestBidder}</span>
            </div>
            <div className="text-2xl font-bold text-green-600 mt-2">
              ¥{currentPrice.toLocaleString()}
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

import React, { useEffect, useRef, useState } from 'react';

interface ARTryOnProps {
  productImage: string;
  productName: string;
  category: string;
}

export const ARTryOn: React.FC<ARTryOnProps> = ({
  productImage,
  productName: _productName,
  category,
}) => {
  const videoRef = useRef<HTMLVideoElement>(null);
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const [isARActive, setIsARActive] = useState(false);
  const [supported, setSupported] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const animationFrameRef = useRef<number | null>(null);

  useEffect(() => {
    // Check if WebXR or camera API is supported
    if (typeof navigator !== 'undefined' && navigator.mediaDevices && 'getUserMedia' in navigator.mediaDevices) {
      setSupported(true);
    }
  }, []);

  const startAR = async () => {
    try {
      setError(null);
      const stream = await navigator.mediaDevices.getUserMedia({
        video: { facingMode: 'user', width: { ideal: 1280 }, height: { ideal: 720 } },
      });

      if (videoRef.current) {
        videoRef.current.srcObject = stream;
        await videoRef.current.play();
        setIsARActive(true);
      }
    } catch (err) {
      console.error('Error accessing camera:', err);
      setError('カメラへのアクセスが拒否されました。ブラウザの設定を確認してください。');
    }
  };

  const stopAR = () => {
    if (animationFrameRef.current) {
      cancelAnimationFrame(animationFrameRef.current);
      animationFrameRef.current = null;
    }

    if (videoRef.current && videoRef.current.srcObject) {
      const stream = videoRef.current.srcObject as MediaStream;
      stream.getTracks().forEach((track) => track.stop());
      videoRef.current.srcObject = null;
    }

    setIsARActive(false);
  };

  // Start AR rendering when video is ready
  useEffect(() => {
    if (!isARActive) return;

    const video = videoRef.current;
    const canvas = canvasRef.current;

    if (!video || !canvas) return;

    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    // Preload product image
    const productImg = new Image();
    productImg.src = productImage;
    productImg.crossOrigin = 'anonymous';

    const renderFrame = () => {
      if (!video || !canvas || !ctx) return;

      // Set canvas size to match video
      if (video.videoWidth > 0 && video.videoHeight > 0) {
        canvas.width = video.videoWidth;
        canvas.height = video.videoHeight;

        // Draw video frame
        ctx.drawImage(video, 0, 0, canvas.width, canvas.height);

        // AR overlay logic
        if (productImg.complete && productImg.naturalWidth > 0) {
          if (category === 'fashion' || category === 'clothing') {
            // Draw clothing overlay
            const overlayWidth = canvas.width * 0.6;
            const overlayHeight = (overlayWidth * productImg.height) / productImg.width;
            const x = (canvas.width - overlayWidth) / 2;
            const y = canvas.height * 0.3;

            ctx.globalAlpha = 0.8;
            ctx.drawImage(productImg, x, y, overlayWidth, overlayHeight);
            ctx.globalAlpha = 1.0;
          } else {
            // Default: center overlay for other categories
            const overlayWidth = canvas.width * 0.5;
            const overlayHeight = (overlayWidth * productImg.height) / productImg.width;
            const x = (canvas.width - overlayWidth) / 2;
            const y = (canvas.height - overlayHeight) / 2;

            ctx.globalAlpha = 0.85;
            ctx.drawImage(productImg, x, y, overlayWidth, overlayHeight);
            ctx.globalAlpha = 1.0;
          }
        }
      }

      animationFrameRef.current = requestAnimationFrame(renderFrame);
    };

    // Start rendering when video starts playing
    const handleVideoPlay = () => {
      renderFrame();
    };

    video.addEventListener('playing', handleVideoPlay);

    return () => {
      video.removeEventListener('playing', handleVideoPlay);
      if (animationFrameRef.current) {
        cancelAnimationFrame(animationFrameRef.current);
      }
    };
  }, [isARActive, productImage, category]);

  const capturePhoto = () => {
    const canvas = canvasRef.current;
    if (!canvas) return;

    const link = document.createElement('a');
    link.download = `ar-tryOn-${Date.now()}.png`;
    link.href = canvas.toDataURL('image/png');
    link.click();
  };

  if (!supported) {
    return (
      <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
        <p className="text-yellow-800">
          AR試着機能はこのブラウザではサポートされていません。
        </p>
      </div>
    );
  }

  return (
    <div className="bg-white rounded-lg shadow-lg p-6">
      <h3 className="text-xl font-bold mb-4 flex items-center gap-2">
        <span>📸</span>
        <span>AR試着機能</span>
        {isARActive && (
          <span className="ml-auto text-sm font-normal text-green-600 flex items-center gap-1">
            <span className="inline-block w-2 h-2 bg-green-600 rounded-full animate-pulse"></span>
            カメラ起動中
          </span>
        )}
      </h3>

      {error && (
        <div className="bg-red-50 border border-red-200 rounded-lg p-4 mb-4">
          <p className="text-red-800">{error}</p>
          <p className="text-sm text-red-600 mt-2">
            ブラウザのアドレスバーのカメラアイコンをクリックして、カメラのアクセスを許可してください。
          </p>
        </div>
      )}

      <div className="relative bg-gray-900 rounded-lg overflow-hidden mb-4" style={{ minHeight: '400px' }}>
        {!isARActive ? (
          <div className="aspect-video flex flex-col items-center justify-center gap-4 p-8">
            <div className="text-6xl mb-4">👕</div>
            <button
              onClick={startAR}
              className="bg-green-600 text-white px-8 py-4 rounded-lg hover:bg-green-700 transition text-lg font-semibold shadow-lg hover:shadow-xl transform hover:scale-105"
            >
              📸 AR試着を開始
            </button>
            <p className="text-gray-400 text-sm text-center max-w-md">
              カメラを使って商品を試着できます。<br />
              「AR試着を開始」をクリックしてカメラへのアクセスを許可してください。
            </p>
          </div>
        ) : (
          <>
            <video
              ref={videoRef}
              className="hidden"
              playsInline
              muted
              autoPlay
            />
            <canvas
              ref={canvasRef}
              className="w-full h-auto"
              style={{ maxHeight: '600px' }}
            />

            <div className="absolute bottom-4 left-1/2 transform -translate-x-1/2 flex gap-3">
              <button
                onClick={capturePhoto}
                className="bg-white text-gray-800 px-6 py-3 rounded-lg shadow-lg hover:bg-gray-100 transition font-medium flex items-center gap-2"
              >
                📷 <span>写真を撮る</span>
              </button>
              <button
                onClick={stopAR}
                className="bg-red-600 text-white px-6 py-3 rounded-lg shadow-lg hover:bg-red-700 transition font-medium flex items-center gap-2"
              >
                ✖️ <span>終了</span>
              </button>
            </div>
          </>
        )}
      </div>

      <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
        <p className="font-medium text-blue-900 mb-2">💡 使い方</p>
        <ul className="text-sm text-blue-800 space-y-1">
          <li>• 「AR試着を開始」をクリックしてカメラを起動</li>
          <li>• カメラに顔や体を映すと、商品画像が重ねて表示されます</li>
          <li>• 「写真を撮る」で試着画像を保存できます</li>
          <li>• ファッション系の商品は体の中央に、その他は画面中央に表示されます</li>
        </ul>
      </div>
    </div>
  );
};

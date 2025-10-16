import React, { useEffect, useRef, useState } from 'react';

interface LiveStreamProps {
  streamId: string;
  productId: string;
  isHost: boolean;
}

export const LiveStream: React.FC<LiveStreamProps> = ({
  streamId,
  productId: _productId,
  isHost,
}) => {
  const localVideoRef = useRef<HTMLVideoElement>(null);
  const remoteVideoRef = useRef<HTMLVideoElement>(null);
  const peerConnectionRef = useRef<RTCPeerConnection | null>(null);
  const [isLive, setIsLive] = useState(false);
  const [viewerCount, setViewerCount] = useState(0);
  const [comments, setComments] = useState<{ user: string; message: string }[]>([]);
  const [commentInput, setCommentInput] = useState('');

  useEffect(() => {
    const ws = new WebSocket(`ws://localhost:8080/v1/ws/live/${streamId}`);

    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);

      switch (data.type) {
        case 'viewer_count':
          setViewerCount(data.count);
          break;
        case 'comment':
          setComments((prev) => [...prev, { user: data.user, message: data.message }]);
          break;
        case 'offer':
          handleOffer(data.offer);
          break;
        case 'answer':
          handleAnswer(data.answer);
          break;
        case 'ice_candidate':
          handleIceCandidate(data.candidate);
          break;
      }
    };

    return () => {
      ws.close();
      if (peerConnectionRef.current) {
        peerConnectionRef.current.close();
      }
    };
  }, [streamId]);

  const startStream = async () => {
    try {
      const stream = await navigator.mediaDevices.getUserMedia({
        video: true,
        audio: true,
      });

      if (localVideoRef.current) {
        localVideoRef.current.srcObject = stream;
      }

      // Setup WebRTC peer connection
      const peerConnection = new RTCPeerConnection({
        iceServers: [
          { urls: 'stun:stun.l.google.com:19302' },
        ],
      });

      stream.getTracks().forEach((track) => {
        peerConnection.addTrack(track, stream);
      });

      peerConnectionRef.current = peerConnection;
      setIsLive(true);
    } catch (error) {
      console.error('Error starting stream:', error);
    }
  };

  const stopStream = () => {
    if (localVideoRef.current && localVideoRef.current.srcObject) {
      const stream = localVideoRef.current.srcObject as MediaStream;
      stream.getTracks().forEach((track) => track.stop());
      localVideoRef.current.srcObject = null;
    }

    if (peerConnectionRef.current) {
      peerConnectionRef.current.close();
      peerConnectionRef.current = null;
    }

    setIsLive(false);
  };

  const handleOffer = async (offer: RTCSessionDescriptionInit) => {
    if (!peerConnectionRef.current) return;
    await peerConnectionRef.current.setRemoteDescription(new RTCSessionDescription(offer));
    const answer = await peerConnectionRef.current.createAnswer();
    await peerConnectionRef.current.setLocalDescription(answer);
    // Send answer via WebSocket
  };

  const handleAnswer = async (answer: RTCSessionDescriptionInit) => {
    if (!peerConnectionRef.current) return;
    await peerConnectionRef.current.setRemoteDescription(new RTCSessionDescription(answer));
  };

  const handleIceCandidate = async (candidate: RTCIceCandidateInit) => {
    if (!peerConnectionRef.current) return;
    await peerConnectionRef.current.addIceCandidate(new RTCIceCandidate(candidate));
  };

  const sendComment = () => {
    if (commentInput.trim()) {
      // Send via WebSocket
      setCommentInput('');
    }
  };

  return (
    <div className="bg-gray-900 rounded-lg overflow-hidden">
      <div className="relative">
        {isHost ? (
          <video
            ref={localVideoRef}
            autoPlay
            muted
            playsInline
            className="w-full aspect-video bg-black"
          />
        ) : (
          <video
            ref={remoteVideoRef}
            autoPlay
            playsInline
            className="w-full aspect-video bg-black"
          />
        )}

        {isLive && (
          <div className="absolute top-4 left-4 flex items-center gap-2">
            <div className="bg-red-600 text-white px-3 py-1 rounded-full flex items-center gap-2">
              <span className="w-2 h-2 bg-white rounded-full animate-pulse"></span>
              <span className="font-semibold">LIVE</span>
            </div>
            <div className="bg-black bg-opacity-50 text-white px-3 py-1 rounded-full">
              👁️ {viewerCount}
            </div>
          </div>
        )}

        {isHost && !isLive && (
          <div className="absolute inset-0 flex items-center justify-center bg-black bg-opacity-50">
            <button
              onClick={startStream}
              className="bg-red-600 text-white px-8 py-4 rounded-lg text-xl font-bold hover:bg-red-700 transition"
            >
              🎥 配信を開始
            </button>
          </div>
        )}

        {isHost && isLive && (
          <div className="absolute bottom-4 left-1/2 transform -translate-x-1/2">
            <button
              onClick={stopStream}
              className="bg-red-600 text-white px-6 py-2 rounded-lg hover:bg-red-700 transition"
            >
              配信を終了
            </button>
          </div>
        )}
      </div>

      <div className="bg-gray-800 p-4">
        <div className="h-64 overflow-y-auto mb-4 space-y-2">
          {comments.map((comment, index) => (
            <div key={index} className="bg-gray-700 rounded p-2">
              <span className="text-green-400 font-semibold">{comment.user}: </span>
              <span className="text-white">{comment.message}</span>
            </div>
          ))}
        </div>

        <div className="flex gap-2">
          <input
            type="text"
            value={commentInput}
            onChange={(e) => setCommentInput(e.target.value)}
            onKeyPress={(e) => e.key === 'Enter' && sendComment()}
            placeholder="コメントを入力..."
            className="flex-1 px-4 py-2 bg-gray-700 text-white rounded-lg focus:outline-none focus:ring-2 focus:ring-green-500"
          />
          <button
            onClick={sendComment}
            className="px-6 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition"
          >
            送信
          </button>
        </div>
      </div>
    </div>
  );
};

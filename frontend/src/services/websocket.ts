import { useEffect, useRef, useState } from 'react';

export interface Message {
  id: string;
  conversation_id: string;
  sender_id: string;
  content: string;
  created_at: string;
  read: boolean;
}

export interface WebSocketMessage {
  type: 'message' | 'typing' | 'read' | 'online' | 'offline';
  data: any;
}

export const useWebSocket = (conversationId: string | null, token: string | null) => {
  const [messages, setMessages] = useState<Message[]>([]);
  const [connected, setConnected] = useState(false);
  const [typing, setTyping] = useState<string[]>([]);
  const wsRef = useRef<WebSocket | null>(null);

  useEffect(() => {
    if (!conversationId || !token) return;

    const ws = new WebSocket(
      `ws://localhost:8080/v1/ws/conversations/${conversationId}?token=${token}`
    );

    ws.onopen = () => {
      console.log('WebSocket connected');
      setConnected(true);
    };

    ws.onmessage = (event) => {
      const wsMessage: WebSocketMessage = JSON.parse(event.data);

      switch (wsMessage.type) {
        case 'message':
          setMessages((prev) => [...prev, wsMessage.data as Message]);
          break;
        case 'typing':
          setTyping((prev) => {
            if (wsMessage.data.typing) {
              return [...prev, wsMessage.data.user_id];
            } else {
              return prev.filter((id) => id !== wsMessage.data.user_id);
            }
          });
          break;
        case 'read':
          setMessages((prev) =>
            prev.map((msg) =>
              msg.id === wsMessage.data.message_id ? { ...msg, read: true } : msg
            )
          );
          break;
        default:
          console.log('Unknown message type:', wsMessage.type);
      }
    };

    ws.onerror = (error) => {
      console.error('WebSocket error:', error);
    };

    ws.onclose = () => {
      console.log('WebSocket disconnected');
      setConnected(false);
    };

    wsRef.current = ws;

    return () => {
      ws.close();
    };
  }, [conversationId, token]);

  const sendMessage = (content: string) => {
    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      wsRef.current.send(
        JSON.stringify({
          type: 'message',
          data: { content },
        })
      );
    }
  };

  const sendTyping = (isTyping: boolean) => {
    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      wsRef.current.send(
        JSON.stringify({
          type: 'typing',
          data: { typing: isTyping },
        })
      );
    }
  };

  const markAsRead = (messageId: string) => {
    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      wsRef.current.send(
        JSON.stringify({
          type: 'read',
          data: { message_id: messageId },
        })
      );
    }
  };

  return {
    messages,
    connected,
    typing,
    sendMessage,
    sendTyping,
    markAsRead,
  };
};

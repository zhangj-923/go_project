import { ref } from 'vue';
import type { WsMessage } from '../types/game';

const wsUrl = `ws://${window.location.hostname}:18189/ws`;
let socket: WebSocket | null = null;
const isConnected = ref(false);

type MsgHandler = (payload: any) => void;
const handlers = new Map<string, MsgHandler[]>();

export function useWebSocket() {
  const connect = () => {
    if (socket && (socket.readyState === WebSocket.OPEN || socket.readyState === WebSocket.CONNECTING)) {
      return;
    }

    socket = new WebSocket(wsUrl);

    socket.onopen = () => {
      console.log('WS Connected');
      isConnected.value = true;
    };

    socket.onmessage = (event) => {
      try {
        const msg: WsMessage = JSON.parse(event.data);
        const typeHandlers = handlers.get(msg.type);
        if (typeHandlers) {
          typeHandlers.forEach(h => h(msg.payload));
        }
      } catch (e) {
        console.error('Invalid message', e);
      }
    };

    socket.onclose = () => {
      console.log('WS Disconnected');
      isConnected.value = false;
      setTimeout(connect, 2000); // Reconnect
    };
  };

  const onMessage = (type: string, handler: MsgHandler) => {
    if (!handlers.has(type)) {
      handlers.set(type, []);
    }
    handlers.get(type)!.push(handler);
  };

  const removeHandler = (type: string, handler: MsgHandler) => {
    const typeHandlers = handlers.get(type);
    if (typeHandlers) {
      handlers.set(type, typeHandlers.filter(h => h !== handler));
    }
  };

  const send = (type: string, payload?: any) => {
    if (socket && socket.readyState === WebSocket.OPEN) {
      socket.send(JSON.stringify({ type, payload: payload || {} }));
    } else {
      console.error('Socket not connected');
    }
  };

  return {
    isConnected,
    connect,
    onMessage,
    removeHandler,
    send
  };
}

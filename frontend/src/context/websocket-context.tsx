import React, { createContext, useContext } from 'react';
import useWebSocket, { ReadyState } from 'react-use-websocket';

interface WebSocketContextType {
    readyState: ReadyState;
    lastMessage: MessageEvent<any> | null;
    isConnected: boolean;
    sendJsonMessage: (jsonMessage: any) => void;
}

const WebSocketContext = createContext<WebSocketContextType | null>(null);

export const WebSocketProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    // In dev, usage with Vite proxy: ws://localhost:5173/ws -> http://localhost:8082/ws
    // The proxy configuration in vite.config.ts handles the /ws path
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const socketUrl = `${protocol}//${window.location.host}/ws`;
    console.log(`[WebSocket] Initializing connection to ${socketUrl} (Secure: ${window.location.protocol === 'https:'})`);

    const {
        sendJsonMessage,
        lastMessage,
        readyState,
    } = useWebSocket(socketUrl, {
        shouldReconnect: () => true,
        reconnectAttempts: 20,
        reconnectInterval: 3000,
        share: true,
        onOpen: () => console.log('WebSocket Connected'),
        onClose: () => console.log('WebSocket Disconnected'),
        onError: (e) => console.error('WebSocket Error:', e),
    });

    const isConnected = readyState === ReadyState.OPEN;

    return (
        <WebSocketContext.Provider value={{ readyState, lastMessage, isConnected, sendJsonMessage }}>
            {children}
        </WebSocketContext.Provider>
    );
};

export const useWebSocketContext = () => {
    const context = useContext(WebSocketContext);
    if (!context) {
        throw new Error('useWebSocketContext must be used within a WebSocketProvider');
    }
    return context;
};

import { useState, useEffect } from 'react';
import { useWebSocketContext } from '@/context/websocket-context';
import { SourceState, ArbitrageOpportunity, BasisTradeOpportunity } from '@/types';

export function useMarketData() {
    const { lastMessage } = useWebSocketContext();

    // Use Ref for high-frequency price updates to avoid re-rendering entire app on every tick if possible
    // But for this UI we probably want to render the source list updates.
    // We'll trust React 18 batching for now.
    const [prices, setPrices] = useState<Record<string, SourceState>>({});
    const [opportunities, setOpportunities] = useState<ArbitrageOpportunity[]>([]);
    const [basisTrades, setBasisTrades] = useState<BasisTradeOpportunity[]>([]);
    const [spreads, setSpreads] = useState<any>(null); // Define proper type later

    useEffect(() => {
        if (lastMessage && lastMessage.data) {
            try {
                const data = JSON.parse(lastMessage.data);

                if (data.type === 'prices') {
                    // Handle bulk price update
                    const newPrices: Record<string, SourceState> = {};
                    const priceMap = data.prices["TONUSDT"] || {}; // Default to TONUSDT for now as per main.go

                    // Using a timestamp for all updates in this batch
                    const now = Date.now();

                    Object.entries(priceMap).forEach(([source, price]) => {
                        const p = Number(price);

                        // Calculate change relative to previous state if available
                        const previous = prices[source];
                        const previousPrice = previous ? previous.price : p;
                        const change = p - previousPrice;
                        const changePercent = previousPrice !== 0 ? (change / previousPrice) * 100 : 0;

                        newPrices[source] = {
                            price: p,
                            previousPrice: previousPrice,
                            change,
                            changePercent,
                            lastUpdate: now
                        };
                    });

                    setPrices(prev => ({
                        ...prev,
                        ...newPrices
                    }));
                } else if (data.type === 'price_update') {
                    // Update specific source
                    setPrices(prev => {
                        const previous = prev[data.source];
                        const previousPrice = previous ? previous.price : data.price;
                        const change = data.price - previousPrice;
                        // Avoid NaN on init
                        const changePercent = previousPrice !== 0 ? (change / previousPrice) * 100 : 0;

                        return {
                            ...prev,
                            [data.source]: {
                                price: data.price,
                                previousPrice: previousPrice,
                                change,
                                changePercent,
                                lastUpdate: Date.now()
                            }
                        };
                    });
                } else if (data.type === 'arbitrage') {
                    const opp = data.opportunity;
                    opp.id = opp.id || `${Date.now()}-${Math.random()}`; // Ensure ID

                    setOpportunities(prev => {
                        const newOpps = [opp, ...prev];
                        return newOpps.slice(0, 50); // Keep last 50
                    });
                } else if (data.type === 'basis_trade') {
                    const opp = data.opportunity;
                    opp.id = opp.id || `${Date.now()}-${Math.random()}`; // Ensure ID

                    setBasisTrades(prev => {
                        const newOpps = [opp, ...prev];
                        return newOpps.slice(0, 50); // Keep last 50
                    });
                } else if (data.type === 'spreads') {
                    setSpreads(data);
                }
            } catch (e) {
                console.error("Failed to parse websocket message", e);
            }
        }
    }, [lastMessage]);

    return {
        prices,
        opportunities,
        basisTrades,
        spreads
    };
}

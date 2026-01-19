export interface PriceData {
    symbol: string;
    source: string;
    price: number;
    timestamp: number;
}

export interface ArbitrageOpportunity {
    id?: string;
    symbol: string;
    buy_source: string;
    sell_source: string;
    buy_price: number;
    sell_price: number;
    profit_pct: number;
    timestamp: number;
}

export interface BasisTradeOpportunity {
    id?: string;
    symbol: string;
    dedust_price: number;
    short_source: string;
    short_price: number;
    profit_pct: number;
    timestamp: number;
}

export interface SourceState {
    price: number;
    previousPrice: number;
    change: number;
    changePercent: number;
    lastUpdate: number;
}

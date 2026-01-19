import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card"

import { SourceState } from "@/types"
import { cn } from "@/lib/utils"

interface SourcesPanelProps {
    prices: Record<string, SourceState>
}

// Map backend source keys to friendly names & colors (simplified mapping)
const sourceConfig: Record<string, { label: string, color: string }> = {
    'binance_futures': { label: 'Binance Futures', color: '#F0B90B' },
    'bybit_futures': { label: 'Bybit Futures', color: '#F7931A' },
    'hyperliquid_futures': { label: 'Hyperliquid', color: '#97FCE4' },
    'kraken_futures': { label: 'Kraken', color: '#5a5aff' },
    'okx_futures': { label: 'OKX Futures', color: '#1890ff' },
    'gate_futures': { label: 'Gate.io', color: '#6c5ce7' },
    'paradex_futures': { label: 'Paradex', color: '#ff6b6b' },
    'binance_spot': { label: 'Binance Spot', color: '#F0B90B' },
    'bybit_spot': { label: 'Bybit Spot', color: '#F7931A' },
    'pyth': { label: 'Pyth Oracle', color: '#00ff88' },
}

export function SourcesPanel({ prices }: SourcesPanelProps) {
    // Sort sources to keep stable order

    // Let's show active ones primarily, or just list all from config

    return (
        <Card className="h-full bg-zinc-900/50 border-zinc-800">
            <CardHeader className="py-3 px-4">
                <CardTitle>Data Sources</CardTitle>
            </CardHeader>
            <CardContent className="p-0 overflow-y-auto max-h-[calc(100%-50px)]">
                <div className="divide-y divide-zinc-800/50">
                    {Object.entries(sourceConfig).map(([key, config]) => {
                        const data = prices[key];
                        const price = data ? data.price : 0;
                        const change = data ? data.changePercent : 0;
                        const isUp = change >= 0;

                        return (
                            <div key={key} className="flex items-center justify-between py-2 px-4 hover:bg-zinc-800/30 transition-colors cursor-pointer group">
                                <div className="flex items-center gap-2">
                                    <div className="w-2 h-2 rounded-full" style={{ backgroundColor: config.color }}></div>
                                    <span className="text-xs text-zinc-400 group-hover:text-zinc-200 transition-colors uppercase font-medium">{config.label}</span>
                                </div>
                                <div className="text-right">
                                    <div className="text-xs font-mono font-medium text-zinc-200">
                                        ${price > 0 ? price.toFixed(price > 100 ? 2 : 4) : '---'}
                                    </div>
                                    <div className={cn("text-[10px] font-mono", isUp ? "text-primary" : "text-error")}>
                                        {data ? `${isUp ? '+' : ''}${change.toFixed(3)}%` : '--'}
                                    </div>
                                </div>
                            </div>
                        )
                    })}
                </div>
            </CardContent>
        </Card>
    )
}


import { cn } from "@/lib/utils"

interface SpreadsMatrixProps {
    data: any // using any for complex nested object for now, define properly later
}

// Short names for headers
const shortNames: Record<string, string> = {
    'binance_futures': 'BIN-F',
    'bybit_futures': 'BYB-F',
    'hyperliquid_futures': 'HYP',
    'kraken_futures': 'KRK',
    'okx_futures': 'OKX',
    'gate_futures': 'GATE',
    'paradex_futures': 'PRD',
    'lighter_futures': 'LGHT',
    'extended_futures': 'EXT',
    'vest_futures': 'VEST',
}

export function SpreadsMatrix({ data }: SpreadsMatrixProps) {
    if (!data || !data.spreads) {
        return (
            <div className="h-full flex items-center justify-center text-zinc-500 text-xs italic">
                Waiting for spread data...
            </div>
        )
    }

    const sources = Object.keys(data.spreads);
    // Filter to only show futures or configured sources
    const displaySources = sources.filter(s => shortNames[s]);
    const matrix = data.spreads;

    return (
        <div className="h-full w-full overflow-hidden flex flex-col">
            <div className="flex-1 overflow-auto scrollbar-hide">
                <div className="w-full flex flex-col min-w-[600px]">
                    {/* Header Row */}
                    <div className="flex w-full">
                        <div className="w-16 h-8 flex-shrink-0"></div> {/* Corner */}
                        {displaySources.map(source => (
                            <div key={source} className="flex-1 h-8 flex items-center justify-center text-[10px] font-bold text-zinc-500 uppercase bg-zinc-900/50 border-b border-zinc-800">
                                {shortNames[source] || source.substring(0, 3)}
                            </div>
                        ))}
                    </div>

                    {/* Rows */}
                    {displaySources.map(buySource => (
                        <div key={buySource} className="flex w-full flex-1 min-h-[30px]">
                            {/* Row Header */}
                            <div className="w-16 flex items-center justify-center text-[10px] font-bold text-zinc-500 uppercase flex-shrink-0 bg-zinc-900/50 border-r border-zinc-800">
                                {shortNames[buySource]}
                            </div>

                            {/* Cells */}
                            {displaySources.map(sellSource => {
                                if (buySource === sellSource) {
                                    return <div key={sellSource} className="flex-1 bg-zinc-900/30 flex items-center justify-center text-zinc-700">-</div>;
                                }

                                const val = matrix[buySource]?.[sellSource] || 0;
                                // Formatting 
                                let colorClass = "text-zinc-500";
                                let bgClass = "";

                                if (val > 0.05) {
                                    colorClass = "text-primary font-bold";
                                    bgClass = "bg-primary/10";
                                } else if (val < -0.05) {
                                    colorClass = "text-error";
                                    bgClass = "bg-error/5";
                                }

                                return (
                                    <div key={sellSource} className={cn("flex-1 flex items-center justify-center text-[11px] font-mono border border-zinc-900/20 hover:bg-zinc-800/50 transition-colors", bgClass)}>
                                        <span className={colorClass}>{val.toFixed(2)}%</span>
                                    </div>
                                )
                            })}
                        </div>
                    ))}
                </div>
            </div>
        </div>

    )
}

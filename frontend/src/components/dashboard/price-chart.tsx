import { useEffect, useRef } from 'react';
import { createChart, ColorType, IChartApi, ISeriesApi, LineSeries, LineStyle } from 'lightweight-charts';
import { useMarketData } from '@/hooks/use-market-data';


interface PriceChartProps {
    symbol?: string;
}

const sourceColors: Record<string, string> = {
    'binance_futures': '#F0B90B',
    'bybit_futures': '#F7931A',
    'hyperliquid_futures': '#97FCE4',
    'kraken_futures': '#5a5aff',
    'okx_futures': '#1890ff',
    'gate_futures': '#6c5ce7',
    'paradex_futures': '#ff6b6b',
    'binance_spot': '#F0B90B', // Dashed usually
};

export function PriceChart({ symbol = "BTCUSDT" }: PriceChartProps) {
    const chartContainerRef = useRef<HTMLDivElement>(null);
    const chartRef = useRef<IChartApi | null>(null);
    const seriesRef = useRef<Map<string, ISeriesApi<"Line">>>(new Map());

    // We need to access the LATEST prices for the chart. 
    // The `useMarketData` hook gives us current state, but for a chart we want to append.
    // Ideally the hook would provide a stream or we listen to the socket directly here, 
    // OR we just watch the `prices` object change. Watching `prices` object change is risky if it changes too fast.
    // But let's try watching it for now.
    const { prices } = useMarketData();

    useEffect(() => {
        if (!chartContainerRef.current) return;

        const chart = createChart(chartContainerRef.current, {
            layout: {
                background: { type: ColorType.Solid, color: 'transparent' },
                textColor: '#71717a',
                fontFamily: 'JetBrains Mono',
            },
            grid: {
                vertLines: { color: '#27272a' },
                horzLines: { color: '#27272a' },
            },
            width: chartContainerRef.current.clientWidth,
            height: chartContainerRef.current.clientHeight,
            timeScale: {
                timeVisible: true,
                secondsVisible: true,
                borderColor: '#27272a',
            },
            rightPriceScale: {
                borderColor: '#27272a',
            },
        });

        chartRef.current = chart;

        // Initialize series for known keys (or dynamic)
        Object.entries(sourceColors).forEach(([key, color]) => {
            const series = chart.addSeries(LineSeries, {
                color: color,
                lineWidth: 2,
                title: key.replace('_futures', '').toUpperCase(),
                lineStyle: key.includes('spot') ? LineStyle.Dashed : LineStyle.Solid,
                lastValueVisible: false,
                priceLineVisible: false,
            });
            seriesRef.current.set(key, series);
        });

        const handleResize = () => {
            if (chartContainerRef.current) {
                chart.applyOptions({ width: chartContainerRef.current.clientWidth });
            }
        };

        window.addEventListener('resize', handleResize);

        return () => {
            window.removeEventListener('resize', handleResize);
            chart.remove();
        };
    }, []);

    // Update chart data when prices change
    useEffect(() => {
        if (!prices || !chartRef.current) return;

        const now = Date.now() / 1000; // Time in seconds

        Object.entries(prices).forEach(([source, data]) => {
            const series = seriesRef.current.get(source);
            if (series) {
                // Determine timestamp logic: use data.lastUpdate or local time?
                // data.lastUpdate might not be set for all, so use local time if needed
                // But generally data.price is what we want.
                series.update({
                    time: now as any, // casting for lightweight-charts strict types
                    value: data.price
                });
            }
        });

    }, [prices]); // This depends on `prices` object identity changing. `useMarketData` should update it immutably.

    return (
        <div className="relative w-full h-full">
            <div ref={chartContainerRef} className="w-full h-full" />
            <div className="absolute top-2 right-2 text-xs text-zinc-500 pointer-events-none">
                {symbol} LIVE
            </div>
        </div>
    );
}

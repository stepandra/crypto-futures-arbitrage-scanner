import { useState } from 'react'
import { MainLayout } from '@/components/layout/main-layout'
import { useWebSocketContext } from '@/context/websocket-context'
import { useMarketData } from '@/hooks/use-market-data'
import { SourcesPanel } from './sources-panel'
import { OpportunitiesTable } from './opportunities-table'
import { SpreadsMatrix } from './spreads-matrix'
import { PriceChart } from './price-chart'
import { BasisTradeTable } from './basis-trade-table'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'

export function Dashboard() {
    const { isConnected } = useWebSocketContext()
    const { prices, opportunities, spreads, basisTrades } = useMarketData()
    const [activeView, setActiveView] = useState('dashboard')

    const renderContent = () => {
        switch (activeView) {
            case 'charts':
                return (
                    <div className="h-full">
                        <Card className="border-zinc-800 bg-zinc-900/50 flex flex-col h-full overflow-hidden">
                            <CardHeader className="py-2 px-4 border-zinc-800 h-10 min-h-0 flex-shrink-0"><CardTitle>Live Chart - TON/USD</CardTitle></CardHeader>
                            <CardContent className="flex-1 p-0 relative">
                                <PriceChart />
                            </CardContent>
                        </Card>
                    </div>
                )
            case 'spreads':
                return (
                    <div className="h-full">
                        <Card className="border-zinc-800 bg-zinc-900/50 flex flex-col h-full overflow-hidden">
                            <CardHeader className="py-2 px-4 border-zinc-800 h-10 min-h-0 flex-shrink-0"><CardTitle>Spreads Analysis</CardTitle></CardHeader>
                            <CardContent className="flex-1 p-2 overflow-auto">
                                <SpreadsMatrix data={spreads} />
                            </CardContent>
                        </Card>
                    </div>
                )
            case 'basis-trade':
                return (
                    <div className="h-full">
                        <BasisTradeTable data={basisTrades} />
                    </div>
                )
            case 'dashboard':
            default:
                return (
                    <div className="grid grid-cols-12 gap-4 h-full">
                        {/* Left Sidebar: Sources */}
                        <div className="col-span-12 md:col-span-3 lg:col-span-2 h-full flex flex-col gap-4">
                            <SourcesPanel prices={prices} />
                        </div>


                        {/* Center: Main Workspace */}
                        <div className="col-span-12 md:col-span-9 lg:col-span-10 grid grid-rows-12 gap-4 h-full">

                            {/* Top Row: Spreads */}
                            <div className="row-span-4 grid grid-cols-1 gap-4">
                                <Card className="border-zinc-800 bg-zinc-900/50 flex flex-col overflow-hidden">
                                    <CardHeader className="py-2 px-4 border-zinc-800 h-10 min-h-0 flex-shrink-0"><CardTitle>Spreads Analysis</CardTitle></CardHeader>
                                    <CardContent className="flex-1 p-2 overflow-auto">
                                        <SpreadsMatrix data={spreads} />
                                    </CardContent>
                                </Card>
                            </div>

                            {/* Bottom Row: Opportunities Table (Dominant) */}
                            <div className="row-span-8 h-full overflow-hidden">
                                <OpportunitiesTable data={opportunities} />
                            </div>
                        </div>
                    </div>
                )
        }
    }

    return (
        <MainLayout isConnected={isConnected} activeView={activeView} onNavigate={setActiveView}>
            {renderContent()}
        </MainLayout>
    )
}

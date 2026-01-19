import { LayoutDashboard, Settings, Activity, LineChart, Layers, Radio, ArrowLeftRight } from "lucide-react"
import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"

interface SidebarProps extends React.HTMLAttributes<HTMLDivElement> {
    className?: string
    activeView?: string
    onNavigate?: (view: string) => void
}

export function Sidebar({ className, activeView = 'dashboard', onNavigate }: SidebarProps) {
    return (
        <div className={cn("pb-12 border-r border-border bg-background flex flex-col h-full", className)}>
            <div className="space-y-4 py-4">
                <div className="px-3 py-2">
                    <h2 className="mb-2 px-4 text-lg font-semibold tracking-tight text-primary flex items-center gap-2">
                        <Activity className="h-5 w-5" />
                        ALGO TON
                    </h2>
                    <div className="space-y-1">
                        <Button
                            variant={activeView === 'dashboard' ? "secondary" : "ghost"}
                            className="w-full justify-start gap-2"
                            onClick={() => onNavigate?.('dashboard')}
                        >
                            <LayoutDashboard className="h-4 w-4" />
                            Dashboard
                        </Button>
                        <Button
                            variant={activeView === 'charts' ? "secondary" : "ghost"}
                            className="w-full justify-start gap-2"
                            onClick={() => onNavigate?.('charts')}
                        >
                            <LineChart className="h-4 w-4" />
                            Charts
                        </Button>
                        <Button
                            variant={activeView === 'spreads' ? "secondary" : "ghost"}
                            className="w-full justify-start gap-2"
                            onClick={() => onNavigate?.('spreads')}
                        >
                            <Layers className="h-4 w-4" />
                            Spreads
                        </Button>
                        <Button
                            variant={activeView === 'basis-trade' ? "secondary" : "ghost"}
                            className="w-full justify-start gap-2"
                            onClick={() => onNavigate?.('basis-trade')}
                        >
                            <ArrowLeftRight className="h-4 w-4" />
                            Basis Trade
                        </Button>
                    </div>
                </div>
                <div className="px-3 py-2">
                    <h2 className="mb-2 px-4 text-xs font-semibold uppercase tracking-wider text-zinc-500">
                        Settings
                    </h2>
                    <div className="space-y-1">
                        <Button variant="ghost" className="w-full justify-start gap-2">
                            <Radio className="h-4 w-4" />
                            Data Sources
                        </Button>
                        <Button variant="ghost" className="w-full justify-start gap-2">
                            <Settings className="h-4 w-4" />
                            Configuration
                        </Button>
                    </div>
                </div>
            </div>
        </div>
    )
}

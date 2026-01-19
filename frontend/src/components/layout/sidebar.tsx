import { LayoutDashboard, LineChart, Layers, ArrowLeftRight } from "lucide-react"
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

            </div>
        </div>
    )
}

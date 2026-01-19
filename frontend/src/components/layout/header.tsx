import { Badge } from "@/components/ui/badge"
import { Wifi, WifiOff } from "lucide-react"

interface HeaderProps {
    isConnected: boolean;
}

export function Header({ isConnected }: HeaderProps) {
    return (
        <header className="h-14 border-b border-border bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 flex items-center justify-between px-6 sticky top-0 z-50">
            <div className="flex items-center gap-4">
                <div className="font-bold text-sm tracking-widest text-primary">
                    ALGO TON
                </div>
                <div className="h-4 w-[1px] bg-zinc-800"></div>
                <div className="flex items-center gap-2 text-xs text-zinc-500 font-mono">
                    <span>TON/USD</span>
                    <span className="text-primary">+1.24%</span>
                </div>
            </div>

            <div className="flex items-center gap-4">
                <Badge variant={isConnected ? "success" : "destructive"} className="gap-1.5 transition-all duration-300">
                    {isConnected ? <Wifi className="h-3 w-3" /> : <WifiOff className="h-3 w-3" />}
                    {isConnected ? "CONNECTED" : "DISCONNECTED"}
                </Badge>
            </div>
        </header>
    )
}

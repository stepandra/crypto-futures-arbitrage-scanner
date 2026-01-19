import { ReactNode } from "react"
import { Sidebar } from "./sidebar"
import { Header } from "./header"

interface MainLayoutProps {
    children: ReactNode
    isConnected?: boolean
    activeView?: string
    onNavigate?: (view: string) => void
}

export function MainLayout({ children, isConnected = false, activeView, onNavigate }: MainLayoutProps) {
    return (
        <div className="min-h-screen bg-background flex flex-col">
            <Header isConnected={isConnected} />
            <div className="flex-1 flex overflow-hidden h-[calc(100vh-3.5rem)]">
                <Sidebar className="w-64 hidden md:flex shrink-0" activeView={activeView} onNavigate={onNavigate} />
                <main className="flex-1 overflow-auto p-2 md:p-4 bg-zinc-950/50">
                    {children}
                </main>
            </div>
        </div>
    )
}

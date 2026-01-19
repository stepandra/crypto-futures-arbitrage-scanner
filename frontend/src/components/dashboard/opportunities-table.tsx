import { useState } from "react"
import {
    Table, TableBody, TableCell, TableHead, TableHeader, TableRow
} from "@/components/ui/table"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input" // Need to create this or use standard input
import { ArbitrageOpportunity } from "@/types"
import { Trash2, Filter } from "lucide-react"


interface OpportunitiesTableProps {
    data: ArbitrageOpportunity[]
}

export function OpportunitiesTable({ data }: OpportunitiesTableProps) {
    const [minProfit, setMinProfit] = useState<number>(0.05)

    // Filter logic
    const filteredData = data.filter(opp => opp.profit_pct >= minProfit)

    return (
        <div className="rounded-md border border-border bg-surface overflow-hidden flex flex-col h-full">
            <div className="p-2 border-b border-border flex items-center justify-between bg-zinc-900/50">
                <div className="flex items-center gap-4">
                    <h3 className="text-xs font-bold uppercase tracking-wider text-primary flex items-center gap-2">
                        Arbitrage Signals
                    </h3>
                    <div className="flex items-center gap-2">
                        <Filter className="h-3 w-3 text-zinc-500" />
                        <span className="text-[10px] text-zinc-500 uppercase">Min Profit %</span>
                        <Input
                            type="number"
                            step="0.01"
                            value={minProfit}
                            onChange={(e) => setMinProfit(parseFloat(e.target.value) || 0)}
                            className="bg-zinc-950 border-zinc-800 h-6 w-16 px-2 text-xs focus-visible:ring-primary/50"
                        />
                    </div>
                </div>
                <div className="flex items-center gap-2">
                    <Badge variant="outline" className="text-[10px] h-5 gap-1">
                        {filteredData.length} <span className="text-zinc-600">/ {data.length}</span>
                    </Badge>
                    <Button variant="ghost" size="icon" className="h-6 w-6 text-zinc-500 hover:text-error">
                        <Trash2 className="h-3 w-3" />
                    </Button>
                </div>
            </div>
            <div className="flex-1 overflow-auto">
                <Table>
                    <TableHeader className="bg-zinc-900/80 backdrop-blur">
                        <TableRow className="hover:bg-transparent border-zinc-800">
                            <TableHead className="w-[100px] text-zinc-400">Pair</TableHead>
                            <TableHead className="text-zinc-400">Profit</TableHead>
                            <TableHead className="text-zinc-400">Buy @</TableHead>
                            <TableHead className="text-zinc-400">Sell @</TableHead>
                            <TableHead className="text-right text-zinc-400">Time</TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        {filteredData.length === 0 ? (
                            <TableRow>
                                <TableCell colSpan={5} className="h-24 text-center text-zinc-500 italic">
                                    {data.length > 0 ? "No signals match filter" : "Waiting for signals..."}
                                </TableCell>
                            </TableRow>
                        ) : (
                            filteredData.map((opp) => (
                                <TableRow key={opp.id} className="hover:bg-zinc-800/40 border-zinc-800/50 group animate-fade-in text-[11px]">
                                    <TableCell className="font-bold text-zinc-200 py-1">
                                        <div className="flex items-center gap-2">
                                            {opp.symbol}
                                        </div>
                                    </TableCell>
                                    <TableCell className="py-1">
                                        <Badge variant={opp.profit_pct > 0.5 ? "destructive" : opp.profit_pct > 0.2 ? "warning" : "default"} className="font-mono text-[10px] h-5">
                                            {opp.profit_pct.toFixed(3)}%
                                        </Badge>
                                    </TableCell>
                                    <TableCell className="py-1">
                                        <div className="flex flex-col">
                                            <span className="text-[9px] text-zinc-500 uppercase">{opp.buy_source.replace('_futures', '')}</span>
                                            <span className="text-zinc-300 font-mono">${opp.buy_price.toFixed(4)}</span>
                                        </div>
                                    </TableCell>
                                    <TableCell className="py-1">
                                        <div className="flex flex-col">
                                            <span className="text-[9px] text-zinc-500 uppercase">{opp.sell_source.replace('_futures', '')}</span>
                                            <span className="text-zinc-300 font-mono">${opp.sell_price.toFixed(4)}</span>
                                        </div>
                                    </TableCell>
                                    <TableCell className="text-right text-zinc-500 text-[10px] py-1 font-mono">
                                        {new Date(opp.timestamp).toLocaleTimeString()}
                                    </TableCell>
                                </TableRow>
                            ))
                        )}
                    </TableBody>
                </Table>
            </div>
        </div>
    )
}

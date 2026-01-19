import { BasisTradeOpportunity } from "@/types"
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"

interface BasisTradeTableProps {
    data: BasisTradeOpportunity[]
}

export function BasisTradeTable({ data }: BasisTradeTableProps) {
    return (
        <Card className="border-1 border-white/5 bg-zinc-900/50 flex flex-col h-full overflow-hidden">
            <CardHeader className="py-2 px-4 border-b border-white/5 h-10 min-h-0 flex-shrink-0 flex flex-row items-center justify-between">
                <CardTitle className="text-sm font-medium text-zinc-400">Basis Trade Opportunities</CardTitle>
                <Badge variant="outline" className="bg-emerald-500/10 text-emerald-500 border-emerald-500/20">
                    Live
                </Badge>
            </CardHeader>
            <CardContent className="flex-1 p-0 overflow-auto">
                <Table>
                    <TableHeader className="bg-zinc-900/90 sticky top-0 z-10">
                        <TableRow className="border-white/5 hover:bg-transparent">
                            <TableHead className="text-xs font-semibold text-zinc-500 uppercase w-[100px]">Time</TableHead>
                            <TableHead className="text-xs font-semibold text-zinc-500 uppercase">Symbol</TableHead>
                            <TableHead className="text-xs font-semibold text-zinc-500 uppercase text-right text-emerald-500">BUY on DeDust DEX</TableHead>
                            <TableHead className="text-xs font-semibold text-zinc-500 uppercase text-right text-red-500">SHORT on Perps</TableHead>
                            <TableHead className="text-xs font-semibold text-zinc-500 uppercase text-right">Spread %</TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        {data.length === 0 ? (
                            <TableRow>
                                <TableCell colSpan={5} className="text-center text-zinc-500 py-8">
                                    No basis trade opportunities detected yet...
                                </TableCell>
                            </TableRow>
                        ) : (
                            data.map((opp, i) => (
                                <TableRow key={opp.id || i} className="border-white/5 hover:bg-white/5 transition-colors">
                                    <TableCell className="font-mono text-xs text-zinc-500">
                                        {new Date(opp.timestamp).toLocaleTimeString()}
                                    </TableCell>
                                    <TableCell className="font-medium text-zinc-300">
                                        {opp.symbol}
                                    </TableCell>
                                    <TableCell className="text-right font-mono text-emerald-400">
                                        {opp.dedust_price.toFixed(4)}
                                    </TableCell>
                                    <TableCell className="text-right font-mono">
                                        <span className="text-red-400 mr-2">{opp.short_price.toFixed(4)}</span>
                                        <Badge variant="secondary" className="px-1 py-0 h-5 text-[10px] bg-zinc-800 text-zinc-400 border-zinc-700">
                                            {opp.short_source}
                                        </Badge>
                                    </TableCell>
                                    <TableCell className="text-right font-mono font-bold text-emerald-500">
                                        +{opp.profit_pct.toFixed(2)}%
                                    </TableCell>
                                </TableRow>
                            ))
                        )}
                    </TableBody>
                </Table>
            </CardContent>
        </Card>
    )
}

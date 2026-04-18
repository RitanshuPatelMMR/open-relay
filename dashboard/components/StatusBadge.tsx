import { Badge } from "@/components/ui/badge"
import { cn } from "@/lib/utils"

export function StatusBadge({ status }: { status: string }) {
  const styles: Record<string, string> = {
    delivered: "bg-emerald-500/15 text-emerald-400 border-emerald-500/30 hover:bg-emerald-500/20",
    failed:    "bg-red-500/15 text-red-400 border-red-500/30 hover:bg-red-500/20",
    pending:   "bg-amber-500/15 text-amber-400 border-amber-500/30 hover:bg-amber-500/20",
  }
  return (
    <Badge
      variant="outline"
      className={cn("text-xs font-medium capitalize", styles[status] ?? "bg-slate-500/15 text-slate-400")}
    >
      {status}
    </Badge>
  )
}
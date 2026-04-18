"use client"
import { useEffect, useState } from "react"
import { useRouter } from "next/navigation"
import { useWebSocket } from "@/hooks/useWebSocket"
import { getEvents } from "@/lib/api"
import { StatusBadge } from "@/components/StatusBadge"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { Skeleton } from "@/components/ui/skeleton"
import { Badge } from "@/components/ui/badge"
import { formatDistanceToNow } from "date-fns"

export function EventFeed() {
  const [events, setEvents] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  const { messages, connected } = useWebSocket()
  const router = useRouter()

  useEffect(() => {
    getEvents().then((data) => {
      setEvents(Array.isArray(data) ? data : [])
      setLoading(false)
    })
  }, [])

  // refresh list on WS update
  useEffect(() => {
    if (messages.length > 0) {
      getEvents().then((data) => setEvents(Array.isArray(data) ? data : []))
    }
  }, [messages])

  if (loading) return (
    <div className="space-y-3">
      {[...Array(6)].map((_, i) => <Skeleton key={i} className="h-14 w-full bg-slate-800" />)}
    </div>
  )

  return (
    <div className="rounded-xl border border-slate-800 overflow-hidden">
      <Table>
        <TableHeader>
          <TableRow className="border-slate-800 hover:bg-transparent">
            <TableHead className="text-slate-400 font-medium">Event ID</TableHead>
            <TableHead className="text-slate-400 font-medium">Method</TableHead>
            <TableHead className="text-slate-400 font-medium">Status</TableHead>
            <TableHead className="text-slate-400 font-medium">Attempts</TableHead>
            <TableHead className="text-slate-400 font-medium">Received</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {events.length === 0 && (
            <TableRow>
              <TableCell colSpan={5} className="text-center text-slate-500 py-12">
                No events yet. Send a webhook to get started.
              </TableCell>
            </TableRow>
          )}
          {events.map((e) => (
            <TableRow
              key={e.id}
              onClick={() => router.push(`/events/${e.id}`)}
              className="border-slate-800 cursor-pointer hover:bg-slate-800/50 transition-colors"
            >
              <TableCell className="font-mono text-xs text-slate-300">{e.id.slice(0, 8)}…</TableCell>
              <TableCell>
                <Badge variant="outline" className="text-xs border-slate-700 text-slate-300">{e.method}</Badge>
              </TableCell>
              <TableCell><StatusBadge status={e.status} /></TableCell>
              <TableCell className="text-slate-400 text-sm">{e.attempt_count}</TableCell>
              <TableCell className="text-slate-500 text-xs">
                {formatDistanceToNow(new Date(e.created_at), { addSuffix: true })}
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  )
}
import { getEvent, replayEvent } from "@/lib/api"
import { StatusBadge } from "@/components/StatusBadge"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { ScrollArea } from "@/components/ui/scroll-area"
import { ReplayButton } from "@/components/ReplayButton"
import { ArrowLeft, Clock, Globe, Hash } from "lucide-react"
import Link from "next/link"
import { formatDistanceToNow } from "date-fns"

export default async function EventDetailPage({ params }: { params: { id: string } }) {
  const { event, delivery_attempts } = await getEvent(params.id)

  return (
    <div className="p-8 max-w-5xl">
      {/* Back */}
      <Link href="/" className="inline-flex items-center gap-2 text-slate-400 hover:text-slate-50 text-sm mb-6 transition-colors">
        <ArrowLeft className="w-4 h-4" /> Back to Events
      </Link>

      {/* Header */}
      <div className="flex items-start justify-between mb-8">
        <div>
          <div className="flex items-center gap-3 mb-2">
            <h1 className="text-xl font-bold font-mono text-slate-50">{event?.id}</h1>
            <StatusBadge status={event?.status} />
          </div>
          <div className="flex items-center gap-4 text-xs text-slate-500">
            <span className="flex items-center gap-1"><Clock className="w-3 h-3" />
              {event?.created_at && formatDistanceToNow(new Date(event.created_at), { addSuffix: true })}
            </span>
            <span className="flex items-center gap-1"><Globe className="w-3 h-3" />{event?.source_ip}</span>
            <span className="flex items-center gap-1"><Hash className="w-3 h-3" />{event?.attempt_count} attempt(s)</span>
          </div>
        </div>
        <ReplayButton eventId={event?.id} />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Payload */}
        <Card className="bg-slate-900 border-slate-800">
          <CardHeader className="pb-3">
            <CardTitle className="text-sm text-slate-300 font-medium">Payload</CardTitle>
          </CardHeader>
          <CardContent>
            <ScrollArea className="h-72">
              <pre className="text-xs text-emerald-400 font-mono leading-relaxed whitespace-pre-wrap break-all">
                {event?.payload
                  ? JSON.stringify(JSON.parse(event.payload), null, 2)
                  : "No payload"}
              </pre>
            </ScrollArea>
          </CardContent>
        </Card>

        {/* Headers */}
        <Card className="bg-slate-900 border-slate-800">
          <CardHeader className="pb-3">
            <CardTitle className="text-sm text-slate-300 font-medium">Headers</CardTitle>
          </CardHeader>
          <CardContent>
            <ScrollArea className="h-72">
              <pre className="text-xs text-violet-400 font-mono leading-relaxed whitespace-pre-wrap break-all">
                {event?.headers
                  ? JSON.stringify(JSON.parse(event.headers), null, 2)
                  : "No headers"}
              </pre>
            </ScrollArea>
          </CardContent>
        </Card>
      </div>

      {/* Retry Timeline */}
      <Card className="bg-slate-900 border-slate-800 mt-6">
        <CardHeader className="pb-3">
          <CardTitle className="text-sm text-slate-300 font-medium">Delivery Attempts</CardTitle>
        </CardHeader>
        <CardContent>
          {(!delivery_attempts || delivery_attempts.length === 0) && (
            <p className="text-slate-500 text-sm">No delivery attempts yet.</p>
          )}
          <div className="space-y-3">
            {delivery_attempts?.map((a: any, i: number) => (
              <div key={a.id} className="flex items-start gap-4 p-3 rounded-lg bg-slate-800/50 border border-slate-700/50">
                <div className={`mt-0.5 w-2 h-2 rounded-full flex-shrink-0 ${a.status_code >= 200 && a.status_code < 300 ? "bg-emerald-500" : "bg-red-500"}`} />
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-3 mb-1">
                    <span className="text-xs font-medium text-slate-300">Attempt #{i + 1}</span>
                    <span className={`text-xs font-mono font-bold ${a.status_code >= 200 && a.status_code < 300 ? "text-emerald-400" : "text-red-400"}`}>
                      HTTP {a.status_code ?? "—"}
                    </span>
                    <span className="text-xs text-slate-500">{a.duration_ms}ms</span>
                  </div>
                  {a.error_message && (
                    <p className="text-xs text-red-400 font-mono truncate">{a.error_message}</p>
                  )}
                  <p className="text-xs text-slate-600 mt-1">
                    {formatDistanceToNow(new Date(a.attempted_at), { addSuffix: true })}
                  </p>
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
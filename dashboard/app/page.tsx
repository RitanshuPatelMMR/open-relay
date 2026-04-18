import { EventFeed } from "@/components/EventFeed"
import { Activity } from "lucide-react"

export default function EventsPage() {
  return (
    <div className="p-8">
      {/* Header */}
      <div className="mb-8">
        <div className="flex items-center gap-3 mb-2">
          <Activity className="w-5 h-5 text-violet-400" />
          <h1 className="text-2xl font-bold text-slate-50">Events</h1>
        </div>
        <p className="text-slate-400 text-sm">
          Live feed of all incoming webhooks. Click any event to inspect payload and delivery attempts.
        </p>
      </div>
      <EventFeed />
    </div>
  )
}
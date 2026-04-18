"use client"
import { useState } from "react"
import { Button } from "@/components/ui/button"
import { replayEvent } from "@/lib/api"
import { RotateCcw } from "lucide-react"

export function ReplayButton({ eventId }: { eventId: string }) {
  const [loading, setLoading] = useState(false)
  const [done, setDone] = useState(false)

  async function handleReplay() {
    setLoading(true)
    await replayEvent(eventId)
    setLoading(false)
    setDone(true)
    setTimeout(() => setDone(false), 3000)
  }

  return (
    <Button
      onClick={handleReplay}
      disabled={loading}
      variant="outline"
      className="border-slate-700 text-slate-300 hover:bg-violet-600/20 hover:border-violet-500 hover:text-violet-300 transition-all"
    >
      <RotateCcw className={`w-4 h-4 mr-2 ${loading ? "animate-spin" : ""}`} />
      {done ? "Requeued!" : loading ? "Replaying…" : "Replay"}
    </Button>
  )
}
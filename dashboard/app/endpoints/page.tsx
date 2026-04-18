"use client"
import { useEffect, useState } from "react"
import { getEndpoints, createEndpoint, deleteEndpoint } from "@/lib/api"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog"
import { Badge } from "@/components/ui/badge"
import { Skeleton } from "@/components/ui/skeleton"
import { Plus, Trash2, Webhook, ExternalLink } from "lucide-react"

export default function EndpointsPage() {
  const [endpoints, setEndpoints] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  const [open, setOpen] = useState(false)
  const [url, setUrl] = useState("")
  const [desc, setDesc] = useState("")
  const [saving, setSaving] = useState(false)

  useEffect(() => {
    getEndpoints().then((d) => { setEndpoints(Array.isArray(d) ? d : []); setLoading(false) })
  }, [])

  async function handleCreate() {
    if (!url) return
    setSaving(true)
    const ep = await createEndpoint(url, desc)
    setEndpoints((prev) => [ep, ...prev])
    setUrl(""); setDesc(""); setOpen(false); setSaving(false)
  }

  async function handleDelete(id: string) {
    await deleteEndpoint(id)
    setEndpoints((prev) => prev.filter((e) => e.id !== id))
  }

  return (
    <div className="p-8 max-w-4xl">
      <div className="flex items-center justify-between mb-8">
        <div>
          <div className="flex items-center gap-3 mb-2">
            <Webhook className="w-5 h-5 text-violet-400" />
            <h1 className="text-2xl font-bold text-slate-50">Endpoints</h1>
          </div>
          <p className="text-slate-400 text-sm">Destination URLs where webhooks get forwarded.</p>
        </div>
        <Dialog open={open} onOpenChange={setOpen}>
          <DialogTrigger asChild>
            <Button className="bg-violet-600 hover:bg-violet-700 text-white">
              <Plus className="w-4 h-4 mr-2" /> Add Endpoint
            </Button>
          </DialogTrigger>
          <DialogContent className="bg-slate-900 border-slate-800 text-slate-50">
            <DialogHeader>
              <DialogTitle className="text-slate-50">New Endpoint</DialogTitle>
            </DialogHeader>
            <div className="space-y-4 mt-2">
              <div>
                <label className="text-xs text-slate-400 mb-1.5 block">Destination URL *</label>
                <input
                  value={url}
                  onChange={(e) => setUrl(e.target.value)}
                  placeholder="https://yourapp.com/webhooks"
                  className="w-full bg-slate-800 border border-slate-700 rounded-lg px-3 py-2.5 text-sm text-slate-50 placeholder-slate-500 focus:outline-none focus:border-violet-500 transition-colors"
                />
              </div>
              <div>
                <label className="text-xs text-slate-400 mb-1.5 block">Description</label>
                <input
                  value={desc}
                  onChange={(e) => setDesc(e.target.value)}
                  placeholder="Production webhook"
                  className="w-full bg-slate-800 border border-slate-700 rounded-lg px-3 py-2.5 text-sm text-slate-50 placeholder-slate-500 focus:outline-none focus:border-violet-500 transition-colors"
                />
              </div>
              <Button
                onClick={handleCreate}
                disabled={saving || !url}
                className="w-full bg-violet-600 hover:bg-violet-700 text-white"
              >
                {saving ? "Creating…" : "Create Endpoint"}
              </Button>
            </div>
          </DialogContent>
        </Dialog>
      </div>

      {loading ? (
        <div className="space-y-3">
          {[...Array(3)].map((_, i) => <Skeleton key={i} className="h-20 w-full bg-slate-800" />)}
        </div>
      ) : endpoints.length === 0 ? (
        <Card className="bg-slate-900 border-slate-800 border-dashed">
          <CardContent className="flex flex-col items-center justify-center py-16">
            <Webhook className="w-10 h-10 text-slate-700 mb-3" />
            <p className="text-slate-500 text-sm">No endpoints yet. Add one to start forwarding webhooks.</p>
          </CardContent>
        </Card>
      ) : (
        <div className="space-y-3">
          {endpoints.map((ep) => (
            <Card key={ep.id} className="bg-slate-900 border-slate-800 hover:border-slate-700 transition-colors">
              <CardContent className="flex items-center justify-between p-4">
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-2 mb-1">
                    <a href={ep.destination_url} target="_blank" rel="noreferrer"
                      className="text-sm text-violet-400 hover:text-violet-300 font-mono truncate flex items-center gap-1">
                      {ep.destination_url}
                      <ExternalLink className="w-3 h-3 flex-shrink-0" />
                    </a>
                    <Badge variant="outline" className={ep.is_active
                      ? "border-emerald-500/30 text-emerald-400 text-xs"
                      : "border-slate-600 text-slate-500 text-xs"}>
                      {ep.is_active ? "active" : "inactive"}
                    </Badge>
                  </div>
                  {ep.description && <p className="text-xs text-slate-500">{ep.description}</p>}
                  <p className="text-xs text-slate-600 mt-1 font-mono">{ep.id}</p>
                </div>
                <Button
                  variant="ghost"
                  size="icon"
                  onClick={() => handleDelete(ep.id)}
                  className="text-slate-600 hover:text-red-400 hover:bg-red-500/10 ml-4 flex-shrink-0"
                >
                  <Trash2 className="w-4 h-4" />
                </Button>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  )
}
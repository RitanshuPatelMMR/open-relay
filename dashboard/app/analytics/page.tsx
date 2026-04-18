"use client"
import { useEffect, useState } from "react"
import { getAnalytics } from "@/lib/api"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Skeleton } from "@/components/ui/skeleton"
import { BarChart3, TrendingUp, CheckCircle, XCircle } from "lucide-react"
import {
  AreaChart, Area, BarChart, Bar, XAxis, YAxis, CartesianGrid,
  Tooltip, ResponsiveContainer, Legend
} from "recharts"

export default function AnalyticsPage() {
  const [data, setData] = useState<any[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    getAnalytics().then((d) => {
      const formatted = (d.hourly ?? []).map((b: any) => ({
        ...b,
        hour: new Date(b.hour).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" }),
      }))
      setData(formatted)
      setLoading(false)
    })
  }, [])

  const totalDelivered = data.reduce((s, d) => s + d.delivered, 0)
  const totalFailed = data.reduce((s, d) => s + d.failed, 0)
  const totalEvents = data.reduce((s, d) => s + d.total, 0)
  const rate = totalEvents > 0 ? Math.round((totalDelivered / totalEvents) * 100) : 0

  return (
    <div className="p-8 max-w-5xl">
      <div className="flex items-center gap-3 mb-2">
        <BarChart3 className="w-5 h-5 text-violet-400" />
        <h1 className="text-2xl font-bold text-slate-50">Analytics</h1>
      </div>
      <p className="text-slate-400 text-sm mb-8">Delivery stats for the last 24 hours.</p>

      {/* Stat cards */}
      <div className="grid grid-cols-3 gap-4 mb-8">
        {[
          { label: "Total Events", value: totalEvents, icon: TrendingUp, color: "text-violet-400" },
          { label: "Delivered", value: totalDelivered, icon: CheckCircle, color: "text-emerald-400" },
          { label: "Failed", value: totalFailed, icon: XCircle, color: "text-red-400" },
        ].map(({ label, value, icon: Icon, color }) => (
          <Card key={label} className="bg-slate-900 border-slate-800">
            <CardContent className="p-5">
              <div className="flex items-center justify-between mb-3">
                <span className="text-xs text-slate-500 font-medium">{label}</span>
                <Icon className={`w-4 h-4 ${color}`} />
              </div>
              {loading
                ? <Skeleton className="h-8 w-16 bg-slate-800" />
                : <p className={`text-3xl font-bold ${color}`}>{value}</p>
              }
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Delivery rate */}
      <Card className="bg-slate-900 border-slate-800 mb-6">
        <CardHeader className="pb-2">
          <CardTitle className="text-sm text-slate-300 font-medium">Delivery Rate (24h)</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex items-end gap-3 mb-4">
            <span className="text-4xl font-bold text-emerald-400">{rate}%</span>
            <span className="text-slate-500 text-sm mb-1">success rate</span>
          </div>
          {loading
            ? <Skeleton className="h-48 w-full bg-slate-800" />
            : <ResponsiveContainer width="100%" height={200}>
                <AreaChart data={data}>
                  <defs>
                    <linearGradient id="delivered" x1="0" y1="0" x2="0" y2="1">
                      <stop offset="5%" stopColor="#10b981" stopOpacity={0.3} />
                      <stop offset="95%" stopColor="#10b981" stopOpacity={0} />
                    </linearGradient>
                    <linearGradient id="failed" x1="0" y1="0" x2="0" y2="1">
                      <stop offset="5%" stopColor="#ef4444" stopOpacity={0.3} />
                      <stop offset="95%" stopColor="#ef4444" stopOpacity={0} />
                    </linearGradient>
                  </defs>
                  <CartesianGrid strokeDasharray="3 3" stroke="#1e293b" />
                  <XAxis dataKey="hour" stroke="#475569" tick={{ fontSize: 11 }} />
                  <YAxis stroke="#475569" tick={{ fontSize: 11 }} />
                  <Tooltip contentStyle={{ background: "#0f172a", border: "1px solid #1e293b", borderRadius: "8px", color: "#f1f5f9" }} />
                  <Legend />
                  <Area type="monotone" dataKey="delivered" stroke="#10b981" fill="url(#delivered)" strokeWidth={2} />
                  <Area type="monotone" dataKey="failed" stroke="#ef4444" fill="url(#failed)" strokeWidth={2} />
                </AreaChart>
              </ResponsiveContainer>
          }
        </CardContent>
      </Card>

      {/* Bar chart */}
      <Card className="bg-slate-900 border-slate-800">
        <CardHeader className="pb-2">
          <CardTitle className="text-sm text-slate-300 font-medium">Events by Hour</CardTitle>
        </CardHeader>
        <CardContent>
          {loading
            ? <Skeleton className="h-48 w-full bg-slate-800" />
            : <ResponsiveContainer width="100%" height={200}>
                <BarChart data={data}>
                  <CartesianGrid strokeDasharray="3 3" stroke="#1e293b" />
                  <XAxis dataKey="hour" stroke="#475569" tick={{ fontSize: 11 }} />
                  <YAxis stroke="#475569" tick={{ fontSize: 11 }} />
                  <Tooltip contentStyle={{ background: "#0f172a", border: "1px solid #1e293b", borderRadius: "8px", color: "#f1f5f9" }} />
                  <Bar dataKey="total" fill="#7c3aed" radius={[4, 4, 0, 0]} />
                </BarChart>
              </ResponsiveContainer>
          }
        </CardContent>
      </Card>
    </div>
  )
}
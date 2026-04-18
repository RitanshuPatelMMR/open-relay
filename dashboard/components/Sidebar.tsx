"use client"
import Link from "next/link"
import { usePathname } from "next/navigation"
import { Activity, Webhook, BarChart3, Settings, Zap } from "lucide-react"
import { cn } from "@/lib/utils"

const nav = [
  { href: "/",           label: "Events",    icon: Activity },
  { href: "/endpoints",  label: "Endpoints", icon: Webhook },
  { href: "/analytics",  label: "Analytics", icon: BarChart3 },
]

export function Sidebar() {
  const path = usePathname()
  return (
    <aside className="fixed left-0 top-0 h-screen w-60 bg-slate-950 border-r border-slate-800 flex flex-col z-50">
      {/* Logo */}
      <div className="flex items-center gap-2 px-6 py-5 border-b border-slate-800">
        <div className="w-8 h-8 bg-violet-600 rounded-lg flex items-center justify-center">
          <Zap className="w-4 h-4 text-white" />
        </div>
        <span className="text-slate-50 font-bold text-lg tracking-tight">OpenRelay</span>
      </div>

      {/* Nav */}
      <nav className="flex-1 px-3 py-4 space-y-1">
        {nav.map(({ href, label, icon: Icon }) => (
          <Link
            key={href}
            href={href}
            className={cn(
              "flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-all",
              path === href
                ? "bg-violet-600/20 text-violet-400 border border-violet-500/30"
                : "text-slate-400 hover:text-slate-50 hover:bg-slate-800"
            )}
          >
            <Icon className="w-4 h-4" />
            {label}
          </Link>
        ))}
      </nav>

      {/* Footer */}
      <div className="px-6 py-4 border-t border-slate-800">
        <p className="text-xs text-slate-600">v0.1.0 — local dev</p>
      </div>
    </aside>
  )
}
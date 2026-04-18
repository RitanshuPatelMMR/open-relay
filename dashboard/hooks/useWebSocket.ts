"use client"
import { useEffect, useRef, useState } from "react"

export type WSMessage = {
  event_id: string
  status: string
  timestamp?: string
}

export function useWebSocket() {
  const [messages, setMessages] = useState<WSMessage[]>([])
  const [connected, setConnected] = useState(false)
  const wsRef = useRef<WebSocket | null>(null)

  useEffect(() => {
    const apiKey = process.env.NEXT_PUBLIC_API_KEY!
    const wsUrl = `${process.env.NEXT_PUBLIC_WS_URL}/ws?api_key=${apiKey}`

    function connect() {
      const ws = new WebSocket(wsUrl)
      wsRef.current = ws

      ws.onopen = () => setConnected(true)
      ws.onclose = () => {
        setConnected(false)
        setTimeout(connect, 3000) // reconnect
      }
      ws.onmessage = (e) => {
        try {
          const msg = JSON.parse(e.data)
          setMessages((prev) => [msg, ...prev].slice(0, 100))
        } catch {}
      }
    }

    connect()
    return () => wsRef.current?.close()
  }, [])

  return { messages, connected }
}
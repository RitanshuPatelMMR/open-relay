const API_URL = process.env.NEXT_PUBLIC_API_URL!
const API_KEY = process.env.NEXT_PUBLIC_API_KEY!

const headers = {
  "Content-Type": "application/json",
  "X-API-Key": API_KEY,
}

export async function getEvents(status?: string) {
  const url = status
    ? `${API_URL}/api/events?status=${status}&limit=100`
    : `${API_URL}/api/events?limit=100`
  const res = await fetch(url, { headers, cache: "no-store" })
  return res.json()
}

export async function getEvent(id: string) {
  const res = await fetch(`${API_URL}/api/events/${id}`, { headers, cache: "no-store" })
  return res.json()
}

export async function replayEvent(id: string) {
  const res = await fetch(`${API_URL}/api/events/${id}/replay`, {
    method: "POST",
    headers,
  })
  return res.json()
}

export async function getAnalytics() {
  const res = await fetch(`${API_URL}/api/analytics`, { headers, cache: "no-store" })
  return res.json()
}

export async function getEndpoints() {
  const res = await fetch(`${API_URL}/api/endpoints`, { headers, cache: "no-store" })
  return res.json()
}

export async function createEndpoint(destination_url: string, description: string) {
  const res = await fetch(`${API_URL}/api/endpoints`, {
    method: "POST",
    headers,
    body: JSON.stringify({ destination_url, description }),
  })
  return res.json()
}

export async function deleteEndpoint(id: string) {
  await fetch(`${API_URL}/api/endpoints/${id}`, { method: "DELETE", headers })
}
export interface Project {
  id: string;
  name: string;
  api_key: string;
  created_at: string;
}

export interface Endpoint {
  id: string;
  project_id: string;
  destination_url: string;
  description: string | null;
  is_active: boolean;
  created_at: string;
}

export interface DeliveryAttempt {
  id: string;
  event_id: string;
  attempted_at: string;
  status_code: number | null;
  response_body: string | null;
  duration_ms: number | null;
  error_message: string | null;
  next_retry_at: string | null;
}

export interface Event {
  id: string;
  project_id: string;
  endpoint_id: string | null;
  idempotency_key: string | null;
  source_ip: string | null;
  method: string;
  path: string | null;
  headers: Record<string, string> | null;
  payload: unknown;
  status: "pending" | "delivered" | "failed";
  attempt_count: number;
  created_at: string;
  delivered_at: string | null;
  delivery_attempts?: DeliveryAttempt[];
}

export interface ListEventsParams {
  status?: "pending" | "delivered" | "failed";
  limit?: number;
  offset?: number;
}

export interface CreateEndpointParams {
  destination_url: string;
  description?: string;
}

export interface UpdateEndpointParams {
  destination_url?: string;
  description?: string;
  is_active?: boolean;
}

export interface CreateProjectParams {
  name: string;
}

export interface OpenRelayClientConfig {
  apiKey: string;
  baseUrl?: string;
}
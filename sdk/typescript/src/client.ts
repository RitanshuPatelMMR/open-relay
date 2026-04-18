import {
  OpenRelayClientConfig,
  Event,
  ListEventsParams,
  Endpoint,
  CreateEndpointParams,
  UpdateEndpointParams,
  Project,
  CreateProjectParams,
} from "./types";

export class OpenRelayClient {
  private apiKey: string;
  private baseUrl: string;

  constructor(config: OpenRelayClientConfig) {
    if (!config.apiKey) throw new Error("apiKey required");
    this.apiKey = config.apiKey;
    this.baseUrl = (config.baseUrl ?? "http://localhost:8081").replace(/\/$/, "");
  }

  private async request<T>(method: string, path: string, body?: unknown): Promise<T> {
    const res = await fetch(`${this.baseUrl}${path}`, {
      method,
      headers: {
        "Content-Type": "application/json",
        "X-API-Key": this.apiKey,
      },
      body: body ? JSON.stringify(body) : undefined,
    });

    if (!res.ok) {
      const text = await res.text();
      throw new Error(`OpenRelay API error ${res.status}: ${text}`);
    }

    if (res.status === 204) return undefined as T;
    return res.json() as Promise<T>;
  }

  // Events
  async getEvents(params: ListEventsParams = {}): Promise<Event[]> {
    const query = new URLSearchParams();
    if (params.status) query.set("status", params.status);
    if (params.limit !== undefined) query.set("limit", String(params.limit));
    if (params.offset !== undefined) query.set("offset", String(params.offset));
    const qs = query.toString();
    return this.request<Event[]>("GET", `/api/events${qs ? `?${qs}` : ""}`);
  }

  async getEvent(eventId: string): Promise<Event> {
    return this.request<Event>("GET", `/api/events/${eventId}`);
  }

  async replayEvent(eventId: string): Promise<void> {
    return this.request<void>("POST", `/api/events/${eventId}/replay`);
  }

  // Endpoints
  async getEndpoints(): Promise<Endpoint[]> {
    return this.request<Endpoint[]>("GET", "/api/endpoints");
  }

  async createEndpoint(params: CreateEndpointParams): Promise<Endpoint> {
    return this.request<Endpoint>("POST", "/api/endpoints", params);
  }

  async updateEndpoint(endpointId: string, params: UpdateEndpointParams): Promise<Endpoint> {
    return this.request<Endpoint>("PUT", `/api/endpoints/${endpointId}`, params);
  }

  async deleteEndpoint(endpointId: string): Promise<void> {
    return this.request<void>("DELETE", `/api/endpoints/${endpointId}`);
  }

  // Projects
  async getProjects(): Promise<Project[]> {
    return this.request<Project[]>("GET", "/api/projects");
  }

  async createProject(params: CreateProjectParams): Promise<Project> {
    return this.request<Project>("POST", "/api/projects", params);
  }
}
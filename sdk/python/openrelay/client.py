from __future__ import annotations
from typing import List, Optional
import httpx
from .models import Event, Endpoint, Project


class OpenRelayClient:
    def __init__(self, api_key: str, base_url: str = "http://localhost:8081"):
        if not api_key:
            raise ValueError("api_key required")
        self._headers = {
            "X-API-Key": api_key,
            "Content-Type": "application/json",
        }
        self._base = base_url.rstrip("/")

    def _request(self, method: str, path: str, **kwargs) -> httpx.Response:
        res = httpx.request(method, f"{self._base}{path}", headers=self._headers, **kwargs)
        res.raise_for_status()
        return res

    # Events
    def get_events(
        self,
        status: Optional[str] = None,
        limit: int = 50,
        offset: int = 0,
    ) -> List[Event]:
        params = {"limit": limit, "offset": offset}
        if status:
            params["status"] = status
        res = self._request("GET", "/api/events", params=params)
        return [Event(**e) for e in res.json()]

    def get_event(self, event_id: str) -> Event:
        res = self._request("GET", f"/api/events/{event_id}")
        return Event(**res.json())

    def replay_event(self, event_id: str) -> None:
        self._request("POST", f"/api/events/{event_id}/replay")

    # Endpoints
    def get_endpoints(self) -> List[Endpoint]:
        res = self._request("GET", "/api/endpoints")
        return [Endpoint(**e) for e in res.json()]

    def create_endpoint(self, destination_url: str, description: Optional[str] = None) -> Endpoint:
        body = {"destination_url": destination_url}
        if description:
            body["description"] = description
        res = self._request("POST", "/api/endpoints", json=body)
        return Endpoint(**res.json())

    def update_endpoint(
        self,
        endpoint_id: str,
        destination_url: Optional[str] = None,
        description: Optional[str] = None,
        is_active: Optional[bool] = None,
    ) -> Endpoint:
        body = {}
        if destination_url is not None:
            body["destination_url"] = destination_url
        if description is not None:
            body["description"] = description
        if is_active is not None:
            body["is_active"] = is_active
        res = self._request("PUT", f"/api/endpoints/{endpoint_id}", json=body)
        return Endpoint(**res.json())

    def delete_endpoint(self, endpoint_id: str) -> None:
        self._request("DELETE", f"/api/endpoints/{endpoint_id}")

    # Projects
    def get_projects(self) -> List[Project]:
        res = self._request("GET", "/api/projects")
        return [Project(**p) for p in res.json()]

    def create_project(self, name: str) -> Project:
        res = self._request("POST", "/api/projects", json={"name": name})
        return Project(**res.json())
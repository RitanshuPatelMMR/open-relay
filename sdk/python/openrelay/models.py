from __future__ import annotations
from typing import Any, Dict, Optional, Union
import json
from pydantic import BaseModel, field_validator


class Project(BaseModel):
    id: str
    name: str
    api_key: str
    created_at: str


class Endpoint(BaseModel):
    id: str
    project_id: str
    destination_url: str
    description: Optional[str] = None
    is_active: bool
    created_at: str


class DeliveryAttempt(BaseModel):
    id: str
    event_id: str
    attempted_at: str
    status_code: Optional[int] = None
    response_body: Optional[str] = None
    duration_ms: Optional[int] = None
    error_message: Optional[str] = None
    next_retry_at: Optional[str] = None


class Event(BaseModel):
    id: str
    project_id: str
    endpoint_id: Optional[str] = None
    idempotency_key: Optional[str] = None
    source_ip: Optional[str] = None
    method: str
    path: Optional[str] = None
    headers: Optional[Any] = None
    payload: Any = None
    status: str
    attempt_count: int
    created_at: str
    delivered_at: Optional[str] = None
    delivery_attempts: Optional[list[DeliveryAttempt]] = None

    @field_validator("headers", "payload", mode="before")
    @classmethod
    def parse_json_string(cls, v: Any) -> Any:
        if isinstance(v, str):
            try:
                return json.loads(v)
            except Exception:
                return v
        return v
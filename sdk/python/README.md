# openrelay-sdk (Python)

Python SDK for [OpenRelay](https://github.com/RitanshuPatelMMR/open-relay) — open source webhook delivery infrastructure.

## Install

```bash
pip install openrelay-sdk
```

## Quick Start

```python
from openrelay import OpenRelayClient

client = OpenRelayClient(
    api_key="your-api-key",
    base_url="https://your-openrelay-instance.com",
)

# List failed events
events = client.get_events(status="failed", limit=20)
for e in events:
    print(e.id, e.status)

# Replay an event
client.replay_event("event-id-here")

# Create endpoint
endpoint = client.create_endpoint(
    destination_url="https://myapp.com/webhooks",
    description="Production endpoint",
)
```

## Methods

| Method | Description |
|--------|-------------|
| `get_events(status?, limit?, offset?)` | List events |
| `get_event(id)` | Get event with delivery attempts |
| `replay_event(id)` | Re-queue failed event |
| `get_endpoints()` | List endpoints |
| `create_endpoint(url, description?)` | Create endpoint |
| `update_endpoint(id, ...)` | Update endpoint |
| `delete_endpoint(id)` | Delete endpoint |
| `get_projects()` | List projects |
| `create_project(name)` | Create project |

## License

MIT
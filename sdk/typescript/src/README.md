# openrelay-sdk

TypeScript SDK for [OpenRelay](https://github.com/RitanshuPatelMMR/open-relay) — open source webhook delivery infrastructure.

## Install

```bash
npm install openrelay-sdk
```

## Quick Start

```typescript
import { OpenRelayClient } from 'openrelay-sdk';

const client = new OpenRelayClient({
  apiKey: 'your-api-key',
  baseUrl: 'https://your-openrelay-instance.com',
});

// List recent failed events
const events = await client.getEvents({ status: 'failed', limit: 20 });

// Replay a failed event
await client.replayEvent('event-id-here');

// Create an endpoint
const endpoint = await client.createEndpoint({
  destinationUrl: 'https://myapp.com/webhooks',
  description: 'Production webhook endpoint',
});
```

## Methods

| Method | Description |
|--------|-------------|
| `getEvents(params?)` | List events. Filter by `status`, `limit`, `offset` |
| `getEvent(id)` | Get single event with delivery attempts |
| `replayEvent(id)` | Re-queue a failed event for delivery |
| `getEndpoints()` | List all endpoints |
| `createEndpoint(params)` | Create delivery endpoint |
| `updateEndpoint(id, params)` | Update endpoint |
| `deleteEndpoint(id)` | Delete endpoint |
| `getProjects()` | List projects |
| `createProject(params)` | Create new project |

## License

MIT
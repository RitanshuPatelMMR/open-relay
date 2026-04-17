-- projects: one project = one team / one app
CREATE TABLE projects (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name        VARCHAR(255) NOT NULL,
  api_key     VARCHAR(255) UNIQUE NOT NULL DEFAULT gen_random_uuid()::text,
  created_at  TIMESTAMPTZ DEFAULT NOW()
);

-- endpoints: where to forward webhooks
CREATE TABLE endpoints (
  id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  project_id      UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
  destination_url TEXT NOT NULL,
  description     VARCHAR(255),
  is_active       BOOLEAN DEFAULT TRUE,
  created_at      TIMESTAMPTZ DEFAULT NOW()
);

-- events: every incoming webhook
CREATE TABLE events (
  id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  project_id       UUID NOT NULL REFERENCES projects(id),
  endpoint_id      UUID REFERENCES endpoints(id),
  idempotency_key  VARCHAR(255),
  source_ip        VARCHAR(45),
  method           VARCHAR(10) NOT NULL,
  path             TEXT,
  headers          JSONB,
  payload          JSONB,
  status           VARCHAR(20) NOT NULL DEFAULT 'pending',
  attempt_count    INTEGER DEFAULT 0,
  created_at       TIMESTAMPTZ DEFAULT NOW(),
  delivered_at     TIMESTAMPTZ,
  UNIQUE(project_id, idempotency_key)
);

-- delivery_attempts: log of every delivery try
CREATE TABLE delivery_attempts (
  id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  event_id       UUID NOT NULL REFERENCES events(id),
  attempted_at   TIMESTAMPTZ DEFAULT NOW(),
  status_code    INTEGER,
  response_body  TEXT,
  duration_ms    INTEGER,
  error_message  TEXT,
  next_retry_at  TIMESTAMPTZ
);
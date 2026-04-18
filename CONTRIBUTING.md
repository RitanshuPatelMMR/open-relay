# Contributing to OpenRelay

Thanks for your interest in contributing. This guide gets you from zero to first PR.

## Setup

```bash
git clone https://github.com/RitanshuPatelMMR/open-relay.git
cd open-relay
cp .env.example .env
docker compose up -d
```

Verify:
```bash
curl http://localhost:8080/health  # ingestion
curl http://localhost:8081/health  # api
```

## Project Structure
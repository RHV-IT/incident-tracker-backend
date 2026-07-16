# Incident Tracker

RESTful API for managing workplace incidents and safety reports. Built with Go, Gin, and PostgreSQL.

## Prerequisites

- Docker and Docker Compose
- Go 1.22+

## Setup

```bash
git clone <repository-url>
cd incident-tracker-backend
cp .env.example .env
go mod download
```

Start with Docker:

```bash
docker compose up -d
```

API available at `http://localhost:3002`.

## Development

```bash
air
```

API available at `http://localhost:3001`.

## Tests

```bash
go test -v -tags=test ./...
```

Or run the helper script:

```bash
./scripts/runtests.sh
```

## Quick Usage

Login to get a token:

```bash
TOKEN=$(curl -s -X POST http://localhost:3002/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"yourpassword"}' | jq -r '.token')
```

Report an incident (no auth required):

```bash
curl -X POST http://localhost:3002/api/v1/incidents \
  -H "Content-Type: application/json" \
  -d '{ ... }'
```

List incidents:

```bash
curl http://localhost:3002/api/v1/incidents -H "Authorization: Bearer $TOKEN"
```

For full API documentation, request/response schemas, and role permissions, see [SYSTEM_DESIGN.md](SYSTEM_DESIGN.md).

For architecture, layering, and design decisions, see [ARCHITECTURE.md](ARCHITECTURE.md).

For database schema details, see `tables.sql`.

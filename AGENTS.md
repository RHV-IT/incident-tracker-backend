# Agent Instructions for Issue Tracker

## Project Overview

The Issue Tracker is a RESTful API for managing workplace incidents and safety reports built with Go, Gin, and PostgreSQL.

## Development Commands

```bash
# Start development server with live reload
air

# Run application directly
go run ./cmd/main.go

# Run tests
go test ./...

# Format code
go fmt ./...

# Run linter
go vet ./...
```

## Docker Commands

```bash
# Start all services
docker compose up -d

# Stop services
docker compose down

# Remove volumes (fresh database)
docker compose down -v

# View logs
docker compose logs -f
```

## Database Access

```bash
# Access PostgreSQL shell
./login.sh
```

## API Testing

```bash
# Health check
curl http://localhost:3002/api/v1/ping

# Login (save token)
TOKEN=$(curl -s -X POST http://localhost:3002/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"yourpassword"}' | jq -r '.token')

# Report incident (no auth required)
curl -X POST http://localhost:3002/api/v1/incidents \
  -H "Content-Type: application/json" \
  -d '{"reporterName":"John","department":"IT","position":"Dev","contactInfo":"john@example.com","dateOfIncident":"2026-06-09","timeOfIncident":"14:00","locationOfIncident":"Office","typeOfIncident":"Slip","peopleInvolved":"None","descriptionOfIncident":"Test","immediateActionTaken":"Clean","injuryOrDamage":"None","severityLevel":"minor","supervisorNotified":"Yes","recommendedPreventiveAction":"None"}'

# Get incidents
curl http://localhost:3002/api/v1/incidents -H "Authorization: Bearer $TOKEN"
```

## Role Permissions

| Role | Permissions |
|------|-------------|
| superadmin | All endpoints, user management |
| admin | Report incidents, view department incidents |
| supervisor | Report incidents, view own department incidents |
| reporter | Report incidents via public endpoint only |

## Default Credentials

A superadmin user is created by default:
- Email: `admin@example.com`
- Password: Check the database or reset via code
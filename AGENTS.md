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

# Register a new user (requires superadmin token)
curl -X POST http://localhost:3002/api/v1/auth/register \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"email":"newuser@example.com","name":"New User","password":"password123","role":"admin","department":"IT"}'

# Report incident (no auth required)
curl -X POST http://localhost:3002/api/v1/incidents \
  -H "Content-Type: application/json" \
  -d '{"reporterName":"John","department":"IT","position":"Dev","contactInfo":"john@example.com","dateOfIncident":"2026-06-09","timeOfIncident":"14:00","locationOfIncident":"Office","typeOfIncident":"Slip","peopleInvolved":"None","descriptionOfIncident":"Test","immediateActionTaken":"Clean","injuryOrDamage":"None","severityLevel":"minor","supervisorNotified":"Yes","recommendedPreventiveAction":"None"}'

# Get incidents (requires auth)
curl http://localhost:3002/api/v1/incidents -H "Authorization: Bearer $TOKEN"

# Get incidents with pagination
curl "http://localhost:3002/api/v1/incidents?page=1&limit=20" -H "Authorization: Bearer $TOKEN"

# Update incident status (requires auth)
curl -X PATCH http://localhost:3002/api/v1/incidents/1/status \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status":"resolved"}'

# Get user info (requires auth)
curl "http://localhost:3002/api/v1/user?email=test@example.com" -H "Authorization: Bearer $TOKEN"

# Disable user (requires superadmin)
curl -X PUT http://localhost:3002/api/v1/auth/disable \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com"}'

# Enable user (requires superadmin)
curl -X PUT http://localhost:3002/api/v1/auth/enable \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com"}'

# Reset user password (requires superadmin)
curl -X PUT http://localhost:3002/api/v1/auth/resetpassword \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"newpassword123"}'
```

## Role Permissions

| Role | Permissions |
|------|-------------|
| superadmin | All endpoints including user management (register, update, disable, enable, reset password, get user), report incidents, view all incidents, update incident status |
| admin | Report incidents, view all incidents, update incident status |
| supervisor | Report incidents, view own department incidents, update own department incidents |
| reporter | Report incidents via public endpoint only |

## Default Credentials

A superadmin user is created by default:
- Email: `admin@example.com`
- Password: The default password is hashed with bcrypt. Check the database or reset via code to set a known password.
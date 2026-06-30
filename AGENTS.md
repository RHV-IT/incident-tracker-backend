# Agent Instructions for Issue Tracker

## Project Overview

The Issue Tracker is a RESTful API for managing workplace incidents and safety reports built with Go, Gin, and PostgreSQL.

**Code Metrics:**
- Total Go code: 1443 lines
- 16 Go source files
- Architecture: Clean layered (presentation → application → data → infrastructure)

## Development Commands

```bash
# Start development server with live reload
air

# Run application directly
go run ./cmd/

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
  -d '{
    "principalName": "John Doe",
    "principalGender": "Male",
    "principalDob": "1990-01-15",
    "principalType": "patient",
    "patientId": "P12345",
    "patientWardDept": "Ward A",
    "peopleInvolved": "Nurse Smith",
    "dateOfIncident": "2026-06-09",
    "timeOfIncident": "14:00",
    "locationOfIncident": "Ward A, Room 3",
    "incidentWardDept": "Ward A",
    "witnesses": "Dr. Brown",
    "witnessType": "Staff",
    "witnessWardDept": "Ward A",
    "witnessJobTitle": "Doctor",
    "witnessPhone": "555-0100",
    "isNearMiss": false,
    "causeGroup": "Fall",
    "causes": "Wet floor",
    "prescribingDoctor": "Dr. Brown",
    "treatmentReceived": "First Aid",
    "equipmentInvolved": false,
    "equipmentSentForRepair": false,
    "equipmentWithdrawn": false,
    "equipmentRetained": false,
    "isMedicalDevice": false,
    "reporterName": "Jane Reporter",
    "reporterDesignation": "Nurse",
    "signature": true,
    "reporterInfo": "jane@example.com",
    "reporterDate": "2026-06-09",
    "severityLevel": "minor"
  }'

# Get incidents (requires auth)
curl http://localhost:3002/api/v1/incidents -H "Authorization: Bearer $TOKEN"

# Update incident status (requires auth; reporter role forbidden)
curl -X PATCH http://localhost:3002/api/v1/incidents/1/status \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status":"resolved"}'

# Get user info (requires auth)
curl "http://localhost:3002/api/v1/user?email=test@example.com" -H "Authorization: Bearer $TOKEN"
```

## Role Permissions

| Role | Permissions |
|------|-------------|
| superadmin | All endpoints including user management (register, update, disable, enable, reset password, get user), report incidents, view all incidents, update any incident status, submit incident management reports |
| admin | Report incidents, view all incidents, update any incident status, submit incident management reports |
| supervisor | Report incidents, view own department incidents (via `incident_ward_dept`), update own department incident status |
| reporter | Report incidents via public endpoint only, view own department incidents |

## Default Credentials

A superadmin user is created by default:
- Email: `admin@example.com`
- Password: The default password is hashed with bcrypt. Check the database or reset via code to set a known password.

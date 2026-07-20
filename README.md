# Incident Tracker

RESTful API for managing workplace incidents and safety reports. Built with Go, Gin, and PostgreSQL.

## Prerequisites

- Docker and Docker Compose
- Go 1.26+

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

Or run directly:

```bash
go run ./cmd/
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

Format code:

```bash
go fmt ./...
```

Run linter:

```bash
go vet ./...
```

## Scripts

```bash
# Access PostgreSQL shell
./scripts/login.sh

# Reset database (drop and recreate tables)
./scripts/resetdb.sh

# Recreate tables without dropping data
./scripts/createtable.sh
```

## Default Credentials

A superadmin user is seeded by default:

- Email: `admin@example.com`
- Password: The password is stored as a bcrypt hash in `tables.sql`. Use the database or reset it via code to set a known password.

## Quick Usage

Login to get a token:

```bash
TOKEN=$(curl -s -X POST http://localhost:3002/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"yourpassword"}' | jq -r '.token')
```

Register a new user (requires superadmin token):

```bash
curl -X POST http://localhost:3002/api/v1/auth/register \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"email":"newuser@example.com","name":"New User","password":"password123","role":"admin","department":"IT"}'
```

Report an incident (no auth required):

```bash
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
    "witenssPhone": "555-0100",
    "isNearMiss": false,
    "causeGroup": "Fall",
    "causes": "Wet floor",
    "prescribingDoctor": "Dr. Brown",
    "treatmentReceived": "First Aid",
    "equipmentInvolved": "No",
    "equipmentSentForRepair": false,
    "equipmentWithdrawn": false,
    "equipmentRetained": false,
    "isMedicalDevice": "No",
    "reporterName": "Jane Reporter",
    "reporterDesignation": "Nurse",
    "signature": true,
    "reporterInfo": "jane@example.com",
    "date": "2026-06-09",
    "severityLevel": "minor"
  }'
```

List incidents (requires auth):

```bash
curl http://localhost:3002/api/v1/incidents -H "Authorization: Bearer $TOKEN"
```

Add comment to incident (requires manager or admin):

```bash
curl -X POST http://localhost:3002/api/v1/incidents/comments \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"incidentId": 1, "userId": 2, "comment": "Follow up needed"}'
```

Update incident status (requires auth; reporter/supervisor/manager roles forbidden):

```bash
curl -X PATCH http://localhost:3002/api/v1/incidents/1/status \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status":"resolved"}'
```

Get user info (requires superadmin role):

```bash
curl "http://localhost:3002/api/v1/user?email=test@example.com" -H "Authorization: Bearer $TOKEN"
```

Get comments for incident (requires admin or manager):

```bash
curl "http://localhost:3002/api/v1/incidents/comments?incidentId=1" -H "Authorization: Bearer $TOKEN"
```

Get incident management logs (requires admin role):

```bash
curl "http://localhost:3002/api/v1/incidents/1/managementlogs" -H "Authorization: Bearer $TOKEN"
```

## Role Permissions

| Role | Permissions |
|------|-------------|
| superadmin | All endpoints including user management (register, update, disable, enable, reset password, get user), report incidents, view all incidents, update any incident status, submit incident management reports, update incident management reports, add comments, view comments |
| admin | Report incidents, view all incidents, update any incident status, submit incident management reports, update incident management reports, add comments, view comments |
| supervisor | Report incidents, view own department incidents (matched via `incident_ward_dept`, `patient_ward_dept`, or `staff_place_of_work`) |
| manager | Add comments, submit incident management reports, update incident management reports, view all incidents |
| reporter | Report incidents via public endpoint only, view own department incidents |

## Docker Commands

```bash
# Start all services (API at localhost:3002)
docker compose up -d

# Stop services
docker compose down

# Remove volumes (fresh database)
docker compose down -v

# View logs
docker compose logs -f
```

For full API documentation, request/response schemas, and role permissions, see [SYSTEM_DESIGN.md](SYSTEM_DESIGN.md).

For architecture, layering, and design decisions, see [ARCHITECTURE.md](ARCHITECTURE.md).

For database schema details, see `tables.sql`.

For user API routes, roles, inputs, and outputs, see [users.md](users.md).

# Agent Instructions for Issue Tracker

## Project Overview

The Issue Tracker is a RESTful API for managing workplace incidents and safety reports built with Go, Gin, and PostgreSQL.

**Code Metrics:**
- Total Go code: ~1800 lines
- 20 Go source files
- Architecture: Clean layered (presentation → application → data → infrastructure)

## Development Commands

```bash
# Start development server with live reload
air

# Run application directly
go run ./cmd/

# Run tests
go test ./...
# Or for verbose output:
./scripts/runtests.sh

# Format code
go fmt ./...

# Run linter
go vet ./...
```

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

## Database Access

```bash
# Access PostgreSQL shell
./scripts/login.sh
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

# Add comment to incident (requires manager or admin)
curl -X POST http://localhost:3002/api/v1/incidents/comments \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"incidentId": 1, "userId": 2, "comment": "Follow up needed"}'

# Get incidents (requires auth)
curl http://localhost:3002/api/v1/incidents -H "Authorization: Bearer $TOKEN"

# Update incident status (requires auth; reporter/supervisor/manager roles forbidden)
curl -X PATCH http://localhost:3002/api/v1/incidents/1/status \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status":"resolved"}'

# Get user info (requires superadmin role)
curl "http://localhost:3002/api/v1/user?email=test@example.com" -H "Authorization: Bearer $TOKEN"

# Get comments for incident (requires admin or manager)
curl "http://localhost:3002/api/v1/incidents/comments?incidentId=1" -H "Authorization: Bearer $TOKEN"

# Get incident management logs (requires admin role)
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

## Incident Management Form

The incident management report captures follow-up documentation after an incident occurs. The form includes:

### Form Sections

1. **Operational Evaluation Metrics**
   - Impact on Service
   - Contributory Factors
   - Actions Taken / Outcomes
   - Recommendations
   - Lessons Learned

2. **Stakeholder Notifications Log**
   - Patient Informed
   - Relative Informed
   - Senior Manager Notified
   - Pharmacist Informed
   - Police Incident Number
   - Other Informed Parties

3. **Risk Factor Assessment**
   - Risk Severity Score (1-5)
   - Risk Likelihood Score (1-5)
   - Risk Rating (submitted value; typically Severity × Likelihood)

4. **Occupational Health & Safety Compliance**
   - Staff Absence Over 3 Days
   - Act of Violence or Peril Danger
   - Hospitalization Over 24 Hours
   - OHS Impacted Staff Name
   - Staff Date of Birth
   - Staff Home Address

5. **Executive Authorization Sign-Off**
   - Manager Name
   - Corporate Designation
   - Authorization Date
   - Manager Signature (required, legally binding)

### API Endpoints

**POST /api/v1/incidents/:id/management** - Create management report
- Requires: admin or manager role
- Request body: `IncidentManagement` struct

**GET /api/v1/incidents/:id/management** - Retrieve management report
- Requires: Authentication (any authenticated user)
- Response: `IncidentManagement` struct

**PUT /api/v1/incidents/:id/management** - Update existing report
- Requires: manager or admin role
- Request body: `IncidentManagement` struct

### IncidentManagement Struct

```go
type IncidentManagement struct {
    ID                              int    `json:"id"`
    IncidentID                      int    `json:"incidentId"`
    ImpactOnService                 string `json:"impactOnService"`
    ContributoryFactors             string `json:"contributoryFactors"`
    ActionsTakenOutcomes            string `json:"actionsTakenOutcomes"`
    Recommendations                 string `json:"recommendations"`
    LessonsLearned                  string `json:"lessonsLearned"`
    InformedPatient                 bool   `json:"informedPatient"`
    InformedRelative                bool   `json:"informedRelative"`
    InformedSeniorManager           bool   `json:"informedSeniorManager"`
    InformedPharmacist              bool   `json:"informedPharmacist"`
    PoliceIncidentNumber            string `json:"policeIncidentNumber"`
    InformedOther                   string `json:"informedOther"`
    RiskSeverity                    int    `json:"riskSeverity"`
    RiskLikelihood                  int    `json:"riskLikelihood"`
    RiskRating                      int    `json:"riskRating"`
    OhsAbsenceOver3Days             bool   `json:"ohsAbsenceOver3Days"`
    OhsActOfViolenceOrDanger        bool   `json:"ohsActOfViolenceOrDanger"`
    OhsHospitalizationOver24Hours   bool   `json:"ohsHospitalizationOver24Hours"`
    OhsStaffName                    string `json:"ohsStaffName"`
    OhsStaffDob                     string `json:"ohsStaffDob"`
    OhsStaffAddress                 string `json:"ohsStaffAddress"`
    ManagerName                     string `json:"managerName"`
    ManagerSignature                bool   `json:"managerSignature"`
    ManagerDesignation              string `json:"managerDesignation"`
    ManagerDate                     string `json:"managerDate"`
}
```

## Default Credentials

A superadmin user is created by default:
- Email: `admin@example.com`
- Password: The default password is hashed with bcrypt. Check the database or reset via code to set a known password.

# Incident Tracker

A RESTful API for managing workplace incidents and safety reports built with Go, Gin, and PostgreSQL.

**Code Metrics:**
- Total Go code: ~1800 lines
- 20 Go source files
- Architecture: Clean layered (presentation → application → data → infrastructure)

## Overview

The Incident Tracker is a web application designed to help organizations (particularly healthcare settings) track and manage workplace incidents, safety reports, and clinical events. It provides user authentication, role-based access control, comprehensive incident reporting, and incident management follow-up documentation.

## Features

- **User Authentication**
  - User registration (requires superadmin role)
  - Secure login with JWT tokens
  - User profile updates (superadmin only)
  - Account disable/enable functionality (superadmin only)
  - Password reset (superadmin only)
  - Disabled account enforcement at login

- **Logging & Observability**
  - Structured logging with separate log files for errors and incident updates
  - Background error logging to `backgroundErrors.log`
  - Incident update logging to `updateIncidents.log`
  - Automatic audit logging for incident updates (stored in `incident_logs` table)

- **Incident Management**
   - Report new incidents via public endpoint (no auth required)
   - Comprehensive 37-field clinical incident form
   - Track incident severity levels (Near Miss, Minor, Major, Critical)
   - Track incident lifecycle status (Unresolved → In Progress → Resolved)
   - Principal person involved details (patient, staff, consultant, other)
   - Witness information (names, types, departments, contact)
   - Equipment involvement tracking (models, serial numbers, disposition)
   - Treatment received and prescribing doctor fields
   - Reporter details section (name, designation, signature, date)
   - Paginated incident listing with metadata
   - Follow-up incident management reports (admin/manager only)
   - Comments on incidents (manager/admin only)

- **Role-Based Access Control**
   - Five distinct roles: Reporter, Supervisor, Manager, Admin, Superadmin
   - Role-based endpoint protection via JWT middleware
   - Department-based data scoping for supervisors and reporters
   - Superadmin privileges for user management

- **Development Experience**
   - Docker Compose setup for easy development
   - Hot reload with Air for live reloading
   - Comprehensive environment variable configuration
   - Helper scripts for common operations
   - Unit tests for handlers and routes

## Technology Stack

- **Language**: Go 1.26.3
- **Web Framework**: Gin-Gonic
- **Database**: PostgreSQL 16 with PGX driver (connection pool)
- **Authentication**: JWT (HS256, 72-hour expiry) with bcrypt password hashing
- **Development Tools**: Air (live reload), Docker Compose
- **CORS**: gin-contrib/cors middleware for Cross-Origin Resource Sharing
- **Validation**: go-playground/validator via Gin binding

## Getting Started

### Prerequisites

- Docker and Docker Compose (for containerized setup)
- Go 1.22+ (if running locally without Docker)
- Git (for version control)

## Setup

### Clone and Install

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd incident-tracker-backend
   ```

2. Install Go dependencies:
   ```bash
   go mod download
   ```

3. Copy the environment template and configure it:
   ```bash
   cp .env.example .env
   ```
   - Edit `.env` and update `jwtSecret` (minimum 32 characters) and `allowedOrigins` as needed.
   - If running locally (non-Docker), ensure `dbConnStr` points to your local PostgreSQL instance.

4. Set up the database schema:
   - **With Docker**: The schema is automatically loaded from `tables.sql` when the PostgreSQL container starts.
   - **Without Docker**: Manually run `tables.sql` against your local PostgreSQL database.

---

## Running with Docker

1. Start all services:
   ```bash
   docker compose up -d
   ```

2. The API will be available at `http://localhost:3002`
   - The server runs on port 3001 internally (PORT env var)
   - Port 3002 on host is mapped to port 3001 in container
   - PostgreSQL initializes automatically with `tables.sql`
   - Docker compose mounts `logs/` directory for log output

3. Stop services:
   ```bash
   docker compose down
   ```

4. Remove volumes (fresh database):
   ```bash
   docker compose down -v
   ```

5. View logs:
   ```bash
   docker compose logs -f
   ```

6. Access database shell:
   ```bash
   ./scripts/login.sh
   ```

---

## Local Development (Without Docker)

1. Ensure a local PostgreSQL instance is running on port `5432` with:
   - Database: `issuetracker`
   - User: `tracker_user`
   - Password: `tracker_password`

2. Apply the schema:
   ```bash
   psql -U tracker_user -d issuetracker -f tables.sql
   ```

3. Start the application with live reload:
   ```bash
   air
   ```

   Or run directly:
   ```bash
   go run ./cmd/
   ```

4. The API will be available at `http://localhost:3001`.

---

## Running Tests

### Prerequisites

- Docker must be installed and running locally. Tests use [`testcontainers-go`](https://golang.testcontainers.org/) to spin up a temporary PostgreSQL container automatically.
- Go 1.22+

### Commands

Run all tests (this will automatically start and clean up a PostgreSQL test container):
```bash
go test -v -tags=test ./...
```

Or use the helper script:
```bash
./scripts/runtests.sh
```

Run tests for a specific package:
```bash
go test -v -tags=test ./cmd/
go test -v -tags=test ./internal/db/
```

### Notes

- The `-tags=test` flag is required because the test database setup lives in `internal/db/testhelpers.go`, which is protected by a `//go:build test` build tag.
- Tests truncate all tables (`users`, `incidents`, `incident_logs`, `comments`) between test cases via `TruncateTables`.

---

## API Endpoints

### Health Check

- `GET /api/v1/ping` - Returns a pong message to verify the service is running
  - Response: `{"message": "pong"}`

### Authentication Endpoints

#### User Registration

- `POST /api/v1/auth/register` - Register a new user
  - **Requires**: Superadmin role
  - **Request Body**:
    ```json
    {
      "email": "string (required)",
      "name": "string (required)",
      "password": "string (required, min 8 characters)",
      "role": "string (required, one of: reporter, supervisor, manager, admin, superadmin)",
      "department": "string (required)"
    }
    ```
  - **Responses**:
    - `201 Created`: User successfully created
    - `400 Bad Request`: Invalid input data or invalid role
    - `403 Forbidden`: User is not a superadmin
    - `409 Conflict`: User with email already exists
    - `500 Internal Server Error`: Database or hashing error

#### User Login

- `POST /api/v1/auth/login` - Authenticate user and receive JWT token
  - **Request Body**:
    ```json
    {
      "email": "string (required)",
      "password": "string (required)"
    }
    ```
  - **Responses**:
    - `200 OK`: Authentication successful
      ```json
      {
        "token": "jwt-token-string",
        "user": {
          "id": integer,
          "name": string,
          "email": string,
          "role": string,
          "department": string,
          "disabled": boolean
        }
      }
      ```
    - `401 Unauthorized`: Invalid credentials
    - `403 Forbidden`: Account has been disabled
    - `404 Not Found`: User not found
    - `500 Internal Server Error`: Database error

#### User Management (Superadmin Only)

All user management endpoints require superadmin role and authentication middleware.

- `PUT /api/v1/auth/update` - Update user information
  - **Request Body**:
    ```json
    {
      "email": "string (required)",
      "name": "string (required)",
      "role": "string (required, one of: reporter, supervisor, manager, admin, superadmin)",
      "department": "string (required)"
    }
    ```

- `PUT /api/v1/auth/disable` - Disable a user account
  - **Request Body**:
    ```json
    {
      "email": "string (required)"
    }
    ```

- `PUT /api/v1/auth/enable` - Enable a disabled user account
  - **Request Body**:
    ```json
    {
      "email": "string (required)"
    }
    ```

- `PUT /api/v1/auth/resetpassword` - Reset a user's password
  - **Request Body**:
    ```json
    {
      "email": "string (required)",
      "password": "string (required, min 8 characters)"
    }
    ```

#### Get User

- `GET /api/v1/user` - Get user information by email
  - **Requires**: superadmin role
  - **Query Parameters**:
    - `email`: User's email address (required)
  - **Responses**:
    - `200 OK`: User information
    - `400 Bad Request`: Email parameter missing
    - `403 Forbidden`: User is not a superadmin
    - `500 Internal Server Error`: Database error

### Incident Management Endpoints

#### Report Incident

- `POST /api/v1/incidents` - Submit a new incident report
  - **Requires**: No authentication required (public endpoint)
  - **Request Body** (see full IncidentReport schema below):
    ```json
    {
      "principalName": "string (required)",
      "principalGender": "string (required)",
      "principalDob": "string (required)",
      "principalType": "string (required, one of: patient, staff, visiting consultant, other)",
      "patientId": "string (optional)",
      "patientWardDept": "string (optional)",
      "staffJobTitle": "string (optional)",
      "staffPhone": "string (optional)",
      "staffPlaceOfWork": "string (optional)",
      "staffSite": "string (optional)",
      "peopleInvolved": "string (required)",
      "dateOfIncident": "string (required)",
      "timeOfIncident": "string (required)",
      "locationOfIncident": "string (required)",
      "incidentWardDept": "string (required)",
      "witnesses": "string (optional)",
      "witnessType": "string (optional)",
      "witnessWardDept": "string (optional)",
      "witnessJobTitle": "string (optional)",
      "witenssPhone": "string (optional)",
      "isNearMiss": "boolean (required)",
      "causeGroup": "string (required)",
      "causes": "string (required)",
      "prescribingDoctor": "string (optional)",
      "treatmentReceived": "string (required)",
      "equipmentInvolved": "string (required)",
      "equipmentModel": "string (optional)",
      "equipmentSentForRepair": "boolean (required)",
      "equipmentWithdrawn": "boolean (required)",
      "equipmentRetained": "boolean (required)",
      "equipmentNumber": "string (optional)",
      "isMedicalDevice": "string (optional)",
      "reporterName": "string (required)",
      "reporterDesignation": "string (required)",
      "signature": "boolean (required)",
      "reporterInfo": "string (required)",
      "date": "string (required)",
      "severityLevel": "string (required, one of: near miss, minor, major, critical)",
      "incidentStatus": "string (optional, one of: unresolved, inprogress, resolved, defaults to unresolved)"
    }
    ```
  - **Responses**:
    - `200 OK`: Incident successfully recorded (returns saved incident with ID)
    - `400 Bad Request`: Invalid input data or invalid severity level
    - `500 Internal Server Error`: Database error

#### Get Incidents

- `GET /api/v1/incidents` - Retrieve paginated list of incidents
  - **Requires**: Authentication (any authenticated user)
  - **Role-specific behavior**:
     - `superadmin` / `admin` / `manager`: See all incidents
     - `supervisor` / `reporter`: See only incidents from their department (matched via `incident_ward_dept`, `patient_ward_dept`, or `staff_place_of_work`)
  - **Query Parameters**:
    - `page`: Page number (default: 1)
    - `limit`: Number of items per page (default: 10, max: 50)
  - **Responses**:
    - `200 OK`: Paginated response:
      ```json
      {
        "data": [ ... ],
        "pagination": {
          "current_page": 1,
          "page_size": 10,
          "total_items": 42,
          "total_pages": 5
        }
      }
      ```
    - `401 Unauthorized`: Missing or invalid authentication token
    - `500 Internal Server Error`: Database error

#### Update Incident Status

- `PATCH /api/v1/incidents/:id/status` - Update incident status
  - **Requires**: Authentication (admin/superadmin only)
  - **Path Parameters**:
    - `id`: Incident ID (required)
  - **Request Body**:
    ```json
    {
      "status": "string (required, one of: unresolved, inprogress, resolved)"
    }
    ```
  - **Role-specific behavior**:
     - `reporter` / `supervisor` / `manager`: Forbidden (403)
     - `admin` / `superadmin`: Can update any incident
  - **Responses**:
    - `200 OK`: Status updated successfully (returns updated incident)
    - `400 Bad Request`: Invalid ID or invalid status value
    - `401 Unauthorized`: Missing or invalid authentication token
    - `403 Forbidden`: Reporter, supervisor, or manager role
    - `404 Not Found`: Incident not found
    - `500 Internal Server Error`: Database error

#### Add Comment

- `POST /api/v1/incidents/comments` - Add a comment to an incident
   - **Requires**: manager or admin role
   - **Request Body**:
     ```json
     {
       "incidentId": "integer (required)",
       "userId": "integer (required)",
       "comment": "string (required)"
     }
     ```
   - **Responses**:
     - `201 Created`: Comment successfully added
     - `403 Forbidden`: User is not a manager or admin
     - `400 Bad Request`: Invalid input data
     - `500 Internal Server Error`: Database error

#### List Comments

- `GET /api/v1/incidents/comments` - Retrieve comments for an incident
   - **Requires**: admin or manager role
   - **Query Parameters**:
     - `incidentId`: Incident ID (required)
   - **Responses**:
     - `200 OK`: List of comments with commenter name and role
     - `403 Forbidden`: User is not an admin or manager

#### Submit Incident Management Report

- `POST /api/v1/incidents/:id/management` - Submit a follow-up management report for an incident
   - **Requires**: admin or manager role
  - **Path Parameters**:
    - `id`: Incident ID (required)
  - **Request Body** (see full IncidentManagement schema below):
    ```json
    {
      "impactOnService": "string (required)",
      "contributoryFactors": "string (required)",
      "actionsTakenOutcomes": "string (required)",
      "recommendations": "string (required)",
      "lessonsLearned": "string (required)",
      "informedPatient": "boolean (optional, default false)",
      "informedRelative": "boolean (optional, default false)",
      "informedSeniorManager": "boolean (optional, default false)",
      "informedPharmacist": "boolean (optional, default false)",
      "policeIncidentNumber": "string (optional)",
      "informedOther": "string (optional)",
      "riskSeverity": "integer (required)",
      "riskLikelihood": "integer (required)",
      "riskRating": "integer (required)",
      "ohsAbsenceOver3Days": "boolean (optional, default false)",
      "ohsActOfViolenceOrDanger": "boolean (optional, default false)",
      "ohsHospitalizationOver24Hours": "boolean (optional, default false)",
      "ohsStaffName": "string (optional)",
      "ohsStaffDob": "string (optional)",
      "ohsStaffAddress": "string (optional)",
      "managerName": "string (required)",
      "managerSignature": "boolean (required)",
      "managerDesignation": "string (required)",
      "managerDate": "string (required)"
    }
    ```
  - **Responses**:
    - `200 OK`: Report submitted successfully (returns saved record with ID)
    - `400 Bad Request`: Invalid input data
    - `403 Forbidden`: User is not an admin
    - `500 Internal Server Error`: Database error

#### Get Incident Management Report

- `GET /api/v1/incidents/:id/management` - Get management report for an incident
  - **Requires**: Authentication (any authenticated user)
  - **Path Parameters**:
    - `id`: Incident ID (required)
  - **Responses**:
    - `200 OK`: Management report data
    - `401 Unauthorized`: Missing or invalid authentication token
    - `404 Not Found`: Report not found
    - `500 Internal Server Error`: Database error

#### Update Incident Management Report

- `PUT /api/v1/incidents/:id/management` - Update an existing management report
  - **Requires**: manager or admin role
  - **Path Parameters**:
    - `id`: Incident ID (required)
  - **Request Body**:
    ```json
    {
      "impactOnService": "string (required)",
      "contributoryFactors": "string (required)",
      "actionsTakenOutcomes": "string (required)",
      "recommendations": "string (required)",
      "lessonsLearned": "string (required)",
      "informedPatient": "boolean (optional, default false)",
      "informedRelative": "boolean (optional, default false)",
      "informedSeniorManager": "boolean (optional, default false)",
      "informedPharmacist": "boolean (optional, default false)",
      "policeIncidentNumber": "string (optional)",
      "informedOther": "string (optional)",
      "riskSeverity": "integer (required)",
      "riskLikelihood": "integer (required)",
      "riskRating": "integer (required)",
      "ohsAbsenceOver3Days": "boolean (optional, default false)",
      "ohsActOfViolenceOrDanger": "boolean (optional, default false)",
      "ohsHospitalizationOver24Hours": "boolean (optional, default false)",
      "ohsStaffName": "string (optional)",
      "ohsStaffDob": "string (optional)",
      "ohsStaffAddress": "string (optional)",
      "managerName": "string (required)",
      "managerSignature": "boolean (required)",
      "managerDesignation": "string (required)",
      "managerDate": "string (required)"
    }
    ```
  - **Responses**:
    - `200 OK`: Report updated successfully
    - `400 Bad Request`: Invalid input data
    - `403 Forbidden`: User is not a supervisor or admin
    - `500 Internal Server Error`: Database error

#### Get Incident Management Logs

- `GET /api/v1/incidents/:id/managementlogs` - Retrieve audit logs for an incident's management reports
  - **Requires**: admin role
  - **Path Parameters**:
    - `id`: Incident ID (required)
  - **Responses**:
    - `200 OK`: List of management report audit logs with user names
    - `401 Unauthorized`: Missing or invalid authentication token
    - `403 Forbidden`: User is not an admin
    - `500 Internal Server Error`: Database error

## Role Permissions

| Role       | Permissions                                                                                                                                                         |
| ---------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| superadmin | All endpoints including user management (register, update, disable, enable, reset password, get user), report incidents, view all incidents, update incident status, submit incident management reports, update incident management reports, add comments, view comments, view incident management reports and logs |
| admin      | Report incidents, view all incidents, update incident status, submit incident management reports, update incident management reports, add comments, view comments, view incident management reports and logs |
| manager    | Report incidents, add comments, submit incident management reports, update incident management reports, view all incidents, view comments, view incident management reports                                                                                           |
| supervisor | Report incidents, view own department incidents                                                                               |
| reporter   | Report incidents via public endpoint only, view own department incidents                                                                                            |

## Database Schema

The application uses four main tables defined in `tables.sql`:

### Users Table

Stores user account information:

| Column     | Type         | Constraints                 | Description                                         |
| ---------- | ------------ | --------------------------- | --------------------------------------------------- |
| id         | SERIAL       | PRIMARY KEY                 | Auto-incrementing unique identifier                 |
| name       | VARCHAR(255) | NOT NULL                    | User's full name                                    |
| email      | VARCHAR(255) | UNIQUE NOT NULL             | User's email address (used for login)               |
| password   | VARCHAR(255) | NOT NULL                    | Bcrypt hashed password                              |
| role       | VARCHAR(50)  | NOT NULL DEFAULT 'reporter' | User role (reporter, supervisor, manager, admin, superadmin) |
| department | VARCHAR(100) | NOT NULL                    | User's department                                   |
| disabled   | BOOLEAN      | NOT NULL DEFAULT FALSE      | Account status (true = disabled)                    |

### Incidents Table

Stores comprehensive clinical incident reports (37 columns):

| Column                     | Type         | Constraints                   | Description                                        |
| -------------------------- | ------------ | ----------------------------- | -------------------------------------------------- |
| id                         | SERIAL       | PRIMARY KEY                   | Auto-incrementing unique identifier                |
| **Principal Person**       |              |                               |                                                    |
| principal_name             | VARCHAR(255) | NOT NULL                      | Name of person the incident happened to           |
| principal_gender           | VARCHAR(50)  | NOT NULL                      | Gender of principal person                         |
| principal_dob              | VARCHAR(50)  | NOT NULL                      | Date of birth of principal person                  |
| principal_type             | VARCHAR(100) | NOT NULL                      | Type: patient, staff, visiting consultant, other  |
| patient_id                 | VARCHAR(100) |                               | Patient ID (if principal is a patient)             |
| patient_ward_dept          | VARCHAR(255) |                               | Patient's ward/department                          |
| staff_job_title            | VARCHAR(255) |                               | Staff job title (if principal is staff)            |
| staff_phone                | VARCHAR(50)  |                               | Staff phone number                                 |
| staff_place_of_work        | VARCHAR(255) |                               | Staff place of work                                |
| staff_site                 | VARCHAR(255) |                               | Staff site                                         |
| **Others Involved**        |              |                               |                                                    |
| people_involved            | TEXT         | NOT NULL                      | Others directly involved in the incident           |
| **When and Where**         |              |                               |                                                    |
| date_of_incident           | VARCHAR(50)  | NOT NULL                      | Date the incident occurred                         |
| time_of_incident           | VARCHAR(50)  | NOT NULL                      | Time the incident occurred                         |
| location_of_incident       | VARCHAR(255) | NOT NULL                      | Location where incident occurred                   |
| incident_ward_dept         | VARCHAR(255) | NOT NULL                      | Ward/department (used for access control scoping)  |
| **Witnesses**              |              |                               |                                                    |
| witnesses                  | TEXT         |                               | Witness names/details                              |
| witness_type               | VARCHAR(100) |                               | Type of witnesses                                  |
| witness_ward_dept          | VARCHAR(255) |                               | Witness ward/department                            |
| witness_job_title          | VARCHAR(255) |                               | Witness job title                                  |
| witness_phone              | VARCHAR(50)  |                               | Witness phone number                               |
| **Factual Description**    |              |                               |                                                    |
| is_near_miss               | BOOLEAN      | NOT NULL DEFAULT FALSE        | Whether this was a near miss                       |
| cause_group                | VARCHAR(255) | NOT NULL                      | Cause group classification                         |
| causes                     | TEXT         | NOT NULL                      | Detailed cause description                         |
| prescribing_doctor         | VARCHAR(255) |                               | Prescribing doctor (for medication incidents)      |
| **Treatment**              |              |                               |                                                    |
| treatment_received         | VARCHAR(255) | NOT NULL                      | Treatment received                                 |
| **Equipment**              |              |                               |                                                    |
| equipment_involved         | VARCHAR(100) | NOT NULL                      | Equipment involved (string for Go alignment)       |
| equipment_model            | VARCHAR(255) |                               | Equipment model                                    |
| equipment_sent_for_repair  | BOOLEAN      | NOT NULL DEFAULT FALSE        | Whether equipment was sent for repair              |
| equipment_withdrawn        | BOOLEAN      | NOT NULL DEFAULT FALSE        | Whether equipment was withdrawn                    |
| equipment_retained         | BOOLEAN      | NOT NULL DEFAULT FALSE        | Whether equipment was retained                     |
| equipment_number           | VARCHAR(100) |                               | Equipment serial number                            |
| is_medical_device          | VARCHAR(50)  |                               | Whether it's a medical device (string for Go align)|
| **Reporter Details**       |              |                               |                                                    |
| reporter_name              | VARCHAR(255) | NOT NULL                      | Reporter's name                                    |
| reporter_designation       | VARCHAR(255) | NOT NULL                      | Reporter's designation                             |
| signature                  | BOOLEAN      | NOT NULL DEFAULT FALSE        | Whether report was signed                          |
| reporter_info              | VARCHAR(255) | NOT NULL                      | Additional reporter information                    |
| reporter_date              | VARCHAR(50)  | NOT NULL                      | Date of report                                     |
| **Status**                 |              |                               |                                                    |
| severity_level             | VARCHAR(50)  | NOT NULL                      | Severity: near miss, minor, major, critical        |
| incident_status            | VARCHAR(50)  | NOT NULL DEFAULT 'unresolved' | Status: unresolved, inprogress, resolved           |

**Index**: `idx_incidents_id_desc` on `incidents(id DESC)` for dashboard performance.

### Incident Management Table

Stores follow-up incident management data linked to incidents:

| Column | Type | Constraints | Description |
| ------ | ---- | ----------- | ----------- |
| id | SERIAL | PRIMARY KEY | Auto-incrementing unique identifier |
| incident_id | INT | UNIQUE NOT NULL REFERENCES incidents(id) ON DELETE CASCADE | Linked incident |
| impact_on_service | TEXT | NOT NULL | Impact on service description |
| contributory_factors | TEXT | NOT NULL | Contributory factors identified |
| actions_taken_outcomes | TEXT | NOT NULL | Actions taken and outcomes |
| recommendations | TEXT | NOT NULL | Recommendations made |
| lessons_learned | TEXT | NOT NULL | Lessons learned |
| informed_patient | BOOLEAN | DEFAULT FALSE | Whether patient was informed |
| informed_relative | BOOLEAN | DEFAULT FALSE | Whether relative was informed |
| informed_senior_manager | BOOLEAN | DEFAULT FALSE | Whether senior manager was informed |
| informed_pharmacist | BOOLEAN | DEFAULT FALSE | Whether pharmacist was informed |
| police_incident_number | VARCHAR(100) | | Police incident number if applicable |
| informed_other | TEXT | | Other parties informed |
| risk_severity | INT | NOT NULL | Risk severity rating |
| risk_likelihood | INT | NOT NULL | Risk likelihood rating |
| risk_rating | INT | NOT NULL | Overall risk rating |
| ohs_absence_over_3_days | BOOLEAN | DEFAULT FALSE | OHS absence over 3 days |
| ohs_act_of_violence_or_danger | BOOLEAN | DEFAULT FALSE | OHS act of violence or danger |
| ohs_hospitalisation_over_24_hours | BOOLEAN | DEFAULT FALSE | OHS hospitalisation over 24 hours |
| ohs_staff_name | VARCHAR(255) | | OHS staff name |
| ohs_staff_dob | VARCHAR(50) | | OHS staff date of birth |
| ohs_staff_address | TEXT | | OHS staff address |
| manager_name | VARCHAR(255) | NOT NULL | Manager name |
| manager_signature | BOOLEAN | NOT NULL DEFAULT FALSE | Manager signature |
| manager_designation | VARCHAR(255) | NOT NULL | Manager designation |
| manager_date | VARCHAR(50) | NOT NULL | Date of management review |

**Index**: `idx_incident_management_incident_id` on `incident_management(incident_id)`.

### Comments Table

Stores comments linked to incidents:

| Column | Type | Constraints | Description |
| ------ | ---- | ----------- | ----------- |
| id | SERIAL | PRIMARY KEY | Auto-incrementing unique identifier |
| incident_id | INT | REFERENCES incidents(id) ON DELETE CASCADE | Linked incident |
| user_id | INT | REFERENCES users(id) ON DELETE CASCADE | Comment author |
| comment | TEXT | | Comment content |

**Index**: `idx_comment` on `comments(id)`.

### Incident Logs Table

Stores audit logs for incident management report changes:

| Column | Type | Constraints | Description |
| ------ | ---- | ----------- | ----------- |
| id | SERIAL | PRIMARY KEY | Auto-incrementing unique identifier |
| incident_id | INT | REFERENCES incidents(id) ON DELETE CASCADE | Linked incident |
| action | VARCHAR(255) | NOT NULL | Action performed (created, updated, etc.) |
| changed_by | VARCHAR(255) | NOT NULL | User who made the change |
| changed_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Timestamp of change |
| details | TEXT | | Details of what changed |

## Project Structure

```
.
├── .air.toml                          # Air configuration for live reloading
├── .dockerignore                      # Docker ignore rules
├── .env                               # Local env overrides (gitignored)
├── .env.example                       # Environment variable template
├── .gitignore                         # Git ignore rules
├── AGENTS.md                          # Agent/development instructions
├── ARCHITECTURE.md                    # Architecture documentation
├── CODE_QUALITY_ASSESSMENT_AND_FIXES.md # Code quality review and improvement guide
├── CONTRIBUTING.md                    # Contribution guidelines
├── Dockerfile                         # Multi-stage Go build (builder + alpine runtime)
├── LICENSE                            # MIT License
├── README.md                          # This file
├── SYSTEM_DESIGN.md                   # System design documentation
├── docker-compose.yml                 # PostgreSQL + server service definitions
├── go.mod                             # Go module definition
├── go.sum                             # Go module checksums
├── scripts/
│   ├── login.sh                       # Script to access database shell
│   └── runtests.sh                    # Helper script to run tests
├── commit.sh                          # Helper script for git operations
├── requirements.txt                   # Feature requirements
├── tables.sql                         # Database schema, seed data, and indexes
├── update.txt                         # Schema/struct evolution notes
│
├── cmd/                               # Application entrypoint and HTTP handlers
│   ├── auth.go                        # Authentication handlers (register, login, reset password)
│   ├── auth_test.go                   # Authentication handler tests
│   ├── comments.go                    # Comment handlers (add comment to incidents)
│   ├── incidentmanagement.go          # Incident management handler (submit follow-up report)
│   ├── incidentmanagement_test.go     # Incident management handler tests
│   ├── incidents.go                   # Incident handlers (report, get, update status)
│   ├── incidents_test.go              # Incident handler tests
│   ├── main.go                        # Application initialization
│   ├── main_test.go                   # Test setup and helpers
│   ├── middleware.go                  # JWT authentication middleware
│   ├── ping_test.go                   # Health check test
│   ├── routes.go                      # API route definitions + CORS configuration
│   ├── server.go                      # HTTP server configuration (timeouts)
│   ├── types.go                       # Request/response DTOs, JWT Claims
│   ├── users.go                       # User management handlers (update, disable, enable, get user)
│   ├── users_test.go                  # User handler tests
│   └── utils.go                       # Utility functions (bcrypt password hashing)
│
├── internal/                          # Private application libraries
│   ├── db/                            # Database models and connection handling
│   │   ├── comments.go                # Comment model + CRUD queries
│   │   ├── db.go                      # Connection pool initialization, Models factory
│   │   ├── incidentmanagement.go      # Incident management model (follow-up data)
│   │   ├── incidents.go               # Incident model + CRUD queries
│   │   ├── users.go                   # User model + CRUD queries
│   │   ├── testhelpers.go             # Test database setup and utilities
│   │   └── models.go                  # Models interface definition
│   ├── env/                           # Environment variable helpers
│   │   └── env.go                     # Environment variable parsing utilities
│   └── logger/                        # Structured logging
│       └── logger.go                  # Logging setup and utilities
│
└── tmp/                               # Temporary directory (used by Air for builds)
```

## Data Flow

1. **Client Request**: HTTP request arrives at Gin router
2. **Routing**: Request directed to appropriate handler based on path and method
3. **Middleware**: CORS middleware processes all requests; JWT auth middleware validates tokens for protected routes
4. **Handler**: Application logic processes request, validates input, performs role/department authorization
5. **Database Layer**: Handler calls model methods which execute parameterized SQL queries via PGX
6. **Response**: Handler returns JSON response with appropriate status code

## Security Features

- **Password Security**: Passwords are hashed using bcrypt before storage
- **Authentication**: JWT-based authentication with 72-hour expiration
- **Authorization**: Role-based access control enforced via middleware; department-based data scoping
- **Input Validation**: All incoming data is validated using Gin's binding mechanism
- **SQL Injection Prevention**: Uses parameterized queries via PGX
- **CORS**: Configured via gin-contrib/cors middleware with allowed origins from environment
- **Disabled Accounts**: Login is rejected for disabled user accounts

## Error Handling

The API follows consistent error response format:

```json
{
  "error": "Descriptive error message"
}
```

For paginated endpoints, data and pagination metadata are returned at the top level.

HTTP status codes are used appropriately:

- 2xx: Success
- 4xx: Client errors (validation, authentication, authorization, not found, conflict)
- 5xx: Server errors (database issues, etc.)

## Extending the Application

### Adding New Endpoints

1. Define handler function in appropriate file under `cmd/`
2. Add route to `routes.go` within the appropriate group
3. Apply middleware as needed (`authMiddleware()` for protected routes)
4. Update corresponding model in `internal/db/` if database changes needed

### Adding New Database Tables/Columns

1. Modify `tables.sql` with schema changes
2. Update corresponding model in `internal/db/` (struct + queries)
3. Update handler functions in `cmd/` to use new fields
4. Rebuild and restart containers to apply changes

### Configuration Changes

- Update `.env` file for environment-specific settings
- Modify `.air.toml` for Air live reload configuration
- Adjust Docker Compose file for service changes

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for contribution guidelines, coding standards, and development workflow.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Gin-Gonic team for the excellent web framework
- PostgreSQL team for the reliable database
- JWT and bcrypt libraries for secure authentication
- Air team for the live reload development experience

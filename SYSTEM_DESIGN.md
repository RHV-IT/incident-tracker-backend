# Issue Tracker - System Design

**Code Metrics:** ~1800 lines of Go, 20 source files

## System Overview

The Issue Tracker is a stateless RESTful API built with Go that provides incident tracking capabilities with role-based access control. The system follows a layered architecture pattern with clear separation between presentation, application, and data layers.

Go version: 1.22+

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              CLIENT LAYER                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐          │
│  │   Web Client    │    │  Mobile App     │    │  API Consumer   │          │
│  └────────┬────────┘    └────────┬────────┘    └────────┬────────┘          │
│           │                      │                      │                    │
└           │                      │                      │                    │
            ▼                      ▼                      ▼                    │
└─────────────────────────────────────────────────────────────────────────────┘
                                     │
┌─────────────────────────────────────────────────────────────────────────────┐
│                           LOAD BALANCER                                     │
│                      (Optional for production)                               │
└─────────────────────────────────────────────────────────────────────────────┘
                                     │
┌─────────────────────────────────────────────────────────────────────────────┐
│                              API LAYER                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                        Gin Web Framework                               │   │
│  │  ┌─────────────────────────────────────────────────────────────┐    │   │
│  │  │                     Router & Middleware                      │    │   │
│  │  │  ┌─────────────┐  ┌─────────────┐  ┌──────────────────┐  │    │   │
│  │  │  │ CORS        │  │ JWT Auth    │  │ Rate Limiting    │  │    │   │
│  │  │  └─────────────┘  └─────────────┘  └──────────────────┘  │    │   │
│  │  └─────────────────────────────────────────────────────────────┘    │   │
│  │                              │                                        │   │
│  │  ┌─────────────────────────────────────────────────────────────┐    │   │
│  │  │                     Route Groups                            │    │   │
│  │  │  /api/v1/ping          → Health Check Handler             │    │   │
│  │  │  /api/v1/auth/*        → Auth Handlers (register, login, reset pwd)  │    │   │
  │  │  │  /api/v1/incidents     → Incident Handlers (public report, auth list/update)  │    │   │
  │  │  │  /api/v1/user          → User Handlers (get user by email)  │    │   │
│  │  └─────────────────────────────────────────────────────────────┘    │   │
│  │                              │                                        │   │
│  │  ┌─────────────────────────────────────────────────────────────┐    │   │
│  │  │                     Handlers                                │    │   │
│  │  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐       │    │   │
│  │  │  │ auth.go     │  │ incidents.go│  │ users.go    │       │    │   │
  │  │  │  │ - register  │  │ - report    │  │ - update    │       │    │   │
  │  │  │  │ - login     │  │ - get       │  │ - disable   │       │    │   │
  │  │  │  │ - resetpwd  │  │ - updateStatus│ - enable  │       │    │   │
  │  │  │  └─────────────┘  └─────────────┘  │ - get user  │       │    │   │
  │  │  │                                     └─────────────┘       │    │   │
│  │  │  ┌─────────────┐                                         │    │   │
│  │  │  │ utils.go    │                                         │    │   │
│  │  │  │ - hashPass  │                                         │    │   │
│  │  │  │ - verifyPass│                                         │    │   │
│  │  │  └─────────────┘                                         │    │   │
│  │  └─────────────────────────────────────────────────────────────┘    │   │
│  │                              │                                        │   │
│  │                              ▼                                        │   │
│  │  ┌─────────────────────────────────────────────────────────────┐    │   │
│  │  │                   Application Layer                         │    │   │
│  │  │  ┌─────────────────────────────────────────────────────┐  │    │   │
│  │  │  │  type application struct {                          │  │    │   │
│  │  │  │    port int                                         │  │    │   │
│  │  │  │    jwtsecret string                                 │  │    │   │
│  │  │  │    db *pgxpool.Pool                                 │  │    │   │
│  │  │  │    models db.Models                                   │  │    │   │
│  │  │  │    origins string                                     │  │    │   │
│  │  │  │  }                                                   │  │    │   │
│  │  │  └─────────────────────────────────────────────────────┘  │    │   │
│  │  └─────────────────────────────────────────────────────────────┘    │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                              │                                              │
│                              ▼                                              │
└─────────────────────────────────────────────────────────────────────────────┘
                                     │
┌─────────────────────────────────────────────────────────────────────────────┐
│                           DATA LAYER                                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                       Models                                        │   │
│  │  ┌─────────────────────────────────────────────────────────────┐    │   │
│  │  │  db.Models {                                               │    │   │
│  │  │    Users UserModel                                          │    │   │
│  │  │    Incidents IncidentsModel                                 │    │   │
│  │  │    IncidentManagement IncidentManagementModel               │    │   │
│  │  │  }                                                          │    │   │
│  │  └─────────────────────────────────────────────────────────────┘    │   │
│  │                              │                                        │   │
│  │  ┌─────────────────┐    ┌─────────────────┐                        │   │
│  │  │   users.go      │    │  incidents.go   │                        │   │
│  │  │  - GetByEmail   │    │  - Insert       │                        │   │
│  │  │  - Insert       │    │  - FetchIncidents                        │   │
│  │  │  - Update       │    │  - FetchBySupervisor                     │   │
  │  │  │  - DisableUser  │    │  - FetchById    │                        │   │
  │  │  │  - EnableUser   │    │  - UpdateIncidentStatus                  │   │
  │  │  │  - ResetPassword│    │                 │                        │   │
  │  │  └─────────────────┘    └─────────────────┘                        │   │
│  │                                                                     │   │
│  │  
│  │  │  ┌─────────────────────────────────────────────────────────────┐    │   │
│  │  │  │  incidentmanagement.go                                     │    │   │
│  │  │  │  - IncidentManagementModel (follow-up data access)         │    │   │
│  │  │  │  └─────────────────────────────────────────────────────────────┘    │   │
│  │  │                                                                     │   │┌─────────────────────────────────────────────────────────────┐    │   │
│  │  │  db.go                                                      │    │   │
│  │  │  - InitPool()  → Creates PGX connection pool              │    │   │
│  │  │  - NewModels() → Factory for model instances               │    │   │
│  │  └─────────────────────────────────────────────────────────────┘    │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                              │                                              │
│                              ▼                                              │
└─────────────────────────────────────────────────────────────────────────────┘
                                     │
┌─────────────────────────────────────────────────────────────────────────────┐
│                        INFRASTRUCTURE LAYER                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐          │
│  │ PostgreSQL 16   │    │ Environment     │    │ Docker          │          │
│  │ - users table   │    │ Variables       │    │ Compose         │          │
│  │ - incidents     │    │ - dbConnStr     │    │ - postgres:16-  │          │
│  │ - indexes       │    │ - jwtSecret     │    │   alpine        │          │
│  │ - constraints   │    │ - PORT          │    │ - Air (hot     │          │
│  │                 │    │ - allowedOrig.  │    │   reload)       │          │
│  └─────────────────┘    └─────────────────┘    └─────────────────┘          │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
  │  │  Scripts                                                             │   │
  │  │  - commit.sh        → Git helper                                    │   │
  │  │  - login.sh         → psql shell into DB container                  │   │
  │  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  Database Initialization                                             │   │
│  │  - tables.sql → Auto-run via Docker initdb mechanism               │   │
│  │  - Schema created on first container start                        │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Component Interaction Flow

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│   Client     │────▶│     API      │────▶│  Database    │
│              │     │              │     │              │
└──────────────┘     └──────────────┘     └──────────────┘
                           │
                           ▼
                    ┌──────────────┐
                    │  Middleware  │
                    │  (JWT Auth)  │
                    └──────────────┘
                           │
                           ▼
                    ┌──────────────┐
                    │   Handler    │
                    │  (Business  │
                    │   Logic)     │
                    └──────────────┘
                           │
                           ▼
                    ┌──────────────┐
                    │    Model     │
                    │ (Data Access)│
                    └──────────────┘
```

## Data Flow Sequence

```
Client Request
      │
      ▼
┌─────────────────┐
│ HTTP Request    │
│ POST /api/v1/auth/login
│ Body: {email, password}
└─────────────────┘
      │
      ▼
┌─────────────────┐
│ Gin Router      │
│ matches route   │
└─────────────────┘
      │
      ▼
┌─────────────────┐
│ Middleware      │
│ (none for login)│
└─────────────────┘
      │
      ▼
┌─────────────────┐
│ Handler         │
│ login()         │
│ 1. Validate     │
│ 2. GetByEmail() │
│ 3. Verify pass  │
│ 4. Create JWT   │
└─────────────────┘
      │
      ▼
┌─────────────────┐
│ Model Layer     │
│ UserModel       │
│ Query: SELECT   │
│ WHERE email=$1  │
└─────────────────┘
      │
      ▼
┌─────────────────┐
│ PostgreSQL      │
│ Execute query   │
└─────────────────┘
      │
      ▼
┌─────────────────┐
│ Response        │
│ {token, user}   │
└─────────────────┘
```

## Request-Response Flow

```
                    ┌─────────────────────────────────────────┐
                    │              REQUEST                  │
                    └──────────────────┬────────────────────┘
                                       │
                                       ▼
┌─────────────────────────────────────────────────────────────────┐
│                         HTTP SERVER                              │
│  ┌──────────────────────────────────────────────────────────┐    │
│  │                     Gin Engine                            │    │
│  │                                                          │    │
│  │  1. CORS Middleware                                      │    │
│  │  2. Route Matching                                       │    │
│  │  3. Auth Middleware (if protected)                       │    │
│  │     - Extract Bearer token                                │    │
│  │     - Validate JWT signature                              │    │
│  │     - Verify expiration                                   │    │
│  │     - Set user context (userId, role, email, dept)      │    │
│  │                                                          │    │
│  │  4. Handler Execution                                    │    │
│  │     - Input validation                                   │    │
│  │     - Role-based authorization                           │    │
│  │     - Business logic                                     │    │
│  │                                                          │    │
│  │  5. Response Serialization                               │    │
│  └──────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────┘
                                       │
                                       ▼
                    ┌─────────────────────────────────────────┐
                    │            DATABASE LAYER               │
                    │                                         │
                    │  ┌────────────────────────────────────┐   │
                    │  │      Connection Pool              │   │
                    │  │  (Min: 2, Max: 10 connections)    │   │
                    │  └────────────────────────────────────┘   │
                    │                   │                       │
                    │                   ▼                       │
                    │  ┌────────────────────────────────────┐   │
                    │  │         PGX Driver                 │   │
                    │  │  - Parameterized queries           │   │
                    │  │  - Connection pooling              │   │
                    │  │  - Row scanning to structs         │   │
                    │  └────────────────────────────────────┘   │
                    │                                         │
                    └─────────────────────────────────────────┘
                                       │
                                       ▼
                    ┌─────────────────────────────────────────┐
                    │            RESPONSE                     │
                    └─────────────────────────────────────────┘
```

## System Components

### 1. Presentation Layer (`cmd/`)
- **Routes**: HTTP endpoint definitions with middleware chains
- **Handlers**: Request processing and response formatting
- **Types**: Request/response DTOs and domain types

### 2. Application Layer
- **Business Logic**: Implemented in handlers
- **Authorization**: Role-based access control
- **Validation**: Input validation using Gin binding

### 3. Data Access Layer (`internal/db/`)
- **Models**: Database interaction logic
- **Connection Pool**: PGX connection management
- **Queries**: Parameterized SQL operations

### 4. Infrastructure Layer
- **Database**: PostgreSQL with connection pooling
- **Configuration**: Environment variables
- **Deployment**: Docker Compose (PostgreSQL + Server containers)
- **Scripts**:
  - `login.sh` → Access DB shell
  - `commit.sh` → Git helper

## Security Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                      SECURITY FLOW                              │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐         │
│  │   Client    │───▶│  Password   │───▶│  Bcrypt     │         │
│  │             │    │  (plain)    │    │  Hash       │         │
│  └─────────────┘    └─────────────┘    └──────┬──────┘         │
│                                                │                │
│                                                ▼                │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐         │
│  │   Login     │───▶│  JWT        │───▶│  Signed     │         │
│  │  Request    │    │  Claims     │    │  Token      │         │
│  └─────────────┘    └─────────────┘    └──────┬──────┘         │
│                                                │                │
│                                                ▼                │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐         │
│  │  Protected  │───▶│  Extract    │───▶│  Validate   │         │
│  │  Endpoint   │    │  Token      │    │  Signature  │         │
│  └─────────────┘    └─────────────┘    └──────┬──────┘         │
│                                                │                │
│                                                ▼                │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐         │
│  │   Access    │◀───│  Check      │◀───│  Verify     │         │
│  │  Granted?   │    │  Role       │    │  Claims     │         │
│  └─────────────┘    └─────────────┘    └─────────────┘         │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

## Role Hierarchy

| Role | Permissions |
|------|-------------|
| **superadmin** | User management (register, update, disable/enable, reset password, get user), report incidents, view all incidents, update any incident status, add comments, view comments, submit and update incident management reports, view incident management reports and logs |
| **admin** | Report incidents, view all incidents, update any incident status, add comments, view comments, submit and update incident management reports, view incident management reports and logs |
| **supervisor** | Report incidents, view own department incidents (via `incident_ward_dept`) |
| **manager** | Report incidents, add comments, submit incident management reports, update incident management reports, view all incidents, view comments, view incident management reports |
| **reporter** | Report incidents via public endpoint only, view own department incidents |

## Deployment Architecture

### Development
```
┌─────────────────────────────────────────────────────────┐
│                  DEVELOPMENT ENVIRONMENT                 │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  Local Machine                                         │
│  ┌─────────────────────────────────────────────────┐   │
│  │  Docker Compose                                 │   │
│  │  ┌─────────────────────────────────────────┐    │   │
│  │  │  Go Application (built from Dockerfile) │   │
│  │  │  - Port: 3002                           │   │
│  │  │  - Hot reload via Air                   │   │
│  │  │  - Volume: .:/app (code mounting)       │   │
│  │  │  - Excludes: scripts/ directory         │   │
│  │  └─────────────────────────────────────────┘   │   │
│  └─────────────────────────────────────────────────┘   │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

### Production
```
┌─────────────────────────────────────────────────────────┐
│                   PRODUCTION ARCHITECTURE               │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  Internet Traffic                                       │
│         │                                                │
│         ▼                                                │
│  ┌─────────────────┐                                      │
│  │  Load Balancer  │                                      │
│  │  (NGINX/HAProxy)│                                      │
│  └────────┬────────┘                                      │
│           │                                                │
│           ▼                                                │
│  ┌─────────────────┐    ┌─────────────────┐               │
│  │  API Instance 1 │    │  API Instance 2 │               │
│  │  (Go + Gin)     │    │  (Go + Gin)     │               │
│  └────────┬────────┘    └────────┬────────┘               │
│           │                      │                          │
│           └───────────┬──────────┘                          │
│                       │                                     │
│                       ▼                                     │
│  ┌────────────────────────────────────────────────────┐    │
│  │               PostgreSQL Cluster                    │    │
│  │  ┌─────────────┐    ┌─────────────┐    ┌───────┐   │    │
│  │  │  Primary    │    │  Replica 1  │    │  ...  │   │    │
│  │  │  (RW)       │    │  (RO)       │    │       │   │    │
│  │  └─────────────┘    └─────────────┘    └───────┘   │    │
│  └────────────────────────────────────────────────────┘    │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

## Performance Characteristics

| Metric | Value |
|--------|-------|
| Max Connections | 10 (configurable) |
| JWT Expiration | 72 hours |
| Request Timeout | 1s read, 5s write |
| Idle Timeout | 30s |
| Pagination Limit | Max 50 items |

## Current Implementation Status

**Implemented:**
- ✅ Clean layered architecture (presentation → application → data → infrastructure)
- ✅ JWT authentication with 72-hour expiry
- ✅ Role-based access control (superadmin, admin, supervisor, manager, reporter)
- ✅ Department-scoped incident access
- ✅ Incident management follow-up reports
- ✅ Incident comments
- ✅ Structured logging (`internal/logger`)
- ✅ Health check endpoint (`/api/v1/ping`)
- ✅ CORS configuration
- ✅ Database connection pooling
- ✅ Docker Compose setup
- ✅ Unit tests for routes and handlers

**Pending:**
- ❌ Unit tests (increase coverage beyond current partial implementation)
- ❌ Rate limiting
- ❌ Error handling package
- ❌ Configuration validation

## Scalability Considerations

1. **Horizontal Scaling**: API instances can be scaled behind a load balancer
2. **Database Scaling**: Connection pooling, read replicas for reporting
3. **Caching**: Redis layer can be added for frequently accessed data
4. **Stateless**: JWT enables horizontal scaling without session storage
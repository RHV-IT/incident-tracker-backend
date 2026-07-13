# Issue Tracker Architecture

**Code Metrics:** ~1800 lines of Go, 20 source files

Go version: 1.22+

## System Overview

The Issue Tracker is a layered RESTful API built with Go that follows a clean separation of concerns. The system is designed to be maintainable, testable, and scalable while providing robust incident tracking capabilities.

## Architectural Layers

### 1. Presentation Layer (HTTP Handlers)
**Location**: `/cmd/` directory
**Purpose**: Handles HTTP requests, validates input, and formats responses

Components:
- **Route Definitions** (`routes.go`): Maps HTTP endpoints to handler functions, configures CORS
- **HTTP Handlers**: Process requests and return responses:
   - `auth.go`: Authentication endpoints (register, login, reset password)
   - `comments.go`: Comment handlers (add comment to incidents)
   - `users.go`: User management endpoints (update, disable, enable, get user)
   - `incidents.go`: Incident handlers (public report, authenticated list with dept scoping, status update)
   - `main.go`: Application entrypoint and initialization
   - `server.go`: HTTP server configuration (timeouts)
   - `types.go`: Request/response DTOs, JWT Claims, pagination types
   - `utils.go`: Utility functions (bcrypt password hashing)
   - `middleware.go`: JWT authentication middleware

### 2. Application Layer
**Location**: Implicit in handler functions
**Purpose**: Contains business logic and orchestrates operations between layers

The application layer is implemented directly within the handler functions, which:
- Validate incoming requests
- Apply business rules (role checking, data validation)
- Coordinate between presentation and data layers
- Format responses for clients

### 3. Data Access Layer
**Location**: `/internal/db/` directory
**Purpose**: Handles all database interactions and data modeling

Components:
- **Database Connection** (`db.go`): Manages PostgreSQL connection pool
- **Data Models**:
   - `users.go`: User model with CRUD operations
   - `comments.go`: Comment model with CRUD operations
   - `incidents.go`: Incident model with CRUD operations
   - `incidentmanagement.go`: Incident management model (follow-up actions, risk assessment, OHS details)
- **Model Factory** (`db.go`): `NewModels()` function creates model instances

### 4. Infrastructure Layer
**Location**: Various supporting files
**Purpose**: Provides foundational services and configuration

Components:
- **Environment Configuration** (`internal/env/env.go`): Typed accessors (string, int) with fallback defaults
- **Configuration Files**:
  - `.air.toml`: Live reload configuration for development
  - `docker-compose.yml`: Container orchestration (PostgreSQL + server)
  - `tables.sql`: Database schema definition, seed data, and indexes
  - `.env.example`: Environment variable template
- **Scripts**:
  - `commit.sh`: Git operations helper
  - `login.sh`: Database access helper (psql shell into container)

## Component Responsibilities

### HTTP Handlers (`/cmd/`)
- Parse HTTP requests and extract data
- Validate input using Gin's binding and custom validation
- Apply authentication and authorization checks
- Call appropriate service/data methods
- Serialize responses to JSON
- Handle errors and return appropriate HTTP status codes

### Middleware
The application uses two middlewares:
- **CORS Middleware**: Configured in `routes.go`, sets up Cross-Origin Resource Sharing.
- **JWT Middleware** (`middleware.go`): Extracts and validates JWT tokens from Authorization headers, verifies token signature and expiration, extracts user claims (ID, role, email, department), sets user context in Gin context for handlers to access, and returns 401 Unauthorized for invalid/missing tokens.

### Database Models (`/internal/db/`)
- Define structs that map to database tables
- Implement CRUD operations using PGX
- Handle database connections and transactions
- Convert between database rows and Go structs
- Handle database-specific errors

The models are:
- `UserModel` — CRUD on users
- `IncidentsModel` — CRUD on incidents
- `IncidentManagementModel` — access to incident management follow-up records
- `CommentModel` — CRUD on incident comments

### Environment Handling (`internal/env/`)
- Load environment variables with fallback defaults
- Provide typed accessors (string, int)
- Centralize configuration access

## Data Flow

### Request Processing Flow
```
HTTP Request → Gin Router → Middleware (if applicable) → Handler → Validation → Model → Database → Model → Handler → Response → Client
```

### Detailed Authentication Flow
1. Client sends request with `Authorization: Bearer <token>` header
2. Middleware extracts token and validates signature using jwtSecret
3. Middleware verifies token expiration and extracts claims
4. Middleware sets user context (ID, role, email, department) in Gin context
5. Handler accesses user context via `c.GetString("userRole")` etc.
6. Handler performs role-based authorization checks
7. Handler processes request and returns response

### Database Operation Flow
1. Handler calls model method (e.g., `a.models.Users.GetByEmail`)
2. Model executes parameterized SQL query using PGX
3. PGX prevents SQL injection through parameter binding
4. Model scans results into Go struct
5. Model returns struct and error (if any) to handler
6. Handler processes result and returns appropriate response

## Key Design Decisions

### 1. Layered Architecture
- Separation of concerns between HTTP handling, business logic, and data access
- Each layer has distinct responsibilities
- Layers communicate through well-defined interfaces
- Enhances testability and maintainability

### 2. Dependency Injection
- Models are injected into the application struct
- Database connection is shared across models
- Application struct holds all dependencies
- Facilitates easy testing with mock dependencies

### 3. JWT-Based Authentication
- Stateless authentication suitable for APIs
- Token contains user identity and role information
- Short expiration time (72 hours) enhances security
- Secret key stored in environment variable

### 4. Role-Based Access Control (RBAC)
- Five distinct roles with permissions
- Role checking performed in handlers after authentication
- Superadmin has highest privileges (user management)
- Incident reporting endpoint is public (no authentication required)
- Incident listing requires authentication

### 5. Environment-Based Configuration
- Configuration via environment variables
- Defaults provided for development
- `.env` file for local overrides
- Separation of config from code

### 6. Docker-First Development
- Docker Compose for consistent development environment
- Containerized PostgreSQL ensures consistency
- Scripts simplify common operations
- Easy reproduction of production-like environment

## Current Implementation Status

**Code Metrics:**
- Total Go code: ~1800 lines
- 20 Go source files
- Logger already implemented in `internal/logger/logger.go`
- Test helpers implemented in `internal/db/testhelpers.go`

**Implemented Features:**
- User authentication with JWT (72-hour expiry)
- Role-based access control (superadmin, admin, supervisor, manager, reporter)
- Department-scoped incident access
- Incident management follow-up reports
- Incident comments
- Health check endpoint
- CORS configuration
- Unit tests for routes and handlers
- Database connection pooling

## Security Considerations

### Authentication Security
- JWT tokens signed with HMAC-SHA256 using secret key
- Tokens expire after 72 hours
- Secure password hashing with bcrypt (default cost)
- HTTPS recommended for production (not implemented in dev)

### Authorization Security
- Role-based endpoint protection
- Superadmin-only endpoints for user management
- Incident reporting endpoint is public (no authentication required)
- Incident listing requires authentication; supervisors and reporters are scoped to their department
- Incident status update blocked for reporters, supervisors, and managers; only admin and superadmin can update any incident status
- Principle of least privilege applied

### Data Security
- Parameterized queries prevent SQL injection
- Input validation on all endpoints
- Passwords never stored or logged in plaintext
- Sensitive fields omitted from JSON responses (`password` field has `json:"-"`)

### Infrastructure Security
- Database credentials in environment variables
- Container isolation via Docker
- Non-root user in PostgreSQL container (default)

## Scalability Considerations

### Horizontal Scaling
- Stateless API instances can be behind load balancer
- Shared PostgreSQL database requires connection pooling
- Redis could be added for distributed caching/sessions
- File uploads would require shared storage solution

### Vertical Scaling
- Database connection pooling (PGX pool with Min/Max conns)
- Efficient JSON serialization
- Minimal memory footprint per request
- Optimized database queries

### Caching Opportunities
- User profile data could be cached
- Frequently accessed reference data
- Rate limiting could be implemented
- Response caching for public endpoints

## Deployment Architecture

### Development Environment
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

### Production Environment
```
Load Balancer (optional)
     ↓
API Instances (multiple for HA)
     ↓
PostgreSQL Database (primary)
     ↓
PostgreSQL Replicas (for read scaling)
     ↓
Backup Storage
     ↓
Monitoring/Logging Services
```

### Key Production Considerations
1. **Database**: 
   - Connection pooling tuned for expected load
   - Regular backups and point-in-time recovery
   - Monitoring for query performance and connection usage
   - SSL/TLS encryption for database connections

2. **Application**:
   - Proper logging and error tracking
   - Health check endpoints for load balancer
   - Resource limits and auto-scaling policies
   - Security scanning and vulnerability management
   - TLS termination at load balancer or ingress

3. **Infrastructure**:
   - Container orchestration (Kubernetes, ECS, etc.)
   - Service discovery and configuration management
   - Secret management for database credentials and JWT secret
   - CI/CD pipeline for automated testing and deployment

## Technology Choices Justification

### Go Language
- **Performance**: Compiled language with efficient memory management
- **Concurrency**: Excellent support for handling multiple requests
- **Deployment**: Single binary deployment simplifies operations
- **Ecosystem**: Strong standard library and growing package ecosystem
- **Maintainability**: Strict typing and clear syntax reduce bugs

### Gin Web Framework
- **Performance**: High performance with minimal overhead
- **Middleware**: Flexible middleware system for cross-cutting concerns
- **Routing**: Powerful routing capabilities with grouping
- **Validation**: Built-in JSON validation and binding
- **Community**: Popular framework with good documentation

### PostgreSQL Database
- **Reliability**: ACID compliance and proven track record
- **Features**: Rich data types, full-text search, JSON support
- **Performance**: Excellent performance with proper indexing
- **Scalability**: Vertical scaling capabilities and read replicas
- **Open Source**: Strong community and enterprise support

### JWT Authentication
- **Stateless**: No server-side session storage required
- **Standard**: Widely adopted standard (RFC 7519)
- **Flexible**: Can carry arbitrary claims in token
- **Interoperable**: Supported across many platforms and languages
- **Secure**: Cryptographically signed when implemented correctly

### PGX Driver
- **Performance**: High-performance PostgreSQL driver for Go
- **Features**: Native support for PostgreSQL features
- **Standards**: Follows database/sql interface
- **Maintenance**: Actively maintained with good documentation
- **Compatibility**: Works with connection pooling and contexts

## Current Implementation Status

## Future Enhancements

### Short Term
1. **Password Reset**: Add forgot password and reset functionality
2. **Email Verification**: Verify user emails during registration
3. **API Documentation**: Add Swagger/OpenAPI documentation
4. **Rate Limiting**: Prevent abuse of endpoints
5. **Request Logging**: Structured logging for debugging and monitoring
6. **Health Checks**: Add readiness and liveness probes

### Medium Term
1. **Role Management**: Dynamic role and permission management
2. **Audit Logging**: Track all changes for compliance
3. **File Attachments**: Allow attaching files to incident reports
4. **Search & Filtering**: Advanced search capabilities for incidents
5. **Analytics Dashboard**: Basic statistics and reporting
6. **Internationalization**: Support for multiple languages

### Long Term
1. **Microservices**: Split into separate services (auth, incidents, users)
2. **Event Streaming**: Use Kafka/RabbitMQ for inter-service communication
3. **Caching Layer**: Add Redis for frequently accessed data
4. **Monitoring**: Comprehensive metrics, tracing, and alerting
5. **Mobile App**: Native mobile applications for incident reporting
6. **Workflow Engine**: Customizable incident resolution workflows

## Diagrams

### Component Interaction Diagram
```
┌─────────────────┐    ┌──────────────────┐    ┌──────────────────┐
│   Client App    │◄──►│   API Gateway    │◄──►│   Load Balancer  │
└─────────────────┘    └──────────────────┘    └──────────────────┘
                              │                         │
                              ▼                         ▼
                    ┌──────────────────┐        ┌──────────────────┐
                    │   API Instance   │        │   API Instance   │
                    └──────────────────┘        └──────────────────┘
                              │                         │
                              ▼                         ▼
                    ┌──────────────────┐        ┌──────────────────┐
                    │  Gin Router      │        │  Gin Router      │
                    └──────────────────┘        └──────────────────┘
                              │                         │
              ┌───────────────┴─────────────┐ ┌───────────────┴─────────────┐
              ▼                             ▼ ▼                             ▼
    ┌──────────────────┐        ┌──────────────────┐  ┌──────────────────┐
    │ Auth Middleware  │        │ Validation Layer │  │   Handler Logic  │
    └──────────────────┘        └──────────────────┘  └──────────────────┘
              │                         │           │
              ▼                         ▼           ▼
    ┌──────────────────┐        ┌──────────────────┐  ┌──────────────────┐
    │   User Models    │◄───────┤   Incident Models│  │   Incident Mgmt  │
    └──────────────────┘        └──────────────────┘  └──────────────────┘
              │                         │           │
              └─────────────┬───────────┘└──────────┬─────────────┘
                            ▼                       ▼
                    ┌─────────────────────────────────────┐
                    │   PostgreSQL Database Connection    │
                    └─────────────────────────────────────┘
                                      │
                                      ▼
                            ┌──────────────────┐
                            │  PostgreSQL DB   │
                            └──────────────────┘
```

### Data Flow Diagram
```
HTTP Request
     ↓
[Route Matching] → [Middleware Auth] → [Handler Function]
     ↓                     ↓                    ↓
[Input Validation] ← [Context User Data]    [Business Rules]
     ↓                     ↓                    ↓
[Call Model Method] → [Parameterized Query] → [Database]
     ↓                     ↓                    ↓
[Process Results] ← [Scan to Struct]     ← [Query Results]
     ↓                     ↓                    ↓
[Format Response] ← [Error Handling]     ← [Database Errors]
     ↓
HTTP Response
```

## Conclusion

This architecture provides a solid foundation for an incident tracking system that is:
- **Maintainable**: Clear separation of concerns and consistent patterns
- **Secure**: Proper authentication, authorization, and data protection
- **Testable**: Loose coupling allows for unit and integration testing
- **Extensible**: Well-defined interfaces make adding features straightforward
- **Operational**: Docker-based deployment and helpful scripts simplify management

The layered approach ensures that changes in one area (e.g., database schema) have minimal impact on other areas (e.g., HTTP handling), making the system easier to evolve over time.

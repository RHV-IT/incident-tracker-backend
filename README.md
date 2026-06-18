# Issue Tracker

A RESTful API for tracking incidents/issues built with Go, Gin, and PostgreSQL.

## Overview

The Issue Tracker is a web application designed to help organizations track and manage workplace incidents, safety reports, and other types of issues. It provides user authentication, role-based access control, and incident reporting capabilities.

## Features

- **User Authentication**
  - User registration (requires superadmin role)
  - Secure login with JWT tokens
  - User profile updates (superadmin only)
  - Account disable/enable functionality (superadmin only)
  
- **Incident Management**
  - Report new incidents with detailed information
  - Track incident severity levels (Near Miss, Minor, Major, Critical)
  - Record incident details including location, time, description, and actions taken

- **Role-Based Access Control**
  - Four distinct roles: Reporter, Supervisor, Admin, Superadmin
  - Role-based endpoint protection
  - Superadmin privileges for user management

- **Development Experience**
  - Docker Compose setup for easy development
  - Hot reconfiguration with Air for live reloading
  - Comprehensive environment variable configuration
  - Helper scripts for common operations

## Technology Stack

- **Language**: Go 1.26.3
- **Web Framework**: Gin-Gonic
- **Database**: PostgreSQL with PGX driver
- **Authentication**: JWT (JSON Web Tokens) with bcrypt password hashing
- **Development Tools**: Air (live reload), Docker Compose
- **API Testing**: Built-in endpoint testing capabilities
- **CORS**: gin-contrib/cors middleware for Cross-Origin Resource Sharing

## Getting Started

### Prerequisites

- Docker and Docker Compose (for containerized setup)
- Go 1.26.3 (if running locally without Docker)
- Git (for version control)

### Running with Docker Compose (Recommended)

1. Clone the repository
   ```bash
   git clone <repository-url>
   cd issueTracking
   ```

2. Configure environment variables (optional, defaults are provided):
   ```bash
   # Or set environment variables directly
   export jwtSecret=your-secret-key
   export allowedOrigins="http://localhost:3000"
   ```

3. Start all services (PostgreSQL and server)
   ```bash
   docker compose up -d
   ```

4. The server will be available at `http://localhost:3002`
   - Database tables are automatically created on first run
   - Server waits for PostgreSQL to be healthy before starting

### Running Locally

```bash
# Install dependencies
go mod download

# Set environment variables (optional)
export PORT=3001
export jwtSecret=your-secret-key
export allowedOrigins="http://localhost:3000"
export dbConnStr="postgres://tracker_user:tracker_password@localhost:5432/issuetracker"

# Run with live reload
air

# Or run directly
go run ./cmd/main.go
```

### Environment Variables

The following environment variables are used:

| Variable | Description | Default Value |
|----------|-------------|---------------|
| `PORT` | The port on which the server runs | `3001` |
| `dbConnStr` | PostgreSQL connection string (use `postgres` hostname in Docker) | `postgres://tracker_user:tracker_password@postgres:5432/issuetracker` |
| `jwtSecret` | Secret key for JWT token signing | `someSecret` |
| `allowedOrigins` | Comma-separated list of allowed origins for CORS | `http://localhost:3000,http://192.168.9.227:3000` |

These can be set in the `docker-compose.yml` environment section, exported in the shell, or via an `.env` file.

> **Note**: For production, change the default `jwtSecret` to a strong, unique value.

### Helper Scripts

- `./login.sh` - Opens an interactive psql shell to the PostgreSQL container for direct database access
- `./commit.sh` - Helper script that stages all changes, prompts for a commit message, commits, and pushes to remote

## API Endpoints

All API endpoints are prefixed with `/api/v1`.

### Health Check

- `GET /api/v1/ping` - Returns a pong message to verify the service is running
  - Response: `{"message": "pong"}`

### Authentication Endpoints

All authentication endpoints require appropriate roles and are protected by authentication middleware where noted.

#### User Registration
- `POST /api/v1/auth/register` - Register a new user
  - **Requires**: Superadmin role
  - **Request Body**:
    ```json
    {
      "email": "string (required)",
      "name": "string (required)",
      "password": "string (required, min 8 characters)",
      "role": "string (required, one of: reporter, supervisor, admin, superadmin)",
      "department": "string (required)"
    }
    ```
  - **Responses**:
    - `201 Created`: User successfully created
    - `400 Bad Request`: Invalid input data
    - `401 Unauthorized`: User is not a superadmin
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
      "role": "string (required, one of: reporter, supervisor, admin, superadmin)",
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

#### Reset Password (Superadmin Only)
- `PUT /api/v1/auth/resetpassword` - Reset a user's password
  - **Requires**: Superadmin role
  - **Request Body**:
  ```json
  {
      "email": "string (required)",
      "password": "string (required, min 8 characters)"
  }
  ```
  - **Responses**:
    - `200 OK`: Password successfully reset
    - `400 Bad Request`: Invalid input data
    - `401 Unauthorized`: User is not a superadmin
    - `404 Not Found`: User not found
    - `500 Internal Server Error`: Database or hashing error

#### Get User
- `GET /api/v1/user` - Get user information by email
  - **Requires**: Authentication (any authenticated user)
  - **Query Parameters**:
    - `email`: User's email address (required)
  - **Responses**:
    - `200 OK`: User information
    - `400 Bad Request`: Email parameter missing
    - `404 Not Found`: User not found
    - `500 Internal Server Error`: Database error

### Incident Management Endpoints

#### Report Incident
- `POST /api/v1/incidents` - Submit a new incident report
  - **Requires**: No authentication required (anyone can report)
  - **Request Body**:
    ```json
    {
      "reporterName": "string (required)",
      "department": "string (required)",
      "position": "string (required)",
      "contactInfo": "string (required)",
      "dateOfIncident": "string (required)",
      "timeOfIncident": "string (required)",
      "locationOfIncident": "string (required)",
      "typeOfIncident": "string (required)",
      "peopleInvolved": "string (required)",
      "descriptionOfIncident": "string (required)",
      "immediateActionTaken": "string (required)",
      "injuryOrDamage": "string (required)",
      "severityLevel": "string (required, one of: near miss, minor, major, critical)",
      "supervisorNotified": "string (required)",
      "recommendedPreventiveAction": "string (required)",
      "incidentStatus": "string (optional, one of: unresolved, inprogress, resolved, defaults to unresolved)"
    }
    ```
  - **Responses**:
    - `200 OK`: Incident successfully recorded
    - `400 Bad Request`: Invalid input data or invalid severity level
    - `500 Internal Server Error`: Database error

#### Update Incident Status
- `PATCH /api/v1/incidents/:id/status` - Update incident status
  - **Requires**: Authentication (any authenticated user except reporter)
  - **Path Parameters**:
    - `id`: Incident ID (required)
  - **Request Body**:
    ```json
    {
      "status": "string (required, one of: unresolved, inprogress, resolved)"
    }
    ```
  - **Responses**:
    - `200 OK`: Status updated successfully
    - `400 Bad Request`: Invalid status value
    - `401 Unauthorized`: Missing or invalid authentication token
    - `403 Forbidden`: Reporter role or supervisor updating incident from another department
    - `404 Not Found`: Incident not found
    - `500 Internal Server Error`: Database error

#### Get Incidents
- `GET /api/v1/incidents` - Retrieve paginated list of incidents
  - **Requires**: Authentication (any authenticated user)
  - **Role-specific behavior**: Supervisors only see incidents from their department
  - **Query Parameters**:
    - `page`: Page number (default: 1)
    - `limit`: Number of items per page (default: 10, max: 50)
  - **Responses**:
    - `200 OK`: Paginated list of incidents
    - `401 Unauthorized`: Missing or invalid authentication token
    - `500 Internal Server Error`: Database error

## Role Permissions

| Role | Permissions |
|------|-------------|
| superadmin | All endpoints including user management (register, update, disable, enable, reset password, get user), report incidents, view all incidents, update incident status |
| admin | Report incidents, view all incidents, update incident status |
| supervisor | Report incidents, view own department incidents, update own department incidents |
| reporter | Report incidents via public endpoint only |

## Database Schema

The application uses two main tables defined in `tables.sql`:

### Users Table
Stores user account information:

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | SERIAL | PRIMARY KEY | Auto-incrementing unique identifier |
| name | VARCHAR(255) | NOT NULL | User's full name |
| email | VARCHAR(255) | UNIQUE NOT NULL | User's email address (used for login) |
| password | VARCHAR(255) | NOT NULL | Bcrypt hashed password |
| role | VARCHAR(50) | NOT NULL DEFAULT 'reporter' | User role (reporter, supervisor, admin, superadmin) |
| department | VARCHAR(100) | NOT NULL | User's department |
| disabled | BOOLEAN | NOT NULL DEFAULT FALSE | Account status (true = disabled) |

### Incidents Table
Stores incident reports:

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | SERIAL | PRIMARY KEY | Auto-incrementing unique identifier |
| reporter_name | VARCHAR(255) | NOT NULL | Name of person reporting the incident |
| department | VARCHAR(100) | NOT NULL | Department where incident occurred |
| position | VARCHAR(100) | NOT NULL | Position/job title of reporter |
| contact_info | VARCHAR(255) | NOT NULL | Contact information for reporter |
| date_of_incident | VARCHAR(50) | NOT NULL | Date of incident (YYYY-MM-DD format) |
| time_of_incident | VARCHAR(50) | NOT NULL | Time of incident (HH:MM format) |
| location_of_incident | VARCHAR(255) | NOT NULL | Location where incident occurred |
| type_of_incident | VARCHAR(150) | NOT NULL | Type/category of incident |
| people_involved | TEXT | NOT NULL | Description of people involved |
| description_of_incident | TEXT | NOT NULL | Detailed description of the incident |
| immediate_action_taken | TEXT | NOT NULL | Actions taken immediately after incident |
| injury_or_damage | TEXT | NOT NULL | Details of any injury or property damage |
| severity_level | VARCHAR(50) | NOT NULL | Severity level (near miss, minor, major, critical) |
| supervisor_notified | VARCHAR(255) | NOT NULL | Whether supervisor was notified |
| recommended_preventive_action | TEXT | NOT NULL | Recommended actions to prevent recurrence |
| incident_status | VARCHAR(50) | NOT NULL DEFAULT 'unresolved' | Current status (unresolved, inprogress, resolved) |

## Project Structure

```
.
├── .air.toml              # Air configuration for live reloading
├── .gitignore             # Git ignore rules
├── commit.sh              # Helper script for git operations
├── docker-compose.yml     # PostgreSQL and server service definitions
├── login.sh               # Script to access database shell
├── README.md              # This file
├── tables.sql             # Database schema definition
│
├── cmd/                   # Application entrypoint and handlers
│   ├── auth.go            # Authentication handlers (register, login, reset password)
│   ├── incidents.go       # Incident reporting handlers
│   ├── main.go            # Application initialization
│   ├── middleware.go      # Authentication middleware (JWT validation)
│   ├── routes.go          # API route definitions
│   ├── server.go          # HTTP server configuration
│   ├── types.go           # Request/response structs and type definitions
│   ├── utils.go           # Utility functions (password hashing)
│   └── users.go           # User management handlers (update, disable, enable, get user)
│
├── internal/              # Private application libraries
│   ├── db/                # Database models and connection handling
│   │   ├── db.go          # Database connection pool initialization
│   │   ├── incidents.go   # Incident database model
│   │   └── users.go       # User database model
│   └── env/               # Environment variable helpers
│       └── env.go
│
└── tmp/                   # Temporary directory (used by Air for builds)
```

## Data Flow

1. **Client Request**: HTTP request arrives at Gin router
2. **Routing**: Request directed to appropriate handler based on path and method
3. **Middleware**: Authentication middleware validates JWT token (for protected routes)
4. **Handler**: Application logic processes request, validates input
5. **Database Layer**: Handler calls appropriate model methods to interact with PostgreSQL
6. **Response**: Handler returns JSON response with appropriate status code

## Security Features

- **Password Security**: Passwords are hashed using bcrypt before storage
- **Authentication**: JWT-based authentication with expiration (72 hours)
- **Authorization**: Role-based access control enforced via middleware
- **Input Validation**: All incoming data is validated using Gin's binding mechanism
- **SQL Injection Prevention**: Uses parameterized queries via PGX
- **CORS**: Configured via gin-contrib/cors middleware with allowed origins from environment

## Error Handling

The API follows consistent error response format:
```json
{
  "error": "Descriptive error message"
}
```

HTTP status codes are used appropriately:
- 2xx: Success
- 4xx: Client errors (validation, authentication, etc.)
- 5xx: Server errors (database issues, etc.)

## Extending the Application

### Adding New Endpoints
1. Define handler function in appropriate file under `cmd/`
2. Add route to `routes.go` within the appropriate group
3. Apply middleware as needed (authMiddleware() for protected routes)
4. Update corresponding model in `internal/db/` if database changes needed

### Adding New Database Tables/Columns
1. Modify `tables.sql` with schema changes
2. Update corresponding model in `internal/db/`
3. Update handler functions in `cmd/` to use new fields
4. Rebuild and restart containers to apply changes

### Configuration Changes
- Update `.env` file for environment-specific settings
- Modify `.air.toml` for Air live reload configuration
- Adjust Docker Compose file for service changes

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Commit your changes (`git commit -m 'Add amazing feature'`)
5. Push to the branch (`git push origin feature/amazing-feature`)
6. Open a Pull Request

Please ensure your code follows:
- Go formatting standards (`go fmt`)
- Clear, descriptive comments
- Consistent error handling
- Proper input validation

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Gin-Gonic team for the excellent web framework
- PostgreSQL team for the reliable database
- JWT and bcrypt libraries for secure authentication
- Air team for the live reload development experience
# Issue Tracker

A RESTful API for tracking incidents/issues built with Go, Gin, and PostgreSQL.

## Features

- User authentication (registration, login, update, disable)
- Incident reporting and tracking
- Role-based access control (implied from User model)
- Docker Compose setup for easy development

## Technologies

- Go 1.22+ (or as per go.mod)
- Gin web framework
- PostgreSQL database
- PGX PostgreSQL driver

## Getting Started

### Prerequisites

- Docker and Docker Compose (for containerized setup)
- Go (if running locally)

### Running with Docker Compose

1. Clone the repository
2. Copy `.env.example` to `.env` (if exists) or create one based on the environment variables below
3. Run `docker-compose up -d` to start the PostgreSQL container
4. Run `go run ./cmd/main.go` to start the application

Alternatively, you can use the provided scripts:

- `./createtables.sh` - creates necessary database tables (if not already created)
- `./login.sh` - opens a psql shell to the PostgreSQL container
- `./commit.sh` - helper script to stage, commit, and push changes

### Environment Variables

The following environment variables are used:

- `PORT`: The port on which the server runs (default: 3001)
- `dbConnStr`: PostgreSQL connection string (default: `postgres://tracker_user:tracker_password@localhost:5432/issuetracker`)

These can be set in a `.env` file or exported in the shell.

## API Endpoints

### Health Check

- `GET /api/v1/ping` - Returns a pong message

### Authentication

- `POST /api/v1/auth/register` - Register a new user (requires superadmin role)
  - Request body: `{ "email": "string", "name": "string", "password": "string", "role": "string" }`
  - Password must be at least 8 characters
  - Role must be one of: reporter, supervisor, admin, superadmin

- `POST /api/v1/auth/login` - Login a user
  - Request body: `{ "email": "string", "password": "string" }`
  - Returns a JWT token and user information

- `PUT /api/v1/auth/update` - Update a user (requires superadmin role)
  - Requires authentication middleware
  - (Implementation in progress)

- `PUT /api/v1/auth/disable` - Disable a user (requires superadmin role)
  - Requires authentication middleware
  - (Implementation in progress)

### Incidents

_(Endpoints not fully implemented in the provided code, but the model exists)_

## Database Schema

The application uses two main tables:

1. `users` - stores user information (id, name, email, password hash, role)
2. `incidents` - stores incident reports (see Incident struct in internal/db/incidents.go for fields)

## Scripts

- `commit.sh` - Helper script for committing changes (stages, commits, and pushes)
- `createtables.sh` - Creates database tables by running the SQL in tables.sql against the database
- `login.sh` - Opens a psql shell to the PostgreSQL container for direct database access

## Project Structure

```
.
├── cmd/                - Application entrypoint and route handlers
│   ├── auth.go         - Authentication handlers (register, login)
│   ├── main.go         - Application initialization and server startup
│   ├── middleware.go   - Authentication middleware
│   ├── routes.go       - Route definitions
│   ├── server.go       - Server configuration
│   ├── types.go        - Type definitions (e.g., request/response structs)
│   ├── users.go        - User handlers (update, disable)
│   └── utils.go        - Utility functions (password hashing, etc.)
├── internal/           - Private application libraries
│   ├── db/             - Database models and connection pooling
│   │   ├── db.go       - Database connection and model initialization
│   │   ├── incidents.go - Incident model and severity levels
│   │   └── users.go    - User model
│   └── env/            - Environment variable helpers
├── docker-compose.yml  - PostgreSQL service definition
├── go.mod              - Go module definition
├── go.sum              - Go module checksums
├── test.sql            - Sample SQL for table creation (same as tables.sql?)
├── tables.sql          - SQL script to create database tables
├── README.md
└── .air.toml           - Configuration for Air (live reload for Go)
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Open a pull request
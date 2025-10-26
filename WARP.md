# WARP.md

This file provides guidance to WARP (warp.dev) when working with code in this repository.

## Project Overview

Vado-server is a dual-protocol Go server providing both HTTP/REST and gRPC endpoints for task management with user authentication. Built using Gin (HTTP), gRPC, GORM (PostgreSQL), and JWT authentication with access/refresh tokens.

## Architecture

### Layered Structure
The codebase follows a clean architecture pattern with clear separation of concerns:

- **cmd/server**: Application entry point (`main.go`) that starts both HTTP and gRPC servers concurrently with graceful shutdown
- **api/proto**: Protocol Buffer definitions for gRPC services (auth, chat, hello, server)
- **api/pb**: Auto-generated gRPC code from `.proto` files (do not edit manually)
- **internal/handler**: Dual protocol handlers
  - `handler/http`: Gin HTTP handlers for web UI and REST endpoints
  - `handler/grpc`: gRPC service implementations
- **internal/services**: Business logic layer
- **internal/repository**: Data access layer (GORM)
- **internal/models**: GORM database models
- **internal/middleware**: HTTP middleware (JWT validation, auth checks)
- **internal/auth**: JWT token creation and parsing logic
- **internal/appcontext**: Shared application context containing logger and DB connection
- **internal/router**: HTTP route configuration
- **web**: Frontend templates and static assets

### Key Architectural Patterns

**AppContext Pattern**: The `appcontext.AppContext` struct is passed throughout the application, carrying the Zap logger and GORM DB connection.

**Dual Authentication**: 
- HTTP endpoints use cookie-based JWT authentication via `middleware.CheckJWT()`
- gRPC endpoints use `AuthInterceptor` that validates Bearer tokens from metadata
- Public endpoints (Login, Refresh, Ping) bypass authentication

**Database Relations**:
- Users → Tasks (one-to-many with cascade delete)
- Users → Roles (many-to-many via `user_roles` join table)

**JWT Token System**: Uses access tokens with refresh token support. Custom claims include `UserID` and `Roles`.

## Development Commands

### Running the Application

```bash
# Local development (uses .env file)
go run ./cmd/server/main.go

# With Docker Compose (production-like)
docker-compose up --build

# Access points:
# - HTTP: http://localhost:5555 (local) or :5556 (docker)
# - gRPC: localhost:50051
```

### Dependency Management

```bash
# Add/remove unused dependencies
go mod tidy

# Add a specific package
go get -u <package>

# Show all modules
go list -m all

# Show why a module was added
go mod why <package>
```

### gRPC Code Generation

After modifying `.proto` files in `api/proto/`:

```bash
# Generate Go code for a specific proto file
protoc --go_out=./ --go-grpc_out=./ api/proto/auth.proto

# Generate for all proto files
protoc --go_out=./ --go-grpc_out=./ api/proto/*.proto
```

Generated files go to `api/pb/<service>/` directory.

### Database Management

```bash
# Connect to Postgres container
docker exec -it vado-postgres bash

# Access psql
psql -U vadmark -d vadodb

# Inside psql:
\dt              # List tables
\d tasks         # Show table structure

# Reset database (removes all data)
docker volume rm vado-server_postgres-data
docker-compose up postgres
```

Database schema is in `db/01_schema.sql` and seed data in `db/02_inserts.sql`. These run automatically on first container start.

### Process Management

```bash
# Find process using a port
sudo lsof -i:8080

# Kill process
sudo kill -9 <PID>
```

## Environment Configuration

The application reads from `.env` file (local development) or Docker environment variables:

- `GIN_MODE`: `debug` or `release`
- `PORT`: HTTP server port (default: 5555 local, 5556 docker)
- `GRPC_PORT`: gRPC server port (default: 50051)
- `POSTGRES_DSN`: PostgreSQL connection string
- `JWT_SECRET`: Secret key for JWT signing

**Note**: Docker uses `postgres` hostname for database, local development uses `localhost`.

## Code Conventions

### Import Organization
Group imports as: standard library, external packages, internal packages (as seen in `main.go`).

### Logging
Use the Zap sugared logger from `appCtx.Log`. Methods: `Infow()`, `Errorw()`, `Fatalw()`, `Warnw()`.

### Error Handling
- HTTP handlers return error pages via templates
- gRPC handlers return gRPC status errors with codes (e.g., `codes.Unauthenticated`)

### Constants
Centralized in `internal/constants/` subdirectories:
- `code`: Context keys and standard codes
- `env`: Environment variable names
- `role`: Role name constants
- `route`: HTTP route paths

## Testing

No test framework is currently configured. Before adding tests, check with the team on preferred testing approach.

## Stack Reference

- **Web Framework**: Gin
- **Database**: PostgreSQL 15 with GORM
- **RPC**: gRPC with Protocol Buffers
- **Auth**: JWT (golang-jwt/jwt/v5) with access and refresh tokens
- **Logging**: Zap
- **Sessions**: gin-contrib/sessions (cookie store)
- **Password Hashing**: golang.org/x/crypto

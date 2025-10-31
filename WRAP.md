# WRAP.md - Vado Server Project Overview

## Project Summary
**Vado Server** is a Go-based backend service utilizing gRPC, gRPC-Web, and REST APIs with PostgreSQL database integration. The project implements authentication (JWT-based), task management, and real-time communication capabilities with Kafka integration.

## Architecture

### Technology Stack
- **Language**: Go 1.25.1
- **API Protocols**: gRPC, gRPC-Web, REST (Gin)
- **Database**: PostgreSQL (GORM ORM)
- **Authentication**: JWT (access + refresh tokens)
- **Logging**: Zap (Uber)
- **Message Queue**: Kafka
- **Containerization**: Docker, Docker Compose

### Project Structure
```
vado-server/
├── cmd/                    # Application entry points
│   └── server/            # Main server application
├── internal/              # Private application code
│   ├── app/              # Application logic
│   ├── config/           # Configuration management
│   ├── domain/           # Business domain models
│   ├── infra/            # Infrastructure (DB, external services)
│   ├── trasport/         # Transport layer (HTTP, gRPC handlers)
│   └── util/             # Utility functions
├── api/                   # API definitions
│   └── proto/            # Protocol Buffer definitions
├── db/                    # Database scripts and migrations
├── web/                   # Web frontend assets
│   └── static/js/pb/     # Generated gRPC-Web JS files
├── docker-compose.yml     # Main services (server + postgres)
├── docker-compose.kafka.yml  # Kafka services
├── Dockerfile             # Server container image
├── Makefile              # Build and deployment commands
└── .env                  # Environment configuration
```

## Quick Start

### Prerequisites
- Go 1.25.1+
- Docker & Docker Compose
- Protocol Buffer compiler (`protoc`)
- `protoc-gen-go` and `protoc-gen-go-grpc` plugins
- `protoc-gen-grpc-web` plugin (for web clients)

### Development Setup

1. **Start all services** (server, postgres, kafka):
   ```bash
   make all-up
   ```

2. **Start only main services** (server + postgres):
   ```bash
   make up-main
   ```

3. **View logs**:
   ```bash
   make logs
   ```

4. **Access PostgreSQL shell**:
   ```bash
   make psql
   ```

### Build Commands

#### Generate Protocol Buffers
```bash
# Generate Go gRPC files
make go-proto

# Generate JavaScript gRPC-Web files
make web-proto
```

#### Docker Operations
```bash
# Full rebuild (clean database)
make rebuild

# Rebuild only server (keep database)
make rebuild-server

# Stop all services
make all-down

# Clean Docker cache
make clean
```

### Kafka Operations
```bash
# Start Kafka services
make kafka-up

# Stop Kafka services
make kafka-down
```

## Development Workflow

### Running the Server Locally
```bash
go run ./cmd/server/main.go
```

### Managing Dependencies
```bash
# Add/remove unused dependencies
go mod tidy

# Check why a module is included
go mod why <package>

# List all modules
go list -m all
```

### Database Operations

#### Connect to PostgreSQL
```bash
# Via Docker
docker exec -it vado_postgres bash
psql -U vadmark -d vadodb

# Via Make
make psql
```

#### Common psql Commands
```sql
-- List tables
\dt

-- Show table structure
\d tasks

-- Exit
\q
```

#### Reset Database
```bash
# Remove volume to recreate database
docker volume rm vado-server_postgres-data

# Or full rebuild
make rebuild
```

## API Testing

### gRPC Testing with grpcurl
```bash
# Test authentication
grpcurl -plaintext \
  -import-path ./api/proto \
  -proto auth.proto \
  -d '{"username": "test", "password": "test"}' \
  localhost:50051 AuthService/Login

# Test ping service
grpcurl -plaintext \
  -import-path ./api/proto \
  -proto ping.proto \
  localhost:50051 PingService/Ping
```

## Port Management

### Kill Process on Port
```bash
# Find process
sudo lsof -i:8080

# Kill process
sudo kill -9 PID
```

## Environment Configuration
Configuration is stored in `.env` file. Key variables typically include:
- Database connection parameters
- JWT secret keys
- Server ports
- Kafka broker addresses

## Future Enhancements
- [ ] golang-migrate integration for database migrations
- [ ] Additional API endpoints
- [ ] Enhanced monitoring and observability

## Help
Run `make help` to see all available commands:
```bash
make help
```

## Key Dependencies
- `github.com/gin-gonic/gin` - HTTP web framework
- `google.golang.org/grpc` - gRPC framework
- `gorm.io/gorm` - ORM library
- `github.com/golang-jwt/jwt/v5` - JWT authentication
- `go.uber.org/zap` - Structured logging
- `github.com/segmentio/kafka-go` - Kafka client
- `github.com/improbable-eng/grpc-web` - gRPC-Web support

## Notes
- Default PostgreSQL container: `vado_postgres`
- Default database: `vadodb`
- Default user: `vadmark`
- Project name in Docker: `vado-app`

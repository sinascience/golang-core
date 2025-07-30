# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Running the Application
- **With Docker (recommended)**: `docker-compose up --build`
- **Locally**: `go run ./cmd/server/main.go` (after setting up MySQL and running migrations)
- **Hot-reloading with Air**: `air` (automatically configured in Docker)

### Database Operations
- **Run migrations (Docker)**: `docker-compose run --rm app go run ./cmd/migrate/main.go up`
- **Run migrations (local)**: `migrate -database 'mysql://user:pass@tcp(127.0.0.1:3306)/db_name' -path database/migrations up`
- **Create new migration**: `migrate create -ext sql -dir database/migrations -seq migration_name`

### Testing
- **Run tests**: `go test ./...`
- **Run specific package tests**: `go test ./pkg/validator`
- **Test with verbose output**: `go test -v ./...`

### Documentation
- **Generate/update Swagger docs**: `swag init -g ./cmd/server/main.go --parseDependency --parseInternal`
- **View API docs**: Start server and visit `http://localhost:3000/swagger/index.html`

### Build Commands
- **Build binary**: `go build -o ./tmp/server ./cmd/server`
- **Build dependencies**: `go mod tidy`

## Architecture Overview

This is a Go REST API using a **Fat Model (Active Record) pattern** with clean architecture principles:

### Request Flow
1. **HTTP Handler** (`internal/handler/http/`) - Receives requests, validates input, calls services
2. **Service Layer** (`internal/service/`) - Contains business logic, orchestrates operations  
3. **Model Layer** (`internal/model/`) - Database operations and data structures (Fat Model pattern)

### Key Architectural Patterns
- **Dependency Injection**: Services and handlers receive dependencies via constructors
- **Adapter Pattern**: Used for third-party integrations (see `internal/adapter/storage/`)
- **Middleware**: Authentication and other cross-cutting concerns in `internal/middleware/`

### Core Components
- **Server Setup**: `internal/server/server.go` - App initialization with graceful shutdown
- **Routing**: `internal/server/routes.go` - Route registration and dependency wiring
- **Configuration**: `configs/config.go` - Environment variable loading
- **Database**: `internal/database/database.go` - GORM connection management

## Project Structure Patterns

### Adding New Features
Follow this pattern (e.g., for "comments"):
1. Create migration: `migrate create -ext sql -dir database/migrations -seq create_comments_table`
2. Run migration to update schema
3. Create model: `internal/model/comment_model.go` with struct and database methods
4. Create service: `internal/service/comment_service.go` for business logic
5. Create handler: `internal/handler/http/comment_handler.go` with validation tags
6. Register routes in `internal/server/routes.go`

### Model Methods Pattern
Models use the Fat Model pattern with methods like:
- `Create(db *gorm.DB) error`
- `Update(db *gorm.DB) error`
- `Delete(db *gorm.DB) error`
- `GetByID(db *gorm.DB, id string) error`

### Handler Pattern
- Use `pkg/validator.ValidateStruct()` for input validation
- Return responses with `pkg/response` helpers
- Add Swagger annotations above handler functions
- Always use service layer, never call models directly

## Technology Stack

### Core Framework
- **Web Framework**: Fiber v2 (Express-like for Go)
- **ORM**: GORM with MySQL driver
- **Migrations**: golang-migrate
- **Authentication**: JWT with golang-jwt/jwt/v5

### Development Tools
- **Hot Reload**: Air (configured in `.air.toml`)
- **Documentation**: Swaggo for Swagger generation
- **Containerization**: Docker with docker-compose
- **Validation**: go-playground/validator/v10

### Utilities
- **Logging**: Standard library `slog` with JSON output
- **Config**: godotenv for environment variables
- **File Upload**: Custom uploader utility in `pkg/uploader/`
- **Response Formatting**: Standardized API responses in `pkg/response/`

## Environment Configuration

Required `.env` variables:
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` - Database connection
- `JWT_ACCESS_SECRET_KEY`, `JWT_REFRESH_SECRET_KEY` - JWT signing keys
- `JWT_ACCESS_EXPIRATION_IN_MINUTES`, `JWT_REFRESH_EXPIRATION_IN_HOURS` - Token expiration
- `CORS_ALLOWED_ORIGINS` - CORS configuration

## Testing Approach

- Unit tests in `*_test.go` files alongside source code
- Table-driven test pattern (see `pkg/validator/validator_test.go`)
- Manual API testing examples in `curl_tests` file
- Focus on testing business logic in services and utility functions

## Background Processing

- Uses goroutines for async operations (e.g., file uploads)
- `sync.WaitGroup` tracks background jobs for graceful shutdown
- Graceful shutdown waits for all background processes before exiting
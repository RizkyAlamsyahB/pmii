# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

PMII Backend API - A Go REST API built with Gin framework using Clean Architecture patterns.

## Development Commands

```bash
# Run the application (requires .env file and running PostgreSQL)
go run cmd/api/main.go

# Run all tests
go test ./...

# Run a specific test function
go test -run TestGetAllUsers_Success ./internal/service/

# Run tests with verbose output
go test -v ./internal/service/...

# Build the binary
go build -o pmii-backend cmd/api/main.go

# Start with Docker Compose (includes PostgreSQL)
docker-compose up -d

# Rebuild after code changes
docker-compose up -d --build
```

## Environment Setup

Copy `.env.example` to `.env` and configure:
- Database: `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`
- JWT: `JWT_SECRET`, `JWT_EXPIRATION_HOURS`
- Server: `PORT`, `ENV` (development/production), `ALLOWED_ORIGINS`

## Architecture

The codebase follows Clean Architecture with clear separation of concerns:

```
cmd/api/main.go          - Application entry point, dependency injection
config/                  - Configuration loading via Viper (.env)
internal/
├── domain/             - Domain models (User, Post, Category, etc.) - GORM entities
├── dto/                - Request/Response DTOs
│   ├── requests/       - Input validation structs (uses gin binding tags)
│   └── responses/      - API response formatting
├── handlers/           - HTTP handlers (transport layer)
├── middleware/         - Auth, CORS, RBAC, rate limiting, recovery
├── repository/         - Data access layer (interfaces + GORM implementations)
├── routes/             - Route definitions and middleware wiring
└── service/            - Business logic layer
pkg/
├── database/           - Database connection, migrations, seeding
├── logger/             - Application logging
└── utils/              - JWT, password hashing, token blacklist
migrations/             - SQL migration files (golang-migrate format)
```

## Key Patterns

**Dependency Flow**: main.go → repository → service → handler → routes

**Adding a New Entity** (e.g., "Widget"):
1. Create domain model: `internal/domain/widget.go`
2. Create DTOs: `internal/dto/requests/widget_request.go`, `internal/dto/responses/widget_response.go`
3. Create repository interface + implementation: `internal/repository/widget_repository.go`
4. Create service: `internal/service/widget_service.go`
5. Create handler: `internal/handlers/widget_handler.go`
6. Wire up in `cmd/api/main.go` (instantiate repo → service → handler)
7. Add routes in `internal/routes/routes.go`

**Authentication**: JWT-based with token blacklisting for logout. JWT claims contain `user_id` and `user_role`.

**RBAC**: Role-based access control via middleware:
- `RequireRole("1")` - Admin only (role=1)
- `RequireAnyRole("1", "2")` - Multiple roles allowed
- `RequireOwnerOrAdmin("id")` - Resource owner or admin access

**User Roles**: Role 1 = Admin, Role 2 = Author

**Request Validation**: Uses Gin binding tags in request DTOs. Example:
```go
type CreateUserRequest struct {
    FullName string `json:"full_name" binding:"required,min=2,max=100"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
}
```

**Migrations**: Auto-run on startup via `database.RunMigrations()`. Files in `migrations/` folder follow `NNNNNN_description.{up,down}.sql` naming.

**Seeding**: Default users seeded automatically on startup via `database.SeedDefaultUsers()`.

## Testing Patterns

Tests use manual mock structs with function fields. Each repository method has a corresponding `*Func` field:

```go
type MockUserRepository struct {
    FindAllFunc  func() ([]domain.User, error)
    FindByIDFunc func(id int) (*domain.User, error)
    // ...
}

func (m *MockUserRepository) FindAll() ([]domain.User, error) {
    if m.FindAllFunc != nil {
        return m.FindAllFunc()
    }
    return nil, errors.New("mock not configured")
}
```

Test naming convention: `Test<Function>_<Scenario>` (e.g., `TestCreateUser_EmailAlreadyExists`)

## API Structure

Base URL: `/v1`

**Auth Routes:**
- `POST /v1/auth/login` - Login (rate limited: 60 req/min)
- `POST /v1/auth/logout` - Logout (requires auth)

**Admin Routes** (`/v1/admin/*` - requires admin role):
- `GET /v1/admin/dashboard` - Admin dashboard
- CRUD for testimonials: `/v1/admin/testimonials`
- CRUD for members: `/v1/admin/members`

**User Routes** (`/v1/users` - requires auth):
- `GET /v1/users` - List all (admin only)
- `POST /v1/users` - Create (admin only)
- `GET /v1/users/:id` - Get by ID (owner or admin)
- `PUT /v1/users/:id` - Update (admin only)
- `DELETE /v1/users/:id` - Delete (admin only)

## Database

PostgreSQL with GORM ORM. Entities use soft deletes via `gorm.DeletedAt`. Repository pattern with interfaces enables unit testing via mocking.

# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

PMII Backend API - A Go REST API built with Gin framework using Clean Architecture patterns.

## Development Commands

```bash
# Run the application (requires .env file and running PostgreSQL)
go run cmd/api/main.go

# Run tests
go test ./...

# Run specific test file
go test ./internal/service/auth_service_test.go

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
│   ├── requests/       - Input validation structs
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

**Authentication**: JWT-based with token blacklisting for logout. JWT claims contain `user_id` and `user_role`.

**RBAC**: Role-based access control via middleware:
- `RequireRole("1")` - Admin only (role=1)
- `RequireAnyRole("1", "2")` - Multiple roles allowed
- `RequireOwnerOrAdmin("id")` - Resource owner or admin access

**User Roles**: Role 1 = Admin, Role 2 = Author

**Migrations**: Auto-run on startup via `database.RunMigrations()`. Files in `migrations/` folder follow `NNNNNN_description.{up,down}.sql` naming.

**Seeding**: Default users seeded automatically on startup via `database.SeedDefaultUsers()`.

## API Structure

Base URL: `/v1`

- `POST /v1/auth/login` - Login (rate limited)
- `POST /v1/auth/logout` - Logout (requires auth)
- `GET /v1/admin/dashboard` - Admin only
- `GET /v1/users` - List users (admin only)
- `GET /v1/users/:id` - Get user (owner or admin)

## Database

PostgreSQL with GORM ORM. Entities use soft deletes via `gorm.DeletedAt`. Repository pattern with interfaces enables unit testing via mocking.

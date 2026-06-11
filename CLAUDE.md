# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

Go REST web application for a book store backed by PostgreSQL. Uses Gin (HTTP), uber-go/fx (IoC/DI), and pgx/v5 via the `database/sql` stdlib adapter. No ORM — raw SQL throughout.

## Commands

```bash
# Run
go run .

# Build
go build -o book-store .

# Run all tests
go test ./...

# Run a single test
go test ./... -run TestFunctionName

# Add/update dependencies
go mod tidy

# Regenerate Swagger docs after changing handler annotations
$(go env GOPATH)/bin/swag init --parseInternal
```

## Architecture

Strict layered architecture — each layer depends only on the layer below it via interfaces:

```
main.go → fx wiring
internal/
  config/       Config struct + LoadConfig() + DSN()
  domain/       Pure structs and request DTOs (no external deps)
  repository/   DB interfaces + postgres impls (*sql.DB, raw SQL)
  service/      Business logic; UUID generation lives here
  handler/      Gin handlers; binds JSON, calls service, maps errors
  server/       NewServer(): gin.Engine + all routes registered
migrations/
  001_init_schema.sql
```

**fx graph** (dependency resolution order):
`LoadConfig` → `newDatabase (*sql.DB)` → repositories → services → handlers → `NewServer` → `registerHooks (lifecycle)`

`ErrNotFound` is defined once in `repository/customer_repository.go` and used by all repository implementations. Handlers check `errors.Is(err, repository.ErrNotFound)` to decide 404 vs 500.

`order_items.price` stores a price snapshot at order creation time (the service resolves the current product price before writing).

## API Endpoints (prefix `/api/v1`)

| Method | Path | Description |
|--------|------|-------------|
| GET/POST | `/customers` | list / create |
| GET/PUT/DELETE | `/customers/:id` | get / replace / delete |
| GET | `/customers/:id/orders` | orders for a customer |
| GET/POST | `/products` | list / create |
| GET/PUT/DELETE | `/products/:id` | get / replace / delete |
| GET/POST | `/orders` | list / create (body includes items array) |
| GET | `/orders/:id` | get order with items |
| PATCH | `/orders/:id/status` | update status only |
| DELETE | `/orders/:id` | delete |

## Environment Variables

| Variable | Default |
|---|---|
| `DB_HOST` | `localhost` |
| `DB_PORT` | `5432` |
| `DB_USER` | `postgres` |
| `DB_PASSWORD` | `postgres` |
| `DB_NAME` | `bookstore` |
| `DB_SSLMODE` | `disable` |
| `SERVER_PORT` | `8080` |

## Running Locally

```bash
# Start Postgres
docker run -d --name bookstore-db \
  -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=bookstore \
  -p 5432:5432 postgres:16-alpine

# Apply schema
psql postgresql://postgres:postgres@localhost:5432/bookstore -f migrations/001_init_schema.sql

# Run
go run .
```

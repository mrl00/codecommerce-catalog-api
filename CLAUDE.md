# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**codecommerce-catalog-api** — a REST API for a product catalog with categories, built in Go using gorilla/mux, PostgreSQL, and raw SQL.

## Go Project Structure: cmd/internal Pattern

```
cmd/server/main.go                    # Entry point, wiring, graceful shutdown
internal/
├── entities/                         # Domain models (Category, Product)
│   ├── category.go / category_test.go
│   └── product.go  / product_test.go
├── handler/                          # HTTP handlers (JSON req/res)
│   ├── category_handler.go / category_handler_test.go
│   ├── product_handler.go  / product_handler_test.go
│   ├── health.go           / health_test.go
│   └── helpers.go          / helpers_test.go   # parsePaginationParams, errorResponse
├── router/router.go                  # gorilla/mux route registration
├── service/                          # Business logic + repository interfaces
│   ├── service.go                    # Interfaces, PaginationParams, PaginatedResult, sentinels
│   ├── category_service.go / category_service_test.go
│   └── product_service.go  / product_service_test.go
├── database/                         # PostgreSQL implementations (raw SQL via lib/pq)
│   ├── category_db.go
│   └── product_db.go
└── migrate/                          # Embedded migration runner with schema_migrations tracking
    ├── migrate.go
    └── migrations/001_initial_schema.sql
tests/                                # Hurl integration tests
    ├── health.hurl / categories.hurl / products.hurl
    └── run_tests.sh
```

## Key Dependencies

- **`github.com/gorilla/mux`** — router. Routes: `r.HandleFunc(path, handler).Methods("METHOD")`
- **`github.com/lib/pq`** — PostgreSQL driver via `database/sql`
- **`github.com/google/uuid`** — UUID v4 for entity IDs

## Conventions

### Code Patterns
- Handlers use `errorResponse(w, msg, code)` for all JSON error responses
- Pagination is parsed via `parsePaginationParams(r)` in `handler/helpers.go`
- Service layer owns sentinel errors: `service.ErrCategoryNotFound`, `service.ErrProductNotFound`
- Error checks in handlers use `==` (not `errors.Is`) against sentinel errors
- Repository interfaces defined in `service/service.go`, implemented in `database/`
- Generic `PaginatedResult[T]` for all list endpoints
- Entities generate their own UUID + timestamps in constructors (`NewCategory`, `NewProduct`)
- Logging uses `log/slog` (structured, leveled)
- Server supports graceful shutdown (SIGINT/SIGTERM with 10s drain)

### Database
- **Naming**: `tb_<entity>`, `pk_<entity>`, `ts_<entity>_created_at`, `tx_<field>`, `nr_<field>`, `fk_<field>`
- **FK cascade**: `fk_category` has `ON DELETE CASCADE` on `tb_product`
- **Price**: `BIGINT` (cents) to avoid float precision issues
- **Migrations**: embedded SQL files in `internal/migrate/migrations/`, numbered prefix (e.g., `001_`), tracked in `schema_migrations` table
- **Pagination**: `COUNT(*) + LIMIT/OFFSET`
- **List ordering**: currently `ORDER BY pk_*` (UUID-based, effectively random)

### API
- All endpoints under `/api/` prefix
- JSON request/response with `Content-Type: application/json`
- Create returns `201`, Delete returns `204`, Not Found returns `404`
- Validation errors return `400` with `{"error": "message"}`

## Testing

- **Unit tests** (85 test cases, 37 top-level): `go test ./...`
  - Entities: constructor + ResetUpdatedAt tests
  - Service: mock-based CRUD + pagination tests
  - Handler: HTTP request/response tests (table-driven subtests)
  - Helpers: errorResponse + parsePaginationParams tests
- **Integration tests**: hurl tests in `tests/` covering full HTTP lifecycle
  - `./tests/run_tests.sh` (requires running server + PostgreSQL)
- **No tests**: `database/`, `router/`, `migrate/` packages

## Commands

```bash
go run ./cmd/server/             # run dev server on :8080
go build -o server ./cmd/server/ # build binary
go test ./...                    # run unit tests
go test -v ./... -count=1        # verbose, no cache
go mod tidy                      # sync dependencies
./tests/run_tests.sh             # run hurl integration tests
```

## CI/CD

- **ci.yml** — runs on PRs to `dev`: build → unit tests → start server → hurl integration tests (Postgres service container)
- **docker-publish.yml** — runs on push to `main`: build + push Docker image to DockerHub, sign with cosign
- **branch-check.yml** — PRs to `main` must come from `dev`

## Backlog

See `docs/backlog.md` for pending tasks. Always check it before starting new work.

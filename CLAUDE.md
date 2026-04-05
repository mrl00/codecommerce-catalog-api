# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Go Project Structure: cmd/internal Pattern

This project follows Go's standard project layout:

- **`cmd/server/`** — application entry point (`main.go`). This is the only package compiled into a binary.
- **`internal/handler/`** — HTTP request handlers. Each handler uses `http.ResponseWriter`/`*http.Request` and uses `errorResponse()` for structured JSON errors.
- **`internal/router/`** — route definitions using `gorilla/mux`. All routes are registered in `router.New()`, which returns an `*mux.Router`.
- **`internal/service/`** — service layer with business logic, repository interfaces (`CategoryRepository`, `ProductRepository`), shared types (`PaginationParams`, `PaginatedResult`), and sentinel errors (`ErrCategoryNotFound`, `ErrProductNotFound`).
- **`internal/database/`** — repository implementations using raw SQL via `lib/pq`. Methods returning `PaginatedResult` use `COUNT(*)` + `LIMIT/OFFSET` for pagination.

## Key Dependencies

- **`github.com/gorilla/mux`** — router. Routes: `r.HandleFunc(path, handler).Methods("METHOD")`
- **`github.com/lib/pq`** — PostgreSQL driver via `database/sql`
- **`github.com/google/uuid`** — UUID v4 for entity IDs

## Database

- **Naming convention**: `tb_<entity>`, `pk_<entity>`, `ts_<entity>_created_at`, `tx_<field>`, `nr_<field>`, `bl_<field>`
- **Cascading deletes**: `fk_category` has `ON DELETE CASCADE` on `tb_product`
- **Price stored as**: `BIGINT` (cents) to avoid float precision issues
- **Migration files**: `migrations/` directory, numbered prefix (e.g., `001_initial_schema.sql`)

## Commands

```bash
go run ./cmd/server/          # run development server on :8080
go build -o server ./cmd/server/  # build binary
go mod tidy                   # sync dependencies
```

## Backlog

See `docs/backlog.md` for pending tasks. Always check it before starting new work.

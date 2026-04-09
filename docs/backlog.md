# Backlog

Task list and improvements for the codecommerce-catalog-api project.

## Done

- [x] Refactor entities to use `time.Time` for date fields
- [x] Apply database naming convention (`tb_`, `pk_`, `tx_`, `nr_`, `ts_`, `fk_`)
- [x] Fix `Save*` to use entity ID instead of generating a new UUID
- [x] Add Update and Delete operations for Category and Product
- [x] Create CRUD handlers for Category and Product
- [x] Add route for products by category (`GET /api/categories/{id}/products`)
- [x] Connect database to server via `lib/pq`
- [x] Split `entity.go` into separate files (`category.go`, `product.go`)
- [x] Create SQL database initialization schema (`migrations/001_initial_schema.sql`)
- [x] Add input validation in handlers (name, price, category required; errors as JSON)
- [x] Standardize module name (`goapi` → `codecommerceapi`)
- [x] Fix CI `secrets.DATABASE_URL` — conflicts with the postgres service defined in the workflow
- [x] Add `_test.go` files — entity creation tests and service layer tests with full CRUD + pagination mocks
- [x] Replace `db.Prepare()` with direct `db.Exec()`/`db.Query()` (statements are single-use)
- [x] Change `float64` to `int64` (cents) for price
- [x] Implement cascade delete for products when deleting a category (via `ON DELETE CASCADE` FK constraint)
- [x] Add pagination to `FindAllProducts()`, `FindAllCategories()`, and `FindProductsByCategoryID()` (query params: `page`, `per_page`)
- [x] Optimize `Delete*` in service: use `DELETE WHERE` + check rows affected instead of pre-fetch
- [x] Add embedded migration runner to execute SQL files on server startup
- [x] Add handler unit tests — 47 tests covering CRUD, pagination parsing, input validation, and error responses for all endpoints
- [x] Add branch protection workflow (`branch-check.yml`) — PRs to `main` only allowed from `dev`
- [x] Move CI trigger from `main` to `dev` branch
- [x] Hardcode `DATABASE_URL` in CI workflow — currently uses `secrets.DATABASE_URL` which conflicts with the inline Postgres service
- [x] Add `schema_migrations` tracking table to migration runner — prevents re-executing already applied migrations
- [x] Consolidate duplicate pagination parser — extracted `parsePaginationParams()` into `handler/helpers.go`
- [x] Add structured logging with `log/slog` — replaced `log.Println`/`log.Fatalf` with leveled, structured output
- [x] Implement graceful shutdown — handle OS signals (`SIGINT`/`SIGTERM`) and drain connections before exit
- [x] Return `[]` instead of `null` for empty paginated `Items` — initialize slices with `make([]T, 0)` before marshalling
- [x] Fix Docker image name typo in `docker-publish.yml` — `codecomerce` → `codecommerce`

## Todo

### High

- [ ] Add HTTP middleware — request logging (method, path, status, duration), CORS headers, panic recovery, and request ID injection via `X-Request-Id` header. Implement as `func(http.Handler) http.Handler` middleware chain in `internal/router/` or a new `internal/middleware/` package. Apply globally to the mux router. For CORS, support configurable allowed origins via env var (`CORS_ALLOWED_ORIGINS`).
- [ ] Check `RowsAffected()` in `UpdateCategory`/`UpdateProduct` — currently `database/category_db.go:UpdateCategory` and `database/product_db.go:UpdateProduct` return `nil` even if zero rows matched (row deleted between `FindByID` and `Update`). Fix: call `result.RowsAffected()` after `db.Exec` and return `service.ErrCategoryNotFound`/`service.ErrProductNotFound` if zero rows affected, matching the pattern already used in `DeleteCategory`/`DeleteProduct`.

### Medium

- [ ] Use `errors.Is()` for sentinel error comparisons in handlers — all error checks in `category_handler.go` and `product_handler.go` currently use `err == service.ErrCategoryNotFound` / `err == service.ErrProductNotFound`. Replace with `errors.Is(err, ...)` so comparisons remain correct if errors get wrapped (e.g., via `fmt.Errorf("...: %w", err)`). Affects 5 comparisons across both handler files.
- [ ] Order list queries by `ts_*_created_at DESC` instead of primary key — `FindAllCategories` orders by `pk_category`, `FindAllProducts` and `FindProductsByCategoryID` order by `pk_product`. Since PKs are UUIDs, the order is effectively random. Change to `ORDER BY ts_category_created_at DESC` / `ORDER BY ts_product_created_at DESC` for newest-first ordering. Update hurl integration tests if they depend on order.
- [ ] Add index on `fk_category` in `tb_product` — `FindProductsByCategoryID` filters on `fk_category` but PostgreSQL doesn't auto-index FK columns. Create a new migration `002_add_fk_category_index.sql` with `CREATE INDEX IF NOT EXISTS idx_product_fk_category ON tb_product (fk_category)`.
- [ ] Make health endpoint check DB connectivity — `handler/health.go` currently returns `{"status":"ok"}` without verifying the database is reachable. Inject `*sql.DB` into the health handler and call `db.PingContext(r.Context())`. Return `{"status":"ok"}` on success, `{"status":"unhealthy","error":"..."}` with `503 Service Unavailable` on failure.

### Low

- [ ] Add database repository unit tests — `internal/database/` has no `_test.go` files. Add tests using `sqlmock` or a test-container Postgres instance to verify `SaveCategory`, `FindCategoryByID`, `FindAllCategories`, `UpdateCategory`, `DeleteCategory`, and all equivalent Product methods. Focus on edge cases: empty result sets, duplicate inserts, concurrent updates.
- [ ] Add `context.Context` propagation — handlers currently don't pass request context to service/database calls. Thread `r.Context()` through service methods and use `db.QueryContext`/`db.ExecContext` in the database layer. This enables request-scoped timeouts, cancellation, and tracing.
- [ ] Extract product input validation — `CreateProduct` and `UpdateProduct` handlers repeat identical input struct definition + validation (name, price, category parsing). Extract into a shared `validateProductInput(r) (*productInput, error)` function in `helpers.go`.
- [ ] Add API versioning — prefix all routes with `/api/v1/` instead of `/api/` to support future breaking changes without disrupting existing clients.
- [ ] Add `.claude/CLAUDE.md` symlink — the root `CLAUDE.md` duplicates Claude Code's expected path at `.claude/CLAUDE.md`. Consider symlinking or moving the file to `.claude/CLAUDE.md` and keeping a reference in the root.

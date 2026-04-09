# Backlog

Task list and improvements for the codecommerceapi project.

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

## Todo

### High

- [x] Hardcode `DATABASE_URL` in CI workflow — currently uses `secrets.DATABASE_URL` which conflicts with the inline Postgres service
- [x] Add `schema_migrations` tracking table to migration runner — currently re-executes all migrations on every startup, will break with `ALTER TABLE` migrations

### Medium

- [ ] Consolidate duplicate pagination parser — `parseCategoryPaginationParams()` and `parsePaginationParams()` are identical, extract into a single shared function
- [ ] Add structured logging with `log/slog` — replace `log.Println`/`log.Fatalf` with leveled, structured output
- [ ] Add HTTP middleware — request logging, CORS, panic recovery, request ID injection
- [ ] Implement graceful shutdown — handle OS signals (`SIGINT`/`SIGTERM`) and drain connections before exit
- [ ] Check `RowsAffected()` in `UpdateCategory`/`UpdateProduct` — silent no-op if row is deleted between `FindByID` and `Update`

### Low

- [ ] Use `errors.Is()` for sentinel error comparisons in handlers — current `==` checks are fragile if errors get wrapped
- [ ] Order list queries by `ts_*_created_at DESC` instead of primary key (UUID) — current order is effectively random
- [ ] Return `[]` instead of `null` for empty paginated `Items` — initialize slices with `make([]T, 0)` before marshalling
- [ ] Fix Docker image name typo in `docker-publish.yml` — `codecomerce` → `codecommerce` (missing 'm')
- [ ] Add index on `fk_category` in `tb_product` — `FindProductsByCategoryID` filters on it but Postgres doesn't auto-index FK columns
- [ ] Make health endpoint check DB connectivity — currently returns `ok` without pinging the database

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

## Todo

- [x] Add `_test.go` files — entity creation tests and service layer tests with full CRUD + pagination mocks
- [x] Replace `db.Prepare()` with direct `db.Exec()`/`db.Query()` (statements are single-use)
- [x] Change `float64` to `int64` (cents) for price
- [x] Implement cascade delete for products when deleting a category (via `ON DELETE CASCADE` FK constraint)
- [x] Add pagination to `FindAllProducts()`, `FindAllCategories()`, and `FindProductsByCategoryID()` (query params: `page`, `per_page`)
- [x] Optimize `Delete*` in service: use `DELETE WHERE` + check rows affected instead of pre-fetch
- [x] Add embedded migration runner to execute SQL files on server startup

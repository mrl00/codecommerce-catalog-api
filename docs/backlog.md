# Backlog

Task list and improvements for the goapi project.

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

## Todo

- [ ] Add `_test.go` files (CI runs `go test ./...` but there are no tests)
- [ ] Replace `db.Prepare()` with direct `db.Exec()`/`db.Query()` (statements are single-use)
- [ ] Change `float64` to decimal type or integer (cents) for price
- [ ] Implement cascade delete for products when deleting a category
- [ ] Add pagination to `FindAllProducts()` and `FindAllCategories()`
- [ ] Optimize `Delete*` in service: use `DELETE WHERE` + check rows affected instead of pre-fetch
- [ ] Fix CI `secrets.DATABASE_URL` — conflicts with the postgres service defined in the workflow
- [ ] Standardize module name (`goapi` → `codecommerceapi`)

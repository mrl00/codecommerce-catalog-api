# codecommerceapi

A REST API service built with Go and gorilla/mux, featuring a product catalog with categories and PostgreSQL persistence.

## Project Structure

```
cmd/
└── server/
    └── main.go             # Application entry point
internal/
├── entities/
│   ├── category.go         # Category domain model
│   └── product.go          # Product domain model
├── handler/
│   ├── category_handler.go # Category HTTP handlers
│   ├── product_handler.go  # Product HTTP handlers
│   └── health.go           # Health check handler
├── router/
│   └── router.go           # Route definitions
├── database/
│   ├── category_db.go      # Category repository (raw SQL)
│   └── product_db.go       # Product repository (raw SQL)
└── service/
    ├── service.go          # Repository interfaces
    ├── category_service.go # Category business logic
    └── product_service.go  # Product business logic
```

## Getting Started

### Prerequisites

- Go 1.25+
- PostgreSQL 16

### Configuration

The database connection is configured via `DATABASE_URL` environment variable:

```bash
export DATABASE_URL="postgres://user:password@localhost:5432/codecommerce?sslmode=disable"
```

Defaults to `postgres://localhost:5432/codecommerce?sslmode=disable` if not set.

### Run

```bash
go run ./cmd/server/
```

The server listens on `:8080`.

### Build

```bash
go build -o server ./cmd/server/
```

## API Endpoints

### Health

| Method | Path        | Description  |
|--------|-------------|--------------|
| GET    | /api/health | Health check |

### Categories

| Method | Path              | Description           |
|--------|-------------------|-----------------------|
| GET    | /api/categories   | List all categories   |
| POST   | /api/categories   | Create a category     |
| GET    | /api/categories/:id | Get a category by ID  |
| PUT    | /api/categories/:id | Update a category     |
| DELETE | /api/categories/:id | Delete a category     |

### Products

| Method | Path                        | Description                    |
|--------|-----------------------------|--------------------------------|
| GET    | /api/products               | List all products              |
| POST   | /api/products               | Create a product               |
| GET    | /api/products/:id           | Get a product by ID            |
| PUT    | /api/products/:id           | Update a product               |
| DELETE | /api/products/:id           | Delete a product               |
| GET    | /api/categories/:id/products | List products by category ID  |

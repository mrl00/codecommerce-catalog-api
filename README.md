# goapi

A simple REST API service built with Go and gorilla/mux.

## Project Structure

```
cmd/
└── server/
    └── main.go         # Application entry point
internal/
├── handler/
│   └── handler.go      # HTTP request handlers
└── router/
    └── router.go       # Route definitions
```

## Getting Started

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

| Method | Path        | Description  |
|--------|-------------|--------------|
| GET    | /api/health | Health check |

## Requirements

Go 1.25+

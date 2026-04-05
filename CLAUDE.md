# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Go Project Structure: cmd/internal Pattern

This project follows Go's standard project layout:

- **`cmd/server/`** — application entry point (`main.go`). This is the only package compiled into a binary.
- **`internal/handler/`** — HTTP request handlers. Each handler function follows the `http.HandlerFunc` signature.
- **`internal/router/`** — route definitions using `gorilla/mux`. All routes are registered in `router.New()`, which returns an `*mux.Router`.

Key dependency: `github.com/gorilla/mux` (not `net/http`'s built-in `ServeMux`). New routes are added via `r.HandleFunc(path, handler).Methods("METHOD")` in `router.go`.

## Commands

```bash
go run ./cmd/server/          # run development server on :8080
go build -o server ./cmd/server/  # build binary
go mod tidy                   # sync dependencies
```

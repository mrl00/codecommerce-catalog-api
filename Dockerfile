# ─── Stage 1: build ───────────────────────────────────────────────────────────
FROM golang:1.24-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server ./cmd/server/main.go


# ─── Stage 2: runtime ─────────────────────────────────────────────────────────
FROM bitnami/golang:sha256-3db89b95bd969d4c1e198f68f6ad6dba4ea2c6a4a8a62e6aab751dc332b890f6

COPY --from=build /app/server /server

EXPOSE 8080

ENTRYPOINT ["/server"]

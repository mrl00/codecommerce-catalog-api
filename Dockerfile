# ─── Stage 1: build ───────────────────────────────────────────────────────────
FROM golang:1.25.4 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server ./cmd/server/main.go


# ─── Stage 2: runtime ─────────────────────────────────────────────────────────
FROM gcr.io/distroless/static

COPY --from=build /app/server /server

EXPOSE 8080

ENTRYPOINT ["/server"]

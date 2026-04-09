package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"codecommerceapi/internal/database"
	"codecommerceapi/internal/migrate"
	"codecommerceapi/internal/router"
	"codecommerceapi/internal/service"

	_ "github.com/lib/pq"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://localhost:5432/codecommerce?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	if err := migrate.Run(db); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	catDB := database.NewCategoryDB(db)
	prodDB := database.NewProductDB(db)

	catSvc := service.NewCategoryService(catDB)
	prodSvc := service.NewProductService(prodDB)

	r := router.New(catSvc, prodSvc)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Start server in a goroutine so it doesn't block signal handling.
	go func() {
		log.Println("starting server on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Block until we receive SIGINT or SIGTERM.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// Give in-flight requests up to 10 seconds to complete.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("server stopped")
}

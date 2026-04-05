package main

import (
	"database/sql"
	"embed"
	"log"
	"net/http"
	"os"

	"codecommerceapi/internal/database"
	"codecommerceapi/internal/migrate"
	"codecommerceapi/internal/router"
	"codecommerceapi/internal/service"

	_ "github.com/lib/pq"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

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

	if err := migrate.Run(db, migrationsFS); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	catDB := database.NewCategoryDB(db)
	prodDB := database.NewProductDB(db)

	catSvc := service.NewCategoryService(catDB)
	prodSvc := service.NewProductService(prodDB)

	r := router.New(catSvc, prodSvc)

	log.Println("starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}

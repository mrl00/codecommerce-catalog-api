package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"goapi/internal/database"
	"goapi/internal/router"

	_ "github.com/lib/pq"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://localhost:5432/goapi?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	catDB := database.NewCategoryDB(db)
	prodDB := database.NewProductDB(db)

	r := router.New(catDB, prodDB)

	log.Println("starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}

package migrate

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"path/filepath"
	"sort"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

// Run executes all embedded SQL migration files in alphabetical order,
// skipping any that have already been applied. Applied migrations are
// tracked in the schema_migrations table.
func Run(db *sql.DB) error {
	if err := ensureMigrationsTable(db); err != nil {
		return err
	}

	applied, err := appliedMigrations(db)
	if err != nil {
		return err
	}

	paths, err := fs.Glob(migrationFS, "migrations/*.sql")
	if err != nil {
		return fmt.Errorf("failed to glob migrations: %w", err)
	}
	sort.Strings(paths)

	ran := 0
	for _, path := range paths {
		fn := filepath.Base(path)
		if applied[fn] {
			slog.Info("skipping migration (already applied)", "file", fn)
			continue
		}

		content, err := migrationFS.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", path, err)
		}

		slog.Info("running migration", "file", fn)
		if err := execMigration(db, fn, string(content)); err != nil {
			return fmt.Errorf("failed to execute %s: %w", fn, err)
		}
		ran++
	}

	slog.Info("migrations completed", "applied", ran, "skipped", len(paths)-ran)
	return nil
}

// ensureMigrationsTable creates the schema_migrations table if it does not exist.
func ensureMigrationsTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			filename TEXT PRIMARY KEY,
			applied_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`)
	if err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}
	return nil
}

// appliedMigrations returns a set of filenames that have already been applied.
func appliedMigrations(db *sql.DB) (map[string]bool, error) {
	rows, err := db.Query("SELECT filename FROM schema_migrations ORDER BY filename")
	if err != nil {
		return nil, fmt.Errorf("failed to query schema_migrations: %w", err)
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var fn string
		if err := rows.Scan(&fn); err != nil {
			return nil, fmt.Errorf("failed to scan migration row: %w", err)
		}
		applied[fn] = true
	}
	return applied, rows.Err()
}

// execMigration runs a single migration inside a transaction and records it.
func execMigration(db *sql.DB, filename, content string) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec(content); err != nil {
		return err
	}

	if _, err := tx.Exec(
		"INSERT INTO schema_migrations (filename) VALUES ($1)", filename,
	); err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	return tx.Commit()
}

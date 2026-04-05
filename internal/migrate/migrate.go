package migrate

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"sort"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

// Run executes all embedded SQL migration files in alphabetical order.
func Run(db *sql.DB) error {
	paths, err := fs.Glob(migrationFS, "migrations/*.sql")
	if err != nil {
		return fmt.Errorf("failed to glob migrations: %w", err)
	}
	sort.Strings(paths)

	for _, path := range paths {
		content, err := migrationFS.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", path, err)
		}

		fn := filepath.Base(path)
		log.Printf("running migration: %s", fn)
		_, err = db.Exec(string(content))
		if err != nil {
			return fmt.Errorf("failed to execute %s: %w", path, err)
		}
	}

	log.Printf("migrations completed: %d file(s)", len(paths))
	return nil
}

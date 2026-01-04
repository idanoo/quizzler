package database

import (
	"embed"
	"log/slog"
	"path/filepath"
	"sort"
	"strings"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func Migrate() error {
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	entries, err := migrationFiles.ReadDir("migrations")
	if err != nil {
		return err
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			files = append(files, entry.Name())
		}
	}
	sort.Strings(files)

	for _, file := range files {
		version := strings.TrimSuffix(file, filepath.Ext(file))

		var count int
		err := DB.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = ?", version).Scan(&count)
		if err != nil {
			return err
		}
		if count > 0 {
			continue
		}

		content, err := migrationFiles.ReadFile("migrations/" + file)
		if err != nil {
			return err
		}

		slog.Info("Running migration", "version", version)
		_, err = DB.Exec(string(content))
		if err != nil {
			slog.Error("Migration failed", "version", version, "error", err)
			return err
		}

		// Record migration
		_, err = DB.Exec("INSERT INTO schema_migrations (version) VALUES (?)", version)
		if err != nil {
			return err
		}
	}

	slog.Info("Migrations completed")
	return nil
}

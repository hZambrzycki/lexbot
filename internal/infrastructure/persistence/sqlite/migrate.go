package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
)

func RunMigrations(db *sql.DB, path string) error {
	if err := ensureSchemaMigrationsTable(db); err != nil {
		return err
	}

	migrationName := filepath.Base(path)

	applied, err := isMigrationApplied(db, migrationName)
	if err != nil {
		return err
	}
	if applied {
		return nil
	}

	sqlBytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(string(sqlBytes)); err != nil {
		return fmt.Errorf("apply migration %s: %w", migrationName, err)
	}

	if _, err := tx.Exec(
		`INSERT INTO schema_migrations (name) VALUES (?)`,
		migrationName,
	); err != nil {
		return fmt.Errorf("record migration %s: %w", migrationName, err)
	}

	return tx.Commit()
}

func ensureSchemaMigrationsTable(db *sql.DB) error {
	const query = `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			name TEXT PRIMARY KEY,
			applied_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := db.Exec(query)
	return err
}

func isMigrationApplied(db *sql.DB, name string) (bool, error) {
	const query = `
		SELECT 1
		FROM schema_migrations
		WHERE name = ?
		LIMIT 1
	`

	var exists int
	err := db.QueryRow(query, name).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

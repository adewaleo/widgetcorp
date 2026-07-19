package db

import (
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5" // registers the pgx5:// driver
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

// migrationFiles embeds every .sql file in the migrations/ directory into the
// compiled binary, so the server is self-contained and needs no external files
// at runtime.
//
//go:embed migrations/*.sql
var migrationFiles embed.FS

// RunMigrations applies all pending "up" migrations against the database at dsn.
// It is safe to call on every startup: if the schema is already current it
// returns nil (migrate.ErrNoChange is treated as success).
func RunMigrations(dsn string) error {
	// Source driver: read migration SQL from the embedded filesystem.
	source, err := iofs.New(migrationFiles, "migrations")
	if err != nil {
		return fmt.Errorf("load embedded migrations: %w", err)
	}

	// golang-migrate wants a "pgx5://" URL to select the pgx/v5 database driver.
	// Our dsn is a normal "postgres://..." string, so swap the scheme.
	migrateURL := "pgx5://" + trimScheme(dsn)

	m, err := migrate.NewWithSourceInstance("iofs", source, migrateURL)
	if err != nil {
		return fmt.Errorf("init migrate: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("apply migrations: %w", err)
	}
	return nil
}

// trimScheme strips a leading "postgres://" or "postgresql://" so we can prefix
// the pgx5 scheme instead.
func trimScheme(dsn string) string {
	for _, p := range []string{"postgres://", "postgresql://"} {
		if len(dsn) >= len(p) && dsn[:len(p)] == p {
			return dsn[len(p):]
		}
	}
	return dsn
}

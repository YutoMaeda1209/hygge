package model

import (
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

// Schema and seed data are managed as versioned SQL migrations; GORM tags are not used to create tables.
//
//go:embed migrations/*.sql
var migrationFiles embed.FS

func runMigrations() error {
	source, err := iofs.New(migrationFiles, "migrations")
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	sqlDb, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}
	driver, err := pgx.WithInstance(sqlDb, &pgx.Config{})
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", source, "pgx5", driver)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to migrate database: %w", err)
	}
	return nil
}

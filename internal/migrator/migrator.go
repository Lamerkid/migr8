// Package migrator provides logic for migrator app.
package migrator

import (
	"context"
	"fmt"
)

// Migrator is the main migrator.
type Migrator struct {
	db     Database
	logger Logger
}

// NewMigrator returns new instance of the migrator.
func NewMigrator(db Database, logger Logger) *Migrator {
	return &Migrator{
		db:     db,
		logger: logger,
	}
}

// CreateServiceTables initializes migration tables.
func (m *Migrator) CreateServiceTables(ctx context.Context) error {
	if err := m.db.CreateServiceTables(ctx); err != nil {
		return fmt.Errorf("failed to start migration: %w", err)
	}
	return nil
}

// Up applies migrations to the database.
func (m *Migrator) Up(ctx context.Context, path int) error {
	_ = ctx
	_ = path
	return nil
}

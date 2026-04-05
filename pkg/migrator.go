// Package migr8 provides functions to use.
package migr8

import (
	"context"
	"database/sql"

	"github.com/Lamerkid/migr8/internal/core"
)

// Migrator is the public API.
type Migrator struct {
	core *core.Migrator
}

// New creates a new migrator instance.
func New(db *sql.DB, migrationPath string, opts ...core.Option) *Migrator {
	m := core.NewMigrator(db, migrationPath, opts...)

	return &Migrator{core: m}
}

// NewWithDB creates a new migrator instance with custom database.
func NewWithDB(db core.Database, migrationPath string, opts ...core.Option) *Migrator {
	m := core.NewMigratorWithDB(db, migrationPath, opts...)

	return &Migrator{core: m}
}

// Up applies all pending migrations.
func (m *Migrator) Up(ctx context.Context) error {
	return m.core.Up(ctx)
}

// Down rolls back the last migration.
func (m *Migrator) Down(ctx context.Context) error {
	return m.core.Down(ctx)
}

// Redo reapplies the last migration.
func (m *Migrator) Redo(ctx context.Context) error {
	return m.core.Redo(ctx)
}

// Status returns the current migration status.
func (m *Migrator) Status(ctx context.Context) error {
	return m.core.Status(ctx)
}

// DBversion returns the current migration status.
func (m *Migrator) DBversion(ctx context.Context) error {
	return m.core.DBversion(ctx)
}

// Close closes the database connection.
func (m *Migrator) Close() error {
	return m.core.Close()
}

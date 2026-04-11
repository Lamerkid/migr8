// Package migration provides applying go migrations to database.
package migration

import (
	"context"
	"database/sql"

	migr8 "github.com/Lamerkid/migr8/pkg"
)

func init() {
	migr8.RegisterMigration(20260406204250, &TestGoMigration{})
}

// TestGoMigration struct represent up and down command for single file.
type TestGoMigration struct{}

// Up command lets apply migrations.
func (m *TestGoMigration) Up(ctx context.Context, tx *sql.Tx) error {
	// Rewrite with desired functions.
	query := `CREATE TABLE IF NOT EXISTS example(id SERIAL PRIMARY KEY)`
	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}
	return nil
}

// Down command lets rollback migration.
func (m *TestGoMigration) Down(ctx context.Context, tx *sql.Tx) error {
	// Rewrite with desired functions.
	query := `DROP TABLE IF EXISTS example`
	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}
	return nil
}

package core

import (
	"context"
	"database/sql"
)

func init() {
	RegisterMigration(123, &MockGoMigration{})
}

// MockGoMigration struct represent up and down command for single file.
type MockGoMigration struct{}

// Up command lets apply migrations.
func (m *MockGoMigration) Up(ctx context.Context, tx *sql.Tx) error {
	query := `CREATE TABLE IF NOT EXISTS test (id UUID PRIMARY KEY)`
	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}
	return nil
}

// Down command lets rollback migration.
func (m *MockGoMigration) Down(ctx context.Context, tx *sql.Tx) error {
	query := `DROP TABLE IF EXISTS test`
	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}
	return nil
}

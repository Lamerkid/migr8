package migrator

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
	if _, err := tx.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS test (id UUID PRIMARY KEY);`); err != nil {
		return err
	}
	return nil
}

// Down command lets rollback migration.
func (m *MockGoMigration) Down(ctx context.Context, tx *sql.Tx) error {
	if _, err := tx.ExecContext(ctx, `DROP TABLE IF EXISTS test`); err != nil {
		return err
	}
	return nil
}

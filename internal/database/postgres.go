// Package postgres provides functions for Postgres database.
package postgres

import (
	"context"
	"database/sql"
)

// Database implements Database using standard *sql.DB.
type Database struct {
	db *sql.DB
}

// NewDatabase creates default database instance.
func NewDatabase(db *sql.DB) *Database {
	return &Database{db: db}
}

// Close the connection to the database.
func (d *Database) Close() error {
	return d.db.Close()
}

// BeginTx starts a new transaction.
func (d *Database) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return d.db.BeginTx(ctx, opts)
}

// QueryRowContext executes a query that returns a single row.
func (d *Database) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return d.db.QueryRowContext(ctx, query, args...)
}

// CreateServiceTable for migrations.
func (d *Database) CreateServiceTable(ctx context.Context) error {
	query := `CREATE TABLE IF NOT EXISTS changelog (
		version			BIGINT PRIMARY KEY,
		name				TEXT,
		applied_at	TIMESTAMP WITH TIME ZONE DEFAULT NOW()
  )`

	_, err := d.db.ExecContext(ctx, query)
	return err
}

// GetAppliedVersions returns all applied migration versions.
func (d *Database) GetAppliedVersions(ctx context.Context) (map[int64]bool, error) {
	query := "SELECT version FROM changelog ORDER BY version"
	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[int64]bool)
	for rows.Next() {
		var version int64
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}
	return applied, rows.Err()
}

// GetLatestVersion returns the most recent applied version.
func (d *Database) GetLatestVersion(ctx context.Context) (int64, error) {
	query := "SELECT COALESCE(MAX(version), 0) FROM changelog"
	var version int64
	err := d.db.QueryRowContext(ctx, query).Scan(&version)
	return version, err
}

// InsertApplied records a migration as applied.
func (d *Database) InsertApplied(ctx context.Context, tx *sql.Tx, version int64, name string) error {
	query := "INSERT INTO changelog (version, name) VALUES ($1, $2)"
	_, err := tx.ExecContext(ctx, query, version, name)
	return err
}

// DeleteApplied removes a migration record.
func (d *Database) DeleteApplied(ctx context.Context, tx *sql.Tx, version int64) error {
	query := "DELETE FROM changelog WHERE version = $1"
	_, err := tx.ExecContext(ctx, query, version)
	return err
}

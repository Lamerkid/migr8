// Package postgres provides functions for postgres database.
package postgres

import (
	"context"
	"database/sql"
	"fmt"

	// Pgx driver.
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Database is the main database.
type Database struct {
	DB *sql.DB
}

// NewDatabase returns new instance of the database.
func NewDatabase() *Database {
	return &Database{}
}

// Connect to postgres database.
func (d *Database) Connect(ctx context.Context, dsn string) (err error) {
	d.DB, err = sql.Open("pgx", dsn)
	if err != nil {
		return err
	}

	return d.DB.PingContext(ctx)
}

// Close the connection to the database.
func (d *Database) Close() error {
	return d.DB.Close()
}

// BeginTx starts a new transaction.
func (d *Database) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	if d.DB == nil {
		return nil, fmt.Errorf("database not connected")
	}
	return d.DB.BeginTx(ctx, opts)
}

// QueryRowContext executes a query that returns a single row.
func (d *Database) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	if d.DB == nil {
		return &sql.Row{}
	}
	return d.DB.QueryRowContext(ctx, query, args...)
}

// CreateServiceTable for migrations.
func (d *Database) CreateServiceTable(ctx context.Context) error {
	logTable := `CREATE TABLE IF NOT EXISTS changelog (
		version			BIGINT PRIMARY KEY,
		name				TEXT,
		applied_at	TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	)`

	_, err := d.DB.ExecContext(ctx, logTable)
	if err != nil {
		return fmt.Errorf("error creating migration log table: %w", err)
	}

	return nil
}

// GetAppliedVersions returns all applied migration versions.
func (d *Database) GetAppliedVersions(ctx context.Context) (map[int64]bool, error) {
	rows, err := d.DB.QueryContext(ctx, "SELECT version FROM changelog ORDER BY version")
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
	var version int64
	err := d.DB.QueryRowContext(ctx, "SELECT COALESCE(MAX(version), 0) FROM changelog").Scan(&version)
	return version, err
}

// InsertApplied records a migration as applied.
func (d *Database) InsertApplied(ctx context.Context, tx *sql.Tx, version int64, name string) error {
	query := `INSERT INTO changelog (version, name) VALUES ($1, $2)`
	_, err := tx.ExecContext(ctx, query, version, name)
	return err
}

// DeleteApplied removes a migration record.
func (d *Database) DeleteApplied(ctx context.Context, tx *sql.Tx, version int64) error {
	query := `DELETE FROM changelog WHERE version = $1`
	_, err := tx.ExecContext(ctx, query, version)
	return err
}

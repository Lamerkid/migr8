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

// CreateServiceTables for migrations.
func (d *Database) CreateServiceTables(ctx context.Context) error {
	logTable := `CREATE TABLE IF NOT EXISTS changelog (
		id 					VARCHAR(255) PRIMARY KEY,
		applied_at	TIMESTAMP WITH TIME ZONE,
		status			VARCHAR(10),
		checksum		TEXT,
		UNIQUE(id)
	)`

	lockTable := `CREATE TABLE IF NOT EXISTS migrationlocks (
		id 					INTEGER PRIMARY KEY,
		locked_by		TEXT NOT NULL,
		locked_at			TIMESTAMP WITH TIME ZONE,
		UNIQUE(id)
	)`

	_, err := d.DB.ExecContext(ctx, logTable)
	if err != nil {
		return fmt.Errorf("error creating migration log table: %w", err)
	}

	_, err = d.DB.ExecContext(ctx, lockTable)
	if err != nil {
		return fmt.Errorf("error creating migration lock table: %w", err)
	}

	return nil
}

package migrator

import (
	"context"
	"database/sql"
)

// Logger interface implements custom logger.
type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

// Database interface implements database functions.
type Database interface {
	Close() error
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	CreateServiceTable(ctx context.Context) error
	GetAppliedVersions(ctx context.Context) (map[int64]bool, error)
	InsertApplied(ctx context.Context, tx *sql.Tx, version int64, name string) error
	DeleteApplied(ctx context.Context, tx *sql.Tx, version int64) error
	GetLatestVersion(ctx context.Context) (int64, error)
}

// GoMigration is the interface that Go migration files must implement.
type GoMigration interface {
	Up(ctx context.Context, tx *sql.Tx) error
	Down(ctx context.Context, tx *sql.Tx) error
}

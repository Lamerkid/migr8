package migrator

import "context"

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
	CreateServiceTables(ctx context.Context) error
}

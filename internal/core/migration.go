package core

import (
	"context"
	"database/sql"
	"sync"
)

var (
	migrations     = make(map[int64]GoMigration)
	migrationsLock sync.RWMutex
)

// MigrationType represents type of passed migration.
type MigrationType string

const (
	// TypeSQL represents files with sql migrations.
	TypeSQL MigrationType = "sql"
	// TypeGo represents files with go migrations.
	TypeGo MigrationType = "go"
)

// Migration represents a single migration file.
type Migration struct {
	Version  int64
	Name     string
	Type     MigrationType
	UpSQL    string
	DownSQL  string
	UpFunc   func(context.Context, *sql.Tx) error
	DownFunc func(context.Context, *sql.Tx) error
}

// RegisterMigration registers a Go migration.
func RegisterMigration(version int64, migration GoMigration) {
	migrationsLock.Lock()
	defer migrationsLock.Unlock()

	if _, exists := migrations[version]; exists {
		return
	}

	migrations[version] = migration
}

// GetRegisteredMigration retrieves a registered Go migration.
func GetRegisteredMigration(version int64) (GoMigration, bool) {
	migrationsLock.RLock()
	defer migrationsLock.RUnlock()

	mig, exists := migrations[version]
	return mig, exists
}

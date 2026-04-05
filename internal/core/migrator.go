// Package core provides core logic for migrator.
package core

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	postgres "github.com/Lamerkid/migr8/internal/database"
	"github.com/Lamerkid/migr8/internal/logger"
)

// Migrator is the main migrator.
type Migrator struct {
	source       *MigrationSource
	db           Database
	logger       Logger
	advisoryLock AdvisoryLock
	mu           sync.Mutex
}

// Option pattern for flexible configuration.
type Option func(*Migrator)

// WithLogger allows to pass custom logger to migrator.
func WithLogger(logg Logger) Option {
	return func(m *Migrator) {
		m.logger = logg
		m.source.Logger = logg
	}
}

// NewMigrator returns new instance of the migrator.
func NewMigrator(db *sql.DB, migrationPath string, opts ...Option) *Migrator {
	defaultDB := postgres.NewDatabase(db)
	defaultLogger := logger.NewLogger("INFO")

	m := &Migrator{
		source:       &MigrationSource{Logger: defaultLogger, Path: migrationPath},
		db:           defaultDB,
		logger:       defaultLogger,
		advisoryLock: AdvisoryLock{DB: defaultDB},
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

// NewMigratorWithDB creates a migrator with custom Database implementation.
func NewMigratorWithDB(db Database, migrationsDir string, opts ...Option) *Migrator {
	defaultLogger := logger.NewLogger("INFO")

	m := &Migrator{
		source:       &MigrationSource{Logger: defaultLogger, Path: migrationsDir},
		db:           db,
		logger:       defaultLogger,
		advisoryLock: AdvisoryLock{DB: db},
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

// Close closes database connection.
func (m *Migrator) Close() error {
	if m.db != nil {
		return m.db.Close()
	}
	return nil
}

// Up applies all migrations.
func (m *Migrator) Up(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	acquired, err := m.advisoryLock.Acquire(ctx, 123456789)
	if err != nil {
		return err
	}
	if !acquired {
		return fmt.Errorf("another migration is in progress")
	}
	defer m.advisoryLock.Release(ctx, 123456789)

	if err := m.ensureTables(ctx); err != nil {
		return fmt.Errorf("failed to ensure migration table: %w", err)
	}

	m.logger.Info("loading migrations")
	migrations, err := m.source.LoadMigrations()
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	m.logger.Info("checking for applied versions")
	applied, err := m.db.GetAppliedVersions(ctx)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Apply only pending migrations.
	var pending []*Migration
	for _, mig := range migrations {
		if !applied[mig.Version] {
			pending = append(pending, mig)
		}
	}

	if len(pending) == 0 {
		return nil
	}

	for i := range pending {
		mig := pending[i]
		m.logger.Info("applying migration %d: %s", mig.Version, mig.Name)

		if err := m.applyUp(ctx, mig); err != nil {
			return fmt.Errorf("failed to apply migration %d: %w", mig.Version, err)
		}
	}

	return nil
}

// Down rolls back applied migrations.
func (m *Migrator) Down(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	acquired, err := m.advisoryLock.Acquire(ctx, 123456789)
	if err != nil {
		return err
	}
	if !acquired {
		return fmt.Errorf("another migration is in progress")
	}
	defer m.advisoryLock.Release(ctx, 123456789)

	if err := m.ensureTables(ctx); err != nil {
		return fmt.Errorf("failed to ensure migration table: %w", err)
	}

	m.logger.Info("checking for latest applied version")
	version, err := m.db.GetLatestVersion(ctx)
	if err != nil {
		return err
	}
	if version == 0 {
		return nil
	}

	m.logger.Info("loading latests migration")
	mig, err := m.loadMigration(version)
	if err != nil {
		return err
	}

	m.logger.Info("rolling back %d\n", version)
	return m.applyDown(ctx, mig)
}

// Redo reapplies the last migration.
func (m *Migrator) Redo(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	acquired, err := m.advisoryLock.Acquire(ctx, 123456789)
	if err != nil {
		return err
	}
	if !acquired {
		return fmt.Errorf("another migration is in progress")
	}
	defer m.advisoryLock.Release(ctx, 123456789)

	if err := m.ensureTables(ctx); err != nil {
		return fmt.Errorf("failed to ensure migration table: %w", err)
	}

	m.logger.Info("checking for latest applied version")
	version, err := m.db.GetLatestVersion(ctx)
	if err != nil {
		return err
	}
	if version == 0 {
		return nil
	}

	m.logger.Info("loading latests migration")
	mig, err := m.loadMigration(version)
	if err != nil {
		return err
	}

	m.logger.Info("rolling back %d for redo\n", version)
	if err := m.applyDown(ctx, mig); err != nil {
		return err
	}

	m.logger.Info("reapplying %d\n", version)
	return m.applyUp(ctx, mig)
}

// Status shows current migration state.
func (m *Migrator) Status(ctx context.Context) error {
	if err := m.ensureTables(ctx); err != nil {
		return fmt.Errorf("failed to ensure migration table: %w", err)
	}

	m.logger.Info("loading migrations")
	migrations, err := m.source.LoadMigrations()
	if err != nil {
		return err
	}

	m.logger.Info("checking for applied versions")
	applied, err := m.db.GetAppliedVersions(ctx)
	if err != nil {
		return err
	}

	fmt.Println("\nMigration Status:")

	for _, mig := range migrations {
		marker := "-"
		if applied[mig.Version] {
			marker = "+"
		}
		fmt.Printf("%s %d %s\n", marker, mig.Version, mig.Name)
	}

	fmt.Printf("\nTotal: %d, Applied: %d, Pending: %d\n",
		len(migrations),
		len(applied),
		len(migrations)-len(applied))

	return nil
}

// DBversion shows current migration state.
func (m *Migrator) DBversion(ctx context.Context) error {
	if err := m.ensureTables(ctx); err != nil {
		return fmt.Errorf("failed to ensure migration table: %w", err)
	}

	m.logger.Info("checking for latest applied version")
	version, err := m.db.GetLatestVersion(ctx)
	if err != nil {
		return err
	}
	if version == 0 {
		return nil
	}

	m.logger.Info("loading latests migration")
	mig, err := m.loadMigration(version)
	if err != nil {
		return err
	}

	fmt.Println("\nDatabase version:")
	fmt.Printf("%d - applied migration: %s\n", mig.Version, mig.Name)

	return nil
}

func (m *Migrator) ensureTables(ctx context.Context) error {
	return m.db.CreateServiceTable(ctx)
}

func (m *Migrator) applyUp(ctx context.Context, mig *Migration) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	switch mig.Type {
	case TypeSQL:
		if mig.UpSQL == "" {
			return fmt.Errorf("migration %d has no up SQL", mig.Version)
		}
		if _, err := tx.ExecContext(ctx, mig.UpSQL); err != nil {
			return fmt.Errorf("up failed: %w", err)
		}

	case TypeGo:
		if mig.UpFunc == nil {
			return fmt.Errorf("go migration %d has no Up function", mig.Version)
		}
		if err := mig.UpFunc(ctx, tx); err != nil {
			return fmt.Errorf("go up failed: %w", err)
		}

	default:
		return fmt.Errorf("unknown migration type: %s", mig.Type)
	}

	if err := m.db.InsertApplied(ctx, tx, mig.Version, mig.Name); err != nil {
		return err
	}

	return tx.Commit()
}

func (m *Migrator) applyDown(ctx context.Context, mig *Migration) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	switch mig.Type {
	case TypeSQL:
		if mig.DownSQL == "" {
			return fmt.Errorf("migration %d has no down SQL", mig.Version)
		}
		if _, err := tx.ExecContext(ctx, mig.DownSQL); err != nil {
			return fmt.Errorf("down failed: %w", err)
		}

	case TypeGo:
		if mig.DownFunc == nil {
			return fmt.Errorf("go migration %d has no Down function", mig.Version)
		}
		if err := mig.DownFunc(ctx, tx); err != nil {
			return fmt.Errorf("go up failed: %w", err)
		}

	default:
		return fmt.Errorf("unknown migration type: %s", mig.Type)
	}

	if err := m.db.DeleteApplied(ctx, tx, mig.Version); err != nil {
		return err
	}

	return tx.Commit()
}

func (m *Migrator) loadMigration(version int64) (*Migration, error) {
	migrations, err := m.source.LoadMigrations()
	if err != nil {
		return nil, err
	}

	for _, mig := range migrations {
		if mig.Version == version {
			return mig, nil
		}
	}

	return nil, fmt.Errorf("migration %d not found", version)
}

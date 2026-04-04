package cli

import (
	"context"
	"strings"

	"github.com/Lamerkid/migr8/internal/config"
	postgres "github.com/Lamerkid/migr8/internal/database"
	"github.com/Lamerkid/migr8/internal/logger"
	"github.com/Lamerkid/migr8/internal/migrator"
)

var exampleSQL = `-- +migr8:up
CREATE TABLE IF EXISTS example (
	id SERIAL PRIMARY KEY
);

-- +migr8:down
DROP TABLE IF EXISTS example;
`

var exampleGo = `// Package migrations provides usage of go migrations to database.
package migrations

import (
	"context"
	"database/sql"

	"github.com/Lamerkid/migr8/internal/migrator"
)

func init() {
	migrator.RegisterMigration(%s, &%s{})
}

// %s struct represent up and down command for single file.
type %s struct{}

// Up command lets apply migrations.
func (m *%s) Up(ctx context.Context, tx *sql.Tx) error {
	// Write your up migration here
	return nil
}

// Down command lets rollback migration.
func (m *%s) Down(ctx context.Context, tx *sql.Tx) error {
	// Write your down migration here
	return nil
}
`

func createMigratorInstance(ctx context.Context, flags map[string]string) (*migrator.Migrator, error) {
	cfg, err := config.BuildFromFlags(flags)
	if err != nil {
		return nil, err
	}

	logger := logger.NewLogger(cfg.Logger.Level)

	db := postgres.NewDatabase()
	if err := db.Connect(ctx, cfg.Database.DSN); err != nil {
		return nil, err
	}

	return migrator.NewMigrator(db, logger, cfg.Migration.Dir), nil
}

func toCamelCase(s string) string {
	parts := strings.Split(s, "_")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return strings.Join(parts, "")
}

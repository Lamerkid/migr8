package cli

import (
	"database/sql"
	"strings"

	"github.com/Lamerkid/migr8/internal/config"
	migr8 "github.com/Lamerkid/migr8/pkg"
	_ "github.com/jackc/pgx/v5/stdlib" // Pgx driver for Postgres.
)

var exampleSQL = `-- +migr8:up

CREATE TABLE IF NOT EXISTS example (
	id SERIAL PRIMARY KEY
);

-- +migr8:down

DROP TABLE IF EXISTS example;
`

var exampleGo = `// Package migration provides applying go migrations to database.
package migration

import (
	"context"
	"database/sql"

	migr8 "github.com/Lamerkid/migr8/pkg"
)

func init() {
	migr8.RegisterMigration(%s, &%s{})
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

func createMigratorInstance(flags map[string]string) (*migr8.Migrator, error) {
	cfg, err := config.BuildFromFlags(flags)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("pgx", cfg.Database.DSN)
	if err != nil {
		return nil, err
	}

	return migr8.New(db, cfg.Migration.Dir), nil
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

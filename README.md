# migr8

[![CI tests](https://github.com/Lamerkid/migr8/actions/workflows/tests.yml/badge.svg)](https://github.com/Lamerkid/migr8/actions/workflows/tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/Lamerkid/migr8)](https://goreportcard.com/report/github.com/Lamerkid/migr8)

Migr8 is a PostgreSQL database migration tool. It can be used as library or CLI.
Allows to manage your PostgreSQL database schema using SQL changes or Go functions.

## CLI usage

## Commands

Supported commands:

- [version](#version)
- [create](#create)
- [up](#up)
- [down](#down)
- [redo](#redo)
- [status](#status)
- [dbversion](#dbversion)

### version

Prints a version of a CLI app.

### create

Creates a new migration file.

```sh
$ migr8 create add_table
created new file: 20260406184939_add_table.sql
```

### up

Apply all available migrations.

```sh
$ migr8 up
[INFO] applying 20260406184939: 20260406184939_add_table.sql
[INFO] applying 20260406195043: 20260406195043_add_column.go
[INFO] applied successfully
```

### down

Roll back a single migration from current version.

```sh
$ migr8 down
[INFO] successfully rolled back 20260406195043
```

### redo

Reapply a single migration from current version.

```sh
$ migr8 create add_table
[INFO] successfully reapplied 20260406184939
```

### status

Show status of applied and pending versions.

```sh
$ migr8 status

Migration Status:
+ 20260406184939: 20260406184939_add_table.sql
- 20260406195043: 20260406195043_add_column.go

Total: 2, Applied: 1, Pending: 1
```

### dbversion

Show version of the database (last applied migration).

```sh
$ migr8 dbversion

Database version:
20260406184939 - applied migration: 20260406184939_add_table.sql
```

## Options

Example usage:

```sh
migr8 up -cfg=./path/config.json

migr8 -log INFO down

migr8 create add_column -type=go
```

Supported options:

- **-cfg** - path to JSON config file.
- **-log** - select a preferred logger level (**only UPPERCASE**) (DEBUG/INFO/WARN/ERROR)
- **-dir** - pass a path to a directory with migrations.
- **-dsn** - specify DSN connection to PostgreSQL database.
- **-type** - choose a type of created migration (**only lowercase**) (sql/go)

## Environment variables

Instead of using options (*dir*, *dsn*), you can set the following environment variables:

```sh
export M8_DSN postgres://postgres:postgres@localhost:5432/migr8_test?sslmode=disable
export M8_DIR ../../testsqlmigration
```

## Default values

For options `-log` and `-type` declared default values:

- logger lever    - INFO
- migration type  - sql

If you need different value, you need to explicitly declare a preference using options.

## Migrations

Migrations file types to same database may vary.
Migrations must begin with numeric value and not end with `*_test.go`:
  CORRECT -   `20260315230925_add_column.sql`
  WRONG   -   `add_column_20260301.sql`

Note that migr8 sorts migrations by first numeric value.
If you create `2026_*.sql` (first) and `0001_*.sql` (second) sorted:
  LATEST - `2026_*.sql`

### SQL migrations

Sample SQL migration:

```SQL
-- +migr8:up

CREATE TABLE IF NOT EXISTS example (
 id SERIAL PRIMARY KEY
);

-- +migr8:down

DROP TABLE IF EXISTS example;

```

SQL migrations use custom comments `-- +migr8:up` and `-- +migr8:down`
This comments separate up and down statements inside one file.

For both sections there are checks:

- up, redo    checks for Up section in migration files
- down, redo  checks for Down section in migration files

### Go migrations

To use Go migrations you need:

1. Move **main.go** into your `cmd/` directory

2. Import migrations directory from your custom [cmd/main.go](cmd/migr8/main.go):

   ```go
   import (
       // Invoke init() functions within migrations pkg.
       _ "github.com/Lamerkid/migr8/example"
   )
   ```

Sample Go migration:

```go
package migration

import (
  "context"
  "database/sql"

  migr8 "github.com/Lamerkid/migr8/pkg"
)

func init() {
  migr8.RegisterMigration(20260406195043, &YourName{})
}

// YourName struct represent up and down command for single file.
type YourName struct{}

// Up command lets apply migrations.
func (m *YourName) Up(ctx context.Context, tx *sql.Tx) error {
  query := `CREATE TABLE IF NOT EXISTS test (id UUID PRIMARY KEY)`
  if _, err := tx.ExecContext(ctx, query); err != nil {
    return err
  }
  return nil
}

// Down command lets rollback migration.
func (m *YourName) Down(ctx context.Context, tx *sql.Tx) error {
  query := `DROP TABLE IF EXISTS test`
  if _, err := tx.ExecContext(ctx, query); err != nil {
    return err
  }
  return nil
}
```

Go migrations generates a RegisterMigration command within itself.

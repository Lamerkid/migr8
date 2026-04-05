package cli

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Lamerkid/migr8/internal/config"
)

type command struct {
	Name        string
	Description string
	Action      func(context.Context, []string, map[string]string) error
}

func (a *App) addCommand(cmd *command) {
	a.commands[cmd.Name] = cmd
}

// RegisterCommands adds commands to the CLI app.
func RegisterCommands(app *App) {
	app.addCommand(&command{
		Name:        "version",
		Description: "Prints a version of an app",
		Action: func(_ context.Context, _ []string, _ map[string]string) error {
			fmt.Printf("migr8 version: %s\n", app.Version)
			return nil
		},
	})

	app.addCommand(&command{
		Name:        "create",
		Description: "Creates a migration file",
		Action: func(_ context.Context, args []string, flags map[string]string) error {
			if len(args) == 0 {
				return fmt.Errorf("migration name required")
			}

			cfg, err := config.BuildFromFlags(flags)
			if err != nil {
				return err
			}

			timeStamp := time.Now().Format("20060102150405")
			fileName := fmt.Sprintf("%s_%s", timeStamp, args[0])
			var content string

			if cfg.Migration.Type == "sql" {
				fileName += ".sql"
				content = exampleSQL
			}
			if cfg.Migration.Type == "go" {
				fileName += ".go"
				content = fmt.Sprintf(exampleGo,
					timeStamp,
					toCamelCase(args[0]),
					toCamelCase(args[0]),
					toCamelCase(args[0]),
					toCamelCase(args[0]),
					toCamelCase(args[0]))
			}

			fullPath := filepath.Join(cfg.Migration.Dir, fileName)

			// #nosec G306
			if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
				return err
			}

			log.Printf("created new file: %s", fileName)

			return nil
		},
	})

	app.addCommand(&command{
		Name:        "up",
		Description: "Apply migrations",
		Action: func(ctx context.Context, _ []string, flags map[string]string) error {
			m, err := createMigratorInstance(flags)
			if err != nil {
				return err
			}
			defer m.Close()

			fmt.Printf("Applying migration(s)...\n")
			return m.Up(ctx)
		},
	})

	app.addCommand(&command{
		Name:        "down",
		Description: "Rollback last migration",
		Action: func(ctx context.Context, _ []string, flags map[string]string) error {
			m, err := createMigratorInstance(flags)
			if err != nil {
				return err
			}
			defer m.Close()

			fmt.Printf("Rollback previous migration...\n")
			return m.Down(ctx)
		},
	})

	app.addCommand(&command{
		Name:        "redo",
		Description: "Redo last migration",
		Action: func(ctx context.Context, _ []string, flags map[string]string) error {
			m, err := createMigratorInstance(flags)
			if err != nil {
				return err
			}
			defer m.Close()

			fmt.Printf("Reapplying previous migration...\n")
			return m.Redo(ctx)
		},
	})

	app.addCommand(&command{
		Name:        "status",
		Description: "Migrations status",
		Action: func(ctx context.Context, _ []string, flags map[string]string) error {
			m, err := createMigratorInstance(flags)
			if err != nil {
				return err
			}
			defer m.Close()

			return m.Status(ctx)
		},
	})

	app.addCommand(&command{
		Name:        "dbversion",
		Description: "Database version (last applied migration)",
		Action: func(ctx context.Context, _ []string, flags map[string]string) error {
			m, err := createMigratorInstance(flags)
			if err != nil {
				return err
			}
			defer m.Close()

			return m.DBversion(ctx)
		},
	})
}

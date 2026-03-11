package cli

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Lamerkid/migr8/internal/config"
	postgres "github.com/Lamerkid/migr8/internal/database"
	"github.com/Lamerkid/migr8/internal/logger"
	"github.com/Lamerkid/migr8/internal/migrator"
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
			fileName := fmt.Sprintf("%s_%s.sql", timeStamp, args[0])

			fullPath := filepath.Join(cfg.Migration.Path, fileName)

			if err := os.WriteFile(fullPath, []byte("-- +mig8:up\n\n-- +migr8:down\n"), 0o600); err != nil {
				return err
			}

			log.Printf("created new file: %s", fileName)

			return nil
		},
	})

	app.addCommand(&command{
		Name:        "up",
		Description: "Apply migrations",
		Action: func(ctx context.Context, args []string, flags map[string]string) error {
			m, err := createMigratorInstance(ctx, flags)
			if err != nil {
				return err
			}

			_ = args
			_ = m

			return nil
		},
	})
}

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

	return migrator.NewMigrator(db, logger), nil
}

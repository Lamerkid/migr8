//go:build integration

package integration

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	// Pgx driver.
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
)

func TestIntegration(t *testing.T) {
	ctx := context.Background()

	dsn := os.Getenv("M8_DSN")

	waitForPostgres(ctx, t, dsn)

	cmd := exec.CommandContext(ctx, "go", "build", "-o", "../../bin/", "../../cmd/migr8")
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to build binary: %v", err)
	}

	t.Run("create migration", func(t *testing.T) {
		cmd := exec.CommandContext(ctx, "../../bin/migr8", "create", "test_migration", "-dir", "../../migrations")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("create failed: %v\n%s", err, output)
		}

		// Check if file was created.
		files, _ := filepath.Glob(filepath.Join("../../migrations", "*.sql"))
		if len(files) == 0 {
			t.Error("no migration file created")
		}
	})

	require.NotEmpty(t, dsn)
}

func waitForPostgres(ctx context.Context, t *testing.T, dsn string) {
	t.Helper()

	maxRetries := 10
	for i := range maxRetries {
		db, err := sql.Open("pgx", dsn)
		if err != nil {
			time.Sleep(time.Second)
			continue
		}

		err = db.PingContext(ctx)
		_ = db.Close()

		if err == nil {
			fmt.Println("PostgreSQL is ready")
			return
		}

		fmt.Printf("Waiting for PostgreSQL... (%d/%d)\n", i+1, maxRetries)
		time.Sleep(time.Second)
	}

	t.Fatal("PostgreSQL not ready after 10 seconds")
}

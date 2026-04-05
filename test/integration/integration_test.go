//go:build integration

package integration

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"testing"
	"time"

	// Pgx driver.
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
)

var (
	binaryPath = "../../bin/migr8"
	dsn        = os.Getenv("M8_DSN")
	dir        = os.Getenv("M8_DIR")
)

func TestIntegration(t *testing.T) {
	ctx := context.Background()

	waitForPostgres(ctx, t, dsn)

	cmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, "../../cmd/migr8")
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to build binary: %v", err)
	}

	testCreateMigration(t)

	t.Run("concurrent up migrations", func(t *testing.T) {
		testConcurrentUp(t)
	})

	t.Run("concurrent up and down", func(t *testing.T) {
		testConcurrentUpAndDown(t)
	})

	t.Run("concurrent redo operations", func(t *testing.T) {
		testConcurrentRedo(t)
	})
}

func testCreateMigration(t *testing.T) {
	t.Helper()

	ctx := context.Background()

	cmd := exec.CommandContext(ctx, binaryPath, "create", "test_migration")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("create failed: %v\n%s", err, output)
	}

	// Check if file was created.
	files, _ := filepath.Glob(filepath.Join(dir, "*.sql"))
	if len(files) == 0 {
		t.Error("no migration file created")
	}
}

func testConcurrentUp(t *testing.T) {
	t.Helper()

	ctx := context.Background()

	var wg sync.WaitGroup
	errChan := make(chan error, 5)

	for i := range 5 {
		wg.Go(func() {
			cmd := exec.CommandContext(ctx, binaryPath, "up")

			output, err := cmd.CombinedOutput()
			if err != nil {
				errChan <- fmt.Errorf("instance %d failed: %w\n%s", i, err, output)
			} else {
				errChan <- nil
				t.Logf("Instance %d succeeded", i)
			}
		})
	}

	wg.Wait()
	close(errChan)

	// Count successes and failures
	var successCount int
	for err := range errChan {
		if err == nil {
			successCount++
		}
	}

	require.Equal(t, 1, successCount, "Expected 1 success")
	verifyMigrationCount(t, dsn, 1)
}

func testConcurrentUpAndDown(t *testing.T) {
	t.Helper()

	ctx := context.Background()

	var wg sync.WaitGroup
	errChan := make(chan error, 10)

	for i := range 5 {
		wg.Go(func() {
			cmd := exec.CommandContext(ctx, binaryPath, "up")

			output, err := cmd.CombinedOutput()
			if err != nil {
				errChan <- fmt.Errorf("up-%d failed: %w\n%s", i, err, output)
			} else {
				errChan <- nil
				t.Logf("Instance %d succeeded", i)
			}
		})

		wg.Go(func() {
			cmd := exec.CommandContext(ctx, binaryPath, "down")
			output, err := cmd.CombinedOutput()
			if err != nil {
				errChan <- fmt.Errorf("down-%d failed: %w\n%s", i, err, output)
			} else {
				errChan <- nil
				t.Logf("Instance %d succeeded", i)
			}
		})
	}

	wg.Wait()
	close(errChan)

	// At least one operation should succeed.
	var anySuccess bool
	for err := range errChan {
		if err == nil {
			anySuccess = true
		}
	}

	require.True(t, anySuccess, "No operations succeeded")
}

func testConcurrentRedo(t *testing.T) {
	t.Helper()

	ctx := context.Background()

	cmd := exec.CommandContext(ctx, binaryPath, "up")
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to apply initial migration: %v", err)
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 5)

	for i := range 5 {
		wg.Go(func() {
			cmd := exec.CommandContext(ctx, binaryPath, "redo")

			output, err := cmd.CombinedOutput()
			if err != nil {
				errChan <- fmt.Errorf("instance %d failed: %w\n%s", i, err, output)
			} else {
				errChan <- nil
				t.Logf("Instance %d succeeded", i)
			}
		})
	}

	wg.Wait()
	close(errChan)

	var successCount int
	for err := range errChan {
		if err == nil {
			successCount++
		}
	}

	require.Equal(t, 1, successCount, "Expected 1 successful redo")
}

func verifyMigrationCount(t *testing.T, dsn string, expectedCount int) {
	t.Helper()

	ctx := context.Background()

	db, err := sql.Open("pgx", dsn)
	require.NoError(t, err)
	defer db.Close()

	var count int
	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM changelog").Scan(&count)
	require.NoError(t, err)

	require.Equal(t, expectedCount, count)
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

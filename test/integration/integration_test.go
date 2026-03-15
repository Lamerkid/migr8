//go:build integration

package integration

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegration(t *testing.T) {
	ctx := context.Background()

	dsn := os.Getenv("M8_DSN")
	dir := os.Getenv("M8_DIR")

	t.Run("create migration", func(t *testing.T) {
		cmd := exec.CommandContext(ctx, "migr8", "create", "test_migration")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("create failed: %v\n%s", err, output)
		}

		// Check if file was created.
		files, _ := filepath.Glob(filepath.Join(dir, "*.sql"))
		if len(files) == 0 {
			t.Error("no migration file created")
		}
	})

	require.NotEmpty(t, dsn)
}

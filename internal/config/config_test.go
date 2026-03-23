package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	flags := make(map[string]string)

	err := os.Setenv("M8_DIR", "./migrations/test")
	if err != nil {
		t.Fatalf("failed to set environment variables: %v", err)
	}

	err = os.Setenv("M8_DSN", "postgres://env")
	if err != nil {
		t.Fatalf("failed to set environment variables: %v", err)
	}

	// Flag should reassign config for dsn.
	flags["-dsn"] = "postgres://test"
	config, err := BuildFromFlags(flags)
	if err != nil {
		t.Fatalf("failed to build config: %v", err)
	}

	require.Equal(t, "postgres://test", config.Database.DSN)
	require.Equal(t, "./migrations/test", config.Migration.Dir)

	// Load config from json file.
	flags["-config"] = "./config.json"
	config2, err := BuildFromFlags(flags)
	if err != nil {
		t.Fatalf("failed to build config: %v", err)
	}

	require.NotNil(t, config2)
}

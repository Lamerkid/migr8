package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	flags := make(map[string]string)

	err := os.Setenv("M8_DSN", "postgres://env")
	if err != nil {
		t.Fatalf("failed to set environment variables: %v", err)
	}

	err = os.Setenv("M8_DIR", "./migrations/env")
	if err != nil {
		t.Fatalf("failed to set environment variables: %v", err)
	}

	// Load config from json file with expanded env variables.
	flags["-cfg"] = "./config.json"
	config, err := BuildFromFlags(flags)
	if err != nil {
		t.Fatalf("failed to build config: %v", err)
	}

	require.Equal(t, "postgres://env", config.Database.DSN)
	require.Equal(t, "./migrations/env", config.Migration.Dir)

	// Flag should reassign config for dsn.
	flags["-dsn"] = "postgres://test"
	flags["-dir"] = "./migrations/test"
	config2, err := BuildFromFlags(flags)
	if err != nil {
		t.Fatalf("failed to build config: %v", err)
	}

	require.Equal(t, "postgres://test", config2.Database.DSN)
	require.Equal(t, "./migrations/test", config2.Migration.Dir)
}

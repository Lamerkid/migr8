// Package config provides functions for configuration.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config is the main configuration.
type Config struct {
	Logger    *loggerConf    `json:"logger,omitempty"`
	Database  *databaseConf  `json:"database,omitempty"`
	Migration *migrationConf `json:"migration,omitempty"`
}

type loggerConf struct {
	Level string `json:"level"`
}

type databaseConf struct {
	DSN string `json:"dsn"`
}

type migrationConf struct {
	Type string `json:"type"`
	Path string `json:"path"`
}

// BuildFromFlags creates a config from parsed CLI flags.
func BuildFromFlags(flags map[string]string) (*Config, error) {
	config := defaultConfig()

	// Override with config file first.
	if configFile, ok := flags["-config"]; ok && configFile != "" {
		fileConfig, err := loadFromFile(configFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load config file: %w", err)
		}
		config = mergeConfig(config, fileConfig)
	}

	// Override with CLI flags.
	if dbURL, ok := flags["-dsn"]; ok && dbURL != "" {
		config.Database.DSN = dbURL
	}

	if path, ok := flags["-path"]; ok && path != "" {
		config.Migration.Path = path
	}

	return config, nil
}

func defaultConfig() *Config {
	return &Config{
		Logger: &loggerConf{
			Level: "INFO",
		},
		Database: &databaseConf{
			DSN: os.Getenv("M8_DSN"),
		},
		Migration: &migrationConf{
			Type: "sql",
			Path: os.Getenv("M8_MIG_DIR"),
		},
	}
}

func loadFromFile(path string) (*Config, error) {
	var config *Config
	cleanedPath := filepath.Clean(path)

	confFile, err := os.Open(cleanedPath)
	if err != nil {
		return &Config{}, err
	}
	defer confFile.Close()

	decoder := json.NewDecoder(confFile)
	if err = decoder.Decode(&config); err != nil {
		return &Config{}, err
	}

	return config, nil
}

func mergeConfig(base, override *Config) *Config {
	if override.Database.DSN != "" {
		base.Database.DSN = override.Database.DSN
	}
	if override.Migration.Path != "" {
		base.Migration.Path = override.Migration.Path
	}
	if override.Logger.Level != "" {
		base.Logger.Level = override.Logger.Level
	}
	return base
}

package core

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// MigrationSource reads migrations from the file.
type MigrationSource struct {
	Logger Logger
	Path   string
}

// LoadMigrations reads all migration files from the directory.
func (s *MigrationSource) LoadMigrations() ([]*Migration, error) {
	files, err := os.ReadDir(s.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var migrations []*Migration

	for _, file := range files {
		if file.IsDir() {
			s.Logger.Error("%s is a directory", file.Name())
			continue
		}

		ext := strings.ToLower(filepath.Ext(file.Name()))

		version, err := s.parseVersion(file.Name())
		if err != nil {
			s.Logger.Error("error parsing version for %s", file.Name())
			continue
		}

		switch ext {
		case ".sql":
			s.Logger.Debug("parsing sql file %s", file.Name())

			content, err := os.ReadFile(filepath.Join(s.Path, file.Name()))
			if err != nil {
				return nil, fmt.Errorf("failed to read %s: %w", file.Name(), err)
			}

			upSQL, downSQL := s.splitMigrationSQL(string(content))

			migrations = append(migrations, &Migration{
				Version: version,
				Name:    file.Name(),
				Type:    TypeSQL,
				UpSQL:   upSQL,
				DownSQL: downSQL,
			})

		case ".go":
			s.Logger.Debug("parsing go file %s", file.Name())

			mig, exists := GetRegisteredMigration(version)
			if !exists {
				return nil, fmt.Errorf("go migration %s not registered. Make sure to register it in init()", file.Name())
			}

			migrations = append(migrations, &Migration{
				Version:  version,
				Name:     file.Name(),
				Type:     TypeGo,
				UpFunc:   mig.Up,
				DownFunc: mig.Down,
			})
		default:
			continue
		}
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

func (s *MigrationSource) splitMigrationSQL(content string) (up, down string) {
	parts := strings.Split(content, "-- +migr8:down")
	if len(parts) != 2 {
		s.Logger.Debug("file has no down sql section")
		return strings.TrimSpace(content), ""
	}

	upPart := strings.Split(parts[0], "-- +mig8:up")
	if len(upPart) == 2 {
		up = strings.TrimSpace(upPart[1])
	} else {
		up = strings.TrimSpace(parts[0])
	}

	down = strings.TrimSpace(parts[1])
	return
}

func (s *MigrationSource) parseVersion(filename string) (int64, error) {
	var version int64
	_, err := fmt.Sscanf(filename, "%d_", &version)
	return version, err
}

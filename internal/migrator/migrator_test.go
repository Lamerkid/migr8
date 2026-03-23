package migrator

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	postgres "github.com/Lamerkid/migr8/internal/database"
	"github.com/Lamerkid/migr8/internal/logger"
)

func TestMigrator(t *testing.T) {
	ctx := context.Background()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}

	database := &postgres.Database{DB: db}
	defer database.Close()

	logger := logger.NewLogger("INFO")

	m := NewMigrator(database, logger)

	// Check for successuful passed command to db.
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS changelog").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS migrationlocks").
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err = m.CreateServiceTables(ctx); err != nil {
		t.Fatalf("error creating service tables for migrations: %v", err)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations not met: %v", err)
	}
}

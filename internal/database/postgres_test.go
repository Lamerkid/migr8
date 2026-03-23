package postgres

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreateServiceTables(t *testing.T) {
	ctx := context.Background()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}

	database := &Database{DB: db}
	defer database.Close()

	// Set expectations.
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS changelog").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS migrationlocks").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = database.CreateServiceTables(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify expectations.
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations not met: %v", err)
	}
}

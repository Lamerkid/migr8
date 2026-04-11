package postgres

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestCreateServiceTables(t *testing.T) {
	ctx := context.Background()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}

	database := &Database{db: db}
	defer database.Close()

	// Set expectations.
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS changelog").
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectBegin()

	mock.ExpectExec("INSERT INTO changelog").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	rows := sqlmock.NewRows([]string{"version"})
	mock.ExpectQuery("SELECT version FROM changelog").
		WillReturnRows(rows)

	versionRow := sqlmock.NewRows([]string{"version"}).
		AddRow(111)
	mock.ExpectQuery("SELECT COALESCE\\(MAX\\(version\\), 0\\) FROM changelog").
		WillReturnRows(versionRow)

	mock.ExpectBegin()

	mock.ExpectExec("DELETE FROM changelog WHERE version = \\$1").
		WithArgs(111).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectCommit()

	// Begin commands execution.
	if err = database.CreateServiceTable(ctx); err != nil {
		t.Fatalf("error creating service table: %v", err)
	}

	tx, err := database.BeginTx(ctx, nil)
	require.NoError(t, err)

	if err = database.InsertApplied(ctx, tx, 111, "111_test.sql"); err != nil {
		t.Fatalf("error inserting data to service table: %v", err)
	}

	err = tx.Commit()
	require.NoError(t, err)

	if _, err = database.GetAppliedVersions(ctx); err != nil {
		t.Fatalf("could not get applied versions: %v", err)
	}

	if _, err = database.GetLatestVersion(ctx); err != nil {
		t.Fatalf("could not get latest version: %v", err)
	}

	tx, err = database.BeginTx(ctx, nil)
	require.NoError(t, err)

	if err = database.DeleteApplied(ctx, tx, 111); err != nil {
		t.Fatalf("error inserting data to service table: %v", err)
	}

	err = tx.Commit()
	require.NoError(t, err)

	// Verify expectations.
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations not met: %v", err)
	}
}

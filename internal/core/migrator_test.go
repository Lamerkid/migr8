package core

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestMigrator(t *testing.T) {
	ctx := context.Background()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	m := NewMigrator(db, ".")

	// Set expectations.
	lockExpectations(t, mock)
	upExpectations(t, mock)
	unlockExpectations(t, mock)

	lockExpectations(t, mock)
	redoExpectations(t, mock)
	unlockExpectations(t, mock)

	lockExpectations(t, mock)
	downExpectations(t, mock)
	unlockExpectations(t, mock)

	statusExpectations(t, mock)

	// Begin commands execution.
	if err = m.Up(ctx); err != nil {
		t.Fatalf("error applying migrations: %v", err)
	}

	if err = m.Redo(ctx); err != nil {
		t.Fatalf("error redo last migration: %v", err)
	}

	if err = m.Down(ctx); err != nil {
		t.Fatalf("error rollback last migration: %v", err)
	}

	if err = m.Status(ctx); err != nil {
		t.Fatalf("error checking migration status: %v", err)
	}

	// Verify expectations.
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations not met: %v", err)
	}
}

func TestMigratorConcurrently(t *testing.T) {
	ctx := context.Background()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	m := NewMigrator(db, ".")

	lockExpectations(t, mock)
	upExpectations(t, mock)
	unlockExpectations(t, mock)

	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	// Run two migrations concurrently
	for range 2 {
		wg.Go(func() {
			errChan <- m.Up(ctx)
		})
	}

	wg.Wait()
	close(errChan)

	var successCount, failCount int
	for err := range errChan {
		if err == nil {
			successCount++
		} else {
			t.Log(err)
			failCount++
		}
	}

	require.Equal(t, 1, successCount, "Expected 1 success")
	require.Equal(t, 1, failCount, "Expected 1 failure")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations not met: %v", err)
	}
}

func TestMigratorConcurrentDifferentCommands(t *testing.T) {
	ctx := context.Background()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	m := NewMigrator(db, ".")

	lockExpectations(t, mock)
	upExpectations(t, mock)
	unlockExpectations(t, mock)

	lockExpectationsFail(t, mock)

	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	wg.Go(func() {
		errChan <- m.Up(ctx)
	})

	time.Sleep(10 * time.Millisecond)

	wg.Go(func() {
		errChan <- m.Redo(ctx)
	})

	wg.Wait()
	close(errChan)

	var successCount, failCount int
	for err := range errChan {
		if err == nil {
			successCount++
		} else {
			failCount++
		}
	}

	require.Equal(t, 1, successCount, "Expected 1 success")
	require.Equal(t, 1, failCount, "Expected 1 failure")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations not met: %v", err)
	}
}

func lockExpectations(t *testing.T, mock sqlmock.Sqlmock) {
	t.Helper()

	mock.ExpectQuery("SELECT pg_try_advisory_lock\\(\\$1\\)").
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"pg_try_advisory_lock"}).AddRow(true))
}

func lockExpectationsFail(t *testing.T, mock sqlmock.Sqlmock) {
	t.Helper()

	mock.ExpectQuery("SELECT pg_try_advisory_lock\\(\\$1\\)").
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"pg_try_advisory_lock"}).AddRow(false))
}

func unlockExpectations(t *testing.T, mock sqlmock.Sqlmock) {
	t.Helper()

	mock.ExpectQuery("SELECT pg_advisory_unlock\\(\\$1\\)").
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"pg_advisory_unlock"}).AddRow(true))
}

func upExpectations(t *testing.T, mock sqlmock.Sqlmock) {
	t.Helper()

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS changelog").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery("SELECT version FROM changelog").
		WillReturnRows(sqlmock.NewRows([]string{"version"}))

	mock.ExpectBegin()
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS test").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("INSERT INTO changelog").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
}

func redoExpectations(t *testing.T, mock sqlmock.Sqlmock) {
	t.Helper()

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS changelog").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery("SELECT COALESCE\\(MAX\\(version\\), 0\\) FROM changelog").
		WillReturnRows(sqlmock.NewRows([]string{"coalesce"}).AddRow(123))

	mock.ExpectBegin()
	mock.ExpectExec("DROP TABLE IF EXISTS test").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("DELETE FROM changelog WHERE version = \\$1").
		WithArgs(123).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS test").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("INSERT INTO changelog").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
}

func downExpectations(t *testing.T, mock sqlmock.Sqlmock) {
	t.Helper()

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS changelog").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery("SELECT COALESCE\\(MAX\\(version\\), 0\\) FROM changelog").
		WillReturnRows(sqlmock.NewRows([]string{"coalesce"}).AddRow(123))

	mock.ExpectBegin()
	mock.ExpectExec("DROP TABLE IF EXISTS test").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("DELETE FROM changelog WHERE version = \\$1").
		WithArgs(123).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
}

func statusExpectations(t *testing.T, mock sqlmock.Sqlmock) {
	t.Helper()

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS changelog").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery("SELECT version FROM changelog").
		WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow(123))
}

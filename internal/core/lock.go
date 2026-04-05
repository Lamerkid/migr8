package core

import (
	"context"
	"fmt"
)

// AdvisoryLock using PostgreSQL advisory locks.
type AdvisoryLock struct {
	DB Database
}

// Acquire lock with hardcoded ID.
func (alm *AdvisoryLock) Acquire(ctx context.Context, lockID int64) (bool, error) {
	var acquired bool
	err := alm.DB.QueryRowContext(ctx, "SELECT pg_try_advisory_lock($1)", lockID).Scan(&acquired)
	return acquired, err
}

// Release lock with hardcoded ID.
func (alm *AdvisoryLock) Release(ctx context.Context, lockID int64) error {
	var released bool
	err := alm.DB.QueryRowContext(ctx, "SELECT pg_advisory_unlock($1)", lockID).Scan(&released)
	if err != nil {
		return err
	}
	if !released {
		return fmt.Errorf("lock was not held")
	}
	return nil
}

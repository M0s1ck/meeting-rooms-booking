package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// AdvisoryLocker runs a job only if no other process currently holds the same
// PostgreSQL session-level advisory lock. It lets several app replicas run the same
// background jobs safely: only one replica does the work at a time (12-factor, VIII. Concurrency).
// If the process dies, the connection closes and PostgreSQL releases the lock.
type AdvisoryLocker struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NewAdvisoryLocker(db *pgxpool.Pool, logger *slog.Logger) *AdvisoryLocker {
	return &AdvisoryLocker{db: db, logger: logger}
}

// WithLock wraps job: it is executed only when the lock with the given key is acquired,
// otherwise it is skipped (another replica is already doing it).
func (l *AdvisoryLocker) WithLock(key int64, job func(context.Context) error) func(context.Context) error {
	return func(ctx context.Context) error {
		conn, err := l.db.Acquire(ctx)
		if err != nil {
			return fmt.Errorf("advisory lock: acquire conn: %w", err)
		}
		defer conn.Release()

		var locked bool
		if err := conn.QueryRow(ctx, "SELECT pg_try_advisory_lock($1)", key).Scan(&locked); err != nil {
			return fmt.Errorf("advisory lock: try lock: %w", err)
		}

		if !locked {
			l.logger.Info("job skipped: advisory lock is held by another process", "lock_key", key)
			return nil
		}

		defer func() {
			// use a fresh context: ctx may already be cancelled on shutdown
			unlockCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_, _ = conn.Exec(unlockCtx, "SELECT pg_advisory_unlock($1)", key)
		}()

		return job(ctx)
	}
}

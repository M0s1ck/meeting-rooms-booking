package repository

import (
	"context"
	"errors"
	"fmt"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/schedule"
)

type ScheduleRepo struct {
	db     *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func NewScheduleRepo(db *pgxpool.Pool, txGetter *trmpgx.CtxGetter) *ScheduleRepo {
	return &ScheduleRepo{
		db:     db,
		getter: txGetter,
	}
}

func (r *ScheduleRepo) Create(ctx context.Context, sch *schedule.Schedule) error {
	querier := r.getter.DefaultTrOrDB(ctx, r.db)

	const query = `
        INSERT INTO room_schedules (id, room_id, days_of_week, start_time, end_time, created_at)
        VALUES ($1, $2, $3, $4::time, $5::time, $6)
    `

	_, err := querier.Exec(ctx, query,
		sch.ID,
		sch.RoomID,
		sch.DaysOfWeek,
		sch.StartTime.String(),
		sch.EndTime.String(),
		sch.CreatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return schedule.ErrAlreadyExists
			case "23503":
				return schedule.ErrRoomNotFound
			}
		}

		return fmt.Errorf("create schedule: %w", err)
	}

	return nil
}

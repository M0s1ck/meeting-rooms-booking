package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

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

func (r *ScheduleRepo) List(ctx context.Context) ([]schedule.Schedule, error) {
	querier := r.getter.DefaultTrOrDB(ctx, r.db)

	const query = `SELECT id, room_id, days_of_week, start_time, end_time, created_at
        FROM room_schedules`

	rows, err := querier.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list schedules: %w", err)
	}

	defer rows.Close()
	schedules := make([]schedule.Schedule, 0)

	for rows.Next() {
		var (
			sch          schedule.Schedule
			daysOfWeek   []int16
			startTimeRaw time.Time
			endTimeRaw   time.Time
		)

		if err := rows.Scan(
			&sch.ID, &sch.RoomID, &daysOfWeek,
			&startTimeRaw, &endTimeRaw, &sch.CreatedAt,
		); err != nil {

			return nil, fmt.Errorf("list schedules: %w", err)
		}

		startTime, err := schedule.ParseTimeOfDay(startTimeRaw.Format("15:04"))
		if err != nil {
			return nil, fmt.Errorf("list schedules parse start time: %w", err)
		}

		endTime, err := schedule.ParseTimeOfDay(endTimeRaw.Format("15:04"))
		if err != nil {
			return nil, fmt.Errorf("list schedules parse end time: %w", err)
		}

		sch.DaysOfWeek = smallInt16SliceToInt(daysOfWeek)
		sch.StartTime = startTime
		sch.EndTime = endTime

		schedules = append(schedules, sch)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list schedules: %w", err)
	}

	return schedules, nil
}

func smallInt16SliceToInt(values []int16) []int {
	if len(values) == 0 {
		return nil
	}

	res := make([]int, len(values))
	for i, value := range values {
		res[i] = int(value)
	}

	return res
}

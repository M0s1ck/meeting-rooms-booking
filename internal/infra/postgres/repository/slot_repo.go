package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/slot"
)

type SlotRepo struct {
	db     *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func NewSlotRepo(db *pgxpool.Pool, txGetter *trmpgx.CtxGetter) *SlotRepo {
	return &SlotRepo{
		db:     db,
		getter: txGetter,
	}
}

func (r *SlotRepo) Add(ctx context.Context, slots []slot.Slot) error {
	if len(slots) == 0 {
		return nil
	}

	querier := r.getter.DefaultTrOrDB(ctx, r.db)

	const (
		columnsPerRow = 5
		chunkSize     = 500
	)

	for start := 0; start < len(slots); start += chunkSize {
		end := min(start+chunkSize, len(slots))
		chunk := slots[start:end]

		var builder strings.Builder
		builder.WriteString(`
        INSERT INTO slots (id, room_id, start_at, end_at, created_at)
        VALUES `)

		args := make([]any, 0, len(chunk)*columnsPerRow)
		for i, item := range chunk {
			if i > 0 {
				builder.WriteString(", ")
			}

			argPos := i*columnsPerRow + 1
			_, _ = fmt.Fprintf(
				&builder,
				"($%d, $%d, $%d, $%d, $%d)",
				argPos,
				argPos+1,
				argPos+2,
				argPos+3,
				argPos+4,
			)

			args = append(args,
				item.ID,
				item.RoomID,
				item.StartAt,
				item.EndAt,
				item.CreatedAt,
			)
		}

		builder.WriteString(" ON CONFLICT (room_id, start_at) DO NOTHING")

		if _, err := querier.Exec(ctx, builder.String(), args...); err != nil {
			return fmt.Errorf("add slots batch [%d:%d): %w", start, end, err)
		}

	}

	return nil
}

func (r *SlotRepo) GetLastEndAtByRooms(
	ctx context.Context,
	roomIDS []uuid.UUID,
) (map[uuid.UUID]time.Time, error) {

	querier := r.getter.DefaultTrOrDB(ctx, r.db)

	query := `SELECT room_id, MAX(end_at)
		FROM slots
		WHERE room_id = ANY($1)
		GROUP BY room_id`

	rows, err := querier.Query(ctx, query, roomIDS)
	if err != nil {
		return nil, fmt.Errorf("get last end at by rooms: %w", err)
	}

	defer rows.Close()

	roomEnd := make(map[uuid.UUID]time.Time, len(roomIDS))

	for rows.Next() {
		var roomID uuid.UUID
		var end time.Time

		if err := rows.Scan(&roomID, &end); err != nil {
			return nil, fmt.Errorf("get last end at by rooms scan: %w", err)
		}

		roomEnd[roomID] = end
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("get last end at by rooms rows: %w", err)
	}

	return roomEnd, nil
}

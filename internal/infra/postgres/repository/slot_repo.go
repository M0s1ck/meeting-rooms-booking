package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

func (r *SlotRepo) GetByID(ctx context.Context, slotID uuid.UUID) (*slot.Slot, error) {
	querier := r.getter.DefaultTrOrDB(ctx, r.db)

	const query = `
		SELECT id, room_id, start_at, end_at, created_at
		FROM slots
		WHERE id = $1
	`

	row := querier.QueryRow(ctx, query, slotID)
	slt, err := scanSlot(row)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, slot.ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("get slot by id: %w", err)
	}

	return slt, nil
}

func (r *SlotRepo) ListAvailableByRoomAndDate(
	ctx context.Context,
	roomID uuid.UUID,
	date time.Time,
) ([]slot.Slot, error) {
	querier := r.getter.DefaultTrOrDB(ctx, r.db)

	const query = `
		SELECT s.id, s.room_id, s.start_at, s.end_at, s.created_at
		FROM slots s
		LEFT JOIN bookings b
			ON b.slot_id = s.id
			AND b.status = 'active'
		WHERE s.room_id = $1
			AND s.start_at >= $2
			AND s.start_at < $3
			AND b.id IS NULL
		ORDER BY s.start_at, s.id
	`

	lower := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	upper := lower.AddDate(0, 0, 1)

	rows, err := querier.Query(ctx, query, roomID, lower, upper)
	if err != nil {
		return nil, fmt.Errorf("list available slots by room and date: %w", err)
	}
	defer rows.Close()

	slots := make([]slot.Slot, 0)
	for rows.Next() {
		item, scanErr := scanSlot(rows)
		if scanErr != nil {
			return nil, scanErr
		}

		slots = append(slots, *item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("list available slots by room and date rows: %w", err)
	}

	return slots, nil
}

func (r *SlotRepo) GetLastEndAtByRooms(
	ctx context.Context,
	roomIDS []uuid.UUID,
) (map[uuid.UUID]time.Time, error) {
	if len(roomIDS) == 0 {
		return map[uuid.UUID]time.Time{}, nil
	}

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

type slotScanner interface {
	Scan(dest ...any) error
}

func scanSlot(scanner slotScanner) (*slot.Slot, error) {
	var item slot.Slot

	if err := scanner.Scan(
		&item.ID,
		&item.RoomID,
		&item.StartAt,
		&item.EndAt,
		&item.CreatedAt,
	); err != nil {
		return nil, fmt.Errorf("scan slot: %w", err)
	}

	return &item, nil
}

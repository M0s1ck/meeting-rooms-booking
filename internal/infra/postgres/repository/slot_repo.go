package repository

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/slot"
	"github.com/jackc/pgx/v5/pgxpool"
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
	slots = deduplicateSlots(slots)
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

func deduplicateSlots(slots []slot.Slot) []slot.Slot {
	if len(slots) < 2 {
		return slots
	}

	unique := make(map[string]slot.Slot, len(slots))
	for _, item := range slots {
		startAt := item.StartAt.UTC()
		key := item.RoomID.String() + "|" + startAt.Format(time.RFC3339)

		if _, exists := unique[key]; exists {
			continue
		}

		item.StartAt = startAt
		item.EndAt = item.EndAt.UTC()
		item.CreatedAt = item.CreatedAt.UTC()
		unique[key] = item
	}

	result := make([]slot.Slot, 0, len(unique))
	for _, item := range unique {
		result = append(result, item)
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].RoomID == result[j].RoomID {
			return result[i].StartAt.Before(result[j].StartAt)
		}

		return result[i].RoomID.String() < result[j].RoomID.String()
	})

	return result
}

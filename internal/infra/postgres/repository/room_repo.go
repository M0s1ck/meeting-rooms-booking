package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/room"
)

type RoomRepo struct {
	db     *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func NewRoomRepo(db *pgxpool.Pool, txGetter *trmpgx.CtxGetter) *RoomRepo {
	return &RoomRepo{
		db:     db,
		getter: txGetter,
	}
}

func (r *RoomRepo) Create(ctx context.Context, rm *room.Room) error {
	querier := r.getter.DefaultTrOrDB(ctx, r.db)

	const query = `INSERT INTO rooms (id, name, description, capacity, created_at)
		VALUES ($1, $2, $3, $4, $5)`

	_, err := querier.Exec(ctx, query,
		rm.ID, rm.Name, rm.Description, rm.Capacity, rm.CreatedAt)

	if err != nil {
		return fmt.Errorf("create room: %w", err)
	}

	return nil
}

func (r *RoomRepo) List(ctx context.Context) ([]*room.Room, error) {
	querier := r.getter.DefaultTrOrDB(ctx, r.db)

	const query = `
		SELECT id, name, description, capacity, created_at, updated_at
		FROM rooms
		ORDER BY created_at ASC, id ASC
	`

	rows, err := querier.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list rooms query: %w", err)
	}
	defer rows.Close()

	rooms := make([]*room.Room, 0)
	for rows.Next() {
		item, scanErr := scanRoom(rows)
		if scanErr != nil {
			return nil, scanErr
		}

		rooms = append(rooms, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("list rooms rows: %w", err)
	}

	return rooms, nil
}

type roomScanner interface {
	Scan(dest ...any) error
}

func scanRoom(scanner roomScanner) (*room.Room, error) {
	var (
		id          uuid.UUID
		name        string
		description sql.NullString
		capacity    sql.NullInt32
		createdAt   time.Time
		updatedAt   time.Time
	)

	if err := scanner.Scan(
		&id,
		&name,
		&description,
		&capacity,
		&createdAt,
		&updatedAt,
	); err != nil {
		return nil, fmt.Errorf("scan room: %w", err)
	}

	return &room.Room{
		ID:          id,
		Name:        name,
		Description: nullStringToPtr(description),
		Capacity:    nullInt32ToPtr(capacity),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}, nil
}

func nullStringToPtr(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}

	v := value.String
	return &v
}

func nullInt32ToPtr(value sql.NullInt32) *int {
	if !value.Valid {
		return nil
	}

	v := int(value.Int32)
	return &v
}

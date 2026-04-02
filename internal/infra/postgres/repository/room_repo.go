package repository

import (
	"context"
	"fmt"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
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
		return fmt.Errorf("create link: %w", err)
	}

	return nil
}

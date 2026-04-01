package repository

import (
	"context"
	"fmt"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/room"
	"github.com/internships-backend/test-backend-M0s1ck/internal/infra/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RoomRepo struct {
	db *pgxpool.Pool
}

func NewRoomRepo(db *pgxpool.Pool) *RoomRepo {
	return &RoomRepo{
		db: db,
	}
}

func (r *RoomRepo) Create(ctx context.Context, rm *room.Room) error {
	querier := postgres.GetQuerier(ctx, r.db)

	const query = `INSERT INTO rooms (id, name, description, capacity, created_at)
		VALUES ($1, $2, $3, $4, $5)`

	_, err := querier.Exec(ctx, query,
		rm.ID, rm.Name, rm.Description, rm.Capacity, rm.CreatedAt)

	if err != nil {
		return fmt.Errorf("create link: %w", err)
	}

	return nil
}

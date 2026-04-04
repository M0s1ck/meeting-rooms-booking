package repository

import (
	"context"
	"fmt"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/user"
)

type UserRepo struct {
	db     *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func NewUserRepo(db *pgxpool.Pool, txGetter *trmpgx.CtxGetter) *UserRepo {
	return &UserRepo{
		db:     db,
		getter: txGetter,
	}
}

func (r *UserRepo) Ensure(ctx context.Context, userID uuid.UUID, role user.Role) error {
	querier := r.getter.DefaultTrOrDB(ctx, r.db)

	const query = `
		INSERT INTO users (id, email, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $4)
		ON CONFLICT (id) DO NOTHING
	`

	now := time.Now().UTC()
	email := fmt.Sprintf("dummy-%s-%s@example.test", role, userID.String())

	if _, err := querier.Exec(ctx, query, userID, email, role, now); err != nil {
		return fmt.Errorf("ensure user: %w", err)
	}

	return nil
}

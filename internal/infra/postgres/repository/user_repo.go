package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

func (r *UserRepo) Create(ctx context.Context, usr *user.User) error {
	querier := r.getter.DefaultTrOrDB(ctx, r.db)

	const query = `
        INSERT INTO users (id, email, password_hash, role, updated_at, created_at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `

	_, err := querier.Exec(ctx, query,
		usr.ID,
		usr.Email,
		usr.PassHash,
		usr.Role,
		usr.CreatedAt,
		usr.CreatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return user.ErrEmailAlreadyTaken
		}

		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email user.Email) (*user.User, error) {
	querier := r.getter.DefaultTrOrDB(ctx, r.db)

	const query = `
		SELECT id, email, password_hash, role, created_at
		FROM users
		WHERE email = $1
	`

	row := querier.QueryRow(ctx, query, email)
	usr, err := scanUser(row)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, user.ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	return usr, nil
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

type userScanner interface {
	Scan(dest ...interface{}) error
}

func scanUser(scanner userScanner) (*user.User, error) {
	var usr user.User

	if err := scanner.Scan(
		&usr.ID,
		&usr.Email,
		&usr.PassHash,
		&usr.Role,
		&usr.CreatedAt); err != nil {
		return nil, fmt.Errorf("scan user: %w", err)
	}

	return &usr, nil
}

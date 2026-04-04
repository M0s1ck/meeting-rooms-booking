package repository

import (
	"context"
	"errors"
	"fmt"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/booking"
)

type BookingRepo struct {
	db     *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func NewBookingRepo(db *pgxpool.Pool, txGetter *trmpgx.CtxGetter) *BookingRepo {
	return &BookingRepo{
		db:     db,
		getter: txGetter,
	}
}

func (r *BookingRepo) Create(ctx context.Context, b *booking.Booking) error {
	querier := r.getter.DefaultTrOrDB(ctx, r.db)

	const query = `
		INSERT INTO bookings (id, slot_id, user_id, status, conference_link, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := querier.Exec(ctx, query,
		b.ID,
		b.SlotID,
		b.UserID,
		b.Status,
		b.ConferenceLink,
		b.CreatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return booking.ErrAlreadyBooked
		}

		return fmt.Errorf("create booking: %w", err)
	}

	return nil
}

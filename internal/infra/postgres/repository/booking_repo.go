package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
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

func (r *BookingRepo) List(ctx context.Context, page, pageSize int) ([]booking.Booking, int, error) {
	querier := r.getter.DefaultTrOrDB(ctx, r.db)

	const countQuery = `SELECT COUNT(*) FROM bookings`

	var total int
	if err := querier.QueryRow(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count bookings: %w", err)
	}

	offset := (page - 1) * pageSize

	const listQuery = `
		SELECT id, slot_id, user_id, status, conference_link, created_at
		FROM bookings
		ORDER BY created_at DESC, id DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := querier.Query(ctx, listQuery, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list bookings: %w", err)
	}
	defer rows.Close()

	bookings, err := collectBookings(rows)
	if err != nil {
		return nil, 0, fmt.Errorf("list bookings: %w", err)
	}

	return bookings, total, nil
}

func (r *BookingRepo) ListFutureByUser(ctx context.Context, userID uuid.UUID, now time.Time) ([]booking.Booking, error) {
	querier := r.getter.DefaultTrOrDB(ctx, r.db)

	const query = `
		SELECT b.id, b.slot_id, b.user_id, b.status, b.conference_link, b.created_at
		FROM bookings b
		JOIN slots s ON s.id = b.slot_id
		WHERE b.user_id = $1
		  AND s.start_at >= $2
		  AND b.status = 'active'	
		ORDER BY s.start_at, b.id
	`

	rows, err := querier.Query(ctx, query, userID, now.UTC())
	if err != nil {
		return nil, fmt.Errorf("list future bookings by user: %w", err)
	}
	defer rows.Close()

	items, err := collectBookings(rows)
	if err != nil {
		return nil, fmt.Errorf("list future bookings by user: %w", err)
	}

	return items, nil
}

type bookingScanner interface {
	Scan(dest ...any) error
}

func scanBooking(scanner bookingScanner) (*booking.Booking, error) {
	var b booking.Booking
	var statusRaw string

	if err := scanner.Scan(
		&b.ID,
		&b.SlotID,
		&b.UserID,
		&statusRaw,
		&b.ConferenceLink,
		&b.CreatedAt,
	); err != nil {
		return nil, err
	}

	b.Status = booking.Status(statusRaw)
	return &b, nil
}

type pgRows interface {
	Next() bool
	Err() error
	Scan(dest ...any) error
}

func collectBookings(rows pgRows) ([]booking.Booking, error) {
	bookings := make([]booking.Booking, 0)

	for rows.Next() {
		item, scanErr := scanBooking(rows)
		if scanErr != nil {
			return nil, scanErr
		}

		bookings = append(bookings, *item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return bookings, nil
}

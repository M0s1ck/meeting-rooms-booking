package app

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/booking"
	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/slot"
	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/user"
	"github.com/internships-backend/test-backend-M0s1ck/internal/infra/metrics"
	"github.com/internships-backend/test-backend-M0s1ck/internal/service/authjwt"
	"github.com/internships-backend/test-backend-M0s1ck/internal/usecase/auth/dummylogin"
	"github.com/internships-backend/test-backend-M0s1ck/internal/usecase/auth/login"
	"github.com/internships-backend/test-backend-M0s1ck/internal/usecase/auth/register"
	createbooking "github.com/internships-backend/test-backend-M0s1ck/internal/usecase/booking/create"
	createroom "github.com/internships-backend/test-backend-M0s1ck/internal/usecase/room/create"
)

// Decorators below add business metrics around usecases without touching
// the usecases themselves (and their unit tests): each one wraps the real
// usecase, calls it and records the outcome.

// ---------- create booking ----------

type instrumentedCreateBooking struct {
	next *createbooking.Usecase
	m    *metrics.Metrics
}

func (d *instrumentedCreateBooking) Execute(
	ctx context.Context, req *createbooking.Request, identity *authjwt.Identity,
) (*createbooking.Response, error) {
	resp, err := d.next.Execute(ctx, req, identity)
	if err != nil {
		d.m.BookingFailures.WithLabelValues(bookingFailureReason(err)).Inc()
		return nil, err
	}

	d.m.BookingsCreated.WithLabelValues(strconv.FormatBool(resp.ConferenceLink != nil)).Inc()
	return resp, nil
}

func bookingFailureReason(err error) string {
	switch {
	case errors.Is(err, booking.ErrAlreadyBooked):
		return metrics.BookingFailSlotAlreadyBooked
	case errors.Is(err, booking.ErrSlotInPast):
		return metrics.BookingFailSlotInPast
	case errors.Is(err, slot.ErrNotFound):
		return metrics.BookingFailSlotNotFound
	case errors.Is(err, user.ErrUserRoleRequired):
		return metrics.BookingFailForbidden
	case errors.Is(err, booking.ErrSlotIDRequired), errors.Is(err, booking.ErrUserIDRequired):
		return metrics.BookingFailInvalidRequest
	default:
		return metrics.BookingFailInternal
	}
}

// ---------- cancel booking ----------

// cancelBookingRepo is what the cancel usecase needs from the repository.
type cancelBookingRepo interface {
	GetByID(ctx context.Context, id uuid.UUID) (*booking.Booking, error)
	Cancel(ctx context.Context, id uuid.UUID, now time.Time) error
}

// instrumentedCancelRepo counts real cancellations. Cancel endpoint is
// idempotent (repeated cancel returns 200 too), but the usecase calls
// repo.Cancel only for an active booking, so this is the exact place
// where an active booking turns into a cancelled one.
type instrumentedCancelRepo struct {
	cancelBookingRepo
	m *metrics.Metrics
}

func (d *instrumentedCancelRepo) Cancel(ctx context.Context, id uuid.UUID, now time.Time) error {
	if err := d.cancelBookingRepo.Cancel(ctx, id, now); err != nil {
		return err
	}

	d.m.BookingsCancelled.Inc()
	return nil
}

// ---------- login (email + password) ----------

type instrumentedLogin struct {
	next *login.Usecase
	m    *metrics.Metrics
}

func (d *instrumentedLogin) Execute(ctx context.Context, req *login.Request) (*login.Response, error) {
	resp, err := d.next.Execute(ctx, req)
	d.m.Logins.WithLabelValues(metrics.LoginMethodPassword, resultLabel(err)).Inc()
	return resp, err
}

// ---------- dummy login ----------

type instrumentedDummyLogin struct {
	next *dummylogin.Usecase
	m    *metrics.Metrics
}

func (d *instrumentedDummyLogin) Execute(ctx context.Context, req *dummylogin.Request) (*dummylogin.Response, error) {
	resp, err := d.next.Execute(ctx, req)
	d.m.Logins.WithLabelValues(metrics.LoginMethodDummy, resultLabel(err)).Inc()
	return resp, err
}

// ---------- register ----------

type instrumentedRegister struct {
	next *register.Usecase
	m    *metrics.Metrics
}

func (d *instrumentedRegister) Execute(ctx context.Context, req *register.Request) (*register.Response, error) {
	resp, err := d.next.Execute(ctx, req)
	if err != nil {
		return nil, err
	}

	d.m.UsersRegistered.WithLabelValues(string(resp.Role)).Inc()
	return resp, nil
}

// ---------- create room ----------

type instrumentedCreateRoom struct {
	next *createroom.Usecase
	m    *metrics.Metrics
}

func (d *instrumentedCreateRoom) Execute(
	ctx context.Context, req *createroom.Request, identity *authjwt.Identity,
) (*createroom.Response, error) {
	resp, err := d.next.Execute(ctx, req, identity)
	if err != nil {
		return nil, err
	}

	d.m.RoomsCreated.Inc()
	return resp, nil
}

// ---------- conference service client ----------

type confLinkProvider interface {
	CreateLink(ctx context.Context, slotID uuid.UUID, userID uuid.UUID) (string, error)
}

type instrumentedConfLinkProvider struct {
	next confLinkProvider
	m    *metrics.Metrics
}

func (d *instrumentedConfLinkProvider) CreateLink(ctx context.Context, slotID uuid.UUID, userID uuid.UUID) (string, error) {
	link, err := d.next.CreateLink(ctx, slotID, userID)

	result := metrics.ConfLinkSuccess
	switch {
	case err == nil:
	case errors.Is(err, createbooking.ErrConferenceServiceUnavailable):
		result = metrics.ConfLinkUnavailable
	default:
		result = metrics.ConfLinkInternal
	}
	d.m.ConfLinkRequests.WithLabelValues(result).Inc()

	return link, err
}

func resultLabel(err error) string {
	if err != nil {
		return metrics.ResultFailure
	}
	return metrics.ResultSuccess
}

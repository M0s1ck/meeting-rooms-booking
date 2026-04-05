package create

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/booking"
	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/user"
	"github.com/internships-backend/test-backend-M0s1ck/internal/service/authjwt"
)

type Usecase struct {
	bookingRepo      bookingRepo
	slotRepo         slotRepo
	confLinkProvider conferenceLinkProvider
	logger           *slog.Logger
}

func NewUsecase(
	bookingRepo bookingRepo,
	slotRepo slotRepo,
	confLinkProvider conferenceLinkProvider,
	logger *slog.Logger,
) *Usecase {

	return &Usecase{
		bookingRepo:      bookingRepo,
		slotRepo:         slotRepo,
		confLinkProvider: confLinkProvider,
		logger:           logger,
	}
}

func (u *Usecase) Execute(ctx context.Context, req *Request, identity *authjwt.Identity) (*Response, error) {
	if err := u.authorize(identity); err != nil {
		return nil, err
	}

	now := time.Now().UTC()

	sl, err := u.slotRepo.GetByID(ctx, req.SlotID)
	if err != nil {
		return nil, err
	}

	if sl.StartAt.Before(now) {
		return nil, booking.ErrSlotInPast
	}

	book, err := booking.NewActive(req.SlotID, identity.UserID, now)
	if err != nil {
		return nil, err
	}

	if req.CreateConferenceLink {
		u.provideConfLink(ctx, book, req.SlotID, identity.UserID)
	}

	if err := u.bookingRepo.Create(ctx, book); err != nil {
		return nil, err
	}

	return &Response{
		ID:             book.ID,
		SlotID:         book.SlotID,
		UserID:         book.UserID,
		Status:         book.Status,
		ConferenceLink: book.ConferenceLink,
		CreatedAt:      book.CreatedAt,
	}, nil
}

func (u *Usecase) authorize(identity *authjwt.Identity) error {
	if identity == nil || identity.Role != user.RoleUser {
		return user.ErrUserRoleRequired
	}

	return nil
}

// in case of errors just log them, it makes booking still available
func (u *Usecase) provideConfLink(ctx context.Context, book *booking.Booking, slotID uuid.UUID, userID uuid.UUID) {
	link, err := u.confLinkProvider.CreateLink(ctx, slotID, userID)
	if err != nil {
		u.logger.Warn("failed to create conference link",
			"booking_id", book.ID,
			"slot_id", slotID,
			"err", err)
		return
	}

	err = book.SetConferenceLink(link)
	if err != nil {
		u.logger.Warn("conference is invalid",
			"booking_id", book.ID,
			"slot_id", slotID,
			"err", err)
		return
	}
}

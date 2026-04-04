package create

import (
	"context"
	"time"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/booking"
	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/user"
	"github.com/internships-backend/test-backend-M0s1ck/internal/service/authjwt"
)

type Usecase struct {
	bookingRepo bookingRepo
	slotRepo    slotRepo
}

func NewUsecase(
	bookingRepo bookingRepo,
	slotRepo slotRepo,
) *Usecase {

	return &Usecase{
		bookingRepo: bookingRepo,
		slotRepo:    slotRepo,
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

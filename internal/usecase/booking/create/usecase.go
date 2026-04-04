package create

import (
	"context"
	"time"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/booking"
	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/user"
	"github.com/internships-backend/test-backend-M0s1ck/internal/service/authjwt"
	"github.com/internships-backend/test-backend-M0s1ck/internal/usecase/common"
)

type Usecase struct {
	bookingRepo bookingRepo
	slotRepo    slotRepo
	userRepo    userRepo
	txManager   common.TxManager
}

func NewUsecase(
	bookingRepo bookingRepo,
	slotRepo slotRepo,
	userRepo userRepo,
	txManager common.TxManager,
) *Usecase {
	return &Usecase{
		bookingRepo: bookingRepo,
		slotRepo:    slotRepo,
		userRepo:    userRepo,
		txManager:   txManager,
	}
}

func (u *Usecase) Execute(ctx context.Context, req *Request, identity *authjwt.Identity) (*Response, error) {
	if err := u.authorize(identity); err != nil {
		return nil, err
	}

	now := time.Now().UTC()

	var created *booking.Booking
	err := u.txManager.Do(ctx, func(ctx context.Context) error {
		//
		if err := u.userRepo.Ensure(ctx, identity.UserID, identity.Role); err != nil {
			return err
		}

		sl, err := u.slotRepo.GetByID(ctx, req.SlotID)
		if err != nil {
			return err
		}

		if sl.StartAt.Before(now) {
			return booking.ErrSlotInPast
		}

		created, err = booking.NewActive(req.SlotID, identity.UserID, now)
		if err != nil {
			return err
		}

		return u.bookingRepo.Create(ctx, created)
	})
	if err != nil {
		return nil, err
	}

	return &Response{
		ID:             created.ID,
		SlotID:         created.SlotID,
		UserID:         created.UserID,
		Status:         created.Status,
		ConferenceLink: created.ConferenceLink,
		CreatedAt:      created.CreatedAt,
	}, nil
}

func (u *Usecase) authorize(identity *authjwt.Identity) error {
	if identity == nil || identity.Role != user.RoleUser {
		return user.ErrUserRoleRequired
	}

	return nil
}

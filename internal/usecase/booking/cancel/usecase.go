package cancel

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/booking"
	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/user"
	"github.com/internships-backend/test-backend-M0s1ck/internal/service/authjwt"
)

type Usecase struct {
	repo bookingRepo
}

func NewUsecase(repo bookingRepo) *Usecase {
	return &Usecase{repo: repo}
}

func (u *Usecase) Execute(ctx context.Context, bookingID uuid.UUID, identity *authjwt.Identity) (*Response, error) {
	if err := authorize(identity); err != nil {
		return nil, err
	}

	now := time.Now().UTC()

	book, err := u.repo.GetByID(ctx, bookingID)
	if err != nil {
		return nil, err
	}

	if book.UserID != identity.UserID {
		return nil, booking.ErrCancelForbidden
	}

	if book.Status == booking.StatusActive {
		if err := u.repo.Cancel(ctx, bookingID, now); err != nil {
			return nil, err
		}

		book, err = u.repo.GetByID(ctx, bookingID)
		if err != nil {
			return nil, err
		}
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

func authorize(identity *authjwt.Identity) error {
	if identity == nil || identity.Role != user.RoleUser {
		return user.ErrUserRoleRequired
	}

	return nil
}

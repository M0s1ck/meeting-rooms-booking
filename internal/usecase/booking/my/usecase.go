package my

import (
	"context"
	"time"

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

func (u *Usecase) Execute(ctx context.Context, identity *authjwt.Identity) (*Response, error) {
	if err := authorize(identity); err != nil {
		return nil, err
	}

	bookings, err := u.repo.ListFutureByUser(ctx, identity.UserID, time.Now().UTC())
	if err != nil {
		return nil, err
	}

	resp := buildResp(bookings)
	return resp, nil
}

func authorize(identity *authjwt.Identity) error {
	if identity == nil || identity.Role != user.RoleUser {
		return user.ErrUserRoleRequired
	}

	return nil
}

func buildResp(bookings []booking.Booking) *Response {
	resp := &Response{
		Bookings: make([]Booking, 0, len(bookings)),
	}

	for _, item := range bookings {
		resp.Bookings = append(resp.Bookings, Booking{
			ID:             item.ID,
			SlotID:         item.SlotID,
			UserID:         item.UserID,
			Status:         item.Status,
			ConferenceLink: item.ConferenceLink,
			CreatedAt:      item.CreatedAt,
		})
	}

	return resp
}

package list

import (
	"context"

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

func (u *Usecase) Execute(ctx context.Context, req *Request, identity *authjwt.Identity) (*Response, error) {
	if err := authorize(identity); err != nil {
		return nil, err
	}

	page, pageSize, err := normalizePagination(req)
	if err != nil {
		return nil, err
	}

	bookings, total, err := u.repo.List(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	resp := buildResp(bookings, page, pageSize, total)
	return resp, nil
}

func authorize(identity *authjwt.Identity) error {
	if identity == nil || identity.Role != user.RoleAdmin {
		return user.ErrAdminRoleRequired
	}

	return nil
}

func buildResp(bookings []booking.Booking, page, pageSize, total int) *Response {
	resp := &Response{
		Bookings: make([]Booking, 0, len(bookings)),
		Pagination: Pagination{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
		},
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

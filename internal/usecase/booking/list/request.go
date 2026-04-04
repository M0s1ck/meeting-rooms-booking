package list

import "github.com/internships-backend/test-backend-M0s1ck/internal/domain/booking"

const (
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 100
)

type Request struct {
	Page     int
	PageSize int
}

func normalizePagination(req *Request) (int, int, error) {
	page := defaultPage
	pageSize := defaultPageSize

	if req != nil {
		if req.Page != 0 {
			page = req.Page
		}
		if req.PageSize != 0 {
			pageSize = req.PageSize
		}
	}

	if page < 1 {
		return 0, 0, booking.ErrInvalidPage
	}

	if pageSize < 1 || pageSize > maxPageSize {
		return 0, 0, booking.ErrInvalidPageSize
	}

	return page, pageSize, nil
}

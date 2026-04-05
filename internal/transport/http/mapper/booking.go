package mapper

import (
	"errors"

	"github.com/internships-backend/test-backend-M0s1ck/internal/transport/http/oapi"
	cancelbooking "github.com/internships-backend/test-backend-M0s1ck/internal/usecase/booking/cancel"
	createbooking "github.com/internships-backend/test-backend-M0s1ck/internal/usecase/booking/create"
	listbooking "github.com/internships-backend/test-backend-M0s1ck/internal/usecase/booking/list"
	mybooking "github.com/internships-backend/test-backend-M0s1ck/internal/usecase/booking/my"
)

var ErrNilCreateBookingBody = errors.New("request body is required")

func ToCreateBookingRequest(body *oapi.PostBookingsCreateJSONRequestBody) (*createbooking.Request, error) {
	if body == nil {
		return nil, ErrNilCreateBookingBody
	}

	return &createbooking.Request{
		SlotID:               body.SlotId,
		CreateConferenceLink: body.CreateConferenceLink != nil && *body.CreateConferenceLink,
	}, nil
}

func ToCreateBookingResponse(resp *createbooking.Response) oapi.PostBookingsCreate201JSONResponse {
	if resp == nil {
		return oapi.PostBookingsCreate201JSONResponse{}
	}

	createdAt := resp.CreatedAt.UTC()

	return oapi.PostBookingsCreate201JSONResponse{
		Booking: &oapi.Booking{
			Id:             resp.ID,
			SlotId:         resp.SlotID,
			UserId:         resp.UserID,
			Status:         oapi.BookingStatus(resp.Status),
			ConferenceLink: resp.ConferenceLink,
			CreatedAt:      &createdAt,
		},
	}
}

func ToCancelBookingResponse(resp *cancelbooking.Response) oapi.PostBookingsBookingIdCancel200JSONResponse {
	if resp == nil {
		return oapi.PostBookingsBookingIdCancel200JSONResponse{}
	}

	createdAt := resp.CreatedAt.UTC()

	return oapi.PostBookingsBookingIdCancel200JSONResponse{
		Booking: &oapi.Booking{
			Id:             resp.ID,
			SlotId:         resp.SlotID,
			UserId:         resp.UserID,
			Status:         oapi.BookingStatus(resp.Status),
			ConferenceLink: resp.ConferenceLink,
			CreatedAt:      &createdAt,
		},
	}
}

func ToListBookingsRequest(params oapi.GetBookingsListParams) *listbooking.Request {
	req := &listbooking.Request{}
	if params.Page != nil {
		req.Page = *params.Page
	}
	if params.PageSize != nil {
		req.PageSize = *params.PageSize
	}

	return req
}

func ToListBookingsResponse(resp *listbooking.Response) oapi.GetBookingsList200JSONResponse {
	if resp == nil {
		return oapi.GetBookingsList200JSONResponse{}
	}

	bookings := make([]oapi.Booking, 0, len(resp.Bookings))
	for _, item := range resp.Bookings {
		createdAt := item.CreatedAt.UTC()
		bookings = append(bookings, oapi.Booking{
			Id:             item.ID,
			SlotId:         item.SlotID,
			UserId:         item.UserID,
			Status:         oapi.BookingStatus(item.Status),
			ConferenceLink: item.ConferenceLink,
			CreatedAt:      &createdAt,
		})
	}

	return oapi.GetBookingsList200JSONResponse{
		Bookings: &bookings,
		Pagination: &oapi.Pagination{
			Page:     resp.Pagination.Page,
			PageSize: resp.Pagination.PageSize,
			Total:    resp.Pagination.Total,
		},
	}
}

func ToMyBookingsResponse(resp *mybooking.Response) oapi.GetBookingsMy200JSONResponse {
	if resp == nil {
		return oapi.GetBookingsMy200JSONResponse{}
	}

	bookings := make([]oapi.Booking, 0, len(resp.Bookings))
	for _, item := range resp.Bookings {
		createdAt := item.CreatedAt.UTC()
		bookings = append(bookings, oapi.Booking{
			Id:             item.ID,
			SlotId:         item.SlotID,
			UserId:         item.UserID,
			Status:         oapi.BookingStatus(item.Status),
			ConferenceLink: item.ConferenceLink,
			CreatedAt:      &createdAt,
		})
	}

	return oapi.GetBookingsMy200JSONResponse{
		Bookings: &bookings,
	}
}

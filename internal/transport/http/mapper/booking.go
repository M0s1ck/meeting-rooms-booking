package mapper

import (
	"errors"

	"github.com/internships-backend/test-backend-M0s1ck/internal/transport/http/oapi"
	createbooking "github.com/internships-backend/test-backend-M0s1ck/internal/usecase/booking/create"
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

	createdAt := resp.CreatedAt

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

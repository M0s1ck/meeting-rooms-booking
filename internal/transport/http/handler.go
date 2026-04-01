package http

import (
	"context"

	"github.com/internships-backend/test-backend-M0s1ck/internal/transport/http/oapi"
)

type StrictHandler struct {
}

func NewHandler() *StrictHandler {
	return &StrictHandler{}
}

func (s StrictHandler) GetRoomsList(ctx context.Context, request oapi.GetRoomsListRequestObject) (oapi.GetRoomsListResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (s StrictHandler) PostRoomsCreate(ctx context.Context, request oapi.PostRoomsCreateRequestObject) (oapi.PostRoomsCreateResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (s StrictHandler) PostRoomsRoomIdScheduleCreate(ctx context.Context, request oapi.PostRoomsRoomIdScheduleCreateRequestObject) (oapi.PostRoomsRoomIdScheduleCreateResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (s StrictHandler) GetRoomsRoomIdSlotsList(ctx context.Context, request oapi.GetRoomsRoomIdSlotsListRequestObject) (oapi.GetRoomsRoomIdSlotsListResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (s StrictHandler) PostBookingsCreate(ctx context.Context, request oapi.PostBookingsCreateRequestObject) (oapi.PostBookingsCreateResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (s StrictHandler) GetBookingsList(ctx context.Context, request oapi.GetBookingsListRequestObject) (oapi.GetBookingsListResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (s StrictHandler) GetBookingsMy(ctx context.Context, request oapi.GetBookingsMyRequestObject) (oapi.GetBookingsMyResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (s StrictHandler) PostBookingsBookingIdCancel(ctx context.Context, request oapi.PostBookingsBookingIdCancelRequestObject) (oapi.PostBookingsBookingIdCancelResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (s StrictHandler) PostDummyLogin(ctx context.Context, request oapi.PostDummyLoginRequestObject) (oapi.PostDummyLoginResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (s StrictHandler) PostLogin(ctx context.Context, request oapi.PostLoginRequestObject) (oapi.PostLoginResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (s StrictHandler) PostRegister(ctx context.Context, request oapi.PostRegisterRequestObject) (oapi.PostRegisterResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

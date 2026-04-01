package http

import (
	"context"
	"errors"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/room"
	httphelpers "github.com/internships-backend/test-backend-M0s1ck/internal/transport/http/helpers"
	httpmapper "github.com/internships-backend/test-backend-M0s1ck/internal/transport/http/mapper"
	"github.com/internships-backend/test-backend-M0s1ck/internal/transport/http/oapi"
	createroom "github.com/internships-backend/test-backend-M0s1ck/internal/usecase/room/create"
)

type createRoomUsecase interface {
	Execute(ctx context.Context, req *createroom.Request) (*createroom.Response, error)
}

type StrictHandler struct {
	createRoom createRoomUsecase
}

func NewHandler(
	createRoomUC createRoomUsecase,
) *StrictHandler {

	return &StrictHandler{
		createRoom: createRoomUC,
	}
}

func (s StrictHandler) GetRoomsList(ctx context.Context, request oapi.GetRoomsListRequestObject) (oapi.GetRoomsListResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (s StrictHandler) PostRoomsCreate(ctx context.Context, request oapi.PostRoomsCreateRequestObject) (oapi.PostRoomsCreateResponseObject, error) {
	ucReq, err := httpmapper.ToCreateRoomRequest(request.Body)
	if err != nil {
		return oapi.PostRoomsCreate400JSONResponse(
			httphelpers.NewErrorResponse(oapi.INVALIDREQUEST, err.Error()),
		), nil
	}

	ucResp, err := s.createRoom.Execute(ctx, ucReq)
	if err != nil {
		switch {

		case errors.Is(err, room.ErrEmptyName),
			errors.Is(err, room.ErrInvalidCapacity):

			return oapi.PostRoomsCreate400JSONResponse(
				httphelpers.NewErrorResponse(oapi.INVALIDREQUEST, err.Error()),
			), nil

		default:
			return oapi.PostRoomsCreate500JSONResponse(
				httphelpers.NewInternalErrorResponse("create room internal server error"),
			), nil
		}
	}

	return httpmapper.ToCreateRoomResponse(ucResp), nil
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

func newPostRoomsCreateInternalError(message string) oapi.PostRoomsCreate500JSONResponse {
	return oapi.PostRoomsCreate500JSONResponse(
		httphelpers.NewInternalErrorResponse(message),
	)
}

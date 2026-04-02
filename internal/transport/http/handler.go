package http

import (
	"context"
	"errors"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/room"
	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/user"
	"github.com/internships-backend/test-backend-M0s1ck/internal/service/authjwt"
	"github.com/internships-backend/test-backend-M0s1ck/internal/transport/http/helpers"
	"github.com/internships-backend/test-backend-M0s1ck/internal/transport/http/mapper"
	"github.com/internships-backend/test-backend-M0s1ck/internal/transport/http/middleware"
	"github.com/internships-backend/test-backend-M0s1ck/internal/transport/http/oapi"
	dummylogin "github.com/internships-backend/test-backend-M0s1ck/internal/usecase/auth/dummylogin"
	createroom "github.com/internships-backend/test-backend-M0s1ck/internal/usecase/room/create"
)

type createRoomUsecase interface {
	Execute(ctx context.Context, req *createroom.Request, identity *authjwt.Identity) (*createroom.Response, error)
}

type dummyLoginUsecase interface {
	Execute(ctx context.Context, req *dummylogin.Request) (*dummylogin.Response, error)
}

type StrictHandler struct {
	createRoom createRoomUsecase
	dummyLogin dummyLoginUsecase
}

func NewHandler(
	createRoom createRoomUsecase,
	dummyLogin dummyLoginUsecase,
) *StrictHandler {

	return &StrictHandler{
		createRoom: createRoom,
		dummyLogin: dummyLogin,
	}
}

func (s *StrictHandler) PostRoomsCreate(ctx context.Context, request oapi.PostRoomsCreateRequestObject) (oapi.PostRoomsCreateResponseObject, error) {
	ucReq, err := mapper.ToCreateRoomRequest(request.Body)
	if err != nil {
		return oapi.PostRoomsCreate400JSONResponse(
			helpers.NewErrorResponse(oapi.INVALIDREQUEST, err.Error()),
		), nil
	}

	identity := middleware.MustIdentityFromContext(ctx)

	ucResp, err := s.createRoom.Execute(ctx, ucReq, identity)

	if err != nil {
		switch {

		case errors.Is(err, user.ErrAdminRoleRequired):
			return oapi.PostRoomsCreate403JSONResponse(
				helpers.NewErrorResponse(oapi.FORBIDDEN, err.Error()),
			), nil

		case errors.Is(err, room.ErrEmptyName),
			errors.Is(err, room.ErrInvalidCapacity):
			return oapi.PostRoomsCreate400JSONResponse(
				helpers.NewErrorResponse(oapi.INVALIDREQUEST, err.Error()),
			), nil

		default:
			return oapi.PostRoomsCreate500JSONResponse(
				helpers.NewInternalErrorResponse("create room internal server error"),
			), nil
		}
	}

	return mapper.ToCreateRoomResponse(ucResp), nil
}

func (s *StrictHandler) GetRoomsList(ctx context.Context, request oapi.GetRoomsListRequestObject) (oapi.GetRoomsListResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (s *StrictHandler) PostRoomsRoomIdScheduleCreate(ctx context.Context, request oapi.PostRoomsRoomIdScheduleCreateRequestObject) (oapi.PostRoomsRoomIdScheduleCreateResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (s *StrictHandler) GetRoomsRoomIdSlotsList(ctx context.Context, request oapi.GetRoomsRoomIdSlotsListRequestObject) (oapi.GetRoomsRoomIdSlotsListResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (s *StrictHandler) PostBookingsCreate(ctx context.Context, request oapi.PostBookingsCreateRequestObject) (oapi.PostBookingsCreateResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (s *StrictHandler) GetBookingsList(ctx context.Context, request oapi.GetBookingsListRequestObject) (oapi.GetBookingsListResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (s *StrictHandler) GetBookingsMy(ctx context.Context, request oapi.GetBookingsMyRequestObject) (oapi.GetBookingsMyResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (s *StrictHandler) PostBookingsBookingIdCancel(ctx context.Context, request oapi.PostBookingsBookingIdCancelRequestObject) (oapi.PostBookingsBookingIdCancelResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (s *StrictHandler) PostDummyLogin(ctx context.Context, request oapi.PostDummyLoginRequestObject) (oapi.PostDummyLoginResponseObject, error) {
	ucReq, err := mapper.ToDummyLoginRequest(request.Body)
	if err != nil {
		return oapi.PostDummyLogin400JSONResponse(
			helpers.NewErrorResponse(oapi.INVALIDREQUEST, err.Error()),
		), nil
	}

	ucResp, err := s.dummyLogin.Execute(ctx, ucReq)

	if err != nil {
		switch {
		case errors.Is(err, user.ErrInvalidRole):
			return oapi.PostDummyLogin400JSONResponse(
				helpers.NewErrorResponse(oapi.INVALIDREQUEST, err.Error()),
			), nil
		default:
			return oapi.PostDummyLogin500JSONResponse(
				helpers.NewInternalErrorResponse("dummy login internal server error"),
			), nil
		}
	}

	return mapper.ToDummyLoginResponse(ucResp), nil
}

func (s *StrictHandler) PostLogin(ctx context.Context, request oapi.PostLoginRequestObject) (oapi.PostLoginResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (s *StrictHandler) PostRegister(ctx context.Context, request oapi.PostRegisterRequestObject) (oapi.PostRegisterResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

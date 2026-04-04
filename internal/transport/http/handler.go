package http

import (
	"context"
	"errors"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/booking"
	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/room"
	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/schedule"
	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/slot"
	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/user"
	"github.com/internships-backend/test-backend-M0s1ck/internal/service/authjwt"
	"github.com/internships-backend/test-backend-M0s1ck/internal/transport/http/helpers"
	"github.com/internships-backend/test-backend-M0s1ck/internal/transport/http/mapper"
	"github.com/internships-backend/test-backend-M0s1ck/internal/transport/http/middleware"
	"github.com/internships-backend/test-backend-M0s1ck/internal/transport/http/oapi"
	"github.com/internships-backend/test-backend-M0s1ck/internal/usecase/auth/dummylogin"
	createbooking "github.com/internships-backend/test-backend-M0s1ck/internal/usecase/booking/create"
	createroom "github.com/internships-backend/test-backend-M0s1ck/internal/usecase/room/create"
	listroom "github.com/internships-backend/test-backend-M0s1ck/internal/usecase/room/list"
	createschedule "github.com/internships-backend/test-backend-M0s1ck/internal/usecase/schedule/create"
	listslot "github.com/internships-backend/test-backend-M0s1ck/internal/usecase/slot/list"
)

type createRoomUsecase interface {
	Execute(ctx context.Context, req *createroom.Request, identity *authjwt.Identity) (*createroom.Response, error)
}

type listRoomUsecase interface {
	Execute(ctx context.Context) (*listroom.Response, error)
}

type createScheduleUsecase interface {
	Execute(ctx context.Context, req *createschedule.Request, identity *authjwt.Identity) (*createschedule.Response, error)
}

type listSlotUsecase interface {
	Execute(ctx context.Context, req *listslot.Request) (*listslot.Response, error)
}

type createBookingUsecase interface {
	Execute(ctx context.Context, req *createbooking.Request, identity *authjwt.Identity) (*createbooking.Response, error)
}

type dummyLoginUsecase interface {
	Execute(ctx context.Context, req *dummylogin.Request) (*dummylogin.Response, error)
}

type StrictHandler struct {
	createRoom     createRoomUsecase
	listRoom       listRoomUsecase
	createSchedule createScheduleUsecase
	listSlot       listSlotUsecase
	createBooking  createBookingUsecase
	dummyLogin     dummyLoginUsecase
}

func NewHandler(
	createRoom createRoomUsecase,
	listRoom listRoomUsecase,
	createSchedule createScheduleUsecase,
	listSlot listSlotUsecase,
	createBooking createBookingUsecase,
	dummyLogin dummyLoginUsecase,
) *StrictHandler {

	return &StrictHandler{
		createRoom:     createRoom,
		listRoom:       listRoom,
		createSchedule: createSchedule,
		listSlot:       listSlot,
		createBooking:  createBooking,
		dummyLogin:     dummyLogin,
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

func (s *StrictHandler) GetRoomsList(ctx context.Context, _ oapi.GetRoomsListRequestObject) (oapi.GetRoomsListResponseObject, error) {
	ucResp, err := s.listRoom.Execute(ctx)
	if err != nil {
		return oapi.GetRoomsList500JSONResponse(
			helpers.NewInternalErrorResponse("list rooms internal server error"),
		), nil
	}

	return mapper.ToListRoomsResponse(ucResp), nil
}

func (s *StrictHandler) PostRoomsRoomIdScheduleCreate(ctx context.Context, request oapi.PostRoomsRoomIdScheduleCreateRequestObject) (oapi.PostRoomsRoomIdScheduleCreateResponseObject, error) {
	ucReq, err := mapper.ToCreateScheduleRequest(request.RoomId, request.Body)
	if err != nil {
		return oapi.PostRoomsRoomIdScheduleCreate400JSONResponse(
			helpers.NewErrorResponse(oapi.INVALIDREQUEST, err.Error()),
		), nil
	}

	identity := middleware.MustIdentityFromContext(ctx)

	ucResp, err := s.createSchedule.Execute(ctx, ucReq, identity)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrAdminRoleRequired):
			return oapi.PostRoomsRoomIdScheduleCreate403JSONResponse(
				helpers.NewErrorResponse(oapi.FORBIDDEN, err.Error()),
			), nil

		case errors.Is(err, room.ErrNotFound):
			return oapi.PostRoomsRoomIdScheduleCreate404JSONResponse(
				helpers.NewErrorResponse(oapi.ROOMNOTFOUND, err.Error()),
			), nil

		case errors.Is(err, schedule.ErrAlreadyExists):
			return oapi.PostRoomsRoomIdScheduleCreate409JSONResponse(
				helpers.NewErrorResponse(oapi.SCHEDULEEXISTS, err.Error()),
			), nil

		case errors.Is(err, schedule.ErrRoomIDRequired),
			errors.Is(err, schedule.ErrEmptyDaysOfWeek),
			errors.Is(err, schedule.ErrInvalidDayOfWeek),
			errors.Is(err, schedule.ErrInvalidTimeFmt),
			errors.Is(err, schedule.ErrInvalidTimeRange):
			return oapi.PostRoomsRoomIdScheduleCreate400JSONResponse(
				helpers.NewErrorResponse(oapi.INVALIDREQUEST, err.Error()),
			), nil

		default:
			return oapi.PostRoomsRoomIdScheduleCreate500JSONResponse(
				helpers.NewInternalErrorResponse("create schedule internal server error"),
			), nil
		}
	}

	return mapper.ToCreateScheduleResponse(ucResp), nil
}

func (s *StrictHandler) GetRoomsRoomIdSlotsList(ctx context.Context, request oapi.GetRoomsRoomIdSlotsListRequestObject) (oapi.GetRoomsRoomIdSlotsListResponseObject, error) {
	ucReq := mapper.ToListSlotsRequest(request.RoomId, request.Params.Date)

	ucResp, err := s.listSlot.Execute(ctx, ucReq)
	if err != nil {
		switch {
		case errors.Is(err, room.ErrNotFound):
			return oapi.GetRoomsRoomIdSlotsList404JSONResponse(
				helpers.NewErrorResponse(oapi.ROOMNOTFOUND, err.Error()),
			), nil
		default:
			return oapi.GetRoomsRoomIdSlotsList500JSONResponse(
				helpers.NewInternalErrorResponse("list slots internal server error"),
			), nil
		}
	}

	return mapper.ToListSlotsResponse(ucResp), nil
}

func (s *StrictHandler) PostBookingsCreate(ctx context.Context, request oapi.PostBookingsCreateRequestObject) (oapi.PostBookingsCreateResponseObject, error) {
	ucReq, err := mapper.ToCreateBookingRequest(request.Body)
	if err != nil {
		return oapi.PostBookingsCreate400JSONResponse(
			helpers.NewErrorResponse(oapi.INVALIDREQUEST, err.Error()),
		), nil
	}

	identity := middleware.MustIdentityFromContext(ctx)

	ucResp, err := s.createBooking.Execute(ctx, ucReq, identity)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrUserRoleRequired):
			return oapi.PostBookingsCreate403JSONResponse(
				helpers.NewErrorResponse(oapi.FORBIDDEN, err.Error()),
			), nil

		case errors.Is(err, booking.ErrSlotIDRequired),
			errors.Is(err, booking.ErrUserIDRequired),
			errors.Is(err, booking.ErrSlotInPast):
			return oapi.PostBookingsCreate400JSONResponse(
				helpers.NewErrorResponse(oapi.INVALIDREQUEST, err.Error()),
			), nil

		case errors.Is(err, slot.ErrNotFound):
			return oapi.PostBookingsCreate404JSONResponse(
				helpers.NewErrorResponse(oapi.SLOTNOTFOUND, err.Error()),
			), nil

		case errors.Is(err, booking.ErrAlreadyBooked):
			return oapi.PostBookingsCreate409JSONResponse(
				helpers.NewErrorResponse(oapi.SLOTALREADYBOOKED, err.Error()),
			), nil

		default:
			return oapi.PostBookingsCreate500JSONResponse(
				helpers.NewInternalErrorResponse("create booking internal server error"),
			), nil
		}
	}

	return mapper.ToCreateBookingResponse(ucResp), nil
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

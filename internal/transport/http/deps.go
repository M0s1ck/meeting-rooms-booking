package http

import (
	"context"

	"github.com/internships-backend/test-backend-M0s1ck/internal/service/authjwt"
	"github.com/internships-backend/test-backend-M0s1ck/internal/usecase/auth/dummylogin"
	createbooking "github.com/internships-backend/test-backend-M0s1ck/internal/usecase/booking/create"
	listbooking "github.com/internships-backend/test-backend-M0s1ck/internal/usecase/booking/list"
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

type listBookingUsecase interface {
	Execute(ctx context.Context, req *listbooking.Request, identity *authjwt.Identity) (*listbooking.Response, error)
}

type dummyLoginUsecase interface {
	Execute(ctx context.Context, req *dummylogin.Request) (*dummylogin.Response, error)
}

type HandlerDeps struct {
	CreateRoom     createRoomUsecase
	ListRoom       listRoomUsecase
	CreateSchedule createScheduleUsecase
	ListSlot       listSlotUsecase
	CreateBooking  createBookingUsecase
	ListBooking    listBookingUsecase
	DummyLogin     dummyLoginUsecase
}

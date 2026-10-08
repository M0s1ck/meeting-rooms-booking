package http

import (
	"context"

	"github.com/google/uuid"

	"github.com/internships-backend/test-backend-M0s1ck/internal/service/authjwt"
	"github.com/internships-backend/test-backend-M0s1ck/internal/usecase/auth/dummylogin"
	"github.com/internships-backend/test-backend-M0s1ck/internal/usecase/auth/login"
	"github.com/internships-backend/test-backend-M0s1ck/internal/usecase/auth/register"
	cancelbooking "github.com/internships-backend/test-backend-M0s1ck/internal/usecase/booking/cancel"
	createbooking "github.com/internships-backend/test-backend-M0s1ck/internal/usecase/booking/create"
	listbooking "github.com/internships-backend/test-backend-M0s1ck/internal/usecase/booking/list"
	mybooking "github.com/internships-backend/test-backend-M0s1ck/internal/usecase/booking/my"
	createroom "github.com/internships-backend/test-backend-M0s1ck/internal/usecase/room/create"
	listroom "github.com/internships-backend/test-backend-M0s1ck/internal/usecase/room/list"
	createschedule "github.com/internships-backend/test-backend-M0s1ck/internal/usecase/schedule/create"
	listslot "github.com/internships-backend/test-backend-M0s1ck/internal/usecase/slot/list"
	listrangeslot "github.com/internships-backend/test-backend-M0s1ck/internal/usecase/slot/listrange"
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

type listSlotRangeUsecase interface {
	Execute(ctx context.Context, req *listrangeslot.Request) (*listrangeslot.Response, error)
}

type createBookingUsecase interface {
	Execute(ctx context.Context, req *createbooking.Request, identity *authjwt.Identity) (*createbooking.Response, error)
}

type listBookingUsecase interface {
	Execute(ctx context.Context, req *listbooking.Request, identity *authjwt.Identity) (*listbooking.Response, error)
}

type myBookingUsecase interface {
	Execute(ctx context.Context, identity *authjwt.Identity) (*mybooking.Response, error)
}

type cancelBookingUsecase interface {
	Execute(ctx context.Context, bookingID uuid.UUID, identity *authjwt.Identity) (*cancelbooking.Response, error)
}

type dummyLoginUsecase interface {
	Execute(ctx context.Context, req *dummylogin.Request) (*dummylogin.Response, error)
}

type registerUsecase interface {
	Execute(ctx context.Context, req *register.Request) (*register.Response, error)
}

type loginUsecase interface {
	Execute(ctx context.Context, req *login.Request) (*login.Response, error)
}

type HandlerDeps struct {
	CreateRoom     createRoomUsecase
	ListRoom       listRoomUsecase
	CreateSchedule createScheduleUsecase
	ListSlot       listSlotUsecase
	ListSlotRange  listSlotRangeUsecase
	CreateBooking  createBookingUsecase
	ListBooking    listBookingUsecase
	MyBooking      myBookingUsecase
	CancelBooking  cancelBookingUsecase
	DummyLogin     dummyLoginUsecase
	Register       registerUsecase
	Login          loginUsecase
}

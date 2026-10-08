package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/sync/errgroup"

	"github.com/internships-backend/test-backend-M0s1ck/internal/config"
	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/slot"
	infrabcrypt "github.com/internships-backend/test-backend-M0s1ck/internal/infra/bcrypt"
	"github.com/internships-backend/test-backend-M0s1ck/internal/infra/cron"
	"github.com/internships-backend/test-backend-M0s1ck/internal/infra/postgres"
	"github.com/internships-backend/test-backend-M0s1ck/internal/infra/postgres/repository"
	"github.com/internships-backend/test-backend-M0s1ck/internal/infra/thirdparty/conferenceservice"
	"github.com/internships-backend/test-backend-M0s1ck/internal/service/authjwt"
	httpapi "github.com/internships-backend/test-backend-M0s1ck/internal/transport/http"
	httphelpers "github.com/internships-backend/test-backend-M0s1ck/internal/transport/http/helpers"
	appmiddleware "github.com/internships-backend/test-backend-M0s1ck/internal/transport/http/middleware"
	"github.com/internships-backend/test-backend-M0s1ck/internal/transport/http/oapi"
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
	"github.com/internships-backend/test-backend-M0s1ck/internal/usecase/slot/ensureroomdate"
	"github.com/internships-backend/test-backend-M0s1ck/internal/usecase/slot/fillhorizon"
	listslot "github.com/internships-backend/test-backend-M0s1ck/internal/usecase/slot/list"
	listrangeslot "github.com/internships-backend/test-backend-M0s1ck/internal/usecase/slot/listrange"
)

// Keys of PostgreSQL advisory locks for background jobs.
// Any constant int64 works; they only have to be unique within the database.
const (
	lockKeySlotHorizonFillTail   int64 = 1001
	lockKeySlotHorizonFullRepair int64 = 1002
)

type App struct {
	HttpSrv           *http.Server
	db                *pgxpool.Pool
	cronSched         *cron.Scheduler
	slotHorizonRepair func(context.Context) error
	logger            *slog.Logger
}

func Build(ctx context.Context, cfg *config.Config, logger *slog.Logger) (*App, error) {
	db, err := postgres.Connect(ctx, cfg.Psg)
	if err != nil {
		return nil, err
	}

	txManager := manager.Must(trmpgx.NewDefaultFactory(db))
	txGetter := trmpgx.DefaultCtxGetter

	tokenManager := authjwt.NewManager(cfg.JwtCfg)
	passHasher := infrabcrypt.NewHasher(cfg.BcryptCfg)
	confLinkProvider := conferenceservice.NewMockProvider(cfg.ConferenceCfg.MockFailureRate)

	roomRepo := repository.NewRoomRepo(db, txGetter)
	shedRepo := repository.NewScheduleRepo(db, txGetter)
	slotRepo := repository.NewSlotRepo(db, txGetter)
	userRepo := repository.NewUserRepo(db, txGetter)
	bookingRepo := repository.NewBookingRepo(db, txGetter)

	slotGen := slot.NewGenerator()

	slotHorizon := 14 * 24 * time.Hour
	if cfg.SlotCfg != nil && cfg.SlotCfg.Horizon > 0 {
		slotHorizon = cfg.SlotCfg.Horizon
	}

	createRoom := createroom.NewUsecase(roomRepo)
	listRoom := listroom.NewUsecase(roomRepo)
	createSchedule := createschedule.NewUsecase(slotGen, shedRepo, slotRepo, txManager, slotHorizon)
	fillSlotHorizon := fillhorizon.NewUsecase(slotGen, slotRepo, shedRepo, txManager, slotHorizon)
	ensureRoomDate := ensureroomdate.NewUsecase(slotGen, shedRepo, slotRepo, txManager)
	listSlots := listslot.NewUsecase(roomRepo, slotRepo, ensureRoomDate, slotHorizon)
	listSlotsRange := listrangeslot.NewUsecase(roomRepo, slotRepo, ensureRoomDate, slotHorizon)
	createBooking := createbooking.NewUsecase(bookingRepo, slotRepo, confLinkProvider, logger)
	listBooking := listbooking.NewUsecase(bookingRepo)
	myBooking := mybooking.NewUsecase(bookingRepo)
	cancelBooking := cancelbooking.NewUsecase(bookingRepo)
	dummyLogin := dummylogin.NewUsecase(tokenManager, userRepo)
	reg := register.NewUsecase(userRepo, passHasher)
	logIn := login.NewUsecase(userRepo, tokenManager, passHasher)

	// background jobs run in every replica, but the advisory lock guarantees
	// that only one of them does the work at a time
	locker := postgres.NewAdvisoryLocker(db, logger)
	fillTailJob := locker.WithLock(lockKeySlotHorizonFillTail, fillSlotHorizon.FillTail)
	fullRepairJob := locker.WithLock(lockKeySlotHorizonFullRepair, fillSlotHorizon.FullRepair)

	cronSched, err := cron.New(logger, cfg.SlotCfg.HorizonRepairTZ)
	if err != nil {
		return nil, err
	}

	err = cronSched.AddJob(ctx,
		cfg.SlotCfg.HorizonFillTailInterval,
		"fill slot horizon tail",
		fillTailJob,
	)
	if err != nil {
		return nil, err
	}

	err = cronSched.AddDailyJobAt(ctx,
		cfg.SlotCfg.HorizonRepairHour, cfg.SlotCfg.HorizonRepairMinute,
		"slot horizon full repair",
		fullRepairJob,
	)
	if err != nil {
		return nil, err
	}

	handler := httpapi.NewHandler(httpapi.HandlerDeps{
		CreateRoom:     createRoom,
		ListRoom:       listRoom,
		CreateSchedule: createSchedule,
		ListSlot:       listSlots,
		ListSlotRange:  listSlotsRange,
		CreateBooking:  createBooking,
		ListBooking:    listBooking,
		MyBooking:      myBooking,
		CancelBooking:  cancelBooking,
		DummyLogin:     dummyLogin,
		Register:       reg,
		Login:          logIn,
	})

	router := chi.NewRouter()

	router.Use(
		middleware.RequestID,
		middleware.RealIP,
		appmiddleware.RequestLogger(logger),
	)

	strictMiddlewares := appmiddleware.NewStrictMiddlewares(tokenManager, logger)

	oapi.HandlerFromMux(
		oapi.NewStrictHandler(handler, strictMiddlewares),
		router,
	)

	httphelpers.AddInfo(router)

	httpSrv := &http.Server{
		Addr:         cfg.HttpCfg.Addr,
		Handler:      router,
		ReadTimeout:  cfg.HttpCfg.ReadTimeout,
		WriteTimeout: cfg.HttpCfg.WriteTimeout,
	}

	return &App{
		HttpSrv:           httpSrv,
		db:                db,
		slotHorizonRepair: fullRepairJob,
		cronSched:         cronSched,
		logger:            logger,
	}, nil
}

func (app *App) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		app.logger.Info("initial slot horizon repair started")
		if err := app.slotHorizonRepair(ctx); err != nil {
			return err
		}

		app.logger.Info("initial slot horizon repair finished")
		app.cronSched.Start()
		<-ctx.Done()
		return nil
	})

	g.Go(func() error {
		app.logger.Info("http server starting", "addr", app.HttpSrv.Addr)

		err := app.HttpSrv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}

		return nil
	})

	return g.Wait()
}

func (app *App) Shutdown(shutdownCtx context.Context) error {
	app.logger.Info("Gracefully shutting down...")

	// stop accepting requests and let in-flight ones finish
	httpErr := app.HttpSrv.Shutdown(shutdownCtx)

	// wait for running background jobs, then release DB connections
	cronErr := app.cronSched.Shutdown()
	app.db.Close()

	if err := errors.Join(httpErr, cronErr); err != nil {
		return err
	}

	app.logger.Info("server stopped")
	return nil
}

package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/internships-backend/test-backend-M0s1ck/internal/config"
	"github.com/internships-backend/test-backend-M0s1ck/internal/infra/postgres"
	"github.com/internships-backend/test-backend-M0s1ck/internal/infra/postgres/repository"
	"github.com/internships-backend/test-backend-M0s1ck/internal/service/authjwt"
	httpapi "github.com/internships-backend/test-backend-M0s1ck/internal/transport/http"
	appmiddleware "github.com/internships-backend/test-backend-M0s1ck/internal/transport/http/middleware"
	"github.com/internships-backend/test-backend-M0s1ck/internal/transport/http/oapi"
	"github.com/internships-backend/test-backend-M0s1ck/internal/usecase/auth/dummylogin"
	createroom "github.com/internships-backend/test-backend-M0s1ck/internal/usecase/room/create"
	listroom "github.com/internships-backend/test-backend-M0s1ck/internal/usecase/room/list"
)

type App struct {
	httpSrv *http.Server

	logger *slog.Logger
}

func Build(ctx context.Context, cfg *config.Config, logger *slog.Logger) (*App, error) {
	db, err := postgres.Connect(ctx, cfg.Psg)
	if err != nil {
		return nil, err
	}

	//  txManager := manager.Must(trmpgx.NewDefaultFactory(db))
	txGetter := trmpgx.DefaultCtxGetter

	tokenManager := authjwt.NewManager(cfg.JwtCfg)

	roomRepo := repository.NewRoomRepo(db, txGetter)

	dummyLogin := dummylogin.NewUsecase(tokenManager)
	createRoom := createroom.NewUsecase(roomRepo)
	listRoom := listroom.NewUsecase(roomRepo)

	handler := httpapi.NewHandler(
		createRoom,
		listRoom,
		dummyLogin,
	)

	router := chi.NewRouter()

	router.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
	)

	strictMiddlewares := appmiddleware.NewStrictMiddlewares(tokenManager, logger)

	oapi.HandlerFromMux(
		oapi.NewStrictHandler(handler, strictMiddlewares),
		router,
	)

	httpSrv := &http.Server{
		Addr:         cfg.HttpCfg.Addr,
		Handler:      router,
		ReadTimeout:  cfg.HttpCfg.ReadTimeout,
		WriteTimeout: cfg.HttpCfg.WriteTimeout,
	}

	return &App{
		httpSrv: httpSrv,
		logger:  logger,
	}, nil
}

func (app *App) Run(ctx context.Context) error {
	app.logger.Info("http server starting", "addr", app.httpSrv.Addr)

	err := app.httpSrv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (app *App) Shutdown(shutdownCtx context.Context) error {
	app.logger.Info("Gracefully shutting down...")

	if err := app.httpSrv.Shutdown(shutdownCtx); err != nil {
		return err
	}

	app.logger.Info("server stopped")
	return nil
}

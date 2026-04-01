package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/internships-backend/test-backend-M0s1ck/internal/config"
	"github.com/internships-backend/test-backend-M0s1ck/internal/infra/postgres"
	httpapi "github.com/internships-backend/test-backend-M0s1ck/internal/transport/http"
	"github.com/internships-backend/test-backend-M0s1ck/internal/transport/http/oapi"
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

	_ = db

	handler := httpapi.NewHandler()

	router := chi.NewRouter()

	router.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
	)

	oapi.HandlerFromMux(
		oapi.NewStrictHandler(handler, nil),
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

package app

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/internships-backend/test-backend-M0s1ck/internal/config"
	httpapi "github.com/internships-backend/test-backend-M0s1ck/internal/transport/http"
	"github.com/internships-backend/test-backend-M0s1ck/internal/transport/http/oapi"
)

type App struct {
	httpSrv *http.Server
}

func Build(ctx context.Context, cfg *config.Config, logger *slog.Logger) *App {
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
	}
}

package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	// embed the IANA time zone database into the binary, so the app does not
	// depend on tzdata being installed in the OS image (12-factor, II. Dependencies)
	_ "time/tzdata"

	application "github.com/internships-backend/test-backend-M0s1ck/internal/app"
	"github.com/internships-backend/test-backend-M0s1ck/internal/config"
	infralog "github.com/internships-backend/test-backend-M0s1ck/internal/infra/log"
)

func main() {
	conf, err := config.Load()
	if err != nil {
		infralog.NewSlogger(slog.LevelInfo).Error("couldn't load config", "error", err)
		os.Exit(1)
	}

	logger := infralog.NewSlogger(conf.AppCfg.LogLevel)
	logger.Info("Service is starting...")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	app, err := application.Build(ctx, conf, logger)
	if err != nil {
		logger.Error("couldn't build app", "error", err)
		os.Exit(1)
	}

	// buffered, so Run's goroutine never blocks if we exit by signal
	runErrCh := make(chan error, 1)

	go func() {
		runErrCh <- app.Run(ctx)
	}()

	exitCode := 0

	select {
	case <-ctx.Done():
		logger.Info("shutdown signal received")
	case runErr := <-runErrCh:
		if runErr != nil {
			logger.Error("app stopped with error", "err", runErr)
			exitCode = 1
		}
		stop()
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), conf.AppCfg.ShutdownTimeout)
	defer cancel()

	err = app.Shutdown(shutdownCtx)
	if err != nil {
		logger.Error("app stopped with error", "error", err)
		exitCode = 1
	}

	if exitCode != 0 {
		os.Exit(exitCode)
	}
}

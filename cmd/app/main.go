package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	application "github.com/internships-backend/test-backend-M0s1ck/internal/app"
	"github.com/internships-backend/test-backend-M0s1ck/internal/config"
	infralog "github.com/internships-backend/test-backend-M0s1ck/internal/infra/log"
)

func main() {
	logger := infralog.NewSlogger()
	logger.Info("Service is starting...")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	conf := config.Load()

	app, err := application.Build(ctx, conf, logger)
	if err != nil {
		logger.Error("couldn't build app", "error", err)
		os.Exit(1)
	}

	runErrCh := make(chan error)

	go func() {
		runErrCh <- app.Run(ctx)
	}()

	exitCode := 0

	select {
	case <-ctx.Done():
	case runErr := <-runErrCh:
		if runErr != nil {
			logger.Error("app stopped with error", "error", runErr)
			exitCode = 1
		}
		stop()
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

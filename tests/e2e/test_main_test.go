package e2e_test

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"golang.org/x/crypto/bcrypt"

	application "github.com/internships-backend/test-backend-M0s1ck/internal/app"
	"github.com/internships-backend/test-backend-M0s1ck/internal/config"
)

var (
	app     *application.App
	baseURL string
)

func TestMain(m *testing.M) {
	code := run(m)
	os.Exit(code)
}

func run(m *testing.M) int {
	ctx := context.Background()

	logger := slog.Default()

	dbName := "test_db"
	dbUser := "test_user"
	dbPass := "test_pass"

	postgresContainer, err := postgres.Run(ctx,
		"postgres:17-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPass),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		logger.Error("failed to start postgres container", "err", err)
		return 1
	}

	defer func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			logger.Warn("failed to terminate postgres container", "err", err)
		}
	}()

	pgHost, err := postgresContainer.Host(ctx)
	if err != nil {
		logger.Error("failed to get postgres host", "err", err)
		return 1
	}

	pgPort, err := postgresContainer.MappedPort(ctx, "5432/tcp")
	if err != nil {
		logger.Error("failed to get postgres port", "err", err)
		return 1
	}

	pgSPort := pgPort.Port()

	logger.Info("postgres ready", "host", pgHost, "port", pgSPort)

	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		dbUser, dbPass, net.JoinHostPort(pgHost, pgSPort), dbName,
	)

	if err := runMigrations(dbURL); err != nil {
		logger.Error("failed to run postgres migrations", "err", err)
		return 1
	}

	loc, _ := time.LoadLocation("Europe/Moscow")

	conf := &config.Config{
		HttpCfg: &config.HttpCfg{
			Addr:         ":8080",
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		Psg: &config.Postgres{
			Host:         pgHost,
			Port:         pgSPort,
			Name:         dbName,
			User:         dbUser,
			Password:     dbPass,
			SSLMode:      "disable",
			MaxPoolConns: 10,
			MinPoolConns: 2,
		},
		JwtCfg: &config.JwtCfg{
			Secret: "secret",
			Ttl:    30 * time.Minute,
		},
		SlotCfg: &config.SlotCfg{
			Horizon:                 0,
			HorizonRepairTZ:         *loc,
			HorizonRepairHour:       1,
			HorizonRepairMinute:     0,
			HorizonFillTailInterval: time.Hour,
		},
		BcryptCfg: &config.BcryptCfg{
			Cost: bcrypt.DefaultCost,
		},
	}

	app, err = application.Build(ctx, conf, logger)
	if err != nil {
		logger.Error("failed to build app", "err", err)
		return 1
	}

	server := httptest.NewServer(app.HttpSrv.Handler)
	baseURL = server.URL

	logger.Info("test environment ready")

	code := m.Run()

	server.Close()
	return code
}

func runMigrations(dbURL string) error {
	m, err := migrate.New(
		"file://../../internal/infra/postgres/migrations",
		dbURL,
	)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

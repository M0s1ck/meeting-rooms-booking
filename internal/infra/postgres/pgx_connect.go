package postgres

import (
	"context"
	"crypto/tls"
	"fmt"
	"strconv"

	"github.com/internships-backend/test-backend-M0s1ck/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(ctx context.Context, cfg *config.Postgres) (*pgxpool.Pool, error) {
	pgxCfg, err := buildPgxConfFromAppConf(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to build pgxpool config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxCfg)
	if err != nil {
		return nil, fmt.Errorf("postgres pgx pool init: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("postgres pgx ping: %w", err)
	}

	return pool, nil
}

func buildPgxConfFromAppConf(cfg *config.Postgres) (*pgxpool.Config, error) {
	pgxCfg, err := pgxpool.ParseConfig("")
	if err != nil {
		return nil, err
	}

	cc := pgxCfg.ConnConfig

	port, err := strconv.Atoi(cfg.Port)
	if err != nil {
		port = 5432
	}

	if cfg.SSLMode == "disable" {
		cc.TLSConfig = nil
	} else {
		cc.TLSConfig = &tls.Config{}
	}

	cc.User = cfg.User
	cc.Password = cfg.Password
	cc.Host = cfg.Host
	cc.Port = uint16(port)
	cc.Database = cfg.Name

	pgxCfg.MaxConns = int32(cfg.MaxPoolConns)
	pgxCfg.MinConns = int32(cfg.MinPoolConns)

	return pgxCfg, nil
}

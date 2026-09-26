package config

import (
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func setRequired(t *testing.T) {
	t.Setenv("JWT_SECRET", "secret")
	t.Setenv("POSTGRES_HOST", "db")
	t.Setenv("POSTGRES_DB", "db")
	t.Setenv("POSTGRES_USER", "u")
	t.Setenv("POSTGRES_PASSWORD", "p")
}

func TestLoad_Defaults(t *testing.T) {
	setRequired(t)

	cfg, err := Load()
	require.NoError(t, err)
	require.Equal(t, ":8080", cfg.HttpCfg.Addr)
	require.Equal(t, 14*24*time.Hour, cfg.SlotCfg.Horizon)
	require.Equal(t, "Europe/Moscow", cfg.SlotCfg.HorizonRepairTZ.String())
	require.Equal(t, uint(3), cfg.SlotCfg.HorizonRepairHour)
	require.Equal(t, slog.LevelInfo, cfg.AppCfg.LogLevel)
	require.InDelta(t, 0.05, cfg.ConferenceCfg.MockFailureRate, 1e-9)
}

func TestLoad_Overrides(t *testing.T) {
	setRequired(t)
	t.Setenv("HTTP_ADDR", ":9090")
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("SLOT_HORIZON_REPAIR_TIME", "04:30")
	t.Setenv("CONFERENCE_MOCK_FAILURE_RATE", "0")

	cfg, err := Load()
	require.NoError(t, err)
	require.Equal(t, ":9090", cfg.HttpCfg.Addr)
	require.Equal(t, slog.LevelDebug, cfg.AppCfg.LogLevel)
	require.Equal(t, uint(4), cfg.SlotCfg.HorizonRepairHour)
	require.Equal(t, uint(30), cfg.SlotCfg.HorizonRepairMinute)
	require.Zero(t, cfg.ConferenceCfg.MockFailureRate)
}

func TestLoad_MissingRequired(t *testing.T) {
	t.Setenv("JWT_SECRET", "")
	t.Setenv("POSTGRES_HOST", "")

	_, err := Load()
	require.ErrorContains(t, err, "JWT_SECRET is required")
	require.ErrorContains(t, err, "POSTGRES_HOST is required")
}

func TestLoad_Malformed(t *testing.T) {
	setRequired(t)
	t.Setenv("JWT_TTL", "thirty minutes")
	t.Setenv("SLOT_HORIZON_REPAIR_TZ", "Mars/Olympus")

	_, err := Load()
	require.ErrorContains(t, err, "JWT_TTL")
	require.ErrorContains(t, err, "SLOT_HORIZON_REPAIR_TZ")
}

package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"
)

// Config is the whole application configuration.
// It is read only from environment variables (12-factor, III. Config).
type Config struct {
	AppCfg        *AppCfg
	HttpCfg       *HttpCfg
	Psg           *Postgres
	JwtCfg        *JwtCfg
	SlotCfg       *SlotCfg
	BcryptCfg     *BcryptCfg
	ConferenceCfg *ConferenceCfg
}

// Load reads the configuration from the environment.
// Missing required variables and malformed values are reported as an error,
// so the process fails fast instead of silently running with a wrong config.
func Load() (*Config, error) {
	r := &envReader{}

	cfg := &Config{
		AppCfg:        loadAppCfg(r),
		HttpCfg:       loadHttpCfg(r),
		Psg:           loadPsgCfg(r),
		JwtCfg:        loadJwtCfg(r),
		SlotCfg:       loadSlotCfg(r),
		BcryptCfg:     loadBcryptCfg(r),
		ConferenceCfg: loadConferenceCfg(r),
	}

	if err := errors.Join(r.errs...); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// envReader reads typed values from environment variables and collects all errors.
type envReader struct {
	errs []error
}

func (r *envReader) fail(key, val string, err error) {
	r.errs = append(r.errs, fmt.Errorf("%s=%q: %w", key, val, err))
}

func (r *envReader) str(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func (r *envReader) required(key string) string {
	v := os.Getenv(key)
	if v == "" {
		r.errs = append(r.errs, fmt.Errorf("%s is required", key))
	}
	return v
}

func (r *envReader) int(key string, def int) int {
	val := os.Getenv(key)
	if val == "" {
		return def
	}

	num, err := strconv.Atoi(val)
	if err != nil {
		r.fail(key, val, err)
		return def
	}
	return num
}

func (r *envReader) float(key string, def float64) float64 {
	val := os.Getenv(key)
	if val == "" {
		return def
	}

	num, err := strconv.ParseFloat(val, 64)
	if err != nil {
		r.fail(key, val, err)
		return def
	}
	return num
}

func (r *envReader) duration(key string, def time.Duration) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		return def
	}

	d, err := time.ParseDuration(val)
	if err != nil {
		r.fail(key, val, err)
		return def
	}
	return d
}

func (r *envReader) location(key, def string) time.Location {
	val := r.str(key, def)

	loc, err := time.LoadLocation(val)
	if err != nil {
		r.fail(key, val, err)
		return *time.UTC
	}
	return *loc
}

// clock parses "HH:MM".
func (r *envReader) clock(key, def string) (hour, minute uint) {
	val := r.str(key, def)

	t, err := time.Parse("15:04", val)
	if err != nil {
		r.fail(key, val, err)
		return 0, 0
	}
	return uint(t.Hour()), uint(t.Minute())
}

func (r *envReader) logLevel(key string, def slog.Level) slog.Level {
	val := os.Getenv(key)
	if val == "" {
		return def
	}

	var lvl slog.Level
	if err := lvl.UnmarshalText([]byte(val)); err != nil {
		r.fail(key, val, err)
		return def
	}
	return lvl
}

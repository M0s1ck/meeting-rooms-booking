package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	HttpCfg   *HttpCfg
	Psg       *Postgres
	JwtCfg    *JwtCfg
	SlotCfg   *SlotCfg
	BcryptCfg *BcryptCfg
}

func Load() *Config {
	return &Config{
		HttpCfg:   loadHttpCfg(),
		Psg:       loadPsgCfg(),
		JwtCfg:    loadJwtCfg(),
		SlotCfg:   loadSlotCfg(),
		BcryptCfg: loadBcryptCfg(),
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	if str := os.Getenv(key); str != "" {
		num, err := strconv.Atoi(str)
		if err != nil {
			return def
		}
		return num
	}
	return def
}

func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}

	d, err := time.ParseDuration(val)
	if err != nil {
		return defaultVal
	}

	return d
}

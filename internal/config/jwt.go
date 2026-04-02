package config

import (
	"os"
	"time"
)

type JwtCfg struct {
	Secret string
	Issuer string
	Ttl    time.Duration
}

func loadJwtCfg() *JwtCfg {
	return &JwtCfg{
		Secret: os.Getenv("JWT_SECRET"),
		Ttl:    getEnvDuration("JWT_TTL", time.Minute*30),
		Issuer: getEnv("JWT_ISSUER", "meeting-scheduler"),
	}
}

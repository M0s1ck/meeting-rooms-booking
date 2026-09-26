package config

import (
	"time"
)

type JwtCfg struct {
	Secret string
	Issuer string
	Ttl    time.Duration
}

func loadJwtCfg(r *envReader) *JwtCfg {
	return &JwtCfg{
		Secret: r.required("JWT_SECRET"),
		Ttl:    r.duration("JWT_TTL", 30*time.Minute),
		Issuer: r.str("JWT_ISSUER", "meeting-scheduler"),
	}
}

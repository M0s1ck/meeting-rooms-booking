package config

import (
	"time"
)

type HttpCfg struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func loadHttpCfg() *HttpCfg {
	return &HttpCfg{
		Addr:         getEnv("HTTP_ADDR", ":8080"),
		ReadTimeout:  getEnvDuration("HTTP_READ_TIMEOUT", 5*time.Second),
		WriteTimeout: getEnvDuration("HTTP_WRITE_TIMEOUT", 10*time.Second),
	}
}

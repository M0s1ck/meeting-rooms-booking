package config

import (
	"time"
)

type HttpCfg struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func loadHttpCfg(r *envReader) *HttpCfg {
	return &HttpCfg{
		Addr:         r.str("HTTP_ADDR", ":8080"),
		ReadTimeout:  r.duration("HTTP_READ_TIMEOUT", 5*time.Second),
		WriteTimeout: r.duration("HTTP_WRITE_TIMEOUT", 10*time.Second),
	}
}

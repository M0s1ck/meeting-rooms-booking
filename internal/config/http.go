package config

import "time"

type HttpCfg struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

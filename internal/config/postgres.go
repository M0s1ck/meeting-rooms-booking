package config

import "os"

type Postgres struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string

	MaxPoolConns int
	MinPoolConns int
}

func loadPsgCfg() *Postgres {
	return &Postgres{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     getEnv("POSTGRES_PORT", "5432"),
		Name:     os.Getenv("POSTGRES_DB"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		SSLMode:  getEnv("POSTGRES_SSL_MODE", "disable"),

		MaxPoolConns: getEnvInt("POSTGRES_MAX_POOL_CONNS", 10),
		MinPoolConns: getEnvInt("POSTGRES_MIN_POOL_CONNS", 2),
	}
}

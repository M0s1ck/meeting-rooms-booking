package config

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

func loadPsgCfg(r *envReader) *Postgres {
	return &Postgres{
		Host:     r.required("POSTGRES_HOST"),
		Port:     r.str("POSTGRES_PORT", "5432"),
		Name:     r.required("POSTGRES_DB"),
		User:     r.required("POSTGRES_USER"),
		Password: r.required("POSTGRES_PASSWORD"),
		SSLMode:  r.str("POSTGRES_SSL_MODE", "disable"),

		MaxPoolConns: r.int("POSTGRES_MAX_POOL_CONNS", 10),
		MinPoolConns: r.int("POSTGRES_MIN_POOL_CONNS", 2),
	}
}

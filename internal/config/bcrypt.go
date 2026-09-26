package config

type BcryptCfg struct {
	Cost int
}

func loadBcryptCfg(r *envReader) *BcryptCfg {
	return &BcryptCfg{
		Cost: r.int("BCRYPT_COST", 10),
	}
}

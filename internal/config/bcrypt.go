package config

type BcryptCfg struct {
	Cost int
}

func loadBcryptCfg() *BcryptCfg {
	return &BcryptCfg{
		Cost: getEnvInt("BCRYPT_COST", 10),
	}
}

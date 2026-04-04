package config

import "time"

type SlotCfg struct {
	Horizon time.Duration
}

func loadSlotCfg() *SlotCfg {
	return &SlotCfg{
		Horizon: getEnvDuration("SLOT_HORIZON", 14*24*time.Hour),
	}
}

package config

import (
	"time"
)

type SlotCfg struct {
	Horizon                 time.Duration
	HorizonRepairTZ         time.Location
	HorizonRepairHour       uint
	HorizonRepairMinute     uint
	HorizonFillTailInterval time.Duration
}

func loadSlotCfg() *SlotCfg {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		loc = time.Local
	}

	hour, minute, err := getEnvClock("SLOT_HORIZON_REPAIR_TIME")
	if err != nil {
		hour = 3 // 03:00
		minute = 0
	}

	return &SlotCfg{
		Horizon:                 getEnvDuration("SLOT_HORIZON", 14*24*time.Hour),
		HorizonRepairTZ:         getEnvTimeZone("SLOT_HORIZON_REPAIR_TZ", *loc),
		HorizonRepairHour:       hour,
		HorizonRepairMinute:     minute,
		HorizonFillTailInterval: getEnvDuration("SLOT_HORIZON_FILLTAIL_INTERVAL", time.Hour),
	}
}

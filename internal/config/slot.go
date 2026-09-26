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

func loadSlotCfg(r *envReader) *SlotCfg {
	hour, minute := r.clock("SLOT_HORIZON_REPAIR_TIME", "03:00")

	return &SlotCfg{
		Horizon:                 r.duration("SLOT_HORIZON", 14*24*time.Hour),
		HorizonRepairTZ:         r.location("SLOT_HORIZON_REPAIR_TZ", "Europe/Moscow"),
		HorizonRepairHour:       hour,
		HorizonRepairMinute:     minute,
		HorizonFillTailInterval: r.duration("SLOT_HORIZON_FILLTAIL_INTERVAL", time.Hour),
	}
}

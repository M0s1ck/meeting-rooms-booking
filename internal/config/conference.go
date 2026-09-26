package config

type ConferenceCfg struct {
	// MockFailureRate is the share of requests (0..1) the mock conference service fails.
	MockFailureRate float64
}

func loadConferenceCfg(r *envReader) *ConferenceCfg {
	return &ConferenceCfg{
		MockFailureRate: r.float("CONFERENCE_MOCK_FAILURE_RATE", 0.05),
	}
}

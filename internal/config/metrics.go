package config

// MetricsCfg configures the internal metrics endpoint.
// It runs on a separate port from the public API (HTTP_ADDR).
type MetricsCfg struct {
	Addr string
}

func loadMetricsCfg(r *envReader) *MetricsCfg {
	return &MetricsCfg{
		Addr: r.str("METRICS_ADDR", ":9090"),
	}
}

package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// NewServer builds a separate HTTP server that only serves /metrics.
// It listens on its own port, so the endpoint is not routed through the public
// API (Gateway/HTTPRoute) and is reachable only from inside the cluster.
func NewServer(addr string, m *Metrics) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(m.Registry, promhttp.HandlerOpts{
		Registry:          m.Registry, // exposes promhttp_metric_handler_errors_total
		EnableOpenMetrics: true,       // OpenMetrics format when the scraper asks for it
	}))

	return &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
}

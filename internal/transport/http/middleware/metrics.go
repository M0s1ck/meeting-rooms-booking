package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/internships-backend/test-backend-M0s1ck/internal/infra/metrics"
)

// HTTPMetrics counts requests and measures latency per route.
// The chi route pattern (e.g. /rooms/{roomId}/slots/list) is used as a label
// instead of the raw path, so ids do not blow up metric cardinality.
func HTTPMetrics(m *metrics.Metrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == healthPath {
				next.ServeHTTP(w, r)
				return
			}

			start := time.Now()
			ww := chimw.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			route := "unmatched"
			if rc := chi.RouteContext(r.Context()); rc != nil {
				if p := rc.RoutePattern(); p != "" {
					route = p
				}
			}

			status := ww.Status()
			if status == 0 {
				status = http.StatusOK
			}

			m.HTTPRequests.WithLabelValues(r.Method, route, strconv.Itoa(status)).Inc()
			m.HTTPRequestDuration.WithLabelValues(r.Method, route).Observe(time.Since(start).Seconds())
		})
	}
}

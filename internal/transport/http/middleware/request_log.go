package middleware

import (
	"log/slog"
	"net/http"
	"time"

	chimw "github.com/go-chi/chi/v5/middleware"
)

// healthPath is hit by k8s probes every few seconds, logging it would drown real requests.
const healthPath = "/_info"

// RequestLogger logs every HTTP request with the application's slog logger,
// so access logs go to the same stdout JSON stream as the rest of the logs.
func RequestLogger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == healthPath {
				next.ServeHTTP(w, r)
				return
			}

			start := time.Now()
			ww := chimw.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			logger.Log(r.Context(), levelForStatus(ww.Status()), "http request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", ww.Status(),
				"bytes", ww.BytesWritten(),
				"duration_ms", float64(time.Since(start).Microseconds())/1000,
				"remote_addr", r.RemoteAddr,
				"request_id", chimw.GetReqID(r.Context()),
			)
		})
	}
}

// levelForStatus makes server failures searchable by level, not only by status.
func levelForStatus(status int) slog.Level {
	switch {
	case status >= 500:
		return slog.LevelError
	case status >= 400:
		return slog.LevelWarn
	default:
		return slog.LevelInfo
	}
}

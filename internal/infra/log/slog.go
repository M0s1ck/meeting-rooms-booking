package log

import (
	"log/slog"
	"os"
	"strings"
	"time"
)

// ServiceName is attached to every record as service.name,
// so logs of this service can be told apart in a shared log storage.
const ServiceName = "meeting-rooms-booking"

// NewSlogger writes structured JSON logs to stdout as an unbuffered event stream
// (12-factor, XI. Logs). Routing and storage are the environment's job.
//
// Record shape: {"timestamp":"...","level":"info","message":"...","service.name":"...",...attrs}
func NewSlogger(level slog.Level) *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:       level,
		ReplaceAttr: renameBuiltinAttrs,
	})

	return slog.New(handler).With(slog.String("service.name", ServiceName))
}

// renameBuiltinAttrs maps slog's built-in keys to the common log schema:
// time -> timestamp (RFC3339 in UTC), msg -> message, level -> lowercase level.
func renameBuiltinAttrs(groups []string, a slog.Attr) slog.Attr {
	if len(groups) > 0 {
		return a
	}

	switch a.Key {
	case slog.TimeKey:
		return slog.String("timestamp", a.Value.Time().UTC().Format(time.RFC3339Nano))
	case slog.MessageKey:
		return slog.Attr{Key: "message", Value: a.Value}
	case slog.LevelKey:
		return slog.String(slog.LevelKey, strings.ToLower(a.Value.String()))
	}

	return a
}

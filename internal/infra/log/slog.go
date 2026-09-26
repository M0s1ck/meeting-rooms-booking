package log

import (
	"log/slog"
	"os"
)

// NewSlogger writes structured JSON logs to stdout as an unbuffered event stream
// (12-factor, XI. Logs). Routing and storage are the environment's job.
func NewSlogger(level slog.Level) *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})

	return slog.New(handler)
}

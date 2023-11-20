package log

import (
	"io"
	"log/slog"
)

// NewLogger returns a new instance of slog.Logger that handles passing values via context and outputting opentelemetry trace information.
func NewLogger(out io.Writer) *slog.Logger {
	handler := &TracingHandler{
		handler: slog.NewJSONHandler(out, &slog.HandlerOptions{AddSource: true}),
	}
	logger := slog.New(handler)

	slog.SetDefault(logger)
	return logger
}

package log

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/trace"
)

var _ slog.Handler = (*TracingHandler)(nil)

const (
	traceIDKey = "trace_id"
	spanIDKey  = "span_id"
)

type TracingHandler struct {
	handler slog.Handler
}

func (h *TracingHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *TracingHandler) Handle(ctx context.Context, r slog.Record) error {
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		sc := span.SpanContext()
		r.AddAttrs(
			slog.String(traceIDKey, sc.TraceID().String()),
			slog.String(spanIDKey, sc.SpanID().String()),
		)
	}
	return h.handler.Handle(ctx, r)
}

func (h *TracingHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &TracingHandler{h.handler.WithAttrs(attrs)}
}

func (h *TracingHandler) WithGroup(name string) slog.Handler {
	return &TracingHandler{h.handler.WithGroup(name)}
}

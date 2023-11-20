package log

import (
	"context"
	"log/slog"
)

type ctxKey struct{}

func ContextWithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, logger)
}

func Ctx(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(ctxKey{}).(*slog.Logger); ok {
		return l
	}

	return slog.New(&noOpHandler{})
}

var _ slog.Handler = (*noOpHandler)(nil)

type noOpHandler struct{}

func (h *noOpHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return true
}

func (h *noOpHandler) Handle(ctx context.Context, r slog.Record) error {
	return nil
}

func (h *noOpHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &noOpHandler{}
}

func (h *noOpHandler) WithGroup(name string) slog.Handler {
	return &noOpHandler{}
}

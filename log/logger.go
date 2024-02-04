package log

import (
	"context"
	"io"
	"log/slog"

	"github.com/actatum/tk/errs"
	"go.opentelemetry.io/otel/trace"
)

type ctxKey int

const slogFields ctxKey = iota

// NewLogger returns a new instance of slog.Logger that handles passing values via context and outputting opentelemetry trace information.
func NewLogger(out io.Writer) *slog.Logger {
	handler := ContextHandler{
		Handler: slog.NewJSONHandler(out, &slog.HandlerOptions{
			AddSource:   true,
			ReplaceAttr: marshalError,
		}),
	}

	logger := slog.New(handler)

	slog.SetDefault(logger)

	return logger
}

type ContextHandler struct {
	slog.Handler
}

func (h ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if attrs, ok := ctx.Value(slogFields).([]slog.Attr); ok {
		for _, v := range attrs {
			r.AddAttrs(v)
		}
	}

	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.IsValid() {
		r.AddAttrs(
			slog.String("trace_id", spanCtx.TraceID().String()),
			slog.String("span_id", spanCtx.SpanID().String()),
		)
	}

	return h.Handler.Handle(ctx, r)
}

func AppendCtx(ctx context.Context, attrs ...slog.Attr) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	if v, ok := ctx.Value(slogFields).([]slog.Attr); ok {
		v = append(v, attrs...)
		return context.WithValue(ctx, slogFields, v)
	}

	v := []slog.Attr{}
	v = append(v, attrs...)
	return context.WithValue(ctx, slogFields, v)
}

func marshalError(_ []string, a slog.Attr) slog.Attr {
	switch a.Key {
	case "error":
		err, ok := a.Value.Any().(error)
		if ok {
			kind := errs.ErrorKind(err)
			message := errs.ErrorMessage(err)
			return slog.Group(
				"error",
				slog.String("message", message),
				slog.String("kind", kind.String()),
			)
		}
	}

	return a
}

package trace

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
)

const (
	TraceIDHeader = "X-Trace-Id"
)

// HTTPMiddleware parses incoming traceparent header and uses it to start a span.
// If no traceparent header exists the middleware will create a new root span.
// The span is then propagated to the request context.
func HTTPMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		sc := trace.SpanContextFromContext(r.Context())
		w.Header().Add(TraceIDHeader, sc.TraceID().String())

		next.ServeHTTP(w, r)
	}

	return otelhttp.NewHandler(http.HandlerFunc(fn), "", otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
		return r.URL.Path
	}))
}

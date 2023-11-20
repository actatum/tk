package log

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-chi/chi/v5/middleware"
)

// HTTPMiddleware logs incoming http requests.
func HTTPMiddleware(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			defer func(start time.Time) {
				status := ww.Status()

				l := logger.WithGroup("request").With(
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
					slog.String("proto", r.Proto),
					slog.String("from", r.RemoteAddr),
					slog.Int("status", status),
					slog.Int("size", ww.BytesWritten()),
					slog.String("duration", time.Since(start).String()),
				)

				switch {
				case status >= http.StatusInternalServerError:
					l.ErrorContext(r.Context(), "Internal Server Error")
				case status >= http.StatusBadRequest:
					l.WarnContext(r.Context(), "Client Error")
				case status >= http.StatusMultipleChoices:
					l.InfoContext(r.Context(), "Redirection")
				default:
					l.InfoContext(r.Context(), "success")
				}
			}(time.Now())

			next.ServeHTTP(ww, r)
		})
	}
}

// WatermillMiddleware logs incoming event handling.
func WatermillMiddleware(logger *slog.Logger) message.HandlerMiddleware {
	return func(next message.HandlerFunc) message.HandlerFunc {
		return func(msg *message.Message) ([]*message.Message, error) {
			start := time.Now()

			msgs, err := next(msg)
			l := logger.WithGroup("event").With(
				slog.String("id", msg.UUID),
				slog.String("duration", time.Since(start).String()),
			)
			if err != nil {
				l.ErrorContext(msg.Context(), "error handling event")
			} else {
				l.InfoContext(msg.Context(), "handled event")
			}

			return msgs, err
		}
	}
}

package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/actatum/tk/log"
	"github.com/actatum/tk/trace"
	"github.com/go-chi/chi/v5"
)

func main() {
	tp, err := trace.NewTraceProvider("test")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer tp.Shutdown(context.Background())

	ctx, _ := tp.Tracer("").Start(context.Background(), "test")

	logger := log.NewLogger(os.Stderr)

	logger.InfoContext(ctx, "test")

	slog.InfoContext(ctx, "default")

	r := chi.NewRouter()
	r.Use(log.HTTPMiddleware(logger))

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println(err)
	}
}

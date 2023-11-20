package trace

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

// NewTraceProvider returns a new trace provider with the name and function options applied.
func NewTraceProvider(service string, opts ...Option) (*sdktrace.TracerProvider, error) {
	options := options{}

	for _, opt := range opts {
		opt.apply(&options)
	}

	resourceAttrs := []attribute.KeyValue{
		semconv.ServiceNameKey.String(service),
	}

	if options.env != "" {
		resourceAttrs = append(resourceAttrs, attribute.String("environment", options.env))
	}

	res, err := resource.New(
		context.Background(),
		resource.WithTelemetrySDK(),
		resource.WithAttributes(resourceAttrs...),
	)
	if err != nil {
		return nil, fmt.Errorf("creating resource: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tp, nil
}

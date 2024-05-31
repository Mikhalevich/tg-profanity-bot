package tracing

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

var (
	serviceName string
)

func SetupTracer(
	endpoint string,
	serviceName string,
	serviceVersion string,
) error {
	exporter, err := otlptrace.New(
		context.Background(),
		otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint(endpoint),
			otlptracehttp.WithHeaders(map[string]string{
				"content-type": "application/json",
			}),
			otlptracehttp.WithInsecure(),
		),
	)

	if err != nil {
		return fmt.Errorf("creating exporter: %w", err)
	}

	tracerprovider := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceName(serviceName),
				semconv.ServiceVersion(serviceVersion),
			),
		),
	)

	otel.SetTracerProvider(tracerprovider)

	return nil
}

func TraceFn(ctx context.Context, name string, fn func(ctx context.Context) error) error {
	t := otel.Tracer(serviceName)

	ctx, span := t.Start(ctx, name)
	defer span.End()

	if err := fn(ctx); err != nil {
		return fmt.Errorf("trace fn: %w", err)
	}

	return nil
}

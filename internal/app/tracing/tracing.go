package tracing

import (
	"context"
	"fmt"
	"runtime"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	serviceName string
)

func SetupTracer(
	endpoint string,
	name string,
	version string,
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

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(name),
			semconv.ServiceVersion(version),
		),
	)

	if err != nil {
		return fmt.Errorf("merge resource: %w", err)
	}

	tracerprovider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	serviceName = name

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

func Trace(ctx context.Context) (context.Context, trace.Span) {
	t := otel.Tracer(serviceName)
	return t.Start(ctx, funcName())
}

func funcName() string {
	pc, _, _, ok := runtime.Caller(2)
	if !ok {
		return ""
	}

	return runtime.FuncForPC(pc).Name()
}

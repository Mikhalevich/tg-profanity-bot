package tracing

import (
	"context"
	"fmt"
	"runtime"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	std Tracer = NewNoopTracer()
)

type Tracer interface {
	StartSpan(ctx context.Context) (context.Context, trace.Span)
}

type OtelTracer struct {
	Name string
}

func NewOtelTracer(name string) *OtelTracer {
	return &OtelTracer{
		Name: name,
	}
}

func (t *OtelTracer) StartSpan(ctx context.Context) (context.Context, trace.Span) {
	tr := otel.Tracer(t.Name)
	//nolint:spancheck
	return tr.Start(ctx, callingFuncName())
}

func callingFuncName() string {
	pc, _, _, ok := runtime.Caller(2)
	if !ok {
		return ""
	}

	return runtime.FuncForPC(pc).Name()
}

func StartSpan(ctx context.Context) (context.Context, trace.Span) {
	return std.StartSpan(ctx)
}

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

	std = NewOtelTracer(name)

	otel.SetTracerProvider(tracerprovider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return nil
}
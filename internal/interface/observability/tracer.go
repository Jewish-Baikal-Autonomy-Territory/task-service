package observability

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.41.0"
)

type TracerOptions struct {
	ServiceName      string
	ServiceVersion   string
	Env              string
	CollectorAddress string
	TracerTimeout    time.Duration
}

var samplerStrategies = map[string]func() sdktrace.Sampler{
	"dev": sdktrace.AlwaysSample,
	"prod": func() sdktrace.Sampler {
		return sdktrace.ParentBased(sdktrace.TraceIDRatioBased(0.1))
	},
}

func InitTracerProvider(ctx context.Context, opts TracerOptions) (*sdktrace.TracerProvider, error) {
	exporter, err := otlptrace.New(
		ctx,
		otlptracegrpc.NewClient(
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint(opts.CollectorAddress),
			otlptracegrpc.WithTimeout(opts.TracerTimeout),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("otlptrace.New: %w", err)
	}

	resources, err := resource.New(
		ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(opts.ServiceName),
			semconv.ServiceVersionKey.String(opts.ServiceVersion),
			semconv.DeploymentEnvironmentNameKey.String(opts.Env),
			semconv.TelemetrySDKLanguageGo,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("resource.New: %w", err)
	}

	samplerProducer, ok := samplerStrategies[opts.Env]
	if !ok {
		return nil, fmt.Errorf("samplerStrategies not supported by env: %s", opts.Env)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(samplerProducer()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resources),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	),
	)

	return tp, nil
}

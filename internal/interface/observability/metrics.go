package observability

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"

	semconv "go.opentelemetry.io/otel/semconv/v1.41.0"
)

type MetricsProviderOptions struct {
	ServiceName                 string
	ServiceVersion              string
	Env                         string
	MetricsExporterAddress      string
	MetricsExporterTimeout      time.Duration
	MeterPeriodicReaderInterval time.Duration
	MeterPeriodicReaderTimeout  time.Duration
}

func InitMetricsProvider(ctx context.Context, opts MetricsProviderOptions) (*metric.MeterProvider, error) {
	metricExporter, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(opts.MetricsExporterAddress),
		otlpmetricgrpc.WithTimeout(opts.MetricsExporterTimeout),
	)
	if err != nil {
		return nil, fmt.Errorf("otlpmetricgrpc.New: %w", err)
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

	meterProvider := metric.NewMeterProvider(
		metric.WithReader(
			metric.NewPeriodicReader(
				metricExporter,
				metric.WithInterval(opts.MeterPeriodicReaderInterval),
				metric.WithTimeout(opts.MeterPeriodicReaderTimeout),
			),
		),
		metric.WithResource(resources),
	)

	return meterProvider, nil
}

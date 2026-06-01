package observability

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.41.0"
)

type LoggerProviderOptions struct {
	ServiceName           string
	ServiceVersion        string
	Env                   string
	LoggerExporterAddress string
	LoggerExporterTimeout time.Duration
}

func InitLoggerProvider(ctx context.Context, opts LoggerProviderOptions) (*sdklog.LoggerProvider, error) {
	loggerExporter, err := otlploggrpc.New(
		ctx,
		otlploggrpc.WithInsecure(),
		otlploggrpc.WithEndpoint(opts.LoggerExporterAddress),
		otlploggrpc.WithTimeout(opts.LoggerExporterTimeout),
	)
	if err != nil {
		return nil, fmt.Errorf("otlploggrpc.New: %w", err)
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

	processor := sdklog.NewBatchProcessor(loggerExporter)

	loggerProvider := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(processor),
		sdklog.WithResource(resources),
	)

	global.SetLoggerProvider(loggerProvider)

	return loggerProvider, nil
}

package interceptor

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type MetricsInterceptor struct {
	rpcDuration metric.Float64Histogram
	rpcTotal    metric.Int64Counter
}

func (i *MetricsInterceptor) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		elapsed := time.Since(start).Seconds()

		st, _ := status.FromError(err)
		codeStr := st.Code().String()

		attrs := metric.WithAttributes(
			attribute.String("rpc.method", info.FullMethod),
			attribute.String("rpc.status_code", codeStr),
		)

		i.rpcTotal.Add(ctx, 1, attrs)
		i.rpcDuration.Record(ctx, elapsed, attrs)

		return resp, err
	}
}

func NewMetricsInterceptor(meter metric.Meter) (*MetricsInterceptor, error) {
	duration, err := meter.Float64Histogram(
		"rpc.server.duration",
		metric.WithDescription("Measures the total latency of gRPC unary invocations."),
		metric.WithUnit("s"),
	)

	if err != nil {
		return nil, err
	}

	total, err := meter.Int64Counter(
		"rpc.server.request_count",
		metric.WithDescription("Total number of handled RPC requests."),
	)

	if err != nil {
		return nil, err
	}

	return &MetricsInterceptor{
		rpcDuration: duration,
		rpcTotal:    total,
	}, nil
}

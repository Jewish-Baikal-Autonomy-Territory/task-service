package interceptor

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

func NewLoggingInterceptor(logger *zerolog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start)

		if err != nil {
			logger.Error().
				Ctx(ctx).
				Str("rpc.method", info.FullMethod).
				Dur("duration", duration).
				Err(err).
				Msg("gRPC failure")
		} else {
			logger.Info().
				Ctx(ctx).
				Str("rpc.method", info.FullMethod).
				Dur("duration", duration).
				Msg("gRPC success")
		}

		return resp, err
	}
}

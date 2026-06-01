package interceptor

import (
	"context"
	"runtime/debug"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewRecoverInterceptor(logger *zerolog.Logger) grpc.UnaryServerInterceptor {
	recoveryFunc := func(ctx context.Context, req any) error {
		stackTrace := debug.Stack()

		logger.Error().
			Any("panic_value", req).
			Str("stack_trace", string(stackTrace)).
			Msg("gRPC server panicked and recovered")

		return status.Error(codes.Internal, "internal server error")
	}

	opts := []recovery.Option{
		recovery.WithRecoveryHandlerContext(recoveryFunc),
	}

	return recovery.UnaryServerInterceptor(opts...)
}

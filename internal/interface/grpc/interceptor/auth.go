package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func NewAuthInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Internal, "internal server error")
		}

		ownerID := md.Get("api-user-id")
		if len(ownerID) == 0 || ownerID[0] == "" {
			return nil, status.Error(codes.Internal, "internal server error")
		}

		valueCtx := context.WithValue(ctx, "requester-id", ownerID[0])

		return handler(valueCtx, req)
	}
}

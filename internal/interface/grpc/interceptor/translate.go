package interceptor

import (
	"context"
	"errors"
	apptask "task-service/internal/application/task"
	domaintask "task-service/internal/domain/task"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func translateError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, domaintask.ErrFailedPrecondition) {
		return status.Error(codes.FailedPrecondition, "failed precondition")
	}

	if errors.Is(err, domaintask.ErrInvalidData) {
		return status.Error(codes.InvalidArgument, "invalid data")
	}

	if errors.Is(err, domaintask.ErrNotFound) {
		return status.Error(codes.NotFound, "not found")
	}

	if errors.Is(err, domaintask.ErrAlreadyDeleted) {
		return status.Error(codes.FailedPrecondition, "already deleted")
	}

	if errors.Is(err, domaintask.ErrNotDeleted) {
		return status.Error(codes.FailedPrecondition, "not deleted")
	}

	if errors.Is(err, apptask.ErrUnsupported) {
		return status.Error(codes.Unimplemented, "unsupported")
	}

	if errors.Is(err, apptask.ErrPermissionDenied) {
		return status.Error(codes.PermissionDenied, "permission denied")
	}

	return status.Error(codes.Internal, "internal error")
}

func NewErrTranslateInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		resp, err := handler(ctx, req)
		return resp, translateError(err)
	}
}

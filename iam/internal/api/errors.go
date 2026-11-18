package api

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ivmerk/goshop/iam/internal/model"
)

// ToGRPCError переводит доменные ошибки в gRPC-статусы.
// Внутренние ошибки наружу не раскрываются.
func ToGRPCError(err error) error {
	switch {
	case errors.Is(err, model.ErrInvalidArgument):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, model.ErrUserAlreadyExists):
		return status.Error(codes.AlreadyExists, "user already exists")
	case errors.Is(err, model.ErrUserNotFound):
		return status.Error(codes.NotFound, "user not found")
	case errors.Is(err, model.ErrInvalidCredentials):
		return status.Error(codes.Unauthenticated, "invalid login or password")
	case errors.Is(err, model.ErrSessionNotFound):
		return status.Error(codes.Unauthenticated, "session not found or expired")
	default:
		return status.Error(codes.Internal, "internal error")
	}
}

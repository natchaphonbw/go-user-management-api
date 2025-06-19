package usecases

import (
	"context"

	"github.com/google/uuid"

	"github.com/natchaphonbw/usermanagement/modules/users/dtos"
	app_errors "github.com/natchaphonbw/usermanagement/pkg/errors"
)

type AuthUsecase interface {
	RegisterUser(ctx context.Context, input dtos.RegisterRequest) (*dtos.UserResponse, *app_errors.AppError)
	Login(ctx context.Context, input dtos.LoginRequest) (*dtos.LoginResponse, *app_errors.AppError)
	GetProfile(ctx context.Context, userID uuid.UUID) (*dtos.UserResponse, *app_errors.AppError)
}

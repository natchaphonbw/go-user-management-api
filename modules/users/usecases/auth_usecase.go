package usecases

import (
	"context"

	"github.com/google/uuid"

	"github.com/natchaphonbw/usermanagement/modules/users/dtos"
)

type AuthUsecase interface {
	RegisterUser(ctx context.Context, input dtos.RegisterRequest) (*dtos.UserResponse, error)
	Login(ctx context.Context, input dtos.LoginRequest) (*dtos.LoginResponse, error)
	GetProfile(ctx context.Context, userID uuid.UUID) (*dtos.UserResponse, error)
}

package usecases

import (
	"context"

	"github.com/google/uuid"

	"github.com/natchaphonbw/usermanagement/modules/users/dtos"
	app_errors "github.com/natchaphonbw/usermanagement/pkg/errors"
)

type UserUsecase interface {
	CreateUser(ctx context.Context, input dtos.CreateUserRequest) (*dtos.UserResponse, *app_errors.AppError)
	GetAllUsers(ctx context.Context) ([]*dtos.UserResponse, *app_errors.AppError)
	GetUserByID(ctx context.Context, id uuid.UUID) (*dtos.UserResponse, *app_errors.AppError)
	UpdateUserByID(ctx context.Context, id uuid.UUID, input dtos.UpdateUserRequest) (*dtos.UserResponse, *app_errors.AppError)
	DeleteUserByID(ctx context.Context, id uuid.UUID) (*dtos.UserResponse, *app_errors.AppError)
}

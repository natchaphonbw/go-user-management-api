package usecases

import (
	"context"

	"github.com/google/uuid"

	"github.com/natchaphonbw/usermanagement/modules/users/dtos"
)

type UserUsecase interface {
	CreateUser(ctx context.Context, input dtos.CreateUserRequest) (*dtos.UserResponse, error)
	GetAllUsers(ctx context.Context) ([]*dtos.UserResponse, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*dtos.UserResponse, error)
	UpdateUserByID(ctx context.Context, id uuid.UUID, input dtos.UpdateUserRequest) (*dtos.UserResponse, error)
	DeleteUserByID(ctx context.Context, id uuid.UUID) (*dtos.UserResponse, error)
}

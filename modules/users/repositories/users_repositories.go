package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/natchaphonbw/usermanagement/modules/users/entities"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *entities.User) error
	GetAllUsers(ctx context.Context) ([]entities.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	UpdateUserByID(ctx context.Context, id uuid.UUID, user *entities.User) (*entities.User, error)
	DeleteUserByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
}

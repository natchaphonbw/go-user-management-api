package usecases

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/natchaphonbw/usermanagement/modules/users/dtos"
	"github.com/natchaphonbw/usermanagement/modules/users/entities"
	"github.com/natchaphonbw/usermanagement/modules/users/repositories"
	"github.com/natchaphonbw/usermanagement/pkg/utils"
)

type userUsecaseImpl struct {
	repo repositories.UserRepository
}

func NewUserUseCase(repo repositories.UserRepository) UserUsecase {
	return &userUsecaseImpl{
		repo: repo,
	}
}

// Create User
func (u *userUsecaseImpl) CreateUser(ctx context.Context, input dtos.CreateUserRequest) (*dtos.UserResponse, error) {
	hash, salt, err := utils.GeneratePasswordHash(input.Password, &utils.DefaultArgon2Config)
	if err != nil {
		return nil, err
	}

	user := &entities.User{
		ID:           uuid.New(),
		Name:         input.Name,
		Email:        input.Email,
		Age:          input.Age,
		PasswordHash: hash,
		Salt:         salt,
		Created_at:   time.Now(),
		Updated_at:   time.Now(),
	}

	if err := u.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return dtos.FromUserEntity(user), nil
}

// Get All Users
func (u *userUsecaseImpl) GetAllUsers(ctx context.Context) ([]*dtos.UserResponse, error) {
	users, err := u.repo.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	return dtos.FromUserEntities(users), nil
}

// Get User By ID
func (u *userUsecaseImpl) GetUserByID(ctx context.Context, id uuid.UUID) (*dtos.UserResponse, error) {
	user, err := u.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return dtos.FromUserEntity(user), nil

}

// Update User By ID
func (u *userUsecaseImpl) UpdateUserByID(ctx context.Context, id uuid.UUID, input dtos.UpdateUserRequest) (*dtos.UserResponse, error) {

	user, err := u.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	updated := false
	// Update user fields
	if input.Name != nil && user.Name != *input.Name {
		user.Name = *input.Name
		updated = true
	}

	if input.Email != nil && user.Email != *input.Email {
		user.Email = *input.Email
		updated = true
	}

	if input.Age != nil && user.Age != *input.Age {
		user.Age = *input.Age
		updated = true
	}

	if !updated {
		return dtos.FromUserEntity(user), nil
	}

	user.Updated_at = time.Now()

	// Update user in repository
	user, err = u.repo.UpdateUserByID(ctx, id, user)
	if err != nil {
		return nil, err
	}

	return dtos.FromUserEntity(user), nil

}

// Delete User By ID
func (u *userUsecaseImpl) DeleteUserByID(ctx context.Context, id uuid.UUID) (*dtos.UserResponse, error) {
	user, err := u.repo.DeleteUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return dtos.FromUserEntity(user), nil

}

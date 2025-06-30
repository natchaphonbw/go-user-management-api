package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/natchaphonbw/usermanagement/modules/users/dtos"
	"github.com/natchaphonbw/usermanagement/modules/users/entities"
	"github.com/natchaphonbw/usermanagement/modules/users/repositories"
	app_errors "github.com/natchaphonbw/usermanagement/pkg/errors"
	"github.com/natchaphonbw/usermanagement/pkg/utils"
)

type userUsecaseImpl struct {
	userRepo repositories.UserRepository
}

func NewUserUseCase(userRepo repositories.UserRepository) UserUsecase {
	return &userUsecaseImpl{
		userRepo: userRepo,
	}
}

// Create User
func (u *userUsecaseImpl) CreateUser(ctx context.Context, input dtos.CreateUserRequest) (*dtos.UserResponse, *app_errors.AppError) {
	hash, salt, err := utils.GeneratePasswordHash(input.Password, &utils.DefaultArgon2Config)
	if err != nil {
		return nil, app_errors.InternalServer("Failed to hash password", err)
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

	if err := u.userRepo.CreateUser(ctx, user); err != nil {
		return nil, app_errors.InternalServer("Failed to create user", err)
	}

	return dtos.FromUserEntity(user), nil
}

// Get All Users
func (u *userUsecaseImpl) GetAllUsers(ctx context.Context) ([]*dtos.UserResponse, *app_errors.AppError) {
	users, err := u.userRepo.GetAllUsers(ctx)
	if err != nil {
		return nil, app_errors.InternalServer("Failed to get users", err)
	}

	return dtos.FromUserEntities(users), nil
}

// Get User By ID
func (u *userUsecaseImpl) GetUserByID(ctx context.Context, id uuid.UUID) (*dtos.UserResponse, *app_errors.AppError) {
	user, err := u.userRepo.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, app_errors.NotFound("User not found", err)
		}
		return nil, app_errors.InternalServer("Failed to get user", err)
	}
	return dtos.FromUserEntity(user), nil
}

// Update User By ID
func (u *userUsecaseImpl) UpdateUserByID(ctx context.Context, id uuid.UUID, input dtos.UpdateUserRequest) (*dtos.UserResponse, *app_errors.AppError) {

	user, err := u.userRepo.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, app_errors.NotFound("User not found", err)
		}
		return nil, app_errors.InternalServer("Failed to get user for update", err)
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

	// Update user in userRepository
	user, err = u.userRepo.UpdateUserByID(ctx, id, user)
	if err != nil {
		return nil, app_errors.InternalServer("Failed to update user", err)
	}

	return dtos.FromUserEntity(user), nil

}

// Delete User By ID
func (u *userUsecaseImpl) DeleteUserByID(ctx context.Context, id uuid.UUID) (*dtos.UserResponse, *app_errors.AppError) {
	user, err := u.userRepo.DeleteUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, app_errors.NotFound("User not found", err)
		}
		return nil, app_errors.InternalServer("Failed to delete user", err)
	}

	return dtos.FromUserEntity(user), nil

}

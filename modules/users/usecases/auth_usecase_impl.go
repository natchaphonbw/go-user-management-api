package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/natchaphonbw/usermanagement/modules/users/dtos"
	"github.com/natchaphonbw/usermanagement/modules/users/repositories"
	"github.com/natchaphonbw/usermanagement/modules/users/validator"
	app_errors "github.com/natchaphonbw/usermanagement/pkg/errors"
	"github.com/natchaphonbw/usermanagement/pkg/jwt"
	"github.com/natchaphonbw/usermanagement/pkg/utils"
	"gorm.io/gorm"
)

type AuthUsecaseImpl struct {
	userUsecase UserUsecase
	repo        repositories.UserRepository
}

func NewAuthUseCase(userUsecase UserUsecase, repo repositories.UserRepository) AuthUsecase {
	return &AuthUsecaseImpl{
		userUsecase: userUsecase,
		repo:        repo,
	}
}

// register
func (a *AuthUsecaseImpl) RegisterUser(ctx context.Context, input dtos.RegisterRequest) (*dtos.UserResponse, *app_errors.AppError) {
	// validate pwd
	if err := validator.ValidatePassword(input.Password); err != nil {
		return nil, app_errors.BadRequest("Invalid password", err)
	}

	createReq := dtos.CreateUserRequest{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
		Age:      input.Age,
	}

	userResp, err := a.userUsecase.CreateUser(ctx, createReq)
	if err != nil {
		return nil, err
	}

	return userResp, nil
}

// login
func (a *AuthUsecaseImpl) Login(ctx context.Context, req dtos.LoginRequest) (*dtos.LoginResponse, *app_errors.AppError) {
	user, err := a.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, app_errors.Unautherized("Invalid credentials", err)
		}
		return nil, app_errors.InternalServer("Failed to get user", err)
	}

	// check pwd
	match, err := utils.VerifyPassword(req.Password, user.PasswordHash, user.Salt, &utils.DefaultArgon2Config)
	if err != nil || !match {
		return nil, app_errors.Unautherized("Invalid credentials", fmt.Errorf("password mismatch"))
	}

	// gen jwt token
	token, err := jwt.GenerateToken(user.ID.String())
	if err != nil {
		return nil, app_errors.InternalServer("Failed to generate token", err)
	}

	return &dtos.LoginResponse{Token: token}, nil

}

// Get Profile
func (a *AuthUsecaseImpl) GetProfile(ctx context.Context, userID uuid.UUID) (*dtos.UserResponse, *app_errors.AppError) {
	user, err := a.repo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, app_errors.NotFound("User not found", err)
		}
		return nil, app_errors.InternalServer("Failed to get user", err)
	}

	return dtos.FromUserEntity(user), nil
}

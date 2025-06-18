package usecases

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/natchaphonbw/usermanagement/modules/users/dtos"
	"github.com/natchaphonbw/usermanagement/modules/users/repositories"
	"github.com/natchaphonbw/usermanagement/modules/users/validations"
	"github.com/natchaphonbw/usermanagement/pkg/jwt"
	"github.com/natchaphonbw/usermanagement/pkg/utils"
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
func (a *AuthUsecaseImpl) RegisterUser(ctx context.Context, input dtos.RegisterRequest) (*dtos.UserResponse, error) {
	// validate pwd
	if err := validations.ValidatePassword(input.Password); err != nil {
		return nil, err
	}

	createReq := dtos.CreateUserRequest{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
		Age:      input.Age,
	}

	return a.userUsecase.CreateUser(ctx, createReq)
}

// login
func (a *AuthUsecaseImpl) Login(ctx context.Context, req dtos.LoginRequest) (*dtos.LoginResponse, error) {
	user, err := a.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// check pwd
	match, err := utils.VerifyPassword(req.Password, user.PasswordHash, user.Salt, &utils.DefaultArgon2Config)
	if err != nil || !match {
		return nil, fmt.Errorf("invalid credentials")
	}

	// gen jwt token
	token, err := jwt.GenerateToken(user.ID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %v", err)
	}

	return &dtos.LoginResponse{Token: token}, nil

}

// Get Profile
func (a *AuthUsecaseImpl) GetProfile(ctx context.Context, userID uuid.UUID) (*dtos.UserResponse, error) {
	user, err := a.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return dtos.FromUserEntity(user), nil
}

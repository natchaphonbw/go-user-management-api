package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/natchaphonbw/usermanagement/pkg/utils"
	"gorm.io/gorm"

	"github.com/natchaphonbw/usermanagement/modules/users/dtos"
	"github.com/natchaphonbw/usermanagement/modules/users/repositories"
	"github.com/natchaphonbw/usermanagement/modules/users/validator"
	app_errors "github.com/natchaphonbw/usermanagement/pkg/errors"
	"github.com/natchaphonbw/usermanagement/pkg/jwt"
)

type AuthUsecaseImpl struct {
	userUsecase    UserUsecase
	refreshUsecase RefreshTokenUsecase

	userRepo    repositories.UserRepository
	refreshRepo repositories.RefreshTokenRepository
}

func NewAuthUseCase(userUsecase UserUsecase, refreshUsecase RefreshTokenUsecase, userRepo repositories.UserRepository, refreshRepo repositories.RefreshTokenRepository) AuthUsecase {
	return &AuthUsecaseImpl{
		userUsecase:    userUsecase,
		refreshUsecase: refreshUsecase,

		userRepo:    userRepo,
		refreshRepo: refreshRepo,
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
func (a *AuthUsecaseImpl) Login(ctx context.Context, req dtos.LoginRequest, deviceIP, deviceUA, deviceID string) (*dtos.LoginResponse, *app_errors.AppError) {
	user, err := a.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, app_errors.Unautherized("Invalid credentials", err)
		}
		return nil, app_errors.InternalServer("Failed to get user", err)
	}

	// verify pwd
	match, err := utils.VerifyPassword(req.Password, user.PasswordHash, user.Salt, &utils.DefaultArgon2Config)
	if err != nil || !match {
		return nil, app_errors.Unautherized("Invalid credentials", fmt.Errorf("password mismatch"))
	}

	// gen jwt token
	tokenPair, pairErr := a.refreshUsecase.IssueTokenPair(ctx, user.ID, deviceIP, deviceUA, deviceID)
	if pairErr != nil {
		return nil, app_errors.InternalServer("Failed to issue token pair", pairErr).WithDetails(pairErr.Details)
	}
	return &dtos.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}

// Log out
func (a *AuthUsecaseImpl) Logout(ctx context.Context, userID uuid.UUID, refreshToken, deviceID, deviceUA string) *app_errors.AppError {
	claims, verifyErr := jwt.VerifyRefreshToken(refreshToken)
	if verifyErr != nil {
		return app_errors.Unautherized("Invalid refresh token", verifyErr)
	}

	if claims.UserID != userID.String() {
		return app_errors.Unautherized("Refresh token does not belong to user", nil)
	}

	// Load token record from DB
	tokenRecord, err := a.refreshRepo.GetByToken(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return app_errors.NotFound("Refresh token not found", err)
		}
		return app_errors.InternalServer("Failed to retrieve refresh token", err)
	}

	// Check device info
	if tokenRecord.DeviceID != deviceID || tokenRecord.DeviceUA != deviceUA {
		return app_errors.Unautherized("Device info mismatch", nil)
	}

	if err := a.refreshRepo.DeleteByToken(ctx, refreshToken); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return app_errors.NotFound("Refresh token not found", err)
		}
		return app_errors.InternalServer("Failed to revoke refresh token", err)
	}
	return nil
}

// Log out all devices
func (a *AuthUsecaseImpl) LogoutAll(ctx context.Context, userID uuid.UUID) *app_errors.AppError {
	if err := a.refreshRepo.DeleteByUserID(ctx, userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return app_errors.NotFound("No refresh tokens found for user", err)
		}
		return app_errors.InternalServer("Failed to revoke all refresh tokens", err)
	}
	return nil
}

// Get Profile
func (a *AuthUsecaseImpl) GetProfile(ctx context.Context, userID uuid.UUID) (*dtos.UserResponse, *app_errors.AppError) {
	user, err := a.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, app_errors.NotFound("User not found", err)
		}
		return nil, app_errors.InternalServer("Failed to get user", err)
	}

	return dtos.FromUserEntity(user), nil
}

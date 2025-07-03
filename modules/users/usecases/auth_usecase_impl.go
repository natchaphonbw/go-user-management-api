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
)

type AuthUsecaseImpl struct {
	userUsecase    UserUsecase
	sessionUsecase SessionUsecase

	userRepo    repositories.UserRepository
	sessionRepo repositories.SessionRepository
}

func NewAuthUseCase(userUsecase UserUsecase, sessionUsecase SessionUsecase, userRepo repositories.UserRepository, sessionRepo repositories.SessionRepository) AuthUsecase {
	return &AuthUsecaseImpl{
		userUsecase:    userUsecase,
		sessionUsecase: sessionUsecase,

		userRepo:    userRepo,
		sessionRepo: sessionRepo,
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
	tokenPair, pairErr := a.sessionUsecase.IssueTokenPair(ctx, user.ID, deviceIP, deviceUA, deviceID)
	if pairErr != nil {
		return nil, app_errors.InternalServer("Failed to issue token pair", pairErr).WithDetails(pairErr.Details)
	}
	return &dtos.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}

// Log out
func (a *AuthUsecaseImpl) Logout(ctx context.Context, sessionID uuid.UUID, deviceID, deviceUA string) *app_errors.AppError {

	// Load token record from DB
	session, err := a.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return app_errors.NotFound("Session not found", err)
		}
		return app_errors.InternalServer("Failed to get session", err)
	}

	// check revoked
	if session.Revoked {
		return app_errors.Unautherized("Refresh token already revoked", nil)
	}

	// check device info
	if session.DeviceID != deviceID || session.DeviceUA != deviceUA {
		return app_errors.Unautherized("Device info mismatch", nil)
	}

	// revoke
	if revokeErr := a.sessionRepo.MarkRevoked(ctx, sessionID); revokeErr != nil {
		if errors.Is(revokeErr, gorm.ErrRecordNotFound) {
			return app_errors.NotFound("session not found", revokeErr)
		}
		return app_errors.InternalServer("Failed to revoke sessions", revokeErr)
	}

	return nil
}

// Log out all devices
func (a *AuthUsecaseImpl) LogoutAll(ctx context.Context, userID uuid.UUID) *app_errors.AppError {

	// revoke
	if revokeErr := a.sessionRepo.MarkRevokedByUserID(ctx, userID); revokeErr != nil {
		if errors.Is(revokeErr, gorm.ErrRecordNotFound) {
			return app_errors.NotFound("session not found", revokeErr)
		}
		return app_errors.InternalServer("Failed to revoke sessions for user", revokeErr)
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

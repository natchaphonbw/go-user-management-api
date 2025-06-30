package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/natchaphonbw/usermanagement/modules/users/dtos"
	"github.com/natchaphonbw/usermanagement/modules/users/entities"
	"github.com/natchaphonbw/usermanagement/modules/users/repositories"
	"github.com/natchaphonbw/usermanagement/modules/users/validator"
	app_errors "github.com/natchaphonbw/usermanagement/pkg/errors"
	"github.com/natchaphonbw/usermanagement/pkg/jwt"
	"gorm.io/gorm"
)

type RefreshTokenUsecaseImpl struct {
	repo repositories.RefreshTokenRepository
}

func NewRefreshTokenUsecase(repo repositories.RefreshTokenRepository) RefreshTokenUsecase {
	return &RefreshTokenUsecaseImpl{
		repo: repo,
	}
}

func (u *RefreshTokenUsecaseImpl) IssueTokenPair(ctx context.Context, userID uuid.UUID, deviceIP, deviceUA, deviceID string) (*dtos.TokenPair, *app_errors.AppError) {
	accessToken, err := jwt.GenerateAccessToken(userID.String())
	if err != nil {
		return nil, app_errors.InternalServer("Failed to generate token", err)
	}

	refreshToken, issuedAt, expiresAt, err := jwt.GenerateRefreshToken(userID.String())
	if err != nil {
		return nil, app_errors.InternalServer("Failed to generate token", err)
	}

	// save refresh token
	refresh := &entities.RefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     refreshToken,
		DeviceIP:  deviceIP,
		DeviceUA:  deviceUA,
		DeviceID:  deviceID,
		IssuedAt:  issuedAt,
		ExpiresAt: expiresAt,
	}

	// validate refresh token
	if errs := validator.ValidateStruct(refresh); errs != nil {
		return nil, app_errors.BadRequest("Invalid refresh token data", nil).WithDetails(errs)
	}
	if err := u.repo.Upsert(ctx, refresh); err != nil {
		return nil, app_errors.InternalServer("Failed to save refresh token", err)
	}

	return &dtos.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (u *RefreshTokenUsecaseImpl) Refresh(ctx context.Context, refreshToken, deviceIP, deviceUA, deviceID string) (*dtos.TokenPair, *app_errors.AppError) {
	_, err := jwt.VerifyRefreshToken(refreshToken)
	if err != nil {
		return nil, app_errors.Unautherized("Invalid refresh token", err)
	}

	// find in db
	tokenRecord, err := u.repo.GetByToken(ctx, refreshToken)
	if err != nil {
		return nil, app_errors.Unautherized("Refresh token not found", err)
	}

	if time.Now().After(tokenRecord.ExpiresAt) {
		return nil, app_errors.Unautherized("Refresh token expired", nil)
	}

	// check device info
	if tokenRecord.DeviceID != deviceID || tokenRecord.DeviceUA != deviceUA {
		return nil, app_errors.Unautherized("Device info mismatch", nil)
	}

	// revoke old
	if revokeErr := u.repo.DeleteByToken(ctx, refreshToken); revokeErr != nil {
		if errors.Is(revokeErr, gorm.ErrRecordNotFound) {
			return nil, app_errors.NotFound("Refresh token not found", revokeErr)
		}
		return nil, app_errors.InternalServer("Failed to revoke old refresh token", revokeErr)
	}

	// gen new
	return u.IssueTokenPair(ctx, tokenRecord.UserID, deviceIP, deviceUA, deviceID)

}

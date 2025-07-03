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

type SessionUsecaseImpl struct {
	repo repositories.SessionRepository
}

func NewSessionUsecase(repo repositories.SessionRepository) SessionUsecase {
	return &SessionUsecaseImpl{
		repo: repo,
	}
}

func (u *SessionUsecaseImpl) IssueTokenPair(ctx context.Context, userID uuid.UUID, deviceIP, deviceUA, deviceID string) (*dtos.TokenPair, *app_errors.AppError) {
	// gen sessionID
	sessionID := uuid.New()

	// gen tokens
	accessToken, err := jwt.GenerateAccessToken(userID.String(), sessionID.String())
	if err != nil {
		return nil, app_errors.InternalServer("Failed to generate token", err)
	}
	// gen refresh
	refreshToken, issuedAt, expiresAt, err := jwt.GenerateRefreshToken(userID.String(), sessionID.String())
	if err != nil {
		return nil, app_errors.InternalServer("Failed to generate token", err)
	}

	// hash refresh token
	hashedToken, hashErr := jwt.HashRefreshToken(refreshToken)
	if hashErr != nil {
		return nil, app_errors.InternalServer("Failed to hash refresh token", hashErr)
	}

	session := &entities.Session{
		ID:          sessionID,
		UserID:      userID,
		HashedToken: hashedToken,
		DeviceIP:    deviceIP,
		DeviceID:    deviceID,
		DeviceUA:    deviceUA,
		IssuedAt:    issuedAt,
		ExpiresAt:   expiresAt,
		Revoked:     false,
	}

	// validate
	if errs := validator.ValidateStruct(session); errs != nil {
		return nil, app_errors.BadRequest("Invalid session data", nil).WithDetails(errs)
	}
	// save session
	if err := u.repo.Insert(ctx, session); err != nil {
		return nil, app_errors.InternalServer("Failed to save refresh token", err)
	}

	return &dtos.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (u *SessionUsecaseImpl) Refresh(ctx context.Context, refreshToken, deviceIP, deviceUA, deviceID string, sessionID uuid.UUID) (*dtos.TokenPair, *app_errors.AppError) {
	// get from db
	session, err := u.repo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, app_errors.Unautherized("Refresh token not found", err)
	}

	// check revoked or expired
	if session.Revoked {
		return nil, app_errors.Unautherized("Refresh token has been revoked", nil)
	}
	if time.Now().After(session.ExpiresAt) {
		return nil, app_errors.Unautherized("Refresh token expired", nil)
	}

	// compare
	match, verifyErr := jwt.VerifyRefreshTokenHash(refreshToken, session.HashedToken)
	if verifyErr != nil || !match {
		return nil, app_errors.Unautherized("Refresh token hash mismatch", verifyErr)
	}

	// check device info
	if session.DeviceID != deviceID || session.DeviceUA != deviceUA {
		return nil, app_errors.Unautherized("Device info mismatch", nil)
	}

	// revoke old
	if revokeErr := u.repo.MarkRevoked(ctx, sessionID); revokeErr != nil {
		if errors.Is(revokeErr, gorm.ErrRecordNotFound) {
			return nil, app_errors.NotFound("Refresh token not found", revokeErr)
		}
		return nil, app_errors.InternalServer("Failed to revoke old refresh token", revokeErr)
	}

	// gen new
	return u.IssueTokenPair(ctx, session.UserID, deviceIP, deviceUA, deviceID)

}

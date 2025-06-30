package usecases

import (
	"context"

	"github.com/google/uuid"

	"github.com/natchaphonbw/usermanagement/modules/users/dtos"
	app_errors "github.com/natchaphonbw/usermanagement/pkg/errors"
)

type RefreshTokenUsecase interface {
	IssueTokenPair(ctx context.Context, userID uuid.UUID, deviceIP, deviceUA, deviceID string) (*dtos.TokenPair, *app_errors.AppError)
	Refresh(ctx context.Context, refreshToken, deviceIP, deviceUA, deviceID string) (*dtos.TokenPair, *app_errors.AppError)
	// Revoke(ctx context.Context, refreshToken string) *app_errors.AppError
	// RevokeAll(ctx context.Context, userID uuid.UUID) *app_errors.AppError
}

package usecases

import (
	"context"

	"github.com/google/uuid"

	"github.com/natchaphonbw/usermanagement/modules/users/dtos"
	app_errors "github.com/natchaphonbw/usermanagement/pkg/errors"
)

type SessionUsecase interface {
	IssueTokenPair(ctx context.Context, userID uuid.UUID, deviceIP, deviceUA, deviceID string) (*dtos.TokenPair, *app_errors.AppError)
	Refresh(ctx context.Context, refreshToken, deviceIP, deviceUA, deviceID string, sessionID uuid.UUID) (*dtos.TokenPair, *app_errors.AppError)	
}

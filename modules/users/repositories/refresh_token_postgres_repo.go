package repositories

import (
	"context"

	"github.com/google/uuid"

	"github.com/natchaphonbw/usermanagement/modules/users/entities"
)

type RefreshTokenRepository interface {
	Upsert(ctx context.Context, token *entities.RefreshToken) error

	GetByToken(ctx context.Context, token string) (*entities.RefreshToken, error)
	GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.RefreshToken, error)

	DeleteByToken(ctx context.Context, token string) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}

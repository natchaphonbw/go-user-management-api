package repositories

import (
	"context"

	"github.com/google/uuid"

	"github.com/natchaphonbw/usermanagement/modules/users/entities"
)

type SessionRepository interface {
	Insert(ctx context.Context, session *entities.Session) error
	GetByID(ctx context.Context, sessionID uuid.UUID) (*entities.Session, error)
	MarkRevoked(ctx context.Context, sessionID uuid.UUID) error
	MarkRevokedByUserID(ctx context.Context, userID uuid.UUID) error
}

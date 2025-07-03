package repositories

import (
	"context"

	"gorm.io/gorm"

	"github.com/google/uuid"
	"github.com/natchaphonbw/usermanagement/modules/users/entities"
)

type sessionPostgresRepository struct {
	db *gorm.DB
}

func NewSessionPostgresRepository(db *gorm.DB) SessionRepository {
	return &sessionPostgresRepository{db: db}

}

// Insert
func (r *sessionPostgresRepository) Insert(ctx context.Context, session *entities.Session) error {
	return r.db.WithContext(ctx).Create(session).Error
}

// GetByID
func (r *sessionPostgresRepository) GetByID(ctx context.Context, sessionID uuid.UUID) (*entities.Session, error) {
	var session entities.Session
	err := r.db.WithContext(ctx).First(&session, "id = ?", sessionID).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// GetAllByUserID
func (r *sessionPostgresRepository) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.Session, error) {
	var sessions []*entities.Session
	result := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&sessions)

	if result.Error != nil {
		return nil, result.Error
	}

	return sessions, nil
}

// MarkRevoked
func (r *sessionPostgresRepository) MarkRevoked(ctx context.Context, sessionID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&entities.Session{}).
		Where("id = ?", sessionID).
		Update("revoked", true).Error
}

// MarkAllRevokedByUserID
func (r *sessionPostgresRepository) MarkRevokedByUserID(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&entities.Session{}).
		Where("user_id = ?", userID).
		Update("revoked", true).Error
}

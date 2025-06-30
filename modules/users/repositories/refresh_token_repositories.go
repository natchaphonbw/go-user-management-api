package repositories

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/google/uuid"
	"github.com/natchaphonbw/usermanagement/modules/users/entities"
)

type refreshTokenPostgresRepository struct {
	db *gorm.DB
}

func NewRefreshTokenPostgresRepository(db *gorm.DB) RefreshTokenRepository {
	return &refreshTokenPostgresRepository{db: db}

}

// Upsert
func (r *refreshTokenPostgresRepository) Upsert(ctx context.Context, token *entities.RefreshToken) error {
	result := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "user_id"},
			{Name: "device_id"},
			{Name: "device_ua"},
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"token", "issued_at", "expires_at", "device_ip",
		}),
	}).Create(token)

	return result.Error
}

// GetByToken
func (r *refreshTokenPostgresRepository) GetByToken(ctx context.Context, token string) (*entities.RefreshToken, error) {
	var refreshToken entities.RefreshToken
	result := r.db.WithContext(ctx).First(&refreshToken, "token = ?", token)

	if result.Error != nil {
		return nil, result.Error
	}

	return &refreshToken, nil
}

// GetAllByUserID
func (r *refreshTokenPostgresRepository) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.RefreshToken, error) {
	var refreshTokens []*entities.RefreshToken
	result := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&refreshTokens)

	if result.Error != nil {
		return nil, result.Error
	}

	return refreshTokens, nil
}

// DeleteByToken
func (r *refreshTokenPostgresRepository) DeleteByToken(ctx context.Context, token string) error {
	var refreshToken entities.RefreshToken
	findResult := r.db.WithContext(ctx).First(&refreshToken, "token = ?", token)

	if findResult.Error != nil {
		return findResult.Error
	}

	// delete
	DeleteResult := r.db.WithContext(ctx).Delete(&refreshToken)
	if DeleteResult.Error != nil {
		return DeleteResult.Error
	}

	return nil
}

// DeleteByUserID
func (r *refreshTokenPostgresRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	var refreshTokens []*entities.RefreshToken
	findResult := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&refreshTokens)

	if findResult.Error != nil {
		return findResult.Error
	}

	// delete
	DeleteResult := r.db.WithContext(ctx).Delete(&refreshTokens)
	if DeleteResult.Error != nil {
		return DeleteResult.Error
	}

	return nil
}

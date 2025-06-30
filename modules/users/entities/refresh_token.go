package entities

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key" json:"id" validate:"required"`
	User      User      `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE;" json:"user"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_device_user_id_ua" json:"user_id" validate:"required"`
	Token     string    `json:"token" validate:"required"`
	DeviceID  string    `gorm:"not null;uniqueIndex:idx_device_user_id_ua" json:"device_id" validate:"required"`
	DeviceUA  string    `gorm:"not null;uniqueIndex:idx_device_user_id_ua" json:"device_ua" validate:"required"`
	DeviceIP  string    `json:"device_ip"`
	ExpiresAt time.Time `json:"expires_at" validate:"required"`
	IssuedAt  time.Time `json:"issued_at" validate:"required"`
}

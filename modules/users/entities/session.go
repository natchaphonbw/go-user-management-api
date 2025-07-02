package entities

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id" validate:"required"`
	User         User      `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE;" json:"user"`
	UserID       uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id" validate:"required"`

	HashedToken  string    `gorm:"not null" json:"-"`
	DeviceID     string    `gorm:"not null" json:"device_id" validate:"required"`
	DeviceUA     string    `gorm:"not null" json:"device_ua" validate:"required"`
	DeviceIP     string    `json:"device_ip"`

	IssuedAt     time.Time `gorm:"not null" json:"issued_at" validate:"required"`
	ExpiresAt    time.Time `gorm:"not null" json:"expires_at" validate:"required"`
	Revoked      bool      `gorm:"default:false" json:"revoked"`
}
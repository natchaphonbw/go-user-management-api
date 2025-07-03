package entities

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name         string    `gorm:"type:varchar(100);not null" json:"name"`
	Email        string    `gorm:"type:varchar(100);unique;not null" json:"email"`
	PasswordHash string    `gorm:"not null" json:"-"`
	Salt         string    `gorm:"not null" json:"-"`
	Age          int       `gorm:"type:int;not null" json:"age"`
	Created_at   time.Time `gorm:"type:timestamp;default:current_timestamp" json:"created_at"`
	Updated_at   time.Time `gorm:"type:timestamp;default:current_timestamp" json:"updated_at"`
}

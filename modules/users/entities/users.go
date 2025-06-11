package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name       string    `gorm:"type:varchar(100);not null" json:"name"`
	Email      string    `gorm:"type:varchar(100);unique;not null" json:"email"`
	Age        int       `gorm:"type:int;not null" json:"age"`
	Created_at time.Time `gorm:"type:timestamp;default:current_timestamp" json:"created_at"`
	Updated_at time.Time `gorm:"type:timestamp;default:current_timestamp" json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}

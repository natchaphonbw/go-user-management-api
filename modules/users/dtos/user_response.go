package dtos

import (
	"time"

	"github.com/google/uuid"
	"github.com/natchaphonbw/usermanagement/modules/users/entities"
)

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Age       int       `json:"age"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func FromUserEntity(user *entities.User) *UserResponse {
	return &UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Age:       user.Age,
		CreatedAt: user.Created_at,
		UpdatedAt: user.Updated_at,
	}
}

func FromUserEntities(users []entities.User) []*UserResponse {
	var userResponse []*UserResponse
	for _, user := range users {
		userResponse = append(userResponse, FromUserEntity(&user))
	}
	return userResponse
}

package dtos

import (
	"time"

	"github.com/google/uuid"
	"github.com/natchaphonbw/usermanagement/modules/users/entities"
)

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Age      int    `json:"age" validate:"required,min=13"`
	Password string `json:"password" validate:"required,min=6"`
}

type UpdateUserRequest struct {
	Name  *string `json:"name" validate:"omitempty,max=100"`
	Email *string `json:"email" validate:"omitempty,email"`
	Age   *int    `json:"age" validate:"omitempty,min=13"`
}

// Response

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

package dtos

type CreateUserRequest struct {
	Name  string `json:"name" validate:"required,max=100"`
	Email string `json:"email" validate:"required,email"`
	Age   int    `json:"age" validate:"required,min=13"`
}

type UpdateUserRequest struct {
	Name  *string `json:"name" validate:"omitempty,max=100"`
	Email *string `json:"email" validate:"omitempty,email"`
	Age   *int    `json:"age" validate:"omitempty,min=13"`
}

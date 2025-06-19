package errors

import "github.com/gofiber/fiber/v2"

type ErrorResponse struct {
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

func Send(c *fiber.Ctx, appErr *AppError) error {
	return c.Status(appErr.Code).JSON(ErrorResponse{
		Message: appErr.Message,
	})
}

func SendWithDetail(c *fiber.Ctx, appErr *AppError, details interface{}) error {
	return c.Status(appErr.Code).JSON(ErrorResponse{
		Message: appErr.Message,
		Details: details,
	})
}

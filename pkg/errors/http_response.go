package errors

import "github.com/gofiber/fiber/v2"

type ErrorResponse struct {
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

func Send(c *fiber.Ctx, appErr *AppError) error {

	resp := ErrorResponse{
		Message: appErr.Message,
	}
	if appErr.Details != nil {
		resp.Details = appErr.Details
	}
	return c.Status(appErr.Code).JSON(resp)
}

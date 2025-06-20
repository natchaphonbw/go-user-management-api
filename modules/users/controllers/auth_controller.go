package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/natchaphonbw/usermanagement/modules/users/dtos"
	"github.com/natchaphonbw/usermanagement/modules/users/usecases"
	"github.com/natchaphonbw/usermanagement/modules/users/validator"
	app_errors "github.com/natchaphonbw/usermanagement/pkg/errors"
)

type AuthController struct {
	authUseCase usecases.AuthUsecase
}

func NewAuthController(u usecases.AuthUsecase) *AuthController {
	return &AuthController{
		authUseCase: u,
	}
}

// Register
func (a *AuthController) Register(c *fiber.Ctx) error {
	var req dtos.RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return app_errors.Send(c, app_errors.BadRequest("Invalid request body", err))

	}

	if errs := validator.ValidateStruct(req); len(errs) > 0 {
		return app_errors.SendWithDetail(c, app_errors.BadRequest("Validation failed", nil), errs)
	}

	userResp, respErr := a.authUseCase.RegisterUser(c.Context(), req)
	if respErr != nil {
		return app_errors.Send(c, respErr)
	}

	return c.Status(fiber.StatusCreated).JSON(userResp)

}

// Login
func (a *AuthController) Login(c *fiber.Ctx) error {
	var req dtos.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return app_errors.Send(c, app_errors.BadRequest("Invalid request body", err))
	}

	if errs := validator.ValidateStruct(req); len(errs) > 0 {
		return app_errors.SendWithDetail(c, app_errors.BadRequest("Validation failed", nil), errs)
	}

	loginResp, respErr := a.authUseCase.Login(c.Context(), req)
	if respErr != nil {
		return app_errors.Send(c, respErr)
	}

	return c.Status(fiber.StatusOK).JSON(loginResp)
}

// GetProfile
func (a *AuthController) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	profile, err := a.authUseCase.GetProfile(c.Context(), userID)
	if err != nil {
		return app_errors.Send(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(profile)

}

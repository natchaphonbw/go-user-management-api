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
	authUseCase    usecases.AuthUsecase
	refreshUseCase usecases.RefreshTokenUsecase
}

func NewAuthController(u usecases.AuthUsecase, r usecases.RefreshTokenUsecase) *AuthController {
	return &AuthController{
		authUseCase:    u,
		refreshUseCase: r,
	}
}

// Register
func (a *AuthController) Register(c *fiber.Ctx) error {
	var req dtos.RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return app_errors.Send(c, app_errors.BadRequest("Invalid request body", err))

	}

	if errs := validator.ValidateStruct(req); len(errs) > 0 {
		return app_errors.Send(c, app_errors.BadRequest("Validation failed", nil).WithDetails(errs))
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
		return app_errors.Send(c, app_errors.BadRequest("Validation failed", nil).WithDetails(errs))
	}

	deviceIP := c.IP()
	deviceUA := c.Get("User-Agent")
	deviceID := c.Get("X-Device-ID")

	loginResp, respErr := a.authUseCase.Login(c.Context(), req, deviceIP, deviceUA, deviceID)
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

// RefreshToken
func (a *AuthController) RefreshToken(c *fiber.Ctx) error {
	var req dtos.RefreshTokenRequest

	if err := c.BodyParser(&req); err != nil {
		return app_errors.Send(c, app_errors.BadRequest("Invalid request body", err))
	}

	if errs := validator.ValidateStruct(req); len(errs) > 0 {
		return app_errors.Send(c, app_errors.BadRequest("Validation failed", nil).WithDetails(errs))
	}

	deviceIP := c.IP()
	deviceUA := c.Get("User-Agent")
	deviceID := c.Get("X-Device-ID")

	tokenPair, respErr := a.refreshUseCase.Refresh(c.Context(), req.RefreshToken, deviceIP, deviceUA, deviceID)
	if respErr != nil {
		return app_errors.Send(c, respErr)
	}

	return c.Status(fiber.StatusOK).JSON(tokenPair)
}

// Logout
func (a *AuthController) Logout(c *fiber.Ctx) error {
	var req dtos.RefreshTokenRequest

	if err := c.BodyParser(&req); err != nil {
		return app_errors.Send(c, app_errors.BadRequest("Invalid request body", err))
	}

	if errs := validator.ValidateStruct(req); len(errs) > 0 {
		return app_errors.Send(c, app_errors.BadRequest("Validation failed", nil).WithDetails(errs))
	}

	userID := c.Locals("userID").(uuid.UUID)
	deviceUA := c.Get("User-Agent")
	deviceID := c.Get("X-Device-ID")

	if logoutErr := a.authUseCase.Logout(c.Context(), userID, req.RefreshToken, deviceID, deviceUA); logoutErr != nil {
		return app_errors.Send(c, logoutErr)
	}

	return c.SendStatus(fiber.StatusNoContent)

}

// LogoutAll
func (a *AuthController) LogoutAll(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	if logoutErr := a.authUseCase.LogoutAll(c.Context(), userID); logoutErr != nil {
		return app_errors.Send(c, logoutErr)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

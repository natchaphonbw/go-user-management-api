package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/natchaphonbw/usermanagement/modules/users/dtos"
	"github.com/natchaphonbw/usermanagement/modules/users/usecases"
	"github.com/natchaphonbw/usermanagement/modules/users/validations"
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})

	}

	if err := validations.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Validation failed", "detail": err.Error()})
	}

	userResp, err := a.authUseCase.RegisterUser(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(userResp)

}

// Login
func (a *AuthController) Login(c *fiber.Ctx) error {
	var req dtos.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if err := validations.ValidateStruct(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	loginResp, err := a.authUseCase.Login(c.Context(), req)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(loginResp)
}

// GetProfile
func (a *AuthController) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	profile, err := a.authUseCase.GetProfile(c.Context(), userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(profile)

}

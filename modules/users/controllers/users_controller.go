package controllers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/natchaphonbw/usermanagement/modules/users/dtos"
	"github.com/natchaphonbw/usermanagement/modules/users/usecases"
	"github.com/natchaphonbw/usermanagement/modules/users/validations"
	"gorm.io/gorm"
)

type UserController struct {
	userUsecase usecases.UserUsecase
}

// instance
func NewUserController(u usecases.UserUsecase) *UserController {
	return &UserController{
		userUsecase: u,
	}
}

// handler

// CreateUser
func (ctrl *UserController) CreateUser(c *fiber.Ctx) error {
	var req dtos.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// validate
	if err := validations.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":  "Validation failed",
			"detail": err.Error(),
		})
	}

	// call usecase
	userResp, err := ctrl.userUsecase.CreateUser(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error()})

	}

	return c.Status(fiber.StatusCreated).JSON(userResp)
}

// GetAllUsers
func (ctrl *UserController) GetAllUsers(c *fiber.Ctx) error {
	usersResp, err := ctrl.userUsecase.GetAllUsers(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(usersResp)
}

// GetUserByID
func (ctrl *UserController) GetUserByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	log.Printf("Fetching user with ID: %s", id)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	userResp, err := ctrl.userUsecase.GetUserByID(c.Context(), id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(userResp)
}

// UpdateUserByID
func (ctrl *UserController) UpdateUserByID(c *fiber.Ctx) error {
	// parse user ID
	id, err := uuid.Parse(c.Params("id"))
	log.Printf("Fetching user with ID: %s", id)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// parse request body
	var req dtos.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// validate
	if err := validations.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":  "Validation failed",
			"detail": err.Error(),
		})
	}

	// call usecase
	userResp, err := ctrl.userUsecase.UpdateUserByID(c.Context(), id, req)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error()})

	}

	return c.Status(fiber.StatusCreated).JSON(userResp)
}

// DeleteUserByID
func (ctrl *UserController) DeleteUserByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	log.Printf("Fetching user with ID: %s", id)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}
	// Delete user
	deletedUser, err := ctrl.userUsecase.DeleteUserByID(c.Context(), id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})

	}

	return c.Status(fiber.StatusOK).JSON(deletedUser)
}

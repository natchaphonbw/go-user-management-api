package controllers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/natchaphonbw/usermanagement/modules/users/dtos"
	"github.com/natchaphonbw/usermanagement/modules/users/usecases"
	"github.com/natchaphonbw/usermanagement/modules/users/validator"
	app_errors "github.com/natchaphonbw/usermanagement/pkg/errors"
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
		return app_errors.Send(c, app_errors.BadRequest("Invalid request body", err))
	}

	// validate
	if errs := validator.ValidateStruct(req); len(errs) > 0 {
		return app_errors.Send(c, app_errors.BadRequest("Validation failed", nil).WithDetails(errs))
	}

	// call usecase
	userResp, err := ctrl.userUsecase.CreateUser(c.Context(), req)
	if err != nil {
		return app_errors.Send(c, err)

	}

	return c.Status(fiber.StatusCreated).JSON(userResp)
}

// GetAllUsers
func (ctrl *UserController) GetAllUsers(c *fiber.Ctx) error {
	usersResp, err := ctrl.userUsecase.GetAllUsers(c.Context())
	if err != nil {
		return app_errors.Send(c, err)
	}

	return c.JSON(usersResp)
}

// GetUserByID
func (ctrl *UserController) GetUserByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	log.Printf("Fetching user with ID: %s", id)

	if err != nil {
		return app_errors.Send(c, app_errors.BadRequest("Invalid user ID", err))
	}

	userResp, respErr := ctrl.userUsecase.GetUserByID(c.Context(), id)
	if respErr != nil {
		return app_errors.Send(c, respErr)
	}

	return c.Status(fiber.StatusOK).JSON(userResp)
}

// UpdateUserByID
func (ctrl *UserController) UpdateUserByID(c *fiber.Ctx) error {
	// parse user ID
	id, err := uuid.Parse(c.Params("id"))
	log.Printf("Fetching user with ID: %s", id)

	if err != nil {
		return app_errors.Send(c, app_errors.BadRequest("Invalid user ID", err))
	}

	// parse request body
	var req dtos.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return app_errors.Send(c, app_errors.BadRequest("Invalid request body", err))
	}

	// validate
	if errs := validator.ValidateStruct(req); len(errs) > 0 {
		return app_errors.Send(c, app_errors.BadRequest("Validation failed", nil).WithDetails(errs))
	}

	// call usecase
	userResp, respErr := ctrl.userUsecase.UpdateUserByID(c.Context(), id, req)
	if respErr != nil {
		return app_errors.Send(c, respErr)

	}

	return c.Status(fiber.StatusCreated).JSON(userResp)
}

// DeleteUserByID
func (ctrl *UserController) DeleteUserByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	log.Printf("Fetching user with ID: %s", id)

	if err != nil {
		return app_errors.Send(c, app_errors.BadRequest("Invalid user ID", err))
	}
	// Delete user
	deletedUser, delErr := ctrl.userUsecase.DeleteUserByID(c.Context(), id)
	if delErr != nil {
		return app_errors.Send(c, delErr)

	}

	return c.JSON(deletedUser)
}

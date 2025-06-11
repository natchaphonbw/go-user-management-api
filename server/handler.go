package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/natchaphonbw/usermanagement/modules/users/controllers"
	"github.com/natchaphonbw/usermanagement/modules/users/repositories"
	"github.com/natchaphonbw/usermanagement/modules/users/usecases"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {

	userRepo := repositories.NewUserPostgresRepository(db)
	userUseCase := usecases.NewUserUseCase(userRepo)
	userController := controllers.NewUserController(userUseCase)

	userGroup := app.Group("/")
	userGroup.Post("/users", userController.CreateUser)
	userGroup.Get("/users", userController.GetAllUsers)
	userGroup.Get("/users/:id", userController.GetUserByID)
	userGroup.Put("/users/:id", userController.UpdateUserByID)
	userGroup.Delete("/users/:id", userController.DeleteUserByID)

}

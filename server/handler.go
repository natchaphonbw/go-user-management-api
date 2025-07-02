package server

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/natchaphonbw/usermanagement/modules/users/controllers"
	"github.com/natchaphonbw/usermanagement/modules/users/repositories"
	"github.com/natchaphonbw/usermanagement/modules/users/usecases"
	"github.com/natchaphonbw/usermanagement/pkg/middlewares"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {

	userRepo := repositories.NewUserPostgresRepository(db)
	userUseCase := usecases.NewUserUseCase(userRepo)
	sessionRepo := repositories.NewSessionPostgresRepository(db)
	sessionUseCase := usecases.NewSessionUsecase(sessionRepo)
	authUseCase := usecases.NewAuthUseCase(userUseCase, sessionUseCase, userRepo, sessionRepo)

	userController := controllers.NewUserController(userUseCase)
	authController := controllers.NewAuthController(authUseCase, sessionUseCase)

	userGroup := app.Group("/users")
	userGroup.Post("/", userController.CreateUser)
	userGroup.Get("/", userController.GetAllUsers)
	userGroup.Get("/:id", userController.GetUserByID)
	userGroup.Put("/:id", userController.UpdateUserByID)
	userGroup.Delete("/:id", userController.DeleteUserByID)

	authPublic := app.Group("/auth")
	authPublic.Post("/register", authController.Register)
	authPublic.Post("/login", authController.Login)
	authPublic.Post("/refresh", authController.RefreshToken)

	authProtect := app.Group("/auth", middlewares.JWTAuthMiddleware())
	authProtect.Get("/me", authController.GetProfile)
	authProtect.Post("/logout", authController.Logout)
	authProtect.Post("/logout/all", authController.LogoutAll)

}

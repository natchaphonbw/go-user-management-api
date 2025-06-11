package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/natchaphonbw/usermanagement/config"
	"github.com/natchaphonbw/usermanagement/pkg/databases"
	"github.com/natchaphonbw/usermanagement/pkg/databases/migrations"
	"github.com/natchaphonbw/usermanagement/server"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()
	// Connect to the database
	db := databases.Connect()
	// Run migrations
	migrations.Migrate(db)

	app := server.NewFiberApp()

	server.SetupRoutes(app, db)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"fiberHost": cfg.FiberHost,
			"fiberPort": cfg.FiberPort,
			"dbHost":    cfg.DBHost,
			"dbPort":    cfg.DBPort,
		})
	})

	addr := fmt.Sprintf("%s:%s", cfg.FiberHost, cfg.FiberPort)
	app.Listen(addr)

}

package main

import (
	"fmt"

	"github.com/natchaphonbw/usermanagement/config"
	"github.com/natchaphonbw/usermanagement/pkg/databases"
	"github.com/natchaphonbw/usermanagement/pkg/databases/migrations"
	"github.com/natchaphonbw/usermanagement/server"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()
	// Connect to the database
	db := databases.Connect(cfg)
	// Run migrations
	migrations.Migrate(db)

	app := server.NewFiberApp()

	server.SetupRoutes(app, db)

	addr := fmt.Sprintf("%s:%s", cfg.FiberHost, cfg.FiberPort)
	app.Listen(addr)

}

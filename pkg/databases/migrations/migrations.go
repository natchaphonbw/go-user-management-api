package migrations

import (
	"log"

	"github.com/natchaphonbw/usermanagement/modules/users/entities"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&entities.User{},
	)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Migrations completed successfully")

}

package databases

import (
	"fmt"
	"log"

	"github.com/natchaphonbw/usermanagement/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() *gorm.DB {
	cfg := config.LoadConfig()

	host := cfg.DBHost
	port := cfg.DBPort
	user := cfg.DBUser
	password := cfg.DBPassword
	dbName := cfg.DBName

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)

	}

	log.Println("Connected to the PostgreSQL database successfully")
	return db
}

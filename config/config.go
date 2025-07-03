package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	FiberHost string
	FiberPort string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	JWT_ACCESS_SECRET  string
	JWT_REFRESH_SECRET string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using default values")
	}

	return &Config{
		FiberHost: getEnv("FIBER_HOST", "0.0.0.0"),
		FiberPort: getEnv("FIBER_PORT", "5000"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "user"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "user_db"),
	}
}

func getEnv(key, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}

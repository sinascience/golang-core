package configs

import (
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	JWTSecretKey string
}

// LoadConfig loads application configuration from .env file
func LoadConfig() (config Config, err error) {
	err = godotenv.Load()
	if err != nil {
		return
	}

	config.DBHost = os.Getenv("DB_HOST")
	config.DBPort = os.Getenv("DB_PORT")
	config.DBUser = os.Getenv("DB_USER")
	config.DBPassword = os.Getenv("DB_PASSWORD")
	config.DBName = os.Getenv("DB_NAME")

	config.JWTSecretKey = os.Getenv("JWT_SECRET_KEY")
	return
}

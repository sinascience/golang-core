package configs

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	JWTAccessSecret          string
	JWTRefreshSecret         string
	JWTAccessTokenExpiresIn  time.Duration
	JWTRefreshTokenExpiresIn time.Duration

	CORSAllowedOrigins string
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

	config.CORSAllowedOrigins = os.Getenv("CORS_ALLOWED_ORIGINS")

	config.JWTAccessSecret = os.Getenv("JWT_ACCESS_SECRET_KEY")
	config.JWTRefreshSecret = os.Getenv("JWT_REFRESH_SECRET_KEY")

	accessExp, _ := strconv.Atoi(os.Getenv("JWT_ACCESS_EXPIRATION_IN_MINUTES"))
	config.JWTAccessTokenExpiresIn = time.Duration(accessExp) * time.Minute

	refreshExp, _ := strconv.Atoi(os.Getenv("JWT_REFRESH_EXPIRATION_IN_HOURS"))
	config.JWTRefreshTokenExpiresIn = time.Duration(refreshExp) * time.Hour
	return
}

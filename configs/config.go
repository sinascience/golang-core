package configs

import (
	"fmt"
	"os"
	"strconv"
	"strings"
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

	// Storage configuration
	StorageBucketName string
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

	config.StorageBucketName = os.Getenv("STORAGE_BUCKET_NAME")

	// Validate required configuration
	if err := validateConfig(&config); err != nil {
		return config, err
	}

	return config, nil
}

// validateConfig validates that all required configuration values are present
func validateConfig(config *Config) error {
	var missing []string

	// Required database configuration
	if config.DBHost == "" {
		missing = append(missing, "DB_HOST")
	}
	if config.DBPort == "" {
		missing = append(missing, "DB_PORT")
	}
	if config.DBUser == "" {
		missing = append(missing, "DB_USER")
	}
	if config.DBName == "" {
		missing = append(missing, "DB_NAME")
	}

	// Required JWT configuration
	if config.JWTAccessSecret == "" {
		missing = append(missing, "JWT_ACCESS_SECRET_KEY")
	}
	if config.JWTRefreshSecret == "" {
		missing = append(missing, "JWT_REFRESH_SECRET_KEY")
	}

	// Check JWT secret strength (minimum 32 characters)
	if len(config.JWTAccessSecret) < 32 {
		missing = append(missing, "JWT_ACCESS_SECRET_KEY (must be at least 32 characters)")
	}
	if len(config.JWTRefreshSecret) < 32 {
		missing = append(missing, "JWT_REFRESH_SECRET_KEY (must be at least 32 characters)")
	}

	// Check token expiration values
	if config.JWTAccessTokenExpiresIn <= 0 {
		missing = append(missing, "JWT_ACCESS_EXPIRATION_IN_MINUTES (must be positive)")
	}
	if config.JWTRefreshTokenExpiresIn <= 0 {
		missing = append(missing, "JWT_REFRESH_EXPIRATION_IN_HOURS (must be positive)")
	}

	// Optional but recommended configuration warnings
	if config.StorageBucketName == "" {
		fmt.Println("⚠️  Warning: STORAGE_BUCKET_NAME not set - file uploads may fail")
	}
	if config.CORSAllowedOrigins == "" {
		fmt.Println("⚠️  Warning: CORS_ALLOWED_ORIGINS not set - using wildcard (*)")
		config.CORSAllowedOrigins = "*"
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}

	return nil
}

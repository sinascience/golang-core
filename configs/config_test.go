package configs

import (
	"os"
	"testing"
	"time"
)

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "Valid configuration",
			config: Config{
				DBHost:                   "localhost",
				DBPort:                   "3306",
				DBUser:                   "root",
				DBName:                   "test_db",
				JWTAccessSecret:          "this-is-a-32-character-secret-key",
				JWTRefreshSecret:         "this-is-another-32-char-secret-key",
				JWTAccessTokenExpiresIn:  15 * time.Minute,
				JWTRefreshTokenExpiresIn: 24 * time.Hour,
				CORSAllowedOrigins:       "*",
				StorageBucketName:        "test-bucket",
			},
			wantErr: false,
		},
		{
			name: "Missing database host",
			config: Config{
				DBPort:                   "3306",
				DBUser:                   "root",
				DBName:                   "test_db",
				JWTAccessSecret:          "this-is-a-32-character-secret-key",
				JWTRefreshSecret:         "this-is-another-32-char-secret-key",
				JWTAccessTokenExpiresIn:  15 * time.Minute,
				JWTRefreshTokenExpiresIn: 24 * time.Hour,
			},
			wantErr: true,
		},
		{
			name: "JWT secret too short",
			config: Config{
				DBHost:                   "localhost",
				DBPort:                   "3306",
				DBUser:                   "root",
				DBName:                   "test_db",
				JWTAccessSecret:          "short", // Too short
				JWTRefreshSecret:         "this-is-another-32-char-secret-key",
				JWTAccessTokenExpiresIn:  15 * time.Minute,
				JWTRefreshTokenExpiresIn: 24 * time.Hour,
			},
			wantErr: true,
		},
		{
			name: "Invalid token expiration",
			config: Config{
				DBHost:                   "localhost",
				DBPort:                   "3306",
				DBUser:                   "root",
				DBName:                   "test_db",
				JWTAccessSecret:          "this-is-a-32-character-secret-key",
				JWTRefreshSecret:         "this-is-another-32-char-secret-key",
				JWTAccessTokenExpiresIn:  0, // Invalid
				JWTRefreshTokenExpiresIn: 24 * time.Hour,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(&tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoadConfig_MissingEnvFile(t *testing.T) {
	// Temporarily remove .env file if it exists
	envFile := ".env"
	if _, err := os.Stat(envFile); err == nil {
		// File exists, rename it temporarily
		tempName := ".env.backup"
		os.Rename(envFile, tempName)
		defer os.Rename(tempName, envFile) // Restore after test
	}

	// Set required environment variables manually
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "3306")
	os.Setenv("DB_USER", "root")
	os.Setenv("DB_NAME", "test_db")
	os.Setenv("JWT_ACCESS_SECRET_KEY", "this-is-a-32-character-secret-key")
	os.Setenv("JWT_REFRESH_SECRET_KEY", "this-is-another-32-char-secret-key")
	os.Setenv("JWT_ACCESS_EXPIRATION_IN_MINUTES", "15")
	os.Setenv("JWT_REFRESH_EXPIRATION_IN_HOURS", "24")

	defer func() {
		// Clean up environment variables
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("JWT_ACCESS_SECRET_KEY")
		os.Unsetenv("JWT_REFRESH_SECRET_KEY")
		os.Unsetenv("JWT_ACCESS_EXPIRATION_IN_MINUTES")
		os.Unsetenv("JWT_REFRESH_EXPIRATION_IN_HOURS")
	}()

	config, err := LoadConfig()
	if err != nil {
		t.Errorf("LoadConfig() should work with environment variables, got error: %v", err)
	}

	// Verify some key values
	if config.DBHost != "localhost" {
		t.Errorf("Expected DB_HOST=localhost, got %s", config.DBHost)
	}
	if config.JWTAccessTokenExpiresIn != 15*time.Minute {
		t.Errorf("Expected access token expiration=15m, got %v", config.JWTAccessTokenExpiresIn)
	}
}
package service

import (
	"context"
	"testing"
	"time"
	"venturo-core/configs"
	"venturo-core/internal/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&model.User{}, &model.RefreshToken{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func setupTestConfig() *configs.Config {
	return &configs.Config{
		JWTAccessSecret:          "test-access-secret-key-32-chars-long",
		JWTRefreshSecret:         "test-refresh-secret-key-32-chars-long",
		JWTAccessTokenExpiresIn:  15 * time.Minute,
		JWTRefreshTokenExpiresIn: 24 * time.Hour,
	}
}

func TestAuthService_Register(t *testing.T) {
	db := setupTestDB(t)
	config := setupTestConfig()
	authService := NewAuthService(db, config)

	tests := []struct {
		name     string
		email    string
		password string
		wantErr  bool
	}{
		{
			name:     "Valid registration",
			email:    "test@example.com",
			password: "validpassword123",
			wantErr:  false,
		},
		{
			name:     "Duplicate email registration",
			email:    "test@example.com", // Same as above
			password: "anotherpassword123",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := authService.Register(ctx, "Test User", tt.email, tt.password)

			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	db := setupTestDB(t)
	config := setupTestConfig()
	authService := NewAuthService(db, config)

	// First register a user
	ctx := context.Background()
	email := "test@example.com"
	password := "validpassword123"

	err := authService.Register(ctx, "Test User", email, password)
	if err != nil {
		t.Fatalf("Failed to register test user: %v", err)
	}

	tests := []struct {
		name     string
		email    string
		password string
		wantErr  bool
	}{
		{
			name:     "Valid login",
			email:    email,
			password: password,
			wantErr:  false,
		},
		{
			name:     "Invalid email",
			email:    "nonexistent@example.com",
			password: password,
			wantErr:  true,
		},
		{
			name:     "Invalid password",
			email:    email,
			password: "wrongpassword",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := authService.Login(ctx, tt.email, tt.password)

			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if tokens["access_token"] == "" || tokens["refresh_token"] == "" {
					t.Error("Login() should return both access and refresh tokens")
				}
			}
		})
	}
}

func TestAuthService_RefreshToken(t *testing.T) {
	db := setupTestDB(t)
	config := setupTestConfig()
	authService := NewAuthService(db, config)

	// Setup: Register and login a user
	ctx := context.Background()
	email := "test@example.com"
	password := "validpassword123"

	err := authService.Register(ctx, "Test User", email, password)
	if err != nil {
		t.Fatalf("Failed to register test user: %v", err)
	}

	tokens, err := authService.Login(ctx, email, password)
	if err != nil {
		t.Fatalf("Failed to login test user: %v", err)
	}

	tests := []struct {
		name         string
		refreshToken string
		wantErr      bool
	}{
		{
			name:         "Valid refresh token",
			refreshToken: tokens["refresh_token"],
			wantErr:      false,
		},
		{
			name:         "Invalid refresh token",
			refreshToken: "invalid-token",
			wantErr:      true,
		},
		{
			name:         "Empty refresh token",
			refreshToken: "",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newTokens, err := authService.RefreshToken(ctx, tt.refreshToken)

			if (err != nil) != tt.wantErr {
				t.Errorf("RefreshToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if newTokens["access_token"] == "" || newTokens["refresh_token"] == "" {
					t.Error("RefreshToken() should return both access and refresh tokens")
				}

				// Verify new tokens are different from old ones
				if newTokens["access_token"] == tokens["access_token"] {
					t.Error("RefreshToken() should return a new access token")
				}
				if newTokens["refresh_token"] == tokens["refresh_token"] {
					t.Error("RefreshToken() should return a new refresh token")
				}
			}
		})
	}
}

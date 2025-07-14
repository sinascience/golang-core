package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// GenerateToken creates a new JWT for a given user.
func GenerateToken(userID uuid.UUID, secretKey string) (string, error) {
	// Create the claims
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // Token expires in 72 hours
		"iat":     time.Now().Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and return it
	return token.SignedString([]byte(secretKey))
}

// Generate refresh token for a given user.
func GenerateRefreshToken(userId uuid.UUID, secretKey string) (string, error) {
	// create the claims
	claims := jwt.MapClaims{
		"user_id": userId.String(),
		"exp":     time.Now().Add(time.Hour * 168).Unix(), // Token expires in 72 hours
		"iat":     time.Now().Unix(),
	}

	// create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// generate encoded token and return it
	return token.SignedString([]byte(secretKey))
}

// Generate new token for a given user.
func GenerateNewToken(userId string, secretKey string) (string, error) {
	// create the claims
	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // Token expires in 72 hours
		"iat":     time.Now().Unix(),
	}

	// create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// generate encoded token and return it
	return token.SignedString([]byte(secretKey))
}

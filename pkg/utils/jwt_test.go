package utils

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestGenerateAndParseTokens(t *testing.T) {
	// Define a sample user ID for testing.
	userID := uuid.New()

	// Define test cases for both access and refresh tokens.
	testCases := []struct {
		name     string
		genFunc  func(userID uuid.UUID, secretKey string, expiresAt time.Duration) (string, error)
		secret   string
		duration time.Duration
	}{
		{
			name:     "AccessToken",
			genFunc:  GenerateAccessToken,
			secret:   "test-secret-for-access",
			duration: time.Minute * 15,
		},
		{
			name:     "RefreshToken",
			genFunc:  GenerateRefreshToken,
			secret:   "test-secret-for-refresh",
			duration: time.Hour * 24,
		},
	}

	for _, tc := range testCases {
		// Use t.Run to create a sub-test for each case.
		t.Run(tc.name, func(t *testing.T) {
			// 1. Generate the token.
			tokenString, err := tc.genFunc(userID, tc.secret, tc.duration)
			if err != nil {
				t.Fatalf("Failed to generate token: %v", err)
			}

			// 2. Parse the token back to validate it.
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// We must return the same secret key used for signing.
				return []byte(tc.secret), nil
			})

			if err != nil {
				t.Fatalf("Failed to parse token: %v", err)
			}

			// 3. Check if the token is valid.
			if !token.Valid {
				t.Errorf("Generated token is not valid")
			}

			// 4. Check the claims to ensure the user ID is correct.
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				parsedUserID, err := uuid.Parse(claims["user_id"].(string))
				if err != nil {
					t.Fatalf("Could not parse user_id from token claims: %v", err)
				}

				if parsedUserID != userID {
					t.Errorf("Expected user ID %v, but got %v", userID, parsedUserID)
				}
			} else {
				t.Errorf("Could not read token claims")
			}
		})
	}
}

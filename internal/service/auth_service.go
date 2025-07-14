package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log/slog"
	"time"
	"venturo-core/configs"
	"venturo-core/internal/model"
	"venturo-core/pkg/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	db   *gorm.DB
	conf *configs.Config
}

// NewAuthService creates a new auth service.
func NewAuthService(db *gorm.DB, conf *configs.Config) *AuthService {
	return &AuthService{db: db, conf: conf}
}

// Register creates a new user.
func (s *AuthService) Register(ctx context.Context, name, email, password string) error {
	// Check if user already exists
	var existingUser model.User
	if err := s.db.WithContext(ctx).Where("email = ?", email).First(&existingUser).Error; err == nil {
		return errors.New("user with this email already exists")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Create new user
	newUser := model.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	}

	// Save user to the database
	if err := newUser.Save(s.db.WithContext(ctx)); err != nil {
		return err
	}

	return nil
}

// Login validates user credentials and returns a JWT.
func (s *AuthService) Login(ctx context.Context, email, password string) (map[string]string, error) {
	// Find user by email
	var user model.User
	if err := s.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Compare password with the hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// 1. Generate new access and refresh tokens
	accessToken, err := utils.GenerateAccessToken(user.ID, s.conf.JWTAccessSecret, s.conf.JWTAccessTokenExpiresIn)
	if err != nil {
		return nil, errors.New("could not generate access token")
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, s.conf.JWTRefreshSecret, s.conf.JWTRefreshTokenExpiresIn)
	if err != nil {
		return nil, errors.New("could not generate refresh token")
	}

	// 2. Hash the refresh token before saving
	hasher := sha256.New()
	hasher.Write([]byte(refreshToken))
	hashedRefreshToken := hex.EncodeToString(hasher.Sum(nil))

	// 3. Create and save the new refresh token record
	refreshTokenRecord := model.RefreshToken{
		UserID:    user.ID,
		Token:     hashedRefreshToken, // Save the SHA-256 hash
		ExpiresAt: time.Now().Add(s.conf.JWTRefreshTokenExpiresIn),
	}
	if err := refreshTokenRecord.Create(s.db); err != nil {
		return nil, errors.New("could not save session")
	}

	// 4. Return both tokens to the handler
	tokens := map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}

	return tokens, nil
}

// RefreshToken validates a refresh token and issues a new pair of tokens.
func (s *AuthService) RefreshToken(ctx context.Context, tokenString string) (map[string]string, error) {
	// 1. Parse the token to get the user ID from its claims
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.conf.JWTRefreshSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	userID, err := uuid.Parse(claims["user_id"].(string))
	if err != nil {
		return nil, errors.New("invalid user id in token")
	}

	// 2. Find all refresh tokens for the user
	var rtModel model.RefreshToken
	tokenRecords, err := rtModel.FindAllByUserID(s.db, userID)
	if err != nil || len(tokenRecords) == 0 {
		return nil, errors.New("session not found, please log in again")
	}

	// 3. Find the matching token and validate it
	incomingTokenHasher := sha256.New()
	incomingTokenHasher.Write([]byte(tokenString))
	incomingTokenHash := hex.EncodeToString(incomingTokenHasher.Sum(nil))

	var currentTokenRecord *model.RefreshToken
	for i, record := range tokenRecords {
		// Compare the hashes directly
		if record.Token == incomingTokenHash {
			currentTokenRecord = &tokenRecords[i]
			break
		}
	}

	if currentTokenRecord == nil {
		return nil, errors.New("invalid refresh token, please log in again")
	}

	// 4. (Token Rotation) Delete the used token
	if err := currentTokenRecord.Delete(s.db); err != nil {
		slog.Error("failed to delete used refresh token", "error", err)
	}

	// 5. Issue a new pair of tokens
	newAccessToken, err := utils.GenerateAccessToken(userID, s.conf.JWTAccessSecret, s.conf.JWTAccessTokenExpiresIn)
	if err != nil {
		return nil, errors.New("could not generate access token")
	}

	newRefreshToken, err := utils.GenerateRefreshToken(userID, s.conf.JWTRefreshSecret, s.conf.JWTRefreshTokenExpiresIn)
	if err != nil {
		return nil, errors.New("could not generate refresh token")
	}

	// 6. Save the new refresh token
	hashedRefreshToken, _ := bcrypt.GenerateFromPassword([]byte(newRefreshToken), bcrypt.DefaultCost)
	newRefreshTokenRecord := model.RefreshToken{
		UserID:    userID,
		Token:     string(hashedRefreshToken),
		ExpiresAt: time.Now().Add(s.conf.JWTRefreshTokenExpiresIn),
	}
	if err := newRefreshTokenRecord.Create(s.db); err != nil {
		return nil, errors.New("could not save new session")
	}

	// 7. Return the new tokens
	tokens := map[string]string{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	}

	return tokens, nil
}

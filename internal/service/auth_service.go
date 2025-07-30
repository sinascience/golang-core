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
	"venturo-core/pkg/logger"
	"venturo-core/pkg/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	db            *gorm.DB
	conf          *configs.Config
	securityLogger *logger.SecurityLogger
}

// NewAuthService creates a new auth service.
func NewAuthService(db *gorm.DB, conf *configs.Config) *AuthService {
	return &AuthService{
		db:            db,
		conf:          conf,
		securityLogger: logger.NewSecurityLogger(),
	}
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
	if err := newUser.Save(ctx, s.db); err != nil {
		return err
	}

	// Log successful registration
	s.securityLogger.LogUserRegistration(ctx, newUser.ID, email, getIPFromContext(ctx))

	return nil
}

// Login validates user credentials and returns a JWT.
func (s *AuthService) Login(ctx context.Context, email, password string) (map[string]string, error) {
	// Find user by email
	var user model.User
	if err := s.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		s.securityLogger.LogFailedLogin(ctx, email, getIPFromContext(ctx), "user not found")
		return nil, errors.New("invalid credentials")
	}

	// Compare password with the hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		s.securityLogger.LogFailedLogin(ctx, email, getIPFromContext(ctx), "invalid password")
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
	if err := refreshTokenRecord.Create(ctx, s.db); err != nil {
		return nil, errors.New("could not save session")
	}

	// Log successful login
	s.securityLogger.LogSuccessfulLogin(ctx, user.ID, email, getIPFromContext(ctx))

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

	// 2. Hash the incoming token and find it directly in the database
	incomingTokenHasher := sha256.New()
	incomingTokenHasher.Write([]byte(tokenString))
	incomingTokenHash := hex.EncodeToString(incomingTokenHasher.Sum(nil))

	// 3. Find the specific token in one query
	var rtModel model.RefreshToken
	currentTokenRecord, err := rtModel.FindByUserIDAndToken(ctx, s.db, userID, incomingTokenHash)
	if err != nil {
		s.securityLogger.LogTokenRefresh(ctx, userID, getIPFromContext(ctx), false)
		return nil, errors.New("invalid refresh token, please log in again")
	}

	// 4. (Token Rotation) Delete the used token
	if err := currentTokenRecord.Delete(ctx, s.db); err != nil {
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

	// 6. Save the new refresh token (use SHA-256 for consistency)
	newTokenHasher := sha256.New()
	newTokenHasher.Write([]byte(newRefreshToken))
	hashedNewRefreshToken := hex.EncodeToString(newTokenHasher.Sum(nil))
	
	newRefreshTokenRecord := model.RefreshToken{
		UserID:    userID,
		Token:     hashedNewRefreshToken,
		ExpiresAt: time.Now().Add(s.conf.JWTRefreshTokenExpiresIn),
	}
	if err := newRefreshTokenRecord.Create(ctx, s.db); err != nil {
		return nil, errors.New("could not save new session")
	}

	// Log successful token refresh
	s.securityLogger.LogTokenRefresh(ctx, userID, getIPFromContext(ctx), true)

	// 7. Return the new tokens
	tokens := map[string]string{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	}

	return tokens, nil
}

// getIPFromContext extracts IP address from context (set by middleware)
func getIPFromContext(ctx context.Context) string {
	if ip := ctx.Value("client_ip"); ip != nil {
		return ip.(string)
	}
	return "unknown"
}

package service

import (
	"context"
	"errors"
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
func (s *AuthService) Login(ctx context.Context, email, password string) (string, string, error) {
	// Find user by email
	var user model.User
	if err := s.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return "", "", errors.New("invalid credentials")
	}

	// Compare password with the hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	// Generate JWT
	token, err := utils.GenerateToken(user.ID, s.conf.JWTSecretKey)
	if err != nil {
		return "", "", errors.New("could not generate token")
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, s.conf.JWTSecretKey)
	if err != nil {
		return "", "", errors.New("could not generate refresh token")
	}

	return token, refreshToken, nil
}

// New Access Token from Refresh Token
func (s *AuthService) RefreshToken(c context.Context, refresh_token string) (string, error) {
	// Check if refresh token already exists
	var existingToken model.Token
	if err := s.db.WithContext(c).Where("refresh_token = ?", refresh_token).First(&existingToken).Error; err != nil {
		return "", errors.New("invalid refresh token")
	}
	
	// Parse and validate the token
	token, err := jwt.Parse(refresh_token, func(token *jwt.Token) (interface{}, error) {
		// Validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.conf.JWTSecretKey), nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}

	// Get claims and extract user ID
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid jwt claims")
	}

	userId, ok := claims["user_id"].(string)
	if !ok {
		return "", errors.New("invalid user ID")
	}

	// Generate JWT new access token
	accessToken, err := utils.GenerateNewToken(userId, s.conf.JWTSecretKey)
	if err != nil {
		return "", errors.New("could not generate token")
	}

	return accessToken, nil
}

// CreateToken creates a new token.
func (s *AuthService) CreateToken(ctx context.Context, access_token, refreshToken string) error {
	// Parse and validate the token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		// Validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.conf.JWTSecretKey), nil
	})

	if err != nil || !token.Valid {
		return errors.New("invalid token")
	}

	// Get claims and extract user ID
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return errors.New("invalid jwt claims")
	}

	// Get user ID
	userId, ok := claims["user_id"].(string)
	if !ok {
		return errors.New("invalid user ID")
	}

	// Parse user ID to UUID
	userIdUUID, err := uuid.Parse(userId)
	if err != nil {
		return errors.New("invalid user ID")
	}

	// Check by user ID
	var existingToken model.Token
	if err := s.db.WithContext(ctx).Where("user_id = ?", userIdUUID).First(&existingToken).Error; err == nil {
		// Update token
		existingToken.Token = access_token
		if err := existingToken.Save(s.db); err != nil {
			return err
		}
		return errors.New("token has updated")
	}

	// Create new token
	tokenInput := model.Token{
		UserID: userIdUUID, 
		Token: access_token, 
		RefreshToken: refreshToken,
	}

	// Save in token database
	if err := tokenInput.Save(s.db); err != nil {
		return err
	}

	return nil
}
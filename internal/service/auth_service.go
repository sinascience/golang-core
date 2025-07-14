package service

import (
	"context"
	"errors"
	"venturo-core/configs"
	"venturo-core/internal/model"
	"venturo-core/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"


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
func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	// Find user by email
	var user model.User
	if err := s.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return "", errors.New("invalid credentials")
	}

	// Compare password with the hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Generate JWT
	token, err := utils.GenerateToken(user.ID, s.conf.JWTSecretKey)
	if err != nil {
		return "", errors.New("could not generate token")
	}

	return token, nil
}

// RefreshToken generates a new JWT for a user.
func (s *AuthService) RefreshToken(ctx context.Context, email string) (string, error) {
	// Find user by email
	var user model.User

	// Generate JWT
	accessToken, err := utils.GenerateRefreshToken(user.ID, s.conf.JWTSecretKey)
	if err != nil {
		return "", errors.New("could not generate token")
	}

	return accessToken, nil
}

// New Access Token from Refresh Token
func (s *AuthService) NewAccessToken(c *fiber.Ctx, secretKey string) (string, error) {
	if secretKey == "" {
		return "", c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid refresh token"})
	}
	
	// Parse and validate the token
	token, err := jwt.Parse(secretKey, func(token *jwt.Token) (interface{}, error) {
		// Validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "Unexpected signing method")
		}
		return []byte(secretKey), nil
	})

	if err != nil || !token.Valid {
		return "", c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired JWT"})
	}

	// Get claims and extract user ID
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid JWT claims"})
	}

	userId, ok := claims["user_id"].(string)
	if !ok {
		return "", c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID in token"})
	}

	// Generate JWT new access token
	accessToken, err := utils.GenerateNewToken(userId, s.conf.JWTSecretKey)
	if err != nil {
		return "", c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error":"could not generate token"})
	}

	return accessToken, nil
}

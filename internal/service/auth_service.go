package service

import (
	"context"
	"errors"
	"venturo-core/configs"
	"venturo-core/internal/model"
	"venturo-core/pkg/utils"

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

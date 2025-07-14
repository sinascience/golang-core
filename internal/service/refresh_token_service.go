package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log/slog"
	"venturo-core/configs"
	"venturo-core/internal/model"
	"venturo-core/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshTokenService struct {
	db   *gorm.DB
	conf *configs.Config
}

// NewAuthService creates a new auth service.
func NewRefreshTokenService(db *gorm.DB, conf *configs.Config) *RefreshTokenService {
	return &RefreshTokenService{db: db, conf: conf}
}

func (s *RefreshTokenService) Refresh(ctx context.Context, userID uuid.UUID, refreshToken string) (string, string, error) {
	var refreshModel model.RefreshToken

	err := refreshModel.FindByHashedToken(s.db, HashRefreshToken(refreshToken))
	if err != nil {
		return "", "", errors.New("refresh token is not exist")
	}
	token, err := utils.GenerateToken(userID, s.conf.JWTSecretKey)
	if err != nil {
		return "", "", errors.New("could not generate token")
	}

	refresh, err := utils.GenerateRefresh(userID, s.conf.JWTSecretKey)
	if err != nil {
		return "", "", errors.New("could not generate token")
	}

	err = refreshModel.Delete(s.db, refreshModel.ID)
	if err != nil {
		slog.Error("Failed to delete refresh token", "userId", userID, "error", err)
		return "", "", errors.New("could not refresh token")
	}

	return token, refresh, nil
}

func (s *RefreshTokenService) Store(ctx context.Context, userID uuid.UUID, refreshTokenRaw string) error {
	refreshToken := model.RefreshToken{
		UserID:      userID,
		HashedToken: HashRefreshToken(refreshTokenRaw),
	}
	// Store Refresh Token
	if err := refreshToken.Save(s.db); err != nil {
		slog.Error("Failed to save refresh token", "userId", userID, "error", err)
		return errors.New("Error Saving New Refresh Token")
	}

	return nil
}

func HashRefreshToken(token string) string {
	hash := sha256.New()
	hash.Write([]byte(token))
	return hex.EncodeToString(hash.Sum(nil))
}

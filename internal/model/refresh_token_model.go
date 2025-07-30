package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RefreshToken defines the refresh token model.
// This struct is used internally and is never sent in an API response.
type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:char(36);primary_key"`
	UserID    uuid.UUID `gorm:"type:char(36);not null"`
	Token     string    `gorm:"type:text;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
}

// BeforeCreate is a GORM hook that runs before a new record is created.
func (rt *RefreshToken) BeforeCreate(tx *gorm.DB) (err error) {
	rt.ID = uuid.New()
	return
}

// Create saves a new refresh token record to the database.
func (rt *RefreshToken) Create(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).Create(rt).Error
}

// FindAllByUserID finds all refresh tokens associated with a user.
func (rt *RefreshToken) FindAllByUserID(ctx context.Context, db *gorm.DB, userID uuid.UUID) ([]RefreshToken, error) {
	var tokens []RefreshToken
	err := db.WithContext(ctx).Where("user_id = ?", userID).Find(&tokens).Error
	return tokens, err
}

// FindByUserIDAndToken finds a specific refresh token by user ID and token hash.
func (rt *RefreshToken) FindByUserIDAndToken(ctx context.Context, db *gorm.DB, userID uuid.UUID, tokenHash string) (*RefreshToken, error) {
	var token RefreshToken
	err := db.WithContext(ctx).Where("user_id = ? AND token = ?", userID, tokenHash).First(&token).Error
	return &token, err
}

// Delete removes a refresh token record from the database by its ID.
func (rt *RefreshToken) Delete(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).Delete(rt).Error
}

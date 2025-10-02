package model 

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserPoint struct {
	ID uuid.UUID `gorm:"type:char(36);primary_key"`
	UserID uuid.UUID `gorm:"type:char(36);not null"`
	Point int32 `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// BeforeCreate is a GORM hook.
func (u *UserPoint) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}

// Save creates or updates a user record.
func (u *UserPoint) Save(db *gorm.DB) error {
	return db.WithContext(context.Background()).Save(u).Error
}

// FindByUserID retrieves a single user by their ID.
func (u *UserPoint) FindByUserID(db *gorm.DB, userID uuid.UUID) (*UserPoint, error) {
	var userPoint UserPoint
	err := db.WithContext(context.Background()).Where("user_id = ?", userID).First(&userPoint).Error
	return &userPoint, err
}
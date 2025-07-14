package model

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User defines the user model.
type RefreshToken struct {
	ID          uuid.UUID `gorm:"type:char(36);primary_key" json:"id"`
	UserID      uuid.UUID `gorm:"type:char(36);" json:"user_id"`
	HashedToken string    `gorm:"size:255;not null;unique" json:"hashed_token"`
}

// BeforeCreate is a GORM hook that runs before a new record is created.
func (u *RefreshToken) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}

func (u *RefreshToken) Save(db *gorm.DB) error {
	return db.WithContext(context.Background()).Save(u).Error
}
func (u *RefreshToken) Delete(db *gorm.DB, id uuid.UUID) error {
	return db.WithContext(context.Background()).Where("id = ?", id).Delete(&RefreshToken{}).Error
}

func (r *RefreshToken) FindByHashedToken(db *gorm.DB, hashedToken string) error {
	return db.Where("hashed_token = ?", hashedToken).First(r).Error
}

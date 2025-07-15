package model

import (
	"time"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Token defines the token model.
type Token struct {
	ID        		uuid.UUID `gorm:"type:char(36);primary_key" json:"id"`
	UserID    		uuid.UUID `gorm:"type:char(36);not null;unique" json:"user_id"`
	Token     		string    `gorm:"size:255;not null" json:"token"`
	RefreshToken	string    `gorm:"size:255;not null" json:"refresh_token"`
	CreatedAt 		time.Time `json:"created_at"`
	UpdatedAt 		time.Time `json:"updated_at"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// BeforeCreate is a GORM hook thath runs before a new record is created.
func (t *Token) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.New()
	return
}

// Save creates or updates a token record.
func (t *Token) Save(db *gorm.DB) error {
	return db.WithContext(context.Background()).Save(t).Error
}

// FindByRefreshToken is a custom finder method.
func (t *Token) FindByRefreshToken(db *gorm.DB, refreshToken string) (*Token, error) {
	var token Token
	err := db.WithContext(context.Background()).Where("refresh_token = ?", refreshToken).First(&token).Error
	return &token, err
}

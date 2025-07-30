package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User defines the user model.
type User struct {
	ID          uuid.UUID `gorm:"type:char(36);primary_key" json:"id"`
	Name        string    `gorm:"size:255;not null" json:"name"`
	Email       string    `gorm:"size:255;not null;unique" json:"email"`
	Password    string    `gorm:"size:255;not null" json:"-"`
	AvatarURL   string    `gorm:"size:255;null" json:"avatar_url,omitempty"`              // New field
	ImageStatus string    `gorm:"size:20;not null;default:'default'" json:"image_status"` // New field
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// BeforeCreate is a GORM hook that runs before a new record is created.
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}

// --- Standard CRUD Methods ---

// Save creates or updates a user record.
func (u *User) Save(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).Save(u).Error
}

// FindAll retrieves all user records.
func (u *User) FindAll(ctx context.Context, db *gorm.DB) ([]User, error) {
	var users []User
	err := db.WithContext(ctx).Find(&users).Error
	return users, err
}

// FindByID retrieves a single user by their ID.
func (u *User) FindByID(ctx context.Context, db *gorm.DB, id uuid.UUID) (*User, error) {
	var user User
	err := db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	return &user, err
}

// Delete removes a user record by their ID.
func (u *User) Delete(ctx context.Context, db *gorm.DB, id uuid.UUID) error {
	return db.WithContext(ctx).Where("id = ?", id).Delete(&User{}).Error
}

// FindByEmail is a custom finder method.
func (u *User) FindByEmail(ctx context.Context, db *gorm.DB, email string) (*User, error) {
	var user User
	err := db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return &user, err
}

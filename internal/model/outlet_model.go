package model

import (
	"gorm.io/gorm"
	"time"
	"context"
	"github.com/google/uuid"
)

// Outlet defines the outlet model.
type Outlet struct {
	ID        uuid.UUID `gorm:"type:char(36);primary_key" json:"id"`
	Name      string    `gorm:"size:255;not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time	`json:"updated_at"`

	Transactions []Transaction `gorm:"foreignKey:OutletID"`
}

// BeforeCreate is a GORM hook that runs before a new record is created.
func (outlet *Outlet) BeforeCreate(tx *gorm.DB) (err error) {
	outlet.ID = uuid.New()
	return nil
}

// Save creates or updates an outlet record.
func (outlet *Outlet) Save(db *gorm.DB) error {
	return db.WithContext(context.Background()).Save(outlet).Error
}
package model

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Transaction defines the main transaction model.
type Outlet struct {
	ID          uuid.UUID     `gorm:"type:char(36);primary_key"`
	Name        string        `gorm:"size:255;not null"`
	Transaction []Transaction `gorm:"foreignKey:OutletID"`
}

// BeforeCreate is a GORM hook.
func (t *Outlet) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.New()
	return
}

// Save creates or updates a record.
func (t *Outlet) Save(db *gorm.DB) error {
	return db.WithContext(context.Background()).Save(t).Error
}

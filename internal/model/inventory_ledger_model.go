package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Transaction defines the main transaction model.
type InventoryLedger struct {
	ID             uuid.UUID  `gorm:"type:char(36);primary_key"`
	ItemID         uuid.UUID  `gorm:"type:char(36);not null"`
	OutletID       uuid.UUID  `gorm:"type:char(36);not null"`
	TransactionID  *uuid.UUID `gorm:"type:char(36);"`
	QuantityChange int8       `gorm:"not null"`
	CreatedAt      time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP;<-:create"`
	UpdatedAt      time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP;autoUpdateTime"`

	// Relationships
	Item        Product     `gorm:"foreignKey:ItemID"`
	Transaction Transaction `gorm:"foreignKey:TransactionID"`
	Outlet      Outlet      `gorm:"foreignKey:OutletID"`
}

// BeforeCreate is a GORM hook.
func (t *InventoryLedger) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.New()
	return
}

// Save creates or updates a record.
func (t *InventoryLedger) Save(db *gorm.DB) error {
	return db.WithContext(context.Background()).Save(t).Error
}

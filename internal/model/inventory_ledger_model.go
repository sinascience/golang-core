package model

import (
	"context"
	"time"
	"gorm.io/gorm"
	"github.com/google/uuid"
)

// Inventory defines the inventory model
type Inventory struct {
	ID         uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	ProductID  uuid.UUID `gorm:"type:char(36);not null" json:"product_id"`
	OutletID   uuid.UUID `gorm:"type:char(36);not null" json:"outlet_id"`
	TransactionID *uuid.UUID `gorm:"type:char(36)" json:"transaction_id"`
	Quantity   int       `gorm:"not null" json:"quantity"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	Outlet     Outlet  `gorm:"foreignKey:OutletID"`
	Product    Product    `gorm:"foreignKey:product_id"`
	Transaction Transaction `gorm:"foreignKey:transaction_id"`
}

// Override default table name
func (Inventory) TableName() string {
	return "inventory_ledgers"
}

// BeforeCreate is a GORM hook.
func (i *Inventory) BeforeCreate(tx *gorm.DB) (err error){
	i.ID = uuid.New()
	return
}

// Save creates or updates an inventory record.
func (i *Inventory) Save(db *gorm.DB) error {
	return db.WithContext(context.Background()).Save(i).Error
}
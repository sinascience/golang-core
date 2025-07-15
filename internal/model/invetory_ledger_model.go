package model

import (
	"time"

	"github.com/google/uuid"
)

type InventoryLedger struct {
	ID         uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	ItemID	   uuid.UUID `gorm:"type:char(36);not null" json:"item_id"`
	OutletID   uuid.UUID `gorm:"type:char(36);not null" json:"outlet_id"`
	TransactionID uuid.UUID `gorm:"type:char(36);null" json:"transaction_id"`
	Quantity   int       `gorm:"not null" json:"quantity"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	Outlet  Outlet  `gorm:"foreignKey:OutletID"`
	// Item    Item    `gorm:"foreignKey:ItemID"`
	Transaction Transaction `gorm:"foreignKey:TransactionID"`
}

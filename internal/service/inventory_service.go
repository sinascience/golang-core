package service

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"sync"
	"context"
	"venturo-core/internal/model"
)

// InventoryService handles the business logic for inventory.
type InventoryService struct {
	db *gorm.DB
	wg *sync.WaitGroup
}

// NewInventoryService creates a new inventory service
func NewInventoryService(db *gorm.DB, wg *sync.WaitGroup) *InventoryService {
	return &InventoryService{db: db, wg: wg}
}

// Create Stock In method for inventory
func (i *InventoryService) StockIn(ctx context.Context, quantityIn int, product_id, outled_id uuid.UUID, transaction_id *uuid.UUID) (*model.Inventory, error){	
	// Create the inventory
	inventoryInput := model.Inventory{
		ProductID : product_id,
		OutletID : outled_id,
		TransactionID : transaction_id,
		Quantity : quantityIn,
	}

	// Save the inventory
	if err := inventoryInput.Save(i.db); err != nil {
		return nil,err
	}

	return &inventoryInput, nil
}

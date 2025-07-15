package service

import (
	"context"
	"venturo-core/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InventoryService struct {
	db *gorm.DB
}

type CreateInventoryInput struct {
	ItemID        uuid.UUID
	OutletID      uuid.UUID
	TransactionID uuid.UUID
	Quantity      int8
}

// NewPostService creates a new post service.
func NewInventoryService(db *gorm.DB) *InventoryService {
	return &InventoryService{db: db}
}

func (s *InventoryService) StockIN(ctx context.Context, input CreateInventoryInput) (*model.InventoryLedger, error) {
	stock := model.InventoryLedger{
		ItemID:         input.ItemID,
		OutletID:       input.OutletID,
		QuantityChange: input.Quantity,
	}

	if err := stock.Save(s.db); err != nil {
		return nil, err
	}
	return &stock, nil
}

func (s *InventoryService) StockOUT(ctx context.Context, input CreateInventoryInput) (*model.InventoryLedger, error) {
	stock := model.InventoryLedger{
		ItemID:         input.ItemID,
		OutletID:       input.OutletID,
		TransactionID:  &input.TransactionID,
		QuantityChange: input.Quantity,
	}

	if err := stock.Save(s.db); err != nil {
		return nil, err
	}
	return &stock, nil
}

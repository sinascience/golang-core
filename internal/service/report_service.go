package service

import (
	"context"
	"gorm.io/gorm"
	"github.com/google/uuid"
	"time"
)

// ReportService handles the business logic for reports.
type ReportService struct {
	db *gorm.DB
}

// NewReportService creats a new ReportService.
func NewReportService(db *gorm.DB) *ReportService {
	return &ReportService{db: db}
}

// Report is the struct for a report.
type InventoryReport struct {
	ProductID  uuid.UUID `json:"product_id"`
	ProductName string    `json:"product_name"`
	OutletID   uuid.UUID `json:"outlet_id"`
	OutletName string    `json:"outlet_name"`
	Quantity   int64     `json:"on_hand_quantity"`
	History    []InventoryHistory `json:"history_transaction" gorm:"-"`
}

// InventoryHistory is the struct for a transaction report.
type InventoryHistory struct {
	TransactionID uuid.UUID `json:"transaction_id"`
	TransactionDate time.Time `json:"transaction_date"`
}

// Show all and groupingreport by outlet_id and product_id
func (s *ReportService) Show(ctx context.Context) ([]InventoryReport, error) {
	var report []InventoryReport

	// Get report
	if err := s.db.WithContext(ctx).
		Table("inventory_ledgers").
		Select("inventory_ledgers.outlet_id, o.name as outlet_name, inventory_ledgers.product_id, p.name as product_name, SUM(quantity) as quantity").
		Joins("LEFT JOIN outlets o ON o.id = inventory_ledgers.outlet_id").
		Joins("LEFT JOIN products p ON p.id = inventory_ledgers.product_id").
		Group("inventory_ledgers.outlet_id, o.name, inventory_ledgers.product_id, p.name").
		Scan(&report).Error; err != nil {
		return nil, err
	}

	// Get transaction history
	for i := range report {
		var history []InventoryHistory
		if err := s.db.WithContext(ctx).
			Table("inventory_ledgers").
			Select("transaction_id, t.created_at as transaction_date").
			Joins("LEFT JOIN transactions t ON t.id = transaction_id").
			Where("inventory_ledgers.product_id = ? AND inventory_ledgers.outlet_id = ?", report[i].ProductID, report[i].OutletID).
			Scan(&history).Error; err != nil {
			return nil, err
		}
		report[i].History = history
	}

	return report, nil
}

package service

import (
	"context"
	"time"
	"venturo-core/configs"
	"venturo-core/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReportService struct {
	db   *gorm.DB
	conf *configs.Config
}

func NewReportService(db *gorm.DB, conf *configs.Config) *ReportService {
	return &ReportService{db: db, conf: conf}
}

type InventoryReport struct {
	ItemID      uuid.UUID          `json:"item_id"`
	ProductName string             `json:"product_name"`
	Quantity    int64              `json:"quantity"`
	History     []InventoryHistory `json:"history"`
}

type InventoryHistory struct {
	OutletID       uuid.UUID `json:"outlet_id"`
	OutletName     string    `json:"outlet_name"`
	QuantityChange int8      `json:"quantity_change"`
	InvoiceCode    *string   `json:"invoice_code"`
	Date           time.Time `json:"date"`
}

func (s *ReportService) Report(ctx context.Context, itemID, outletID *uuid.UUID) ([]InventoryReport, error) {
	var ledgers []model.InventoryLedger
	query := s.db.
		Select("id", "item_id", "outlet_id", "transaction_id", "quantity_change", "created_at").
		Preload("Item", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name")
		}).
		Preload("Outlet", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name")
		}).
		Preload("Transaction", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "invoice_code")
		})

	if itemID != nil {
		query = query.Where("item_id = ?", *itemID)
	}
	if outletID != nil {
		query = query.Where("outlet_id = ?", *outletID)
	}

	err := query.Find(&ledgers).Error
	if err != nil {
		return nil, err
	}

	reportMap := make(map[uuid.UUID]*InventoryReport)
	for _, row := range ledgers {
		report, ok := reportMap[row.ItemID]
		if !ok {
			report = &InventoryReport{
				ItemID:      row.ItemID,
				ProductName: row.Item.Name,
				Quantity:    0,
				History:     []InventoryHistory{},
			}
			reportMap[row.ItemID] = report
		}

		var invoiceCode string
		if row.Transaction.ID != uuid.Nil {
			invoiceCode = row.Transaction.InvoiceCode
		} else {
			invoiceCode = "INV-STOCK-IN"
		}
		report.Quantity += int64(row.QuantityChange)
		report.History = append(report.History, InventoryHistory{
			OutletID:       row.OutletID,
			OutletName:     row.Outlet.Name,
			QuantityChange: int8(row.QuantityChange),
			InvoiceCode:    &invoiceCode,
			Date:           row.CreatedAt,
		})
	}

	var finalReports []InventoryReport
	for _, r := range reportMap {
		finalReports = append(finalReports, *r)
	}
	return finalReports, nil
}

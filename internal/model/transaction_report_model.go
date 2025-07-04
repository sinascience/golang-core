package model

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

// CategorySummary defines the structure for our JSON data.
type CategorySummary map[string]int64

// Value makes CategorySummary implement the driver.Valuer interface.
// This is how GORM knows how to save our map into a JSON column.
func (cs CategorySummary) Value() (driver.Value, error) {
	return json.Marshal(cs)
}

// Scan makes CategorySummary implement the sql.Scanner interface.
// This is how GORM knows how to read data from a JSON column into our map.
func (cs *CategorySummary) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &cs)
}

// TransactionReport defines the model for our single-row reporting table.
type TransactionReport struct {
	ID                      uint8           `gorm:"primary_key"`
	TotalRevenue            uint64
	TotalPaidTransactions   uint64
	TotalProductsSold       uint64
	TotalUniqueCustomers    uint64
	CategorySummary         CategorySummary `gorm:"type:json"`
	UpdatedAt               time.Time
}

// Save updates the single report row.
func (tr *TransactionReport) Save(db *gorm.DB) error {
	return db.WithContext(context.Background()).Save(tr).Error
}
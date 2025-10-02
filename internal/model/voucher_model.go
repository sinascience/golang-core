package model

import (
	"context"
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Voucher struct {
	ID uuid.UUID `gorm:"type:char(36);primary_key"`
	Code string `gorm:"size:20;not null;unique"`
	Discount int32 `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// BeforeCreate is a GORM hook.
func (v *Voucher) BeforeCreate(tx *gorm.DB) (err error) {
	v.ID = uuid.New()
	return
}

// Save creates or updates a voucher record.
func (v *Voucher) Save(db *gorm.DB) error {
	return db.WithContext(context.Background()).Save(v).Error
}

// FindByCode retrieves a single voucher by their code.
func (v *Voucher) FindByCode(db *gorm.DB, code string) (*Voucher, error) {
	var voucher Voucher
	return &voucher, db.WithContext(context.Background()).Where("code = ?", code).First(&voucher).Error
}
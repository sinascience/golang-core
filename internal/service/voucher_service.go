package service

import (
	"venturo-core/internal/model"
	"gorm.io/gorm"
)

type VoucherService struct {
	db *gorm.DB
}

// NewVoucherService creates a new voucher service.
func NewVoucherService(db *gorm.DB) *VoucherService {
	return &VoucherService{db: db}
}

func (v *VoucherService) CreateVoucher(discount int32, code string) (*model.Voucher, error) {
	voucher := model.Voucher{Discount: discount, Code: code}
	if err := voucher.Save(v.db); err != nil {
		return nil, err
	}
	return &voucher, nil
}

func (v *VoucherService) GetVoucherByCode(code string) (*model.Voucher, error) {
	voucher, err := (&model.Voucher{}).FindByCode(v.db, code)
	if err != nil {
		return nil, err
	}
	return voucher, nil
}
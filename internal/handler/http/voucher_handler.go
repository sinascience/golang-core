package http

import (
	"venturo-core/internal/service"
	"venturo-core/pkg/response"
	"github.com/gofiber/fiber/v2"
)

type VoucherHandler struct {
	voucherService *service.VoucherService
}

type VoucherInput struct {
	Discount int32 `json:"discount" validate:"required"`
	Code string `json:"code" validate:"required"`
}

func NewVoucherHandler(voucherService *service.VoucherService) *VoucherHandler {
	return &VoucherHandler{voucherService: voucherService}
}

func (v *VoucherHandler) CreateVoucher(c *fiber.Ctx) error {
	var input VoucherInput
	if err := c.BodyParser(&input); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err)
	}

	voucher, err := v.voucherService.CreateVoucher(input.Discount, input.Code)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err)
	}
	return response.Success(c, fiber.StatusCreated, voucher)
}

func (v *VoucherHandler) GetVoucherByCode(c *fiber.Ctx) error {
	var code string = c.FormValue("code")
	voucher, err := v.voucherService.GetVoucherByCode(code)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err)
	}
	return response.Success(c, fiber.StatusOK, voucher)
}
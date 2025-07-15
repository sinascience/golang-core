package http

import (
	"errors"
	"venturo-core/internal/service"
	"venturo-core/pkg/response"
	"venturo-core/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type InventoryHandler struct {
	inventoryService *service.InventoryService
}

// NewTransactionHandler creates a new handler.
func NewInventoryHandler(s *service.InventoryService) *InventoryHandler {
	return &InventoryHandler{inventoryService: s}
}

// CreateInventoryPayload defines the expected JSON from the client.
type CreateInventoryPayload struct {
	ItemID        uuid.UUID `json:"item_id" validate:"required"`
	OutletID      uuid.UUID `json:"outlet_id" validate:"required"`
	TransactionID uuid.UUID `json:"transaction_id"`
	Quantity      int8      `json:"quantity" validate:"required"`
}

func (h *InventoryHandler) parseCreateInventoryPayload(c *fiber.Ctx) (*service.CreateInventoryInput, error) {
	_, ok := c.Locals("current_user_id").(uuid.UUID)
	if !ok {
		return nil, errors.New("unauthorized")
	}

	payload := new(CreateInventoryPayload)
	if err := c.BodyParser(payload); err != nil {
		return nil, errors.New("invalid JSON")
	}

	if errs := validator.ValidateStruct(payload); errs != nil {
		return nil, response.ValidationError(c, errs)
	}

	return &service.CreateInventoryInput{
		ItemID:        payload.ItemID,
		OutletID:      payload.OutletID,
		TransactionID: payload.TransactionID,
		Quantity:      payload.Quantity,
	}, nil
}

// StockIN handler
func (h *InventoryHandler) StockIN(c *fiber.Ctx) error {
	input, err := h.parseCreateInventoryPayload(c)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err)
	}

	result, err := h.inventoryService.StockIN(c.Context(), *input)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err)
	}

	return response.Success(c, fiber.StatusCreated, result)
}

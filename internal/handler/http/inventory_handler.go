package http

import (
	"errors"
	"venturo-core/internal/service"
	"venturo-core/pkg/response"
	"venturo-core/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// InventoryHandler handles inventory-related HTTP requests.
type InventoryHandler struct {
	inventoryService *service.InventoryService
}

// NewInventoryHandler creates a new inventory handler.
func NewInventoryHandler(inventoryService *service.InventoryService) *InventoryHandler {
	return &InventoryHandler{inventoryService: inventoryService}
}

// InventoryPayload defines the payload for creating an inventory item.
type InventoryPayload struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	OutletID  uuid.UUID `json:"outlet_id" validate:"required"`
	TransactionID *uuid.UUID `json:"transaction_id"`
	Quantity int `json:"quantity" validate:"required"` // The quantity of the inventory item.
}

// StockIn is a handler for stock in operation
func (i *InventoryHandler) StockIn(ctx *fiber.Ctx) error {
	payload := new(InventoryPayload)
	// Parse the request body
	if err := ctx.BodyParser(payload); err != nil {
		return response.Error(ctx, fiber.StatusBadRequest, errors.New("cannot parse JSON"))
	}

	// Replace manual checks with a single call to the validator
	if err := validator.ValidateStruct(payload); err != nil {
		return response.ValidationError(ctx, err)
	}

	// Call the service to stock in the inventory
	stockIn, err := i.inventoryService.StockIn(ctx.Context(), payload.Quantity, payload.ProductID, payload.OutletID, payload.TransactionID)
	if err != nil {
		return response.Error(ctx, fiber.StatusInternalServerError, err)
	}

	return response.Success(ctx, fiber.StatusCreated, stockIn)
}
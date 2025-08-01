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
/**
@api {post} /api/v1/inventories Stock In
@apiName StockIn
@apiGroup Inventory
@apiDescription Stocks in a specific quantity of a product into a specific outlet.

@apiHeader {string} Authorization "Bearee <access_token>"

@apiBody {UUID} product_id The ID of the product to stock in.
@apiBody {UUID} outlet_id The ID of the outlet to stock in the product into.
@apiBody {Number} quantity The quantity of the product to stock in.

@apiSuccessExample {json} Success-Response:
{
    "status_code": 201,
    "data": {
        "id": "4a83054d-fc50-4aef-9584-8341acbf610a",
        "product_id": "b53d6068-62d2-11f0-9182-e120fc92e2ba",
        "outlet_id": "32b4ba07-61ee-11f0-b1d6-9c12216b8fcb",
        "transaction_id": null,
        "quantity": 30,
        "created_at": "2025-08-01T14:03:37.121+07:00",
        "updated_at": "2025-08-01T14:03:37.121+07:00",
        "Outlet": {
            "id": "00000000-0000-0000-0000-000000000000",
            "name": "",
            "created_at": "0001-01-01T00:00:00Z",
            "updated_at": "0001-01-01T00:00:00Z",
            "Transactions": null
        },
        "Product": {
            "ID": "00000000-0000-0000-0000-000000000000",
            "Name": "",
            "Price": 0,
            "Stock": 0,
            "ImageURL": "",
            "ImageStatus": "",
            "CreatedAt": "0001-01-01T00:00:00Z",
            "UpdatedAt": "0001-01-01T00:00:00Z"
        },
        "Transaction": {
            "ID": "00000000-0000-0000-0000-000000000000",
            "UserID": "00000000-0000-0000-0000-000000000000",
            "OutletID": "00000000-0000-0000-0000-000000000000",
            "InvoiceCode": "",
            "Total": 0,
            "is_paid": null,
            "Note": "",
            "CreatedAt": "0001-01-01T00:00:00Z",
            "UpdatedAt": "0001-01-01T00:00:00Z",
            "User": {
                "id": "00000000-0000-0000-0000-000000000000",
                "name": "",
                "email": "",
                "image_status": "",
                "created_at": "0001-01-01T00:00:00Z",
                "updated_at": "0001-01-01T00:00:00Z"
            },
            "TransactionDetails": null,
            "Outlet": {
                "id": "00000000-0000-0000-0000-000000000000",
                "name": "",
                "created_at": "0001-01-01T00:00:00Z",
                "updated_at": "0001-01-01T00:00:00Z",
                "Transactions": null
            }
        }
    }
}
*/
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
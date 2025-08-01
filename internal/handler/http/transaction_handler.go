package http

import (
	"errors"
	"strings"
	"venturo-core/internal/model"
	"venturo-core/internal/service"
	"venturo-core/pkg/response"
	"venturo-core/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type TransactionHandler struct {
	transactionService *service.TransactionService
}

// NewTransactionHandler creates a new handler.
func NewTransactionHandler(s *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{transactionService: s}
}

// CreateTransactionPayload defines the expected JSON from the client.
type CreateTransactionPayload struct {
	OutletID uuid.UUID `json:"outlet_id" validate:"required"`
	Items []struct {
		ProductID   uuid.UUID `json:"product_id" validate:"required"`
		ProductName string    `json:"product_name" validate:"required"`
		Category    uint8     `json:"category" validate:"required,min=1,max=3"`
		Qty         int8      `json:"qty" validate:"required,min=1"`
		Price       int32     `json:"price" validate:"required,min=0"`
	} `json:"items" validate:"required,min=1"`
	Note string `json:"note"`
}

// CreateTransaction is the handler for creating a new transaction.
// @Summary      Create a new transaction
// @Description  Creates a transaction with multiple detail items for the authenticated user.
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        payload  body      CreateTransactionPayload  true  "Transaction Payload"
// @Success      201      {object}  response.ApiResponse{data=model.Transaction} "Successfully created transaction"
// @Failure      400      {object}  response.ApiResponse "Bad Request"
// @Failure      401      {object}  response.ApiResponse "Unauthorized"
// @Router       /transactions [post]

/**
@api {post} /transactions Create a Transaction
@apiName CreateTransaction
@apiGroup Transactions
@apiDescription Creates a transaction with multiple detail items for the authenticated user.

@apiHeader {string} Authorization "Bearer <access_token>"

@apiBody {UUID} outlet_id The ID of the outlet
@apiBody {Object[]} items An array of items
@apiBody {UUID} items.product_id The ID of the product
@apiBody {string} items.product_name The name of the product
@apiBody {number} items.category The category of the product
@apiBody {number} items.qty The quantity of the product
@apiBody {number} items.price The price of the product
@apiBody {string} [note] A note about the transaction

@apiSuccess {UUID} id The ID of the created transaction
@apiSuccess {UUID} outlet_id The ID of the outlet
@apiSuccess {Object[]} items An array of items
@apiSuccess {UUID} items.product_id The ID of the product
@apiSuccess {string} items.product_name The name of the product
@apiSuccess {number} items.category The category of the product
@apiSuccess {number} items.qty The quantity of the product
@apiSuccess {number} items.price The price of the product
@apiSuccess {string} [note] A note about the transaction

@apiSuccessExample {json} Success-Response:
{
	"status_code": 201,
    "data": {
        "ID": "6684001c-c4f0-43d2-8749-f4622a8ad931",
        "UserID": "5b5aee18-95c3-436c-b8b9-b66fe82c1480",
        "OutletID": "32b4ba07-61ee-11f0-b1d6-9c12216b8fcb",
        "InvoiceCode": "INV-2025-6529",
        "Total": 15000,
        "is_paid": false,
        "Note": "INV INV-2025-6529 includes: BUKU. Additional notes: Pesanan untuk pelanggan reguler",
        "CreatedAt": "2025-08-01T13:48:28.471+07:00",
        "UpdatedAt": "2025-08-01T13:48:28.471+07:00",
        "User": {
            "id": "00000000-0000-0000-0000-000000000000",
            "name": "",
            "email": "",
            "image_status": "",
            "created_at": "0001-01-01T00:00:00Z",
            "updated_at": "0001-01-01T00:00:00Z"
        },
        "TransactionDetails": [
            {
                "ID": "5b719b68-1450-4c1e-8578-702322d1d479",
                "TransactionID": "6684001c-c4f0-43d2-8749-f4622a8ad931",
                "ProductID": "30d7404c-47cb-440b-be33-447e10211b41",
                "ProductName": "BUKU",
                "Category": 1,
                "Qty": 1,
                "Price": 15000,
                "CreatedAt": "2025-08-01T13:48:28.487+07:00",
                "UpdatedAt": "2025-08-01T13:48:28.487+07:00"
            }
        ],
        "Outlet": {
            "id": "00000000-0000-0000-0000-000000000000",
            "name": "",
            "created_at": "0001-01-01T00:00:00Z",
            "updated_at": "0001-01-01T00:00:00Z",
            "Transactions": null
        }
    }
}
*/
func (h *TransactionHandler) CreateTransaction(c *fiber.Ctx) error {
	userID, ok := c.Locals("current_user_id").(uuid.UUID)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, errors.New("unauthorized"))
	}

	payload := new(CreateTransactionPayload)
	if err := c.BodyParser(payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, errors.New("cannot parse JSON"))
	}

	if errs := validator.ValidateStruct(payload); errs != nil {
		return response.ValidationError(c, errs)
	}

	// Map payload to the service input struct
	serviceInput := service.CreateTransactionInput{
		UserID: userID,
		OutletID: payload.OutletID,
		Note:   payload.Note,
	}
	for _, item := range payload.Items {
		serviceInput.Items = append(serviceInput.Items, struct {
			ProductID   uuid.UUID
			ProductName string
			Category    model.ProductCategory
			Qty         int8
			Price       int32
		}{
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Category:    model.ProductCategory(item.Category),
			Qty:         item.Qty,
			Price:       item.Price,
		})
	}

	// Call the service
	transaction, err := h.transactionService.CreateTransaction(c.Context(), serviceInput)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err)
	}

	return response.Success(c, fiber.StatusCreated, transaction)
}

// MarkAsPaid handles the request to mark a transaction as paid.
// @Summary      Pay for a Transaction
// @Description  Marks a transaction as paid and triggers a background report update.
// @Tags         Transactions
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      string  true  "Transaction ID"
// @Success      200  {object}  response.ApiResponse "Successfully paid"
// @Failure      401  {object}  response.ApiResponse "Unauthorized"
// @Failure      404  {object}  response.ApiResponse "Transaction not found"
// @Router       /transactions/{id}/pay [post]

/**
@api {post} /transactions/:id/pay Pay for a Transaction
@apiName PayForTransaction
@apiGroup Transactions
@apiParam {uuid} id Transaction ID
@apiHeader {string} Authorization "Bearer <access_token>"
@apiSuccess {string} message "Transaction marked as paid. Report is updating."

@apiSuccessExample {json} Success-Response:
{
	"status_code" : 200,
	"message" : "Transaction marked as paid. Report is updating."
}
*/
func (h *TransactionHandler) MarkAsPaid(c *fiber.Ctx) error {
	idParam := c.Params("id")
	transactionID, err := uuid.Parse(idParam)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, errors.New("invalid ID format"))
	}

	err = h.transactionService.MarkAsPaid(c.Context(), transactionID)
	if err != nil {
		// Differentiate between not found and other errors
		if strings.Contains(err.Error(), "not found") {
			return response.Error(c, fiber.StatusNotFound, err)
		}
		return response.Error(c, fiber.StatusInternalServerError, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"message": "Transaction marked as paid. Report is updating."})
}
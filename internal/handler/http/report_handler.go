package http

import (
	"github.com/gofiber/fiber/v2"
	"venturo-core/internal/service"
)

// ReportHandler
type ReportHandler struct {
	reportService *service.ReportService
}

// NewReportHandler creates a new handler
func NewReportHandler(rs *service.ReportService) *ReportHandler {
	return &ReportHandler{reportService: rs}
}

// GetInventoryReport is a handler for invenrotry report
/**
@api {get} /api/v1/reports Get Inventory Report
@apiName GetInventoryReport
@apiGroup Inventory
@apiDescription Get Inventory Report

@apiHeader {string} Authorization "Bearer <access_token>"

@apiSuccessExample {json} Success-Response:
[	
	{
        "product_id": "30d7404c-47cb-440b-be33-447e10211b41",
        "product_name": "Buku",
        "outlet_id": "32b4ba07-61ee-11f0-b1d6-9c12216b8fcb",
        "outlet_name": "outlet-1",
        "on_hand_quantity": 2,
        "history_transaction": [
            {
                "transaction_id": "00000000-0000-0000-0000-000000000000",
                "transaction_date": "0001-01-01T00:00:00Z"
            },
            {
                "transaction_id": "00000000-0000-0000-0000-000000000000",
                "transaction_date": "0001-01-01T00:00:00Z"
            },
            {
                "transaction_id": "89d2bf47-ef68-42ef-96a5-07de4c400752",
                "transaction_date": "2025-07-17T09:14:13+07:00"
            },
            {
                "transaction_id": "6684001c-c4f0-43d2-8749-f4622a8ad931",
                "transaction_date": "2025-08-01T13:48:28+07:00"
            },
            {
                "transaction_id": "63e3c795-d95f-48dd-8f86-564dff955f23",
                "transaction_date": "2025-07-30T15:50:16+07:00"
            },
            {
                "transaction_id": "00000000-0000-0000-0000-000000000000",
                "transaction_date": "0001-01-01T00:00:00Z"
            },
            {
                "transaction_id": "89d2bf47-ef68-42ef-96a5-07de4c400752",
                "transaction_date": "2025-07-17T09:14:13+07:00"
            }
        ]
    }
]

*/
func (r *ReportHandler) GetInventoryReport(c *fiber.Ctx) error {
	// Get the report
	report, err := r.reportService.Show(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(report)
}
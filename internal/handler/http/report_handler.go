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
func (r *ReportHandler) GetInventoryReport(c *fiber.Ctx) error {
	// Get the report
	report, err := r.reportService.Show(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(report)
}
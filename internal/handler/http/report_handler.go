package http

import (
	"errors"
	"venturo-core/internal/service"
	"venturo-core/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ReportHandler struct {
	reportService *service.ReportService
}

func NewReportHandler(s *service.ReportService) *ReportHandler {
	return &ReportHandler{reportService: s}
}

func (h *ReportHandler) GetReport(c *fiber.Ctx) error {
	itemIDStr := c.Query("item_id")
	outletIDStr := c.Query("outlet_id")

	var itemID, outletID *uuid.UUID

	if itemIDStr != "" {
		id, err := uuid.Parse(itemIDStr)
		if err == nil {
			itemID = &id
		}
	}

	if outletIDStr != "" {
		id, err := uuid.Parse(outletIDStr)
		if err == nil {
			outletID = &id
		}
	}
	report, err := h.reportService.Report(c.Context(), itemID, outletID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, errors.New("Internal Server Error"))
	}
	return response.Success(c, fiber.StatusOK, report)
}

package response

import (
	"math"

	"github.com/gofiber/fiber/v2"
)

// ApiResponse is the standard structure for all API responses.
type ApiResponse struct {
	StatusCode int         `json:"status_code"`
	Data       interface{} `json:"data,omitempty"`
	Meta       *Meta       `json:"meta,omitempty"`
	Errors     interface{} `json:"errors,omitempty"`
}

// Meta holds the pagination metadata.
type Meta struct {
	TotalRecords int64 `json:"total_records"`
	CurrentPage  int   `json:"current_page"`
	PerPage      int   `json:"per_page"`
	TotalPages   int   `json:"total_pages"`
}

// Success sends a standard success response.
func Success(c *fiber.Ctx, statusCode int, data interface{}) error {
	return c.Status(statusCode).JSON(ApiResponse{
		StatusCode: statusCode,
		Data:       data,
	})
}

// Pagination sends a standard paginated response.
func Pagination(c *fiber.Ctx, data interface{}, page, limit int, total int64) error {
	meta := &Meta{
		TotalRecords: total,
		CurrentPage:  page,
		PerPage:      limit,
		TotalPages:   int(math.Ceil(float64(total) / float64(limit))),
	}
	return c.Status(fiber.StatusOK).JSON(ApiResponse{
		StatusCode: fiber.StatusOK,
		Data:       data,
		Meta:       meta,
	})
}

// Error sends a standard error response.
func Error(c *fiber.Ctx, statusCode int, err error) error {
	return c.Status(statusCode).JSON(ApiResponse{
		StatusCode: statusCode,
		Errors:     err.Error(),
	})
}

// ValidationError sends a structured validation error response.
func ValidationError(c *fiber.Ctx, errors map[string]string) error {
	return c.Status(fiber.StatusBadRequest).JSON(ApiResponse{
		StatusCode: fiber.StatusBadRequest,
		Errors:     errors,
	})
}

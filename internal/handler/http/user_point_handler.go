package http

import (
	"venturo-core/internal/service"
	"venturo-core/pkg/response"
	// "strconv"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserPointHandler struct {
	userPointService *service.UserPointService
}

type CreateUserPointInput struct {
	Point int32 `json:"point"`
}

// NewUserPointHandler creates a new user point handler.
func NewUserPointHandler(userPointService *service.UserPointService) *UserPointHandler {
	return &UserPointHandler{userPointService: userPointService}
}

func (h *UserPointHandler) GetUserPointByUserID(c *fiber.Ctx) error {
	userID, ok := c.Locals("current_user_id").(uuid.UUID)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, errors.New("unauthorized"))
	}

	userPoint, err := h.userPointService.FindPointByUserID(userID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err)
	}
	return response.Success(c, fiber.StatusOK, userPoint)
}

func (h *UserPointHandler) CreateUserPoint(c *fiber.Ctx) error {
	userID, ok := c.Locals("current_user_id").(uuid.UUID)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, errors.New("unauthorized"))
	}

	userPoint, err := h.userPointService.CreateUserPoint(userID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err)
	}

	return response.Success(c, fiber.StatusCreated, userPoint)
}

func (h *UserPointHandler) UpdateUserPoint(c *fiber.Ctx) error {
    pointID, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return response.Error(c, fiber.StatusBadRequest, errors.New("invalid ID format"))
    }

    userID, ok := c.Locals("current_user_id").(uuid.UUID)
    if !ok {
        return response.Error(c, fiber.StatusUnauthorized, errors.New("unauthorized"))
    }

    var input CreateUserPointInput
    if err := c.BodyParser(&input); err != nil {
        return response.Error(c, fiber.StatusBadRequest, err)
    }

    updatedPoint, err := h.userPointService.UpdateUserPoint(userID, pointID, input.Point)
    if err != nil {
        return response.Error(c, fiber.StatusInternalServerError, err)
    }

    return response.Success(c, fiber.StatusOK, updatedPoint)
}

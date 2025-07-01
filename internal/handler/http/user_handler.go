package http

import (
	"errors"
	"venturo-core/internal/service"
	"venturo-core/pkg/response"
	"venturo-core/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetProfile is the handler for the get user profile endpoint.
// @Summary      Get User Profile
// @Description  Retrieves the profile of the currently authenticated user.
// @Tags         User
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {object}  response.ApiResponse{data=model.User} "Successfully retrieved profile"
// @Failure      401  {object}  response.ApiResponse "Unauthorized"
// @Router       /profile [get]
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	// Retrieve the user ID from the locals, set by the auth middleware
	userID, ok := c.Locals("current_user_id").(uuid.UUID)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, errors.New("unauthorized"))
	}

	// Get user profile from the service
	user, err := h.userService.GetUserProfile(userID)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, errors.New("user not found"))
	}

	return response.Success(c, fiber.StatusOK, user)
}

// UpdateProfile correctly handles both form values and file uploads with validation.
// @Summary      Update User Profile
// @Description  Updates the name and/or avatar of the currently authenticated user.
// @Tags         User
// @Accept       multipart/form-data
// @Produce      json
// @Security     ApiKeyAuth
// @Param        name    formData  string  false  "New name for the user"
// @Param        avatar  formData  file    false  "New avatar image file"
// @Success      200  {object}  response.ApiResponse{data=model.User} "Successfully updated profile"
// @Failure      400  {object}  response.ApiResponse "Bad Request"
// @Failure      401  {object}  response.ApiResponse "Unauthorized"
// @Router       /profile [put]
type UpdateProfilePayload struct {
	Name string `json:"name" validate:"required,min=2"`
}

// UpdateProfile is the handler for updating the user's profile.
func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	// Retrieve the user ID from the locals
	userID, ok := c.Locals("current_user_id").(uuid.UUID)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, errors.New("unauthorized"))
	}

	payload := new(UpdateProfilePayload)
	payload.Name = c.FormValue("name")

	if errs := validator.ValidateStruct(payload); errs != nil {
		return response.ValidationError(c, errs)
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		if err.Error() != "there is no uploaded file associated with the given key" {
			return response.Error(c, fiber.StatusBadRequest, errors.New("invalid file upload"))
		}
		file = nil
	}

	updatedUser, err := h.userService.UpdateUserProfile(c.Context(), userID, payload.Name, file)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err)
	}

	return response.Success(c, fiber.StatusOK, updatedUser)
}

package http

import (
	"errors"
	"venturo-core/internal/service"
	"venturo-core/pkg/response"
	"venturo-core/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

// AuthHandler handles authentication-related HTTP requests.
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// RegisterPayload defines the expected JSON for registration.
type RegisterPayload struct {
	Name     string `json:"name" validate:"required,min=2"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Register is the handler for the user registration endpoint.
// @Summary      Register a new user
// @Description  Creates a new user account with the provided details.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        payload  body      RegisterPayload     true  "User Registration Payload"
// @Success      201      {object}  response.ApiResponse "User registered successfully"
// @Failure      400      {object}  response.ApiResponse "Bad Request - Invalid input"
// @Failure      500      {object}  response.ApiResponse "Internal Server Error"
// @Router       /register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	payload := new(RegisterPayload)

	// Parse the request body
	if err := c.BodyParser(payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, errors.New("cannot parse JSON"))
	}

	// Replace manual checks with a single call to the validator
	if errs := validator.ValidateStruct(payload); errs != nil {
		return response.ValidationError(c, errs)
	}

	// Call the service to register the user
	err := h.authService.Register(c.Context(), payload.Name, payload.Email, payload.Password)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err)
	}

	// Return success response
	return response.Success(c, fiber.StatusCreated, fiber.Map{"message": "User registered successfully"})
}

// Login is the handler for the user login endpoint.
// @Summary      Log in a user
// @Description  Authenticates a user and returns a JWT token.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        payload  body      LoginPayload        true  "User Login Payload"
// @Success      200      {object}  response.ApiResponse "Successfully logged in"
// @Failure      400      {object}  response.ApiResponse "Bad Request - Cannot parse JSON"
// @Failure      401      {object}  response.ApiResponse "Unauthorized - Invalid credentials"
// @Router       /login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	payload := new(LoginPayload)

	if err := c.BodyParser(payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, errors.New("cannot parse JSON"))
	}

	token, err := h.authService.Login(c.Context(), payload.Email, payload.Password)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"token": token})
}

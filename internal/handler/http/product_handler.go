package http

import (
	"errors"
	"strconv"
	"venturo-core/internal/service"
	"venturo-core/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	productService *service.ProductService
}

func NewProductHandler(s *service.ProductService) *ProductHandler {
	return &ProductHandler{productService: s}
}

// CreateProduct handles the multipart/form-data request to create a product.
//
// @Summary Create a new product
// @Description Create a new product with name, price, stock, and optional image.
// @Tags Products
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Product name"
// @Param price formData int true "Product price"
// @Param stock formData int true "Product stock"
// @Param image formData file false "Product image"
// @Success 201 {object} response.SuccessResponse{data=service.Product}
// @Failure 400 {object} response.ErrorResponse "Invalid input format"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /products [post]
func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	price, err := strconv.Atoi(c.FormValue("price"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, errors.New("invalid price format"))
	}
	stock, err := strconv.Atoi(c.FormValue("stock"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, errors.New("invalid stock format"))
	}

	input := service.CreateProductInput{
		Name:  c.FormValue("name"),
		Price: int32(price),
		Stock: int16(stock),
	}

	file, err := c.FormFile("image")
	if err == nil {
		input.Image = file
	}

	product, err := h.productService.CreateProduct(c.Context(), input)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err)
	}

	return response.Success(c, fiber.StatusCreated, product)
}

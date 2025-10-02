package http

import (
	"errors"
	"strconv"
	"venturo-core/internal/model"
	"venturo-core/internal/service"
	"venturo-core/pkg/response"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ProductHandler struct {
	productService *service.ProductService
}

func NewProductHandler(s *service.ProductService) *ProductHandler {
	return &ProductHandler{productService: s}
}

// GetAllProducts returns a list of all products.
func (h *ProductHandler) GetAllProducts(c *fiber.Ctx) error {
	products, err := h.productService.GetAllProducts()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err)
	}
	return response.Success(c, fiber.StatusOK, products)
}


// GetProductByID returns a single product by its ID.
func (h *ProductHandler) GetProductByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, errors.New("invalid ID format"))
	}
	product, err := h.productService.FindProductByID(id)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, errors.New("product not found"))
	}
	return response.Success(c, fiber.StatusOK, product)
}

// GetProductBySlug return a single product by its Slug.
func (h *ProductHandler) GetProductBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")
	product, err := h.productService.FindProductBySLug(slug)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, errors.New("product not found"))
	}
	return response.Success(c, fiber.StatusOK, product)
}

// CreateProduct handles the multipart/form-data request to create a product.
// ... (Add Swagger annotations here) ...
/**
 @api {post} api/v1/products Create Product
 @apiName PostProducts
 @apiDescription Creates or Insert a new product
 @apiGroup Product

 @apiBody {String} name Nama Product
 @apiBody {Number} price Harga Product
 @apiBody {Number} stock Stock Product

 @apiSuccess {String} id Product ID
 @apiSuccess {String} name Nama Product
 @apiSuccess {Number} price Harga Product
 @apiSuccess {Number} stock Stock Product
 @apiSuccessExample Success-Response:
   {
     "id": "1",
     "name": "Product A",
     "price": 10000
 	 "stock": 100
   }
*/
func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	price, err := strconv.Atoi(c.FormValue("price"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, errors.New("invalid price format"))
	}
	// stock, err := strconv.Atoi(c.FormValue("stock"))
	// if err != nil {
	// 	return response.Error(c, fiber.StatusBadRequest, errors.New("invalid stock format"))
	// }

	var productType []model.ProductType
	if typeStr := c.FormValue("type"); typeStr != "" {
		if err := json.Unmarshal([]byte(typeStr), &productType); err != nil {
			return response.Error(c, fiber.StatusBadRequest, errors.New("invalid type format"))
		}
	}

	input := service.CreateProductInput{
		Name:  c.FormValue("name"),
		Price: int32(price),
		// Stock: int16(stock),
		Type : productType,
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
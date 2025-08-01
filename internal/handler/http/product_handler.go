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
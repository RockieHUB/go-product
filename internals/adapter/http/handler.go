package http

import (
	"goproduct/internals/core/product/domain"
	"goproduct/internals/core/product/port"
	"strconv"

	"net/http"

	fiber "github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductHandlers struct {
	productService port.ProductService
}

var _ port.ProductHandlers = (*ProductHandlers)(nil)

func NewProductHandlers(productService port.ProductService) *ProductHandlers {
	return &ProductHandlers{
		productService: productService,
	}
}

// CreateProduct handles the creation of a new product
func (h *ProductHandlers) CreateProduct(c *fiber.Ctx) error {
	var product domain.Product

	if err := c.BodyParser(&product); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	err := h.productService.CreateProduct(&product)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create product: " + err.Error(),
		})
	}

	type ProductResponse struct {
		ProductName string  `json:"product_name"`
		Price       float64 `json:"price"`
		Stock       int     `json:"stock"`
	}

	response := ProductResponse{
		ProductName: product.ProductName,
		Price:       product.Price,
		Stock:       product.Stock,
	}
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"status_code": http.StatusCreated,
		"message":     "Product created successfully",
		"data":        response,
	})
}

// GetProduct handles retrieving a product by its ID
func (h *ProductHandlers) GetProduct(c *fiber.Ctx) error {
	productIDStr := c.Params("id")
	var productID interface{}
	var err error

	productID, err = primitive.ObjectIDFromHex(productIDStr)
	if err != nil {
		productID, err = strconv.Atoi(productIDStr)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid product ID",
			})
		}
	}

	product, err := h.productService.GetProductByID(productID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to get product"})
	}

	if product == nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"message": "Product not found"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status_code": http.StatusOK,
		"message":     "Get data success!",
		"data":        product,
	})
}

// GetAllProducts handles retrieving all products
func (h *ProductHandlers) GetAllProducts(c *fiber.Ctx) error {
	products, err := h.productService.GetAllProducts()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get all products",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status_code": http.StatusOK,
		"message":     "Get all data success!",
		"data":        products,
		"total":       len(products),
	})
}

// UpdateProduct handles updating an existing product
func (h *ProductHandlers) UpdateProduct(c *fiber.Ctx) error {
	productIDStr := c.Params("id")
	var productID interface{}
	var err error

	productID, err = primitive.ObjectIDFromHex(productIDStr)
	if err != nil {
		productID, err = strconv.Atoi(productIDStr)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid product ID",
			})
		}
	}

	product, err := h.productService.GetProductByID(productID)
	if err != nil {
		if err.Error() == "product not found" { // Or use a custom error type
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"message": "Product not found"})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to get product"})
	}
	if err := c.BodyParser(&product); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	// Set the ProductID from the URL parameter
	product.ID = &productID

	err = h.productService.UpdateProduct(product)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update product: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status_code": http.StatusOK,
		"message":     "Product updated successfully",
		"data":        product,
	})
}

// DeleteProduct handles deleting a product by its ID
func (h *ProductHandlers) DeleteProduct(c *fiber.Ctx) error {
	productIDStr := c.Params("id")
	var productID interface{}
	var err error

	productID, err = primitive.ObjectIDFromHex(productIDStr)
	if err != nil {
		productID, err = strconv.Atoi(productIDStr)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid product ID",
			})
		}
	}

	err = h.productService.DeleteProduct(productID)
	if err != nil {
		if err.Error() == "product not found" {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"message": "Product not found",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete product: " + err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status_code": http.StatusOK,
		"message":     "Delete product success!",
	})
}

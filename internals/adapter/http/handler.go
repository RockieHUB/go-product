package http

import (
	"goproduct/internals/core/product/domain"
	"goproduct/internals/core/product/port"

	fiber "github.com/gofiber/fiber/v2"
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
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	err := h.productService.CreateProduct(&product)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create product")
	}

	return c.Status(fiber.StatusCreated).JSON(product)
}

// GetProduct handles retrieving a product by its ID
func (h *ProductHandlers) GetProduct(c *fiber.Ctx) error {
	productID := c.Params("id")

	product, err := h.productService.GetProductByID(productID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get product")
	}

	if product == nil {
		return fiber.NewError(fiber.StatusNotFound, "Product not found")
	}

	return c.JSON(product)
}

// GetAllProducts handles retrieving all products
func (h *ProductHandlers) GetAllProducts(c *fiber.Ctx) error {
	products, err := h.productService.GetAllProducts()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get all products")
	}

	return c.JSON(products)
}

// UpdateProduct handles updating an existing product
func (h *ProductHandlers) UpdateProduct(c *fiber.Ctx) error {
	productID := c.Params("id")

	var product domain.Product
	if err := c.BodyParser(&product); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if product.ProductID != productID {
		return fiber.NewError(fiber.StatusBadRequest, "Product ID mismatch")
	}

	err := h.productService.UpdateProduct(&product)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update product")
	}

	return c.JSON(product)
}

// DeleteProduct handles deleting a product by its ID
func (h *ProductHandlers) DeleteProduct(c *fiber.Ctx) error {
	productID := c.Params("id")

	err := h.productService.DeleteProduct(productID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete product")
	}

	return c.SendStatus(fiber.StatusNoContent)
}

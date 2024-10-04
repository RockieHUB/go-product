package port

import (
	"goproduct/internals/core/product/domain"

	fiber "github.com/gofiber/fiber/v2"
)

// ProductService defines the interface for interacting with Product entities
type ProductService interface {
	CreateProduct(product *domain.Product) error
	GetProductByID(productID int) (*domain.Product, error)
	UpdateProduct(product *domain.Product) error
	DeleteProduct(productID int) error
	GetAllProducts() ([]*domain.Product, error)
}

// ProductRepository defines the interface for data access related to Products
type ProductRepository interface {
	SaveProduct(product *domain.Product) error
	FindProductByID(productID int) (*domain.Product, error)
	UpdateProduct(product *domain.Product) error
	DeleteProduct(productID int) error
	GetAllProducts() ([]*domain.Product, error)
}

// ProductHandlers defines the interface for handling HTTP requests related to Products
type ProductHandlers interface {
	CreateProduct(c *fiber.Ctx) error
	GetProduct(c *fiber.Ctx) error
	UpdateProduct(c *fiber.Ctx) error
	DeleteProduct(c *fiber.Ctx) error
	GetAllProducts(c *fiber.Ctx) error
}

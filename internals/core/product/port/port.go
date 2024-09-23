package port

import (
	"goproduct/internals/core/product/domain"

	fiber "github.com/gofiber/fiber/v2"
)

// ProductService defines the interface for interacting with Product entities
type ProductService interface {
	// Example methods you might need:
	CreateProduct(product *domain.Product) error
	GetProductByID(productID string) (*domain.Product, error)
	UpdateProduct(product *domain.Product) error
	DeleteProduct(productID string) error
	// Add more methods as your application requires
	GetAllProducts() ([]*domain.Product, error)
}

// ProductRepository defines the interface for data access related to Products
type ProductRepository interface {
	// Example methods you might need:
	SaveProduct(product *domain.Product) error
	FindProductByID(productID string) (*domain.Product, error)
	UpdateProduct(product *domain.Product) error
	DeleteProduct(productID string) error
	// Add more methods as your application requires
	GetAllProducts() ([]*domain.Product, error)
}

// ProductHandlers defines the interface for handling HTTP requests related to Products
type ProductHandlers interface {
	// Example methods you might need:
	CreateProduct(c *fiber.Ctx) error
	GetProduct(c *fiber.Ctx) error
	UpdateProduct(c *fiber.Ctx) error
	DeleteProduct(c *fiber.Ctx) error
	// Add more methods as your application requires
	GetAllProducts() ([]*domain.Product, error)
}

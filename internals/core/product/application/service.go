package application

import (
	"errors"
	"goproduct/internals/core/product/domain"
	"goproduct/internals/core/product/port"
)

// ProductService implements the ports.ProductService interface
type ProductService struct {
	productRepository port.ProductRepository
}

// Ensure ProductService implements the interface
var _ port.ProductService = (*ProductService)(nil)

// NewProductService creates a new ProductService instance
func NewProductService(repository port.ProductRepository) *ProductService {
	return &ProductService{
		productRepository: repository,
	}
}

// Example service methods (you'll need to implement the actual logic):

// CreateProduct creates a new product
func (s *ProductService) CreateProduct(product *domain.Product) error {
	// You might add validation here (e.g., check for required fields, uniqueness, etc.)
	if product.ProductName == "" {
		return errors.New("product name are required")
	}
	return s.productRepository.SaveProduct(product)
}

// GetProductByID retrieves a product by its ID
func (s *ProductService) GetProductByID(productID interface{}) (*domain.Product, error) {
	if productID == nil {
		return nil, errors.New("product ID is required")
	}

	return s.productRepository.FindProductByID(productID)
}
func (s *ProductService) GetAllProducts() ([]*domain.Product, error) {
	return s.productRepository.GetAllProducts()
}

// UpdateProduct updates an existing product
func (s *ProductService) UpdateProduct(product *domain.Product) error {
	// You might add validation here and ensure the product exists before updating
	if product.ID == nil {
		return errors.New("product ID is required for update")
	}

	return s.productRepository.UpdateProduct(product)
}

// DeleteProduct deletes a product by its ID
func (s *ProductService) DeleteProduct(productID interface{}) error {
	if productID == nil {
		return errors.New("product ID is required for deletion")
	}

	return s.productRepository.DeleteProduct(productID)
}

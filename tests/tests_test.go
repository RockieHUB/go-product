package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"goproduct/internals/adapter/http"
	"goproduct/internals/core/product/application"
	"goproduct/internals/core/product/domain"
	"io"
	netHTTP "net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProductRepository is a mock implementation of the ProductRepository interface
type MockProductRepository struct {
	mock.Mock
}

// SaveProduct mocks the SaveProduct method
func (m *MockProductRepository) SaveProduct(product *domain.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

// FindProductByID mocks the FindProductByID method
func (m *MockProductRepository) FindProductByID(productID int) (*domain.Product, error) {
	args := m.Called(productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

// GetAllProducts mocks the GetAllProducts method
func (m *MockProductRepository) GetAllProducts() ([]*domain.Product, error) {
	args := m.Called()
	return args.Get(0).([]*domain.Product), args.Error(1)
}

// UpdateProduct mocks the UpdateProduct method
func (m *MockProductRepository) UpdateProduct(product *domain.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

// DeleteProduct mocks the DeleteProduct method
func (m *MockProductRepository) DeleteProduct(productID int) error {
	args := m.Called(productID)
	return args.Error(0)
}

func TestAPI(t *testing.T) {
	app := fiber.New()

	// Create mock repository
	mockRepo := new(MockProductRepository)

	// Initialize services and handlers with the mock repository
	productService := application.NewProductService(mockRepo)
	productHandler := http.NewProductHandlers(productService)

	app.Get("/products", productHandler.GetAllProducts)
	app.Post("/products", productHandler.CreateProduct)
	app.Put("/products/:id", productHandler.UpdateProduct)
	app.Get("/products/:id", productHandler.GetProduct)
	app.Delete("/products/:id", productHandler.DeleteProduct)

	t.Run("GET /products", func(t *testing.T) {
		t.Run("returns a list of products when products exist", func(t *testing.T) {
			// Mock the GetAllProducts method to return some test products
			mockProducts := []*domain.Product{
				{ProductID: new(int), ProductName: "Test Product 1", Price: 10.0, Stock: 5},
				{ProductID: new(int), ProductName: "Test Product 2", Price: 20.0, Stock: 10},
			}
			mockRepo.On("GetAllProducts").Return(mockProducts, nil)

			req := httptest.NewRequest(netHTTP.MethodGet, "/products", nil)
			resp, err := app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, netHTTP.StatusOK, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			var responseBody struct {
				Data []domain.Product `json:"data"`
			}
			err = json.Unmarshal(body, &responseBody)
			assert.NoError(t, err)

			assert.Equal(t, len(mockProducts), len(responseBody.Data))
			assert.NotEmpty(t, responseBody.Data)

			mockRepo.AssertExpectations(t)
		})

		t.Run("returns an empty list when no products exist", func(t *testing.T) {
			mockRepo := new(MockProductRepository)
			// Initialize services and handlers with the mock repository
			productService := application.NewProductService(mockRepo)
			productHandler := http.NewProductHandlers(productService)

			app := fiber.New()
			app.Get("/products", productHandler.GetAllProducts)

			// Mock the GetAllProducts method to return an empty list
			mockRepo.On("GetAllProducts").Return([]*domain.Product{}, nil)

			req := httptest.NewRequest(netHTTP.MethodGet, "/products", nil)
			resp, err := app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, netHTTP.StatusOK, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			var responseBody struct {
				Data []domain.Product `json:"data"`
			}
			err = json.Unmarshal(body, &responseBody)
			assert.NoError(t, err)

			assert.Empty(t, responseBody.Data)

			mockRepo.AssertExpectations(t)
		})
	})

	t.Run("POST /products", func(t *testing.T) {
		t.Run("creates a new product", func(t *testing.T) {
			newProduct := domain.Product{
				ProductName: "New Product",
				Price:       19.99,
				Stock:       50,
			}
			requestBody, _ := json.Marshal(newProduct)

			// Expect SaveProduct to be called with the new product and return no error
			mockRepo.On("SaveProduct", mock.AnythingOfType("*domain.Product")).Return(nil)

			req := httptest.NewRequest(netHTTP.MethodPost, "/products", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, netHTTP.StatusCreated, resp.StatusCode)

			// Assert the response body if needed (e.g., check for success message)
			// ...

			mockRepo.AssertExpectations(t)
		})

		t.Run("returns an error if input is invalid", func(t *testing.T) {
			invalidProduct := domain.Product{
				Price: 19.99,
				Stock: 50,
			}
			requestBody, _ := json.Marshal(invalidProduct)

			req := httptest.NewRequest(netHTTP.MethodPost, "/products", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, netHTTP.StatusInternalServerError, resp.StatusCode) // Or your expected error code

			// Assert the response body contains an error message
			// ...
		})
	})

	t.Run("PUT /products/:id", func(t *testing.T) {
		t.Run("updates an existing product", func(t *testing.T) {
			existingProduct := &domain.Product{
				ProductID:   new(int), // Assign a dummy ID
				ProductName: "Existing Product",
				Price:       10.0,
				Stock:       5,
			}
			*existingProduct.ProductID = 1 // Set the ID to 1 for testing

			updatedProduct := domain.Product{
				ProductID:   existingProduct.ProductID,
				ProductName: "Updated Product",
				Price:       12.0,
				Stock:       8,
			}
			requestBody, _ := json.Marshal(updatedProduct)

			// Expect FindProductByID to be called and return the existing product
			mockRepo.On("FindProductByID", 1).Return(existingProduct, nil)
			// Expect UpdateProduct to be called and return no error
			mockRepo.On("UpdateProduct", mock.AnythingOfType("*domain.Product")).Return(nil)

			req := httptest.NewRequest(netHTTP.MethodPut, "/products/1", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, netHTTP.StatusOK, resp.StatusCode)

			// ... assert the response body if needed ...

			mockRepo.AssertExpectations(t)
		})

		t.Run("returns an error if product is not found", func(t *testing.T) {
			mockRepo.On("FindProductByID", 2).Return(nil, fmt.Errorf("product not found")) // Simulate product not found

			req := httptest.NewRequest(netHTTP.MethodPut, "/products/2", nil) // No need for request body in this case
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, netHTTP.StatusNotFound, resp.StatusCode) // Or your expected error code

			// ... assert the response body contains an error message ...
		})
	})

	t.Run("GET /products/:id", func(t *testing.T) {
		t.Run("returns a product by ID", func(t *testing.T) {
			mockProduct := &domain.Product{
				ProductID:   new(int),
				ProductName: "Test Product",
				Price:       10.0,
				Stock:       5,
			}
			*mockProduct.ProductID = 1
			mockRepo.On("FindProductByID", 1).Return(mockProduct, nil)

			req := httptest.NewRequest(netHTTP.MethodGet, "/products/1", nil)
			resp, err := app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, netHTTP.StatusOK, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			var responseBody struct {
				Data *domain.Product `json:"data"`
			}
			err = json.Unmarshal(body, &responseBody)
			assert.NoError(t, err)

			assert.NotNil(t, responseBody.Data)
			assert.Equal(t, *mockProduct.ProductID, *responseBody.Data.ProductID)
			// ... add more assertions to check other fields if needed ...

			mockRepo.AssertExpectations(t)
		})

		t.Run("returns an error if product is not found", func(t *testing.T) {
			mockRepo.On("FindProductByID", 2).Return(nil, fmt.Errorf("product not found"))

			req := httptest.NewRequest(netHTTP.MethodGet, "/products/2", nil)
			resp, err := app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, netHTTP.StatusNotFound, resp.StatusCode) // Or your expected error code

			// ... assert the response body contains an error message ...
		})
	})

	t.Run("DELETE /products/:id", func(t *testing.T) {
		t.Run("deletes a product by ID", func(t *testing.T) {
			// Expect DeleteProduct to be called and return no error
			mockRepo.On("DeleteProduct", 1).Return(nil)

			req := httptest.NewRequest(netHTTP.MethodDelete, "/products/1", nil)
			resp, err := app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, netHTTP.StatusOK, resp.StatusCode) // Or your expected success code

			// ... assert the response body if needed ...

			mockRepo.AssertExpectations(t)
		})

		t.Run("returns an error if product is not found", func(t *testing.T) {
			mockRepo.On("DeleteProduct", 2).Return(fmt.Errorf("product not found"))

			req := httptest.NewRequest(netHTTP.MethodDelete, "/products/2", nil)
			resp, err := app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, netHTTP.StatusNotFound, resp.StatusCode) // Or your expected error code

			// ... assert the response body contains an error message ...
		})
	})
}

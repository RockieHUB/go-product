package main

import (
	"fmt"
	"log"
	"os"

	"goproduct/internals/adapter/http"
	"goproduct/internals/adapter/repository/mysql_repository"

	"goproduct/internals/core/product/application"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Construct the MySQL DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))

	// Create the product repository
	productRepository, err := mysql_repository.NewProductRepository(dsn)
	if err != nil {
		log.Fatal("Error creating product repository:", err)
	}

	// Create the product service
	productService := application.NewProductService(productRepository)

	// Create the product handlers
	productHandlers := http.NewProductHandlers(productService)

	// Initialize Fiber app
	app := fiber.New()

	// Define routes
	v1 := app.Group("/v1")
	productRoutes := v1.Group("/products")
	productRoutes.Post("/", productHandlers.CreateProduct)
	productRoutes.Get("/:id", productHandlers.GetProduct)
	productRoutes.Put("/:id", productHandlers.UpdateProduct)
	productRoutes.Delete("/:id", productHandlers.DeleteProduct)

	// Adapter function for GetAllProducts
	productRoutes.Get("/", func(c *fiber.Ctx) error {
		products, err := productHandlers.GetAllProducts()
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to get all products")
		}
		return c.JSON(products)
	})

	// Start the server
	err = app.Listen(":5000")
	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"fmt"
	"log"

	"goproduct/internals/adapter/http"
	"goproduct/internals/adapter/repository/mysql_repository"
	"goproduct/internals/config"
	"goproduct/internals/core/product/application"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfigFromEnv()
	if err != nil {
		log.Fatal("Error loading configuration:", err)
	}

	// Create the product repository
	productRepository, err := mysql_repository.NewProductRepository(cfg.Database.DSN)
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
	productRoutes.Get("/", productHandlers.GetAllProducts)
	productRoutes.Get("/:id", productHandlers.GetProduct)
	productRoutes.Put("/:id", productHandlers.UpdateProduct)
	productRoutes.Delete("/:id", productHandlers.DeleteProduct)

	// Start the server
	err = app.Listen(fmt.Sprintf(":%d", cfg.Server.Port))
	if err != nil {
		log.Fatal(err)
	}
}

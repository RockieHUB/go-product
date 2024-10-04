package main

import (
	"fmt"
	"log"

	"goproduct/internals/adapter/http"
	"goproduct/internals/adapter/repository/mongodb_repository"
	"goproduct/internals/adapter/repository/mysql_repository"
	"goproduct/internals/config"
	"goproduct/internals/core/product/application"
	"goproduct/internals/core/product/port"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfigFromEnv()
	if err != nil {
		log.Fatal("Error loading configuration:", err)
	}

	// Create the product repository
	var productRepository port.ProductRepository
	switch cfg.Database.Type {
	case "mysql":
		productRepository, err = mysql_repository.NewProductRepository(cfg.Database.MySQL.DSN)
	case "mongodb":
		productRepository, err = mongodb_repository.NewProductRepository(
			cfg.Database.MongoDB.URI,
			cfg.Database.MongoDB.Database,
			cfg.Database.MongoDB.Collection,
		)
	default:
		log.Fatalf("Unsupported database type: %s", cfg.Database.Type)
	}
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

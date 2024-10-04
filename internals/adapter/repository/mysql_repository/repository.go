package mysql_repository

import (
	"database/sql"
	"goproduct/internals/core/product/domain"
	"goproduct/internals/core/product/port"

	_ "github.com/go-sql-driver/mysql"
)

type ProductRepository struct {
	db *sql.DB
}

var _ port.ProductRepository = (*ProductRepository)(nil)

func NewProductRepository(dsn string) (*ProductRepository, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &ProductRepository{db: db}, nil
}

func (r *ProductRepository) SaveProduct(product *domain.Product) error {
	query := "INSERT INTO Product (product_name, price, stock) VALUES (?, ?, ?)"
	_, err := r.db.Exec(query, product.ProductName, product.Price, product.Stock)
	return err
}

func (r *ProductRepository) FindProductByID(productID int) (*domain.Product, error) {
	query := "SELECT product_id, product_name, price, stock FROM Product WHERE product_id = ?"
	var product domain.Product
	err := r.db.QueryRow(query, productID).Scan(&product.ProductID, &product.ProductName, &product.Price, &product.Stock)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, err
	}
	return &product, nil
}
func (r *ProductRepository) GetAllProducts() ([]*domain.Product, error) {
	query := "SELECT product_id, product_name, price, stock FROM Product"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		var product domain.Product
		if err := rows.Scan(&product.ProductID, &product.ProductName, &product.Price, &product.Stock); err != nil {
			return nil, err
		}
		products = append(products, &product)
	}

	return products, nil
}

func (r *ProductRepository) UpdateProduct(product *domain.Product) error {
	query := "UPDATE Product SET product_name = ?, price = ?, stock = ? WHERE product_id = ?"
	_, err := r.db.Exec(query, product.ProductName, product.Price, product.Stock, product.ProductID)
	return err
}

func (r *ProductRepository) DeleteProduct(productID int) error {
	query := "DELETE FROM Product WHERE product_id = ?"
	_, err := r.db.Exec(query, productID)
	return err
}

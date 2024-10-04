package domain

type Product struct {
	ProductID   *int    `json:"product_id" bson:"product_id,omitempty"`
	ProductName string  `json:"product_name"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

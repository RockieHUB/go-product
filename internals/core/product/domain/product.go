package domain

type Product struct {
	ID          interface{} `json:"id" bson:"_id,omitempty"`
	ProductName string      `json:"product_name" bson:"productname"`
	Price       float64     `json:"price" bson:"price"`
	Stock       int         `json:"stock" bson:"stock"`
}

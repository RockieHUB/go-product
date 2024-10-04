package mongodb_repository

import (
	"context"
	"errors"
	"goproduct/internals/core/product/domain"
	"goproduct/internals/core/product/port"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductRepository struct {
	client     *mongo.Client
	database   string
	collection string
}

var _ port.ProductRepository = (*ProductRepository)(nil)

func NewProductRepository(uri, database, collection string) (*ProductRepository, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return &ProductRepository{
		client:     client,
		database:   database,
		collection: collection,
	}, nil
}

func (r *ProductRepository) SaveProduct(product *domain.Product) error {
	coll := r.client.Database(r.database).Collection(r.collection)
	_, err := coll.InsertOne(context.Background(), product)
	return err
}

func (r *ProductRepository) FindProductByID(productID int) (*domain.Product, error) {
	coll := r.client.Database(r.database).Collection(r.collection)
	var product domain.Product
	err := coll.FindOne(context.Background(), bson.M{"product_id": productID}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Not found
		}
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) GetAllProducts() ([]*domain.Product, error) {
	coll := r.client.Database(r.database).Collection(r.collection)
	cursor, err := coll.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var products []*domain.Product
	for cursor.Next(context.Background()) {
		var product domain.Product
		if err := cursor.Decode(&product); err != nil {
			return nil, err
		}
		products = append(products, &product)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductRepository) UpdateProduct(product *domain.Product) error {
	coll := r.client.Database(r.database).Collection(r.collection)
	result, err := coll.ReplaceOne(
		context.Background(),
		bson.M{"product_id": product.ProductID},
		product,
	)
	if err != nil {
		return err
	}
	if result.ModifiedCount == 0 {
		return errors.New("no document was updated")
	}
	return nil
}

func (r *ProductRepository) DeleteProduct(productID int) error {
	coll := r.client.Database(r.database).Collection(r.collection)
	result, err := coll.DeleteOne(context.Background(), bson.M{"product_id": productID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("no document was deleted")
	}
	return nil
}

package repositories

import (
	"gorm.io/gorm"
)

type APIRepository struct {
	db *gorm.DB
}

func NewAPIRepository(db *gorm.DB) *APIRepository {
	return &APIRepository{db: db}
}


type Product struct {
    ID          int     `json:"id"`
    Name        string  `json:"name"`
    Price       float64 `json:"price"`
    Description string  `json:"description"`
    Views       int     `json:"views"`
    Brand       string  `json:"brand"`
}

// Static product data
var products = []Product{
    {ID: 1, Name: "Product 1", Price: 10.99, Description: "This is Product 1", Views: 150, Brand: "BrandA"},
    {ID: 2, Name: "Product 2", Price: 15.99, Description: "This is Product 2", Views: 200, Brand: "BrandB"},
    {ID: 3, Name: "Product 3", Price: 20.99, Description: "This is Product 3", Views: 250, Brand: "BrandC"},
}

// GetProductById retrieves a product by its ID
func (r *APIRepository) GetProductById(id int) (interface{}, error) {
    for _, product := range products {
        if product.ID == id {
            return product, nil
        }
    }
    // If no product matches, return an empty array
    return []Product{}, nil
}
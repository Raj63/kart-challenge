// Package models provides data structures for the Order Food Online service.
package models

// Product represents a product in the catalog.
// It contains basic product information including pricing and categorization.
type Product struct {
	ID       string  `bson:"id" json:"id"`             // Unique identifier for the product
	Name     string  `bson:"name" json:"name"`         // Product name
	Price    float64 `bson:"price" json:"price"`       // Product price
	Category string  `bson:"category" json:"category"` // Product category
}

package models

// Product represents a product in the catalog.
type Product struct {
	ID       string  `bson:"id" json:"id"`
	Name     string  `bson:"name" json:"name"`
	Price    float64 `bson:"price" json:"price"`
	Category string  `bson:"category" json:"category"`
}

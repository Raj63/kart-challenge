package models

// OrderItem represents an item in an order.
type OrderItem struct {
	ProductID string `bson:"productId" json:"productId"`
	Quantity  int    `bson:"quantity" json:"quantity"`
}

// Order represents a placed order.
type Order struct {
	ID       string      `bson:"id" json:"id"`
	Items    []OrderItem `bson:"items" json:"items"`
	Products []Product   `bson:"products" json:"products"`
}

// OrderCreateRequest represents the request body for placing an order.
type OrderCreateRequest struct {
	CouponCode string      `json:"couponCode"`
	Items      []OrderItem `json:"items"`
}

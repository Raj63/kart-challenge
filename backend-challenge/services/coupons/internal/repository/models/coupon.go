// Package models provides data structures for the Coupons processor service.
package models

// Coupon represents a coupon code entry in the database.
// It stores individual coupon codes with metadata about their source file
// and processing status.
type Coupon struct {
	ID         string `bson:"id" json:"id"`                   // Unique identifier for the coupon
	FileName   string `bson:"file_name" json:"file_name"`     // Name of the source file that contained this coupon
	CouponCode string `bson:"coupon_code" json:"coupon_code"` // The actual coupon code string
	Datetime   int64  `bson:"datetime" json:"datetime"`       // Unix timestamp when the coupon was processed
	IsActive   bool   `bson:"isactive" json:"isactive"`       // Whether the coupon is currently active
}

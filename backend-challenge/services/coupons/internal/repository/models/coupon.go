package models

// Coupon represents a coupon code entry in the database.
type Coupon struct {
	ID         string `bson:"id" json:"id"`
	FileName   string `bson:"file_name" json:"file_name"`
	CouponCode string `bson:"coupon_code" json:"coupon_code"`
	Datetime   int64  `bson:"datetime" json:"datetime"`
	IsActive   bool   `bson:"isactive" json:"isactive"`
}

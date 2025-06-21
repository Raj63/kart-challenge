// Package models provides data structures for the Coupons processor service.
package models

// ProcessedCouponFile represents a processed coupon file record in the database.
// It tracks the processing status and metadata of coupon files to support
// resume functionality and prevent duplicate processing.
type ProcessedCouponFile struct {
	ID              string `bson:"id" json:"id"`                                 // Unique identifier for the processed file
	MD5Hash         string `bson:"md5hash" json:"md5hash"`                       // MD5 hash of the file content
	FileName        string `bson:"file_name" json:"file_name"`                   // Name of the processed file
	IsAdd           bool   `bson:"isadd" json:"isadd"`                           // Whether this was an add operation (true) or remove operation (false)
	Size            int64  `bson:"size" json:"size"`                             // Size of the processed file in bytes
	CouponCodeCount int64  `bson:"coupon_code_counts" json:"coupon_code_counts"` // Number of coupon codes processed
	Datetime        int64  `bson:"datetime" json:"datetime"`                     // Unix timestamp when the file was processed
	Status          string `bson:"status" json:"status"`                         // Status of the processing (e.g., "initated", "completed", "failed")
}

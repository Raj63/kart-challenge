package service

const (
	// ProductListingError indicates an error occurred while listing products.
	ProductListingError = "error listing products"

	// FindProductByIDError indicates a failure while fetching a product by its ID.
	FindProductByIDError = "error fetching product by ID"

	// InvalidPromoCode is returned when the provided promo code is not valid.
	InvalidPromoCode = "invalid promo code"

	// InvalidProductOrQuantity is returned when either the product ID or quantity is invalid.
	InvalidProductOrQuantity = "invalid productId or quantity"

	// ProductNotFound is returned when the specified product could not be found.
	ProductNotFound = "product not found: "
)

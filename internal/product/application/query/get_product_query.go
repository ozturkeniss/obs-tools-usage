package query

// GetProductQuery represents a query to get a product by ID
type GetProductQuery struct {
	ID int `json:"id" binding:"required"`
}

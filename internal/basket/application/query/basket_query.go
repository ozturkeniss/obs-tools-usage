package query

// GetBasketQuery represents a query to get a basket
type GetBasketQuery struct {
	UserID string `json:"user_id" binding:"required"`
}

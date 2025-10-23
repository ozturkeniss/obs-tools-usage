package query

// GetBasketQuery represents a query to get a basket
type GetBasketQuery struct {
	UserID string `json:"user_id" binding:"required"`
}

// GetBasketItemsQuery represents a query to get basket items
type GetBasketItemsQuery struct {
	UserID string `json:"user_id" binding:"required"`
}

// GetBasketTotalQuery represents a query to get basket total
type GetBasketTotalQuery struct {
	UserID string `json:"user_id" binding:"required"`
}

// GetBasketItemCountQuery represents a query to get basket item count
type GetBasketItemCountQuery struct {
	UserID string `json:"user_id" binding:"required"`
}

// GetBasketByCategoryQuery represents a query to get basket items by category
type GetBasketByCategoryQuery struct {
	UserID   string `json:"user_id" binding:"required"`
	Category string `json:"category" binding:"required"`
}

// GetBasketStatsQuery represents a query to get basket statistics
type GetBasketStatsQuery struct {
	UserID string `json:"user_id" binding:"required"`
}

// GetBasketExpiryQuery represents a query to get basket expiry info
type GetBasketExpiryQuery struct {
	UserID string `json:"user_id" binding:"required"`
}

// GetBasketHistoryQuery represents a query to get basket history
type GetBasketHistoryQuery struct {
	UserID string `json:"user_id" binding:"required"`
}

// GetBasketRecommendationsQuery represents a query to get basket recommendations
type GetBasketRecommendationsQuery struct {
	UserID string `json:"user_id" binding:"required"`
}

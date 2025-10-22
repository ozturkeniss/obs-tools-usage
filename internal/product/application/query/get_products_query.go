package query

// GetProductsQuery represents a query to get all products
type GetProductsQuery struct {
	// No filters for now, can add pagination/filters later
}

// GetTopMostExpensiveQuery represents a query to get top most expensive products
type GetTopMostExpensiveQuery struct {
	Limit int `json:"limit" binding:"required,min=1"`
}

// GetLowStockProductsQuery represents a query to get low stock products
type GetLowStockProductsQuery struct {
	MaxStock int `json:"max_stock" binding:"required,min=0"`
}

// GetProductsByCategoryQuery represents a query to get products by category
type GetProductsByCategoryQuery struct {
	Category string `json:"category" binding:"required"`
}

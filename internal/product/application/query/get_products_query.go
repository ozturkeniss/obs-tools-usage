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

// GetProductsByPriceRangeQuery represents a query to get products by price range
type GetProductsByPriceRangeQuery struct {
	MinPrice float64 `json:"min_price" binding:"required"`
	MaxPrice float64 `json:"max_price" binding:"required"`
}

// GetProductsByNameQuery represents a query to get products by name
type GetProductsByNameQuery struct {
	Name string `json:"name" binding:"required"`
}

// GetProductStatsQuery represents a query to get product statistics
type GetProductStatsQuery struct{}

// GetCategoriesQuery represents a query to get categories
type GetCategoriesQuery struct{}

// GetProductsByStockQuery represents a query to get products by stock
type GetProductsByStockQuery struct {
	Stock int `json:"stock" binding:"required"`
}

// GetRandomProductsQuery represents a query to get random products
type GetRandomProductsQuery struct {
	Count int `json:"count" binding:"required"`
}

// GetProductsByDateRangeQuery represents a query to get products by date range
type GetProductsByDateRangeQuery struct {
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
}

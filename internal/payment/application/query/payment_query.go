package query

// GetPaymentQuery represents a query to get a payment
type GetPaymentQuery struct {
	PaymentID string `json:"payment_id" binding:"required"`
}

// GetPaymentsByUserQuery represents a query to get payments by user
type GetPaymentsByUserQuery struct {
	UserID string `json:"user_id" binding:"required"`
}

// GetPaymentsByBasketQuery represents a query to get payments by basket
type GetPaymentsByBasketQuery struct {
	BasketID string `json:"basket_id" binding:"required"`
}

// GetPaymentsByStatusQuery represents a query to get payments by status
type GetPaymentsByStatusQuery struct {
	Status string `json:"status" binding:"required"`
}

// GetPaymentStatsQuery represents a query to get payment statistics
type GetPaymentStatsQuery struct {
	UserID string `json:"user_id"`
}

// GetPaymentsByDateRangeQuery represents a query to get payments by date range
type GetPaymentsByDateRangeQuery struct {
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
}

// GetPaymentsByAmountRangeQuery represents a query to get payments by amount range
type GetPaymentsByAmountRangeQuery struct {
	MinAmount float64 `json:"min_amount" binding:"required"`
	MaxAmount float64 `json:"max_amount" binding:"required"`
}

// GetPaymentsByMethodQuery represents a query to get payments by method
type GetPaymentsByMethodQuery struct {
	Method string `json:"method" binding:"required"`
}

// GetPaymentsByProviderQuery represents a query to get payments by provider
type GetPaymentsByProviderQuery struct {
	Provider string `json:"provider" binding:"required"`
}

// GetPaymentItemsQuery represents a query to get payment items
type GetPaymentItemsQuery struct {
	PaymentID string `json:"payment_id" binding:"required"`
}

// GetPaymentAnalyticsQuery represents a query to get payment analytics
type GetPaymentAnalyticsQuery struct{}

// GetPaymentMethodsQuery represents a query to get payment methods
type GetPaymentMethodsQuery struct{}

// GetPaymentProvidersQuery represents a query to get payment providers
type GetPaymentProvidersQuery struct{}

// GetPaymentSummaryQuery represents a query to get payment summary
type GetPaymentSummaryQuery struct{}

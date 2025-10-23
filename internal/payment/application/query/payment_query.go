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

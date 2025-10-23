package dto

import "time"

// CreatePaymentRequest represents the request payload for creating a payment
type CreatePaymentRequest struct {
	UserID      string            `json:"user_id" binding:"required"`
	BasketID    string            `json:"basket_id" binding:"required"`
	Method      string            `json:"method" binding:"required"`
	Provider    string            `json:"provider" binding:"required"`
	Currency    string            `json:"currency"`
	Description string            `json:"description"`
	Metadata    map[string]string `json:"metadata"`
}

// UpdatePaymentRequest represents the request payload for updating a payment
type UpdatePaymentRequest struct {
	Status   string            `json:"status" binding:"required"`
	Metadata map[string]string `json:"metadata"`
}

// ProcessPaymentRequest represents the request payload for processing a payment
type ProcessPaymentRequest struct {
	PaymentID string `json:"payment_id" binding:"required"`
	ProviderID string `json:"provider_id"`
}

// RefundPaymentRequest represents the request payload for refunding a payment
type RefundPaymentRequest struct {
	PaymentID string  `json:"payment_id" binding:"required"`
	Amount    float64 `json:"amount"`
	Reason    string  `json:"reason"`
}

// PaymentItemResponse represents a payment item in response
type PaymentItemResponse struct {
	ID        string  `json:"id"`
	ProductID int     `json:"product_id"`
	Name      string  `json:"name"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	Subtotal  float64 `json:"subtotal"`
	Category  string  `json:"category"`
	CreatedAt time.Time `json:"created_at"`
}

// PaymentResponse represents the response payload for payment operations
type PaymentResponse struct {
	ID          string                `json:"id"`
	UserID      string                `json:"user_id"`
	BasketID    string                `json:"basket_id"`
	Amount      float64               `json:"amount"`
	Currency    string                `json:"currency"`
	Status      string                `json:"status"`
	Method      string                `json:"method"`
	Provider    string                `json:"provider"`
	ProviderID  string                `json:"provider_id"`
	Description string                `json:"description"`
	Metadata    map[string]string     `json:"metadata"`
	Items       []PaymentItemResponse `json:"items"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
	ProcessedAt *time.Time            `json:"processed_at"`
	ExpiresAt   *time.Time            `json:"expires_at"`
}

// PaymentStatsResponse represents payment statistics response
type PaymentStatsResponse struct {
	TotalPayments     int64   `json:"total_payments"`
	TotalAmount       float64 `json:"total_amount"`
	CompletedPayments int64   `json:"completed_payments"`
	FailedPayments    int64   `json:"failed_payments"`
	PendingPayments   int64   `json:"pending_payments"`
	AverageAmount     float64 `json:"average_amount"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// CancelPaymentRequest represents the request payload for cancelling a payment
type CancelPaymentRequest struct {
	PaymentID string `json:"payment_id" binding:"required"`
}

// RetryPaymentRequest represents the request payload for retrying a payment
type RetryPaymentRequest struct {
	PaymentID string `json:"payment_id" binding:"required"`
}

// PaymentAnalyticsResponse represents payment analytics response
type PaymentAnalyticsResponse struct {
	TotalPayments     int64   `json:"total_payments"`
	TotalRevenue      float64 `json:"total_revenue"`
	SuccessRate       float64 `json:"success_rate"`
	AverageAmount     float64 `json:"average_amount"`
	TopPaymentMethod  string  `json:"top_payment_method"`
	TopProvider       string  `json:"top_provider"`
	DailyTransactions int64   `json:"daily_transactions"`
	MonthlyRevenue    float64 `json:"monthly_revenue"`
}

// PaymentMethodsResponse represents payment methods response
type PaymentMethodsResponse struct {
	Methods []string `json:"methods"`
	Count   int      `json:"count"`
}

// PaymentProvidersResponse represents payment providers response
type PaymentProvidersResponse struct {
	Providers []string `json:"providers"`
	Count     int      `json:"count"`
}

// PaymentSummaryResponse represents payment summary response
type PaymentSummaryResponse struct {
	TotalPayments     int64   `json:"total_payments"`
	TotalRevenue      float64 `json:"total_revenue"`
	PendingPayments   int64   `json:"pending_payments"`
	CompletedPayments int64   `json:"completed_payments"`
	FailedPayments    int64   `json:"failed_payments"`
	RefundedPayments  int64   `json:"refunded_payments"`
	SuccessRate       float64 `json:"success_rate"`
	AverageAmount     float64 `json:"average_amount"`
}

// HealthResponse represents a health check response
type HealthResponse struct {
	Service   string `json:"service"`
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
}

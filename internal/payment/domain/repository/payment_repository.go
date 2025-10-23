package repository

import (
	"obs-tools-usage/internal/payment/domain/entity"
)

// PaymentRepository defines the interface for payment data access
type PaymentRepository interface {
	// Basic CRUD operations
	CreatePayment(payment *entity.Payment) error
	GetPayment(paymentID string) (*entity.Payment, error)
	UpdatePayment(payment *entity.Payment) error
	DeletePayment(paymentID string) error
	
	// Query operations
	GetPaymentsByUser(userID string) ([]*entity.Payment, error)
	GetPaymentsByBasket(basketID string) ([]*entity.Payment, error)
	GetPaymentsByStatus(status entity.PaymentStatus) ([]*entity.Payment, error)
	GetPaymentsByDateRange(startDate, endDate string) ([]*entity.Payment, error)
	
	// Payment items
	CreatePaymentItem(item *entity.PaymentItem) error
	GetPaymentItems(paymentID string) ([]*entity.PaymentItem, error)
	DeletePaymentItems(paymentID string) error
	
	// Statistics and analytics
	GetPaymentStats(userID string) (*PaymentStats, error)
	GetTotalRevenue(startDate, endDate string) (float64, error)
	GetPaymentCountByStatus(status entity.PaymentStatus) (int64, error)
	
	// New query methods
	GetPaymentsByAmountRange(minAmount, maxAmount float64) ([]*entity.Payment, error)
	GetPaymentsByMethod(method string) ([]*entity.Payment, error)
	GetPaymentsByProvider(provider string) ([]*entity.Payment, error)
	GetPaymentAnalytics() (*PaymentAnalytics, error)
	GetPaymentMethods() ([]string, error)
	GetPaymentProviders() ([]string, error)
	GetPaymentSummary() (*PaymentSummary, error)
	
	// Health check
	Ping() error
}

// PaymentStats represents payment statistics
type PaymentStats struct {
	TotalPayments     int64   `json:"total_payments"`
	TotalAmount       float64 `json:"total_amount"`
	CompletedPayments int64   `json:"completed_payments"`
	FailedPayments    int64   `json:"failed_payments"`
	PendingPayments   int64   `json:"pending_payments"`
	AverageAmount     float64 `json:"average_amount"`
}

// PaymentAnalytics represents payment analytics
type PaymentAnalytics struct {
	TotalPayments     int64   `json:"total_payments"`
	TotalRevenue      float64 `json:"total_revenue"`
	SuccessRate       float64 `json:"success_rate"`
	AverageAmount     float64 `json:"average_amount"`
	TopPaymentMethod  string  `json:"top_payment_method"`
	TopProvider       string  `json:"top_provider"`
	DailyTransactions int64   `json:"daily_transactions"`
	MonthlyRevenue    float64 `json:"monthly_revenue"`
}

// PaymentSummary represents payment summary
type PaymentSummary struct {
	TotalPayments     int64   `json:"total_payments"`
	TotalRevenue      float64 `json:"total_revenue"`
	PendingPayments   int64   `json:"pending_payments"`
	CompletedPayments int64   `json:"completed_payments"`
	FailedPayments    int64   `json:"failed_payments"`
	RefundedPayments  int64   `json:"refunded_payments"`
	SuccessRate       float64 `json:"success_rate"`
	AverageAmount     float64 `json:"average_amount"`
}

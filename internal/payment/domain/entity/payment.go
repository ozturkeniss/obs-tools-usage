package entity

import (
	"time"
)

// Payment represents a payment transaction
type Payment struct {
	ID          string            `json:"id" gorm:"primaryKey"`
	UserID      string            `json:"user_id" gorm:"not null;index"`
	BasketID    string            `json:"basket_id" gorm:"not null;index"`
	Amount      float64           `json:"amount" gorm:"not null"`
	Currency    string            `json:"currency" gorm:"not null;default:'USD'"`
	Status      PaymentStatus     `json:"status" gorm:"not null;default:'pending'"`
	Method      PaymentMethod     `json:"method" gorm:"not null"`
	Provider    string            `json:"provider" gorm:"not null"`
	ProviderID  string            `json:"provider_id" gorm:"index"`
	Description string            `json:"description"`
	Metadata    map[string]string `json:"metadata" gorm:"type:json"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	ProcessedAt *time.Time        `json:"processed_at"`
	ExpiresAt   *time.Time        `json:"expires_at"`
}

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusProcessing PaymentStatus = "processing"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusCancelled PaymentStatus = "cancelled"
	PaymentStatusRefunded  PaymentStatus = "refunded"
)

// PaymentMethod represents the payment method
type PaymentMethod string

const (
	PaymentMethodCreditCard PaymentMethod = "credit_card"
	PaymentMethodDebitCard  PaymentMethod = "debit_card"
	PaymentMethodPayPal     PaymentMethod = "paypal"
	PaymentMethodStripe     PaymentMethod = "stripe"
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
	PaymentMethodCrypto     PaymentMethod = "crypto"
)

// PaymentItem represents an item in the payment
type PaymentItem struct {
	ID          string  `json:"id" gorm:"primaryKey"`
	PaymentID   string  `json:"payment_id" gorm:"not null;index"`
	ProductID   int     `json:"product_id" gorm:"not null"`
	Name        string  `json:"name" gorm:"not null"`
	Quantity    int     `json:"quantity" gorm:"not null"`
	Price       float64 `json:"price" gorm:"not null"`
	Subtotal    float64 `json:"subtotal" gorm:"not null"`
	Category    string  `json:"category"`
	CreatedAt   time.Time `json:"created_at"`
}

// IsCompleted checks if payment is completed
func (p *Payment) IsCompleted() bool {
	return p.Status == PaymentStatusCompleted
}

// IsFailed checks if payment failed
func (p *Payment) IsFailed() bool {
	return p.Status == PaymentStatusFailed
}

// IsPending checks if payment is pending
func (p *Payment) IsPending() bool {
	return p.Status == PaymentStatusPending
}

// IsProcessing checks if payment is processing
func (p *Payment) IsProcessing() bool {
	return p.Status == PaymentStatusProcessing
}

// CanBeCancelled checks if payment can be cancelled
func (p *Payment) CanBeCancelled() bool {
	return p.Status == PaymentStatusPending || p.Status == PaymentStatusProcessing
}

// CanBeRefunded checks if payment can be refunded
func (p *Payment) CanBeRefunded() bool {
	return p.Status == PaymentStatusCompleted
}

// MarkAsProcessing marks payment as processing
func (p *Payment) MarkAsProcessing() {
	p.Status = PaymentStatusProcessing
	p.UpdatedAt = time.Now()
}

// MarkAsCompleted marks payment as completed
func (p *Payment) MarkAsCompleted() {
	p.Status = PaymentStatusCompleted
	now := time.Now()
	p.ProcessedAt = &now
	p.UpdatedAt = now
}

// MarkAsFailed marks payment as failed
func (p *Payment) MarkAsFailed() {
	p.Status = PaymentStatusFailed
	p.UpdatedAt = time.Now()
}

// MarkAsCancelled marks payment as cancelled
func (p *Payment) MarkAsCancelled() {
	p.Status = PaymentStatusCancelled
	p.UpdatedAt = time.Now()
}

// MarkAsRefunded marks payment as refunded
func (p *Payment) MarkAsRefunded() {
	p.Status = PaymentStatusRefunded
	p.UpdatedAt = time.Now()
}

// IsExpired checks if payment is expired
func (p *Payment) IsExpired() bool {
	if p.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*p.ExpiresAt)
}

// CalculateTotal calculates the total amount from items
func (p *Payment) CalculateTotal(items []PaymentItem) {
	total := 0.0
	for _, item := range items {
		total += item.Subtotal
	}
	p.Amount = total
	p.UpdatedAt = time.Now()
}

// MarkAsPending marks payment as pending
func (p *Payment) MarkAsPending() {
	p.Status = PaymentStatusPending
	p.UpdatedAt = time.Now()
}

// CanBeRetried checks if payment can be retried
func (p *Payment) CanBeRetried() bool {
	return p.Status == PaymentStatusFailed
}

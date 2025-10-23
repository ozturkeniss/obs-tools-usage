package command

import "obs-tools-usage/internal/payment/application/dto"

// CreatePaymentCommand represents a command to create a payment
type CreatePaymentCommand struct {
	UserID      string            `json:"user_id" binding:"required"`
	BasketID    string            `json:"basket_id" binding:"required"`
	Method      string            `json:"method" binding:"required"`
	Provider    string            `json:"provider" binding:"required"`
	Currency    string            `json:"currency"`
	Description string            `json:"description"`
	Metadata    map[string]string `json:"metadata"`
}

// ToDTO converts command to DTO
func (c *CreatePaymentCommand) ToDTO() dto.CreatePaymentRequest {
	return dto.CreatePaymentRequest{
		UserID:      c.UserID,
		BasketID:    c.BasketID,
		Method:      c.Method,
		Provider:    c.Provider,
		Currency:    c.Currency,
		Description: c.Description,
		Metadata:    c.Metadata,
	}
}

// UpdatePaymentCommand represents a command to update a payment
type UpdatePaymentCommand struct {
	PaymentID string            `json:"payment_id" binding:"required"`
	Status    string            `json:"status" binding:"required"`
	Metadata  map[string]string `json:"metadata"`
}

// ToDTO converts command to DTO
func (c *UpdatePaymentCommand) ToDTO() dto.UpdatePaymentRequest {
	return dto.UpdatePaymentRequest{
		Status:   c.Status,
		Metadata: c.Metadata,
	}
}

// ProcessPaymentCommand represents a command to process a payment
type ProcessPaymentCommand struct {
	PaymentID  string `json:"payment_id" binding:"required"`
	ProviderID string `json:"provider_id"`
}

// ToDTO converts command to DTO
func (c *ProcessPaymentCommand) ToDTO() dto.ProcessPaymentRequest {
	return dto.ProcessPaymentRequest{
		PaymentID:  c.PaymentID,
		ProviderID: c.ProviderID,
	}
}

// RefundPaymentCommand represents a command to refund a payment
type RefundPaymentCommand struct {
	PaymentID string  `json:"payment_id" binding:"required"`
	Amount    float64 `json:"amount"`
	Reason    string  `json:"reason"`
}

// ToDTO converts command to DTO
func (c *RefundPaymentCommand) ToDTO() dto.RefundPaymentRequest {
	return dto.RefundPaymentRequest{
		PaymentID: c.PaymentID,
		Amount:    c.Amount,
		Reason:    c.Reason,
	}
}

// CancelPaymentCommand represents a command to cancel a payment
type CancelPaymentCommand struct {
	PaymentID string `json:"payment_id" binding:"required"`
}

// ToDTO converts command to DTO
func (c *CancelPaymentCommand) ToDTO() dto.CancelPaymentRequest {
	return dto.CancelPaymentRequest{
		PaymentID: c.PaymentID,
	}
}

// RetryPaymentCommand represents a command to retry a payment
type RetryPaymentCommand struct {
	PaymentID string `json:"payment_id" binding:"required"`
}

// ToDTO converts command to DTO
func (c *RetryPaymentCommand) ToDTO() dto.RetryPaymentRequest {
	return dto.RetryPaymentRequest{
		PaymentID: c.PaymentID,
	}
}

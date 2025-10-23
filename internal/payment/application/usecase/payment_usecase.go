package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"obs-tools-usage/internal/payment/application/dto"
	"obs-tools-usage/internal/payment/domain/entity"
	"obs-tools-usage/internal/payment/domain/repository"
	"obs-tools-usage/internal/payment/domain/service"
)

// PaymentUseCase handles payment business logic
type PaymentUseCase struct {
	paymentRepo   repository.PaymentRepository
	basketClient  service.BasketClient
	productClient service.ProductClient
	logger        *logrus.Logger
}

// NewPaymentUseCase creates a new payment use case
func NewPaymentUseCase(paymentRepo repository.PaymentRepository, basketClient service.BasketClient, productClient service.ProductClient, logger *logrus.Logger) *PaymentUseCase {
	return &PaymentUseCase{
		paymentRepo:   paymentRepo,
		basketClient:  basketClient,
		productClient: productClient,
		logger:        logger,
	}
}

// CreatePayment creates a new payment
func (uc *PaymentUseCase) CreatePayment(userID, basketID, method, provider, currency, description string, metadata map[string]string) (*dto.PaymentResponse, error) {
	ctx := context.Background()

	// Get basket information
	basketInfo, err := uc.basketClient.GetBasket(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get basket: %w", err)
	}

	if basketInfo.Total <= 0 {
		return nil, fmt.Errorf("basket is empty or invalid")
	}

	// Generate payment ID
	paymentID := fmt.Sprintf("pay_%s_%d", userID, time.Now().Unix())

	// Create payment entity
	payment := &entity.Payment{
		ID:          paymentID,
		UserID:      userID,
		BasketID:    basketInfo.ID,
		Amount:      basketInfo.Total,
		Currency:    currency,
		Status:      entity.PaymentStatusPending,
		Method:      entity.PaymentMethod(method),
		Provider:    provider,
		Description: description,
		Metadata:    metadata,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Set expiration time (30 minutes from now)
	expiresAt := time.Now().Add(30 * time.Minute)
	payment.ExpiresAt = &expiresAt

	// Create payment in database
	if err := uc.paymentRepo.CreatePayment(payment); err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	// Create payment items from basket
	for _, basketItem := range basketInfo.Items {
		itemID := fmt.Sprintf("item_%s_%d", paymentID, basketItem.ProductID)
		paymentItem := &entity.PaymentItem{
			ID:        itemID,
			PaymentID: paymentID,
			ProductID: basketItem.ProductID,
			Name:      basketItem.Name,
			Quantity:  basketItem.Quantity,
			Price:     basketItem.Price,
			Subtotal:  basketItem.Subtotal,
			Category:  basketItem.Category,
			CreatedAt: time.Now(),
		}

		if err := uc.paymentRepo.CreatePaymentItem(paymentItem); err != nil {
			uc.logger.WithError(err).Error("Failed to create payment item")
			// Continue with other items
		}
	}

	// Convert to response
	response := uc.paymentToResponse(payment)
	
	uc.logger.WithFields(logrus.Fields{
		"payment_id": paymentID,
		"user_id":    userID,
		"amount":     payment.Amount,
		"method":     payment.Method,
	}).Info("Created new payment")

	return response, nil
}

// GetPayment retrieves a payment by ID
func (uc *PaymentUseCase) GetPayment(paymentID string) (*dto.PaymentResponse, error) {
	payment, err := uc.paymentRepo.GetPayment(paymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	// Get payment items
	items, err := uc.paymentRepo.GetPaymentItems(paymentID)
	if err != nil {
		uc.logger.WithError(err).Warn("Failed to get payment items")
	}

	response := uc.paymentToResponse(payment)
	response.Items = uc.itemsToResponse(items)

	return response, nil
}

// UpdatePayment updates payment status
func (uc *PaymentUseCase) UpdatePayment(paymentID, status string, metadata map[string]string) (*dto.PaymentResponse, error) {
	payment, err := uc.paymentRepo.GetPayment(paymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	// Update status
	switch entity.PaymentStatus(status) {
	case entity.PaymentStatusProcessing:
		payment.MarkAsProcessing()
	case entity.PaymentStatusCompleted:
		payment.MarkAsCompleted()
	case entity.PaymentStatusFailed:
		payment.MarkAsFailed()
	case entity.PaymentStatusCancelled:
		payment.MarkAsCancelled()
	case entity.PaymentStatusRefunded:
		payment.MarkAsRefunded()
	default:
		return nil, fmt.Errorf("invalid payment status: %s", status)
	}

	// Update metadata if provided
	if metadata != nil {
		payment.Metadata = metadata
	}

	// Save to database
	if err := uc.paymentRepo.UpdatePayment(payment); err != nil {
		return nil, fmt.Errorf("failed to update payment: %w", err)
	}

	response := uc.paymentToResponse(payment)
	
	uc.logger.WithFields(logrus.Fields{
		"payment_id": paymentID,
		"status":     status,
	}).Info("Updated payment status")

	return response, nil
}

// ProcessPayment processes a payment
func (uc *PaymentUseCase) ProcessPayment(paymentID, providerID string) (*dto.PaymentResponse, error) {
	ctx := context.Background()

	payment, err := uc.paymentRepo.GetPayment(paymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	if !payment.CanBeCancelled() {
		return nil, fmt.Errorf("payment cannot be processed, current status: %s", payment.Status)
	}

	if payment.IsExpired() {
		payment.MarkAsFailed()
		uc.paymentRepo.UpdatePayment(payment)
		return nil, fmt.Errorf("payment has expired")
	}

	// Mark as processing
	payment.MarkAsProcessing()
	payment.ProviderID = providerID
	if err := uc.paymentRepo.UpdatePayment(payment); err != nil {
		return nil, fmt.Errorf("failed to update payment: %w", err)
	}

	// Get payment items for stock update
	items, err := uc.paymentRepo.GetPaymentItems(paymentID)
	if err != nil {
		uc.logger.WithError(err).Warn("Failed to get payment items for stock update")
	}

	// Simulate payment processing (in real implementation, call payment provider)
	time.Sleep(1 * time.Second)

	// For demo purposes, mark as completed
	// In real implementation, this would depend on payment provider response
	payment.MarkAsCompleted()
	if err := uc.paymentRepo.UpdatePayment(payment); err != nil {
		return nil, fmt.Errorf("failed to update payment: %w", err)
	}

	// Update product stock
	for _, item := range items {
		if err := uc.productClient.UpdateProductStock(ctx, item.ProductID, item.Quantity); err != nil {
			uc.logger.WithError(err).WithFields(logrus.Fields{
				"product_id": item.ProductID,
				"quantity":   item.Quantity,
			}).Error("Failed to update product stock")
		}
	}

	// Clear basket after successful payment
	if err := uc.basketClient.ClearBasket(ctx, payment.UserID); err != nil {
		uc.logger.WithError(err).Warn("Failed to clear basket after payment")
	}

	response := uc.paymentToResponse(payment)
	
	uc.logger.WithFields(logrus.Fields{
		"payment_id": paymentID,
		"user_id":    payment.UserID,
		"amount":     payment.Amount,
	}).Info("Payment processed successfully")

	return response, nil
}

// RefundPayment refunds a payment
func (uc *PaymentUseCase) RefundPayment(paymentID string, amount float64, reason string) (*dto.PaymentResponse, error) {
	payment, err := uc.paymentRepo.GetPayment(paymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	if !payment.CanBeRefunded() {
		return nil, fmt.Errorf("payment cannot be refunded, current status: %s", payment.Status)
	}

	// Validate refund amount
	if amount <= 0 {
		amount = payment.Amount // Full refund
	}
	if amount > payment.Amount {
		return nil, fmt.Errorf("refund amount cannot exceed payment amount")
	}

	// Mark as refunded
	payment.MarkAsRefunded()
	if err := uc.paymentRepo.UpdatePayment(payment); err != nil {
		return nil, fmt.Errorf("failed to update payment: %w", err)
	}

	response := uc.paymentToResponse(payment)
	
	uc.logger.WithFields(logrus.Fields{
		"payment_id": paymentID,
		"amount":     amount,
		"reason":     reason,
	}).Info("Payment refunded successfully")

	return response, nil
}

// GetPaymentsByUser retrieves payments by user
func (uc *PaymentUseCase) GetPaymentsByUser(userID string) ([]*dto.PaymentResponse, error) {
	payments, err := uc.paymentRepo.GetPaymentsByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payments by user: %w", err)
	}

	var responses []*dto.PaymentResponse
	for _, payment := range payments {
		items, _ := uc.paymentRepo.GetPaymentItems(payment.ID)
		response := uc.paymentToResponse(payment)
		response.Items = uc.itemsToResponse(items)
		responses = append(responses, response)
	}

	return responses, nil
}

// GetPaymentStats retrieves payment statistics
func (uc *PaymentUseCase) GetPaymentStats(userID string) (*dto.PaymentStatsResponse, error) {
	stats, err := uc.paymentRepo.GetPaymentStats(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment stats: %w", err)
	}

	return &dto.PaymentStatsResponse{
		TotalPayments:     stats.TotalPayments,
		TotalAmount:       stats.TotalAmount,
		CompletedPayments: stats.CompletedPayments,
		FailedPayments:    stats.FailedPayments,
		PendingPayments:   stats.PendingPayments,
		AverageAmount:     stats.AverageAmount,
	}, nil
}

// paymentToResponse converts entity.Payment to dto.PaymentResponse
func (uc *PaymentUseCase) paymentToResponse(payment *entity.Payment) *dto.PaymentResponse {
	return &dto.PaymentResponse{
		ID:          payment.ID,
		UserID:      payment.UserID,
		BasketID:    payment.BasketID,
		Amount:      payment.Amount,
		Currency:    payment.Currency,
		Status:      string(payment.Status),
		Method:      string(payment.Method),
		Provider:    payment.Provider,
		ProviderID:  payment.ProviderID,
		Description: payment.Description,
		Metadata:    payment.Metadata,
		Items:       []dto.PaymentItemResponse{}, // Will be filled separately
		CreatedAt:   payment.CreatedAt,
		UpdatedAt:   payment.UpdatedAt,
		ProcessedAt: payment.ProcessedAt,
		ExpiresAt:   payment.ExpiresAt,
	}
}

// itemsToResponse converts entity.PaymentItem slice to dto.PaymentItemResponse slice
func (uc *PaymentUseCase) itemsToResponse(items []*entity.PaymentItem) []dto.PaymentItemResponse {
	var responses []dto.PaymentItemResponse
	for _, item := range items {
		responses = append(responses, dto.PaymentItemResponse{
			ID:        item.ID,
			ProductID: item.ProductID,
			Name:      item.Name,
			Quantity:  item.Quantity,
			Price:     item.Price,
			Subtotal:  item.Subtotal,
			Category:  item.Category,
			CreatedAt: item.CreatedAt,
		})
	}
	return responses
}

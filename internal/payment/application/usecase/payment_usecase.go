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
	"obs-tools-usage/kafka/events"
	"obs-tools-usage/kafka/publisher"
)

// PaymentUseCase handles payment business logic
type PaymentUseCase struct {
	paymentRepo   repository.PaymentRepository
	basketClient  service.BasketClient
	productClient service.ProductClient
	kafkaPublisher *publisher.PaymentPublisher
	logger        *logrus.Logger
}

// NewPaymentUseCase creates a new payment use case
func NewPaymentUseCase(paymentRepo repository.PaymentRepository, basketClient service.BasketClient, productClient service.ProductClient, kafkaPublisher *publisher.PaymentPublisher, logger *logrus.Logger) *PaymentUseCase {
	return &PaymentUseCase{
		paymentRepo:    paymentRepo,
		basketClient:   basketClient,
		productClient:  productClient,
		kafkaPublisher: kafkaPublisher,
		logger:         logger,
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

	// Publish payment completed event
	paymentCompletedEvent := &events.PaymentCompletedEvent{
		PaymentID: payment.ID,
		UserID:    payment.UserID,
		BasketID:  payment.BasketID,
		Amount:    payment.Amount,
		Currency:  payment.Currency,
		Items:     uc.convertToPaymentItemEvents(items),
		Metadata:  payment.Metadata,
	}

	if err := uc.kafkaPublisher.PublishPaymentCompleted(ctx, paymentCompletedEvent); err != nil {
		uc.logger.WithError(err).Error("Failed to publish payment completed event")
	}

	// Publish stock update events for each item
	for _, item := range items {
		stockUpdateEvent := &events.StockUpdateEvent{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Operation: "decrease",
			Reason:    "Payment completed",
			Metadata: map[string]interface{}{
				"payment_id": payment.ID,
				"user_id":    payment.UserID,
			},
		}

		if err := uc.kafkaPublisher.PublishStockUpdate(ctx, stockUpdateEvent); err != nil {
			uc.logger.WithError(err).WithFields(logrus.Fields{
				"product_id": item.ProductID,
				"quantity":   item.Quantity,
			}).Error("Failed to publish stock update event")
		}
	}

	// Publish basket cleared event
	basketClearedEvent := &events.BasketClearedEvent{
		UserID:   payment.UserID,
		BasketID: payment.BasketID,
		Reason:   "Payment completed",
		Metadata: map[string]interface{}{
			"payment_id": payment.ID,
		},
	}

	if err := uc.kafkaPublisher.PublishBasketCleared(ctx, basketClearedEvent); err != nil {
		uc.logger.WithError(err).Error("Failed to publish basket cleared event")
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

// GetPaymentsByStatus retrieves payments by status
func (uc *PaymentUseCase) GetPaymentsByStatus(status string) ([]*dto.PaymentResponse, error) {
	payments, err := uc.paymentRepo.GetPaymentsByStatus(entity.PaymentStatus(status))
	if err != nil {
		return nil, fmt.Errorf("failed to get payments by status: %w", err)
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

// GetPaymentsByDateRange retrieves payments by date range
func (uc *PaymentUseCase) GetPaymentsByDateRange(startDate, endDate string) ([]*dto.PaymentResponse, error) {
	payments, err := uc.paymentRepo.GetPaymentsByDateRange(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get payments by date range: %w", err)
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

// GetPaymentsByAmountRange retrieves payments by amount range
func (uc *PaymentUseCase) GetPaymentsByAmountRange(minAmount, maxAmount float64) ([]*dto.PaymentResponse, error) {
	payments, err := uc.paymentRepo.GetPaymentsByAmountRange(minAmount, maxAmount)
	if err != nil {
		return nil, fmt.Errorf("failed to get payments by amount range: %w", err)
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

// GetPaymentsByMethod retrieves payments by method
func (uc *PaymentUseCase) GetPaymentsByMethod(method string) ([]*dto.PaymentResponse, error) {
	payments, err := uc.paymentRepo.GetPaymentsByMethod(method)
	if err != nil {
		return nil, fmt.Errorf("failed to get payments by method: %w", err)
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

// GetPaymentsByProvider retrieves payments by provider
func (uc *PaymentUseCase) GetPaymentsByProvider(provider string) ([]*dto.PaymentResponse, error) {
	payments, err := uc.paymentRepo.GetPaymentsByProvider(provider)
	if err != nil {
		return nil, fmt.Errorf("failed to get payments by provider: %w", err)
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

// GetPaymentItems retrieves payment items
func (uc *PaymentUseCase) GetPaymentItems(paymentID string) ([]dto.PaymentItemResponse, error) {
	items, err := uc.paymentRepo.GetPaymentItems(paymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment items: %w", err)
	}

	return uc.itemsToResponse(items), nil
}

// GetPaymentAnalytics retrieves payment analytics
func (uc *PaymentUseCase) GetPaymentAnalytics() (*dto.PaymentAnalyticsResponse, error) {
	analytics, err := uc.paymentRepo.GetPaymentAnalytics()
	if err != nil {
		return nil, fmt.Errorf("failed to get payment analytics: %w", err)
	}

	return &dto.PaymentAnalyticsResponse{
		TotalPayments:     analytics.TotalPayments,
		TotalRevenue:      analytics.TotalRevenue,
		SuccessRate:       analytics.SuccessRate,
		AverageAmount:     analytics.AverageAmount,
		TopPaymentMethod:  analytics.TopPaymentMethod,
		TopProvider:       analytics.TopProvider,
		DailyTransactions: analytics.DailyTransactions,
		MonthlyRevenue:    analytics.MonthlyRevenue,
	}, nil
}

// GetPaymentMethods retrieves available payment methods
func (uc *PaymentUseCase) GetPaymentMethods() (*dto.PaymentMethodsResponse, error) {
	methods, err := uc.paymentRepo.GetPaymentMethods()
	if err != nil {
		return nil, fmt.Errorf("failed to get payment methods: %w", err)
	}

	return &dto.PaymentMethodsResponse{
		Methods: methods,
		Count:   len(methods),
	}, nil
}

// GetPaymentProviders retrieves available payment providers
func (uc *PaymentUseCase) GetPaymentProviders() (*dto.PaymentProvidersResponse, error) {
	providers, err := uc.paymentRepo.GetPaymentProviders()
	if err != nil {
		return nil, fmt.Errorf("failed to get payment providers: %w", err)
	}

	return &dto.PaymentProvidersResponse{
		Providers: providers,
		Count:     len(providers),
	}, nil
}

// GetPaymentSummary retrieves payment summary
func (uc *PaymentUseCase) GetPaymentSummary() (*dto.PaymentSummaryResponse, error) {
	summary, err := uc.paymentRepo.GetPaymentSummary()
	if err != nil {
		return nil, fmt.Errorf("failed to get payment summary: %w", err)
	}

	return &dto.PaymentSummaryResponse{
		TotalPayments:     summary.TotalPayments,
		TotalRevenue:      summary.TotalRevenue,
		PendingPayments:   summary.PendingPayments,
		CompletedPayments: summary.CompletedPayments,
		FailedPayments:    summary.FailedPayments,
		RefundedPayments:  summary.RefundedPayments,
		SuccessRate:       summary.SuccessRate,
		AverageAmount:     summary.AverageAmount,
	}, nil
}

// CancelPayment cancels a payment
func (uc *PaymentUseCase) CancelPayment(paymentID string) (*dto.PaymentResponse, error) {
	payment, err := uc.paymentRepo.GetPayment(paymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	if !payment.CanBeCancelled() {
		return nil, fmt.Errorf("payment cannot be cancelled, current status: %s", payment.Status)
	}

	payment.MarkAsCancelled()
	if err := uc.paymentRepo.UpdatePayment(payment); err != nil {
		return nil, fmt.Errorf("failed to update payment: %w", err)
	}

	response := uc.paymentToResponse(payment)
	
	uc.logger.WithFields(logrus.Fields{
		"payment_id": paymentID,
		"user_id":    payment.UserID,
	}).Info("Payment cancelled successfully")

	return response, nil
}

// RetryPayment retries a failed payment
func (uc *PaymentUseCase) RetryPayment(paymentID string) (*dto.PaymentResponse, error) {
	payment, err := uc.paymentRepo.GetPayment(paymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	if !payment.CanBeRetried() {
		return nil, fmt.Errorf("payment cannot be retried, current status: %s", payment.Status)
	}

	// Reset to pending status for retry
	payment.MarkAsPending()
	if err := uc.paymentRepo.UpdatePayment(payment); err != nil {
		return nil, fmt.Errorf("failed to update payment: %w", err)
	}

	// Process the payment again
	return uc.ProcessPayment(paymentID, "")
}

// convertToPaymentItemEvents converts entity.PaymentItem slice to events.PaymentItemEvent slice
func (uc *PaymentUseCase) convertToPaymentItemEvents(items []*entity.PaymentItem) []events.PaymentItemEvent {
	var eventItems []events.PaymentItemEvent
	for _, item := range items {
		eventItems = append(eventItems, events.PaymentItemEvent{
			ProductID: item.ProductID,
			Name:      item.Name,
			Quantity:  item.Quantity,
			Price:     item.Price,
			Subtotal:  item.Subtotal,
			Category:  item.Category,
		})
	}
	return eventItems
}


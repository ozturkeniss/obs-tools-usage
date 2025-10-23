package handler

import (
	"obs-tools-usage/internal/payment/application/dto"
	"obs-tools-usage/internal/payment/application/query"
	"obs-tools-usage/internal/payment/application/usecase"
)

// QueryHandler handles all queries
type QueryHandler struct {
	paymentUseCase *usecase.PaymentUseCase
}

// NewQueryHandler creates a new query handler
func NewQueryHandler(paymentUseCase *usecase.PaymentUseCase) *QueryHandler {
	return &QueryHandler{
		paymentUseCase: paymentUseCase,
	}
}

// HandleGetPayment handles GetPaymentQuery
func (h *QueryHandler) HandleGetPayment(q query.GetPaymentQuery) (*dto.PaymentResponse, error) {
	return h.paymentUseCase.GetPayment(q.PaymentID)
}

// HandleGetPaymentsByUser handles GetPaymentsByUserQuery
func (h *QueryHandler) HandleGetPaymentsByUser(q query.GetPaymentsByUserQuery) ([]*dto.PaymentResponse, error) {
	return h.paymentUseCase.GetPaymentsByUser(q.UserID)
}

// HandleGetPaymentsByBasket handles GetPaymentsByBasketQuery
func (h *QueryHandler) HandleGetPaymentsByBasket(q query.GetPaymentsByBasketQuery) ([]*dto.PaymentResponse, error) {
	return h.paymentUseCase.GetPaymentsByUser(q.BasketID) // Simplified for now
}

// HandleGetPaymentsByStatus handles GetPaymentsByStatusQuery
func (h *QueryHandler) HandleGetPaymentsByStatus(q query.GetPaymentsByStatusQuery) ([]*dto.PaymentResponse, error) {
	return h.paymentUseCase.GetPaymentsByStatus(q.Status)
}

// HandleGetPaymentStats handles GetPaymentStatsQuery
func (h *QueryHandler) HandleGetPaymentStats(q query.GetPaymentStatsQuery) (*dto.PaymentStatsResponse, error) {
	return h.paymentUseCase.GetPaymentStats(q.UserID)
}

// HandleGetPaymentsByDateRange handles GetPaymentsByDateRangeQuery
func (h *QueryHandler) HandleGetPaymentsByDateRange(q query.GetPaymentsByDateRangeQuery) ([]*dto.PaymentResponse, error) {
	return h.paymentUseCase.GetPaymentsByDateRange(q.StartDate, q.EndDate)
}

// HandleGetPaymentsByAmountRange handles GetPaymentsByAmountRangeQuery
func (h *QueryHandler) HandleGetPaymentsByAmountRange(q query.GetPaymentsByAmountRangeQuery) ([]*dto.PaymentResponse, error) {
	return h.paymentUseCase.GetPaymentsByAmountRange(q.MinAmount, q.MaxAmount)
}

// HandleGetPaymentsByMethod handles GetPaymentsByMethodQuery
func (h *QueryHandler) HandleGetPaymentsByMethod(q query.GetPaymentsByMethodQuery) ([]*dto.PaymentResponse, error) {
	return h.paymentUseCase.GetPaymentsByMethod(q.Method)
}

// HandleGetPaymentsByProvider handles GetPaymentsByProviderQuery
func (h *QueryHandler) HandleGetPaymentsByProvider(q query.GetPaymentsByProviderQuery) ([]*dto.PaymentResponse, error) {
	return h.paymentUseCase.GetPaymentsByProvider(q.Provider)
}

// HandleGetPaymentItems handles GetPaymentItemsQuery
func (h *QueryHandler) HandleGetPaymentItems(q query.GetPaymentItemsQuery) ([]dto.PaymentItemResponse, error) {
	return h.paymentUseCase.GetPaymentItems(q.PaymentID)
}

// HandleGetPaymentAnalytics handles GetPaymentAnalyticsQuery
func (h *QueryHandler) HandleGetPaymentAnalytics(q query.GetPaymentAnalyticsQuery) (*dto.PaymentAnalyticsResponse, error) {
	return h.paymentUseCase.GetPaymentAnalytics()
}

// HandleGetPaymentMethods handles GetPaymentMethodsQuery
func (h *QueryHandler) HandleGetPaymentMethods(q query.GetPaymentMethodsQuery) (*dto.PaymentMethodsResponse, error) {
	return h.paymentUseCase.GetPaymentMethods()
}

// HandleGetPaymentProviders handles GetPaymentProvidersQuery
func (h *QueryHandler) HandleGetPaymentProviders(q query.GetPaymentProvidersQuery) (*dto.PaymentProvidersResponse, error) {
	return h.paymentUseCase.GetPaymentProviders()
}

// HandleGetPaymentSummary handles GetPaymentSummaryQuery
func (h *QueryHandler) HandleGetPaymentSummary(q query.GetPaymentSummaryQuery) (*dto.PaymentSummaryResponse, error) {
	return h.paymentUseCase.GetPaymentSummary()
}

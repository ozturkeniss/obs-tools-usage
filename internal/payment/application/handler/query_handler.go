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
	return h.paymentUseCase.GetPaymentsByUser("") // Simplified for now
}

// HandleGetPaymentStats handles GetPaymentStatsQuery
func (h *QueryHandler) HandleGetPaymentStats(q query.GetPaymentStatsQuery) (*dto.PaymentStatsResponse, error) {
	return h.paymentUseCase.GetPaymentStats(q.UserID)
}

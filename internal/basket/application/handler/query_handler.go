package handler

import (
	"obs-tools-usage/internal/basket/application/dto"
	"obs-tools-usage/internal/basket/application/query"
	"obs-tools-usage/internal/basket/application/usecase"
)

// QueryHandler handles all queries
type QueryHandler struct {
	basketUseCase *usecase.BasketUseCase
}

// NewQueryHandler creates a new query handler
func NewQueryHandler(basketUseCase *usecase.BasketUseCase) *QueryHandler {
	return &QueryHandler{
		basketUseCase: basketUseCase,
	}
}

// HandleGetBasket handles GetBasketQuery
func (h *QueryHandler) HandleGetBasket(q query.GetBasketQuery) (*dto.BasketResponse, error) {
	return h.basketUseCase.GetBasket(q.UserID)
}

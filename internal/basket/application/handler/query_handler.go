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

// HandleGetBasketItems handles GetBasketItemsQuery
func (h *QueryHandler) HandleGetBasketItems(q query.GetBasketItemsQuery) ([]dto.BasketItemResponse, error) {
	return h.basketUseCase.GetBasketItems(q.UserID)
}

// HandleGetBasketTotal handles GetBasketTotalQuery
func (h *QueryHandler) HandleGetBasketTotal(q query.GetBasketTotalQuery) (*dto.BasketTotalResponse, error) {
	return h.basketUseCase.GetBasketTotal(q.UserID)
}

// HandleGetBasketItemCount handles GetBasketItemCountQuery
func (h *QueryHandler) HandleGetBasketItemCount(q query.GetBasketItemCountQuery) (*dto.BasketItemCountResponse, error) {
	return h.basketUseCase.GetBasketItemCount(q.UserID)
}

// HandleGetBasketByCategory handles GetBasketByCategoryQuery
func (h *QueryHandler) HandleGetBasketByCategory(q query.GetBasketByCategoryQuery) ([]dto.BasketItemResponse, error) {
	return h.basketUseCase.GetBasketByCategory(q.UserID, q.Category)
}

// HandleGetBasketStats handles GetBasketStatsQuery
func (h *QueryHandler) HandleGetBasketStats(q query.GetBasketStatsQuery) (*dto.BasketStatsResponse, error) {
	return h.basketUseCase.GetBasketStats(q.UserID)
}

// HandleGetBasketExpiry handles GetBasketExpiryQuery
func (h *QueryHandler) HandleGetBasketExpiry(q query.GetBasketExpiryQuery) (*dto.BasketExpiryResponse, error) {
	return h.basketUseCase.GetBasketExpiry(q.UserID)
}

// HandleGetBasketHistory handles GetBasketHistoryQuery
func (h *QueryHandler) HandleGetBasketHistory(q query.GetBasketHistoryQuery) (*dto.BasketHistoryResponse, error) {
	return h.basketUseCase.GetBasketHistory(q.UserID)
}

// HandleGetBasketRecommendations handles GetBasketRecommendationsQuery
func (h *QueryHandler) HandleGetBasketRecommendations(q query.GetBasketRecommendationsQuery) (*dto.BasketRecommendationsResponse, error) {
	return h.basketUseCase.GetBasketRecommendations(q.UserID)
}

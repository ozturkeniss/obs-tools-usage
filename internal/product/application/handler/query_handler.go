package handler

import (
	"obs-tools-usage/internal/product/application/query"
	"obs-tools-usage/internal/product/application/usecase"
	"obs-tools-usage/internal/product/domain/entity"
)

// QueryHandler handles all queries
type QueryHandler struct {
	productUseCase *usecase.ProductUseCase
}

// NewQueryHandler creates a new query handler
func NewQueryHandler(productUseCase *usecase.ProductUseCase) *QueryHandler {
	return &QueryHandler{
		productUseCase: productUseCase,
	}
}

// HandleGetProduct handles GetProductQuery
func (h *QueryHandler) HandleGetProduct(q query.GetProductQuery) (*entity.Product, error) {
	return h.productUseCase.GetProductByID(q.ID)
}

// HandleGetProducts handles GetProductsQuery
func (h *QueryHandler) HandleGetProducts(q query.GetProductsQuery) ([]entity.Product, error) {
	return h.productUseCase.GetAllProducts()
}

// HandleGetTopMostExpensive handles GetTopMostExpensiveQuery
func (h *QueryHandler) HandleGetTopMostExpensive(q query.GetTopMostExpensiveQuery) ([]entity.Product, error) {
	return h.productUseCase.GetTopMostExpensive(q.Limit)
}

// HandleGetLowStockProducts handles GetLowStockProductsQuery
func (h *QueryHandler) HandleGetLowStockProducts(q query.GetLowStockProductsQuery) ([]entity.Product, error) {
	return h.productUseCase.GetLowStockProducts(q.MaxStock)
}

// HandleGetProductsByCategory handles GetProductsByCategoryQuery
func (h *QueryHandler) HandleGetProductsByCategory(q query.GetProductsByCategoryQuery) ([]entity.Product, error) {
	return h.productUseCase.GetProductsByCategory(q.Category)
}

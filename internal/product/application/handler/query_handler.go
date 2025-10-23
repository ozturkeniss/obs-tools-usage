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

// HandleGetProductsByPriceRange handles GetProductsByPriceRangeQuery
func (h *QueryHandler) HandleGetProductsByPriceRange(q query.GetProductsByPriceRangeQuery) ([]entity.Product, error) {
	return h.productUseCase.GetProductsByPriceRange(q.MinPrice, q.MaxPrice)
}

// HandleGetProductsByName handles GetProductsByNameQuery
func (h *QueryHandler) HandleGetProductsByName(q query.GetProductsByNameQuery) ([]entity.Product, error) {
	return h.productUseCase.GetProductsByName(q.Name)
}

// HandleGetProductStats handles GetProductStatsQuery
func (h *QueryHandler) HandleGetProductStats(q query.GetProductStatsQuery) (*entity.ProductStats, error) {
	return h.productUseCase.GetProductStats()
}

// HandleGetCategories handles GetCategoriesQuery
func (h *QueryHandler) HandleGetCategories(q query.GetCategoriesQuery) ([]entity.Category, error) {
	return h.productUseCase.GetCategories()
}

// HandleGetProductsByStock handles GetProductsByStockQuery
func (h *QueryHandler) HandleGetProductsByStock(q query.GetProductsByStockQuery) ([]entity.Product, error) {
	return h.productUseCase.GetProductsByStock(q.Stock)
}

// HandleGetRandomProducts handles GetRandomProductsQuery
func (h *QueryHandler) HandleGetRandomProducts(q query.GetRandomProductsQuery) ([]entity.Product, error) {
	return h.productUseCase.GetRandomProducts(q.Count)
}

// HandleGetProductsByDateRange handles GetProductsByDateRangeQuery
func (h *QueryHandler) HandleGetProductsByDateRange(q query.GetProductsByDateRangeQuery) ([]entity.Product, error) {
	return h.productUseCase.GetProductsByDateRange(q.StartDate, q.EndDate)
}

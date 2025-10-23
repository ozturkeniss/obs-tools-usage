package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"obs-tools-usage/internal/product/application/command"
	"obs-tools-usage/internal/product/application/dto"
	"obs-tools-usage/internal/product/application/handler"
	"obs-tools-usage/internal/product/application/query"
)

// Handler handles HTTP requests using CQRS pattern
type Handler struct {
	commandHandler *handler.CommandHandler
	queryHandler   *handler.QueryHandler
}

// NewHandler creates a new HTTP handler
func NewHandler(commandHandler *handler.CommandHandler, queryHandler *handler.QueryHandler) *Handler {
	return &Handler{
		commandHandler: commandHandler,
		queryHandler:   queryHandler,
	}
}

// GetAllProducts handles GET /products
func (h *Handler) GetAllProducts(c *gin.Context) {
	products, err := h.queryHandler.HandleGetProducts(query.GetProductsQuery{})
	if err != nil {
		HandleError(c, err)
		return
	}

	response := dto.ProductsResponse{
		Products: make([]dto.ProductResponse, len(products)),
		Count:    len(products),
	}

	for i, product := range products {
		response.Products[i] = dto.ProductResponse{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Stock:       product.Stock,
			Category:    product.Category,
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, response)
}

// GetProductByID handles GET /products/:id
func (h *Handler) GetProductByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid product ID",
			Message: "Product ID must be a valid number",
		})
		return
	}

	product, err := h.queryHandler.HandleGetProduct(query.GetProductQuery{ID: id})
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		Category:    product.Category,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	})
}

// CreateProduct handles POST /products
func (h *Handler) CreateProduct(c *gin.Context) {
	var cmd command.CreateProductCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	product, err := h.commandHandler.HandleCreateProduct(cmd)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		Category:    product.Category,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	})
}

// UpdateProduct handles PUT /products/:id
func (h *Handler) UpdateProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid product ID",
			Message: "Product ID must be a valid number",
		})
		return
	}

	var cmd command.UpdateProductCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	cmd.ID = id

	product, err := h.commandHandler.HandleUpdateProduct(cmd)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		Category:    product.Category,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	})
}

// DeleteProduct handles DELETE /products/:id
func (h *Handler) DeleteProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid product ID",
			Message: "Product ID must be a valid number",
		})
		return
	}

	err = h.commandHandler.HandleDeleteProduct(command.DeleteProductCommand{ID: id})
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Product deleted successfully",
	})
}

// GetTop5MostExpensive handles GET /products/top-5
func (h *Handler) GetTop5MostExpensive(c *gin.Context) {
	products, err := h.queryHandler.HandleGetTopMostExpensive(query.GetTopMostExpensiveQuery{Limit: 5})
	if err != nil {
		HandleError(c, err)
		return
	}

	response := dto.ProductsResponse{
		Products: make([]dto.ProductResponse, len(products)),
		Count:    len(products),
	}

	for i, product := range products {
		response.Products[i] = dto.ProductResponse{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Stock:       product.Stock,
			Category:    product.Category,
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, response)
}

// GetTop10MostExpensive handles GET /products/top-10
func (h *Handler) GetTop10MostExpensive(c *gin.Context) {
	products, err := h.queryHandler.HandleGetTopMostExpensive(query.GetTopMostExpensiveQuery{Limit: 10})
	if err != nil {
		HandleError(c, err)
		return
	}

	response := dto.ProductsResponse{
		Products: make([]dto.ProductResponse, len(products)),
		Count:    len(products),
	}

	for i, product := range products {
		response.Products[i] = dto.ProductResponse{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Stock:       product.Stock,
			Category:    product.Category,
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, response)
}

// GetLowStockProducts1 handles GET /products/low-stock-1
func (h *Handler) GetLowStockProducts1(c *gin.Context) {
	products, err := h.queryHandler.HandleGetLowStockProducts(query.GetLowStockProductsQuery{MaxStock: 1})
	if err != nil {
		HandleError(c, err)
		return
	}

	response := dto.ProductsResponse{
		Products: make([]dto.ProductResponse, len(products)),
		Count:    len(products),
	}

	for i, product := range products {
		response.Products[i] = dto.ProductResponse{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Stock:       product.Stock,
			Category:    product.Category,
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, response)
}

// GetLowStockProducts10 handles GET /products/low-stock-10
func (h *Handler) GetLowStockProducts10(c *gin.Context) {
	products, err := h.queryHandler.HandleGetLowStockProducts(query.GetLowStockProductsQuery{MaxStock: 10})
	if err != nil {
		HandleError(c, err)
		return
	}

	response := dto.ProductsResponse{
		Products: make([]dto.ProductResponse, len(products)),
		Count:    len(products),
	}

	for i, product := range products {
		response.Products[i] = dto.ProductResponse{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Stock:       product.Stock,
			Category:    product.Category,
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, response)
}

// GetProductsByCategory handles GET /products/category/:category
func (h *Handler) GetProductsByCategory(c *gin.Context) {
	category := c.Param("category")
	if category == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid category",
			Message: "Category parameter is required",
		})
		return
	}

	products, err := h.queryHandler.HandleGetProductsByCategory(query.GetProductsByCategoryQuery{Category: category})
	if err != nil {
		HandleError(c, err)
		return
	}

	response := dto.ProductsResponse{
		Products: make([]dto.ProductResponse, len(products)),
		Count:    len(products),
	}

	for i, product := range products {
		response.Products[i] = dto.ProductResponse{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Stock:       product.Stock,
			Category:    product.Category,
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, response)
}

// GetProductsByPriceRange handles GET /products/price/:min/:max
func (h *Handler) GetProductsByPriceRange(c *gin.Context) {
	minPrice, err := strconv.ParseFloat(c.Param("min"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid min price",
			Message: "Min price must be a valid number",
		})
		return
	}

	maxPrice, err := strconv.ParseFloat(c.Param("max"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid max price",
			Message: "Max price must be a valid number",
		})
		return
	}

	products, err := h.queryHandler.HandleGetProductsByPriceRange(query.GetProductsByPriceRangeQuery{
		MinPrice: minPrice,
		MaxPrice: maxPrice,
	})
	if err != nil {
		HandleError(c, err)
		return
	}

	response := dto.ProductsResponse{
		Products: make([]dto.ProductResponse, len(products)),
		Count:    len(products),
	}

	for i, product := range products {
		response.Products[i] = dto.ProductResponse{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Stock:       product.Stock,
			Category:    product.Category,
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, response)
}

// GetProductsByName handles GET /products/search/:name
func (h *Handler) GetProductsByName(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid name",
			Message: "Name parameter is required",
		})
		return
	}

	products, err := h.queryHandler.HandleGetProductsByName(query.GetProductsByNameQuery{Name: name})
	if err != nil {
		HandleError(c, err)
		return
	}

	response := dto.ProductsResponse{
		Products: make([]dto.ProductResponse, len(products)),
		Count:    len(products),
	}

	for i, product := range products {
		response.Products[i] = dto.ProductResponse{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Stock:       product.Stock,
			Category:    product.Category,
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, response)
}

// GetProductStats handles GET /products/stats
func (h *Handler) GetProductStats(c *gin.Context) {
	stats, err := h.queryHandler.HandleGetProductStats(query.GetProductStatsQuery{})
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ProductStatsResponse{
		TotalProducts:     stats.TotalProducts,
		TotalCategories:   stats.TotalCategories,
		AveragePrice:      stats.AveragePrice,
		TotalValue:        stats.TotalValue,
		LowStockProducts:  stats.LowStockProducts,
		OutOfStockProducts: stats.OutOfStockProducts,
	})
}

// GetCategories handles GET /products/categories
func (h *Handler) GetCategories(c *gin.Context) {
	categories, err := h.queryHandler.HandleGetCategories(query.GetCategoriesQuery{})
	if err != nil {
		HandleError(c, err)
		return
	}

	response := dto.CategoriesResponse{
		Categories: make([]dto.CategoryResponse, len(categories)),
		Count:      len(categories),
	}

	for i, category := range categories {
		response.Categories[i] = dto.CategoryResponse{
			Name:        category.Name,
			ProductCount: category.ProductCount,
			AveragePrice: category.AveragePrice,
		}
	}

	c.JSON(http.StatusOK, response)
}

// GetProductsByStock handles GET /products/stock/:stock
func (h *Handler) GetProductsByStock(c *gin.Context) {
	stock, err := strconv.Atoi(c.Param("stock"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid stock",
			Message: "Stock must be a valid number",
		})
		return
	}

	products, err := h.queryHandler.HandleGetProductsByStock(query.GetProductsByStockQuery{Stock: stock})
	if err != nil {
		HandleError(c, err)
		return
	}

	response := dto.ProductsResponse{
		Products: make([]dto.ProductResponse, len(products)),
		Count:    len(products),
	}

	for i, product := range products {
		response.Products[i] = dto.ProductResponse{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Stock:       product.Stock,
			Category:    product.Category,
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, response)
}

// GetRandomProducts handles GET /products/random/:count
func (h *Handler) GetRandomProducts(c *gin.Context) {
	count, err := strconv.Atoi(c.Param("count"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid count",
			Message: "Count must be a valid number",
		})
		return
	}

	if count <= 0 || count > 50 {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid count",
			Message: "Count must be between 1 and 50",
		})
		return
	}

	products, err := h.queryHandler.HandleGetRandomProducts(query.GetRandomProductsQuery{Count: count})
	if err != nil {
		HandleError(c, err)
		return
	}

	response := dto.ProductsResponse{
		Products: make([]dto.ProductResponse, len(products)),
		Count:    len(products),
	}

	for i, product := range products {
		response.Products[i] = dto.ProductResponse{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Stock:       product.Stock,
			Category:    product.Category,
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, response)
}

// GetProductsByDateRange handles GET /products/created/:start/:end
func (h *Handler) GetProductsByDateRange(c *gin.Context) {
	startDate := c.Param("start")
	endDate := c.Param("end")

	if startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid date range",
			Message: "Start and end dates are required",
		})
		return
	}

	products, err := h.queryHandler.HandleGetProductsByDateRange(query.GetProductsByDateRangeQuery{
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		HandleError(c, err)
		return
	}

	response := dto.ProductsResponse{
		Products: make([]dto.ProductResponse, len(products)),
		Count:    len(products),
	}

	for i, product := range products {
		response.Products[i] = dto.ProductResponse{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Stock:       product.Stock,
			Category:    product.Category,
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, response)
}

// HealthCheck handles GET /health
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, dto.HealthResponse{
		Service:   "product-service",
		Status:    "healthy",
		Timestamp: time.Now().Format(time.RFC3339),
		Version:   "1.0.0",
	})
}

// SetupRoutes sets up all routes
func SetupRoutes(r *gin.Engine, commandHandler *handler.CommandHandler, queryHandler *handler.QueryHandler) {
	handler := NewHandler(commandHandler, queryHandler)

	// Product routes
	r.GET("/products", handler.GetAllProducts)
	r.GET("/products/:id", handler.GetProductByID)
	r.POST("/products", handler.CreateProduct)
	r.PUT("/products/:id", handler.UpdateProduct)
	r.DELETE("/products/:id", handler.DeleteProduct)

	// Query routes
	r.GET("/products/top-5", handler.GetTop5MostExpensive)
	r.GET("/products/top-10", handler.GetTop10MostExpensive)
	r.GET("/products/low-stock-1", handler.GetLowStockProducts1)
	r.GET("/products/low-stock-10", handler.GetLowStockProducts10)
	r.GET("/products/category/:category", handler.GetProductsByCategory)
	r.GET("/products/price/:min/:max", handler.GetProductsByPriceRange)
	r.GET("/products/search/:name", handler.GetProductsByName)
	r.GET("/products/stats", handler.GetProductStats)
	r.GET("/products/categories", handler.GetCategories)
	r.GET("/products/stock/:stock", handler.GetProductsByStock)
	r.GET("/products/random/:count", handler.GetRandomProducts)
	r.GET("/products/created/:start/:end", handler.GetProductsByDateRange)

	// Health check
	r.GET("/health", handler.HealthCheck)
}
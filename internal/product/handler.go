package product

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// SetupRoutes configures all product routes
func SetupRoutes(r *gin.Engine, service *Service) {
	handler := NewHandler(service)
	
	// Product routes
	products := r.Group("/products")
	{
		products.GET("", handler.GetAllProducts)
		products.GET("/top-5", handler.GetTop5MostExpensive)
		products.GET("/top-10", handler.GetTop10MostExpensive)
		products.GET("/low-stock-1", handler.GetLowStockProducts1)
		products.GET("/low-stock-10", handler.GetLowStockProducts10)
		products.GET("/:id", handler.GetProductByID)
		products.POST("", handler.CreateProduct)
		products.PUT("/:id", handler.UpdateProduct)
		products.DELETE("/:id", handler.DeleteProduct)
	}
	
	// Health check
	r.GET("/health", handler.HealthCheck)
}

// GetAllProducts returns all products
func (h *Handler) GetAllProducts(c *gin.Context) {
	products, err := h.service.GetAllProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"products": products})
}

// GetProductByID returns a product by ID
func (h *Handler) GetProductByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}
	
	product, err := h.service.GetProductByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, product)
}

// CreateProduct creates a new product
func (h *Handler) CreateProduct(c *gin.Context) {
	var product Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	createdProduct, err := h.service.CreateProduct(product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, createdProduct)
}

// UpdateProduct updates an existing product
func (h *Handler) UpdateProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}
	
	var product Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	product.ID = id
	updatedProduct, err := h.service.UpdateProduct(product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, updatedProduct)
}

// DeleteProduct deletes a product
func (h *Handler) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}
	
	err = h.service.DeleteProduct(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

// GetTop5MostExpensive returns the 5 most expensive products
func (h *Handler) GetTop5MostExpensive(c *gin.Context) {
	products, err := h.service.GetTopMostExpensive(5)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"count": len(products),
		"description": "Top 5 most expensive products",
	})
}

// GetTop10MostExpensive returns the 10 most expensive products
func (h *Handler) GetTop10MostExpensive(c *gin.Context) {
	products, err := h.service.GetTopMostExpensive(10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"count": len(products),
		"description": "Top 10 most expensive products",
	})
}

// GetLowStockProducts1 returns products with stock = 1
func (h *Handler) GetLowStockProducts1(c *gin.Context) {
	products, err := h.service.GetLowStockProducts(1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"count": len(products),
		"description": "Products with stock = 1",
	})
}

// GetLowStockProducts10 returns products with stock < 10
func (h *Handler) GetLowStockProducts10(c *gin.Context) {
	products, err := h.service.GetLowStockProducts(10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"count": len(products),
		"description": "Products with stock < 10",
	})
}

// HealthCheck returns service health status
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"service": "product-service",
	})
}

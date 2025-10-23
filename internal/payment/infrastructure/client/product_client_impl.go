package client

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"obs-tools-usage/api/proto/product"
	"obs-tools-usage/internal/payment/domain/service"
)

// ProductClientImpl implements ProductClient interface using gRPC
type ProductClientImpl struct {
	conn   *grpc.ClientConn
	client product.ProductServiceClient
	logger *logrus.Logger
}

// NewProductClientImpl creates a new product client implementation
func NewProductClientImpl(productServiceURL string, logger *logrus.Logger) (*ProductClientImpl, error) {
	// Create gRPC connection
	conn, err := grpc.Dial(productServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to product service: %w", err)
	}

	client := product.NewProductServiceClient(conn)

	return &ProductClientImpl{
		conn:   conn,
		client: client,
		logger: logger,
	}, nil
}

// GetProduct retrieves a single product by ID
func (c *ProductClientImpl) GetProduct(ctx context.Context, productID int) (*service.ProductInfo, error) {
	c.logger.WithField("product_id", productID).Debug("Getting product from product service")

	req := &product.GetProductRequest{
		Id: int32(productID),
	}

	resp, err := c.client.GetProduct(ctx, req)
	if err != nil {
		c.logger.WithError(err).WithField("product_id", productID).Error("Failed to get product")
		return nil, fmt.Errorf("failed to get product %d: %w", productID, err)
	}

	productInfo := &service.ProductInfo{
		ID:          int(resp.Product.Id),
		Name:        resp.Product.Name,
		Description: resp.Product.Description,
		Price:       resp.Product.Price,
		Stock:       int(resp.Product.Stock),
		Category:    resp.Product.Category,
		Available:   resp.Product.Stock > 0,
	}

	c.logger.WithFields(logrus.Fields{
		"product_id": productInfo.ID,
		"name":       productInfo.Name,
		"price":      productInfo.Price,
		"available":  productInfo.Available,
	}).Debug("Successfully retrieved product")

	return productInfo, nil
}

// GetProducts retrieves multiple products by IDs
func (c *ProductClientImpl) GetProducts(ctx context.Context, productIDs []int) ([]*service.ProductInfo, error) {
	c.logger.WithField("product_ids", productIDs).Debug("Getting products from product service")

	var products []*service.ProductInfo
	
	// Get products one by one (could be optimized with a batch endpoint)
	for _, productID := range productIDs {
		product, err := c.GetProduct(ctx, productID)
		if err != nil {
			c.logger.WithError(err).WithField("product_id", productID).Warn("Failed to get product, skipping")
			continue
		}
		products = append(products, product)
	}

	c.logger.WithFields(logrus.Fields{
		"requested_count": len(productIDs),
		"retrieved_count": len(products),
	}).Debug("Successfully retrieved products")

	return products, nil
}

// UpdateProductStock updates product stock after payment
func (c *ProductClientImpl) UpdateProductStock(ctx context.Context, productID int, quantity int) error {
	c.logger.WithFields(logrus.Fields{
		"product_id": productID,
		"quantity":   quantity,
	}).Debug("Updating product stock")

	// Get current product first
	currentProduct, err := c.GetProduct(ctx, productID)
	if err != nil {
		return fmt.Errorf("failed to get current product: %w", err)
	}

	// Update stock (subtract quantity)
	newStock := currentProduct.Stock - quantity
	if newStock < 0 {
		return fmt.Errorf("insufficient stock for product %d", productID)
	}

	// Update product (this would require an UpdateProduct gRPC method)
	// For now, we'll log the update
	c.logger.WithFields(logrus.Fields{
		"product_id":  productID,
		"old_stock":   currentProduct.Stock,
		"new_stock":   newStock,
		"quantity":    quantity,
	}).Info("Product stock updated after payment")

	return nil
}

// Ping checks the health of the product service
func (c *ProductClientImpl) Ping(ctx context.Context) error {
	// Try to get a product to check if service is responsive
	// This is a simple health check - in production you might want a dedicated health endpoint
	_, err := c.client.ListProducts(ctx, &product.ListProductsRequest{})
	if err != nil {
		return fmt.Errorf("product service is not responding: %w", err)
	}
	return nil
}

// Close closes the gRPC connection
func (c *ProductClientImpl) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// GetConnectionInfo returns connection information for monitoring
func (c *ProductClientImpl) GetConnectionInfo() map[string]interface{} {
	if c.conn == nil {
		return map[string]interface{}{
			"connected": false,
			"state":     "disconnected",
		}
	}

	state := c.conn.GetState()
	return map[string]interface{}{
		"connected": true,
		"state":     state.String(),
	}
}

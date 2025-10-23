package client

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"obs-tools-usage/api/proto/basket"
	"obs-toolsTS/internal/payment/domain/service"
)

// BasketClientImpl implements BasketClient interface using gRPC
type BasketClientImpl struct {
	conn   *grpc.ClientConn
	client basket.BasketServiceClient
	logger *logrus.Logger
}

// NewBasketClientImpl creates a new basket client implementation
func NewBasketClientImpl(basketServiceURL string, logger *logrus.Logger) (*BasketClientImpl, error) {
	// Create gRPC connection
	conn, err := grpc.Dial(basketServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to basket service: %w", err)
	}

	client := basket.NewBasketServiceClient(conn)

	return &BasketClientImpl{
		conn:   conn,
		client: client,
		logger: logger,
	}, nil
}

// GetBasket retrieves basket information
func (c *BasketClientImpl) GetBasket(ctx context.Context, userID string) (*service.BasketInfo, error) {
	c.logger.WithField("user_id", userID).Debug("Getting basket from basket service")

	req := &basket.GetBasketRequest{
		UserId: userID,
	}

	resp, err := c.client.GetBasket(ctx, req)
	if err != nil {
		c.logger.WithError(err).WithField("user_id", userID).Error("Failed to get basket")
		return nil, fmt.Errorf("failed to get basket for user %s: %w", userID, err)
	}

	if !resp.Success {
		return nil, fmt.Errorf("basket service returned error: %s", resp.Message)
	}

	basketInfo := &service.BasketInfo{
		ID:        resp.Basket.Id,
		UserID:    resp.Basket.UserId,
		Total:     resp.Basket.Total,
		ItemCount: int(resp.Basket.ItemCount),
		CreatedAt: resp.Basket.CreatedAt,
		UpdatedAt: resp.Basket.UpdatedAt,
		ExpiresAt: resp.Basket.ExpiresAt,
	}

	// Convert basket items
	for _, item := range resp.Basket.Items {
		basketInfo.Items = append(basketInfo.Items, service.BasketItem{
			ProductID: int(item.ProductId),
			Name:      item.Name,
			Price:     item.Price,
			Quantity:  int(item.Quantity),
			Subtotal:  item.Subtotal,
			Category:  item.Category,
		})
	}

	c.logger.WithFields(logrus.Fields{
		"user_id":    userID,
		"basket_id":  basketInfo.ID,
		"item_count": basketInfo.ItemCount,
		"total":      basketInfo.Total,
	}).Debug("Successfully retrieved basket")

	return basketInfo, nil
}

// ClearBasket clears the basket after successful payment
func (c *BasketClientImpl) ClearBasket(ctx context.Context, userID string) error {
	c.logger.WithField("user_id", userID).Debug("Clearing basket after payment")

	req := &basket.ClearBasketRequest{
		UserId: userID,
	}

	resp, err := c.client.ClearBasket(ctx, req)
	if err != nil {
		c.logger.WithError(err).WithField("user_id", userID).Error("Failed to clear basket")
		return fmt.Errorf("failed to clear basket for user %s: %w", userID, err)
	}

	if !resp.Success {
		return fmt.Errorf("basket service returned error: %s", resp.Message)
	}

	c.logger.WithField("user_id", userID).Info("Successfully cleared basket after payment")
	return nil
}

// Ping checks the health of the basket service
func (c *BasketClientImpl) Ping(ctx context.Context) error {
	req := &basket.HealthCheckRequest{
		Service: "basket-service",
	}

	resp, err := c.client.HealthCheck(ctx, req)
	if err != nil {
		return fmt.Errorf("basket service is not responding: %w", err)
	}

	if !resp.SuccessCompletion {
		return fmt.Errorf("basket service health check failed: %s", resp.Message)
	}

	return nil
}

// Close closes the gRPC connection
func (c *BasketClientImpl) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// GetConnectionInfo returns connection information for monitoring
func (c *BasketClientImpl) GetConnectionInfo() map[string]interface{} {
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

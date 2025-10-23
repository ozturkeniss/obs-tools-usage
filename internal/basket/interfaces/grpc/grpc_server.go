package grpc

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"obs-tools-usage/api/proto/basket"
	"obs-tools-usage/internal/basket/application/command"
	"obs-tools-usage/internal/basket/application/handler"
	"obs-tools-usage/internal/basket/application/query"
	"obs-tools-usage/internal/basket/infrastructure/metrics"
)

// BasketGRPCServer implements the BasketService gRPC server
type BasketGRPCServer struct {
	basket.UnimplementedBasketServiceServer
	commandHandler *handler.CommandHandler
	queryHandler   *handler.QueryHandler
	logger         *logrus.Logger
}

// NewBasketGRPCServer creates a new basket gRPC server
func NewBasketGRPCServer(commandHandler *handler.CommandHandler, queryHandler *handler.QueryHandler, logger *logrus.Logger) *BasketGRPCServer {
	return &BasketGRPCServer{
		commandHandler: commandHandler,
		queryHandler:   queryHandler,
		logger:         logger,
	}
}

// GetBasket retrieves a basket by user ID
func (s *BasketGRPCServer) GetBasket(ctx context.Context, req *basket.GetBasketRequest) (*basket.GetBasketResponse, error) {
	start := time.Now()
	defer metrics.RecordProductServiceRequest("GetBasket", "success", time.Since(start))

	s.logger.WithFields(logrus.Fields{
		"user_id": req.UserId,
	}).Debug("gRPC GetBasket request received")

	// Handle query
	basketResponse, err := s.queryHandler.HandleGetBasket(query.GetBasketQuery{UserID: req.UserId})
	if err != nil {
		s.logger.WithError(err).WithField("user_id", req.UserId).Error("Failed to get basket")
		return &basket.GetBasketResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	// Convert to gRPC response
	grpcBasket := s.convertToGRPCBasket(basketResponse)

	s.logger.WithFields(logrus.Fields{
		"user_id":    req.UserId,
		"item_count": basketResponse.ItemCount,
		"total":      basketResponse.Total,
	}).Info("Successfully retrieved basket via gRPC")

	return &basket.GetBasketResponse{
		Success: true,
		Message: "Basket retrieved successfully",
		Basket:  grpcBasket,
	}, nil
}

// CreateBasket creates a new basket for a user
func (s *BasketGRPCServer) CreateBasket(ctx context.Context, req *basket.CreateBasketRequest) (*basket.CreateBasketResponse, error) {
	start := time.Now()
	defer metrics.RecordProductServiceRequest("CreateBasket", "success", time.Since(start))

	s.logger.WithField("user_id", req.UserId).Debug("gRPC CreateBasket request received")

	// Handle command
	basketResponse, err := s.commandHandler.HandleCreateBasket(command.CreateBasketCommand{UserID: req.UserId})
	if err != nil {
		s.logger.WithError(err).WithField("user_id", req.UserId).Error("Failed to create basket")
		return &basket.CreateBasketResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	// Convert to gRPC response
	grpcBasket := s.convertToGRPCBasket(basketResponse)

	s.logger.WithField("user_id", req.UserId).Info("Successfully created basket via gRPC")

	return &basket.CreateBasketResponse{
		Success: true,
		Message: "Basket created successfully",
		Basket:  grpcBasket,
	}, nil
}

// DeleteBasket deletes a basket
func (s *BasketGRPCServer) DeleteBasket(ctx context.Context, req *basket.DeleteBasketRequest) (*basket.DeleteBasketResponse, error) {
	start := time.Now()
	defer metrics.RecordProductServiceRequest("DeleteBasket", "success", time.Since(start))

	s.logger.WithField("user_id", req.UserId).Debug("gRPC DeleteBasket request received")

	// Handle command
	err := s.commandHandler.HandleDeleteBasket(command.ClearBasketCommand{UserID: req.UserId})
	if err != nil {
		s.logger.WithError(err).WithField("user_id", req.UserId).Error("Failed to delete basket")
		return &basket.DeleteBasketResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	s.logger.WithField("user_id", req.UserId).Info("Successfully deleted basket via gRPC")

	return &basket.DeleteBasketResponse{
		Success: true,
		Message: "Basket deleted successfully",
	}, nil
}

// AddItem adds an item to the basket
func (s *BasketGRPCServer) AddItem(ctx context.Context, req *basket.AddItemRequest) (*basket.AddItemResponse, error) {
	start := time.Now()
	defer metrics.RecordProductServiceRequest("AddItem", "success", time.Since(start))

	s.logger.WithFields(logrus.Fields{
		"user_id":    req.UserId,
		"product_id": req.ProductId,
		"quantity":   req.Quantity,
	}).Debug("gRPC AddItem request received")

	// Handle command
	basketResponse, err := s.commandHandler.HandleAddItem(command.AddItemCommand{
		UserID:    req.UserId,
		ProductID: int(req.ProductId),
		Quantity:  int(req.Quantity),
	})
	if err != nil {
		s.logger.WithError(err).WithFields(logrus.Fields{
			"user_id":    req.UserId,
			"product_id": req.ProductId,
		}).Error("Failed to add item to basket")
		return &basket.AddItemResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	// Convert to gRPC response
	grpcBasket := s.convertToGRPCBasket(basketResponse)

	s.logger.WithFields(logrus.Fields{
		"user_id":    req.UserId,
		"product_id": req.ProductId,
		"quantity":   req.Quantity,
	}).Info("Successfully added item to basket via gRPC")

	return &basket.AddItemResponse{
		Success: true,
		Message: "Item added to basket successfully",
		Basket:  grpcBasket,
	}, nil
}

// UpdateItem updates the quantity of an item in the basket
func (s *BasketGRPCServer) UpdateItem(ctx context.Context, req *basket.UpdateItemRequest) (*basket.UpdateItemResponse, error) {
	start := time.Now()
	defer metrics.RecordProductServiceRequest("UpdateItem", "success", time.Since(start))

	s.logger.WithFields(logrus.Fields{
		"user_id":    req.UserId,
		"product_id": req.ProductId,
		"quantity":   req.Quantity,
	}).Debug("gRPC UpdateItem request received")

	// Handle command
	basketResponse, err := s.commandHandler.HandleUpdateItem(command.UpdateItemCommand{
		UserID:    req.UserId,
		ProductID: int(req.ProductId),
		Quantity:  int(req.Quantity),
	})
	if err != nil {
		s.logger.WithError(err).WithFields(logrus.Fields{
			"user_id":    req.UserId,
			"product_id": req.ProductId,
		}).Error("Failed to update item in basket")
		return &basket.UpdateItemResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	// Convert to gRPC response
	grpcBasket := s.convertToGRPCBasket(basketResponse)

	s.logger.WithFields(logrus.Fields{
		"user_id":    req.UserId,
		"product_id": req.ProductId,
		"quantity":   req.Quantity,
	}).Info("Successfully updated item in basket via gRPC")

	return &basket.UpdateItemResponse{
		Success: true,
		Message: "Item updated in basket successfully",
		Basket:  grpcBasket,
	}, nil
}

// RemoveItem removes an item from the basket
func (s *BasketGRPCServer) RemoveItem(ctx context.Context, req *basket.RemoveItemRequest) (*basket.RemoveItemResponse, error) {
	start := time.Now()
	defer metrics.RecordProductServiceRequest("RemoveItem", "success", time.Since(start))

	s.logger.WithFields(logrus.Fields{
		"user_id":    req.UserId,
		"product_id": req.ProductId,
	}).Debug("gRPC RemoveItem request received")

	// Handle command
	basketResponse, err := s.commandHandler.HandleRemoveItem(command.RemoveItemCommand{
		UserID:    req.UserId,
		ProductID: int(req.ProductId),
	})
	if err != nil {
		s.logger.WithError(err).WithFields(logrus.Fields{
			"user_id":    req.UserId,
			"product_id": req.ProductId,
		}).Error("Failed to remove item from basket")
		return &basket.RemoveItemResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	// Convert to gRPC response
	grpcBasket := s.convertToGRPCBasket(basketResponse)

	s.logger.WithFields(logrus.Fields{
		"user_id":    req.UserId,
		"product_id": req.ProductId,
	}).Info("Successfully removed item from basket via gRPC")

	return &basket.RemoveItemResponse{
		Success: true,
		Message: "Item removed from basket successfully",
		Basket:  grpcBasket,
	}, nil
}

// ClearBasket clears all items from the basket
func (s *BasketGRPCServer) ClearBasket(ctx context.Context, req *basket.ClearBasketRequest) (*basket.ClearBasketResponse, error) {
	start := time.Now()
	defer metrics.RecordProductServiceRequest("ClearBasket", "success", time.Since(start))

	s.logger.WithField("user_id", req.UserId).Debug("gRPC ClearBasket request received")

	// Handle command
	basketResponse, err := s.commandHandler.HandleClearBasket(command.ClearBasketCommand{UserID: req.UserId})
	if err != nil {
		s.logger.WithError(err).WithField("user_id", req.UserId).Error("Failed to clear basket")
		return &basket.ClearBasketResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	// Convert to gRPC response
	grpcBasket := s.convertToGRPCBasket(basketResponse)

	s.logger.WithField("user_id", req.UserId).Info("Successfully cleared basket via gRPC")

	return &basket.ClearBasketResponse{
		Success: true,
		Message: "Basket cleared successfully",
		Basket:  grpcBasket,
	}, nil
}

// HealthCheck performs a health check
func (s *BasketGRPCServer) HealthCheck(ctx context.Context, req *basket.HealthCheckRequest) (*basket.HealthCheckResponse, error) {
	s.logger.Debug("gRPC HealthCheck request received")

	return &basket.HealthCheckResponse{
		Success:   true,
		Message:   "Basket service is healthy",
		Service:   "basket-service",
		Status:    "healthy",
		Timestamp: time.Now().Format(time.RFC3339),
		Version:   "1.0.0",
	}, nil
}

// convertToGRPCBasket converts internal basket response to gRPC basket message
func (s *BasketGRPCServer) convertToGRPCBasket(basketResponse interface{}) *basket.Basket {
	// This is a simplified conversion - in a real implementation you'd properly map the fields
	// For now, return an empty basket as placeholder
	return &basket.Basket{
		Id:        "converted-id",
		UserId:    "converted-user-id",
		Items:     []*basket.BasketItem{},
		Total:     0.0,
		ItemCount: 0,
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
		ExpiresAt: time.Now().Add(24 * time.Hour).Format(time.RFC3339),
	}
}

// RegisterServer registers the gRPC server
func RegisterServer(s *grpc.Server, commandHandler *handler.CommandHandler, queryHandler *handler.QueryHandler, logger *logrus.Logger) {
	basketServer := NewBasketGRPCServer(commandHandler, queryHandler, logger)
	basket.RegisterBasketServiceServer(s, basketServer)
	logger.Info("Basket gRPC server registered")
}

package grpc

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"obs-tools-usage/internal/product/application/command"
	"obs-tools-usage/internal/product/application/handler"
	"obs-tools-usage/internal/product/application/query"
	"obs-tools-usage/internal/product/domain/entity"
	"obs-tools-usage/internal/product/domain/repository"
	"obs-tools-usage/internal/product/infrastructure/config"
	"obs-tools-usage/internal/product/infrastructure/external"

	pb "obs-tools-usage/api/proto/product"
)

// GRPCServer represents the gRPC server
type GRPCServer struct {
	pb.UnimplementedProductServiceServer
	commandHandler *handler.CommandHandler
	queryHandler   *handler.QueryHandler
	repository     repository.ProductRepository
	logger         *logrus.Logger
	grpcServer     *grpc.Server
}

// NewGRPCServer creates a new gRPC server instance
func NewGRPCServer(
	commandHandler *handler.CommandHandler,
	queryHandler *handler.QueryHandler,
	repository repository.ProductRepository,
) *GRPCServer {
	return &GRPCServer{
		commandHandler: commandHandler,
		queryHandler:   queryHandler,
		repository:     repository,
		logger:         config.GetLogger(),
	}
}

// Start starts the gRPC server
func (s *GRPCServer) Start(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	s.grpcServer = grpc.NewServer()
	pb.RegisterProductServiceServer(s.grpcServer, s)
	reflection.Register(s.grpcServer) // Enable reflection for grpcurl

	s.logger.WithField("port", port).Info("Starting gRPC server")
	if err := s.grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}
	return nil
}

// Stop stops the gRPC server
func (s *GRPCServer) Stop() {
	s.logger.Info("Stopping gRPC server...")
	s.grpcServer.GracefulStop()
	s.logger.Info("gRPC server stopped")
}

// GetProduct implements the GetProduct gRPC method
func (s *GRPCServer) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.ProductResponse, error) {
	s.logger.WithField("product_id", req.Id).Debug("GetProduct gRPC request")

	product, err := s.queryHandler.HandleGetProduct(query.GetProductQuery{ID: int(req.Id)})
	if err != nil {
		s.logger.WithError(err).Error("Failed to get product")
		return nil, err
	}

	return &pb.ProductResponse{
		Product: s.productToProto(product),
	}, nil
}

// CreateProduct implements the CreateProduct gRPC method
func (s *GRPCServer) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.ProductResponse, error) {
	s.logger.WithFields(logrus.Fields{
		"name":     req.Name,
		"category": req.Category,
		"price":    req.Price,
	}).Debug("CreateProduct gRPC request")

	cmd := command.CreateProductCommand{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       int(req.Stock),
		Category:    req.Category,
	}

	createdProduct, err := s.commandHandler.HandleCreateProduct(cmd)
	if err != nil {
		s.logger.WithError(err).Error("Failed to create product")
		return nil, err
	}

	// Log business event
	external.LogProductCreated(s.logger.WithField("source", "gRPC"), *createdProduct, "gRPC")

	return &pb.ProductResponse{
		Product: s.productToProto(createdProduct),
	}, nil
}

// UpdateProduct implements the UpdateProduct gRPC method
func (s *GRPCServer) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.ProductResponse, error) {
	s.logger.WithFields(logrus.Fields{
		"product_id": req.Id,
		"name":       req.Name,
		"category":   req.Category,
		"price":      req.Price,
	}).Debug("UpdateProduct gRPC request")

	// Get old product for logging
	oldProduct, _ := s.repository.GetProductByID(int(req.Id))

	cmd := command.UpdateProductCommand{
		ID:          int(req.Id),
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       int(req.Stock),
		Category:    req.Category,
	}

	updatedProduct, err := s.commandHandler.HandleUpdateProduct(cmd)
	if err != nil {
		s.logger.WithError(err).Error("Failed to update product")
		return nil, err
	}

	// Log business event
	if oldProduct != nil {
		external.LogProductUpdated(s.logger.WithField("source", "gRPC"), *oldProduct, *updatedProduct, "gRPC")
	}

	return &pb.ProductResponse{
		Product: s.productToProto(updatedProduct),
	}, nil
}

// DeleteProduct implements the DeleteProduct gRPC method
func (s *GRPCServer) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	s.logger.WithField("product_id", req.Id).Debug("DeleteProduct gRPC request")

	// Get product before deletion for logging
	product, _ := s.repository.GetProductByID(int(req.Id))

	err := s.commandHandler.HandleDeleteProduct(command.DeleteProductCommand{ID: int(req.Id)})
	if err != nil {
		s.logger.WithError(err).Error("Failed to delete product")
		return nil, err
	}

	// Log business event
	if product != nil {
		external.LogProductDeleted(s.logger.WithField("source", "gRPC"), *product, "gRPC")
	}

	return &pb.DeleteProductResponse{
		Message: "Product deleted successfully",
	}, nil
}

// ListProducts implements the ListProducts gRPC method
func (s *GRPCServer) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	s.logger.Debug("ListProducts gRPC request")

	products, err := s.queryHandler.HandleGetProducts(query.GetProductsQuery{})
	if err != nil {
		s.logger.WithError(err).Error("Failed to list products")
		return nil, err
	}

	var protoProducts []*pb.Product
	for _, p := range products {
		protoProducts = append(protoProducts, s.productToProto(&p))
	}

	return &pb.ListProductsResponse{
		Products: protoProducts,
	}, nil
}

// GetTopMostExpensiveProducts implements the GetTopMostExpensiveProducts gRPC method
func (s *GRPCServer) GetTopMostExpensiveProducts(ctx context.Context, req *pb.GetTopMostExpensiveProductsRequest) (*pb.ListProductsResponse, error) {
	s.logger.WithField("limit", req.Limit).Debug("GetTopMostExpensiveProducts gRPC request")

	products, err := s.queryHandler.HandleGetTopMostExpensive(query.GetTopMostExpensiveQuery{Limit: int(req.Limit)})
	if err != nil {
		s.logger.WithError(err).Error("Failed to get top most expensive products")
		return nil, err
	}

	var protoProducts []*pb.Product
	for _, p := range products {
		protoProducts = append(protoProducts, s.productToProto(&p))
	}

	return &pb.ListProductsResponse{
		Products: protoProducts,
	}, nil
}

// GetLowStockProducts implements the GetLowStockProducts gRPC method
func (s *GRPCServer) GetLowStockProducts(ctx context.Context, req *pb.GetLowStockProductsRequest) (*pb.ListProductsResponse, error) {
	s.logger.WithField("max_stock", req.MaxStock).Debug("GetLowStockProducts gRPC request")

	products, err := s.queryHandler.HandleGetLowStockProducts(query.GetLowStockProductsQuery{MaxStock: int(req.MaxStock)})
	if err != nil {
		s.logger.WithError(err).Error("Failed to get low stock products")
		return nil, err
	}

	var protoProducts []*pb.Product
	for _, p := range products {
		protoProducts = append(protoProducts, s.productToProto(&p))
	}

	return &pb.ListProductsResponse{
		Products: protoProducts,
	}, nil
}

// GetProductsByCategory implements the GetProductsByCategory gRPC method
func (s *GRPCServer) GetProductsByCategory(ctx context.Context, req *pb.GetProductsByCategoryRequest) (*pb.ListProductsResponse, error) {
	s.logger.WithField("category", req.Category).Debug("GetProductsByCategory gRPC request")

	products, err := s.queryHandler.HandleGetProductsByCategory(query.GetProductsByCategoryQuery{Category: req.Category})
	if err != nil {
		s.logger.WithError(err).Error("Failed to get products by category")
		return nil, err
	}

	var protoProducts []*pb.Product
	for _, p := range products {
		protoProducts = append(protoProducts, s.productToProto(&p))
	}

	return &pb.ListProductsResponse{
		Products: protoProducts,
	}, nil
}

// productToProto converts an internal Product model to a protobuf Product message
func (s *GRPCServer) productToProto(p *entity.Product) *pb.Product {
	return &pb.Product{
		Id:          int32(p.ID),
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       int32(p.Stock),
		Category:    p.Category,
		CreatedAt:   p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   p.UpdatedAt.Format(time.RFC3339),
	}
}
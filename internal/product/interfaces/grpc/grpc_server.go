package product

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "obs-tools-usage/api/proto/product"
)

// GRPCServer represents the gRPC server
type GRPCServer struct {
	pb.UnimplementedProductServiceServer
	service    *Service
	repository *Repository
	logger     *logrus.Logger
	server     *grpc.Server
}

// NewGRPCServer creates a new gRPC server instance
func NewGRPCServer(service *Service, repository *Repository, logger *logrus.Logger) *GRPCServer {
	return &GRPCServer{
		service:    service,
		repository: repository,
		logger:     logger,
	}
}

// Start starts the gRPC server
func (s *GRPCServer) Start(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen on port %d: %w", port, err)
	}

	// Create gRPC server with options
	s.server = grpc.NewServer(
		grpc.UnaryInterceptor(s.unaryLoggingInterceptor),
	)

	// Register the service
	pb.RegisterProductServiceServer(s.server, s)

	// Enable reflection for debugging
	reflection.Register(s.server)

	s.logger.WithField("port", port).Info("Starting gRPC server")

	// Start serving
	if err := s.server.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve gRPC: %w", err)
	}

	return nil
}

// Stop gracefully stops the gRPC server
func (s *GRPCServer) Stop() {
	if s.server != nil {
		s.logger.Info("Stopping gRPC server...")
		s.server.GracefulStop()
		s.logger.Info("gRPC server stopped")
	}
}

// unaryLoggingInterceptor logs gRPC requests
func (s *GRPCServer) unaryLoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	
	s.logger.WithFields(logrus.Fields{
		"method": info.FullMethod,
		"request": req,
	}).Debug("gRPC request received")

	resp, err := handler(ctx, req)
	
	duration := time.Since(start)
	
	fields := logrus.Fields{
		"method":   info.FullMethod,
		"duration": duration.String(),
	}
	
	if err != nil {
		fields["error"] = err.Error()
		s.logger.WithFields(fields).Error("gRPC request failed")
	} else {
		s.logger.WithFields(fields).Info("gRPC request completed")
	}

	return resp, err
}

// GetProduct implements the GetProduct gRPC method
func (s *GRPCServer) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.ProductResponse, error) {
	s.logger.WithField("product_id", req.Id).Debug("GetProduct gRPC request")

	product, err := s.service.GetProductByID(int(req.Id))
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

	createReq := CreateProductRequest{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       int(req.Stock),
		Category:    req.Category,
	}

	var product Product
	product.FromCreateRequest(createReq)
	
	createdProduct, err := s.service.CreateProduct(product)
	if err != nil {
		s.logger.WithError(err).Error("Failed to create product")
		return nil, err
	}

	// Log business event
	LogProductCreated(s.logger.WithField("source", "gRPC"), *createdProduct, "gRPC")

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

	updateReq := UpdateProductRequest{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       int(req.Stock),
		Category:    req.Category,
	}

	var product Product
	product.ID = int(req.Id)
	product.FromUpdateRequest(updateReq)
	
	updatedProduct, err := s.service.UpdateProduct(product)
	if err != nil {
		s.logger.WithError(err).Error("Failed to update product")
		return nil, err
	}

	// Log business event
	LogProductUpdated(s.logger.WithField("source", "gRPC"), product, *updatedProduct, "gRPC")

	return &pb.ProductResponse{
		Product: s.productToProto(updatedProduct),
	}, nil
}

// DeleteProduct implements the DeleteProduct gRPC method
func (s *GRPCServer) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	s.logger.WithField("product_id", req.Id).Debug("DeleteProduct gRPC request")

	err := s.service.DeleteProduct(int(req.Id))
	if err != nil {
		s.logger.WithError(err).Error("Failed to delete product")
		return nil, err
	}

	// Log business event
	LogProductDeleted(s.logger.WithField("source", "gRPC"), Product{ID: int(req.Id)}, "gRPC")

	return &pb.DeleteProductResponse{
		Message: "Product deleted successfully",
	}, nil
}

// ListProducts implements the ListProducts gRPC method
func (s *GRPCServer) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	s.logger.Debug("ListProducts gRPC request")

	products, err := s.service.GetAllProducts()
	if err != nil {
		s.logger.WithError(err).Error("Failed to list products")
		return nil, err
	}

	protoProducts := make([]*pb.Product, len(products))
	for i, product := range products {
		protoProducts[i] = s.productToProto(&product)
	}

	return &pb.ListProductsResponse{
		Products: protoProducts,
	}, nil
}

// GetTopMostExpensiveProducts implements the GetTopMostExpensiveProducts gRPC method
func (s *GRPCServer) GetTopMostExpensiveProducts(ctx context.Context, req *pb.GetTopMostExpensiveProductsRequest) (*pb.ListProductsResponse, error) {
	s.logger.WithField("limit", req.Limit).Debug("GetTopMostExpensiveProducts gRPC request")

	products, err := s.service.GetTopMostExpensive(int(req.Limit))
	if err != nil {
		s.logger.WithError(err).Error("Failed to get top expensive products")
		return nil, err
	}

	protoProducts := make([]*pb.Product, len(products))
	for i, product := range products {
		protoProducts[i] = s.productToProto(&product)
	}

	return &pb.ListProductsResponse{
		Products: protoProducts,
	}, nil
}

// GetLowStockProducts implements the GetLowStockProducts gRPC method
func (s *GRPCServer) GetLowStockProducts(ctx context.Context, req *pb.GetLowStockProductsRequest) (*pb.ListProductsResponse, error) {
	s.logger.WithField("max_stock", req.MaxStock).Debug("GetLowStockProducts gRPC request")

	products, err := s.service.GetLowStockProducts(int(req.MaxStock))
	if err != nil {
		s.logger.WithError(err).Error("Failed to get low stock products")
		return nil, err
	}

	protoProducts := make([]*pb.Product, len(products))
	for i, product := range products {
		protoProducts[i] = s.productToProto(&product)
	}

	return &pb.ListProductsResponse{
		Products: protoProducts,
	}, nil
}

// GetProductsByCategory implements the GetProductsByCategory gRPC method
func (s *GRPCServer) GetProductsByCategory(ctx context.Context, req *pb.GetProductsByCategoryRequest) (*pb.ListProductsResponse, error) {
	s.logger.WithField("category", req.Category).Debug("GetProductsByCategory gRPC request")

	products, err := s.service.GetProductsByCategory(req.Category)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get products by category")
		return nil, err
	}

	protoProducts := make([]*pb.Product, len(products))
	for i, product := range products {
		protoProducts[i] = s.productToProto(&product)
	}

	return &pb.ListProductsResponse{
		Products: protoProducts,
	}, nil
}

// productToProto converts a Product to protobuf Product
func (s *GRPCServer) productToProto(product *Product) *pb.Product {
	return &pb.Product{
		Id:          int32(product.ID),
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       int32(product.Stock),
		Category:    product.Category,
		CreatedAt:   product.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   product.UpdatedAt.Format(time.RFC3339),
	}
}

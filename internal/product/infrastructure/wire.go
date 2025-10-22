//go:build wireinject
// +build wireinject

package infrastructure

import (
	"obs-tools-usage/internal/product/application/handler"
	"obs-tools-usage/internal/product/application/usecase"
	"obs-tools-usage/internal/product/domain/repository"
	"obs-tools-usage/internal/product/infrastructure/config"
	"obs-tools-usage/internal/product/infrastructure/persistence"
	"obs-tools-usage/internal/product/interfaces/grpc"
	"obs-tools-usage/internal/product/interfaces/http"

	"github.com/google/wire"
	"gorm.io/gorm"
)

// ProviderSet is the provider set for dependency injection
var ProviderSet = wire.NewSet(
	// Config
	config.LoadConfig,
	config.GetLogger,

	// Database
	NewDatabaseProvider,

	// Repository
	NewProductRepositoryProvider,

	// Use Case
	usecase.NewProductUseCase,

	// Handlers
	handler.NewCommandHandler,
	handler.NewQueryHandler,

	// HTTP
	http.NewHandler,
	http.SetupRoutes,

	// gRPC
	grpc.NewGRPCServer,
)

// DatabaseProvider provides database connection
func NewDatabaseProvider(cfg *config.Config) (*persistence.Database, error) {
	return persistence.NewDatabase(&cfg.Database)
}

// ProductRepositoryProvider provides product repository
func NewProductRepositoryProvider(db *gorm.DB) repository.ProductRepository {
	return persistence.NewProductRepositoryImpl(db)
}

// HTTPHandlerProvider provides HTTP handler
func NewHTTPHandlerProvider(
	commandHandler *handler.CommandHandler,
	queryHandler *handler.QueryHandler,
) *http.Handler {
	return http.NewHandler(commandHandler, queryHandler)
}

// GRPCServerProvider provides gRPC server
func NewGRPCServerProvider(
	commandHandler *handler.CommandHandler,
	queryHandler *handler.QueryHandler,
	productRepo repository.ProductRepository,
) *grpc.GRPCServer {
	return grpc.NewGRPCServer(commandHandler, queryHandler, productRepo)
}

//go:build wireinject
// +build wireinject

package main

import (
	"obs-tools-usage/internal/basket/application/handler"
	"obs-tools-usage/internal/basket/application/usecase"
	"obs-tools-usage/internal/basket/domain/repository"
	"obs-tools-usage/internal/basket/domain/service"
	"obs-tools-usage/internal/basket/infrastructure/client"
	"obs-tools-usage/internal/basket/infrastructure/config"
	"obs-tools-usage/internal/basket/infrastructure/persistence"
	httpInterface "obs-tools-usage/internal/basket/interfaces/http"

	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
)

// ProviderSet is the provider set for dependency injection
var ProviderSet = wire.NewSet(
	// Config
	config.LoadConfig,

	// Redis
	NewRedisClient,

	// Product Client
	NewProductClient,

	// Repository
	NewBasketRepository,

	// Use Case
	usecase.NewBasketUseCase,

	// Handlers
	handler.NewCommandHandler,
	handler.NewQueryHandler,

	// HTTP
	httpInterface.NewHandler,
	httpInterface.SetupRoutes,
)

// NewRedisClient provides Redis client
func NewRedisClient(cfg *config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		PoolSize: cfg.Redis.PoolSize,
	})
}

// NewProductClient provides product client
func NewProductClient(cfg *config.Config, redisClient *redis.Client) (service.ProductClient, error) {
	// Note: We need a logger here, but for simplicity we'll use a basic one
	// In a real implementation, you'd inject the logger properly
	return client.NewProductClientImpl(cfg.Product.ServiceURL, nil)
}

// NewBasketRepository provides basket repository
func NewBasketRepository(redisClient *redis.Client) repository.BasketRepository {
	// Note: We need a logger here, but for simplicity we'll use a basic one
	return persistence.NewBasketRepositoryImpl(redisClient, nil)
}

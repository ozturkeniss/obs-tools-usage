package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"obs-tools-usage/internal/basket/application/dto"
	"obs-tools-usage/internal/basket/domain/entity"
	"obs-tools-usage/internal/basket/domain/repository"
	"obs-tools-usage/internal/basket/domain/service"
	"obs-tools-usage/internal/basket/infrastructure/metrics"
)

// BasketUseCase handles basket business logic
type BasketUseCase struct {
	basketRepo    repository.BasketRepository
	productClient service.ProductClient
	logger        *logrus.Logger
}

// NewBasketUseCase creates a new basket use case
func NewBasketUseCase(basketRepo repository.BasketRepository, productClient service.ProductClient, logger *logrus.Logger) *BasketUseCase {
	return &BasketUseCase{
		basketRepo:    basketRepo,
		productClient: productClient,
		logger:        logger,
	}
}

// GetBasket retrieves a basket by user ID
func (uc *BasketUseCase) GetBasket(userID string) (*dto.BasketResponse, error) {
	start := time.Now()
	defer metrics.RecordRedisOperation("GetBasket", "success", time.Since(start))

	basket, err := uc.basketRepo.GetBasket(userID)
	if err != nil {
		metrics.RecordRedisOperation("GetBasket", "error", time.Since(start))
		return nil, fmt.Errorf("failed to get basket: %w", err)
	}

	response := uc.basketToResponse(basket)
	return response, nil
}

// CreateBasket creates a new basket for a user
func (uc *BasketUseCase) CreateBasket(userID string) (*dto.BasketResponse, error) {
	start := time.Now()
	defer metrics.RecordBasketOperation("create_basket")

	// Check if basket already exists
	exists, err := uc.basketRepo.BasketExists(userID)
	if err != nil {
		metrics.RecordRedisOperation("CreateBasket", "error", time.Since(start))
		return nil, fmt.Errorf("failed to check basket existence: %w", err)
	}

	if exists {
		// Return existing basket
		basket, err := uc.basketRepo.GetBasket(userID)
		if err != nil {
			return nil, fmt.Errorf("failed to get existing basket: %w", err)
		}
		return uc.basketToResponse(basket), nil
	}

	// Create new basket
	basket, err := uc.basketRepo.CreateBasket(userID)
	if err != nil {
		metrics.RecordRedisOperation("CreateBasket", "error", time.Since(start))
		return nil, fmt.Errorf("failed to create basket: %w", err)
	}

	metrics.RecordRedisOperation("CreateBasket", "success", time.Since(start))
	response := uc.basketToResponse(basket)
	
	uc.logger.WithField("user_id", userID).Info("Created new basket")
	return response, nil
}

// AddItem adds an item to the basket
func (uc *BasketUseCase) AddItem(userID string, productID int, quantity int) (*dto.BasketResponse, error) {
	start := time.Now()
	defer metrics.RecordBasketOperation("add_item")

	// Get product information from product service
	ctx := context.Background()
	productInfo, err := uc.productClient.GetProduct(ctx, productID)
	if err != nil {
		metrics.RecordProductServiceRequest("GetProduct", "error", time.Since(start))
		return nil, fmt.Errorf("failed to get product information: %w", err)
	}
	metrics.RecordProductServiceRequest("GetProduct", "success", time.Since(start))

	// Check if product is available
	if !productInfo.Available || productInfo.Stock < quantity {
		return nil, fmt.Errorf("product is not available or insufficient stock")
	}

	// Get or create basket
	basket, err := uc.getOrCreateBasket(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get or create basket: %w", err)
	}

	// Add item to basket
	basket.AddItem(productID, productInfo.Name, productInfo.Price, quantity, productInfo.Category)

	// Save basket
	err = uc.basketRepo.UpdateBasket(basket)
	if err != nil {
		metrics.RecordRedisOperation("UpdateBasket", "error", time.Since(start))
		return nil, fmt.Errorf("failed to update basket: %w", err)
	}
	metrics.RecordRedisOperation("UpdateBasket", "success", time.Since(start))

	response := uc.basketToResponse(basket)
	
	uc.logger.WithFields(logrus.Fields{
		"user_id":    userID,
		"product_id": productID,
		"quantity":   quantity,
		"item_count": basket.GetItemCount(),
	}).Info("Added item to basket")

	return response, nil
}

// UpdateItem updates the quantity of an item in the basket
func (uc *BasketUseCase) UpdateItem(userID string, productID int, quantity int) (*dto.BasketResponse, error) {
	start := time.Now()
	defer metrics.RecordBasketOperation("update_item")

	// Get basket
	basket, err := uc.getOrCreateBasket(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get basket: %w", err)
	}

	// Update item quantity
	basket.UpdateItemQuantity(productID, quantity)

	// Save basket
	err = uc.basketRepo.UpdateBasket(basket)
	if err != nil {
		metrics.RecordRedisOperation("UpdateBasket", "error", time.Since(start))
		return nil, fmt.Errorf("failed to update basket: %w", err)
	}
	metrics.RecordRedisOperation("UpdateBasket", "success", time.Since(start))

	response := uc.basketToResponse(basket)
	
	uc.logger.WithFields(logrus.Fields{
		"user_id":    userID,
		"product_id": productID,
		"quantity":   quantity,
	}).Info("Updated item quantity in basket")

	return response, nil
}

// RemoveItem removes an item from the basket
func (uc *BasketUseCase) RemoveItem(userID string, productID int) (*dto.BasketResponse, error) {
	start := time.Now()
	defer metrics.RecordBasketOperation("remove_item")

	// Get basket
	basket, err := uc.getOrCreateBasket(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get basket: %w", err)
	}

	// Remove item
	basket.RemoveItem(productID)

	// Save basket
	err = uc.basketRepo.UpdateBasket(basket)
	if err != nil {
		metrics.RecordRedisOperation("UpdateBasket", "error", time.Since(start))
		return nil, fmt.Errorf("failed to update basket: %w", err)
	}
	metrics.RecordRedisOperation("UpdateBasket", "success", time.Since(start))

	response := uc.basketToResponse(basket)
	
	uc.logger.WithFields(logrus.Fields{
		"user_id":    userID,
		"product_id": productID,
	}).Info("Removed item from basket")

	return response, nil
}

// ClearBasket clears all items from the basket
func (uc *BasketUseCase) ClearBasket(userID string) (*dto.BasketResponse, error) {
	start := time.Now()
	defer metrics.RecordBasketOperation("clear_basket")

	// Get basket
	basket, err := uc.getOrCreateBasket(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get basket: %w", err)
	}

	// Clear basket
	basket.Clear()

	// Save basket
	err = uc.basketRepo.UpdateBasket(basket)
	if err != nil {
		metrics.RecordRedisOperation("UpdateBasket", "error", time.Since(start))
		return nil, fmt.Errorf("failed to update basket: %w", err)
	}
	metrics.RecordRedisOperation("UpdateBasket", "success", time.Since(start))

	response := uc.basketToResponse(basket)
	
	uc.logger.WithField("user_id", userID).Info("Cleared basket")

	return response, nil
}

// DeleteBasket deletes the entire basket
func (uc *BasketUseCase) DeleteBasket(userID string) error {
	start := time.Now()
	defer metrics.RecordBasketOperation("delete_basket")

	err := uc.basketRepo.DeleteBasket(userID)
	if err != nil {
		metrics.RecordRedisOperation("DeleteBasket", "error", time.Since(start))
		return fmt.Errorf("failed to delete basket: %w", err)
	}
	metrics.RecordRedisOperation("DeleteBasket", "success", time.Since(start))

	uc.logger.WithField("user_id", userID).Info("Deleted basket")
	return nil
}

// getOrCreateBasket gets an existing basket or creates a new one
func (uc *BasketUseCase) getOrCreateBasket(userID string) (*entity.Basket, error) {
	// Try to get existing basket
	basket, err := uc.basketRepo.GetBasket(userID)
	if err != nil {
		// If basket doesn't exist, create a new one
		basket, err = uc.basketRepo.CreateBasket(userID)
		if err != nil {
			return nil, fmt.Errorf("failed to create basket: %w", err)
		}
	}
	return basket, nil
}

// basketToResponse converts entity.Basket to dto.BasketResponse
func (uc *BasketUseCase) basketToResponse(basket *entity.Basket) *dto.BasketResponse {
	var items []dto.BasketItemResponse
	for _, item := range basket.Items {
		items = append(items, dto.BasketItemResponse{
			ProductID: item.ProductID,
			Name:      item.Name,
			Price:     item.Price,
			Quantity:  item.Quantity,
			Subtotal:  item.Subtotal,
			Category:  item.Category,
		})
	}

	return &dto.BasketResponse{
		ID:        basket.ID,
		UserID:    basket.UserID,
		Items:     items,
		Total:     basket.Total,
		ItemCount: basket.GetItemCount(),
		CreatedAt: basket.CreatedAt,
		UpdatedAt: basket.UpdatedAt,
		ExpiresAt: basket.ExpiresAt,
	}
}

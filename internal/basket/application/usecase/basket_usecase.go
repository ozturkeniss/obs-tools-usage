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

// GetBasketItems retrieves basket items
func (uc *BasketUseCase) GetBasketItems(userID string) ([]dto.BasketItemResponse, error) {
	start := time.Now()
	defer metrics.RecordRedisOperation("GetBasketItems", "success", time.Since(start))

	basket, err := uc.basketRepo.GetBasket(userID)
	if err != nil {
		metrics.RecordRedisOperation("GetBasketItems", "error", time.Since(start))
		return nil, fmt.Errorf("failed to get basket: %w", err)
	}

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

	return items, nil
}

// GetBasketTotal retrieves basket total
func (uc *BasketUseCase) GetBasketTotal(userID string) (*dto.BasketTotalResponse, error) {
	start := time.Now()
	defer metrics.RecordRedisOperation("GetBasketTotal", "success", time.Since(start))

	basket, err := uc.basketRepo.GetBasket(userID)
	if err != nil {
		metrics.RecordRedisOperation("GetBasketTotal", "error", time.Since(start))
		return nil, fmt.Errorf("failed to get basket: %w", err)
	}

	return &dto.BasketTotalResponse{
		UserID:    userID,
		Total:     basket.Total,
		ItemCount: basket.GetItemCount(),
		Currency:  "USD",
	}, nil
}

// GetBasketItemCount retrieves basket item count
func (uc *BasketUseCase) GetBasketItemCount(userID string) (*dto.BasketItemCountResponse, error) {
	start := time.Now()
	defer metrics.RecordRedisOperation("GetBasketItemCount", "success", time.Since(start))

	basket, err := uc.basketRepo.GetBasket(userID)
	if err != nil {
		metrics.RecordRedisOperation("GetBasketItemCount", "error", time.Since(start))
		return nil, fmt.Errorf("failed to get basket: %w", err)
	}

	return &dto.BasketItemCountResponse{
		UserID:     userID,
		ItemCount:  basket.GetItemCount(),
		UniqueItems: len(basket.Items),
	}, nil
}

// GetBasketByCategory retrieves basket items by category
func (uc *BasketUseCase) GetBasketByCategory(userID, category string) ([]dto.BasketItemResponse, error) {
	start := time.Now()
	defer metrics.RecordRedisOperation("GetBasketByCategory", "success", time.Since(start))

	basket, err := uc.basketRepo.GetBasket(userID)
	if err != nil {
		metrics.RecordRedisOperation("GetBasketByCategory", "error", time.Since(start))
		return nil, fmt.Errorf("failed to get basket: %w", err)
	}

	var items []dto.BasketItemResponse
	for _, item := range basket.Items {
		if item.Category == category {
			items = append(items, dto.BasketItemResponse{
				ProductID: item.ProductID,
				Name:      item.Name,
				Price:     item.Price,
				Quantity:  item.Quantity,
				Subtotal:  item.Subtotal,
				Category:  item.Category,
			})
		}
	}

	return items, nil
}

// GetBasketStats retrieves basket statistics
func (uc *BasketUseCase) GetBasketStats(userID string) (*dto.BasketStatsResponse, error) {
	start := time.Now()
	defer metrics.RecordRedisOperation("GetBasketStats", "success", time.Since(start))

	basket, err := uc.basketRepo.GetBasket(userID)
	if err != nil {
		metrics.RecordRedisOperation("GetBasketStats", "error", time.Since(start))
		return nil, fmt.Errorf("failed to get basket: %w", err)
	}

	totalItems := basket.GetItemCount()
	uniqueItems := len(basket.Items)
	
	var totalValue float64
	var mostExpensive, leastExpensive float64
	categories := make(map[string]bool)
	
	if len(basket.Items) > 0 {
		mostExpensive = basket.Items[0].Price
		leastExpensive = basket.Items[0].Price
	}

	for _, item := range basket.Items {
		totalValue += item.Subtotal
		categories[item.Category] = true
		
		if item.Price > mostExpensive {
			mostExpensive = item.Price
		}
		if item.Price < leastExpensive {
			leastExpensive = item.Price
		}
	}

	averageItemPrice := 0.0
	if uniqueItems > 0 {
		averageItemPrice = totalValue / float64(uniqueItems)
	}

	return &dto.BasketStatsResponse{
		UserID:            userID,
		TotalItems:        totalItems,
		UniqueItems:       uniqueItems,
		TotalValue:        totalValue,
		AverageItemPrice:  averageItemPrice,
		Categories:        len(categories),
		MostExpensiveItem: mostExpensive,
		LeastExpensiveItem: leastExpensive,
	}, nil
}

// GetBasketExpiry retrieves basket expiry information
func (uc *BasketUseCase) GetBasketExpiry(userID string) (*dto.BasketExpiryResponse, error) {
	start := time.Now()
	defer metrics.RecordRedisOperation("GetBasketExpiry", "success", time.Since(start))

	basket, err := uc.basketRepo.GetBasket(userID)
	if err != nil {
		metrics.RecordRedisOperation("GetBasketExpiry", "error", time.Since(start))
		return nil, fmt.Errorf("failed to get basket: %w", err)
	}

	now := time.Now()
	isExpired := now.After(basket.ExpiresAt)
	timeLeft := basket.ExpiresAt.Sub(now)

	return &dto.BasketExpiryResponse{
		UserID:    userID,
		ExpiresAt: basket.ExpiresAt,
		IsExpired: isExpired,
		TimeLeft:  timeLeft.String(),
	}, nil
}

// GetBasketHistory retrieves basket history (simplified)
func (uc *BasketUseCase) GetBasketHistory(userID string) (*dto.BasketHistoryResponse, error) {
	start := time.Now()
	defer metrics.RecordRedisOperation("GetBasketHistory", "success", time.Since(start))

	basket, err := uc.basketRepo.GetBasket(userID)
	if err != nil {
		metrics.RecordRedisOperation("GetBasketHistory", "error", time.Since(start))
		return nil, fmt.Errorf("failed to get basket: %w", err)
	}

	var history []dto.BasketItemResponse
	for _, item := range basket.Items {
		history = append(history, dto.BasketItemResponse{
			ProductID: item.ProductID,
			Name:      item.Name,
			Price:     item.Price,
			Quantity:  item.Quantity,
			Subtotal:  item.Subtotal,
			Category:  item.Category,
		})
	}

	return &dto.BasketHistoryResponse{
		UserID:           userID,
		History:          history,
		TotalOperations:  len(history),
	}, nil
}

// GetBasketRecommendations retrieves basket recommendations (simplified)
func (uc *BasketUseCase) GetBasketRecommendations(userID string) (*dto.BasketRecommendationsResponse, error) {
	start := time.Now()
	defer metrics.RecordRedisOperation("GetBasketRecommendations", "success", time.Since(start))

	// Simplified recommendations - in real implementation, this would use ML or business logic
	recommendations := []dto.BasketItemResponse{
		{
			ProductID: 999,
			Name:      "Recommended Product 1",
			Price:     29.99,
			Quantity:  1,
			Subtotal:  29.99,
			Category:  "Electronics",
		},
		{
			ProductID: 998,
			Name:      "Recommended Product 2",
			Price:     19.99,
			Quantity:  1,
			Subtotal:  19.99,
			Category:  "Accessories",
		},
	}

	return &dto.BasketRecommendationsResponse{
		UserID:         userID,
		Recommendations: recommendations,
		Reason:         "Based on your current basket items",
	}, nil
}

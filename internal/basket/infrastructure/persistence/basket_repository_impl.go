package persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"

	"obs-tools-usage/internal/basket/domain/entity"
	"obs-tools-usage/internal/basket/domain/repository"
)

// BasketRepositoryImpl implements BasketRepository interface using Redis
type BasketRepositoryImpl struct {
	client *redis.Client
	logger *logrus.Logger
}

// NewBasketRepositoryImpl creates a new basket repository implementation
func NewBasketRepositoryImpl(client *redis.Client, logger *logrus.Logger) repository.BasketRepository {
	return &BasketRepositoryImpl{
		client: client,
		logger: logger,
	}
}

// GetBasket retrieves a basket by user ID
func (r *BasketRepositoryImpl) GetBasket(userID string) (*entity.Basket, error) {
	ctx := context.Background()
	
	r.logger.WithField("user_id", userID).Debug("Getting basket from Redis")
	
	data, err := r.client.Get(ctx, r.getBasketKey(userID)).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("basket not found for user %s", userID)
		}
		r.logger.WithError(err).WithField("user_id", userID).Error("Failed to get basket from Redis")
		return nil, fmt.Errorf("failed to get basket: %w", err)
	}

	var basket entity.Basket
	if err := json.Unmarshal([]byte(data), &basket); err != nil {
		r.logger.WithError(err).WithField("user_id", userID).Error("Failed to unmarshal basket data")
		return nil, fmt.Errorf("failed to unmarshal basket data: %w", err)
	}

	// Check if basket is expired
	if basket.IsExpired() {
		r.logger.WithField("user_id", userID).Info("Basket is expired, removing from Redis")
		r.client.Del(ctx, r.getBasketKey(userID))
		return nil, fmt.Errorf("basket is expired")
	}

	r.logger.WithFields(logrus.Fields{
		"user_id":    userID,
		"item_count": basket.GetItemCount(),
		"total":      basket.Total,
	}).Debug("Successfully retrieved basket")

	return &basket, nil
}

// SaveBasket saves a basket to Redis
func (r *BasketRepositoryImpl) SaveBasket(basket *entity.Basket) error {
	ctx := context.Background()
	
	r.logger.WithField("user_id", basket.UserID).Debug("Saving basket to Redis")

	data, err := json.Marshal(basket)
	if err != nil {
		r.logger.WithError(err).WithField("user_id", basket.UserID).Error("Failed to marshal basket data")
		return fmt.Errorf("failed to marshal basket data: %w", err)
	}

	// Calculate TTL (time until expiration)
	ttl := time.Until(basket.ExpiresAt)
	if ttl <= 0 {
		return fmt.Errorf("basket is already expired")
	}

	err = r.client.Set(ctx, r.getBasketKey(basket.UserID), data, ttl).Err()
	if err != nil {
		r.logger.WithError(err).WithField("user_id", basket.UserID).Error("Failed to save basket to Redis")
		return fmt.Errorf("failed to save basket: %w", err)
	}

	r.logger.WithFields(logrus.Fields{
		"user_id":    basket.UserID,
		"item_count": basket.GetItemCount(),
		"total":      basket.Total,
		"ttl":        ttl.String(),
	}).Debug("Successfully saved basket")

	return nil
}

// DeleteBasket deletes a basket from Redis
func (r *BasketRepositoryImpl) DeleteBasket(userID string) error {
	ctx := context.Background()
	
	r.logger.WithField("user_id", userID).Debug("Deleting basket from Redis")

	err := r.client.Del(ctx, r.getBasketKey(userID)).Err()
	if err != nil {
		r.logger.WithError(err).WithField("user_id", userID).Error("Failed to delete basket from Redis")
		return fmt.Errorf("failed to delete basket: %w", err)
	}

	r.logger.WithField("user_id", userID).Debug("Successfully deleted basket")
	return nil
}

// CreateBasket creates a new basket
func (r *BasketRepositoryImpl) CreateBasket(userID string) (*entity.Basket, error) {
	now := time.Now()
	basket := &entity.Basket{
		ID:        fmt.Sprintf("basket_%s_%d", userID, now.Unix()),
		UserID:    userID,
		Items:     []entity.BasketItem{},
		Total:     0.0,
		CreatedAt: now,
		UpdatedAt: now,
		ExpiresAt: now.Add(24 * time.Hour), // 1 day TTL
		Metadata:  make(map[string]string),
	}

	err := r.SaveBasket(basket)
	if err != nil {
		return nil, err
	}

	r.logger.WithField("user_id", userID).Info("Created new basket")
	return basket, nil
}

// UpdateBasket updates an existing basket
func (r *BasketRepositoryImpl) UpdateBasket(basket *entity.Basket) error {
	return r.SaveBasket(basket)
}

// BasketExists checks if a basket exists for the user
func (r *BasketRepositoryImpl) BasketExists(userID string) (bool, error) {
	ctx := context.Background()
	
	exists, err := r.client.Exists(ctx, r.getBasketKey(userID)).Result()
	if err != nil {
		r.logger.WithError(err).WithField("user_id", userID).Error("Failed to check basket existence")
		return false, fmt.Errorf("failed to check basket existence: %w", err)
	}

	return exists > 0, nil
}

// GetAllBaskets retrieves all baskets (for monitoring purposes)
func (r *BasketRepositoryImpl) GetAllBaskets() ([]*entity.Basket, error) {
	ctx := context.Background()
	
	r.logger.Debug("Getting all baskets from Redis")

	keys, err := r.client.Keys(ctx, "basket:*").Result()
	if err != nil {
		r.logger.WithError(err).Error("Failed to get basket keys")
		return nil, fmt.Errorf("failed to get basket keys: %w", err)
	}

	var baskets []*entity.Basket
	for _, key := range keys {
		data, err := r.client.Get(ctx, key).Result()
		if err != nil {
			r.logger.WithError(err).WithField("key", key).Warn("Failed to get basket data, skipping")
			continue
		}

		var basket entity.Basket
		if err := json.Unmarshal([]byte(data), &basket); err != nil {
			r.logger.WithError(err).WithField("key", key).Warn("Failed to unmarshal basket data, skipping")
			continue
		}

		// Skip expired baskets
		if basket.IsExpired() {
			continue
		}

		baskets = append(baskets, &basket)
	}

	r.logger.WithField("count", len(baskets)).Debug("Successfully retrieved all baskets")
	return baskets, nil
}

// ClearExpiredBaskets removes all expired baskets
func (r *BasketRepositoryImpl) ClearExpiredBaskets() error {
	ctx := context.Background()
	
	r.logger.Debug("Clearing expired baskets from Redis")

	keys, err := r.client.Keys(ctx, "basket:*").Result()
	if err != nil {
		r.logger.WithError(err).Error("Failed to get basket keys")
		return fmt.Errorf("failed to get basket keys: %w", err)
	}

	var expiredKeys []string
	now := time.Now()

	for _, key := range keys {
		data, err := r.client.Get(ctx, key).Result()
		if err != nil {
			continue
		}

		var basket entity.Basket
		if err := json.Unmarshal([]byte(data), &basket); err != nil {
			continue
		}

		if now.After(basket.ExpiresAt) {
			expiredKeys = append(expiredKeys, key)
		}
	}

	if len(expiredKeys) > 0 {
		err = r.client.Del(ctx, expiredKeys...).Err()
		if err != nil {
			r.logger.WithError(err).Error("Failed to delete expired baskets")
			return fmt.Errorf("failed to delete expired baskets: %w", err)
		}
	}

	r.logger.WithField("deleted_count", len(expiredKeys)).Info("Successfully cleared expired baskets")
	return nil
}

// Ping checks the Redis connection
func (r *BasketRepositoryImpl) Ping() error {
	ctx := context.Background()
	
	_, err := r.client.Ping(ctx).Result()
	if err != nil {
		r.logger.WithError(err).Error("Redis ping failed")
		return fmt.Errorf("redis ping failed: %w", err)
	}

	return nil
}

// getBasketKey generates the Redis key for a basket
func (r *BasketRepositoryImpl) getBasketKey(userID string) string {
	return fmt.Sprintf("basket:%s", userID)
}

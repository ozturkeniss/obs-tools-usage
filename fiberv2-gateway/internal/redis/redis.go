package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

// Client wraps Redis client with additional functionality
type Client struct {
	client *redis.Client
	logger *logrus.Logger
}

// Config holds Redis configuration
type Config struct {
	Host     string
	Port     string
	Password string
	DB       int
	PoolSize int
	MinIdleConns int
	MaxRetries   int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// NewClient creates a new Redis client
func NewClient(config Config, logger *logrus.Logger) *Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", config.Host, config.Port),
		Password:     config.Password,
		DB:           config.DB,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
		MaxRetries:   config.MaxRetries,
		DialTimeout:  config.DialTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
	})

	return &Client{
		client: rdb,
		logger: logger,
	}
}

// GetClient returns the underlying Redis client
func (c *Client) GetClient() *redis.Client {
	return c.client
}

// Ping tests the connection to Redis
func (c *Client) Ping(ctx context.Context) error {
	pong, err := c.client.Ping(ctx).Result()
	if err != nil {
		c.logger.WithError(err).Error("Failed to ping Redis")
		return fmt.Errorf("failed to ping Redis: %w", err)
	}
	
	c.logger.WithField("response", pong).Debug("Redis ping successful")
	return nil
}

// HealthCheck performs a health check on Redis
func (c *Client) HealthCheck(ctx context.Context) error {
	// Test basic operations
	if err := c.Ping(ctx); err != nil {
		return err
	}
	
	// Test set/get operation
	testKey := "health_check_test"
	testValue := fmt.Sprintf("%d", time.Now().Unix())
	
	err := c.client.Set(ctx, testKey, testValue, time.Minute).Err()
	if err != nil {
		c.logger.WithError(err).Error("Failed to set test key")
		return fmt.Errorf("failed to set test key: %w", err)
	}
	
	result, err := c.client.Get(ctx, testKey).Result()
	if err != nil {
		c.logger.WithError(err).Error("Failed to get test key")
		return fmt.Errorf("failed to get test key: %w", err)
	}
	
	if result != testValue {
		return fmt.Errorf("test key value mismatch: expected %s, got %s", testValue, result)
	}
	
	// Clean up test key
	c.client.Del(ctx, testKey)
	
	c.logger.Debug("Redis health check successful")
	return nil
}

// Close closes the Redis connection
func (c *Client) Close() error {
	return c.client.Close()
}

// GetStats returns Redis connection statistics
func (c *Client) GetStats() *redis.PoolStats {
	return c.client.PoolStats()
}

// Set sets a key-value pair with expiration
func (c *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.client.Set(ctx, key, value, expiration).Err()
}

// Get gets a value by key
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

// Del deletes keys
func (c *Client) Del(ctx context.Context, keys ...string) error {
	return c.client.Del(ctx, keys...).Err()
}

// Exists checks if keys exist
func (c *Client) Exists(ctx context.Context, keys ...string) (int64, error) {
	return c.client.Exists(ctx, keys...).Result()
}

// Expire sets expiration for a key
func (c *Client) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return c.client.Expire(ctx, key, expiration).Err()
}

// TTL returns the time to live for a key
func (c *Client) TTL(ctx context.Context, key string) (time.Duration, error) {
	return c.client.TTL(ctx, key).Result()
}

// Keys returns all keys matching pattern
func (c *Client) Keys(ctx context.Context, pattern string) ([]string, error) {
	return c.client.Keys(ctx, pattern).Result()
}

// Incr increments a key
func (c *Client) Incr(ctx context.Context, key string) (int64, error) {
	return c.client.Incr(ctx, key).Result()
}

// IncrBy increments a key by value
func (c *Client) IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	return c.client.IncrBy(ctx, key, value).Result()
}

// Decr decrements a key
func (c *Client) Decr(ctx context.Context, key string) (int64, error) {
	return c.client.Decr(ctx, key).Result()
}

// DecrBy decrements a key by value
func (c *Client) DecrBy(ctx context.Context, key string, value int64) (int64, error) {
	return c.client.DecrBy(ctx, key, value).Result()
}

// HSet sets hash field
func (c *Client) HSet(ctx context.Context, key string, values ...interface{}) (int64, error) {
	return c.client.HSet(ctx, key, values...).Result()
}

// HGet gets hash field value
func (c *Client) HGet(ctx context.Context, key, field string) (string, error) {
	return c.client.HGet(ctx, key, field).Result()
}

// HGetAll gets all hash fields
func (c *Client) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return c.client.HGetAll(ctx, key).Result()
}

// HDel deletes hash fields
func (c *Client) HDel(ctx context.Context, key string, fields ...string) (int64, error) {
	return c.client.HDel(ctx, key, fields...).Result()
}

// ZAdd adds members to sorted set
func (c *Client) ZAdd(ctx context.Context, key string, members ...*redis.Z) (int64, error) {
	return c.client.ZAdd(ctx, key, members...).Result()
}

// ZRange returns members from sorted set
func (c *Client) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return c.client.ZRange(ctx, key, start, stop).Result()
}

// ZRangeWithScores returns members with scores from sorted set
func (c *Client) ZRangeWithScores(ctx context.Context, key string, start, stop int64) ([]redis.Z, error) {
	return c.client.ZRangeWithScores(ctx, key, start, stop).Result()
}

// ZRem removes members from sorted set
func (c *Client) ZRem(ctx context.Context, key string, members ...interface{}) (int64, error) {
	return c.client.ZRem(ctx, key, members...).Result()
}

// ZRemRangeByScore removes members by score range
func (c *Client) ZRemRangeByScore(ctx context.Context, key, min, max string) (int64, error) {
	return c.client.ZRemRangeByScore(ctx, key, min, max).Result()
}

// ZCard returns the number of members in sorted set
func (c *Client) ZCard(ctx context.Context, key string) (int64, error) {
	return c.client.ZCard(ctx, key).Result()
}

// Eval executes Lua script
func (c *Client) Eval(ctx context.Context, script string, keys []string, args ...interface{}) *redis.Cmd {
	return c.client.Eval(ctx, script, keys, args...)
}

// EvalSha executes Lua script by SHA
func (c *Client) EvalSha(ctx context.Context, sha1 string, keys []string, args ...interface{}) *redis.Cmd {
	return c.client.EvalSha(ctx, sha1, keys, args...)
}

// ScriptLoad loads Lua script
func (c *Client) ScriptLoad(ctx context.Context, script string) (string, error) {
	return c.client.ScriptLoad(ctx, script).Result()
}

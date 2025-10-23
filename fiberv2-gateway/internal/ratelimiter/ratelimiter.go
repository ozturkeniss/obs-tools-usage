package ratelimiter

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

// SlidingWindowRateLimiter implements sliding window rate limiting with Redis
type SlidingWindowRateLimiter struct {
	client *redis.Client
	logger *logrus.Logger
}

// RateLimitConfig holds configuration for rate limiting
type RateLimitConfig struct {
	WindowSize time.Duration // Size of the sliding window
	MaxRequests int          // Maximum requests allowed in the window
	KeyPrefix  string       // Prefix for Redis keys
}

// RateLimitResult holds the result of a rate limit check
type RateLimitResult struct {
	Allowed   bool          // Whether the request is allowed
	Remaining int           // Remaining requests in the window
	ResetTime time.Time     // When the window resets
	RetryAfter time.Duration // How long to wait before retrying (if not allowed)
}

// NewSlidingWindowRateLimiter creates a new sliding window rate limiter
func NewSlidingWindowRateLimiter(redisClient *redis.Client, logger *logrus.Logger) *SlidingWindowRateLimiter {
	return &SlidingWindowRateLimiter{
		client: redisClient,
		logger: logger,
	}
}

// CheckRateLimit checks if a request should be allowed based on sliding window algorithm
func (rl *SlidingWindowRateLimiter) CheckRateLimit(ctx context.Context, config RateLimitConfig, identifier string) (*RateLimitResult, error) {
	now := time.Now()
	windowStart := now.Truncate(config.WindowSize)
	windowEnd := windowStart.Add(config.WindowSize)
	
	// Create Redis key for this identifier and window
	key := fmt.Sprintf("%s:%s:%d", config.KeyPrefix, identifier, windowStart.Unix())
	
	// Use Lua script for atomic operations
	script := `
		local key = KEYS[1]
		local window_start = ARGV[1]
		local window_end = ARGV[2]
		local max_requests = tonumber(ARGV[3])
		local current_time = tonumber(ARGV[4])
		
		-- Remove expired entries
		redis.call('ZREMRANGEBYSCORE', key, '-inf', window_start)
		
		-- Count current requests in window
		local current_count = redis.call('ZCARD', key)
		
		if current_count >= max_requests then
			-- Rate limit exceeded
			return {0, current_count, window_end}
		end
		
		-- Add current request
		redis.call('ZADD', key, current_time, current_time)
		redis.call('EXPIRE', key, math.ceil((window_end - current_time) / 1000))
		
		-- Return success
		return {1, current_count + 1, window_end}
	`
	
	result, err := rl.client.Eval(ctx, script, []string{key}, 
		windowStart.UnixMilli(),
		windowEnd.UnixMilli(),
		config.MaxRequests,
		now.UnixMilli(),
	).Result()
	
	if err != nil {
		rl.logger.WithError(err).Error("Failed to execute rate limit script")
		return nil, fmt.Errorf("failed to check rate limit: %w", err)
	}
	
	// Parse result
	resultSlice := result.([]interface{})
	allowed := resultSlice[0].(int64) == 1
	currentCount := resultSlice[1].(int64)
	resetTimeMs := resultSlice[2].(int64)
	
	remaining := config.MaxRequests - int(currentCount)
	resetTime := time.UnixMilli(resetTimeMs)
	
	var retryAfter time.Duration
	if !allowed {
		retryAfter = time.Until(resetTime)
	}
	
	return &RateLimitResult{
		Allowed:     allowed,
		Remaining:   remaining,
		ResetTime:   resetTime,
		RetryAfter:  retryAfter,
	}, nil
}

// CheckRateLimitWithSlidingWindow implements a more sophisticated sliding window algorithm
func (rl *SlidingWindowRateLimiter) CheckRateLimitWithSlidingWindow(ctx context.Context, config RateLimitConfig, identifier string) (*RateLimitResult, error) {
	now := time.Now()
	
	// Create Redis key for this identifier
	key := fmt.Sprintf("%s:sliding:%s", config.KeyPrefix, identifier)
	
	// Use Lua script for sliding window rate limiting
	script := `
		local key = KEYS[1]
		local window_size_ms = tonumber(ARGV[1])
		local max_requests = tonumber(ARGV[2])
		local current_time_ms = tonumber(ARGV[3])
		local window_start_ms = current_time_ms - window_size_ms
		
		-- Remove expired entries (older than window)
		redis.call('ZREMRANGEBYSCORE', key, '-inf', window_start_ms)
		
		-- Count current requests in window
		local current_count = redis.call('ZCARD', key)
		
		if current_count >= max_requests then
			-- Rate limit exceeded, find the oldest entry to calculate retry time
			local oldest_entries = redis.call('ZRANGE', key, 0, 0, 'WITHSCORES')
			local oldest_time = 0
			if #oldest_entries > 0 then
				oldest_time = tonumber(oldest_entries[2])
			end
			local retry_after_ms = (oldest_time + window_size_ms) - current_time_ms
			if retry_after_ms < 0 then
				retry_after_ms = 0
			end
			return {0, current_count, retry_after_ms}
		end
		
		-- Add current request with unique identifier
		local request_id = current_time_ms .. ':' .. math.random(1000000)
		redis.call('ZADD', key, current_time_ms, request_id)
		
		-- Set expiration for the key
		redis.call('EXPIRE', key, math.ceil(window_size_ms / 1000) + 1)
		
		-- Return success
		return {1, current_count + 1, 0}
	`
	
	result, err := rl.client.Eval(ctx, script, []string{key}, 
		config.WindowSize.Milliseconds(),
		config.MaxRequests,
		now.UnixMilli(),
	).Result()
	
	if err != nil {
		rl.logger.WithError(err).Error("Failed to execute sliding window rate limit script")
		return nil, fmt.Errorf("failed to check sliding window rate limit: %w", err)
	}
	
	// Parse result
	resultSlice := result.([]interface{})
	allowed := resultSlice[0].(int64) == 1
	currentCount := resultSlice[1].(int64)
	retryAfterMs := resultSlice[2].(int64)
	
	remaining := config.MaxRequests - int(currentCount)
	retryAfter := time.Duration(retryAfterMs) * time.Millisecond
	resetTime := now.Add(retryAfter)
	
	return &RateLimitResult{
		Allowed:     allowed,
		Remaining:   remaining,
		ResetTime:   resetTime,
		RetryAfter:  retryAfter,
	}, nil
}

// GetRateLimitStatus gets the current rate limit status for an identifier
func (rl *SlidingWindowRateLimiter) GetRateLimitStatus(ctx context.Context, config RateLimitConfig, identifier string) (*RateLimitResult, error) {
	now := time.Now()
	
	// Create Redis key for this identifier
	key := fmt.Sprintf("%s:sliding:%s", config.KeyPrefix, identifier)
	
	// Use Lua script to get current status
	script := `
		local key = KEYS[1]
		local window_size_ms = tonumber(ARGV[1])
		local max_requests = tonumber(ARGV[2])
		local current_time_ms = tonumber(ARGV[3])
		local window_start_ms = current_time_ms - window_size_ms
		
		-- Remove expired entries
		redis.call('ZREMRANGEBYSCORE', key, '-inf', window_start_ms)
		
		-- Count current requests in window
		local current_count = redis.call('ZCARD', key)
		
		-- Find oldest entry to calculate reset time
		local oldest_entries = redis.call('ZRANGE', key, 0, 0, 'WITHSCORES')
		local oldest_time = 0
		if #oldest_entries > 0 then
			oldest_time = tonumber(oldest_entries[2])
		end
		
		local reset_time_ms = oldest_time + window_size_ms
		if reset_time_ms < current_time_ms then
			reset_time_ms = current_time_ms + window_size_ms
		end
		
		return {current_count, reset_time_ms}
	`
	
	result, err := rl.client.Eval(ctx, script, []string{key}, 
		config.WindowSize.Milliseconds(),
		config.MaxRequests,
		now.UnixMilli(),
	).Result()
	
	if err != nil {
		rl.logger.WithError(err).Error("Failed to get rate limit status")
		return nil, fmt.Errorf("failed to get rate limit status: %w", err)
	}
	
	// Parse result
	resultSlice := result.([]interface{})
	currentCount := resultSlice[0].(int64)
	resetTimeMs := resultSlice[1].(int64)
	
	remaining := config.MaxRequests - int(currentCount)
	resetTime := time.UnixMilli(resetTimeMs)
	allowed := remaining > 0
	
	var retryAfter time.Duration
	if !allowed {
		retryAfter = time.Until(resetTime)
	}
	
	return &RateLimitResult{
		Allowed:     allowed,
		Remaining:   remaining,
		ResetTime:   resetTime,
		RetryAfter:  retryAfter,
	}, nil
}

// ResetRateLimit resets the rate limit for an identifier
func (rl *SlidingWindowRateLimiter) ResetRateLimit(ctx context.Context, config RateLimitConfig, identifier string) error {
	key := fmt.Sprintf("%s:sliding:%s", config.KeyPrefix, identifier)
	
	err := rl.client.Del(ctx, key).Err()
	if err != nil {
		rl.logger.WithError(err).Error("Failed to reset rate limit")
		return fmt.Errorf("failed to reset rate limit: %w", err)
	}
	
	rl.logger.WithField("identifier", identifier).Info("Rate limit reset")
	return nil
}

// GetRateLimitStats gets statistics for rate limiting
func (rl *SlidingWindowRateLimiter) GetRateLimitStats(ctx context.Context, config RateLimitConfig, identifier string) (map[string]interface{}, error) {
	now := time.Now()
	
	// Create Redis key for this identifier
	key := fmt.Sprintf("%s:sliding:%s", config.KeyPrefix, identifier)
	
	// Get current count and oldest entry
	result, err := rl.client.ZRangeWithScores(ctx, key, 0, 0).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get rate limit stats: %w", err)
	}
	
	currentCount := rl.client.ZCard(ctx, key).Val()
	
	var oldestTime time.Time
	if len(result) > 0 {
		oldestTime = time.UnixMilli(int64(result[0].Score))
	}
	
	remaining := config.MaxRequests - int(currentCount)
	allowed := remaining > 0
	
	var resetTime time.Time
	if len(result) > 0 {
		resetTime = oldestTime.Add(config.WindowSize)
	} else {
		resetTime = now.Add(config.WindowSize)
	}
	
	return map[string]interface{}{
		"identifier":    identifier,
		"current_count": currentCount,
		"max_requests":  config.MaxRequests,
		"remaining":     remaining,
		"allowed":       allowed,
		"window_size":   config.WindowSize.String(),
		"reset_time":    resetTime,
		"oldest_request": oldestTime,
	}, nil
}

// CleanupExpiredEntries cleans up expired entries from Redis
func (rl *SlidingWindowRateLimiter) CleanupExpiredEntries(ctx context.Context, config RateLimitConfig) error {
	pattern := fmt.Sprintf("%s:*", config.KeyPrefix)
	
	keys, err := rl.client.Keys(ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("failed to get keys for cleanup: %w", err)
	}
	
	now := time.Now()
	cutoffTime := now.Add(-config.WindowSize).UnixMilli()
	
	for _, key := range keys {
		// Remove expired entries from each key
		err := rl.client.ZRemRangeByScore(ctx, key, "-inf", fmt.Sprintf("%d", cutoffTime)).Err()
		if err != nil {
			rl.logger.WithError(err).WithField("key", key).Error("Failed to cleanup expired entries")
			continue
		}
		
		// If key is empty, delete it
		count, err := rl.client.ZCard(ctx, key).Result()
		if err != nil {
			continue
		}
		
		if count == 0 {
			rl.client.Del(ctx, key)
		}
	}
	
	rl.logger.Info("Rate limit cleanup completed")
	return nil
}

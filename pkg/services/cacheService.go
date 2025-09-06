package services

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// CacheConfig holds configuration for cache client
type CacheConfig struct {
	Addr         string
	Password     string
	DB           int
	MaxRetries   int
	Timeout      time.Duration
	PoolSize     int
	MinIdleConns int
}

// DefaultCacheConfig returns a configuration with sensible defaults
func DefaultCacheConfig() *CacheConfig {
	return &CacheConfig{
		Addr:         "localhost:6379",
		Password:     "",
		DB:           0,
		MaxRetries:   3,
		Timeout:      5 * time.Second,
		PoolSize:     10,
		MinIdleConns: 2,
	}
}

type cacheService struct {
	client *redis.Client
}

type CacheService interface {
	// Basic operations
	GetKey(ctx context.Context, key string) (string, bool, error)
	SetKey(ctx context.Context, key string, value string) error
	SetKeyWithTTL(ctx context.Context, key string, value string, ttl time.Duration) error
	DeleteKey(ctx context.Context, key string) error
	KeyExists(ctx context.Context, key string) (bool, error)

	// Batch operations (useful for redistribution)
	SetBatch(ctx context.Context, keyValues map[string]string) error
	GetBatch(ctx context.Context, keys []string) (map[string]string, error)
	DeleteBatch(ctx context.Context, keys []string) error

	// Utility operations
	Ping(ctx context.Context) error
	GetAllKeys(ctx context.Context, pattern string) ([]string, error)
	FlushDB(ctx context.Context) error
	Close() error
}

// CreateCacheClient creates a new Redis client with the given configuration
func CreateCacheClient(config *CacheConfig) *redis.Client {
	if config == nil {
		config = DefaultCacheConfig()
	}

	client := redis.NewClient(&redis.Options{
		Addr:         config.Addr,
		Password:     config.Password,
		DB:           config.DB,
		MaxRetries:   config.MaxRetries,
		DialTimeout:  config.Timeout,
		ReadTimeout:  config.Timeout,
		WriteTimeout: config.Timeout,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
	})

	return client
}

// NewCacheService creates a new cache service instance
func NewCacheService(client *redis.Client) CacheService {
	return &cacheService{
		client: client,
	}
}

// GetKey retrieves a value by key. Returns (value, found, error)
func (s *cacheService) GetKey(ctx context.Context, key string) (string, bool, error) {
	val, err := s.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", false, nil // Key doesn't exist, not an error
	}
	if err != nil {
		return "", false, fmt.Errorf("failed to get key %s: %w", key, err)
	}
	return val, true, nil
}

// SetKey sets a key-value pair without expiration
func (s *cacheService) SetKey(ctx context.Context, key string, value string) error {
	err := s.client.Set(ctx, key, value, 0).Err()
	if err != nil {
		return fmt.Errorf("failed to set key %s: %w", key, err)
	}
	return nil
}

// SetKeyWithTTL sets a key-value pair with expiration time
func (s *cacheService) SetKeyWithTTL(ctx context.Context, key string, value string, ttl time.Duration) error {
	err := s.client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to set key %s with TTL: %w", key, err)
	}
	return nil
}

// DeleteKey removes a key
func (s *cacheService) DeleteKey(ctx context.Context, key string) error {
	err := s.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete key %s: %w", key, err)
	}
	return nil
}

// KeyExists checks if a key exists
func (s *cacheService) KeyExists(ctx context.Context, key string) (bool, error) {
	result, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check existence of key %s: %w", key, err)
	}
	return result > 0, nil
}

// SetBatch sets multiple key-value pairs in a single transaction
func (s *cacheService) SetBatch(ctx context.Context, keyValues map[string]string) error {
	if len(keyValues) == 0 {
		return nil
	}

	pipe := s.client.Pipeline()
	for key, value := range keyValues {
		pipe.Set(ctx, key, value, 0)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to execute batch set: %w", err)
	}
	return nil
}

// GetBatch retrieves multiple keys at once
func (s *cacheService) GetBatch(ctx context.Context, keys []string) (map[string]string, error) {
	if len(keys) == 0 {
		return make(map[string]string), nil
	}

	pipe := s.client.Pipeline()
	cmds := make([]*redis.StringCmd, len(keys))
	
	for i, key := range keys {
		cmds[i] = pipe.Get(ctx, key)
	}

	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("failed to execute batch get: %w", err)
	}

	result := make(map[string]string)
	for i, cmd := range cmds {
		val, err := cmd.Result()
		if err != redis.Nil && err != nil {
			return nil, fmt.Errorf("failed to get result for key %s: %w", keys[i], err)
		}
		if err != redis.Nil {
			result[keys[i]] = val
		}
	}

	return result, nil
}

// DeleteBatch removes multiple keys in a single transaction
func (s *cacheService) DeleteBatch(ctx context.Context, keys []string) error {
	if len(keys) == 0 {
		return nil
	}

	err := s.client.Del(ctx, keys...).Err()
	if err != nil {
		return fmt.Errorf("failed to delete batch keys: %w", err)
	}
	return nil
}

// Ping checks if the Redis server is reachable
func (s *cacheService) Ping(ctx context.Context) error {
	err := s.client.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("failed to ping Redis server: %w", err)
	}
	return nil
}

// GetAllKeys returns all keys matching the pattern (use with caution in production)
func (s *cacheService) GetAllKeys(ctx context.Context, pattern string) ([]string, error) {
	keys, err := s.client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get keys with pattern %s: %w", pattern, err)
	}
	return keys, nil
}

// FlushDB removes all keys from the current database
func (s *cacheService) FlushDB(ctx context.Context) error {
	err := s.client.FlushDB(ctx).Err()
	if err != nil {
		return fmt.Errorf("failed to flush database: %w", err)
	}
	return nil
}

// Close closes the Redis client connection
func (s *cacheService) Close() error {
	return s.client.Close()
}
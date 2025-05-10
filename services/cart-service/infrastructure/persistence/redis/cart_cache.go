package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"mall-go/services/cart-service/domain/model"
	"mall-go/services/cart-service/domain/repository"
	"mall-go/services/cart-service/infrastructure/config"
)

const (
	// Key prefixes
	cartItemsPrefix = "cart:items:"
	cartCountPrefix = "cart:count:"
	
	// Expiration times
	cartItemsExpiration = 24 * time.Hour
	cartCountExpiration = 24 * time.Hour
)

// RedisCache implements the CartCache interface using Redis
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache creates a new Redis cache
func NewRedisCache(config config.RedisConfig) (repository.CartCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", config.Host, config.Port),
		Password:     config.Password,
		DB:           config.DB,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
	})

	// Ping to verify connection
	ctx := context.Background()
	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisCache{
		client: client,
	}, nil
}

// GetCartItems gets cart items from cache
func (rc *RedisCache) GetCartItems(ctx context.Context, userID string) ([]*model.CartItem, error) {
	key := cartItemsPrefix + userID
	data, err := rc.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var items []*model.CartItem
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, err
	}

	return items, nil
}

// SetCartItems sets cart items in cache
func (rc *RedisCache) SetCartItems(ctx context.Context, userID string, items []*model.CartItem) error {
	key := cartItemsPrefix + userID
	data, err := json.Marshal(items)
	if err != nil {
		return err
	}

	return rc.client.Set(ctx, key, data, cartItemsExpiration).Err()
}

// DeleteCartCache deletes cart cache for a user
func (rc *RedisCache) DeleteCartCache(ctx context.Context, userID string) error {
	itemsKey := cartItemsPrefix + userID
	return rc.client.Del(ctx, itemsKey).Err()
}

// IncrementCartCount increments the cart count for a user
func (rc *RedisCache) IncrementCartCount(ctx context.Context, userID string) error {
	key := cartCountPrefix + userID
	_, err := rc.client.Incr(ctx, key).Result()
	if err != nil {
		return err
	}
	
	// Reset expiration
	return rc.client.Expire(ctx, key, cartCountExpiration).Err()
}

// DecrementCartCount decrements the cart count for a user
func (rc *RedisCache) DecrementCartCount(ctx context.Context, userID string) error {
	key := cartCountPrefix + userID
	
	// First check if the count exists and is greater than 0
	count, err := rc.client.Get(ctx, key).Int()
	if err != nil {
		if err == redis.Nil {
			return nil // Key doesn't exist, nothing to decrement
		}
		return err
	}
	
	if count <= 0 {
		return nil // Count is already 0 or negative, do nothing
	}
	
	_, err = rc.client.Decr(ctx, key).Result()
	if err != nil {
		return err
	}
	
	// Reset expiration
	return rc.client.Expire(ctx, key, cartCountExpiration).Err()
}

// GetCartCount gets the cart count for a user
func (rc *RedisCache) GetCartCount(ctx context.Context, userID string) (int, error) {
	key := cartCountPrefix + userID
	
	count, err := rc.client.Get(ctx, key).Int()
	if err != nil {
		if err == redis.Nil {
			return 0, nil // Key doesn't exist, return 0
		}
		return 0, err
	}
	
	return count, nil
}

// SetCartCount sets the cart count for a user
func (rc *RedisCache) SetCartCount(ctx context.Context, userID string, count int) error {
	key := cartCountPrefix + userID
	return rc.client.Set(ctx, key, count, cartCountExpiration).Err()
}

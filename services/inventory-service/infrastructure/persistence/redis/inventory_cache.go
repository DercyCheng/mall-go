package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"mall-go/services/inventory-service/domain/model"
	"mall-go/services/inventory-service/domain/repository"
	"mall-go/services/inventory-service/infrastructure/config"
)

const (
	// Key prefixes
	inventoryPrefix = "inventory:"
	warehousePrefix = "warehouse:"
	
	// Expiration times
	inventoryExpiration = 1 * time.Hour
	warehouseExpiration = 24 * time.Hour
)

// RedisCache implements the InventoryCache interface using Redis
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache creates a new Redis cache
func NewRedisCache(config config.RedisConfig) (repository.InventoryCache, error) {
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

// GetInventory gets an inventory item from cache
func (r *RedisCache) GetInventory(ctx context.Context, id string) (*model.InventoryItem, error) {
	key := inventoryPrefix + id
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Not found in cache
		}
		return nil, err
	}

	var inventory model.InventoryItem
	if err := json.Unmarshal(data, &inventory); err != nil {
		return nil, err
	}

	return &inventory, nil
}

// SetInventory sets an inventory item in cache
func (r *RedisCache) SetInventory(ctx context.Context, inventory *model.InventoryItem) error {
	key := inventoryPrefix + inventory.ID
	data, err := json.Marshal(inventory)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, data, inventoryExpiration).Err()
}

// DeleteInventory deletes an inventory item from cache
func (r *RedisCache) DeleteInventory(ctx context.Context, id string) error {
	key := inventoryPrefix + id
	return r.client.Del(ctx, key).Err()
}

// GetWarehouse gets a warehouse from cache
func (r *RedisCache) GetWarehouse(ctx context.Context, id string) (*model.Warehouse, error) {
	key := warehousePrefix + id
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Not found in cache
		}
		return nil, err
	}

	var warehouse model.Warehouse
	if err := json.Unmarshal(data, &warehouse); err != nil {
		return nil, err
	}

	return &warehouse, nil
}

// SetWarehouse sets a warehouse in cache
func (r *RedisCache) SetWarehouse(ctx context.Context, warehouse *model.Warehouse) error {
	key := warehousePrefix + warehouse.ID
	data, err := json.Marshal(warehouse)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, data, warehouseExpiration).Err()
}

// DeleteWarehouse deletes a warehouse from cache
func (r *RedisCache) DeleteWarehouse(ctx context.Context, id string) error {
	key := warehousePrefix + id
	return r.client.Del(ctx, key).Err()
}

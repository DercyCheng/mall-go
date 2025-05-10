package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"mall-go/services/payment-service/domain/model"
	"mall-go/services/payment-service/domain/repository"
)

// PaymentCache implements repository.PaymentCache interface using Redis
type PaymentCache struct {
	client *redis.Client
}

// NewPaymentCache creates a new PaymentCache instance
func NewPaymentCache(client *redis.Client) *PaymentCache {
	return &PaymentCache{
		client: client,
	}
}

// getPaymentKey generates a Redis key for a payment by ID
func getPaymentKey(id string) string {
	return fmt.Sprintf("payment:%s", id)
}

// getPaymentOrderKey generates a Redis key for a payment by order ID
func getPaymentOrderKey(orderID string) string {
	return fmt.Sprintf("payment:order:%s", orderID)
}

// SetPayment caches a payment
func (c *PaymentCache) SetPayment(ctx context.Context, payment *model.Payment, ttl int) error {
	data, err := json.Marshal(payment)
	if err != nil {
		return fmt.Errorf("marshal payment error: %w", err)
	}
	
	key := getPaymentKey(payment.ID)
	err = c.client.Set(ctx, key, data, time.Duration(ttl)*time.Second).Err()
	if err != nil {
		return fmt.Errorf("set payment cache error: %w", err)
	}
	
	return nil
}

// GetPayment retrieves a cached payment by ID
func (c *PaymentCache) GetPayment(ctx context.Context, id string) (*model.Payment, error) {
	key := getPaymentKey(id)
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("get payment cache error: %w", err)
	}
	
	var payment model.Payment
	if err := json.Unmarshal(data, &payment); err != nil {
		return nil, fmt.Errorf("unmarshal payment error: %w", err)
	}
	
	return &payment, nil
}

// DeletePayment removes a cached payment
func (c *PaymentCache) DeletePayment(ctx context.Context, id string) error {
	key := getPaymentKey(id)
	err := c.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("delete payment cache error: %w", err)
	}
	
	return nil
}

// SetPaymentByOrderID caches a payment by order ID
func (c *PaymentCache) SetPaymentByOrderID(ctx context.Context, payment *model.Payment, ttl int) error {
	data, err := json.Marshal(payment)
	if err != nil {
		return fmt.Errorf("marshal payment error: %w", err)
	}
	
	key := getPaymentOrderKey(payment.OrderID)
	err = c.client.Set(ctx, key, data, time.Duration(ttl)*time.Second).Err()
	if err != nil {
		return fmt.Errorf("set payment cache by order ID error: %w", err)
	}
	
	return nil
}

// GetPaymentByOrderID retrieves a cached payment by order ID
func (c *PaymentCache) GetPaymentByOrderID(ctx context.Context, orderID string) (*model.Payment, error) {
	key := getPaymentOrderKey(orderID)
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("get payment cache by order ID error: %w", err)
	}
	
	var payment model.Payment
	if err := json.Unmarshal(data, &payment); err != nil {
		return nil, fmt.Errorf("unmarshal payment error: %w", err)
	}
	
	return &payment, nil
}

// DeletePaymentByOrderID removes a cached payment by order ID
func (c *PaymentCache) DeletePaymentByOrderID(ctx context.Context, orderID string) error {
	key := getPaymentOrderKey(orderID)
	err := c.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("delete payment cache by order ID error: %w", err)
	}
	
	return nil
}

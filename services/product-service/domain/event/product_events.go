package event

import (
	"time"

	"github.com/google/uuid"
)

// Event 事件基础接口
type Event interface {
	GetID() string
	GetType() string
	GetOccurredAt() time.Time
	GetPayload() interface{}
}

// BaseEvent 事件基类
type BaseEvent struct {
	ID         string
	Type       string
	OccurredAt time.Time
	Payload    interface{}
}

// GetID 获取事件ID
func (e BaseEvent) GetID() string {
	return e.ID
}

// GetType 获取事件类型
func (e BaseEvent) GetType() string {
	return e.Type
}

// GetOccurredAt 获取事件发生时间
func (e BaseEvent) GetOccurredAt() time.Time {
	return e.OccurredAt
}

// GetPayload 获取事件负载数据
func (e BaseEvent) GetPayload() interface{} {
	return e.Payload
}

// ProductCreatedEvent 商品创建事件
type ProductCreatedEvent struct {
	BaseEvent
}

// NewProductCreatedEvent 创建商品创建事件
func NewProductCreatedEvent(productID, name string, price float64) *ProductCreatedEvent {
	return &ProductCreatedEvent{
		BaseEvent: BaseEvent{
			ID:         uuid.New().String(),
			Type:       "product.created",
			OccurredAt: time.Now(),
			Payload: map[string]interface{}{
				"productID": productID,
				"name":      name,
				"price":     price,
			},
		},
	}
}

// ProductUpdatedEvent 商品更新事件
type ProductUpdatedEvent struct {
	BaseEvent
}

// NewProductUpdatedEvent 创建商品更新事件
func NewProductUpdatedEvent(productID, name string, price float64) *ProductUpdatedEvent {
	return &ProductUpdatedEvent{
		BaseEvent: BaseEvent{
			ID:         uuid.New().String(),
			Type:       "product.updated",
			OccurredAt: time.Now(),
			Payload: map[string]interface{}{
				"productID": productID,
				"name":      name,
				"price":     price,
			},
		},
	}
}

// ProductDeletedEvent 商品删除事件
type ProductDeletedEvent struct {
	BaseEvent
}

// NewProductDeletedEvent 创建商品删除事件
func NewProductDeletedEvent(productID string) *ProductDeletedEvent {
	return &ProductDeletedEvent{
		BaseEvent: BaseEvent{
			ID:         uuid.New().String(),
			Type:       "product.deleted",
			OccurredAt: time.Now(),
			Payload: map[string]interface{}{
				"productID": productID,
			},
		},
	}
}

// ProductStockChangedEvent 商品库存变更事件
type ProductStockChangedEvent struct {
	BaseEvent
}

// NewProductStockChangedEvent 创建商品库存变更事件
func NewProductStockChangedEvent(productID string, oldStock, newStock int) *ProductStockChangedEvent {
	return &ProductStockChangedEvent{
		BaseEvent: BaseEvent{
			ID:         uuid.New().String(),
			Type:       "product.stock_changed",
			OccurredAt: time.Now(),
			Payload: map[string]interface{}{
				"productID": productID,
				"oldStock":  oldStock,
				"newStock":  newStock,
			},
		},
	}
}

// ProductPriceChangedEvent 商品价格变更事件
type ProductPriceChangedEvent struct {
	BaseEvent
}

// NewProductPriceChangedEvent 创建商品价格变更事件
func NewProductPriceChangedEvent(productID string, oldPrice, newPrice float64) *ProductPriceChangedEvent {
	return &ProductPriceChangedEvent{
		BaseEvent: BaseEvent{
			ID:         uuid.New().String(),
			Type:       "product.price_changed",
			OccurredAt: time.Now(),
			Payload: map[string]interface{}{
				"productID": productID,
				"oldPrice":  oldPrice,
				"newPrice":  newPrice,
			},
		},
	}
}

// ProductPublishedEvent 商品发布事件
type ProductPublishedEvent struct {
	BaseEvent
}

// NewProductPublishedEvent 创建商品发布事件
func NewProductPublishedEvent(productID, name string) *ProductPublishedEvent {
	return &ProductPublishedEvent{
		BaseEvent: BaseEvent{
			ID:         uuid.New().String(),
			Type:       "product.published",
			OccurredAt: time.Now(),
			Payload: map[string]interface{}{
				"productID": productID,
				"name":      name,
			},
		},
	}
}

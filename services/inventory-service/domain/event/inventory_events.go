package event

import (
	"time"

	"mall-go/services/inventory-service/domain/model"
)

// InventoryEvent 库存事件基础结构
type InventoryEvent struct {
	ID          string    `json:"id"`
	ProductID   string    `json:"productId"`
	SKU         string    `json:"sku"`
	EventType   string    `json:"eventType"`
	OccurredAt  time.Time `json:"occurredAt"`
	PublishedAt time.Time `json:"publishedAt"`
}

// InventoryCreatedEvent 库存创建事件
type InventoryCreatedEvent struct {
	InventoryEvent
	InitialQuantity int    `json:"initialQuantity"`
	WarehouseID     string `json:"warehouseId"`
}

// InventoryUpdatedEvent 库存更新事件
type InventoryUpdatedEvent struct {
	InventoryEvent
	AvailableQuantity int    `json:"availableQuantity"`
	ReservedQuantity  int    `json:"reservedQuantity"`
	Status            string `json:"status"`
}

// InventoryQuantityChangedEvent 库存数量变更事件
type InventoryQuantityChangedEvent struct {
	InventoryEvent
	OldQuantity    int    `json:"oldQuantity"`
	NewQuantity    int    `json:"newQuantity"`
	ChangeQuantity int    `json:"changeQuantity"`
	OperationType  string `json:"operationType"`
	Reason         string `json:"reason"`
	RelatedOrderID string `json:"relatedOrderId,omitempty"`
}

// InventoryLowStockEvent 库存不足事件
type InventoryLowStockEvent struct {
	InventoryEvent
	CurrentQuantity   int `json:"currentQuantity"`
	LowStockThreshold int `json:"lowStockThreshold"`
}

// InventoryOutOfStockEvent 库存不足事件
type InventoryOutOfStockEvent struct {
	InventoryEvent
	ProductName string `json:"productName,omitempty"`
}

// InventoryLockedEvent 库存锁定事件
type InventoryLockedEvent struct {
	InventoryEvent
	Reason string `json:"reason"`
}

// InventoryUnlockedEvent 库存解锁事件
type InventoryUnlockedEvent struct {
	InventoryEvent
	Reason string `json:"reason"`
}

// NewInventoryCreatedEvent 创建库存创建事件
func NewInventoryCreatedEvent(inventory *model.Inventory) *InventoryCreatedEvent {
	return &InventoryCreatedEvent{
		InventoryEvent: InventoryEvent{
			ID:          inventory.ID,
			ProductID:   inventory.ProductID,
			SKU:         inventory.SKU,
			EventType:   "inventory.created",
			OccurredAt:  inventory.CreatedAt,
			PublishedAt: time.Now(),
		},
		InitialQuantity: inventory.AvailableQuantity,
		WarehouseID:     inventory.WarehouseID,
	}
}

// NewInventoryUpdatedEvent 创建库存更新事件
func NewInventoryUpdatedEvent(inventory *model.Inventory) *InventoryUpdatedEvent {
	return &InventoryUpdatedEvent{
		InventoryEvent: InventoryEvent{
			ID:          inventory.ID,
			ProductID:   inventory.ProductID,
			SKU:         inventory.SKU,
			EventType:   "inventory.updated",
			OccurredAt:  inventory.UpdatedAt,
			PublishedAt: time.Now(),
		},
		AvailableQuantity: inventory.AvailableQuantity,
		ReservedQuantity:  inventory.ReservedQuantity,
		Status:            string(inventory.Status),
	}
}

// NewInventoryQuantityChangedEvent 创建库存数量变更事件
func NewInventoryQuantityChangedEvent(operation model.InventoryOperation) *InventoryQuantityChangedEvent {
	return &InventoryQuantityChangedEvent{
		InventoryEvent: InventoryEvent{
			ID:          operation.ID,
			ProductID:   operation.ProductID,
			SKU:         "", // 需要从库存信息中获取
			EventType:   "inventory.quantity_changed",
			OccurredAt:  operation.CreatedAt,
			PublishedAt: time.Now(),
		},
		OldQuantity:    operation.BeforeStock,
		NewQuantity:    operation.AfterStock,
		ChangeQuantity: operation.Quantity,
		OperationType:  string(operation.Type),
		Reason:         operation.Reason,
		RelatedOrderID: operation.RelatedOrderID,
	}
}

// NewInventoryLowStockEvent 创建库存不足事件
func NewInventoryLowStockEvent(inventory *model.Inventory) *InventoryLowStockEvent {
	return &InventoryLowStockEvent{
		InventoryEvent: InventoryEvent{
			ID:          inventory.ID,
			ProductID:   inventory.ProductID,
			SKU:         inventory.SKU,
			EventType:   "inventory.low_stock",
			OccurredAt:  inventory.UpdatedAt,
			PublishedAt: time.Now(),
		},
		CurrentQuantity:   inventory.AvailableQuantity,
		LowStockThreshold: inventory.LowStockThreshold,
	}
}

// NewInventoryOutOfStockEvent 创建库存不足事件
func NewInventoryOutOfStockEvent(inventory *model.Inventory, productName string) *InventoryOutOfStockEvent {
	return &InventoryOutOfStockEvent{
		InventoryEvent: InventoryEvent{
			ID:          inventory.ID,
			ProductID:   inventory.ProductID,
			SKU:         inventory.SKU,
			EventType:   "inventory.out_of_stock",
			OccurredAt:  inventory.UpdatedAt,
			PublishedAt: time.Now(),
		},
		ProductName: productName,
	}
}

// NewInventoryLockedEvent 创建库存锁定事件
func NewInventoryLockedEvent(inventory *model.Inventory, reason string) *InventoryLockedEvent {
	return &InventoryLockedEvent{
		InventoryEvent: InventoryEvent{
			ID:          inventory.ID,
			ProductID:   inventory.ProductID,
			SKU:         inventory.SKU,
			EventType:   "inventory.locked",
			OccurredAt:  inventory.UpdatedAt,
			PublishedAt: time.Now(),
		},
		Reason: reason,
	}
}

// NewInventoryUnlockedEvent 创建库存解锁事件
func NewInventoryUnlockedEvent(inventory *model.Inventory, reason string) *InventoryUnlockedEvent {
	return &InventoryUnlockedEvent{
		InventoryEvent: InventoryEvent{
			ID:          inventory.ID,
			ProductID:   inventory.ProductID,
			SKU:         inventory.SKU,
			EventType:   "inventory.unlocked",
			OccurredAt:  inventory.UpdatedAt,
			PublishedAt: time.Now(),
		},
		Reason: reason,
	}
}

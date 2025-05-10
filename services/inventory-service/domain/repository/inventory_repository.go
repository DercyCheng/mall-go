package repository

import (
	"context"

	"mall-go/services/inventory-service/domain/model"
)

// InventoryRepository defines the interface for inventory data access
type InventoryRepository interface {
	// CreateInventory creates a new inventory item
	CreateInventory(ctx context.Context, inventory *model.InventoryItem) error

	// GetInventory retrieves an inventory item by ID
	GetInventory(ctx context.Context, id string) (*model.InventoryItem, error)

	// GetInventoryByProduct retrieves an inventory item by product ID and SKU ID
	GetInventoryByProduct(ctx context.Context, productID, skuID string) (*model.InventoryItem, error)

	// GetInventoriesByProductID retrieves all inventory items for a product
	GetInventoriesByProductID(ctx context.Context, productID string) ([]*model.InventoryItem, error)

	// GetInventoriesByWarehouseID retrieves all inventory items in a warehouse
	GetInventoriesByWarehouseID(ctx context.Context, warehouseID string) ([]*model.InventoryItem, error)

	// UpdateInventory updates an existing inventory item
	UpdateInventory(ctx context.Context, inventory *model.InventoryItem) error

	// AddStock adds stock to an inventory item
	AddStock(ctx context.Context, id string, quantity int) error

	// RemoveStock removes stock from an inventory item
	RemoveStock(ctx context.Context, id string, quantity int) error

	// LockStock locks stock for a pending order
	LockStock(ctx context.Context, id string, quantity int) error

	// UnlockStock unlocks previously locked stock
	UnlockStock(ctx context.Context, id string, quantity int) error

	// ConfirmLock confirms a locked stock (e.g., when order is paid)
	ConfirmLock(ctx context.Context, id string, quantity int) error

	// SetStatus sets the inventory status
	SetStatus(ctx context.Context, id string, status model.InventoryStatus) error

	// DeleteInventory deletes an inventory item (soft delete)
	DeleteInventory(ctx context.Context, id string) error

	// ListInventories lists all inventory items with pagination
	ListInventories(ctx context.Context, page, pageSize int) ([]*model.InventoryItem, int, error)
}

// InventoryHistoryRepository defines the interface for inventory history data access
type InventoryHistoryRepository interface {
	// CreateHistory creates a new inventory history record
	CreateHistory(ctx context.Context, history *model.InventoryHistory) error

	// GetHistoryByInventoryID retrieves all history records for an inventory item
	GetHistoryByInventoryID(ctx context.Context, inventoryID string, page, pageSize int) ([]*model.InventoryHistory, int, error)

	// GetHistoryByProductID retrieves all history records for a product
	GetHistoryByProductID(ctx context.Context, productID string, page, pageSize int) ([]*model.InventoryHistory, int, error)

	// GetHistoryByWarehouseID retrieves all history records for a warehouse
	GetHistoryByWarehouseID(ctx context.Context, warehouseID string, page, pageSize int) ([]*model.InventoryHistory, int, error)

	// GetHistoryByOrderID retrieves all history records for an order
	GetHistoryByOrderID(ctx context.Context, orderID string) ([]*model.InventoryHistory, error)

	// ListHistory lists all inventory history records with pagination
	ListHistory(ctx context.Context, page, pageSize int) ([]*model.InventoryHistory, int, error)
}

// WarehouseRepository defines the interface for warehouse data access
type WarehouseRepository interface {
	// CreateWarehouse creates a new warehouse
	CreateWarehouse(ctx context.Context, warehouse *model.Warehouse) error

	// GetWarehouse retrieves a warehouse by ID
	GetWarehouse(ctx context.Context, id string) (*model.Warehouse, error)

	// GetWarehouseByCode retrieves a warehouse by code
	GetWarehouseByCode(ctx context.Context, code string) (*model.Warehouse, error)

	// GetDefaultWarehouse retrieves the default warehouse
	GetDefaultWarehouse(ctx context.Context) (*model.Warehouse, error)

	// UpdateWarehouse updates an existing warehouse
	UpdateWarehouse(ctx context.Context, warehouse *model.Warehouse) error

	// DeleteWarehouse deletes a warehouse (soft delete)
	DeleteWarehouse(ctx context.Context, id string) error

	// ListWarehouses lists all warehouses
	ListWarehouses(ctx context.Context) ([]*model.Warehouse, error)

	// SetDefaultWarehouse sets a warehouse as the default
	SetDefaultWarehouse(ctx context.Context, id string) error

	// EnableWarehouse enables a warehouse
	EnableWarehouse(ctx context.Context, id string) error

	// DisableWarehouse disables a warehouse
	DisableWarehouse(ctx context.Context, id string) error
}

// InventoryCache defines the interface for inventory data caching
type InventoryCache interface {
	// GetInventory gets an inventory item from cache
	GetInventory(ctx context.Context, id string) (*model.InventoryItem, error)

	// SetInventory sets an inventory item in cache
	SetInventory(ctx context.Context, inventory *model.InventoryItem) error

	// DeleteInventory deletes an inventory item from cache
	DeleteInventory(ctx context.Context, id string) error

	// GetWarehouse gets a warehouse from cache
	GetWarehouse(ctx context.Context, id string) (*model.Warehouse, error)

	// SetWarehouse sets a warehouse in cache
	SetWarehouse(ctx context.Context, warehouse *model.Warehouse) error

	// DeleteWarehouse deletes a warehouse from cache
	DeleteWarehouse(ctx context.Context, id string) error
}

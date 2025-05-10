package service

import (
	"context"
	
	"mall-go/services/inventory-service/application/dto"
)

// InventoryService defines the interface for inventory application service
type InventoryService interface {
	// CreateInventory creates a new inventory item
	CreateInventory(ctx context.Context, req dto.CreateInventoryRequest) (*dto.InventoryItemDTO, error)
	
	// GetInventory retrieves an inventory item by ID
	GetInventory(ctx context.Context, id string) (*dto.InventoryItemDTO, error)
	
	// GetInventoryByProduct retrieves an inventory item by product ID and SKU ID
	GetInventoryByProduct(ctx context.Context, productID, skuID string) (*dto.InventoryItemDTO, error)
	
	// GetInventoriesByProductID retrieves all inventory items for a product
	GetInventoriesByProductID(ctx context.Context, productID string) ([]dto.InventoryItemDTO, error)
	
	// GetInventoriesByWarehouseID retrieves all inventory items in a warehouse
	GetInventoriesByWarehouseID(ctx context.Context, warehouseID string, page, pageSize int) ([]dto.InventoryItemDTO, int, error)
	
	// ListInventories lists all inventory items with pagination
	ListInventories(ctx context.Context, page, pageSize int) ([]dto.InventoryItemDTO, int, error)
	
	// AddStock adds stock to an inventory item
	AddStock(ctx context.Context, req dto.AddStockRequest) error
	
	// RemoveStock removes stock from an inventory item
	RemoveStock(ctx context.Context, req dto.RemoveStockRequest) error
	
	// LockStock locks stock for a pending order
	LockStock(ctx context.Context, req dto.LockStockRequest) error
	
	// UnlockStock unlocks previously locked stock
	UnlockStock(ctx context.Context, req dto.UnlockStockRequest) error
	
	// ConfirmLock confirms a locked stock (e.g., when order is paid)
	ConfirmLock(ctx context.Context, req dto.ConfirmLockRequest) error
	
	// SetInventoryStatus sets the inventory status
	SetInventoryStatus(ctx context.Context, req dto.SetInventoryStatusRequest) error
	
	// GetInventoryHistory retrieves history records for an inventory item
	GetInventoryHistory(ctx context.Context, inventoryID string, page, pageSize int) ([]dto.InventoryHistoryDTO, int, error)
	
	// GetProductInventoryHistory retrieves history records for a product
	GetProductInventoryHistory(ctx context.Context, productID string, page, pageSize int) ([]dto.InventoryHistoryDTO, int, error)
	
	// GetWarehouseInventoryHistory retrieves history records for a warehouse
	GetWarehouseInventoryHistory(ctx context.Context, warehouseID string, page, pageSize int) ([]dto.InventoryHistoryDTO, int, error)
	
	// GetOrderInventoryHistory retrieves history records for an order
	GetOrderInventoryHistory(ctx context.Context, orderID string) ([]dto.InventoryHistoryDTO, error)
	
	// CheckStock checks if stock is available for a product
	CheckStock(ctx context.Context, req dto.CheckStockRequest) (*dto.CheckStockResponse, error)
}

// WarehouseService defines the interface for warehouse application service
type WarehouseService interface {
	// CreateWarehouse creates a new warehouse
	CreateWarehouse(ctx context.Context, req dto.CreateWarehouseRequest) (*dto.WarehouseDTO, error)
	
	// GetWarehouse retrieves a warehouse by ID
	GetWarehouse(ctx context.Context, id string) (*dto.WarehouseDTO, error)
	
	// GetWarehouseByCode retrieves a warehouse by code
	GetWarehouseByCode(ctx context.Context, code string) (*dto.WarehouseDTO, error)
	
	// GetDefaultWarehouse retrieves the default warehouse
	GetDefaultWarehouse(ctx context.Context) (*dto.WarehouseDTO, error)
	
	// UpdateWarehouse updates an existing warehouse
	UpdateWarehouse(ctx context.Context, req dto.UpdateWarehouseRequest) error
	
	// DeleteWarehouse deletes a warehouse
	DeleteWarehouse(ctx context.Context, id string) error
	
	// ListWarehouses lists all warehouses
	ListWarehouses(ctx context.Context) ([]dto.WarehouseDTO, error)
	
	// SetDefaultWarehouse sets a warehouse as the default
	SetDefaultWarehouse(ctx context.Context, req dto.SetDefaultWarehouseRequest) error
	
	// EnableWarehouse enables a warehouse
	EnableWarehouse(ctx context.Context, id string) error
	
	// DisableWarehouse disables a warehouse
	DisableWarehouse(ctx context.Context, id string) error
}

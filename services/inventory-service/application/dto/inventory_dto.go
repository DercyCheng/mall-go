package dto

// InventoryItemDTO represents an inventory item data transfer object
type InventoryItemDTO struct {
	ID           string `json:"id"`
	ProductID    string `json:"productId"`
	SkuID        string `json:"skuId"`
	SkuCode      string `json:"skuCode,omitempty"`
	WarehouseID  string `json:"warehouseId"`
	Quantity     int    `json:"quantity"`
	LockedCount  int    `json:"lockedCount"`
	AvailableQty int    `json:"availableQty"`
	Status       string `json:"status"` // "normal", "locked", "sold", "defective"
}

// WarehouseDTO represents a warehouse data transfer object
type WarehouseDTO struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Code         string `json:"code"`
	Address      string `json:"address,omitempty"`
	ContactName  string `json:"contactName,omitempty"`
	ContactPhone string `json:"contactPhone,omitempty"`
	Status       int    `json:"status"` // 0: disabled, 1: enabled
	IsDefault    bool   `json:"isDefault"`
}

// InventoryHistoryDTO represents an inventory history data transfer object
type InventoryHistoryDTO struct {
	ID            string `json:"id"`
	InventoryID   string `json:"inventoryId"`
	ProductID     string `json:"productId"`
	SkuID         string `json:"skuId"`
	WarehouseID   string `json:"warehouseId"`
	OperationType string `json:"operationType"` // "add", "remove", "lock", "unlock", "confirm"
	Quantity      int    `json:"quantity"`
	BeforeQty     int    `json:"beforeQty"`
	AfterQty      int    `json:"afterQty"`
	Operator      string `json:"operator"`
	OrderID       string `json:"orderId,omitempty"`
	Reason        string `json:"reason,omitempty"`
	CreatedAt     string `json:"createdAt"`
}

// CreateInventoryRequest represents a request to create an inventory item
type CreateInventoryRequest struct {
	ProductID   string `json:"productId" binding:"required"`
	SkuID       string `json:"skuId" binding:"required"`
	SkuCode     string `json:"skuCode"`
	WarehouseID string `json:"warehouseId" binding:"required"`
	Quantity    int    `json:"quantity" binding:"min=0"`
}

// AddStockRequest represents a request to add stock to an inventory item
type AddStockRequest struct {
	ID         string `json:"id" binding:"required"`
	Quantity   int    `json:"quantity" binding:"required,min=1"`
	Reason     string `json:"reason"`
	OperatorID string `json:"operatorId"`
}

// RemoveStockRequest represents a request to remove stock from an inventory item
type RemoveStockRequest struct {
	ID         string `json:"id" binding:"required"`
	Quantity   int    `json:"quantity" binding:"required,min=1"`
	Reason     string `json:"reason"`
	OperatorID string `json:"operatorId"`
}

// LockStockRequest represents a request to lock stock for an order
type LockStockRequest struct {
	ID         string `json:"id" binding:"required"`
	Quantity   int    `json:"quantity" binding:"required,min=1"`
	OrderID    string `json:"orderId" binding:"required"`
	OperatorID string `json:"operatorId"`
}

// UnlockStockRequest represents a request to unlock previously locked stock
type UnlockStockRequest struct {
	ID         string `json:"id" binding:"required"`
	Quantity   int    `json:"quantity" binding:"required,min=1"`
	OrderID    string `json:"orderId" binding:"required"`
	Reason     string `json:"reason"`
	OperatorID string `json:"operatorId"`
}

// ConfirmLockRequest represents a request to confirm a locked stock
type ConfirmLockRequest struct {
	ID         string `json:"id" binding:"required"`
	Quantity   int    `json:"quantity" binding:"required,min=1"`
	OrderID    string `json:"orderId" binding:"required"`
	OperatorID string `json:"operatorId"`
}

// SetInventoryStatusRequest represents a request to set the inventory status
type SetInventoryStatusRequest struct {
	ID         string `json:"id" binding:"required"`
	Status     string `json:"status" binding:"required"`
	Reason     string `json:"reason"`
	OperatorID string `json:"operatorId"`
}

// CreateWarehouseRequest represents a request to create a warehouse
type CreateWarehouseRequest struct {
	Name         string `json:"name" binding:"required"`
	Code         string `json:"code" binding:"required"`
	Address      string `json:"address"`
	ContactName  string `json:"contactName"`
	ContactPhone string `json:"contactPhone"`
	IsDefault    bool   `json:"isDefault"`
}

// UpdateWarehouseRequest represents a request to update a warehouse
type UpdateWarehouseRequest struct {
	ID           string `json:"id" binding:"required"`
	Name         string `json:"name"`
	Code         string `json:"code"`
	Address      string `json:"address"`
	ContactName  string `json:"contactName"`
	ContactPhone string `json:"contactPhone"`
	Status       int    `json:"status"`
	IsDefault    bool   `json:"isDefault"`
}

// SetDefaultWarehouseRequest represents a request to set a warehouse as the default
type SetDefaultWarehouseRequest struct {
	ID string `json:"id" binding:"required"`
}

// GetInventoryResponse represents a response containing an inventory item
type GetInventoryResponse struct {
	Success   bool            `json:"success"`
	Message   string          `json:"message"`
	Inventory InventoryItemDTO `json:"inventory"`
}

// ListInventoriesResponse represents a response containing a list of inventory items
type ListInventoriesResponse struct {
	Success  bool              `json:"success"`
	Message  string            `json:"message"`
	Items    []InventoryItemDTO `json:"items"`
	Total    int               `json:"total"`
	Page     int               `json:"page"`
	PageSize int               `json:"pageSize"`
}

// GetWarehouseResponse represents a response containing a warehouse
type GetWarehouseResponse struct {
	Success   bool         `json:"success"`
	Message   string       `json:"message"`
	Warehouse WarehouseDTO `json:"warehouse"`
}

// ListWarehousesResponse represents a response containing a list of warehouses
type ListWarehousesResponse struct {
	Success    bool           `json:"success"`
	Message    string         `json:"message"`
	Warehouses []WarehouseDTO `json:"warehouses"`
}

// ListInventoryHistoryResponse represents a response containing a list of inventory history records
type ListInventoryHistoryResponse struct {
	Success  bool                 `json:"success"`
	Message  string               `json:"message"`
	Records  []InventoryHistoryDTO `json:"records"`
	Total    int                  `json:"total"`
	Page     int                  `json:"page"`
	PageSize int                  `json:"pageSize"`
}

// GenericResponse represents a generic response
type GenericResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// CheckStockRequest represents a request to check stock availability
type CheckStockRequest struct {
	ProductID string `json:"productId" binding:"required"`
	SkuID     string `json:"skuId" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
}

// CheckStockResponse represents a response containing stock availability information
type CheckStockResponse struct {
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	Available    bool   `json:"available"`
	AvailableQty int    `json:"availableQty"`
}

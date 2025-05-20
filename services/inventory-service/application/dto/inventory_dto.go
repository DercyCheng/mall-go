package dto

import (
	"time"
)

// InventoryCreateRequest 创建库存请求DTO
type InventoryCreateRequest struct {
	ProductID         string `json:"productId" binding:"required"`
	SKU               string `json:"sku" binding:"required"`
	InitialQuantity   int    `json:"initialQuantity" binding:"min=0"`
	LowStockThreshold int    `json:"lowStockThreshold" binding:"min=0"`
	WarehouseID       string `json:"warehouseId" binding:"required"`
	ShelfLocation     string `json:"shelfLocation,omitempty"`
}

// InventoryUpdateRequest 更新库存请求DTO
type InventoryUpdateRequest struct {
	ID                string `json:"id" binding:"required"`
	LowStockThreshold int    `json:"lowStockThreshold,omitempty" binding:"min=0"`
	ShelfLocation     string `json:"shelfLocation,omitempty"`
}

// InventoryAdjustRequest 调整库存请求DTO
type InventoryAdjustRequest struct {
	ProductID   string `json:"productId" binding:"required"`
	NewQuantity int    `json:"newQuantity" binding:"min=0"`
	Reason      string `json:"reason" binding:"required"`
	OperatorID  string `json:"operatorId,omitempty"`
}

// InventoryInboundRequest 入库请求DTO
type InventoryInboundRequest struct {
	ProductID  string `json:"productId" binding:"required"`
	Quantity   int    `json:"quantity" binding:"required,min=1"`
	Reason     string `json:"reason" binding:"required"`
	OperatorID string `json:"operatorId,omitempty"`
}

// InventoryOutboundRequest 出库请求DTO
type InventoryOutboundRequest struct {
	ProductID  string `json:"productId" binding:"required"`
	Quantity   int    `json:"quantity" binding:"required,min=1"`
	OrderID    string `json:"orderId" binding:"required"`
	Reason     string `json:"reason" binding:"required"`
	OperatorID string `json:"operatorId,omitempty"`
}

// InventoryReserveRequest 预留库存请求DTO
type InventoryReserveRequest struct {
	ProductID  string `json:"productId" binding:"required"`
	Quantity   int    `json:"quantity" binding:"required,min=1"`
	OrderID    string `json:"orderId" binding:"required"`
	Reason     string `json:"reason,omitempty"`
	OperatorID string `json:"operatorId,omitempty"`
}

// InventoryReleaseRequest 释放预留库存请求DTO
type InventoryReleaseRequest struct {
	ProductID  string `json:"productId" binding:"required"`
	Quantity   int    `json:"quantity" binding:"required,min=1"`
	OrderID    string `json:"orderId" binding:"required"`
	Reason     string `json:"reason,omitempty"`
	OperatorID string `json:"operatorId,omitempty"`
}

// InventoryConfirmRequest 确认预留库存请求DTO
type InventoryConfirmRequest struct {
	ProductID string `json:"productId" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
	OrderID   string `json:"orderId" binding:"required"`
	Reason    string `json:"reason,omitempty"`
}

// InventoryLockRequest 锁定库存请求DTO
type InventoryLockRequest struct {
	ProductID  string `json:"productId" binding:"required"`
	Reason     string `json:"reason" binding:"required"`
	OperatorID string `json:"operatorId,omitempty"`
}

// InventoryUnlockRequest 解锁库存请求DTO
type InventoryUnlockRequest struct {
	ProductID  string `json:"productId" binding:"required"`
	Reason     string `json:"reason,omitempty"`
	OperatorID string `json:"operatorId,omitempty"`
}

// InventoryQueryRequest 库存查询请求DTO
type InventoryQueryRequest struct {
	ProductID   string `form:"productId,omitempty"`
	SKU         string `form:"sku,omitempty"`
	WarehouseID string `form:"warehouseId,omitempty"`
	Status      string `form:"status,omitempty"`
	LowStock    bool   `form:"lowStock,omitempty"`
	Keyword     string `form:"keyword,omitempty"`
	Page        int    `form:"page,default=1" binding:"min=1"`
	Size        int    `form:"size,default=20" binding:"min=1,max=100"`
}

// InventoryOperationResponse 库存操作记录响应DTO
type InventoryOperationResponse struct {
	ID             string    `json:"id"`
	ProductID      string    `json:"productId"`
	Type           string    `json:"type"`
	Quantity       int       `json:"quantity"`
	BeforeStock    int       `json:"beforeStock"`
	AfterStock     int       `json:"afterStock"`
	Reason         string    `json:"reason"`
	RelatedOrderID string    `json:"relatedOrderId,omitempty"`
	OperatorID     string    `json:"operatorId,omitempty"`
	CreatedAt      time.Time `json:"createdAt"`
}

// InventoryResponse 库存响应DTO
type InventoryResponse struct {
	ID                 string                       `json:"id"`
	ProductID          string                       `json:"productId"`
	SKU                string                       `json:"sku"`
	AvailableQuantity  int                          `json:"availableQuantity"`
	ReservedQuantity   int                          `json:"reservedQuantity"`
	TotalQuantity      int                          `json:"totalQuantity"`
	LowStockThreshold  int                          `json:"lowStockThreshold"`
	LockStatus         bool                         `json:"lockStatus"`
	WarehouseID        string                       `json:"warehouseId"`
	ShelfLocation      string                       `json:"shelfLocation,omitempty"`
	Status             string                       `json:"status"`
	LastStockCheckDate time.Time                    `json:"lastStockCheckDate"`
	Operations         []InventoryOperationResponse `json:"operations,omitempty"`
	CreatedAt          time.Time                    `json:"createdAt"`
	UpdatedAt          time.Time                    `json:"updatedAt"`
}

// InventoryBriefResponse 库存简要响应DTO
type InventoryBriefResponse struct {
	ID                string    `json:"id"`
	ProductID         string    `json:"productId"`
	SKU               string    `json:"sku"`
	AvailableQuantity int       `json:"availableQuantity"`
	ReservedQuantity  int       `json:"reservedQuantity"`
	TotalQuantity     int       `json:"totalQuantity"`
	Status            string    `json:"status"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

// InventoryListResponse 库存列表响应DTO
type InventoryListResponse struct {
	Items []InventoryBriefResponse `json:"items"`
	Total int64                    `json:"total"`
	Page  int                      `json:"page"`
	Size  int                      `json:"size"`
}

// InventoryOperationListResponse 库存操作记录列表响应DTO
type InventoryOperationListResponse struct {
	Items []InventoryOperationResponse `json:"items"`
	Total int64                        `json:"total"`
	Page  int                          `json:"page"`
	Size  int                          `json:"size"`
}

// InventoryStockCheckRequest 库存盘点请求DTO
type InventoryStockCheckRequest struct {
	ProductID      string `json:"productId" binding:"required"`
	ActualQuantity int    `json:"actualQuantity" binding:"min=0"`
	Reason         string `json:"reason" binding:"required"`
	OperatorID     string `json:"operatorId,omitempty"`
}

// InventoryStatisticsResponse 库存统计响应DTO
type InventoryStatisticsResponse struct {
	TotalProducts      int64               `json:"totalProducts"`
	TotalQuantity      int64               `json:"totalQuantity"`
	LowStockProducts   int64               `json:"lowStockProducts"`
	OutOfStockProducts int64               `json:"outOfStockProducts"`
	LockedProducts     int64               `json:"lockedProducts"`
	WarehouseStats     []WarehouseStatItem `json:"warehouseStats,omitempty"`
}

// WarehouseStatItem 仓库统计项
type WarehouseStatItem struct {
	WarehouseID        string `json:"warehouseId"`
	WarehouseName      string `json:"warehouseName"`
	ProductCount       int64  `json:"productCount"`
	TotalQuantity      int64  `json:"totalQuantity"`
	LowStockProducts   int64  `json:"lowStockProducts"`
	OutOfStockProducts int64  `json:"outOfStockProducts"`
}

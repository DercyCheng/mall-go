package handler

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"mall-go/services/inventory-service/application/dto"
	"mall-go/services/inventory-service/application/service"
	"mall-go/services/inventory-service/proto/inventorypb"
)

// InventoryHandler handles inventory-related gRPC requests
type InventoryHandler struct {
	inventoryService service.InventoryService
	inventorypb.UnimplementedInventoryServiceServer
}

// NewInventoryHandler creates a new inventory handler
func NewInventoryHandler(inventoryService service.InventoryService) *InventoryHandler {
	return &InventoryHandler{
		inventoryService: inventoryService,
	}
}

// CreateInventory creates a new inventory item
func (h *InventoryHandler) CreateInventory(ctx context.Context, req *inventorypb.CreateInventoryRequest) (*inventorypb.InventoryResponse, error) {
	// Validate request
	if req.ProductId == "" || req.SkuId == "" || req.WarehouseId == "" {
		return nil, status.Error(codes.InvalidArgument, "product ID, SKU ID and warehouse ID are required")
	}

	// Map request to DTO
	createReq := dto.CreateInventoryRequest{
		ProductID:   req.ProductId,
		SkuID:       req.SkuId,
		SkuCode:     req.SkuCode,
		WarehouseID: req.WarehouseId,
		Quantity:    int(req.Quantity),
	}

	// Call service
	inventory, err := h.inventoryService.CreateInventory(ctx, createReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create inventory: %v", err)
	}

	// Map response
	return &inventorypb.InventoryResponse{
		Success:   true,
		Message:   "Inventory item created successfully",
		Inventory: mapToInventoryItemProto(inventory),
	}, nil
}

// GetInventory retrieves an inventory item by ID
func (h *InventoryHandler) GetInventory(ctx context.Context, req *inventorypb.GetInventoryRequest) (*inventorypb.InventoryResponse, error) {
	// Validate request
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "inventory ID is required")
	}

	// Call service
	inventory, err := h.inventoryService.GetInventory(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get inventory: %v", err)
	}

	// Map response
	return &inventorypb.InventoryResponse{
		Success:   true,
		Message:   "Inventory item retrieved successfully",
		Inventory: mapToInventoryItemProto(inventory),
	}, nil
}

// GetInventoryByProduct retrieves an inventory item by product ID and SKU ID
func (h *InventoryHandler) GetInventoryByProduct(ctx context.Context, req *inventorypb.GetInventoryByProductRequest) (*inventorypb.InventoryResponse, error) {
	// Validate request
	if req.ProductId == "" || req.SkuId == "" {
		return nil, status.Error(codes.InvalidArgument, "product ID and SKU ID are required")
	}

	// Call service
	inventory, err := h.inventoryService.GetInventoryByProduct(ctx, req.ProductId, req.SkuId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get inventory by product: %v", err)
	}

	// Map response
	return &inventorypb.InventoryResponse{
		Success:   true,
		Message:   "Inventory item retrieved successfully",
		Inventory: mapToInventoryItemProto(inventory),
	}, nil
}

// GetInventoriesByProductID retrieves all inventory items for a product
func (h *InventoryHandler) GetInventoriesByProductID(ctx context.Context, req *inventorypb.GetInventoriesByProductIDRequest) (*inventorypb.InventoriesResponse, error) {
	// Validate request
	if req.ProductId == "" {
		return nil, status.Error(codes.InvalidArgument, "product ID is required")
	}

	// Call service
	inventories, err := h.inventoryService.GetInventoriesByProductID(ctx, req.ProductId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get inventories by product ID: %v", err)
	}

	// Map response
	items := make([]*inventorypb.InventoryItem, len(inventories))
	for i, inventory := range inventories {
		items[i] = mapToInventoryItemProto(&inventory)
	}

	return &inventorypb.InventoriesResponse{
		Success: true,
		Message: "Inventory items retrieved successfully",
		Items:   items,
	}, nil
}

// GetInventoriesByWarehouseID retrieves all inventory items in a warehouse
func (h *InventoryHandler) GetInventoriesByWarehouseID(ctx context.Context, req *inventorypb.GetInventoriesByWarehouseIDRequest) (*inventorypb.PagedInventoriesResponse, error) {
	// Validate request
	if req.WarehouseId == "" {
		return nil, status.Error(codes.InvalidArgument, "warehouse ID is required")
	}

	page := int(req.Page)
	if page <= 0 {
		page = 1
	}

	pageSize := int(req.PageSize)
	if pageSize <= 0 {
		pageSize = 10
	}

	// Call service
	inventories, total, err := h.inventoryService.GetInventoriesByWarehouseID(ctx, req.WarehouseId, page, pageSize)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get inventories by warehouse ID: %v", err)
	}

	// Map response
	items := make([]*inventorypb.InventoryItem, len(inventories))
	for i, inventory := range inventories {
		items[i] = mapToInventoryItemProto(&inventory)
	}

	return &inventorypb.PagedInventoriesResponse{
		Success:  true,
		Message:  "Inventory items retrieved successfully",
		Items:    items,
		Total:    int32(total),
		Page:     int32(page),
		PageSize: int32(pageSize),
	}, nil
}

// ListInventories lists all inventory items with pagination
func (h *InventoryHandler) ListInventories(ctx context.Context, req *inventorypb.ListInventoriesRequest) (*inventorypb.PagedInventoriesResponse, error) {
	page := int(req.Page)
	if page <= 0 {
		page = 1
	}

	pageSize := int(req.PageSize)
	if pageSize <= 0 {
		pageSize = 10
	}

	// Call service
	inventories, total, err := h.inventoryService.ListInventories(ctx, page, pageSize)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list inventories: %v", err)
	}

	// Map response
	items := make([]*inventorypb.InventoryItem, len(inventories))
	for i, inventory := range inventories {
		items[i] = mapToInventoryItemProto(&inventory)
	}

	return &inventorypb.PagedInventoriesResponse{
		Success:  true,
		Message:  "Inventory items retrieved successfully",
		Items:    items,
		Total:    int32(total),
		Page:     int32(page),
		PageSize: int32(pageSize),
	}, nil
}

// AddStock adds stock to an inventory item
func (h *InventoryHandler) AddStock(ctx context.Context, req *inventorypb.AddStockRequest) (*inventorypb.GenericResponse, error) {
	// Validate request
	if req.Id == "" || req.Quantity <= 0 {
		return nil, status.Error(codes.InvalidArgument, "inventory ID and positive quantity are required")
	}

	// Map request to DTO
	addReq := dto.AddStockRequest{
		ID:         req.Id,
		Quantity:   int(req.Quantity),
		Reason:     req.Reason,
		OperatorID: req.OperatorId,
	}

	// Call service
	if err := h.inventoryService.AddStock(ctx, addReq); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add stock: %v", err)
	}

	return &inventorypb.GenericResponse{
		Success: true,
		Message: "Stock added successfully",
	}, nil
}

// RemoveStock removes stock from an inventory item
func (h *InventoryHandler) RemoveStock(ctx context.Context, req *inventorypb.RemoveStockRequest) (*inventorypb.GenericResponse, error) {
	// Validate request
	if req.Id == "" || req.Quantity <= 0 {
		return nil, status.Error(codes.InvalidArgument, "inventory ID and positive quantity are required")
	}

	// Map request to DTO
	removeReq := dto.RemoveStockRequest{
		ID:         req.Id,
		Quantity:   int(req.Quantity),
		Reason:     req.Reason,
		OperatorID: req.OperatorId,
	}

	// Call service
	if err := h.inventoryService.RemoveStock(ctx, removeReq); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to remove stock: %v", err)
	}

	return &inventorypb.GenericResponse{
		Success: true,
		Message: "Stock removed successfully",
	}, nil
}

// LockStock locks stock for a pending order
func (h *InventoryHandler) LockStock(ctx context.Context, req *inventorypb.LockStockRequest) (*inventorypb.GenericResponse, error) {
	// Validate request
	if req.Id == "" || req.Quantity <= 0 || req.OrderId == "" {
		return nil, status.Error(codes.InvalidArgument, "inventory ID, positive quantity and order ID are required")
	}

	// Map request to DTO
	lockReq := dto.LockStockRequest{
		ID:         req.Id,
		Quantity:   int(req.Quantity),
		OrderID:    req.OrderId,
		OperatorID: req.OperatorId,
	}

	// Call service
	if err := h.inventoryService.LockStock(ctx, lockReq); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to lock stock: %v", err)
	}

	return &inventorypb.GenericResponse{
		Success: true,
		Message: "Stock locked successfully",
	}, nil
}

// UnlockStock unlocks previously locked stock
func (h *InventoryHandler) UnlockStock(ctx context.Context, req *inventorypb.UnlockStockRequest) (*inventorypb.GenericResponse, error) {
	// Validate request
	if req.Id == "" || req.Quantity <= 0 || req.OrderId == "" {
		return nil, status.Error(codes.InvalidArgument, "inventory ID, positive quantity and order ID are required")
	}

	// Map request to DTO
	unlockReq := dto.UnlockStockRequest{
		ID:         req.Id,
		Quantity:   int(req.Quantity),
		OrderID:    req.OrderId,
		Reason:     req.Reason,
		OperatorID: req.OperatorId,
	}

	// Call service
	if err := h.inventoryService.UnlockStock(ctx, unlockReq); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to unlock stock: %v", err)
	}

	return &inventorypb.GenericResponse{
		Success: true,
		Message: "Stock unlocked successfully",
	}, nil
}

// ConfirmLock confirms a locked stock (e.g., when order is paid)
func (h *InventoryHandler) ConfirmLock(ctx context.Context, req *inventorypb.ConfirmLockRequest) (*inventorypb.GenericResponse, error) {
	// Validate request
	if req.Id == "" || req.Quantity <= 0 || req.OrderId == "" {
		return nil, status.Error(codes.InvalidArgument, "inventory ID, positive quantity and order ID are required")
	}

	// Map request to DTO
	confirmReq := dto.ConfirmLockRequest{
		ID:         req.Id,
		Quantity:   int(req.Quantity),
		OrderID:    req.OrderId,
		OperatorID: req.OperatorId,
	}

	// Call service
	if err := h.inventoryService.ConfirmLock(ctx, confirmReq); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to confirm lock: %v", err)
	}

	return &inventorypb.GenericResponse{
		Success: true,
		Message: "Lock confirmed successfully",
	}, nil
}

// SetInventoryStatus sets the inventory status
func (h *InventoryHandler) SetInventoryStatus(ctx context.Context, req *inventorypb.SetInventoryStatusRequest) (*inventorypb.GenericResponse, error) {
	// Validate request
	if req.Id == "" || req.Status == "" {
		return nil, status.Error(codes.InvalidArgument, "inventory ID and status are required")
	}

	// Map request to DTO
	setReq := dto.SetInventoryStatusRequest{
		ID:         req.Id,
		Status:     req.Status,
		Reason:     req.Reason,
		OperatorID: req.OperatorId,
	}

	// Call service
	if err := h.inventoryService.SetInventoryStatus(ctx, setReq); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to set inventory status: %v", err)
	}

	return &inventorypb.GenericResponse{
		Success: true,
		Message: "Inventory status set successfully",
	}, nil
}

// GetInventoryHistory retrieves history records for an inventory item
func (h *InventoryHandler) GetInventoryHistory(ctx context.Context, req *inventorypb.GetInventoryHistoryRequest) (*inventorypb.InventoryHistoryResponse, error) {
	// Validate request
	if req.InventoryId == "" {
		return nil, status.Error(codes.InvalidArgument, "inventory ID is required")
	}

	page := int(req.Page)
	if page <= 0 {
		page = 1
	}

	pageSize := int(req.PageSize)
	if pageSize <= 0 {
		pageSize = 10
	}

	// Call service
	records, total, err := h.inventoryService.GetInventoryHistory(ctx, req.InventoryId, page, pageSize)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get inventory history: %v", err)
	}

	// Map response
	historyRecords := make([]*inventorypb.InventoryHistory, len(records))
	for i, record := range records {
		historyRecords[i] = mapToInventoryHistoryProto(&record)
	}

	return &inventorypb.InventoryHistoryResponse{
		Success:  true,
		Message:  "Inventory history retrieved successfully",
		Records:  historyRecords,
		Total:    int32(total),
		Page:     int32(page),
		PageSize: int32(pageSize),
	}, nil
}

// CheckStock checks if stock is available for a product
func (h *InventoryHandler) CheckStock(ctx context.Context, req *inventorypb.CheckStockRequest) (*inventorypb.CheckStockResponse, error) {
	// Validate request
	if req.ProductId == "" || req.SkuId == "" || req.Quantity <= 0 {
		return nil, status.Error(codes.InvalidArgument, "product ID, SKU ID and positive quantity are required")
	}

	// Map request to DTO
	checkReq := dto.CheckStockRequest{
		ProductID: req.ProductId,
		SkuID:     req.SkuId,
		Quantity:  int(req.Quantity),
	}

	// Call service
	result, err := h.inventoryService.CheckStock(ctx, checkReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check stock: %v", err)
	}

	return &inventorypb.CheckStockResponse{
		Success:      result.Success,
		Message:      result.Message,
		Available:    result.Available,
		AvailableQty: int32(result.AvailableQty),
	}, nil
}

// mapToInventoryItemProto maps a dto.InventoryItemDTO to inventorypb.InventoryItem
func mapToInventoryItemProto(item *dto.InventoryItemDTO) *inventorypb.InventoryItem {
	return &inventorypb.InventoryItem{
		Id:          item.ID,
		ProductId:   item.ProductID,
		SkuId:       item.SkuID,
		SkuCode:     item.SkuCode,
		WarehouseId: item.WarehouseID,
		Quantity:    int32(item.Quantity),
		LockedCount: int32(item.LockedCount),
		AvailableQty: int32(item.AvailableQty),
		Status:      item.Status,
	}
}

// mapToInventoryHistoryProto maps a dto.InventoryHistoryDTO to inventorypb.InventoryHistory
func mapToInventoryHistoryProto(history *dto.InventoryHistoryDTO) *inventorypb.InventoryHistory {
	return &inventorypb.InventoryHistory{
		Id:            history.ID,
		InventoryId:   history.InventoryID,
		ProductId:     history.ProductID,
		SkuId:         history.SkuID,
		WarehouseId:   history.WarehouseID,
		OperationType: history.OperationType,
		Quantity:      int32(history.Quantity),
		BeforeQty:     int32(history.BeforeQty),
		AfterQty:      int32(history.AfterQty),
		Operator:      history.Operator,
		OrderId:       history.OrderID,
		Reason:        history.Reason,
		CreatedAt:     history.CreatedAt,
	}
}

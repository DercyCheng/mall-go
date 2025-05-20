package service

import (
	"context"
	"errors"
	"time"

	"mall-go/services/inventory-service/application/dto"
	"mall-go/services/inventory-service/domain/model"
	"mall-go/services/inventory-service/domain/repository"
	domainService "mall-go/services/inventory-service/domain/service"
)

// InventoryService 库存应用服务接口
type InventoryService interface {
	// CreateInventory 创建库存
	CreateInventory(ctx context.Context, req dto.InventoryCreateRequest) (dto.InventoryResponse, error)

	// UpdateInventory 更新库存信息
	UpdateInventory(ctx context.Context, req dto.InventoryUpdateRequest) error

	// AdjustInventory 手动调整库存数量
	AdjustInventory(ctx context.Context, req dto.InventoryAdjustRequest) error

	// InboundInventory 入库操作
	InboundInventory(ctx context.Context, req dto.InventoryInboundRequest) error

	// OutboundInventory 出库操作
	OutboundInventory(ctx context.Context, req dto.InventoryOutboundRequest) error

	// ReserveInventory 预留库存
	ReserveInventory(ctx context.Context, req dto.InventoryReserveRequest) error

	// ReleaseInventory 释放预留库存
	ReleaseInventory(ctx context.Context, req dto.InventoryReleaseRequest) error

	// LockInventory 锁定库存
	LockInventory(ctx context.Context, req dto.InventoryLockRequest) error

	// UnlockInventory 解锁库存
	UnlockInventory(ctx context.Context, req dto.InventoryUnlockRequest) error

	// GetInventory 获取库存详情
	GetInventory(ctx context.Context, id string) (dto.InventoryResponse, error)

	// GetInventoryByProductID 根据商品ID获取库存
	GetInventoryByProductID(ctx context.Context, productID string) (dto.InventoryResponse, error)

	// GetInventoryBySKU 根据SKU获取库存
	GetInventoryBySKU(ctx context.Context, sku string) (dto.InventoryResponse, error)

	// GetLowStockInventories 获取低库存商品列表
	GetLowStockInventories(ctx context.Context, page, size int) (dto.InventoryListResponse, error)

	// SearchInventories 搜索库存
	SearchInventories(ctx context.Context, req dto.InventoryQueryRequest) (dto.InventoryListResponse, error)

	// GetInventoryOperations 获取库存操作记录
	GetInventoryOperations(ctx context.Context, productID string, page, size int) (dto.InventoryOperationListResponse, error)

	// GetInventoryOperationsByOrderID 根据订单ID获取库存操作记录
	GetInventoryOperationsByOrderID(ctx context.Context, orderID string) (dto.InventoryOperationListResponse, error)
}

// inventoryServiceImpl 库存应用服务实现
type inventoryServiceImpl struct {
	inventoryDomainService *domainService.InventoryDomainService
	inventoryRepo          repository.InventoryRepository
}

// NewInventoryService 创建库存应用服务
func NewInventoryService(
	inventoryDomainService *domainService.InventoryDomainService,
	inventoryRepo repository.InventoryRepository,
) InventoryService {
	return &inventoryServiceImpl{
		inventoryDomainService: inventoryDomainService,
		inventoryRepo:          inventoryRepo,
	}
}

// CreateInventory 创建库存
func (s *inventoryServiceImpl) CreateInventory(ctx context.Context, req dto.InventoryCreateRequest) (dto.InventoryResponse, error) {
	// 验证请求
	if req.ProductID == "" || req.SKU == "" || req.WarehouseID == "" {
		return dto.InventoryResponse{}, errors.New("missing required fields")
	}

	// 调用领域服务创建库存
	inventory, err := s.inventoryDomainService.CreateInventory(
		ctx,
		req.ProductID,
		req.SKU,
		req.InitialQuantity,
		req.LowStockThreshold,
		req.WarehouseID,
		req.ShelfLocation,
	)

	if err != nil {
		return dto.InventoryResponse{}, err
	}

	// 转换为DTO响应
	response := dto.InventoryResponse{
		ID:                 inventory.ID,
		ProductID:          inventory.ProductID,
		SKU:                inventory.SKU,
		AvailableQuantity:  inventory.AvailableQuantity,
		ReservedQuantity:   inventory.ReservedQuantity,
		LowStockThreshold:  inventory.LowStockThreshold,
		Status:             string(inventory.Status),
		LockStatus:         inventory.LockStatus,
		WarehouseID:        inventory.WarehouseID,
		ShelfLocation:      inventory.ShelfLocation,
		LastStockCheckDate: inventory.LastStockCheckDate,
		CreatedAt:          inventory.CreatedAt,
		UpdatedAt:          inventory.UpdatedAt,
	}

	return response, nil
}

// UpdateInventory 更新库存信息
func (s *inventoryServiceImpl) UpdateInventory(ctx context.Context, req dto.InventoryUpdateRequest) error {
	// 验证请求
	if req.ID == "" {
		return errors.New("inventory ID is required")
	}

	// 查找库存
	inventory, err := s.inventoryRepo.FindByID(ctx, req.ID)
	if err != nil {
		return err
	}

	// 更新库存信息
	if req.LowStockThreshold > 0 {
		inventory.LowStockThreshold = req.LowStockThreshold
	}
	if req.ShelfLocation != "" {
		inventory.ShelfLocation = req.ShelfLocation
	}
	inventory.UpdatedAt = time.Now()

	// 保存更新
	return s.inventoryRepo.Update(ctx, inventory)
}

// AdjustInventory 手动调整库存数量
func (s *inventoryServiceImpl) AdjustInventory(ctx context.Context, req dto.InventoryAdjustRequest) error {
	// 验证请求
	if req.ProductID == "" {
		return errors.New("product ID is required")
	}
	if req.Reason == "" {
		return errors.New("adjustment reason is required")
	}

	// 调用领域服务调整库存
	return s.inventoryDomainService.AdjustStock(ctx, req.ProductID, req.NewQuantity, req.Reason, req.OperatorID)
}

// InboundInventory 入库操作
func (s *inventoryServiceImpl) InboundInventory(ctx context.Context, req dto.InventoryInboundRequest) error {
	// 验证请求
	if req.ProductID == "" || req.Quantity <= 0 {
		return errors.New("invalid inbound request")
	}
	if req.Reason == "" {
		return errors.New("inbound reason is required")
	}

	// 调用领域服务入库
	return s.inventoryDomainService.InboundStock(ctx, req.ProductID, req.Quantity, req.Reason, req.OperatorID)
}

// OutboundInventory 出库操作
func (s *inventoryServiceImpl) OutboundInventory(ctx context.Context, req dto.InventoryOutboundRequest) error {
	// 验证请求
	if req.ProductID == "" || req.Quantity <= 0 || req.OrderID == "" {
		return errors.New("invalid outbound request")
	}
	if req.Reason == "" {
		return errors.New("outbound reason is required")
	}

	// 调用领域服务出库
	return s.inventoryDomainService.OutboundStock(ctx, req.ProductID, req.Quantity, req.OrderID, req.Reason, req.OperatorID)
}

// ReserveInventory 预留库存
func (s *inventoryServiceImpl) ReserveInventory(ctx context.Context, req dto.InventoryReserveRequest) error {
	// 验证请求
	if req.ProductID == "" || req.Quantity <= 0 || req.OrderID == "" {
		return errors.New("invalid reserve request")
	}

	reason := "Reserved for order: " + req.OrderID
	if req.Reason != "" {
		reason = req.Reason
	}

	// 调用领域服务预留库存
	return s.inventoryDomainService.ReserveStock(ctx, req.ProductID, req.Quantity, req.OrderID, reason, req.OperatorID)
}

// ReleaseInventory 释放预留库存
func (s *inventoryServiceImpl) ReleaseInventory(ctx context.Context, req dto.InventoryReleaseRequest) error {
	// 验证请求
	if req.ProductID == "" || req.Quantity <= 0 || req.OrderID == "" {
		return errors.New("invalid release request")
	}

	reason := "Released from order: " + req.OrderID
	if req.Reason != "" {
		reason = req.Reason
	}

	// 调用领域服务释放预留库存
	return s.inventoryDomainService.ReleaseReservedStock(ctx, req.ProductID, req.Quantity, req.OrderID, reason, req.OperatorID)
}

// LockInventory 锁定库存
func (s *inventoryServiceImpl) LockInventory(ctx context.Context, req dto.InventoryLockRequest) error {
	// 验证请求
	if req.ProductID == "" {
		return errors.New("product ID is required")
	}

	reason := "Locked for inventory check"
	if req.Reason != "" {
		reason = req.Reason
	}

	// 调用领域服务锁定库存
	return s.inventoryDomainService.LockInventory(ctx, req.ProductID, reason, req.OperatorID)
}

// UnlockInventory 解锁库存
func (s *inventoryServiceImpl) UnlockInventory(ctx context.Context, req dto.InventoryUnlockRequest) error {
	// 验证请求
	if req.ProductID == "" {
		return errors.New("product ID is required")
	}

	reason := "Inventory check completed"
	if req.Reason != "" {
		reason = req.Reason
	}

	// 调用领域服务解锁库存
	return s.inventoryDomainService.UnlockInventory(ctx, req.ProductID, reason, req.OperatorID)
}

// GetInventory 获取库存详情
func (s *inventoryServiceImpl) GetInventory(ctx context.Context, id string) (dto.InventoryResponse, error) {
	inventory, err := s.inventoryRepo.FindByID(ctx, id)
	if err != nil {
		return dto.InventoryResponse{}, err
	}

	return convertToInventoryResponse(inventory), nil
}

// GetInventoryByProductID 根据商品ID获取库存
func (s *inventoryServiceImpl) GetInventoryByProductID(ctx context.Context, productID string) (dto.InventoryResponse, error) {
	inventory, err := s.inventoryRepo.FindByProductID(ctx, productID)
	if err != nil {
		return dto.InventoryResponse{}, err
	}

	return convertToInventoryResponse(inventory), nil
}

// GetInventoryBySKU 根据SKU获取库存
func (s *inventoryServiceImpl) GetInventoryBySKU(ctx context.Context, sku string) (dto.InventoryResponse, error) {
	inventory, err := s.inventoryRepo.FindBySKU(ctx, sku)
	if err != nil {
		return dto.InventoryResponse{}, err
	}

	return convertToInventoryResponse(inventory), nil
}

// GetLowStockInventories 获取低库存商品列表
func (s *inventoryServiceImpl) GetLowStockInventories(ctx context.Context, page, size int) (dto.InventoryListResponse, error) {
	inventories, total, err := s.inventoryRepo.FindLowStock(ctx, page, size)
	if err != nil {
		return dto.InventoryListResponse{}, err
	}

	return convertToInventoryListResponse(inventories, total, page, size), nil
}

// SearchInventories 搜索库存
func (s *inventoryServiceImpl) SearchInventories(ctx context.Context, req dto.InventoryQueryRequest) (dto.InventoryListResponse, error) {
	var inventories []*model.Inventory
	var total int64
	var err error

	// 根据查询条件搜索
	if req.Keyword != "" {
		inventories, total, err = s.inventoryRepo.Search(ctx, req.Keyword, req.Page, req.Size)
	} else if req.Status != "" {
		inventories, total, err = s.inventoryRepo.FindByStatus(ctx, model.InventoryStatus(req.Status), req.Page, req.Size)
	} else if req.WarehouseID != "" {
		inventories, total, err = s.inventoryRepo.FindByWarehouseID(ctx, req.WarehouseID, req.Page, req.Size)
	} else {
		inventories, total, err = s.inventoryRepo.FindAll(ctx, req.Page, req.Size)
	}

	if err != nil {
		return dto.InventoryListResponse{}, err
	}

	return convertToInventoryListResponse(inventories, total, req.Page, req.Size), nil
}

// GetInventoryOperations 获取库存操作记录
func (s *inventoryServiceImpl) GetInventoryOperations(ctx context.Context, productID string, page, size int) (dto.InventoryOperationListResponse, error) {
	operations, total, err := s.inventoryRepo.FindOperationsByProductID(ctx, productID, page, size)
	if err != nil {
		return dto.InventoryOperationListResponse{}, err
	}

	// 转换为DTO响应
	items := make([]dto.InventoryOperationResponse, 0, len(operations))
	for _, op := range operations {
		items = append(items, convertToOperationResponse(op))
	}

	return dto.InventoryOperationListResponse{
		Items: items,
		Total: total,
		Page:  page,
		Size:  size,
	}, nil
}

// GetInventoryOperationsByOrderID 根据订单ID获取库存操作记录
func (s *inventoryServiceImpl) GetInventoryOperationsByOrderID(ctx context.Context, orderID string) (dto.InventoryOperationListResponse, error) {
	operations, err := s.inventoryRepo.FindOperationsByOrderID(ctx, orderID)
	if err != nil {
		return dto.InventoryOperationListResponse{}, err
	}

	// 转换为DTO响应
	items := make([]dto.InventoryOperationResponse, 0, len(operations))
	for _, op := range operations {
		items = append(items, convertToOperationResponse(op))
	}

	return dto.InventoryOperationListResponse{
		Items: items,
		Total: int64(len(items)),
		Page:  1,
		Size:  len(items),
	}, nil
}

// 辅助方法: 转换库存到DTO响应
func convertToInventoryResponse(inventory *model.Inventory) dto.InventoryResponse {
	return dto.InventoryResponse{
		ID:                 inventory.ID,
		ProductID:          inventory.ProductID,
		SKU:                inventory.SKU,
		AvailableQuantity:  inventory.AvailableQuantity,
		ReservedQuantity:   inventory.ReservedQuantity,
		LowStockThreshold:  inventory.LowStockThreshold,
		Status:             string(inventory.Status),
		LockStatus:         inventory.LockStatus,
		WarehouseID:        inventory.WarehouseID,
		ShelfLocation:      inventory.ShelfLocation,
		LastStockCheckDate: inventory.LastStockCheckDate,
		CreatedAt:          inventory.CreatedAt,
		UpdatedAt:          inventory.UpdatedAt,
	}
}

// 辅助方法: 转换库存操作记录到DTO响应
func convertToOperationResponse(operation model.InventoryOperation) dto.InventoryOperationResponse {
	return dto.InventoryOperationResponse{
		ID:             operation.ID,
		ProductID:      operation.ProductID,
		Type:           string(operation.Type),
		Quantity:       operation.Quantity,
		BeforeStock:    operation.BeforeStock,
		AfterStock:     operation.AfterStock,
		Reason:         operation.Reason,
		RelatedOrderID: operation.RelatedOrderID,
		OperatorID:     operation.OperatorID,
		CreatedAt:      operation.CreatedAt,
	}
}

// 辅助方法: 转换库存列表到DTO响应
func convertToInventoryListResponse(inventories []*model.Inventory, total int64, page, size int) dto.InventoryListResponse {
	items := make([]dto.InventoryBriefResponse, 0, len(inventories))
	for _, inventory := range inventories {
		items = append(items, dto.InventoryBriefResponse{
			ID:                inventory.ID,
			ProductID:         inventory.ProductID,
			SKU:               inventory.SKU,
			AvailableQuantity: inventory.AvailableQuantity,
			ReservedQuantity:  inventory.ReservedQuantity,
			TotalQuantity:     inventory.AvailableQuantity + inventory.ReservedQuantity,
			Status:            string(inventory.Status),
			UpdatedAt:         inventory.UpdatedAt,
		})
	}

	return dto.InventoryListResponse{
		Items: items,
		Total: total,
		Page:  page,
		Size:  size,
	}
}

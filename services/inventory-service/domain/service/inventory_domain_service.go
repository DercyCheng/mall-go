package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"mall-go/services/inventory-service/domain/model"
	"mall-go/services/inventory-service/domain/repository"
)

// InventoryDomainService 库存领域服务，处理库存领域的核心业务逻辑
type InventoryDomainService struct {
	inventoryRepo repository.InventoryRepository
	// 这里可以注入其他依赖的仓储或服务
}

// NewInventoryDomainService 创建库存领域服务的工厂方法
func NewInventoryDomainService(inventoryRepo repository.InventoryRepository) *InventoryDomainService {
	return &InventoryDomainService{
		inventoryRepo: inventoryRepo,
	}
}

// CreateInventory 创建新库存
func (s *InventoryDomainService) CreateInventory(ctx context.Context, productID, sku string, initialQuantity, lowStockThreshold int, warehouseID, shelfLocation string) (*model.Inventory, error) {
	// 检查是否已存在该商品的库存
	existingInventory, err := s.inventoryRepo.FindByProductID(ctx, productID)
	if err == nil && existingInventory != nil {
		return nil, errors.New("inventory already exists for this product")
	}

	// 创建新库存
	inventory, err := model.NewInventory(productID, sku, initialQuantity, lowStockThreshold, warehouseID)
	if err != nil {
		return nil, err
	}

	// 设置货架位置
	inventory.ShelfLocation = shelfLocation

	// 保存库存
	err = s.inventoryRepo.Save(ctx, inventory)
	if err != nil {
		return nil, err
	}

	return inventory, nil
}

// AdjustStock 调整库存数量
func (s *InventoryDomainService) AdjustStock(ctx context.Context, productID string, newQuantity int, reason, operatorID string) error {
	inventory, err := s.inventoryRepo.FindByProductID(ctx, productID)
	if err != nil {
		return err
	}

	if inventory.LockStatus {
		return errors.New("inventory is locked, cannot adjust")
	}

	// 记录调整前的库存
	beforeStock := inventory.AvailableQuantity

	// 更新库存数量
	inventory.AvailableQuantity = newQuantity
	inventory.UpdatedAt = time.Now()

	// 更新库存状态
	if newQuantity <= 0 {
		inventory.Status = model.InventoryStatusOutOfStock
	} else if newQuantity <= inventory.LowStockThreshold {
		inventory.Status = model.InventoryStatusLow
	} else {
		inventory.Status = model.InventoryStatusNormal
	}

	// 创建库存操作记录
	operation := model.InventoryOperation{
		ID:          uuid.New().String(),
		ProductID:   productID,
		Type:        model.OperationTypeAdjust,
		Quantity:    newQuantity - beforeStock,
		BeforeStock: beforeStock,
		AfterStock:  newQuantity,
		Reason:      reason,
		OperatorID:  operatorID,
		CreatedAt:   time.Now(),
	}

	// 添加操作记录到库存历史
	inventory.Operations = append(inventory.Operations, operation)

	// 保存更新后的库存
	err = s.inventoryRepo.Update(ctx, inventory)
	if err != nil {
		return err
	}

	// 单独保存操作记录
	return s.inventoryRepo.SaveOperation(ctx, operation)
}

// InboundStock 入库操作
func (s *InventoryDomainService) InboundStock(ctx context.Context, productID string, quantity int, reason, operatorID string) error {
	if quantity <= 0 {
		return errors.New("inbound quantity must be positive")
	}

	inventory, err := s.inventoryRepo.FindByProductID(ctx, productID)
	if err != nil {
		return err
	}

	if inventory.LockStatus {
		return errors.New("inventory is locked, cannot inbound")
	}

	// 记录入库前的库存
	beforeStock := inventory.AvailableQuantity
	afterStock := beforeStock + quantity

	// 更新库存数量
	inventory.AvailableQuantity = afterStock
	inventory.UpdatedAt = time.Now()

	// 更新库存状态
	if afterStock <= inventory.LowStockThreshold {
		inventory.Status = model.InventoryStatusLow
	} else {
		inventory.Status = model.InventoryStatusNormal
	}

	// 创建库存操作记录
	operation := model.InventoryOperation{
		ID:          uuid.New().String(),
		ProductID:   productID,
		Type:        model.OperationTypeInbound,
		Quantity:    quantity,
		BeforeStock: beforeStock,
		AfterStock:  afterStock,
		Reason:      reason,
		OperatorID:  operatorID,
		CreatedAt:   time.Now(),
	}

	// 添加操作记录到库存历史
	inventory.Operations = append(inventory.Operations, operation)

	// 保存更新后的库存
	err = s.inventoryRepo.Update(ctx, inventory)
	if err != nil {
		return err
	}

	// 单独保存操作记录
	return s.inventoryRepo.SaveOperation(ctx, operation)
}

// OutboundStock 出库操作
func (s *InventoryDomainService) OutboundStock(ctx context.Context, productID string, quantity int, orderID, reason, operatorID string) error {
	if quantity <= 0 {
		return errors.New("outbound quantity must be positive")
	}

	inventory, err := s.inventoryRepo.FindByProductID(ctx, productID)
	if err != nil {
		return err
	}

	if inventory.LockStatus {
		return errors.New("inventory is locked, cannot outbound")
	}

	if inventory.AvailableQuantity < quantity {
		return errors.New("insufficient stock")
	}

	// 记录出库前的库存
	beforeStock := inventory.AvailableQuantity
	afterStock := beforeStock - quantity

	// 更新库存数量
	inventory.AvailableQuantity = afterStock
	inventory.UpdatedAt = time.Now()

	// 更新库存状态
	if afterStock <= 0 {
		inventory.Status = model.InventoryStatusOutOfStock
	} else if afterStock <= inventory.LowStockThreshold {
		inventory.Status = model.InventoryStatusLow
	} else {
		inventory.Status = model.InventoryStatusNormal
	}

	// 创建库存操作记录
	operation := model.InventoryOperation{
		ID:             uuid.New().String(),
		ProductID:      productID,
		Type:           model.OperationTypeOutbound,
		Quantity:       quantity,
		BeforeStock:    beforeStock,
		AfterStock:     afterStock,
		Reason:         reason,
		RelatedOrderID: orderID,
		OperatorID:     operatorID,
		CreatedAt:      time.Now(),
	}

	// 添加操作记录到库存历史
	inventory.Operations = append(inventory.Operations, operation)

	// 保存更新后的库存
	err = s.inventoryRepo.Update(ctx, inventory)
	if err != nil {
		return err
	}

	// 单独保存操作记录
	return s.inventoryRepo.SaveOperation(ctx, operation)
}

// ReserveStock 预留库存
func (s *InventoryDomainService) ReserveStock(ctx context.Context, productID string, quantity int, orderID, reason, operatorID string) error {
	if quantity <= 0 {
		return errors.New("reserve quantity must be positive")
	}

	inventory, err := s.inventoryRepo.FindByProductID(ctx, productID)
	if err != nil {
		return err
	}

	if inventory.LockStatus {
		return errors.New("inventory is locked, cannot reserve")
	}

	if inventory.AvailableQuantity < quantity {
		return errors.New("insufficient stock for reservation")
	}

	// 记录预留前的状态
	beforeAvailable := inventory.AvailableQuantity

	// 更新库存状态
	inventory.AvailableQuantity -= quantity
	inventory.ReservedQuantity += quantity
	inventory.UpdatedAt = time.Now()

	if inventory.AvailableQuantity <= 0 {
		inventory.Status = model.InventoryStatusOutOfStock
	} else if inventory.AvailableQuantity <= inventory.LowStockThreshold {
		inventory.Status = model.InventoryStatusLow
	}

	// 创建库存操作记录
	operation := model.InventoryOperation{
		ID:             uuid.New().String(),
		ProductID:      productID,
		Type:           model.OperationTypeReserve,
		Quantity:       quantity,
		BeforeStock:    beforeAvailable,
		AfterStock:     inventory.AvailableQuantity,
		Reason:         reason,
		RelatedOrderID: orderID,
		OperatorID:     operatorID,
		CreatedAt:      time.Now(),
	}

	// 添加操作记录到库存历史
	inventory.Operations = append(inventory.Operations, operation)

	// 保存更新后的库存
	err = s.inventoryRepo.Update(ctx, inventory)
	if err != nil {
		return err
	}

	// 单独保存操作记录
	return s.inventoryRepo.SaveOperation(ctx, operation)
}

// ReleaseReservedStock 释放预留库存
func (s *InventoryDomainService) ReleaseReservedStock(ctx context.Context, productID string, quantity int, orderID, reason, operatorID string) error {
	if quantity <= 0 {
		return errors.New("release quantity must be positive")
	}

	inventory, err := s.inventoryRepo.FindByProductID(ctx, productID)
	if err != nil {
		return err
	}

	if inventory.ReservedQuantity < quantity {
		return errors.New("reserved quantity is less than release quantity")
	}

	// 记录释放前的状态
	beforeAvailable := inventory.AvailableQuantity

	// 更新库存状态
	inventory.AvailableQuantity += quantity
	inventory.ReservedQuantity -= quantity
	inventory.UpdatedAt = time.Now()

	// 更新库存状态
	if inventory.AvailableQuantity > inventory.LowStockThreshold {
		inventory.Status = model.InventoryStatusNormal
	} else if inventory.AvailableQuantity > 0 {
		inventory.Status = model.InventoryStatusLow
	}

	// 创建库存操作记录
	operation := model.InventoryOperation{
		ID:             uuid.New().String(),
		ProductID:      productID,
		Type:           model.OperationTypeRelease,
		Quantity:       quantity,
		BeforeStock:    beforeAvailable,
		AfterStock:     inventory.AvailableQuantity,
		Reason:         reason,
		RelatedOrderID: orderID,
		OperatorID:     operatorID,
		CreatedAt:      time.Now(),
	}

	// 添加操作记录到库存历史
	inventory.Operations = append(inventory.Operations, operation)

	// 保存更新后的库存
	err = s.inventoryRepo.Update(ctx, inventory)
	if err != nil {
		return err
	}

	// 单独保存操作记录
	return s.inventoryRepo.SaveOperation(ctx, operation)
}

// LockInventory 锁定库存，防止操作
func (s *InventoryDomainService) LockInventory(ctx context.Context, productID string, reason, operatorID string) error {
	inventory, err := s.inventoryRepo.FindByProductID(ctx, productID)
	if err != nil {
		return err
	}

	if inventory.LockStatus {
		return errors.New("inventory is already locked")
	}

	// 锁定库存
	inventory.LockStatus = true
	inventory.Status = model.InventoryStatusLocked
	inventory.UpdatedAt = time.Now()

	// 保存更新后的库存
	return s.inventoryRepo.Update(ctx, inventory)
}

// UnlockInventory 解锁库存
func (s *InventoryDomainService) UnlockInventory(ctx context.Context, productID string, reason, operatorID string) error {
	inventory, err := s.inventoryRepo.FindByProductID(ctx, productID)
	if err != nil {
		return err
	}

	if !inventory.LockStatus {
		return errors.New("inventory is not locked")
	}

	// 解锁库存并根据当前库存数量更新状态
	inventory.LockStatus = false
	if inventory.AvailableQuantity <= 0 {
		inventory.Status = model.InventoryStatusOutOfStock
	} else if inventory.AvailableQuantity <= inventory.LowStockThreshold {
		inventory.Status = model.InventoryStatusLow
	} else {
		inventory.Status = model.InventoryStatusNormal
	}
	inventory.UpdatedAt = time.Now()

	// 保存更新后的库存
	return s.inventoryRepo.Update(ctx, inventory)
}

// createOperationRecord 创建库存操作记录并保存
func (s *InventoryDomainService) createOperationRecord(ctx context.Context, inventory *model.Inventory, opType model.OperationType, quantity int, beforeStock, afterStock int, reason, orderID, operatorID string) error {
	// 创建库存操作记录
	operation := model.InventoryOperation{
		ID:             uuid.New().String(),
		ProductID:      inventory.ProductID,
		Type:           opType,
		Quantity:       quantity,
		BeforeStock:    beforeStock,
		AfterStock:     afterStock,
		Reason:         reason,
		RelatedOrderID: orderID,
		OperatorID:     operatorID,
		CreatedAt:      time.Now(),
	}

	// 添加操作记录到库存历史
	inventory.Operations = append(inventory.Operations, operation)

	// 单独保存操作记录
	return s.inventoryRepo.SaveOperation(ctx, operation)
}

package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// InventoryStatus 库存状态枚举
type InventoryStatus string

const (
	InventoryStatusNormal     InventoryStatus = "normal"       // 正常
	InventoryStatusLow        InventoryStatus = "low"          // 库存不足
	InventoryStatusOutOfStock InventoryStatus = "out_of_stock" // 无库存
	InventoryStatusLocked     InventoryStatus = "locked"       // 锁定
)

// OperationType 操作类型枚举
type OperationType string

const (
	OperationTypeInbound  OperationType = "inbound"  // 入库
	OperationTypeOutbound OperationType = "outbound" // 出库
	OperationTypeReserve  OperationType = "reserve"  // 预留
	OperationTypeRelease  OperationType = "release"  // 释放预留
	OperationTypeAdjust   OperationType = "adjust"   // 调整
)

// InventoryOperation 库存操作记录
type InventoryOperation struct {
	ID             string        // 操作ID
	ProductID      string        // 商品ID
	Type           OperationType // 操作类型
	Quantity       int           // 操作数量
	BeforeStock    int           // 操作前库存
	AfterStock     int           // 操作后库存
	Reason         string        // 操作原因
	RelatedOrderID string        // 相关订单ID
	OperatorID     string        // 操作人ID
	CreatedAt      time.Time     // 创建时间
}

// Inventory 库存聚合根
type Inventory struct {
	ID                 string               // 库存ID
	ProductID          string               // 商品ID
	SKU                string               // 商品SKU
	AvailableQuantity  int                  // 可用库存数量
	ReservedQuantity   int                  // 已预留库存数量
	LowStockThreshold  int                  // 低库存阈值
	LockStatus         bool                 // 锁定状态
	WarehouseID        string               // 仓库ID
	ShelfLocation      string               // 货架位置
	Status             InventoryStatus      // 库存状态
	LastStockCheckDate time.Time            // 最后盘点日期
	Operations         []InventoryOperation // 库存操作历史记录
	CreatedAt          time.Time            // 创建时间
	UpdatedAt          time.Time            // 更新时间
}

// NewInventory 创建新库存的工厂方法
func NewInventory(productID, sku string, initialQuantity, lowStockThreshold int, warehouseID string) (*Inventory, error) {
	if productID == "" {
		return nil, errors.New("product ID cannot be empty")
	}
	if initialQuantity < 0 {
		return nil, errors.New("initial quantity cannot be negative")
	}

	status := InventoryStatusNormal
	if initialQuantity == 0 {
		status = InventoryStatusOutOfStock
	} else if initialQuantity <= lowStockThreshold {
		status = InventoryStatusLow
	}

	inventory := &Inventory{
		ID:                 uuid.New().String(),
		ProductID:          productID,
		SKU:                sku,
		AvailableQuantity:  initialQuantity,
		ReservedQuantity:   0,
		LowStockThreshold:  lowStockThreshold,
		LockStatus:         false,
		WarehouseID:        warehouseID,
		Status:             status,
		LastStockCheckDate: time.Now(),
		Operations:         []InventoryOperation{},
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	// 记录初始入库操作
	if initialQuantity > 0 {
		operation := InventoryOperation{
			ID:          uuid.New().String(),
			ProductID:   productID,
			Type:        OperationTypeInbound,
			Quantity:    initialQuantity,
			BeforeStock: 0,
			AfterStock:  initialQuantity,
			Reason:      "Initial inventory setup",
			CreatedAt:   time.Now(),
		}
		inventory.Operations = append(inventory.Operations, operation)
	}

	return inventory, nil
}

// Inbound 入库操作
func (i *Inventory) Inbound(quantity int, reason, operatorID string) error {
	if quantity <= 0 {
		return errors.New("inbound quantity must be positive")
	}
	if i.LockStatus {
		return errors.New("inventory is locked")
	}

	beforeStock := i.AvailableQuantity
	i.AvailableQuantity += quantity
	i.UpdatedAt = time.Now()

	// 更新库存状态
	i.updateStatus()

	// 记录操作
	operation := InventoryOperation{
		ID:          uuid.New().String(),
		ProductID:   i.ProductID,
		Type:        OperationTypeInbound,
		Quantity:    quantity,
		BeforeStock: beforeStock,
		AfterStock:  i.AvailableQuantity,
		Reason:      reason,
		OperatorID:  operatorID,
		CreatedAt:   time.Now(),
	}
	i.Operations = append(i.Operations, operation)

	return nil
}

// Outbound 出库操作
func (i *Inventory) Outbound(quantity int, orderID, reason, operatorID string) error {
	if quantity <= 0 {
		return errors.New("outbound quantity must be positive")
	}
	if i.LockStatus {
		return errors.New("inventory is locked")
	}
	if i.AvailableQuantity < quantity {
		return errors.New("insufficient inventory")
	}

	beforeStock := i.AvailableQuantity
	i.AvailableQuantity -= quantity
	i.UpdatedAt = time.Now()

	// 更新库存状态
	i.updateStatus()

	// 记录操作
	operation := InventoryOperation{
		ID:             uuid.New().String(),
		ProductID:      i.ProductID,
		Type:           OperationTypeOutbound,
		Quantity:       quantity,
		BeforeStock:    beforeStock,
		AfterStock:     i.AvailableQuantity,
		Reason:         reason,
		RelatedOrderID: orderID,
		OperatorID:     operatorID,
		CreatedAt:      time.Now(),
	}
	i.Operations = append(i.Operations, operation)

	return nil
}

// Reserve 预留库存
func (i *Inventory) Reserve(quantity int, orderID, reason string) error {
	if quantity <= 0 {
		return errors.New("reserve quantity must be positive")
	}
	if i.LockStatus {
		return errors.New("inventory is locked")
	}
	if i.AvailableQuantity < quantity {
		return errors.New("insufficient inventory to reserve")
	}

	beforeStock := i.AvailableQuantity
	i.AvailableQuantity -= quantity
	i.ReservedQuantity += quantity
	i.UpdatedAt = time.Now()

	// 更新库存状态
	i.updateStatus()

	// 记录操作
	operation := InventoryOperation{
		ID:             uuid.New().String(),
		ProductID:      i.ProductID,
		Type:           OperationTypeReserve,
		Quantity:       quantity,
		BeforeStock:    beforeStock,
		AfterStock:     i.AvailableQuantity,
		Reason:         reason,
		RelatedOrderID: orderID,
		CreatedAt:      time.Now(),
	}
	i.Operations = append(i.Operations, operation)

	return nil
}

// ReleaseReservation 释放预留的库存
func (i *Inventory) ReleaseReservation(quantity int, orderID, reason string) error {
	if quantity <= 0 {
		return errors.New("release quantity must be positive")
	}
	if i.ReservedQuantity < quantity {
		return errors.New("reserved quantity is less than requested release quantity")
	}

	beforeStock := i.AvailableQuantity
	i.AvailableQuantity += quantity
	i.ReservedQuantity -= quantity
	i.UpdatedAt = time.Now()

	// 更新库存状态
	i.updateStatus()

	// 记录操作
	operation := InventoryOperation{
		ID:             uuid.New().String(),
		ProductID:      i.ProductID,
		Type:           OperationTypeRelease,
		Quantity:       quantity,
		BeforeStock:    beforeStock,
		AfterStock:     i.AvailableQuantity,
		Reason:         reason,
		RelatedOrderID: orderID,
		CreatedAt:      time.Now(),
	}
	i.Operations = append(i.Operations, operation)

	return nil
}

// ConfirmReservation 确认预留库存 (将预留库存转为已出库)
func (i *Inventory) ConfirmReservation(quantity int, orderID, reason string) error {
	if quantity <= 0 {
		return errors.New("confirmation quantity must be positive")
	}
	if i.ReservedQuantity < quantity {
		return errors.New("reserved quantity is less than requested confirmation quantity")
	}

	i.ReservedQuantity -= quantity
	i.UpdatedAt = time.Now()

	// 记录操作，但不改变可用库存，因为在预留时已经减少
	operation := InventoryOperation{
		ID:             uuid.New().String(),
		ProductID:      i.ProductID,
		Type:           OperationTypeOutbound,
		Quantity:       quantity,
		BeforeStock:    i.AvailableQuantity,
		AfterStock:     i.AvailableQuantity,
		Reason:         "Confirm reservation: " + reason,
		RelatedOrderID: orderID,
		CreatedAt:      time.Now(),
	}
	i.Operations = append(i.Operations, operation)

	return nil
}

// AdjustStock 调整库存
func (i *Inventory) AdjustStock(newQuantity int, reason, operatorID string) error {
	if newQuantity < 0 {
		return errors.New("new quantity cannot be negative")
	}
	if i.LockStatus {
		return errors.New("inventory is locked")
	}

	beforeStock := i.AvailableQuantity
	delta := newQuantity - i.AvailableQuantity
	i.AvailableQuantity = newQuantity
	i.UpdatedAt = time.Now()

	// 更新库存状态
	i.updateStatus()

	// 记录操作
	operation := InventoryOperation{
		ID:          uuid.New().String(),
		ProductID:   i.ProductID,
		Type:        OperationTypeAdjust,
		Quantity:    delta,
		BeforeStock: beforeStock,
		AfterStock:  newQuantity,
		Reason:      reason,
		OperatorID:  operatorID,
		CreatedAt:   time.Now(),
	}
	i.Operations = append(i.Operations, operation)

	return nil
}

// Lock 锁定库存
func (i *Inventory) Lock(reason, operatorID string) error {
	if i.LockStatus {
		return errors.New("inventory is already locked")
	}

	i.LockStatus = true
	i.UpdatedAt = time.Now()

	// 记录操作，使用调整类型
	operation := InventoryOperation{
		ID:          uuid.New().String(),
		ProductID:   i.ProductID,
		Type:        OperationTypeAdjust,
		Quantity:    0,
		BeforeStock: i.AvailableQuantity,
		AfterStock:  i.AvailableQuantity,
		Reason:      "Lock inventory: " + reason,
		OperatorID:  operatorID,
		CreatedAt:   time.Now(),
	}
	i.Operations = append(i.Operations, operation)

	return nil
}

// Unlock 解锁库存
func (i *Inventory) Unlock(reason, operatorID string) error {
	if !i.LockStatus {
		return errors.New("inventory is not locked")
	}

	i.LockStatus = false
	i.UpdatedAt = time.Now()

	// 记录操作，使用调整类型
	operation := InventoryOperation{
		ID:          uuid.New().String(),
		ProductID:   i.ProductID,
		Type:        OperationTypeAdjust,
		Quantity:    0,
		BeforeStock: i.AvailableQuantity,
		AfterStock:  i.AvailableQuantity,
		Reason:      "Unlock inventory: " + reason,
		OperatorID:  operatorID,
		CreatedAt:   time.Now(),
	}
	i.Operations = append(i.Operations, operation)

	return nil
}

// StockCheck 执行库存盘点
func (i *Inventory) StockCheck(actualQuantity int, reason, operatorID string) error {
	if i.LockStatus {
		return errors.New("inventory is locked")
	}

	beforeStock := i.AvailableQuantity
	delta := actualQuantity - beforeStock

	// 更新库存数量
	i.AvailableQuantity = actualQuantity
	i.LastStockCheckDate = time.Now()
	i.UpdatedAt = time.Now()

	// 更新库存状态
	i.updateStatus()

	// 记录操作
	operation := InventoryOperation{
		ID:          uuid.New().String(),
		ProductID:   i.ProductID,
		Type:        OperationTypeAdjust,
		Quantity:    delta,
		BeforeStock: beforeStock,
		AfterStock:  actualQuantity,
		Reason:      "Stock check: " + reason,
		OperatorID:  operatorID,
		CreatedAt:   time.Now(),
	}
	i.Operations = append(i.Operations, operation)

	return nil
}

// UpdateShelfLocation 更新货架位置
func (i *Inventory) UpdateShelfLocation(newLocation string, operatorID string) error {
	i.ShelfLocation = newLocation
	i.UpdatedAt = time.Now()

	// 记录操作
	operation := InventoryOperation{
		ID:          uuid.New().String(),
		ProductID:   i.ProductID,
		Type:        OperationTypeAdjust,
		Quantity:    0,
		BeforeStock: i.AvailableQuantity,
		AfterStock:  i.AvailableQuantity,
		Reason:      "Update shelf location to: " + newLocation,
		OperatorID:  operatorID,
		CreatedAt:   time.Now(),
	}
	i.Operations = append(i.Operations, operation)

	return nil
}

// UpdateLowStockThreshold 更新低库存阈值
func (i *Inventory) UpdateLowStockThreshold(threshold int) error {
	if threshold < 0 {
		return errors.New("threshold cannot be negative")
	}

	i.LowStockThreshold = threshold
	i.UpdatedAt = time.Now()

	// 重新评估库存状态
	i.updateStatus()

	return nil
}

// GetTotalQuantity 获取总库存数量 (可用 + 预留)
func (i *Inventory) GetTotalQuantity() int {
	return i.AvailableQuantity + i.ReservedQuantity
}

// 私有方法：更新库存状态
func (i *Inventory) updateStatus() {
	if i.AvailableQuantity <= 0 {
		i.Status = InventoryStatusOutOfStock
	} else if i.AvailableQuantity <= i.LowStockThreshold {
		i.Status = InventoryStatusLow
	} else {
		i.Status = InventoryStatusNormal
	}

	if i.LockStatus {
		i.Status = InventoryStatusLocked
	}
}

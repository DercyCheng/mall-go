package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// InventoryStatus represents the status of inventory
type InventoryStatus int

const (
	// InventoryStatusNormal indicates normal inventory status
	InventoryStatusNormal InventoryStatus = iota
	// InventoryStatusLocked indicates locked inventory status (e.g., in a pending order)
	InventoryStatusLocked
	// InventoryStatusSold indicates sold inventory status
	InventoryStatusSold
	// InventoryStatusDefective indicates defective inventory status
	InventoryStatusDefective
)

// String converts InventoryStatus to string
func (s InventoryStatus) String() string {
	switch s {
	case InventoryStatusNormal:
		return "normal"
	case InventoryStatusLocked:
		return "locked"
	case InventoryStatusSold:
		return "sold"
	case InventoryStatusDefective:
		return "defective"
	default:
		return "unknown"
	}
}

// InventoryItem represents an inventory item in the system
type InventoryItem struct {
	ID           string          `db:"id"`
	ProductID    string          `db:"product_id"`
	SkuID        string          `db:"sku_id"`
	SkuCode      string          `db:"sku_code"`
	WarehouseID  string          `db:"warehouse_id"`
	Quantity     int             `db:"quantity"`
	LockedCount  int             `db:"locked_count"`
	AvailableQty int             `db:"available_qty"`
	Status       InventoryStatus `db:"status"`
	CreatedAt    time.Time       `db:"created_at"`
	UpdatedAt    time.Time       `db:"updated_at"`
}

// NewInventoryItem creates a new inventory item
func NewInventoryItem(productID, skuID, skuCode, warehouseID string, quantity int) (*InventoryItem, error) {
	if productID == "" || skuID == "" || warehouseID == "" {
		return nil, errors.New("product ID, SKU ID, and warehouse ID are required")
	}

	if quantity < 0 {
		return nil, errors.New("quantity cannot be negative")
	}

	now := time.Now()
	return &InventoryItem{
		ID:           uuid.New().String(),
		ProductID:    productID,
		SkuID:        skuID,
		SkuCode:      skuCode,
		WarehouseID:  warehouseID,
		Quantity:     quantity,
		LockedCount:  0,
		AvailableQty: quantity,
		Status:       InventoryStatusNormal,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// AddStock adds stock to the inventory
func (i *InventoryItem) AddStock(quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be positive")
	}

	i.Quantity += quantity
	i.AvailableQty += quantity
	i.UpdatedAt = time.Now()
	return nil
}

// RemoveStock removes stock from the inventory
func (i *InventoryItem) RemoveStock(quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be positive")
	}

	if i.AvailableQty < quantity {
		return errors.New("insufficient available quantity")
	}

	i.Quantity -= quantity
	i.AvailableQty -= quantity
	i.UpdatedAt = time.Now()
	return nil
}

// LockStock locks stock for a pending order
func (i *InventoryItem) LockStock(quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be positive")
	}

	if i.AvailableQty < quantity {
		return errors.New("insufficient available quantity")
	}

	i.LockedCount += quantity
	i.AvailableQty -= quantity
	i.UpdatedAt = time.Now()
	return nil
}

// UnlockStock unlocks previously locked stock
func (i *InventoryItem) UnlockStock(quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be positive")
	}

	if i.LockedCount < quantity {
		return errors.New("cannot unlock more than locked quantity")
	}

	i.LockedCount -= quantity
	i.AvailableQty += quantity
	i.UpdatedAt = time.Now()
	return nil
}

// ConfirmLock confirms a locked stock (e.g., when order is paid)
func (i *InventoryItem) ConfirmLock(quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be positive")
	}

	if i.LockedCount < quantity {
		return errors.New("cannot confirm more than locked quantity")
	}

	i.LockedCount -= quantity
	i.Quantity -= quantity
	i.UpdatedAt = time.Now()
	return nil
}

// SetStatus sets the inventory status
func (i *InventoryItem) SetStatus(status InventoryStatus) {
	i.Status = status
	i.UpdatedAt = time.Now()
}

// InventoryHistory represents an inventory movement history record
type InventoryHistory struct {
	ID           string    `db:"id"`
	InventoryID  string    `db:"inventory_id"`
	ProductID    string    `db:"product_id"`
	SkuID        string    `db:"sku_id"`
	WarehouseID  string    `db:"warehouse_id"`
	OperationType string    `db:"operation_type"` // "add", "remove", "lock", "unlock", "confirm"
	Quantity     int       `db:"quantity"`
	BeforeQty    int       `db:"before_qty"`      // Quantity before operation
	AfterQty     int       `db:"after_qty"`       // Quantity after operation
	Operator     string    `db:"operator"`        // User ID or system
	OrderID      string    `db:"order_id"`        // Associated order ID if applicable
	Reason       string    `db:"reason"`          // Reason for inventory change
	CreatedAt    time.Time `db:"created_at"`
}

// NewInventoryHistory creates a new inventory history record
func NewInventoryHistory(
	inventoryID, productID, skuID, warehouseID, operationType string,
	quantity, beforeQty, afterQty int,
	operator, orderID, reason string,
) *InventoryHistory {
	return &InventoryHistory{
		ID:           uuid.New().String(),
		InventoryID:  inventoryID,
		ProductID:    productID,
		SkuID:        skuID,
		WarehouseID:  warehouseID,
		OperationType: operationType,
		Quantity:     quantity,
		BeforeQty:    beforeQty,
		AfterQty:     afterQty,
		Operator:     operator,
		OrderID:      orderID,
		Reason:       reason,
		CreatedAt:    time.Now(),
	}
}

// Warehouse represents a warehouse in the system
type Warehouse struct {
	ID          string    `db:"id"`
	Name        string    `db:"name"`
	Code        string    `db:"code"`
	Address     string    `db:"address"`
	ContactName string    `db:"contact_name"`
	ContactPhone string   `db:"contact_phone"`
	Status      int       `db:"status"`       // 0: disabled, 1: enabled
	IsDefault   bool      `db:"is_default"`   // Whether it's the default warehouse
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// NewWarehouse creates a new warehouse
func NewWarehouse(name, code, address, contactName, contactPhone string, isDefault bool) (*Warehouse, error) {
	if name == "" || code == "" {
		return nil, errors.New("name and code are required")
	}

	now := time.Now()
	return &Warehouse{
		ID:          uuid.New().String(),
		Name:        name,
		Code:        code,
		Address:     address,
		ContactName: contactName,
		ContactPhone: contactPhone,
		Status:      1, // Enabled by default
		IsDefault:   isDefault,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// Enable enables the warehouse
func (w *Warehouse) Enable() {
	w.Status = 1
	w.UpdatedAt = time.Now()
}

// Disable disables the warehouse
func (w *Warehouse) Disable() {
	w.Status = 0
	w.UpdatedAt = time.Now()
}

// SetAsDefault sets this warehouse as the default
func (w *Warehouse) SetAsDefault() {
	w.IsDefault = true
	w.UpdatedAt = time.Now()
}

// UnsetDefault removes this warehouse from being the default
func (w *Warehouse) UnsetDefault() {
	w.IsDefault = false
	w.UpdatedAt = time.Now()
}

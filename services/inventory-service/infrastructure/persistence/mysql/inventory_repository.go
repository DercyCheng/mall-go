package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"mall-go/services/inventory-service/domain/model"
	"mall-go/services/inventory-service/domain/repository"
	"mall-go/services/inventory-service/infrastructure/config"
)

// InventoryRepository implements the InventoryRepository interface using MySQL
type InventoryRepository struct {
	db *sqlx.DB
}

// NewInventoryRepository creates a new MySQL inventory repository
func NewInventoryRepository(config config.DatabaseConfig) (repository.InventoryRepository, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
	)

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(config.ConnMaxLifetime) * time.Second)

	repo := &InventoryRepository{
		db: db,
	}

	// Ensure tables exist
	if err := repo.createTablesIfNotExist(); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return repo, nil
}

// createTablesIfNotExist creates the required tables if they don't exist
func (r *InventoryRepository) createTablesIfNotExist() error {
	// Create inventory_items table
	query := `
	CREATE TABLE IF NOT EXISTS inventory_items (
		id VARCHAR(50) PRIMARY KEY,
		product_id VARCHAR(50) NOT NULL,
		sku_id VARCHAR(50) NOT NULL,
		sku_code VARCHAR(100),
		warehouse_id VARCHAR(50) NOT NULL,
		quantity INT NOT NULL DEFAULT 0,
		locked_count INT NOT NULL DEFAULT 0,
		available_qty INT NOT NULL DEFAULT 0,
		status INT NOT NULL DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		INDEX idx_product_sku (product_id, sku_id),
		INDEX idx_warehouse (warehouse_id),
		UNIQUE INDEX idx_product_sku_warehouse (product_id, sku_id, warehouse_id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`

	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

// CreateInventory creates a new inventory item
func (r *InventoryRepository) CreateInventory(ctx context.Context, inventory *model.InventoryItem) error {
	query := `
	INSERT INTO inventory_items (
		id, product_id, sku_id, sku_code, warehouse_id,
		quantity, locked_count, available_qty, status,
		created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		inventory.ID, inventory.ProductID, inventory.SkuID, inventory.SkuCode, inventory.WarehouseID,
		inventory.Quantity, inventory.LockedCount, inventory.AvailableQty, inventory.Status,
		inventory.CreatedAt, inventory.UpdatedAt,
	)
	return err
}

// GetInventory retrieves an inventory item by ID
func (r *InventoryRepository) GetInventory(ctx context.Context, id string) (*model.InventoryItem, error) {
	query := `SELECT * FROM inventory_items WHERE id = ?`
	var item model.InventoryItem
	err := r.db.GetContext(ctx, &item, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("inventory item not found with ID: %s", id)
		}
		return nil, err
	}
	return &item, nil
}

// GetInventoryByProduct retrieves an inventory item by product ID and SKU ID
func (r *InventoryRepository) GetInventoryByProduct(ctx context.Context, productID, skuID string) (*model.InventoryItem, error) {
	// Get from default warehouse if warehouse ID not specified
	query := `
	SELECT i.* FROM inventory_items i
	INNER JOIN warehouses w ON i.warehouse_id = w.id
	WHERE i.product_id = ? AND i.sku_id = ? AND w.is_default = 1
	LIMIT 1`

	var item model.InventoryItem
	err := r.db.GetContext(ctx, &item, query, productID, skuID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("inventory item not found for product: %s and SKU: %s", productID, skuID)
		}
		return nil, err
	}
	return &item, nil
}

// GetInventoriesByProductID retrieves all inventory items for a product
func (r *InventoryRepository) GetInventoriesByProductID(ctx context.Context, productID string) ([]*model.InventoryItem, error) {
	query := `SELECT * FROM inventory_items WHERE product_id = ?`
	var items []*model.InventoryItem
	err := r.db.SelectContext(ctx, &items, query, productID)
	return items, err
}

// GetInventoriesByWarehouseID retrieves all inventory items in a warehouse
func (r *InventoryRepository) GetInventoriesByWarehouseID(ctx context.Context, warehouseID string) ([]*model.InventoryItem, error) {
	query := `SELECT * FROM inventory_items WHERE warehouse_id = ?`
	var items []*model.InventoryItem
	err := r.db.SelectContext(ctx, &items, query, warehouseID)
	return items, err
}

// UpdateInventory updates an existing inventory item
func (r *InventoryRepository) UpdateInventory(ctx context.Context, inventory *model.InventoryItem) error {
	query := `
	UPDATE inventory_items SET
		product_id = ?,
		sku_id = ?,
		sku_code = ?,
		warehouse_id = ?,
		quantity = ?,
		locked_count = ?,
		available_qty = ?,
		status = ?,
		updated_at = ?
	WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query,
		inventory.ProductID, inventory.SkuID, inventory.SkuCode, inventory.WarehouseID,
		inventory.Quantity, inventory.LockedCount, inventory.AvailableQty, inventory.Status,
		time.Now(), inventory.ID,
	)
	return err
}

// AddStock adds stock to an inventory item
func (r *InventoryRepository) AddStock(ctx context.Context, id string, quantity int) error {
	// First, get the current inventory item
	item, err := r.GetInventory(ctx, id)
	if err != nil {
		return err
	}

	// Apply changes to the domain model
	if err := item.AddStock(quantity); err != nil {
		return err
	}

	// Update in the database
	query := `
	UPDATE inventory_items SET
		quantity = ?,
		available_qty = ?,
		updated_at = ?
	WHERE id = ?`

	_, err = r.db.ExecContext(ctx, query,
		item.Quantity, item.AvailableQty, time.Now(), id,
	)
	return err
}

// RemoveStock removes stock from an inventory item
func (r *InventoryRepository) RemoveStock(ctx context.Context, id string, quantity int) error {
	// First, get the current inventory item
	item, err := r.GetInventory(ctx, id)
	if err != nil {
		return err
	}

	// Apply changes to the domain model
	if err := item.RemoveStock(quantity); err != nil {
		return err
	}

	// Update in the database
	query := `
	UPDATE inventory_items SET
		quantity = ?,
		available_qty = ?,
		updated_at = ?
	WHERE id = ?`

	_, err = r.db.ExecContext(ctx, query,
		item.Quantity, item.AvailableQty, time.Now(), id,
	)
	return err
}

// LockStock locks stock for a pending order
func (r *InventoryRepository) LockStock(ctx context.Context, id string, quantity int) error {
	// First, get the current inventory item
	item, err := r.GetInventory(ctx, id)
	if err != nil {
		return err
	}

	// Apply changes to the domain model
	if err := item.LockStock(quantity); err != nil {
		return err
	}

	// Update in the database
	query := `
	UPDATE inventory_items SET
		locked_count = ?,
		available_qty = ?,
		updated_at = ?
	WHERE id = ?`

	_, err = r.db.ExecContext(ctx, query,
		item.LockedCount, item.AvailableQty, time.Now(), id,
	)
	return err
}

// UnlockStock unlocks previously locked stock
func (r *InventoryRepository) UnlockStock(ctx context.Context, id string, quantity int) error {
	// First, get the current inventory item
	item, err := r.GetInventory(ctx, id)
	if err != nil {
		return err
	}

	// Apply changes to the domain model
	if err := item.UnlockStock(quantity); err != nil {
		return err
	}

	// Update in the database
	query := `
	UPDATE inventory_items SET
		locked_count = ?,
		available_qty = ?,
		updated_at = ?
	WHERE id = ?`

	_, err = r.db.ExecContext(ctx, query,
		item.LockedCount, item.AvailableQty, time.Now(), id,
	)
	return err
}

// ConfirmLock confirms a locked stock (e.g., when order is paid)
func (r *InventoryRepository) ConfirmLock(ctx context.Context, id string, quantity int) error {
	// First, get the current inventory item
	item, err := r.GetInventory(ctx, id)
	if err != nil {
		return err
	}

	// Apply changes to the domain model
	if err := item.ConfirmLock(quantity); err != nil {
		return err
	}

	// Update in the database
	query := `
	UPDATE inventory_items SET
		locked_count = ?,
		quantity = ?,
		updated_at = ?
	WHERE id = ?`

	_, err = r.db.ExecContext(ctx, query,
		item.LockedCount, item.Quantity, time.Now(), id,
	)
	return err
}

// SetStatus sets the inventory status
func (r *InventoryRepository) SetStatus(ctx context.Context, id string, status model.InventoryStatus) error {
	// First, get the current inventory item
	item, err := r.GetInventory(ctx, id)
	if err != nil {
		return err
	}

	// Apply changes to the domain model
	item.SetStatus(status)

	// Update in the database
	query := `
	UPDATE inventory_items SET
		status = ?,
		updated_at = ?
	WHERE id = ?`

	_, err = r.db.ExecContext(ctx, query, status, time.Now(), id)
	return err
}

// DeleteInventory deletes an inventory item (soft delete)
func (r *InventoryRepository) DeleteInventory(ctx context.Context, id string) error {
	// For now, we'll implement this as a hard delete
	query := `DELETE FROM inventory_items WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// ListInventories lists all inventory items with pagination
func (r *InventoryRepository) ListInventories(ctx context.Context, page, pageSize int) ([]*model.InventoryItem, int, error) {
	// Calculate offset
	offset := (page - 1) * pageSize

	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM inventory_items`
	err := r.db.GetContext(ctx, &total, countQuery)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	query := `SELECT * FROM inventory_items LIMIT ? OFFSET ?`
	var items []*model.InventoryItem
	err = r.db.SelectContext(ctx, &items, query, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}


// WarehouseRepository implements the WarehouseRepository interface using MySQL
type WarehouseRepository struct {
	db *sqlx.DB
}

// NewWarehouseRepository creates a new MySQL warehouse repository
func NewWarehouseRepository(config config.DatabaseConfig) (repository.WarehouseRepository, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
	)

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(config.ConnMaxLifetime) * time.Second)

	repo := &WarehouseRepository{
		db: db,
	}

	// Ensure tables exist
	if err := repo.createTablesIfNotExist(); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return repo, nil
}

// createTablesIfNotExist creates the required tables if they don't exist
func (r *WarehouseRepository) createTablesIfNotExist() error {
	// Create warehouses table
	query := `
	CREATE TABLE IF NOT EXISTS warehouses (
		id VARCHAR(50) PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		code VARCHAR(50) NOT NULL,
		address VARCHAR(255),
		contact_name VARCHAR(100),
		contact_phone VARCHAR(20),
		status INT NOT NULL DEFAULT 1,
		is_default BOOLEAN NOT NULL DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		INDEX idx_code (code),
		INDEX idx_default (is_default)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`

	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

// CreateWarehouse creates a new warehouse
func (r *WarehouseRepository) CreateWarehouse(ctx context.Context, warehouse *model.Warehouse) error {
	// Start transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// If this is the default warehouse, first unset any existing default
	if warehouse.IsDefault {
		unsetQuery := `UPDATE warehouses SET is_default = FALSE WHERE is_default = TRUE`
		_, err := tx.ExecContext(ctx, unsetQuery)
		if err != nil {
			return err
		}
	}

	// Insert new warehouse
	query := `
	INSERT INTO warehouses (
		id, name, code, address, contact_name, contact_phone,
		status, is_default, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err = tx.ExecContext(ctx, query,
		warehouse.ID, warehouse.Name, warehouse.Code, warehouse.Address,
		warehouse.ContactName, warehouse.ContactPhone, warehouse.Status,
		warehouse.IsDefault, warehouse.CreatedAt, warehouse.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetWarehouse retrieves a warehouse by ID
func (r *WarehouseRepository) GetWarehouse(ctx context.Context, id string) (*model.Warehouse, error) {
	query := `SELECT * FROM warehouses WHERE id = ?`
	var warehouse model.Warehouse
	err := r.db.GetContext(ctx, &warehouse, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("warehouse not found with ID: %s", id)
		}
		return nil, err
	}
	return &warehouse, nil
}

// GetWarehouseByCode retrieves a warehouse by code
func (r *WarehouseRepository) GetWarehouseByCode(ctx context.Context, code string) (*model.Warehouse, error) {
	query := `SELECT * FROM warehouses WHERE code = ?`
	var warehouse model.Warehouse
	err := r.db.GetContext(ctx, &warehouse, query, code)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("warehouse not found with code: %s", code)
		}
		return nil, err
	}
	return &warehouse, nil
}

// GetDefaultWarehouse retrieves the default warehouse
func (r *WarehouseRepository) GetDefaultWarehouse(ctx context.Context) (*model.Warehouse, error) {
	query := `SELECT * FROM warehouses WHERE is_default = TRUE LIMIT 1`
	var warehouse model.Warehouse
	err := r.db.GetContext(ctx, &warehouse, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no default warehouse found")
		}
		return nil, err
	}
	return &warehouse, nil
}

// UpdateWarehouse updates an existing warehouse
func (r *WarehouseRepository) UpdateWarehouse(ctx context.Context, warehouse *model.Warehouse) error {
	// Start transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// If this is the default warehouse, first unset any existing default
	if warehouse.IsDefault {
		unsetQuery := `UPDATE warehouses SET is_default = FALSE WHERE is_default = TRUE AND id != ?`
		_, err := tx.ExecContext(ctx, unsetQuery, warehouse.ID)
		if err != nil {
			return err
		}
	}

	// Update warehouse
	query := `
	UPDATE warehouses SET
		name = ?,
		code = ?,
		address = ?,
		contact_name = ?,
		contact_phone = ?,
		status = ?,
		is_default = ?,
		updated_at = ?
	WHERE id = ?`

	_, err = tx.ExecContext(ctx, query,
		warehouse.Name, warehouse.Code, warehouse.Address,
		warehouse.ContactName, warehouse.ContactPhone, warehouse.Status,
		warehouse.IsDefault, time.Now(), warehouse.ID,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// DeleteWarehouse deletes a warehouse (soft delete)
func (r *WarehouseRepository) DeleteWarehouse(ctx context.Context, id string) error {
	// For now, implement as hard delete, but check if it's the default first
	query1 := `SELECT is_default FROM warehouses WHERE id = ?`
	var isDefault bool
	err := r.db.GetContext(ctx, &isDefault, query1, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("warehouse not found with ID: %s", id)
		}
		return err
	}

	// Cannot delete the default warehouse
	if isDefault {
		return fmt.Errorf("cannot delete the default warehouse")
	}

	// Delete the warehouse
	query2 := `DELETE FROM warehouses WHERE id = ?`
	_, err = r.db.ExecContext(ctx, query2, id)
	return err
}

// ListWarehouses lists all warehouses
func (r *WarehouseRepository) ListWarehouses(ctx context.Context) ([]*model.Warehouse, error) {
	query := `SELECT * FROM warehouses ORDER BY is_default DESC, name ASC`
	var warehouses []*model.Warehouse
	err := r.db.SelectContext(ctx, &warehouses, query)
	return warehouses, err
}

// SetDefaultWarehouse sets a warehouse as the default
func (r *WarehouseRepository) SetDefaultWarehouse(ctx context.Context, id string) error {
	// Start transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Unset any existing default
	query1 := `UPDATE warehouses SET is_default = FALSE WHERE is_default = TRUE`
	_, err = tx.ExecContext(ctx, query1)
	if err != nil {
		return err
	}

	// Set new default
	query2 := `UPDATE warehouses SET is_default = TRUE WHERE id = ?`
	result, err := tx.ExecContext(ctx, query2, id)
	if err != nil {
		return err
	}

	// Check if the warehouse exists
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("warehouse not found with ID: %s", id)
	}

	return tx.Commit()
}

// EnableWarehouse enables a warehouse
func (r *WarehouseRepository) EnableWarehouse(ctx context.Context, id string) error {
	query := `UPDATE warehouses SET status = 1, updated_at = ? WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return err
	}

	// Check if the warehouse exists
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("warehouse not found with ID: %s", id)
	}

	return nil
}

// DisableWarehouse disables a warehouse
func (r *WarehouseRepository) DisableWarehouse(ctx context.Context, id string) error {
	// Check if it's the default warehouse
	query1 := `SELECT is_default FROM warehouses WHERE id = ?`
	var isDefault bool
	err := r.db.GetContext(ctx, &isDefault, query1, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("warehouse not found with ID: %s", id)
		}
		return err
	}

	// Cannot disable the default warehouse
	if isDefault {
		return fmt.Errorf("cannot disable the default warehouse")
	}

	// Disable the warehouse
	query2 := `UPDATE warehouses SET status = 0, updated_at = ? WHERE id = ?`
	_, err = r.db.ExecContext(ctx, query2, time.Now(), id)
	return err
}


// InventoryHistoryRepository implements the InventoryHistoryRepository interface using MySQL
type InventoryHistoryRepository struct {
	db *sqlx.DB
}

// NewInventoryHistoryRepository creates a new MySQL inventory history repository
func NewInventoryHistoryRepository(config config.DatabaseConfig) (repository.InventoryHistoryRepository, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
	)

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(config.ConnMaxLifetime) * time.Second)

	repo := &InventoryHistoryRepository{
		db: db,
	}

	// Ensure tables exist
	if err := repo.createTablesIfNotExist(); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return repo, nil
}

// createTablesIfNotExist creates the required tables if they don't exist
func (r *InventoryHistoryRepository) createTablesIfNotExist() error {
	// Create inventory_history table
	query := `
	CREATE TABLE IF NOT EXISTS inventory_history (
		id VARCHAR(50) PRIMARY KEY,
		inventory_id VARCHAR(50) NOT NULL,
		product_id VARCHAR(50) NOT NULL,
		sku_id VARCHAR(50) NOT NULL,
		warehouse_id VARCHAR(50) NOT NULL,
		operation_type VARCHAR(20) NOT NULL,
		quantity INT NOT NULL,
		before_qty INT NOT NULL,
		after_qty INT NOT NULL,
		operator VARCHAR(50) NOT NULL,
		order_id VARCHAR(50),
		reason VARCHAR(255),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		INDEX idx_inventory_id (inventory_id),
		INDEX idx_product_id (product_id),
		INDEX idx_warehouse_id (warehouse_id),
		INDEX idx_order_id (order_id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`

	_, err := r.db.Exec(query)
	return err
}

// CreateHistory creates a new inventory history record
func (r *InventoryHistoryRepository) CreateHistory(ctx context.Context, history *model.InventoryHistory) error {
	query := `
	INSERT INTO inventory_history (
		id, inventory_id, product_id, sku_id, warehouse_id,
		operation_type, quantity, before_qty, after_qty,
		operator, order_id, reason, created_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		history.ID, history.InventoryID, history.ProductID, history.SkuID, history.WarehouseID,
		history.OperationType, history.Quantity, history.BeforeQty, history.AfterQty,
		history.Operator, history.OrderID, history.Reason, history.CreatedAt,
	)
	return err
}

// GetHistoryByInventoryID retrieves all history records for an inventory item
func (r *InventoryHistoryRepository) GetHistoryByInventoryID(ctx context.Context, inventoryID string, page, pageSize int) ([]*model.InventoryHistory, int, error) {
	// Calculate offset
	offset := (page - 1) * pageSize

	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM inventory_history WHERE inventory_id = ?`
	err := r.db.GetContext(ctx, &total, countQuery, inventoryID)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	query := `
	SELECT * FROM inventory_history 
	WHERE inventory_id = ? 
	ORDER BY created_at DESC 
	LIMIT ? OFFSET ?`

	var history []*model.InventoryHistory
	err = r.db.SelectContext(ctx, &history, query, inventoryID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	return history, total, nil
}

// GetHistoryByProductID retrieves all history records for a product
func (r *InventoryHistoryRepository) GetHistoryByProductID(ctx context.Context, productID string, page, pageSize int) ([]*model.InventoryHistory, int, error) {
	// Calculate offset
	offset := (page - 1) * pageSize

	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM inventory_history WHERE product_id = ?`
	err := r.db.GetContext(ctx, &total, countQuery, productID)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	query := `
	SELECT * FROM inventory_history 
	WHERE product_id = ? 
	ORDER BY created_at DESC 
	LIMIT ? OFFSET ?`

	var history []*model.InventoryHistory
	err = r.db.SelectContext(ctx, &history, query, productID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	return history, total, nil
}

// GetHistoryByWarehouseID retrieves all history records for a warehouse
func (r *InventoryHistoryRepository) GetHistoryByWarehouseID(ctx context.Context, warehouseID string, page, pageSize int) ([]*model.InventoryHistory, int, error) {
	// Calculate offset
	offset := (page - 1) * pageSize

	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM inventory_history WHERE warehouse_id = ?`
	err := r.db.GetContext(ctx, &total, countQuery, warehouseID)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	query := `
	SELECT * FROM inventory_history 
	WHERE warehouse_id = ? 
	ORDER BY created_at DESC 
	LIMIT ? OFFSET ?`

	var history []*model.InventoryHistory
	err = r.db.SelectContext(ctx, &history, query, warehouseID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	return history, total, nil
}

// GetHistoryByOrderID retrieves all history records for an order
func (r *InventoryHistoryRepository) GetHistoryByOrderID(ctx context.Context, orderID string) ([]*model.InventoryHistory, error) {
	query := `
	SELECT * FROM inventory_history 
	WHERE order_id = ? 
	ORDER BY created_at DESC`

	var history []*model.InventoryHistory
	err := r.db.SelectContext(ctx, &history, query, orderID)
	return history, err
}

// ListHistory lists all inventory history records with pagination
func (r *InventoryHistoryRepository) ListHistory(ctx context.Context, page, pageSize int) ([]*model.InventoryHistory, int, error) {
	// Calculate offset
	offset := (page - 1) * pageSize

	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM inventory_history`
	err := r.db.GetContext(ctx, &total, countQuery)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	query := `
	SELECT * FROM inventory_history 
	ORDER BY created_at DESC 
	LIMIT ? OFFSET ?`

	var history []*model.InventoryHistory
	err = r.db.SelectContext(ctx, &history, query, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	return history, total, nil
}

package persistence

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"mall-go/services/inventory-service/domain/model"
	"mall-go/services/inventory-service/domain/repository"
)

// InventoryRepository 库存仓储实现
type InventoryRepository struct {
	db *sql.DB
}

// NewInventoryRepository 创建库存仓储
func NewInventoryRepository(db *sql.DB) *InventoryRepository {
	return &InventoryRepository{db: db}
}

// Save 保存库存
func (r *InventoryRepository) Save(ctx context.Context, inventory *model.Inventory) error {
	// 开始事务
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 序列化操作记录
	operationsJSON, err := json.Marshal(inventory.Operations)
	if err != nil {
		return err
	}

	// 插入库存
	query := `
		INSERT INTO inventories (
			id, product_id, sku, 
			available_quantity, reserved_quantity, 
			low_stock_threshold, lock_status,
			warehouse_id, shelf_location,
			status, last_stock_check_date,
			operations, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = tx.ExecContext(
		ctx,
		query,
		inventory.ID, inventory.ProductID, inventory.SKU,
		inventory.AvailableQuantity, inventory.ReservedQuantity,
		inventory.LowStockThreshold, inventory.LockStatus,
		inventory.WarehouseID, inventory.ShelfLocation,
		inventory.Status, inventory.LastStockCheckDate,
		operationsJSON, inventory.CreatedAt, inventory.UpdatedAt,
	)

	if err != nil {
		return err
	}

	// 提交事务
	return tx.Commit()
}

// FindByID 根据ID查询库存
func (r *InventoryRepository) FindByID(ctx context.Context, id string) (*model.Inventory, error) {
	query := `
		SELECT 
			id, product_id, sku, 
			available_quantity, reserved_quantity, 
			low_stock_threshold, lock_status,
			warehouse_id, shelf_location,
			status, last_stock_check_date,
			operations, created_at, updated_at
		FROM inventories
		WHERE id = ? AND deleted_at IS NULL
	`

	var inventory model.Inventory
	var operationsJSON []byte
	var status string

	err := r.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&inventory.ID, &inventory.ProductID, &inventory.SKU,
		&inventory.AvailableQuantity, &inventory.ReservedQuantity,
		&inventory.LowStockThreshold, &inventory.LockStatus,
		&inventory.WarehouseID, &inventory.ShelfLocation,
		&status, &inventory.LastStockCheckDate,
		&operationsJSON, &inventory.CreatedAt, &inventory.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("inventory not found: %w", err)
		}
		return nil, err
	}

	// 设置状态
	inventory.Status = model.InventoryStatus(status)

	// 反序列化操作记录
	if err = json.Unmarshal(operationsJSON, &inventory.Operations); err != nil {
		return nil, err
	}

	return &inventory, nil
}

// FindByProductID 根据商品ID查询库存
func (r *InventoryRepository) FindByProductID(ctx context.Context, productID string) (*model.Inventory, error) {
	query := `
		SELECT id FROM inventories
		WHERE product_id = ? AND deleted_at IS NULL
	`

	var id string
	err := r.db.QueryRowContext(ctx, query, productID).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("inventory not found for product %s: %w", productID, err)
		}
		return nil, err
	}

	return r.FindByID(ctx, id)
}

// FindBySKU 根据SKU查询库存
func (r *InventoryRepository) FindBySKU(ctx context.Context, sku string) (*model.Inventory, error) {
	query := `
		SELECT id FROM inventories
		WHERE sku = ? AND deleted_at IS NULL
	`

	var id string
	err := r.db.QueryRowContext(ctx, query, sku).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("inventory not found for SKU %s: %w", sku, err)
		}
		return nil, err
	}

	return r.FindByID(ctx, id)
}

// FindByWarehouseID 根据仓库ID查询库存列表
func (r *InventoryRepository) FindByWarehouseID(ctx context.Context, warehouseID string, page, size int) ([]*model.Inventory, int64, error) {
	// 计算总数
	var total int64
	countQuery := `
		SELECT COUNT(*) FROM inventories
		WHERE warehouse_id = ? AND deleted_at IS NULL
	`
	err := r.db.QueryRowContext(ctx, countQuery, warehouseID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 如果没有数据，直接返回空列表
	if total == 0 {
		return []*model.Inventory{}, 0, nil
	}

	// 分页查询库存ID
	offset := (page - 1) * size
	query := `
		SELECT id FROM inventories
		WHERE warehouse_id = ? AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, warehouseID, size, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// 收集库存ID
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, 0, err
		}
		ids = append(ids, id)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	// 批量查询库存详情
	inventories := make([]*model.Inventory, 0, len(ids))
	for _, id := range ids {
		inventory, err := r.FindByID(ctx, id)
		if err != nil {
			return nil, 0, err
		}
		inventories = append(inventories, inventory)
	}

	return inventories, total, nil
}

// FindByStatus 根据库存状态查询库存列表
func (r *InventoryRepository) FindByStatus(ctx context.Context, status model.InventoryStatus, page, size int) ([]*model.Inventory, int64, error) {
	// 计算总数
	var total int64
	countQuery := `
		SELECT COUNT(*) FROM inventories
		WHERE status = ? AND deleted_at IS NULL
	`
	err := r.db.QueryRowContext(ctx, countQuery, status).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 如果没有数据，直接返回空列表
	if total == 0 {
		return []*model.Inventory{}, 0, nil
	}

	// 分页查询库存ID
	offset := (page - 1) * size
	query := `
		SELECT id FROM inventories
		WHERE status = ? AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, status, size, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// 收集库存ID
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, 0, err
		}
		ids = append(ids, id)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	// 批量查询库存详情
	inventories := make([]*model.Inventory, 0, len(ids))
	for _, id := range ids {
		inventory, err := r.FindByID(ctx, id)
		if err != nil {
			return nil, 0, err
		}
		inventories = append(inventories, inventory)
	}

	return inventories, total, nil
}

// FindLowStock 查询低库存商品
func (r *InventoryRepository) FindLowStock(ctx context.Context, page, size int) ([]*model.Inventory, int64, error) {
	// 计算总数
	var total int64
	countQuery := `
		SELECT COUNT(*) FROM inventories
		WHERE (status = 'low' OR available_quantity <= low_stock_threshold)
		AND deleted_at IS NULL
	`
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 如果没有数据，直接返回空列表
	if total == 0 {
		return []*model.Inventory{}, 0, nil
	}

	// 分页查询库存ID
	offset := (page - 1) * size
	query := `
		SELECT id FROM inventories
		WHERE (status = 'low' OR available_quantity <= low_stock_threshold)
		AND deleted_at IS NULL
		ORDER BY available_quantity ASC, created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, size, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// 收集库存ID
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, 0, err
		}
		ids = append(ids, id)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	// 批量查询库存详情
	inventories := make([]*model.Inventory, 0, len(ids))
	for _, id := range ids {
		inventory, err := r.FindByID(ctx, id)
		if err != nil {
			return nil, 0, err
		}
		inventories = append(inventories, inventory)
	}

	return inventories, total, nil
}

// FindAll 分页查询所有库存
func (r *InventoryRepository) FindAll(ctx context.Context, page, size int) ([]*model.Inventory, int64, error) {
	// 计算总数
	var total int64
	countQuery := `
		SELECT COUNT(*) FROM inventories
		WHERE deleted_at IS NULL
	`
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 如果没有数据，直接返回空列表
	if total == 0 {
		return []*model.Inventory{}, 0, nil
	}

	// 分页查询库存ID
	offset := (page - 1) * size
	query := `
		SELECT id FROM inventories
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, size, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// 收集库存ID
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, 0, err
		}
		ids = append(ids, id)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	// 批量查询库存详情
	inventories := make([]*model.Inventory, 0, len(ids))
	for _, id := range ids {
		inventory, err := r.FindByID(ctx, id)
		if err != nil {
			return nil, 0, err
		}
		inventories = append(inventories, inventory)
	}

	return inventories, total, nil
}

// Search 搜索库存
func (r *InventoryRepository) Search(ctx context.Context, keyword string, page, size int) ([]*model.Inventory, int64, error) {
	// 计算总数
	var total int64
	countQuery := `
		SELECT COUNT(*) FROM inventories
		WHERE (product_id LIKE ? OR sku LIKE ? OR shelf_location LIKE ?)
		AND deleted_at IS NULL
	`
	likeKeyword := "%" + keyword + "%"
	err := r.db.QueryRowContext(ctx, countQuery, likeKeyword, likeKeyword, likeKeyword).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 如果没有数据，直接返回空列表
	if total == 0 {
		return []*model.Inventory{}, 0, nil
	}

	// 分页查询库存ID
	offset := (page - 1) * size
	query := `
		SELECT id FROM inventories
		WHERE (product_id LIKE ? OR sku LIKE ? OR shelf_location LIKE ?)
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, likeKeyword, likeKeyword, likeKeyword, size, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// 收集库存ID
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, 0, err
		}
		ids = append(ids, id)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	// 批量查询库存详情
	inventories := make([]*model.Inventory, 0, len(ids))
	for _, id := range ids {
		inventory, err := r.FindByID(ctx, id)
		if err != nil {
			return nil, 0, err
		}
		inventories = append(inventories, inventory)
	}

	return inventories, total, nil
}

// Update 更新库存
func (r *InventoryRepository) Update(ctx context.Context, inventory *model.Inventory) error {
	// 开始事务
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 序列化操作记录
	operationsJSON, err := json.Marshal(inventory.Operations)
	if err != nil {
		return err
	}

	// 更新库存
	query := `
		UPDATE inventories SET
			product_id = ?, sku = ?, 
			available_quantity = ?, reserved_quantity = ?, 
			low_stock_threshold = ?, lock_status = ?,
			warehouse_id = ?, shelf_location = ?,
			status = ?, last_stock_check_date = ?,
			operations = ?, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`

	_, err = tx.ExecContext(
		ctx,
		query,
		inventory.ProductID, inventory.SKU,
		inventory.AvailableQuantity, inventory.ReservedQuantity,
		inventory.LowStockThreshold, inventory.LockStatus,
		inventory.WarehouseID, inventory.ShelfLocation,
		inventory.Status, inventory.LastStockCheckDate,
		operationsJSON, time.Now(),
		inventory.ID,
	)

	if err != nil {
		return err
	}

	// 提交事务
	return tx.Commit()
}

// Delete 删除库存
func (r *InventoryRepository) Delete(ctx context.Context, id string) error {
	query := `
		UPDATE inventories SET
			deleted_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	return err
}

// SaveOperation 保存库存操作记录
func (r *InventoryRepository) SaveOperation(ctx context.Context, operation model.InventoryOperation) error {
	// 保存到单独的操作记录表
	query := `
		INSERT INTO inventory_operations (
			id, product_id, type, quantity,
			before_stock, after_stock, reason,
			related_order_id, operator_id, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		operation.ID, operation.ProductID, operation.Type, operation.Quantity,
		operation.BeforeStock, operation.AfterStock, operation.Reason,
		operation.RelatedOrderID, operation.OperatorID, operation.CreatedAt,
	)

	return err
}

// FindOperationsByProductID 根据商品ID查询库存操作记录
func (r *InventoryRepository) FindOperationsByProductID(ctx context.Context, productID string, page, size int) ([]model.InventoryOperation, int64, error) {
	// 计算总数
	var total int64
	countQuery := `
		SELECT COUNT(*) FROM inventory_operations
		WHERE product_id = ?
	`
	err := r.db.QueryRowContext(ctx, countQuery, productID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 如果没有数据，直接返回空列表
	if total == 0 {
		return []model.InventoryOperation{}, 0, nil
	}

	// 分页查询操作记录
	offset := (page - 1) * size
	query := `
		SELECT 
			id, product_id, type, quantity,
			before_stock, after_stock, reason,
			related_order_id, operator_id, created_at
		FROM inventory_operations
		WHERE product_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, productID, size, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// 收集操作记录
	operations := make([]model.InventoryOperation, 0, size)
	for rows.Next() {
		var op model.InventoryOperation
		var opType, relatedOrderID, operatorID sql.NullString

		if err := rows.Scan(
			&op.ID, &op.ProductID, &opType, &op.Quantity,
			&op.BeforeStock, &op.AfterStock, &op.Reason,
			&relatedOrderID, &operatorID, &op.CreatedAt,
		); err != nil {
			return nil, 0, err
		}

		// 设置可空字段
		if opType.Valid {
			op.Type = model.OperationType(opType.String)
		}
		if relatedOrderID.Valid {
			op.RelatedOrderID = relatedOrderID.String
		}
		if operatorID.Valid {
			op.OperatorID = operatorID.String
		}

		operations = append(operations, op)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return operations, total, nil
}

// FindOperationsByOrderID 根据订单ID查询库存操作记录
func (r *InventoryRepository) FindOperationsByOrderID(ctx context.Context, orderID string) ([]model.InventoryOperation, error) {
	query := `
		SELECT 
			id, product_id, type, quantity,
			before_stock, after_stock, reason,
			related_order_id, operator_id, created_at
		FROM inventory_operations
		WHERE related_order_id = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 收集操作记录
	operations := make([]model.InventoryOperation, 0)
	for rows.Next() {
		var op model.InventoryOperation
		var opType, relatedOrderID, operatorID sql.NullString

		if err := rows.Scan(
			&op.ID, &op.ProductID, &opType, &op.Quantity,
			&op.BeforeStock, &op.AfterStock, &op.Reason,
			&relatedOrderID, &operatorID, &op.CreatedAt,
		); err != nil {
			return nil, err
		}

		// 设置可空字段
		if opType.Valid {
			op.Type = model.OperationType(opType.String)
		}
		if relatedOrderID.Valid {
			op.RelatedOrderID = relatedOrderID.String
		}
		if operatorID.Valid {
			op.OperatorID = operatorID.String
		}

		operations = append(operations, op)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return operations, nil
}

// CountByStatus 统计各状态库存数量
func (r *InventoryRepository) CountByStatus(ctx context.Context) (map[model.InventoryStatus]int64, error) {
	result := make(map[model.InventoryStatus]int64)

	// 查询各种状态的库存数量
	var counts []struct {
		Status string
		Count  int64
	}

	// 使用原生SQL查询，因为r.db是*sql.DB而不是*gorm.DB
	rows, err := r.db.QueryContext(ctx,
		"SELECT status, COUNT(*) as count FROM inventories WHERE deleted_at IS NULL GROUP BY status")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 处理结果集
	for rows.Next() {
		var count struct {
			Status string
			Count  int64
		}
		if err := rows.Scan(&count.Status, &count.Count); err != nil {
			return nil, err
		}
		counts = append(counts, count)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// 转换为领域模型状态
	for _, count := range counts {
		result[model.InventoryStatus(count.Status)] = count.Count
	}

	return result, nil
}

// Ensure InventoryRepository implements repository.InventoryRepository
var _ repository.InventoryRepository = (*InventoryRepository)(nil)

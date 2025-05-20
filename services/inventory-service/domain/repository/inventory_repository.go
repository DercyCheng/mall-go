package repository

import (
	"context"

	"mall-go/services/inventory-service/domain/model"
)

// InventoryRepository 定义了库存聚合根的仓储接口
type InventoryRepository interface {
	// Save 保存库存
	Save(ctx context.Context, inventory *model.Inventory) error

	// FindByID 根据ID查询库存
	FindByID(ctx context.Context, id string) (*model.Inventory, error)

	// FindByProductID 根据商品ID查询库存
	FindByProductID(ctx context.Context, productID string) (*model.Inventory, error)

	// FindBySKU 根据SKU查询库存
	FindBySKU(ctx context.Context, sku string) (*model.Inventory, error)

	// FindByWarehouseID 根据仓库ID查询库存列表
	FindByWarehouseID(ctx context.Context, warehouseID string, page, size int) ([]*model.Inventory, int64, error)

	// FindByStatus 根据库存状态查询库存列表
	FindByStatus(ctx context.Context, status model.InventoryStatus, page, size int) ([]*model.Inventory, int64, error)

	// FindLowStock 查询低库存商品
	FindLowStock(ctx context.Context, page, size int) ([]*model.Inventory, int64, error)

	// FindAll 分页查询所有库存
	FindAll(ctx context.Context, page, size int) ([]*model.Inventory, int64, error)

	// Search 搜索库存
	Search(ctx context.Context, keyword string, page, size int) ([]*model.Inventory, int64, error)

	// Update 更新库存
	Update(ctx context.Context, inventory *model.Inventory) error

	// Delete 删除库存
	Delete(ctx context.Context, id string) error

	// SaveOperation 保存库存操作记录
	SaveOperation(ctx context.Context, operation model.InventoryOperation) error

	// FindOperationsByProductID 根据商品ID查询库存操作记录
	FindOperationsByProductID(ctx context.Context, productID string, page, size int) ([]model.InventoryOperation, int64, error)

	// FindOperationsByOrderID 根据订单ID查询库存操作记录
	FindOperationsByOrderID(ctx context.Context, orderID string) ([]model.InventoryOperation, error)

	// CountByStatus 统计各状态库存数量
	CountByStatus(ctx context.Context) (map[model.InventoryStatus]int64, error)
}

package repository

import (
	"context"

	"mall-go/services/order-service/domain/model"
)

// OrderRepository 定义了订单聚合根的仓储接口
type OrderRepository interface {
	// Save 保存订单
	Save(ctx context.Context, order *model.Order) error

	// FindByID 根据ID查询订单
	FindByID(ctx context.Context, id string) (*model.Order, error)

	// FindByOrderSN 根据订单编号查询订单
	FindByOrderSN(ctx context.Context, orderSN string) (*model.Order, error)

	// FindByUserID 根据用户ID查询订单列表
	FindByUserID(ctx context.Context, userID string, page, size int) ([]*model.Order, int64, error)

	// FindByStatus 根据订单状态查询订单列表
	FindByStatus(ctx context.Context, status model.OrderStatus, page, size int) ([]*model.Order, int64, error)

	// FindAll 分页查询所有订单
	FindAll(ctx context.Context, page, size int) ([]*model.Order, int64, error)

	// Search 搜索订单
	Search(ctx context.Context, keyword string, page, size int) ([]*model.Order, int64, error)

	// Update 更新订单
	Update(ctx context.Context, order *model.Order) error

	// Delete 删除订单（软删除）
	Delete(ctx context.Context, id string) error

	// CountByStatus 统计各状态订单数量
	CountByStatus(ctx context.Context) (map[model.OrderStatus]int64, error)

	// FindByDateRange 根据日期范围查询订单
	FindByDateRange(ctx context.Context, startDate, endDate string, page, size int) ([]*model.Order, int64, error)

	// FindByProductID 根据商品ID查询包含该商品的订单
	FindByProductID(ctx context.Context, productID string, page, size int) ([]*model.Order, int64, error)
}

package repository

import (
	"context"

	"mall-go/services/order-service/domain/model"
)

// OrderQuery 订单查询参数
type OrderQuery struct {
	OrderSn         string
	Status          string
	MemberUsername  string
	ReceiverName    string
	ReceiverPhone   string
	CreateTimeBegin string
	CreateTimeEnd   string
	SourceType      int
	OrderType       int
	Page            int
	Size            int
}

// OrderRepository 订单仓储接口
type OrderRepository interface {
	// 保存订单
	Save(ctx context.Context, order *model.Order) error

	// 根据ID查找订单
	FindByID(ctx context.Context, id string) (*model.Order, error)

	// 根据订单号查找订单
	FindByOrderSn(ctx context.Context, orderSn string) (*model.Order, error)

	// 查找会员的订单列表
	FindByMemberID(ctx context.Context, memberID string, page, size int) ([]*model.Order, int64, error)

	// 更新订单
	Update(ctx context.Context, order *model.Order) error

	// 删除订单(逻辑删除)
	Delete(ctx context.Context, id string) error

	// 分页查询订单列表
	List(ctx context.Context, query OrderQuery) ([]*model.Order, int64, error)

	// 更新订单状态
	UpdateStatus(ctx context.Context, id string, status model.OrderStatus) error

	// 更新订单备注
	UpdateNote(ctx context.Context, id string, note string) error

	// 更新收货人信息
	UpdateReceiverInfo(ctx context.Context, id string, receiverInfo map[string]string) error
}

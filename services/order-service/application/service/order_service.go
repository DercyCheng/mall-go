package service

import (
	"context"

	"mall-go/services/order-service/application/dto"
)

// OrderService 订单应用服务接口
type OrderService interface {
	// CreateOrder 创建订单
	CreateOrder(ctx context.Context, req dto.OrderCreateRequest) (dto.OrderResponse, error)

	// PayOrder 支付订单
	PayOrder(ctx context.Context, req dto.OrderPayRequest) error

	// ShipOrder 发货订单
	ShipOrder(ctx context.Context, req dto.OrderShipRequest) error

	// UpdateOrderStatus 更新订单状态
	UpdateOrderStatus(ctx context.Context, req dto.OrderStatusUpdateRequest) error

	// CancelOrder 取消订单
	CancelOrder(ctx context.Context, orderID string, reason string) error

	// GetOrder 获取订单详情
	GetOrder(ctx context.Context, orderID string) (dto.OrderResponse, error)

	// GetOrderByOrderSN 根据订单编号获取订单
	GetOrderByOrderSN(ctx context.Context, orderSN string) (dto.OrderResponse, error)

	// GetUserOrders 获取用户订单列表
	GetUserOrders(ctx context.Context, userID string, page, size int) (dto.OrderListResponse, error)

	// SearchOrders 搜索订单
	SearchOrders(ctx context.Context, req dto.OrderQueryRequest) (dto.OrderListResponse, error)

	// GetOrderStatistics 获取订单统计数据
	GetOrderStatistics(ctx context.Context) (dto.OrderStatisticsResponse, error)

	// DeleteOrder 删除订单
	DeleteOrder(ctx context.Context, orderID string) error
}

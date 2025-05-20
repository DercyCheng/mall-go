package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"mall-go/services/order-service/domain/model"
	"mall-go/services/order-service/domain/repository"
)

// OrderDomainService 订单领域服务，处理订单领域的核心业务逻辑
type OrderDomainService struct {
	orderRepo repository.OrderRepository
	// 这里可以注入其他依赖的仓储或服务
}

// NewOrderDomainService 创建订单领域服务的工厂方法
func NewOrderDomainService(orderRepo repository.OrderRepository) *OrderDomainService {
	return &OrderDomainService{
		orderRepo: orderRepo,
	}
}

// CreateOrder 创建新订单
func (s *OrderDomainService) CreateOrder(ctx context.Context, order *model.Order) error {
	if order == nil {
		return errors.New("order cannot be nil")
	}

	// 生成订单编号
	if order.OrderSN == "" {
		order.OrderSN = s.generateOrderSN()
	}

	// 设置订单初始状态
	order.Status = model.OrderStatusPending
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	// 保存订单
	return s.orderRepo.Save(ctx, order)
}

// PayOrder 支付订单
func (s *OrderDomainService) PayOrder(ctx context.Context, orderID string, paymentMethod model.PaymentMethod, paymentTime time.Time) error {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return err
	}

	if order.Status != model.OrderStatusPending {
		return errors.New("order status is not pending, cannot pay")
	}

	// 更新订单支付信息
	order.Status = model.OrderStatusPaid
	order.PaymentMethod = paymentMethod
	order.PaymentTime = paymentTime
	order.UpdatedAt = time.Now()

	// 保存更新后的订单
	return s.orderRepo.Update(ctx, order)
}

// ShipOrder 订单发货
func (s *OrderDomainService) ShipOrder(ctx context.Context, orderID, deliveryCompany, trackingNo string) error {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return err
	}

	if order.Status != model.OrderStatusPaid {
		return errors.New("order is not paid, cannot ship")
	}

	// 更新订单发货信息
	order.Status = model.OrderStatusShipping
	order.DeliveryCompany = deliveryCompany
	order.DeliveryTrackingNo = trackingNo
	order.DeliveryTime = time.Now()
	order.UpdatedAt = time.Now()

	// 保存更新后的订单
	return s.orderRepo.Update(ctx, order)
}

// CompleteOrder 完成订单
func (s *OrderDomainService) CompleteOrder(ctx context.Context, orderID string) error {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return err
	}

	if order.Status != model.OrderStatusDelivered {
		return errors.New("order is not delivered, cannot complete")
	}

	// 更新订单状态为已完成
	order.Status = model.OrderStatusCompleted
	order.UpdatedAt = time.Now()

	// 保存更新后的订单
	return s.orderRepo.Update(ctx, order)
}

// CancelOrder 取消订单
func (s *OrderDomainService) CancelOrder(ctx context.Context, orderID string, reason string) error {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return err
	}

	// 只有待付款的订单可以取消
	if order.Status != model.OrderStatusPending {
		return errors.New("only pending orders can be cancelled")
	}

	// 更新订单状态为已取消
	order.Status = model.OrderStatusCancelled
	order.Note = order.Note + "\nCancel reason: " + reason
	order.UpdatedAt = time.Now()

	// 保存更新后的订单
	return s.orderRepo.Update(ctx, order)
}

// ReceiveOrder 确认收货
func (s *OrderDomainService) ReceiveOrder(ctx context.Context, orderID string) error {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return err
	}

	if order.Status != model.OrderStatusShipping {
		return errors.New("order is not shipping, cannot confirm receipt")
	}

	// 更新订单状态为已送达
	order.Status = model.OrderStatusDelivered
	order.ReceiveTime = time.Now()
	order.UpdatedAt = time.Now()

	// 保存更新后的订单
	return s.orderRepo.Update(ctx, order)
}

// ApplyRefund 申请退款
func (s *OrderDomainService) ApplyRefund(ctx context.Context, orderID string, reason string) error {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return err
	}

	// 只有已支付的订单可以申请退款
	if order.Status != model.OrderStatusPaid && order.Status != model.OrderStatusShipping {
		return errors.New("only paid or shipping orders can apply for refund")
	}

	// 更新订单状态为退款中
	order.Status = model.OrderStatusRefunding
	order.Note = order.Note + "\nRefund reason: " + reason
	order.UpdatedAt = time.Now()

	// 保存更新后的订单
	return s.orderRepo.Update(ctx, order)
}

// RefundOrder 确认退款
func (s *OrderDomainService) RefundOrder(ctx context.Context, orderID string) error {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return err
	}

	if order.Status != model.OrderStatusRefunding {
		return errors.New("order is not in refunding status")
	}

	// 更新订单状态为已退款
	order.Status = model.OrderStatusRefunded
	order.UpdatedAt = time.Now()

	// 保存更新后的订单
	return s.orderRepo.Update(ctx, order)
}

// 生成订单编号
func (s *OrderDomainService) generateOrderSN() string {
	// 生成格式: 年月日时分秒+4位随机数
	timestamp := time.Now().Format("20060102150405")
	randomID := uuid.New().String()[0:4]
	return timestamp + randomID
}

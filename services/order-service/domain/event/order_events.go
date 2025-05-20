package event

import (
	"time"

	"mall-go/services/order-service/domain/model"
)

// OrderEvent 订单事件基础结构
type OrderEvent struct {
	ID          string    `json:"id"`
	OrderID     string    `json:"orderId"`
	OrderSN     string    `json:"orderSN"`
	UserID      string    `json:"userId"`
	EventType   string    `json:"eventType"`
	OccurredAt  time.Time `json:"occurredAt"`
	PublishedAt time.Time `json:"publishedAt"`
}

// OrderCreatedEvent 订单创建事件
type OrderCreatedEvent struct {
	OrderEvent
	OrderItems  []OrderItemEvent `json:"orderItems"`
	TotalAmount float64          `json:"totalAmount"`
}

// OrderPaidEvent 订单支付事件
type OrderPaidEvent struct {
	OrderEvent
	PaymentMethod string    `json:"paymentMethod"`
	PayAmount     float64   `json:"payAmount"`
	PaymentTime   time.Time `json:"paymentTime"`
}

// OrderShippedEvent 订单发货事件
type OrderShippedEvent struct {
	OrderEvent
	DeliveryCompany    string    `json:"deliveryCompany"`
	DeliveryTrackingNo string    `json:"deliveryTrackingNo"`
	DeliveryTime       time.Time `json:"deliveryTime"`
}

// OrderDeliveredEvent 订单收货事件
type OrderDeliveredEvent struct {
	OrderEvent
	ReceiveTime time.Time `json:"receiveTime"`
}

// OrderCompletedEvent 订单完成事件
type OrderCompletedEvent struct {
	OrderEvent
}

// OrderCancelledEvent 订单取消事件
type OrderCancelledEvent struct {
	OrderEvent
	Reason string `json:"reason"`
}

// OrderRefundingEvent 订单申请退款事件
type OrderRefundingEvent struct {
	OrderEvent
	Reason string `json:"reason"`
}

// OrderRefundedEvent 订单退款完成事件
type OrderRefundedEvent struct {
	OrderEvent
}

// OrderItemEvent 订单项事件数据
type OrderItemEvent struct {
	ProductID  string  `json:"productId"`
	ProductSKU string  `json:"productSku"`
	Quantity   int     `json:"quantity"`
	UnitPrice  float64 `json:"unitPrice"`
	TotalPrice float64 `json:"totalPrice"`
}

// NewOrderCreatedEvent 创建订单创建事件
func NewOrderCreatedEvent(order *model.Order) *OrderCreatedEvent {
	items := make([]OrderItemEvent, 0, len(order.OrderItems))
	for _, item := range order.OrderItems {
		items = append(items, OrderItemEvent{
			ProductID:  item.ProductID,
			ProductSKU: item.ProductSKU,
			Quantity:   item.Quantity,
			UnitPrice:  item.UnitPrice.Amount,
			TotalPrice: item.TotalPrice.Amount,
		})
	}

	return &OrderCreatedEvent{
		OrderEvent: OrderEvent{
			ID:          order.ID,
			OrderID:     order.ID,
			OrderSN:     order.OrderSN,
			UserID:      order.UserID,
			EventType:   "order.created",
			OccurredAt:  order.CreatedAt,
			PublishedAt: time.Now(),
		},
		OrderItems:  items,
		TotalAmount: order.TotalAmount.Amount,
	}
}

// NewOrderPaidEvent 创建订单支付事件
func NewOrderPaidEvent(order *model.Order) *OrderPaidEvent {
	return &OrderPaidEvent{
		OrderEvent: OrderEvent{
			ID:          order.ID,
			OrderID:     order.ID,
			OrderSN:     order.OrderSN,
			UserID:      order.UserID,
			EventType:   "order.paid",
			OccurredAt:  order.PaymentTime,
			PublishedAt: time.Now(),
		},
		PaymentMethod: string(order.PaymentMethod),
		PayAmount:     order.PayAmount.Amount,
		PaymentTime:   order.PaymentTime,
	}
}

// NewOrderShippedEvent 创建订单发货事件
func NewOrderShippedEvent(order *model.Order) *OrderShippedEvent {
	return &OrderShippedEvent{
		OrderEvent: OrderEvent{
			ID:          order.ID,
			OrderID:     order.ID,
			OrderSN:     order.OrderSN,
			UserID:      order.UserID,
			EventType:   "order.shipped",
			OccurredAt:  order.DeliveryTime,
			PublishedAt: time.Now(),
		},
		DeliveryCompany:    order.DeliveryCompany,
		DeliveryTrackingNo: order.DeliveryTrackingNo,
		DeliveryTime:       order.DeliveryTime,
	}
}

// NewOrderDeliveredEvent 创建订单收货事件
func NewOrderDeliveredEvent(order *model.Order) *OrderDeliveredEvent {
	return &OrderDeliveredEvent{
		OrderEvent: OrderEvent{
			ID:          order.ID,
			OrderID:     order.ID,
			OrderSN:     order.OrderSN,
			UserID:      order.UserID,
			EventType:   "order.delivered",
			OccurredAt:  order.ReceiveTime,
			PublishedAt: time.Now(),
		},
		ReceiveTime: order.ReceiveTime,
	}
}

// NewOrderCompletedEvent 创建订单完成事件
func NewOrderCompletedEvent(order *model.Order) *OrderCompletedEvent {
	return &OrderCompletedEvent{
		OrderEvent: OrderEvent{
			ID:          order.ID,
			OrderID:     order.ID,
			OrderSN:     order.OrderSN,
			UserID:      order.UserID,
			EventType:   "order.completed",
			OccurredAt:  order.UpdatedAt,
			PublishedAt: time.Now(),
		},
	}
}

// NewOrderCancelledEvent 创建订单取消事件
func NewOrderCancelledEvent(order *model.Order, reason string) *OrderCancelledEvent {
	return &OrderCancelledEvent{
		OrderEvent: OrderEvent{
			ID:          order.ID,
			OrderID:     order.ID,
			OrderSN:     order.OrderSN,
			UserID:      order.UserID,
			EventType:   "order.cancelled",
			OccurredAt:  order.UpdatedAt,
			PublishedAt: time.Now(),
		},
		Reason: reason,
	}
}

// NewOrderRefundingEvent 创建订单申请退款事件
func NewOrderRefundingEvent(order *model.Order, reason string) *OrderRefundingEvent {
	return &OrderRefundingEvent{
		OrderEvent: OrderEvent{
			ID:          order.ID,
			OrderID:     order.ID,
			OrderSN:     order.OrderSN,
			UserID:      order.UserID,
			EventType:   "order.refunding",
			OccurredAt:  order.UpdatedAt,
			PublishedAt: time.Now(),
		},
		Reason: reason,
	}
}

// NewOrderRefundedEvent 创建订单退款完成事件
func NewOrderRefundedEvent(order *model.Order) *OrderRefundedEvent {
	return &OrderRefundedEvent{
		OrderEvent: OrderEvent{
			ID:          order.ID,
			OrderID:     order.ID,
			OrderSN:     order.OrderSN,
			UserID:      order.UserID,
			EventType:   "order.refunded",
			OccurredAt:  order.UpdatedAt,
			PublishedAt: time.Now(),
		},
	}
}

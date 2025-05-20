package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// OrderStatus 订单状态枚举
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"   // 待付款
	OrderStatusPaid      OrderStatus = "paid"      // 已付款
	OrderStatusShipping  OrderStatus = "shipping"  // 配送中
	OrderStatusDelivered OrderStatus = "delivered" // 已送达
	OrderStatusCompleted OrderStatus = "completed" // 已完成
	OrderStatusCancelled OrderStatus = "cancelled" // 已取消
	OrderStatusRefunding OrderStatus = "refunding" // 退款中
	OrderStatusRefunded  OrderStatus = "refunded"  // 已退款
)

// PaymentMethod 支付方式枚举
type PaymentMethod string

const (
	PaymentMethodAlipay PaymentMethod = "alipay" // 支付宝
	PaymentMethodWechat PaymentMethod = "wechat" // 微信支付
	PaymentMethodCredit PaymentMethod = "credit" // 信用卡
	PaymentMethodDebit  PaymentMethod = "debit"  // 借记卡
	PaymentMethodCash   PaymentMethod = "cash"   // 货到付款
)

// DeliveryMethod 配送方式枚举
type DeliveryMethod string

const (
	DeliveryMethodExpress DeliveryMethod = "express" // 快递
	DeliveryMethodPickup  DeliveryMethod = "pickup"  // 自提
)

// Money 金额值对象
type Money struct {
	Amount   float64 // 金额
	Currency string  // 货币类型，默认CNY
}

// Address 地址值对象
type Address struct {
	Province      string // 省
	City          string // 市
	District      string // 区
	DetailAddress string // 详细地址
	ReceiverName  string // 收件人姓名
	ReceiverPhone string // 收件人电话
	PostCode      string // 邮编
}

// OrderItem 订单项值对象
type OrderItem struct {
	ID              string // 订单项ID
	ProductID       string // 商品ID
	ProductName     string // 商品名称
	ProductImage    string // 商品图片
	ProductSKU      string // 商品SKU
	Quantity        int    // 数量
	UnitPrice       Money  // 单价
	TotalPrice      Money  // 总价
	Discount        Money  // 折扣
	AttributeValues string // 商品属性
}

// Order 订单聚合根
type Order struct {
	ID                 string         // 订单ID
	UserID             string         // 用户ID
	OrderSN            string         // 订单编号
	Status             OrderStatus    // 订单状态
	TotalAmount        Money          // 订单总金额
	PayAmount          Money          // 实际支付金额
	FreightAmount      Money          // 运费
	DiscountAmount     Money          // 折扣金额
	CouponAmount       Money          // 优惠券金额
	PointAmount        Money          // 积分抵扣金额
	PaymentMethod      PaymentMethod  // 支付方式
	PaymentTime        time.Time      // 支付时间
	DeliveryMethod     DeliveryMethod // 配送方式
	DeliveryCompany    string         // 物流公司
	DeliveryTrackingNo string         // 物流单号
	DeliveryTime       time.Time      // 发货时间
	ReceiveTime        time.Time      // 确认收货时间
	CommentTime        time.Time      // 评价时间
	BillingAddress     Address        // 账单地址
	ShippingAddress    Address        // 配送地址
	Note               string         // 订单备注
	OrderItems         []OrderItem    // 订单项列表
	CreatedAt          time.Time      // 创建时间
	UpdatedAt          time.Time      // 更新时间
}

// NewOrder 创建新订单的工厂方法
func NewOrder(userID string, items []OrderItem, shippingAddress Address, totalAmount, freightAmount Money) (*Order, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	if len(items) == 0 {
		return nil, errors.New("order must have at least one item")
	}

	// 计算订单总金额
	totalPrice := totalAmount.Amount

	// 生成订单编号 (简化版，实际可能需要更复杂的生成方式)
	orderSN := generateOrderSN()

	order := &Order{
		ID:              uuid.New().String(),
		UserID:          userID,
		OrderSN:         orderSN,
		Status:          OrderStatusPending,
		TotalAmount:     totalAmount,
		FreightAmount:   freightAmount,
		PayAmount:       Money{Amount: totalPrice + freightAmount.Amount, Currency: "CNY"},
		OrderItems:      items,
		ShippingAddress: shippingAddress,
		BillingAddress:  shippingAddress, // 默认使用相同地址
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	return order, nil
}

// Pay 支付订单
func (o *Order) Pay(paymentMethod PaymentMethod) error {
	if o.Status != OrderStatusPending {
		return errors.New("only pending orders can be paid")
	}

	o.Status = OrderStatusPaid
	o.PaymentMethod = paymentMethod
	o.PaymentTime = time.Now()
	o.UpdatedAt = time.Now()
	return nil
}

// Ship 发货
func (o *Order) Ship(deliveryCompany, trackingNo string) error {
	if o.Status != OrderStatusPaid {
		return errors.New("only paid orders can be shipped")
	}

	o.Status = OrderStatusShipping
	o.DeliveryCompany = deliveryCompany
	o.DeliveryTrackingNo = trackingNo
	o.DeliveryTime = time.Now()
	o.UpdatedAt = time.Now()
	return nil
}

// Receive 确认收货
func (o *Order) Receive() error {
	if o.Status != OrderStatusShipping {
		return errors.New("only shipping orders can be received")
	}

	o.Status = OrderStatusDelivered
	o.ReceiveTime = time.Now()
	o.UpdatedAt = time.Now()
	return nil
}

// Complete 完成订单
func (o *Order) Complete() error {
	if o.Status != OrderStatusDelivered {
		return errors.New("only delivered orders can be completed")
	}

	o.Status = OrderStatusCompleted
	o.UpdatedAt = time.Now()
	return nil
}

// Cancel 取消订单
func (o *Order) Cancel() error {
	if o.Status != OrderStatusPending {
		return errors.New("only pending orders can be cancelled")
	}

	o.Status = OrderStatusCancelled
	o.UpdatedAt = time.Now()
	return nil
}

// RequestRefund 申请退款
func (o *Order) RequestRefund() error {
	if o.Status != OrderStatusPaid && o.Status != OrderStatusShipping && o.Status != OrderStatusDelivered {
		return errors.New("order cannot be refunded in current status")
	}

	o.Status = OrderStatusRefunding
	o.UpdatedAt = time.Now()
	return nil
}

// Refund 确认退款
func (o *Order) Refund() error {
	if o.Status != OrderStatusRefunding {
		return errors.New("only refunding orders can be refunded")
	}

	o.Status = OrderStatusRefunded
	o.UpdatedAt = time.Now()
	return nil
}

// ApplyCoupon 应用优惠券
func (o *Order) ApplyCoupon(couponAmount Money) error {
	if o.Status != OrderStatusPending {
		return errors.New("coupon can only be applied to pending orders")
	}

	if couponAmount.Amount > o.TotalAmount.Amount {
		return errors.New("coupon amount exceeds order total")
	}

	o.CouponAmount = couponAmount
	o.PayAmount.Amount = o.PayAmount.Amount - couponAmount.Amount
	o.UpdatedAt = time.Now()
	return nil
}

// UsePoints 使用积分抵扣
func (o *Order) UsePoints(pointAmount Money) error {
	if o.Status != OrderStatusPending {
		return errors.New("points can only be applied to pending orders")
	}

	if pointAmount.Amount > o.TotalAmount.Amount {
		return errors.New("point amount exceeds order total")
	}

	o.PointAmount = pointAmount
	o.PayAmount.Amount = o.PayAmount.Amount - pointAmount.Amount
	o.UpdatedAt = time.Now()
	return nil
}

// AddOrderItem 添加订单项
func (o *Order) AddOrderItem(item OrderItem) error {
	if o.Status != OrderStatusPending {
		return errors.New("items can only be added to pending orders")
	}

	o.OrderItems = append(o.OrderItems, item)
	o.TotalAmount.Amount += item.TotalPrice.Amount
	o.PayAmount.Amount += item.TotalPrice.Amount
	o.UpdatedAt = time.Now()
	return nil
}

// RemoveOrderItem 移除订单项
func (o *Order) RemoveOrderItem(itemID string) error {
	if o.Status != OrderStatusPending {
		return errors.New("items can only be removed from pending orders")
	}

	for i, item := range o.OrderItems {
		if item.ID == itemID {
			o.TotalAmount.Amount -= item.TotalPrice.Amount
			o.PayAmount.Amount -= item.TotalPrice.Amount

			// 移除订单项
			o.OrderItems = append(o.OrderItems[:i], o.OrderItems[i+1:]...)
			o.UpdatedAt = time.Now()
			return nil
		}
	}

	return errors.New("order item not found")
}

// UpdateShippingAddress 更新配送地址
func (o *Order) UpdateShippingAddress(address Address) error {
	if o.Status != OrderStatusPending {
		return errors.New("shipping address can only be updated for pending orders")
	}

	o.ShippingAddress = address
	o.UpdatedAt = time.Now()
	return nil
}

// UpdateBillingAddress 更新账单地址
func (o *Order) UpdateBillingAddress(address Address) error {
	if o.Status != OrderStatusPending {
		return errors.New("billing address can only be updated for pending orders")
	}

	o.BillingAddress = address
	o.UpdatedAt = time.Now()
	return nil
}

// UpdateNote 更新订单备注
func (o *Order) UpdateNote(note string) error {
	o.Note = note
	o.UpdatedAt = time.Now()
	return nil
}

// CalculateTotal 重新计算订单金额
func (o *Order) CalculateTotal() {
	var total float64
	for _, item := range o.OrderItems {
		total += item.TotalPrice.Amount
	}
	o.TotalAmount.Amount = total

	// 计算实付金额 = 总金额 + 运费 - 折扣 - 优惠券 - 积分
	o.PayAmount.Amount = total + o.FreightAmount.Amount - o.DiscountAmount.Amount - o.CouponAmount.Amount - o.PointAmount.Amount
	o.UpdatedAt = time.Now()
}

// generateOrderSN 生成订单编号
func generateOrderSN() string {
	// 简单实现，实际可能需要包含更多信息，如日期、随机数等
	timeStr := time.Now().Format("20060102150405")
	randomStr := uuid.New().String()[0:8]
	return timeStr + randomStr
}

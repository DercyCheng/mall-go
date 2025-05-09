package model

import (
	"errors"
	"fmt"
	"time"

	"strconv"

	"github.com/google/uuid"
)

// OrderStatus 订单状态
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"   // 待付款
	OrderStatusPaid      OrderStatus = "paid"      // 已付款
	OrderStatusShipping  OrderStatus = "shipping"  // 配送中
	OrderStatusCompleted OrderStatus = "completed" // 已完成
	OrderStatusCancelled OrderStatus = "cancelled" // 已取消
	OrderStatusClosed    OrderStatus = "closed"    // 已关闭
	OrderStatusRefunding OrderStatus = "refunding" // 退款中
	OrderStatusRefunded  OrderStatus = "refunded"  // 已退款
)

// PaymentType 支付方式
type PaymentType string

const (
	PaymentTypeAlipay   PaymentType = "alipay"   // 支付宝
	PaymentTypeWechat   PaymentType = "wechat"   // 微信
	PaymentTypeUnionPay PaymentType = "unionpay" // 银联
	PaymentTypeCash     PaymentType = "cash"     // 现金
	PaymentTypeCredit   PaymentType = "credit"   // 信用卡
)

// Money 金额值对象
type Money struct {
	Amount   float64
	Currency string
}

// Address 收货地址值对象
type Address struct {
	Province      string
	City          string
	District      string
	DetailAddress string
	PostCode      string
	Name          string
	Phone         string
	IsDefault     bool
}

// OrderItem 订单项
type OrderItem struct {
	ID                string
	ProductID         string
	ProductSn         string
	ProductName       string
	ProductPic        string
	ProductPrice      Money
	ProductQuantity   int
	ProductAttr       string
	CouponAmount      float64
	PromotionAmount   float64
	RealAmount        float64
	GiftIntegration   int
	GiftGrowth        int
	ProductCategoryId string
	CreatedAt         time.Time
}

// Order 订单聚合根
type Order struct {
	ID                    string
	MemberID              string
	OrderSn               string
	MemberUsername        string
	TotalAmount           Money
	PayAmount             Money
	FreightAmount         Money
	PromotionAmount       Money
	IntegrationAmount     Money
	CouponAmount          Money
	DiscountAmount        Money
	PayType               PaymentType
	SourceType            int // 0->PC订单；1->app订单
	Status                OrderStatus
	OrderType             int // 0->正常订单；1->秒杀订单
	DeliveryCompany       string
	DeliverySn            string
	AutoConfirmDay        int
	Integration           int
	Growth                int
	PromotionInfo         string
	BillType              int // 0->不开发票；1->电子发票；2->纸质发票
	BillHeader            string
	BillContent           string
	BillReceiverPhone     string
	BillReceiverEmail     string
	ReceiverName          string
	ReceiverPhone         string
	ReceiverPostCode      string
	ReceiverProvince      string
	ReceiverCity          string
	ReceiverDistrict      string
	ReceiverDetailAddress string
	Note                  string
	ConfirmStatus         int // 0->未确认；1->已确认
	DeleteStatus          int // 0->未删除；1->已删除
	UseIntegration        int
	PaymentTime           time.Time
	DeliveryTime          time.Time
	ReceiveTime           time.Time
	CommentTime           time.Time
	ModifyTime            time.Time
	OrderItems            []OrderItem
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

// NewOrder 创建新订单
func NewOrder(memberID string, memberUsername string, address *Address, orderItems []OrderItem, totalAmount, freightAmount float64) (*Order, error) {
	if memberID == "" {
		return nil, errors.New("会员ID不能为空")
	}
	if len(orderItems) == 0 {
		return nil, errors.New("订单项不能为空")
	}

	// 生成订单ID和订单编号
	id := uuid.New().String()
	orderSn := generateOrderSn(memberID)

	now := time.Now()

	// 计算各种金额
	var promotionAmount, integrationAmount, couponAmount float64 = 0, 0, 0
	var realAmount float64 = 0

	// 计算订单项总金额
	for _, item := range orderItems {
		realAmount += item.RealAmount
	}

	// 计算应付金额 = 总金额 + 运费 - 促销折扣 - 积分抵扣 - 优惠券抵扣
	payAmount := totalAmount + freightAmount - promotionAmount - integrationAmount - couponAmount

	// 创建订单对象
	order := &Order{
		ID:                    id,
		MemberID:              memberID,
		OrderSn:               orderSn,
		MemberUsername:        memberUsername,
		TotalAmount:           Money{Amount: totalAmount, Currency: "CNY"},
		PayAmount:             Money{Amount: payAmount, Currency: "CNY"},
		FreightAmount:         Money{Amount: freightAmount, Currency: "CNY"},
		PromotionAmount:       Money{Amount: promotionAmount, Currency: "CNY"},
		IntegrationAmount:     Money{Amount: integrationAmount, Currency: "CNY"},
		CouponAmount:          Money{Amount: couponAmount, Currency: "CNY"},
		DiscountAmount:        Money{Amount: 0, Currency: "CNY"},
		PayType:               PaymentTypeAlipay, // 默认支付宝
		SourceType:            0,                 // 默认PC订单
		Status:                OrderStatusPending,
		OrderType:             0, // 默认正常订单
		ConfirmStatus:         0,
		DeleteStatus:          0,
		ReceiverName:          address.Name,
		ReceiverPhone:         address.Phone,
		ReceiverPostCode:      address.PostCode,
		ReceiverProvince:      address.Province,
		ReceiverCity:          address.City,
		ReceiverDistrict:      address.District,
		ReceiverDetailAddress: address.DetailAddress,
		OrderItems:            orderItems,
		CreatedAt:             now,
		UpdatedAt:             now,
	}

	return order, nil
}

// Pay 支付订单
func (o *Order) Pay(payType PaymentType) error {
	if o.Status != OrderStatusPending {
		return fmt.Errorf("订单状态不正确，当前状态: %s", o.Status)
	}

	o.Status = OrderStatusPaid
	o.PayType = payType
	o.PaymentTime = time.Now()
	o.UpdatedAt = time.Now()

	return nil
}

// Ship 发货
func (o *Order) Ship(deliveryCompany, deliverySn string) error {
	if o.Status != OrderStatusPaid {
		return fmt.Errorf("订单状态不正确，当前状态: %s", o.Status)
	}

	o.Status = OrderStatusShipping
	o.DeliveryCompany = deliveryCompany
	o.DeliverySn = deliverySn
	o.DeliveryTime = time.Now()
	o.UpdatedAt = time.Now()

	return nil
}

// Receive 确认收货
func (o *Order) Receive() error {
	if o.Status != OrderStatusShipping {
		return fmt.Errorf("订单状态不正确，当前状态: %s", o.Status)
	}

	o.Status = OrderStatusCompleted
	o.ReceiveTime = time.Now()
	o.UpdatedAt = time.Now()

	return nil
}

// Cancel 取消订单
func (o *Order) Cancel(reason string) error {
	if o.Status != OrderStatusPending {
		return fmt.Errorf("订单状态不正确，当前状态: %s", o.Status)
	}

	o.Status = OrderStatusCancelled
	o.Note = reason
	o.UpdatedAt = time.Now()

	return nil
}

// Close 关闭订单
func (o *Order) Close() error {
	if o.Status != OrderStatusPending {
		return fmt.Errorf("订单状态不正确，当前状态: %s", o.Status)
	}

	o.Status = OrderStatusClosed
	o.UpdatedAt = time.Now()

	return nil
}

// ApplyRefund 申请退款
func (o *Order) ApplyRefund(reason string) error {
	if o.Status != OrderStatusPaid && o.Status != OrderStatusShipping {
		return fmt.Errorf("订单状态不正确，当前状态: %s", o.Status)
	}

	o.Status = OrderStatusRefunding
	o.Note = reason
	o.UpdatedAt = time.Now()

	return nil
}

// Refund 确认退款
func (o *Order) Refund() error {
	if o.Status != OrderStatusRefunding {
		return fmt.Errorf("订单状态不正确，当前状态: %s", o.Status)
	}

	o.Status = OrderStatusRefunded
	o.UpdatedAt = time.Now()

	return nil
}

// ConfirmOrder 确认订单
func (o *Order) ConfirmOrder() error {
	if o.Status != OrderStatusPending {
		return fmt.Errorf("订单状态不正确，当前状态: %s", o.Status)
	}

	o.ConfirmStatus = 1
	o.UpdatedAt = time.Now()

	return nil
}

// Comment 评论订单
func (o *Order) Comment() error {
	if o.Status != OrderStatusCompleted {
		return fmt.Errorf("订单状态不正确，当前状态: %s", o.Status)
	}

	o.CommentTime = time.Now()
	o.UpdatedAt = time.Now()

	return nil
}

// UpdateNote 更新订单备注
func (o *Order) UpdateNote(note string) {
	o.Note = note
	o.UpdatedAt = time.Now()
}

// generateOrderSn 生成订单号
func generateOrderSn(memberID string) string {
	now := time.Now()
	// 订单号格式: 年月日时分秒 + 会员ID后四位 + 随机4位数
	timestamp := now.Format("20060102150405")

	// 获取会员ID的后四位
	idSuffix := memberID
	if len(memberID) > 4 {
		idSuffix = memberID[len(memberID)-4:]
	}

	// 生成随机4位数
	randomNum := strconv.Itoa(int(now.UnixNano() % 10000))
	for len(randomNum) < 4 {
		randomNum = "0" + randomNum
	}

	return timestamp + idSuffix + randomNum
}

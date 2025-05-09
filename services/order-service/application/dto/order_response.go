package dto

import (
	"time"
)

// OrderItemResponse 订单项响应DTO
type OrderItemResponse struct {
	ID                string    `json:"id"`
	ProductID         string    `json:"productId"`
	ProductSn         string    `json:"productSn"`
	ProductName       string    `json:"productName"`
	ProductPic        string    `json:"productPic"`
	ProductPrice      float64   `json:"productPrice"`
	ProductQuantity   int       `json:"productQuantity"`
	ProductAttr       string    `json:"productAttr"`
	CouponAmount      float64   `json:"couponAmount"`
	PromotionAmount   float64   `json:"promotionAmount"`
	RealAmount        float64   `json:"realAmount"`
	GiftIntegration   int       `json:"giftIntegration"`
	GiftGrowth        int       `json:"giftGrowth"`
	ProductCategoryId string    `json:"productCategoryId"`
	CreatedAt         time.Time `json:"createdAt"`
}

// OrderResponse 订单响应DTO
type OrderResponse struct {
	ID                    string              `json:"id"`
	MemberID              string              `json:"memberId"`
	OrderSn               string              `json:"orderSn"`
	MemberUsername        string              `json:"memberUsername"`
	TotalAmount           float64             `json:"totalAmount"`
	PayAmount             float64             `json:"payAmount"`
	FreightAmount         float64             `json:"freightAmount"`
	PromotionAmount       float64             `json:"promotionAmount"`
	IntegrationAmount     float64             `json:"integrationAmount"`
	CouponAmount          float64             `json:"couponAmount"`
	DiscountAmount        float64             `json:"discountAmount"`
	PayType               string              `json:"payType"`
	SourceType            int                 `json:"sourceType"`
	Status                string              `json:"status"`
	OrderType             int                 `json:"orderType"`
	DeliveryCompany       string              `json:"deliveryCompany"`
	DeliverySn            string              `json:"deliverySn"`
	AutoConfirmDay        int                 `json:"autoConfirmDay"`
	Integration           int                 `json:"integration"`
	Growth                int                 `json:"growth"`
	PromotionInfo         string              `json:"promotionInfo"`
	BillType              int                 `json:"billType"`
	BillHeader            string              `json:"billHeader"`
	BillContent           string              `json:"billContent"`
	BillReceiverPhone     string              `json:"billReceiverPhone"`
	BillReceiverEmail     string              `json:"billReceiverEmail"`
	ReceiverName          string              `json:"receiverName"`
	ReceiverPhone         string              `json:"receiverPhone"`
	ReceiverPostCode      string              `json:"receiverPostCode"`
	ReceiverProvince      string              `json:"receiverProvince"`
	ReceiverCity          string              `json:"receiverCity"`
	ReceiverDistrict      string              `json:"receiverDistrict"`
	ReceiverDetailAddress string              `json:"receiverDetailAddress"`
	Note                  string              `json:"note"`
	ConfirmStatus         int                 `json:"confirmStatus"`
	DeleteStatus          int                 `json:"deleteStatus"`
	UseIntegration        int                 `json:"useIntegration"`
	PaymentTime           *time.Time          `json:"paymentTime"`
	DeliveryTime          *time.Time          `json:"deliveryTime"`
	ReceiveTime           *time.Time          `json:"receiveTime"`
	CommentTime           *time.Time          `json:"commentTime"`
	ModifyTime            *time.Time          `json:"modifyTime"`
	OrderItems            []OrderItemResponse `json:"orderItems"`
	CreatedAt             time.Time           `json:"createdAt"`
	UpdatedAt             time.Time           `json:"updatedAt"`
	StatusName            string              `json:"statusName"`      // 根据Status映射的状态名称，方便前端展示
	PaymentTypeName       string              `json:"paymentTypeName"` // 支付方式名称
	SourceTypeName        string              `json:"sourceTypeName"`  // 订单来源名称
	OrderTypeName         string              `json:"orderTypeName"`   // 订单类型名称
}

// OrderBriefResponse 订单简要信息响应DTO
type OrderBriefResponse struct {
	ID             string    `json:"id"`
	OrderSn        string    `json:"orderSn"`
	MemberUsername string    `json:"memberUsername"`
	TotalAmount    float64   `json:"totalAmount"`
	PayAmount      float64   `json:"payAmount"`
	Status         string    `json:"status"`
	OrderType      int       `json:"orderType"`
	ReceiverName   string    `json:"receiverName"`
	ReceiverPhone  string    `json:"receiverPhone"`
	CreatedAt      time.Time `json:"createdAt"`
}

// OrderListResponse 订单列表响应DTO
type OrderListResponse struct {
	Total int64                `json:"total"`
	List  []OrderBriefResponse `json:"list"`
}

// CreateOrderResponse 创建订单响应DTO
type CreateOrderResponse struct {
	OrderID   string  `json:"orderId"`
	OrderSn   string  `json:"orderSn"`
	PayAmount float64 `json:"payAmount"`
}

// OrderStatsResponse 订单统计响应DTO
type OrderStatsResponse struct {
	TotalCount      int64   `json:"totalCount"`
	TodayCount      int64   `json:"todayCount"`
	YesterdayCount  int64   `json:"yesterdayCount"`
	TotalAmount     float64 `json:"totalAmount"`
	TodayAmount     float64 `json:"todayAmount"`
	YesterdayAmount float64 `json:"yesterdayAmount"`
}

// OrderStatusCount 订单状态统计
type OrderStatusCount struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

// OrderStatisticsResponse 订单统计响应DTO
type OrderStatisticsResponse struct {
	TodayOrderCount      int64   `json:"todayOrderCount"`      // 今日订单数
	TodayOrderAmount     float64 `json:"todayOrderAmount"`     // 今日订单金额
	YesterdayOrderCount  int64   `json:"yesterdayOrderCount"`  // 昨日订单数
	YesterdayOrderAmount float64 `json:"yesterdayOrderAmount"` // 昨日订单金额
	PendingOrderCount    int64   `json:"pendingOrderCount"`    // 待付款订单数
	PaidOrderCount       int64   `json:"paidOrderCount"`       // 已付款订单数
	ShippingOrderCount   int64   `json:"shippingOrderCount"`   // 配送中订单数
	CompletedOrderCount  int64   `json:"completedOrderCount"`  // 已完成订单数
	CancelledOrderCount  int64   `json:"cancelledOrderCount"`  // 已取消订单数
	RefundingOrderCount  int64   `json:"refundingOrderCount"`  // 退款中订单数
	RefundedOrderCount   int64   `json:"refundedOrderCount"`   // 已退款订单数
}

// OrderStatusMap 订单状态映射
var OrderStatusMap = map[string]string{
	"pending":   "待付款",
	"paid":      "已付款",
	"shipping":  "配送中",
	"completed": "已完成",
	"cancelled": "已取消",
	"closed":    "已关闭",
	"refunding": "退款中",
	"refunded":  "已退款",
}

// PaymentTypeMap 支付方式映射
var PaymentTypeMap = map[string]string{
	"alipay":   "支付宝",
	"wechat":   "微信",
	"unionpay": "银联",
	"cash":     "现金",
	"credit":   "信用卡",
}

// OrderSourceTypeMap 订单来源映射
var OrderSourceTypeMap = map[int]string{
	0: "PC订单",
	1: "APP订单",
}

// OrderTypeMap 订单类型映射
var OrderTypeMap = map[int]string{
	0: "正常订单",
	1: "秒杀订单",
}

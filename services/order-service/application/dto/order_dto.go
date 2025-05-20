package dto

import (
	"time"
)

// OrderItemRequest 订单项请求DTO
type OrderItemRequest struct {
	ProductID       string `json:"productId" binding:"required"`
	Quantity        int    `json:"quantity" binding:"required,min=1"`
	AttributeValues string `json:"attributeValues,omitempty"`
}

// AddressRequest 地址请求DTO
type AddressRequest struct {
	Province      string `json:"province" binding:"required"`
	City          string `json:"city" binding:"required"`
	District      string `json:"district" binding:"required"`
	DetailAddress string `json:"detailAddress" binding:"required"`
	ReceiverName  string `json:"receiverName" binding:"required"`
	ReceiverPhone string `json:"receiverPhone" binding:"required"`
	PostCode      string `json:"postCode,omitempty"`
}

// OrderCreateRequest 创建订单请求DTO
type OrderCreateRequest struct {
	UserID          string             `json:"userId" binding:"required"`
	OrderItems      []OrderItemRequest `json:"orderItems" binding:"required,min=1"`
	ShippingAddress AddressRequest     `json:"shippingAddress" binding:"required"`
	BillingAddress  *AddressRequest    `json:"billingAddress,omitempty"`
	DeliveryMethod  string             `json:"deliveryMethod" binding:"required,oneof=express pickup"`
	Note            string             `json:"note,omitempty"`
	CouponCode      string             `json:"couponCode,omitempty"`
	UsePoints       bool               `json:"usePoints,omitempty"`
}

// OrderPayRequest 订单支付请求DTO
type OrderPayRequest struct {
	OrderID       string `json:"orderId" binding:"required"`
	PaymentMethod string `json:"paymentMethod" binding:"required,oneof=alipay wechat credit debit cash"`
	TransactionID string `json:"transactionId,omitempty"`
	Amount        string `json:"amount" binding:"required"`
}

// OrderShipRequest 订单发货请求DTO
type OrderShipRequest struct {
	OrderID         string `json:"orderId" binding:"required"`
	DeliveryCompany string `json:"deliveryCompany" binding:"required"`
	TrackingNo      string `json:"trackingNo" binding:"required"`
}

// OrderStatusUpdateRequest 订单状态更新请求DTO
type OrderStatusUpdateRequest struct {
	OrderID string `json:"orderId" binding:"required"`
	Status  string `json:"status" binding:"required,oneof=paid shipping delivered completed cancelled refunding refunded"`
	Reason  string `json:"reason,omitempty"`
}

// OrderQueryRequest 订单查询请求DTO
type OrderQueryRequest struct {
	UserID    string `form:"userId,omitempty"`
	Status    string `form:"status,omitempty"`
	OrderSN   string `form:"orderSn,omitempty"`
	StartDate string `form:"startDate,omitempty"`
	EndDate   string `form:"endDate,omitempty"`
	ProductID string `form:"productId,omitempty"`
	Keyword   string `form:"keyword,omitempty"`
	Page      int    `form:"page,default=1" binding:"min=1"`
	Size      int    `form:"size,default=20" binding:"min=1,max=100"`
}

// MoneyResponse 金额响应DTO
type MoneyResponse struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

// AddressResponse 地址响应DTO
type AddressResponse struct {
	Province      string `json:"province"`
	City          string `json:"city"`
	District      string `json:"district"`
	DetailAddress string `json:"detailAddress"`
	ReceiverName  string `json:"receiverName"`
	ReceiverPhone string `json:"receiverPhone"`
	PostCode      string `json:"postCode,omitempty"`
}

// OrderItemResponse 订单项响应DTO
type OrderItemResponse struct {
	ID              string        `json:"id"`
	ProductID       string        `json:"productId"`
	ProductName     string        `json:"productName"`
	ProductImage    string        `json:"productImage"`
	ProductSKU      string        `json:"productSku"`
	Quantity        int           `json:"quantity"`
	UnitPrice       MoneyResponse `json:"unitPrice"`
	TotalPrice      MoneyResponse `json:"totalPrice"`
	Discount        MoneyResponse `json:"discount"`
	AttributeValues string        `json:"attributeValues,omitempty"`
}

// OrderResponse 订单响应DTO
type OrderResponse struct {
	ID                 string              `json:"id"`
	UserID             string              `json:"userId"`
	OrderSN            string              `json:"orderSn"`
	Status             string              `json:"status"`
	TotalAmount        MoneyResponse       `json:"totalAmount"`
	PayAmount          MoneyResponse       `json:"payAmount"`
	FreightAmount      MoneyResponse       `json:"freightAmount"`
	DiscountAmount     MoneyResponse       `json:"discountAmount,omitempty"`
	CouponAmount       MoneyResponse       `json:"couponAmount,omitempty"`
	PointAmount        MoneyResponse       `json:"pointAmount,omitempty"`
	PaymentMethod      string              `json:"paymentMethod,omitempty"`
	PaymentTime        *time.Time          `json:"paymentTime,omitempty"`
	DeliveryMethod     string              `json:"deliveryMethod"`
	DeliveryCompany    string              `json:"deliveryCompany,omitempty"`
	DeliveryTrackingNo string              `json:"deliveryTrackingNo,omitempty"`
	DeliveryTime       *time.Time          `json:"deliveryTime,omitempty"`
	ReceiveTime        *time.Time          `json:"receiveTime,omitempty"`
	CommentTime        *time.Time          `json:"commentTime,omitempty"`
	BillingAddress     AddressResponse     `json:"billingAddress"`
	ShippingAddress    AddressResponse     `json:"shippingAddress"`
	Note               string              `json:"note,omitempty"`
	OrderItems         []OrderItemResponse `json:"orderItems"`
	CreatedAt          time.Time           `json:"createdAt"`
	UpdatedAt          time.Time           `json:"updatedAt"`
}

// OrderBriefResponse 订单简要响应DTO
type OrderBriefResponse struct {
	ID            string        `json:"id"`
	OrderSN       string        `json:"orderSn"`
	Status        string        `json:"status"`
	TotalAmount   MoneyResponse `json:"totalAmount"`
	PayAmount     MoneyResponse `json:"payAmount"`
	PaymentMethod string        `json:"paymentMethod,omitempty"`
	ItemCount     int           `json:"itemCount"`
	CreatedAt     time.Time     `json:"createdAt"`
	UpdatedAt     time.Time     `json:"updatedAt"`
}

// OrderListResponse 订单列表响应DTO
type OrderListResponse struct {
	Items []OrderBriefResponse `json:"items"`
	Total int64                `json:"total"`
	Page  int                  `json:"page"`
	Size  int                  `json:"size"`
}

// OrderStatisticsResponse 订单统计响应DTO
type OrderStatisticsResponse struct {
	TotalCount      int64                `json:"totalCount"`
	TodayCount      int64                `json:"todayCount"`
	WeekCount       int64                `json:"weekCount"`
	MonthCount      int64                `json:"monthCount"`
	StatusCounts    map[string]int64     `json:"statusCounts"`
	TotalAmount     MoneyResponse        `json:"totalAmount"`
	TodayAmount     MoneyResponse        `json:"todayAmount"`
	WeekAmount      MoneyResponse        `json:"weekAmount"`
	MonthAmount     MoneyResponse        `json:"monthAmount"`
	TotalOrders     int64                `json:"totalOrders"`
	PendingPayment  int64                `json:"pendingPayment"`
	PendingShipment int64                `json:"pendingShipment"`
	Shipped         int64                `json:"shipped"`
	Completed       int64                `json:"completed"`
	Cancelled       int64                `json:"cancelled"`
	Refunding       int64                `json:"refunding"`
	Refunded        int64                `json:"refunded"`
	CollectionTime  *time.Time           `json:"collectionTime,omitempty"`
	TrendData       []OrderTrendDataItem `json:"trendData,omitempty"`
}

// OrderTrendDataItem 订单趋势数据项
type OrderTrendDataItem struct {
	Date   string        `json:"date"`
	Count  int64         `json:"count"`
	Amount MoneyResponse `json:"amount"`
}

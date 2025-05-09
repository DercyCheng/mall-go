package dto

// OrderItemRequest 订单项请求DTO
type OrderItemRequest struct {
	ProductID         string  `json:"productId" binding:"required"`
	ProductSn         string  `json:"productSn" binding:"required"`
	ProductName       string  `json:"productName" binding:"required"`
	ProductPic        string  `json:"productPic"`
	ProductPrice      float64 `json:"productPrice" binding:"required,gt=0"`
	ProductQuantity   int     `json:"productQuantity" binding:"required,gt=0"`
	ProductAttr       string  `json:"productAttr"`
	CouponAmount      float64 `json:"couponAmount"`
	PromotionAmount   float64 `json:"promotionAmount"`
	RealAmount        float64 `json:"realAmount" binding:"required,gt=0"`
	GiftIntegration   int     `json:"giftIntegration"`
	GiftGrowth        int     `json:"giftGrowth"`
	ProductCategoryId string  `json:"productCategoryId"`
}

// AddressRequest 收货地址请求DTO
type AddressRequest struct {
	Province      string `json:"province" binding:"required"`
	City          string `json:"city" binding:"required"`
	District      string `json:"district" binding:"required"`
	DetailAddress string `json:"detailAddress" binding:"required"`
	PostCode      string `json:"postCode"`
	Name          string `json:"name" binding:"required"`
	Phone         string `json:"phone" binding:"required"`
}

// OrderCreateRequest 创建订单请求DTO
type OrderCreateRequest struct {
	MemberID          string             `json:"memberId" binding:"required"`
	MemberUsername    string             `json:"memberUsername" binding:"required"`
	OrderItems        []OrderItemRequest `json:"orderItems" binding:"required,dive"`
	TotalAmount       float64            `json:"totalAmount" binding:"required,gt=0"`
	FreightAmount     float64            `json:"freightAmount"`
	PayType           string             `json:"payType"`
	SourceType        int                `json:"sourceType"`
	Note              string             `json:"note"`
	UseIntegration    int                `json:"useIntegration"`
	AutoConfirmDay    int                `json:"autoConfirmDay"`
	Address           AddressRequest     `json:"address" binding:"required"`
	BillType          int                `json:"billType"`
	BillHeader        string             `json:"billHeader"`
	BillContent       string             `json:"billContent"`
	BillReceiverPhone string             `json:"billReceiverPhone"`
	BillReceiverEmail string             `json:"billReceiverEmail"`
}

// OrderQueryRequest 订单查询请求DTO
type OrderQueryRequest struct {
	OrderSn         string `form:"orderSn"`
	Status          string `form:"status"`
	MemberUsername  string `form:"memberUsername"`
	ReceiverName    string `form:"receiverName"`
	ReceiverPhone   string `form:"receiverPhone"`
	CreateTimeBegin string `form:"createTimeBegin"`
	CreateTimeEnd   string `form:"createTimeEnd"`
	SourceType      int    `form:"sourceType" binding:"omitempty,oneof=-1 0 1"`
	OrderType       int    `form:"orderType" binding:"omitempty,oneof=-1 0 1"`
	Page            int    `form:"page" binding:"required,min=1"`
	Size            int    `form:"size" binding:"required,min=1,max=100"`
}

// OrderUpdateStatusRequest 更新订单状态请求DTO
type OrderUpdateStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// OrderUpdateNoteRequest 更新订单备注请求DTO
type OrderUpdateNoteRequest struct {
	Note string `json:"note" binding:"required"`
}

// OrderPayRequest 订单支付请求DTO
type OrderPayRequest struct {
	PayType       string  `json:"payType" binding:"required,oneof=alipay wechat unionpay cash credit"`
	PayAmount     float64 `json:"payAmount" binding:"required,gt=0"`
	TransactionId string  `json:"transactionId"`
}

// OrderDeliveryRequest 订单发货请求DTO
type OrderDeliveryRequest struct {
	DeliveryCompany string `json:"deliveryCompany" binding:"required"`
	DeliverySn      string `json:"deliverySn" binding:"required"`
}

// OrderRefundRequest 订单退款请求DTO
type OrderRefundRequest struct {
	Reason string  `json:"reason" binding:"required"`
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

// OrderCancelRequest 取消订单请求DTO
type OrderCancelRequest struct {
	Reason string `json:"reason" binding:"required"`
}

// OrderReceiveRequest 确认收货请求DTO
type OrderReceiveRequest struct {
	OrderId string `json:"orderId" binding:"required"`
}

// OrderUpdateReceiverInfoRequest 更新收货人信息请求DTO
type OrderUpdateReceiverInfoRequest struct {
	ReceiverName          string `json:"receiverName"`
	ReceiverPhone         string `json:"receiverPhone"`
	ReceiverPostCode      string `json:"receiverPostCode"`
	ReceiverProvince      string `json:"receiverProvince"`
	ReceiverCity          string `json:"receiverCity"`
	ReceiverDistrict      string `json:"receiverDistrict"`
	ReceiverDetailAddress string `json:"receiverDetailAddress"`
}

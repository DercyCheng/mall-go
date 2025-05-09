package dto

import (
	"time"
)

// OrderItemCreateRequest 订单项创建请求
type OrderItemCreateRequest struct {
	ProductID         string  `json:"productId" binding:"required"`         // 商品ID
	ProductSn         string  `json:"productSn" binding:"required"`         // 商品货号
	ProductName       string  `json:"productName" binding:"required"`       // 商品名称
	ProductPic        string  `json:"productPic"`                           // 商品主图
	ProductPrice      float64 `json:"productPrice" binding:"required"`      // 销售价格
	ProductQuantity   int     `json:"productQuantity" binding:"required"`   // 购买数量
	ProductAttr       string  `json:"productAttr"`                          // 商品销售属性:[{"key":"颜色","value":"颜色"},{"key":"容量","value":"4G"}]
	CouponAmount      float64 `json:"couponAmount"`                         // 优惠券优惠分解金额
	PromotionAmount   float64 `json:"promotionAmount"`                      // 促销优惠金额
	RealAmount        float64 `json:"realAmount"`                           // 该商品经过优惠后的分解金额
	GiftIntegration   int     `json:"giftIntegration"`                      // 赠送积分
	GiftGrowth        int     `json:"giftGrowth"`                           // 赠送成长值
	ProductCategoryId string  `json:"productCategoryId" binding:"required"` // 商品分类ID
}

// LogisticsResponse 物流信息响应
type LogisticsResponse struct {
	LogisticCode   string             `json:"logisticCode"`   // 物流编号
	ShipperCode    string             `json:"shipperCode"`    // 物流公司编号
	ShipperName    string             `json:"shipperName"`    // 物流公司名称
	Traces         []LogisticsTrace   `json:"traces"`         // 物流轨迹
	ReceiverInfo   AddressRequest     `json:"receiverInfo"`   // 收货人信息
	OrderInfo      LogisticsOrderInfo `json:"orderInfo"`      // 订单信息
	DeliveryTime   time.Time          `json:"deliveryTime"`   // 发货时间
	ShipmentStatus string             `json:"shipmentStatus"` // 物流状态
}

// LogisticsTrace 物流轨迹
type LogisticsTrace struct {
	AcceptTime    string `json:"acceptTime"`    // 接收时间
	AcceptStation string `json:"acceptStation"` // 接收站点
	Remark        string `json:"remark"`        // 备注
}

// LogisticsOrderInfo 物流订单信息
type LogisticsOrderInfo struct {
	OrderSn      string    `json:"orderSn"`      // 订单编号
	OrderAmount  float64   `json:"orderAmount"`  // 订单金额
	OrderTime    time.Time `json:"orderTime"`    // 下单时间
	OrderStatus  string    `json:"orderStatus"`  // 订单状态
	PaymentTime  time.Time `json:"paymentTime"`  // 支付时间
	DeliveryTime time.Time `json:"deliveryTime"` // 发货时间
}

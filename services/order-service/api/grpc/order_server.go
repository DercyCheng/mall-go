package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"mall-go/services/order-service/application/dto"
	applicationservice "mall-go/services/order-service/application/service"
	orderpb "mall-go/services/order-service/proto"
)

// OrderServer 实现OrderService gRPC接口
type OrderServer struct {
	orderpb.UnimplementedOrderServiceServer
	orderService *applicationservice.OrderService
}

// NewOrderServer 创建新的OrderServer
func NewOrderServer(orderService *applicationservice.OrderService) *OrderServer {
	return &OrderServer{
		orderService: orderService,
	}
}

// CreateOrder 创建订单
func (s *OrderServer) CreateOrder(ctx context.Context, req *orderpb.CreateOrderRequest) (*orderpb.CreateOrderResponse, error) {
	// 转换请求
	orderItems := make([]dto.OrderItemRequest, len(req.OrderItems))
	for i, item := range req.OrderItems {
		orderItems[i] = dto.OrderItemRequest{
			ProductID:         item.ProductId,
			ProductSn:         item.ProductSn,
			ProductName:       item.ProductName,
			ProductPic:        item.ProductPic,
			ProductPrice:      item.ProductPrice.Amount,
			ProductQuantity:   int(item.ProductQuantity),
			ProductAttr:       item.ProductAttr,
			CouponAmount:      item.CouponAmount,
			PromotionAmount:   item.PromotionAmount,
			RealAmount:        item.RealAmount,
			GiftIntegration:   int(item.GiftIntegration),
			GiftGrowth:        int(item.GiftGrowth),
			ProductCategoryId: item.ProductCategoryId,
		}
	}

	dtoReq := &dto.OrderCreateRequest{
		MemberID:       req.MemberId,
		MemberUsername: req.MemberUsername,
		Address: dto.AddressRequest{
			Province:      req.Address.Province,
			City:          req.Address.City,
			District:      req.Address.District,
			DetailAddress: req.Address.DetailAddress,
			PostCode:      req.Address.PostCode,
			Name:          req.Address.Name,
			Phone:         req.Address.Phone,
		},
		OrderItems:    orderItems,
		TotalAmount:   req.TotalAmount,
		FreightAmount: req.FreightAmount,
	}

	// 调用应用服务
	result, err := s.orderService.CreateOrder(ctx, dtoReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "创建订单失败: %v", err)
	}

	// 返回响应
	return &orderpb.CreateOrderResponse{
		Success: true,
		Message: "订单创建成功",
		OrderId: result.OrderID,
		OrderSn: result.OrderSn,
	}, nil
}

// GetOrder 获取订单详情
func (s *OrderServer) GetOrder(ctx context.Context, req *orderpb.GetOrderRequest) (*orderpb.GetOrderResponse, error) {
	// 调用应用服务
	result, err := s.orderService.GetOrder(ctx, req.OrderId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "获取订单失败: %v", err)
	}

	// 转换订单项
	orderItems := make([]*orderpb.OrderItem, len(result.OrderItems))
	for i, item := range result.OrderItems {
		orderItems[i] = &orderpb.OrderItem{
			Id:                item.ID,
			ProductId:         item.ProductID,
			ProductSn:         item.ProductSn,
			ProductName:       item.ProductName,
			ProductPic:        item.ProductPic,
			ProductPrice:      &orderpb.Money{Amount: item.ProductPrice, Currency: "CNY"},
			ProductQuantity:   int32(item.ProductQuantity),
			ProductAttr:       item.ProductAttr,
			CouponAmount:      item.CouponAmount,
			PromotionAmount:   item.PromotionAmount,
			RealAmount:        item.RealAmount,
			GiftIntegration:   int32(item.GiftIntegration),
			GiftGrowth:        int32(item.GiftGrowth),
			ProductCategoryId: item.ProductCategoryId,
			CreatedAt:         item.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	// 转换订单状态
	status, err := parseOrderStatus(result.Status)
	if err != nil {
		status = orderpb.OrderStatus_PENDING
	}

	// 转换支付方式
	payType, err := parsePaymentType(result.PayType)
	if err != nil {
		payType = orderpb.PaymentType_ALIPAY
	}

	// 构造响应
	order := &orderpb.Order{
		Id:                    result.ID,
		MemberId:              result.MemberID,
		OrderSn:               result.OrderSn,
		MemberUsername:        result.MemberUsername,
		TotalAmount:           &orderpb.Money{Amount: result.TotalAmount, Currency: "CNY"},
		PayAmount:             &orderpb.Money{Amount: result.PayAmount, Currency: "CNY"},
		FreightAmount:         &orderpb.Money{Amount: result.FreightAmount, Currency: "CNY"},
		PromotionAmount:       &orderpb.Money{Amount: result.PromotionAmount, Currency: "CNY"},
		IntegrationAmount:     &orderpb.Money{Amount: result.IntegrationAmount, Currency: "CNY"},
		CouponAmount:          &orderpb.Money{Amount: result.CouponAmount, Currency: "CNY"},
		DiscountAmount:        &orderpb.Money{Amount: result.DiscountAmount, Currency: "CNY"},
		PayType:               payType,
		SourceType:            int32(result.SourceType),
		Status:                status,
		OrderType:             int32(result.OrderType),
		DeliveryCompany:       result.DeliveryCompany,
		DeliverySn:            result.DeliverySn,
		AutoConfirmDay:        int32(result.AutoConfirmDay),
		Integration:           int32(result.Integration),
		Growth:                int32(result.Growth),
		PromotionInfo:         result.PromotionInfo,
		BillType:              int32(result.BillType),
		BillHeader:            result.BillHeader,
		BillContent:           result.BillContent,
		BillReceiverPhone:     result.BillReceiverPhone,
		BillReceiverEmail:     result.BillReceiverEmail,
		ReceiverName:          result.ReceiverName,
		ReceiverPhone:         result.ReceiverPhone,
		ReceiverPostCode:      result.ReceiverPostCode,
		ReceiverProvince:      result.ReceiverProvince,
		ReceiverCity:          result.ReceiverCity,
		ReceiverDistrict:      result.ReceiverDistrict,
		ReceiverDetailAddress: result.ReceiverDetailAddress,
		Note:                  result.Note,
		ConfirmStatus:         int32(result.ConfirmStatus),
		DeleteStatus:          int32(result.DeleteStatus),
		UseIntegration:        int32(result.UseIntegration),
		OrderItems:            orderItems,
		CreatedAt:             result.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if result.PaymentTime != nil {
		order.PaymentTime = result.PaymentTime.Format("2006-01-02 15:04:05")
	}
	if result.DeliveryTime != nil {
		order.DeliveryTime = result.DeliveryTime.Format("2006-01-02 15:04:05")
	}
	if result.ReceiveTime != nil {
		order.ReceiveTime = result.ReceiveTime.Format("2006-01-02 15:04:05")
	}
	if result.CommentTime != nil {
		order.CommentTime = result.CommentTime.Format("2006-01-02 15:04:05")
	}
	if result.ModifyTime != nil {
		order.ModifyTime = result.ModifyTime.Format("2006-01-02 15:04:05")
	}

	return &orderpb.GetOrderResponse{
		Success: true,
		Message: "获取订单成功",
		Order:   order,
	}, nil
}

// ListOrders 查询用户订单列表
func (s *OrderServer) ListOrders(ctx context.Context, req *orderpb.ListOrdersRequest) (*orderpb.ListOrdersResponse, error) {
	// 调用应用服务
	result, err := s.orderService.ListMemberOrders(ctx, req.MemberId, int(req.Page), int(req.Size))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "查询订单列表失败: %v", err)
	}

	// 转换订单列表
	orders := make([]*orderpb.Order, len(result.List))
	for i, item := range result.List {
		// 转换订单状态
		orderStatus, err := parseOrderStatus(item.Status)
		if err != nil {
			orderStatus = orderpb.OrderStatus_PENDING
		}

		orders[i] = &orderpb.Order{
			Id:             item.ID,
			OrderSn:        item.OrderSn,
			MemberId:       req.MemberId, // Use the member ID from the request instead
			MemberUsername: item.MemberUsername,
			TotalAmount:    &orderpb.Money{Amount: item.TotalAmount, Currency: "CNY"},
			PayAmount:      &orderpb.Money{Amount: item.PayAmount, Currency: "CNY"},
			Status:         orderStatus,
			OrderType:      int32(item.OrderType),
			ReceiverName:   item.ReceiverName,
			ReceiverPhone:  item.ReceiverPhone,
			CreatedAt:      item.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	return &orderpb.ListOrdersResponse{
		Success: true,
		Message: "查询订单列表成功",
		Orders:  orders,
		Total:   int32(result.Total),
	}, nil
}

// PayOrder 支付订单
func (s *OrderServer) PayOrder(ctx context.Context, req *orderpb.PayOrderRequest) (*orderpb.PayOrderResponse, error) {
	// 转换支付方式
	payType := "ALIPAY"
	switch req.PayType {
	case orderpb.PaymentType_ALIPAY:
		payType = "ALIPAY"
	case orderpb.PaymentType_WECHAT:
		payType = "WECHAT"
	case orderpb.PaymentType_UNIONPAY:
		payType = "UNIONPAY"
	case orderpb.PaymentType_CASH:
		payType = "CASH"
	case orderpb.PaymentType_CREDIT:
		payType = "CREDIT"
	}

	// 调用应用服务
	err := s.orderService.PayOrder(ctx, req.OrderId, &dto.OrderPayRequest{
		PayType: payType,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "支付订单失败: %v", err)
	}

	return &orderpb.PayOrderResponse{
		Success: true,
		Message: "订单支付成功",
	}, nil
}

// ShipOrder 订单发货
func (s *OrderServer) ShipOrder(ctx context.Context, req *orderpb.ShipOrderRequest) (*orderpb.ShipOrderResponse, error) {
	// 调用应用服务
	err := s.orderService.ShipOrder(ctx, req.OrderId, &dto.OrderDeliveryRequest{
		DeliveryCompany: req.DeliveryCompany,
		DeliverySn:      req.DeliverySn,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "订单发货失败: %v", err)
	}

	return &orderpb.ShipOrderResponse{
		Success: true,
		Message: "订单发货成功",
	}, nil
}

// ReceiveOrder 确认收货
func (s *OrderServer) ReceiveOrder(ctx context.Context, req *orderpb.ReceiveOrderRequest) (*orderpb.ReceiveOrderResponse, error) {
	// 调用应用服务
	err := s.orderService.ReceiveOrder(ctx, req.OrderId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "确认收货失败: %v", err)
	}

	return &orderpb.ReceiveOrderResponse{
		Success: true,
		Message: "确认收货成功",
	}, nil
}

// CancelOrder 取消订单
func (s *OrderServer) CancelOrder(ctx context.Context, req *orderpb.CancelOrderRequest) (*orderpb.CancelOrderResponse, error) {
	// 调用应用服务
	err := s.orderService.CancelOrder(ctx, req.OrderId, &dto.OrderCancelRequest{
		Reason: req.Reason,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "取消订单失败: %v", err)
	}

	return &orderpb.CancelOrderResponse{
		Success: true,
		Message: "取消订单成功",
	}, nil
}

// ApplyRefund 申请退款
func (s *OrderServer) ApplyRefund(ctx context.Context, req *orderpb.RefundRequest) (*orderpb.RefundResponse, error) {
	// 调用应用服务
	err := s.orderService.ApplyRefund(ctx, req.OrderId, &dto.OrderRefundRequest{
		Reason: req.Reason,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "申请退款失败: %v", err)
	}

	return &orderpb.RefundResponse{
		Success: true,
		Message: "申请退款成功",
	}, nil
}

// ConfirmRefund 确认退款
func (s *OrderServer) ConfirmRefund(ctx context.Context, req *orderpb.ConfirmRefundRequest) (*orderpb.ConfirmRefundResponse, error) {
	// 调用应用服务
	err := s.orderService.ConfirmRefund(ctx, req.OrderId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "确认退款失败: %v", err)
	}

	return &orderpb.ConfirmRefundResponse{
		Success: true,
		Message: "确认退款成功",
	}, nil
}

// 辅助函数 - 解析订单状态
func parseOrderStatus(orderStatus string) (orderpb.OrderStatus, error) {
	switch orderStatus {
	case "PENDING":
		return orderpb.OrderStatus_PENDING, nil
	case "PAID":
		return orderpb.OrderStatus_PAID, nil
	case "SHIPPING":
		return orderpb.OrderStatus_SHIPPING, nil
	case "COMPLETED":
		return orderpb.OrderStatus_COMPLETED, nil
	case "CANCELLED":
		return orderpb.OrderStatus_CANCELLED, nil
	case "CLOSED":
		return orderpb.OrderStatus_CLOSED, nil
	case "REFUNDING":
		return orderpb.OrderStatus_REFUNDING, nil
	case "REFUNDED":
		return orderpb.OrderStatus_REFUNDED, nil
	default:
		return orderpb.OrderStatus_PENDING, status.Errorf(codes.InvalidArgument, "无效的订单状态: %s", orderStatus)
	}
}

// 辅助函数 - 解析支付方式
func parsePaymentType(payType string) (orderpb.PaymentType, error) {
	switch payType {
	case "ALIPAY":
		return orderpb.PaymentType_ALIPAY, nil
	case "WECHAT":
		return orderpb.PaymentType_WECHAT, nil
	case "UNIONPAY":
		return orderpb.PaymentType_UNIONPAY, nil
	case "CASH":
		return orderpb.PaymentType_CASH, nil
	case "CREDIT":
		return orderpb.PaymentType_CREDIT, nil
	default:
		return orderpb.PaymentType_ALIPAY, status.Errorf(codes.InvalidArgument, "无效的支付方式: %s", payType)
	}
}

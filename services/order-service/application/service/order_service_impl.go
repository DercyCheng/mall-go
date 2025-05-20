package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"mall-go/services/order-service/application/dto"
	"mall-go/services/order-service/domain/model"
	"mall-go/services/order-service/domain/repository"
	domainService "mall-go/services/order-service/domain/service"
)

// orderServiceImpl 订单应用服务实现
type orderServiceImpl struct {
	orderDomainService *domainService.OrderDomainService
	orderRepo          repository.OrderRepository
	// 可以添加其他依赖，如产品服务、库存服务、用户服务等
}

// NewOrderService 创建订单应用服务
func NewOrderService(
	orderDomainService *domainService.OrderDomainService,
	orderRepo repository.OrderRepository,
) OrderService {
	return &orderServiceImpl{
		orderDomainService: orderDomainService,
		orderRepo:          orderRepo,
	}
}

// CreateOrder 创建订单
func (s *orderServiceImpl) CreateOrder(ctx context.Context, req dto.OrderCreateRequest) (dto.OrderResponse, error) {
	// 验证请求
	if req.UserID == "" || len(req.OrderItems) == 0 {
		return dto.OrderResponse{}, errors.New("missing required fields")
	}

	// 创建订单对象
	order := &model.Order{
		ID:              uuid.New().String(),
		UserID:          req.UserID,
		Status:          model.OrderStatusPending,
		DeliveryMethod:  model.DeliveryMethod(req.DeliveryMethod),
		Note:            req.Note,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		ShippingAddress: convertToAddressModel(req.ShippingAddress),
	}

	// 设置账单地址，如果没有提供则使用配送地址
	if req.BillingAddress != nil {
		order.BillingAddress = convertToAddressModel(*req.BillingAddress)
	} else {
		order.BillingAddress = order.ShippingAddress
	}

	// 处理订单项
	var totalAmount float64 = 0
	orderItems := make([]model.OrderItem, 0, len(req.OrderItems))

	for _, item := range req.OrderItems {
		// 这里应该调用产品服务获取商品信息
		// 简化处理，假设已获取商品信息
		productName := "Product " + item.ProductID // 应该从产品服务获取
		productImage := "image_url"                // 应该从产品服务获取
		productSKU := "SKU_" + item.ProductID      // 应该从产品服务获取
		unitPrice := 100.0                         // 应该从产品服务获取

		// 计算订单项金额
		itemTotalPrice := unitPrice * float64(item.Quantity)
		totalAmount += itemTotalPrice

		// 创建订单项
		orderItem := model.OrderItem{
			ID:              uuid.New().String(),
			ProductID:       item.ProductID,
			ProductName:     productName,
			ProductImage:    productImage,
			ProductSKU:      productSKU,
			Quantity:        item.Quantity,
			UnitPrice:       model.Money{Amount: unitPrice, Currency: "CNY"},
			TotalPrice:      model.Money{Amount: itemTotalPrice, Currency: "CNY"},
			Discount:        model.Money{Amount: 0, Currency: "CNY"},
			AttributeValues: item.AttributeValues,
		}

		orderItems = append(orderItems, orderItem)
	}

	// 设置订单金额
	order.TotalAmount = model.Money{Amount: totalAmount, Currency: "CNY"}
	order.PayAmount = model.Money{Amount: totalAmount, Currency: "CNY"} // 简化处理，暂不计算折扣
	order.OrderItems = orderItems

	// 调用领域服务创建订单
	err := s.orderDomainService.CreateOrder(ctx, order)
	if err != nil {
		return dto.OrderResponse{}, err
	}

	// 注意：在实际项目中，这里还应该调用库存服务预留商品库存

	// 转换为DTO响应
	return convertToOrderResponse(order), nil
}

// PayOrder 支付订单
func (s *orderServiceImpl) PayOrder(ctx context.Context, req dto.OrderPayRequest) error {
	// 验证请求
	if req.OrderID == "" || req.PaymentMethod == "" {
		return errors.New("missing required fields")
	}

	// 调用领域服务支付订单
	return s.orderDomainService.PayOrder(
		ctx,
		req.OrderID,
		model.PaymentMethod(req.PaymentMethod),
		time.Now(),
	)

	// 注意：在实际项目中，这里还应该调用库存服务从预留状态转为实际扣减
}

// ShipOrder 发货订单
func (s *orderServiceImpl) ShipOrder(ctx context.Context, req dto.OrderShipRequest) error {
	// 验证请求
	if req.OrderID == "" || req.DeliveryCompany == "" || req.TrackingNo == "" {
		return errors.New("missing required fields")
	}

	// 调用领域服务发货订单
	return s.orderDomainService.ShipOrder(
		ctx,
		req.OrderID,
		req.DeliveryCompany,
		req.TrackingNo,
	)
}

// UpdateOrderStatus 更新订单状态
func (s *orderServiceImpl) UpdateOrderStatus(ctx context.Context, req dto.OrderStatusUpdateRequest) error {
	// 验证请求
	if req.OrderID == "" || req.Status == "" {
		return errors.New("missing required fields")
	}

	// 获取订单
	order, err := s.orderRepo.FindByID(ctx, req.OrderID)
	if err != nil {
		return err
	}

	// 根据目标状态调用不同的领域服务方法
	switch model.OrderStatus(req.Status) {
	case model.OrderStatusPaid:
		return s.orderDomainService.PayOrder(ctx, req.OrderID, order.PaymentMethod, time.Now())
	case model.OrderStatusShipping:
		return s.orderDomainService.ShipOrder(ctx, req.OrderID, order.DeliveryCompany, order.DeliveryTrackingNo)
	case model.OrderStatusDelivered:
		return s.orderDomainService.ReceiveOrder(ctx, req.OrderID)
	case model.OrderStatusCompleted:
		return s.orderDomainService.CompleteOrder(ctx, req.OrderID)
	case model.OrderStatusCancelled:
		return s.orderDomainService.CancelOrder(ctx, req.OrderID, req.Reason)
	case model.OrderStatusRefunding:
		return s.orderDomainService.ApplyRefund(ctx, req.OrderID, req.Reason)
	case model.OrderStatusRefunded:
		return s.orderDomainService.RefundOrder(ctx, req.OrderID)
	default:
		return errors.New("invalid order status")
	}
}

// CancelOrder 取消订单
func (s *orderServiceImpl) CancelOrder(ctx context.Context, orderID string, reason string) error {
	if orderID == "" {
		return errors.New("order ID is required")
	}

	// 调用领域服务取消订单
	err := s.orderDomainService.CancelOrder(ctx, orderID, reason)
	if err == nil {
		// 注意：在实际项目中，这里还应该调用库存服务释放预留的库存
	}
	return err
}

// GetOrder 获取订单详情
func (s *orderServiceImpl) GetOrder(ctx context.Context, orderID string) (dto.OrderResponse, error) {
	if orderID == "" {
		return dto.OrderResponse{}, errors.New("order ID is required")
	}

	// 查询订单
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return dto.OrderResponse{}, err
	}

	// 转换为DTO响应
	return convertToOrderResponse(order), nil
}

// GetOrderByOrderSN 根据订单编号获取订单
func (s *orderServiceImpl) GetOrderByOrderSN(ctx context.Context, orderSN string) (dto.OrderResponse, error) {
	if orderSN == "" {
		return dto.OrderResponse{}, errors.New("order SN is required")
	}

	// 查询订单
	order, err := s.orderRepo.FindByOrderSN(ctx, orderSN)
	if err != nil {
		return dto.OrderResponse{}, err
	}

	// 转换为DTO响应
	return convertToOrderResponse(order), nil
}

// GetUserOrders 获取用户订单列表
func (s *orderServiceImpl) GetUserOrders(ctx context.Context, userID string, page, size int) (dto.OrderListResponse, error) {
	if userID == "" {
		return dto.OrderListResponse{}, errors.New("user ID is required")
	}

	// 查询用户订单
	orders, total, err := s.orderRepo.FindByUserID(ctx, userID, page, size)
	if err != nil {
		return dto.OrderListResponse{}, err
	}

	// 转换为DTO响应
	return convertToOrderListResponse(orders, total, page, size), nil
}

// SearchOrders 搜索订单
func (s *orderServiceImpl) SearchOrders(ctx context.Context, req dto.OrderQueryRequest) (dto.OrderListResponse, error) {
	var orders []*model.Order
	var total int64
	var err error

	// 根据查询条件搜索
	if req.Keyword != "" {
		orders, total, err = s.orderRepo.Search(ctx, req.Keyword, req.Page, req.Size)
	} else if req.Status != "" {
		orders, total, err = s.orderRepo.FindByStatus(ctx, model.OrderStatus(req.Status), req.Page, req.Size)
	} else if req.StartDate != "" && req.EndDate != "" {
		orders, total, err = s.orderRepo.FindByDateRange(ctx, req.StartDate, req.EndDate, req.Page, req.Size)
	} else if req.ProductID != "" {
		orders, total, err = s.orderRepo.FindByProductID(ctx, req.ProductID, req.Page, req.Size)
	} else {
		orders, total, err = s.orderRepo.FindAll(ctx, req.Page, req.Size)
	}

	if err != nil {
		return dto.OrderListResponse{}, err
	}

	// 转换为DTO响应
	return convertToOrderListResponse(orders, total, req.Page, req.Size), nil
}

// GetOrderStatistics 获取订单统计数据
func (s *orderServiceImpl) GetOrderStatistics(ctx context.Context) (dto.OrderStatisticsResponse, error) {
	// 调用仓储方法统计各状态订单数量
	statusCounts, err := s.orderRepo.CountByStatus(ctx)
	if err != nil {
		return dto.OrderStatisticsResponse{}, err
	}

	// 设置收集时间
	now := time.Now()
	collectionTime := &now

	// 转换为DTO响应
	response := dto.OrderStatisticsResponse{
		TotalOrders:     0,
		PendingPayment:  0,
		PendingShipment: 0,
		Shipped:         0,
		Completed:       0,
		Cancelled:       0,
		Refunding:       0,
		Refunded:        0,
		StatusCounts:    make(map[string]int64),
		CollectionTime:  collectionTime,
	}

	// 填充状态统计
	for status, count := range statusCounts {
		response.StatusCounts[string(status)] = count
		response.TotalOrders += count

		// 根据状态填充各分类统计
		switch status {
		case model.OrderStatusPending:
			response.PendingPayment = count
		case model.OrderStatusPaid:
			response.PendingShipment = count
		case model.OrderStatusShipping:
			response.Shipped = count
		case model.OrderStatusCompleted:
			response.Completed = count
		case model.OrderStatusCancelled:
			response.Cancelled = count
		case model.OrderStatusRefunding:
			response.Refunding = count
		case model.OrderStatusRefunded:
			response.Refunded = count
		}
	}

	return response, nil
}

// DeleteOrder 删除订单
func (s *orderServiceImpl) DeleteOrder(ctx context.Context, orderID string) error {
	if orderID == "" {
		return errors.New("order ID is required")
	}

	// 执行软删除
	return s.orderRepo.Delete(ctx, orderID)
}

// 辅助方法：将地址DTO转换为地址模型
func convertToAddressModel(address dto.AddressRequest) model.Address {
	return model.Address{
		Province:      address.Province,
		City:          address.City,
		District:      address.District,
		DetailAddress: address.DetailAddress,
		ReceiverName:  address.ReceiverName,
		ReceiverPhone: address.ReceiverPhone,
		PostCode:      address.PostCode,
	}
}

// 辅助方法：将订单模型转换为订单响应DTO
func convertToOrderResponse(order *model.Order) dto.OrderResponse {
	// 转换订单项
	orderItems := make([]dto.OrderItemResponse, 0, len(order.OrderItems))
	for _, item := range order.OrderItems {
		orderItems = append(orderItems, dto.OrderItemResponse{
			ID:              item.ID,
			ProductID:       item.ProductID,
			ProductName:     item.ProductName,
			ProductImage:    item.ProductImage,
			ProductSKU:      item.ProductSKU,
			Quantity:        item.Quantity,
			UnitPrice:       dto.MoneyResponse{Amount: item.UnitPrice.Amount, Currency: item.UnitPrice.Currency},
			TotalPrice:      dto.MoneyResponse{Amount: item.TotalPrice.Amount, Currency: item.TotalPrice.Currency},
			Discount:        dto.MoneyResponse{Amount: item.Discount.Amount, Currency: item.Discount.Currency},
			AttributeValues: item.AttributeValues,
		})
	}

	// 转换地址
	shippingAddress := dto.AddressResponse{
		Province:      order.ShippingAddress.Province,
		City:          order.ShippingAddress.City,
		District:      order.ShippingAddress.District,
		DetailAddress: order.ShippingAddress.DetailAddress,
		ReceiverName:  order.ShippingAddress.ReceiverName,
		ReceiverPhone: order.ShippingAddress.ReceiverPhone,
		PostCode:      order.ShippingAddress.PostCode,
	}

	billingAddress := dto.AddressResponse{
		Province:      order.BillingAddress.Province,
		City:          order.BillingAddress.City,
		District:      order.BillingAddress.District,
		DetailAddress: order.BillingAddress.DetailAddress,
		ReceiverName:  order.BillingAddress.ReceiverName,
		ReceiverPhone: order.BillingAddress.ReceiverPhone,
		PostCode:      order.BillingAddress.PostCode,
	}

	// 构建响应DTO
	return dto.OrderResponse{
		ID:                 order.ID,
		UserID:             order.UserID,
		OrderSN:            order.OrderSN,
		Status:             string(order.Status),
		TotalAmount:        dto.MoneyResponse{Amount: order.TotalAmount.Amount, Currency: order.TotalAmount.Currency},
		PayAmount:          dto.MoneyResponse{Amount: order.PayAmount.Amount, Currency: order.PayAmount.Currency},
		FreightAmount:      dto.MoneyResponse{Amount: order.FreightAmount.Amount, Currency: order.FreightAmount.Currency},
		DiscountAmount:     dto.MoneyResponse{Amount: order.DiscountAmount.Amount, Currency: order.DiscountAmount.Currency},
		CouponAmount:       dto.MoneyResponse{Amount: order.CouponAmount.Amount, Currency: order.CouponAmount.Currency},
		PointAmount:        dto.MoneyResponse{Amount: order.PointAmount.Amount, Currency: order.PointAmount.Currency},
		PaymentMethod:      string(order.PaymentMethod),
		PaymentTime:        &order.PaymentTime,
		DeliveryMethod:     string(order.DeliveryMethod),
		DeliveryCompany:    order.DeliveryCompany,
		DeliveryTrackingNo: order.DeliveryTrackingNo,
		DeliveryTime:       &order.DeliveryTime,
		ReceiveTime:        &order.ReceiveTime,
		CommentTime:        &order.CommentTime,
		BillingAddress:     billingAddress,
		ShippingAddress:    shippingAddress,
		Note:               order.Note,
		OrderItems:         orderItems,
		CreatedAt:          order.CreatedAt,
		UpdatedAt:          order.UpdatedAt,
	}
}

// 辅助方法：将订单列表转换为订单列表响应DTO
func convertToOrderListResponse(orders []*model.Order, total int64, page, size int) dto.OrderListResponse {
	items := make([]dto.OrderBriefResponse, 0, len(orders))
	for _, order := range orders {
		// 创建简要的订单响应
		briefResponse := dto.OrderBriefResponse{
			ID:            order.ID,
			OrderSN:       order.OrderSN,
			Status:        string(order.Status),
			TotalAmount:   dto.MoneyResponse{Amount: order.TotalAmount.Amount, Currency: order.TotalAmount.Currency},
			PaymentMethod: string(order.PaymentMethod),
			CreatedAt:     order.CreatedAt,
			UpdatedAt:     order.UpdatedAt,
		}
		items = append(items, briefResponse)
	}

	return dto.OrderListResponse{
		Items: items,
		Total: total,
		Page:  page,
		Size:  size,
	}
}

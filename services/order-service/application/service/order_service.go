package applicationservice

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"mall-go/services/order-service/application/dto"
	"mall-go/services/order-service/domain/model"
	"mall-go/services/order-service/domain/repository"
	domainservice "mall-go/services/order-service/domain/service"
)

// EventPublisher 事件发布接口
type EventPublisher interface {
	Publish(ctx context.Context, eventName string, data interface{}) error
}

// OrderService 订单应用服务
type OrderService struct {
	orderRepository repository.OrderRepository
	orderDomainSvc  *domainservice.OrderDomainService
	eventPublisher  EventPublisher
}

// NewOrderService 创建订单应用服务
func NewOrderService(
	orderRepository repository.OrderRepository,
	orderDomainSvc *domainservice.OrderDomainService,
	eventPublisher EventPublisher,
) *OrderService {
	return &OrderService{
		orderRepository: orderRepository,
		orderDomainSvc:  orderDomainSvc,
		eventPublisher:  eventPublisher,
	}
}

// CreateOrder 创建订单
func (s *OrderService) CreateOrder(ctx context.Context, req *dto.OrderCreateRequest) (*dto.CreateOrderResponse, error) {
	// 1. 参数验证
	if req == nil {
		return nil, errors.New("请求参数不能为空")
	}

	// 2. 转换为领域模型
	// 2.1 创建订单项
	orderItems := make([]model.OrderItem, len(req.OrderItems))
	for i, item := range req.OrderItems {
		orderItems[i] = model.OrderItem{
			ID:                uuid.New().String(),
			ProductID:         item.ProductID,
			ProductSn:         item.ProductSn,
			ProductName:       item.ProductName,
			ProductPic:        item.ProductPic,
			ProductPrice:      model.Money{Amount: item.ProductPrice, Currency: "CNY"},
			ProductQuantity:   item.ProductQuantity,
			ProductAttr:       item.ProductAttr,
			CouponAmount:      item.CouponAmount,
			PromotionAmount:   item.PromotionAmount,
			RealAmount:        item.RealAmount,
			GiftIntegration:   item.GiftIntegration,
			GiftGrowth:        item.GiftGrowth,
			ProductCategoryId: item.ProductCategoryId,
			CreatedAt:         time.Now(),
		}
	}

	// 2.2 创建收货地址
	address := &model.Address{
		Province:      req.Address.Province,
		City:          req.Address.City,
		District:      req.Address.District,
		DetailAddress: req.Address.DetailAddress,
		PostCode:      req.Address.PostCode,
		Name:          req.Address.Name,
		Phone:         req.Address.Phone,
	}

	// 2.3 创建订单
	order, err := model.NewOrder(
		req.MemberID,
		req.MemberUsername,
		address,
		orderItems,
		req.TotalAmount,
		req.FreightAmount,
	)
	if err != nil {
		return nil, err
	}

	// 3. 设置额外的订单属性
	if req.PayType != "" {
		order.PayType = model.PaymentType(req.PayType)
	}
	if req.SourceType != 0 {
		order.SourceType = req.SourceType
	}
	if req.Note != "" {
		order.Note = req.Note
	}
	if req.UseIntegration > 0 {
		order.UseIntegration = req.UseIntegration
	}
	if req.AutoConfirmDay > 0 {
		order.AutoConfirmDay = req.AutoConfirmDay
	}
	// 设置发票信息
	order.BillType = req.BillType
	order.BillHeader = req.BillHeader
	order.BillContent = req.BillContent
	order.BillReceiverPhone = req.BillReceiverPhone
	order.BillReceiverEmail = req.BillReceiverEmail

	// 4. 计算订单金额
	err = s.orderDomainSvc.CalculateOrderAmount(ctx, order)
	if err != nil {
		return nil, err
	}

	// 5. 确认订单(执行业务规则检查、库存锁定等)
	err = s.orderDomainSvc.ConfirmOrder(ctx, order)
	if err != nil {
		return nil, err
	}

	// 6. 发布订单创建事件
	event := map[string]interface{}{
		"orderId":  order.ID,
		"orderSn":  order.OrderSn,
		"memberId": order.MemberID,
		"status":   order.Status,
		"amount":   order.PayAmount.Amount,
	}
	if err := s.eventPublisher.Publish(ctx, "order.created", event); err != nil {
		// 仅记录日志，不影响业务流程
		// s.logger.Error("Failed to publish order.created event", zap.Error(err))
	}

	// 7. 返回响应
	return &dto.CreateOrderResponse{
		OrderID:   order.ID,
		OrderSn:   order.OrderSn,
		PayAmount: order.PayAmount.Amount,
	}, nil
}

// GetOrder 获取订单详情
func (s *OrderService) GetOrder(ctx context.Context, id string) (*dto.OrderResponse, error) {
	// 1. 查询订单
	order, err := s.orderRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("订单不存在")
	}

	// 2. 转换为DTO
	return convertOrderToDTO(order), nil
}

// GetOrderByOrderSn 根据订单编号获取订单详情
func (s *OrderService) GetOrderByOrderSn(ctx context.Context, orderSn string) (*dto.OrderResponse, error) {
	// 1. 查询订单
	order, err := s.orderRepository.FindByOrderSn(ctx, orderSn)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("订单不存在")
	}

	// 2. 转换为DTO
	return convertOrderToDTO(order), nil
}

// ListMemberOrders 查询会员订单列表
func (s *OrderService) ListMemberOrders(ctx context.Context, memberID string, page, size int) (*dto.OrderListResponse, error) {
	// 1. 查询会员订单列表
	orders, total, err := s.orderRepository.FindByMemberID(ctx, memberID, page, size)
	if err != nil {
		return nil, err
	}

	// 2. 转换为DTO
	orderList := make([]dto.OrderBriefResponse, len(orders))
	for i, order := range orders {
		orderList[i] = dto.OrderBriefResponse{
			ID:             order.ID,
			OrderSn:        order.OrderSn,
			MemberUsername: order.MemberUsername,
			TotalAmount:    order.TotalAmount.Amount,
			PayAmount:      order.PayAmount.Amount,
			Status:         string(order.Status),
			OrderType:      order.OrderType,
			ReceiverName:   order.ReceiverName,
			ReceiverPhone:  order.ReceiverPhone,
			CreatedAt:      order.CreatedAt,
		}
	}

	// 3. 返回响应
	return &dto.OrderListResponse{
		Total: total,
		List:  orderList,
	}, nil
}

// ListOrders 分页查询订单列表
func (s *OrderService) ListOrders(ctx context.Context, req *dto.OrderQueryRequest) (*dto.OrderListResponse, error) {
	// 1. 构造查询条件
	query := repository.OrderQuery{
		OrderSn:         req.OrderSn,
		Status:          req.Status,
		MemberUsername:  req.MemberUsername,
		ReceiverName:    req.ReceiverName,
		ReceiverPhone:   req.ReceiverPhone,
		CreateTimeBegin: req.CreateTimeBegin,
		CreateTimeEnd:   req.CreateTimeEnd,
		SourceType:      req.SourceType,
		OrderType:       req.OrderType,
		Page:            req.Page,
		Size:            req.Size,
	}

	// 2. 查询订单列表
	orders, total, err := s.orderRepository.List(ctx, query)
	if err != nil {
		return nil, err
	}

	// 3. 转换为DTO
	orderList := make([]dto.OrderBriefResponse, len(orders))
	for i, order := range orders {
		orderList[i] = dto.OrderBriefResponse{
			ID:             order.ID,
			OrderSn:        order.OrderSn,
			MemberUsername: order.MemberUsername,
			TotalAmount:    order.TotalAmount.Amount,
			PayAmount:      order.PayAmount.Amount,
			Status:         string(order.Status),
			OrderType:      order.OrderType,
			ReceiverName:   order.ReceiverName,
			ReceiverPhone:  order.ReceiverPhone,
			CreatedAt:      order.CreatedAt,
		}
	}

	// 4. 返回响应
	return &dto.OrderListResponse{
		Total: total,
		List:  orderList,
	}, nil
}

// PayOrder 支付订单
func (s *OrderService) PayOrder(ctx context.Context, id string, req *dto.OrderPayRequest) error {
	// 1. 查询订单
	order, err := s.orderRepository.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.New("订单不存在")
	}

	// 2. 支付订单
	err = order.Pay(model.PaymentType(req.PayType))
	if err != nil {
		return err
	}

	// 3. 更新订单
	err = s.orderRepository.Update(ctx, order)
	if err != nil {
		return err
	}

	// 4. 支付成功后的业务处理
	go func() {
		// 使用新的上下文，避免原上下文取消导致处理失败
		newCtx := context.Background()

		// 完成已支付订单的后续处理(积分、优惠券等)
		if err := s.orderDomainSvc.CompletePaidOrder(newCtx, order.ID); err != nil {
			// 记录错误日志，可能需要人工干预或重试机制
			// s.logger.Error("Failed to complete paid order", zap.Error(err), zap.String("orderId", order.ID))
		}
	}()

	// 5. 发布订单支付事件
	event := map[string]interface{}{
		"orderId":       order.ID,
		"orderSn":       order.OrderSn,
		"memberId":      order.MemberID,
		"status":        order.Status,
		"payType":       req.PayType,
		"payAmount":     req.PayAmount,
		"transactionId": req.TransactionId,
		"paymentTime":   order.PaymentTime,
	}
	if err := s.eventPublisher.Publish(ctx, "order.paid", event); err != nil {
		// 仅记录日志，不影响业务流程
		// s.logger.Error("Failed to publish order.paid event", zap.Error(err))
	}

	return nil
}

// ShipOrder 订单发货
func (s *OrderService) ShipOrder(ctx context.Context, id string, req *dto.OrderDeliveryRequest) error {
	// 1. 查询订单
	order, err := s.orderRepository.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.New("订单不存在")
	}

	// 2. 发货
	err = order.Ship(req.DeliveryCompany, req.DeliverySn)
	if err != nil {
		return err
	}

	// 3. 更新订单
	err = s.orderRepository.Update(ctx, order)
	if err != nil {
		return err
	}

	// 4. 发布订单发货事件
	event := map[string]interface{}{
		"orderId":         order.ID,
		"orderSn":         order.OrderSn,
		"memberId":        order.MemberID,
		"status":          order.Status,
		"deliveryCompany": req.DeliveryCompany,
		"deliverySn":      req.DeliverySn,
		"deliveryTime":    order.DeliveryTime,
	}
	if err := s.eventPublisher.Publish(ctx, "order.shipped", event); err != nil {
		// 仅记录日志，不影响业务流程
		// s.logger.Error("Failed to publish order.shipped event", zap.Error(err))
	}

	return nil
}

// ReceiveOrder 确认收货
func (s *OrderService) ReceiveOrder(ctx context.Context, id string) error {
	// 1. 查询订单
	order, err := s.orderRepository.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.New("订单不存在")
	}

	// 2. 确认收货
	err = order.Receive()
	if err != nil {
		return err
	}

	// 3. 更新订单
	err = s.orderRepository.Update(ctx, order)
	if err != nil {
		return err
	}

	// 4. 发布订单完成事件
	event := map[string]interface{}{
		"orderId":     order.ID,
		"orderSn":     order.OrderSn,
		"memberId":    order.MemberID,
		"status":      order.Status,
		"receiveTime": order.ReceiveTime,
	}
	if err := s.eventPublisher.Publish(ctx, "order.completed", event); err != nil {
		// 仅记录日志，不影响业务流程
		// s.logger.Error("Failed to publish order.completed event", zap.Error(err))
	}

	return nil
}

// CancelOrder 取消订单
func (s *OrderService) CancelOrder(ctx context.Context, id string, req *dto.OrderCancelRequest) error {
	// 1. 调用领域服务取消订单(包含跨服务操作)
	err := s.orderDomainSvc.CancelOrder(ctx, id, req.Reason)
	if err != nil {
		return err
	}

	// 2. 发布订单取消事件
	order, err := s.orderRepository.FindByID(ctx, id)
	if err != nil {
		// 记录错误日志，但不中断流程
		// s.logger.Error("Failed to find order after cancel", zap.Error(err))
	} else {
		event := map[string]interface{}{
			"orderId":  order.ID,
			"orderSn":  order.OrderSn,
			"memberId": order.MemberID,
			"status":   order.Status,
			"reason":   req.Reason,
		}
		if err := s.eventPublisher.Publish(ctx, "order.cancelled", event); err != nil {
			// 仅记录日志，不影响业务流程
			// s.logger.Error("Failed to publish order.cancelled event", zap.Error(err))
		}
	}

	return nil
}

// CloseOrder 关闭订单
func (s *OrderService) CloseOrder(ctx context.Context, id string) error {
	// 1. 查询订单
	order, err := s.orderRepository.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.New("订单不存在")
	}

	// 2. 关闭订单
	err = order.Close()
	if err != nil {
		return err
	}

	// 3. 更新订单
	err = s.orderRepository.Update(ctx, order)
	if err != nil {
		return err
	}

	// 4. 发布订单关闭事件
	event := map[string]interface{}{
		"orderId":  order.ID,
		"orderSn":  order.OrderSn,
		"memberId": order.MemberID,
		"status":   order.Status,
	}
	if err := s.eventPublisher.Publish(ctx, "order.closed", event); err != nil {
		// 仅记录日志，不影响业务流程
		// s.logger.Error("Failed to publish order.closed event", zap.Error(err))
	}

	return nil
}

// DeleteOrder 删除订单
func (s *OrderService) DeleteOrder(ctx context.Context, id string) error {
	// 1. 查询订单
	order, err := s.orderRepository.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.New("订单不存在")
	}

	// 2. 只允许删除已取消、已关闭或已完成的订单
	if order.Status != model.OrderStatusCancelled &&
		order.Status != model.OrderStatusClosed &&
		order.Status != model.OrderStatusCompleted {
		return errors.New("订单状态不允许删除")
	}

	// 3. 删除订单(逻辑删除)
	err = s.orderRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	// 4. 发布订单删除事件
	event := map[string]interface{}{
		"orderId":  order.ID,
		"orderSn":  order.OrderSn,
		"memberId": order.MemberID,
	}
	if err := s.eventPublisher.Publish(ctx, "order.deleted", event); err != nil {
		// 仅记录日志，不影响业务流程
		// s.logger.Error("Failed to publish order.deleted event", zap.Error(err))
	}

	return nil
}

// UpdateOrderNote 更新订单备注
func (s *OrderService) UpdateOrderNote(ctx context.Context, id string, req *dto.OrderUpdateNoteRequest) error {
	// 1. 查询订单
	order, err := s.orderRepository.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.New("订单不存在")
	}

	// 2. 更新订单备注
	order.UpdateNote(req.Note)

	// 3. 更新订单
	return s.orderRepository.UpdateNote(ctx, id, req.Note)
}

// UpdateOrderReceiverInfo 更新收货人信息
func (s *OrderService) UpdateOrderReceiverInfo(ctx context.Context, id string, req *dto.OrderUpdateReceiverInfoRequest) error {
	// 1. 查询订单
	order, err := s.orderRepository.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.New("订单不存在")
	}

	// 2. 构建更新信息
	receiverInfo := map[string]string{}

	if req.ReceiverName != "" {
		receiverInfo["receiverName"] = req.ReceiverName
	}
	if req.ReceiverPhone != "" {
		receiverInfo["receiverPhone"] = req.ReceiverPhone
	}
	if req.ReceiverPostCode != "" {
		receiverInfo["receiverPostCode"] = req.ReceiverPostCode
	}
	if req.ReceiverProvince != "" {
		receiverInfo["receiverProvince"] = req.ReceiverProvince
	}
	if req.ReceiverCity != "" {
		receiverInfo["receiverCity"] = req.ReceiverCity
	}
	if req.ReceiverDistrict != "" {
		receiverInfo["receiverDistrict"] = req.ReceiverDistrict
	}
	if req.ReceiverDetailAddress != "" {
		receiverInfo["receiverDetailAddress"] = req.ReceiverDetailAddress
	}

	// 3. 更新收货人信息
	return s.orderRepository.UpdateReceiverInfo(ctx, id, receiverInfo)
}

// ApplyRefund 申请退款
func (s *OrderService) ApplyRefund(ctx context.Context, id string, req *dto.OrderRefundRequest) error {
	// 1. 查询订单
	order, err := s.orderRepository.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.New("订单不存在")
	}

	// 2. 申请退款
	err = order.ApplyRefund(req.Reason)
	if err != nil {
		return err
	}

	// 3. 更新订单
	err = s.orderRepository.Update(ctx, order)
	if err != nil {
		return err
	}

	// 4. 发布退款申请事件
	event := map[string]interface{}{
		"orderId":  order.ID,
		"orderSn":  order.OrderSn,
		"memberId": order.MemberID,
		"status":   order.Status,
		"reason":   req.Reason,
		"amount":   req.Amount,
	}
	if err := s.eventPublisher.Publish(ctx, "order.refund.applied", event); err != nil {
		// 仅记录日志，不影响业务流程
		// s.logger.Error("Failed to publish order.refund.applied event", zap.Error(err))
	}

	return nil
}

// ConfirmRefund 确认退款
func (s *OrderService) ConfirmRefund(ctx context.Context, id string) error {
	// 1. 查询订单
	order, err := s.orderRepository.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.New("订单不存在")
	}

	// 2. 确认退款
	err = order.Refund()
	if err != nil {
		return err
	}

	// 3. 更新订单
	err = s.orderRepository.Update(ctx, order)
	if err != nil {
		return err
	}

	// 4. 发布退款确认事件
	event := map[string]interface{}{
		"orderId":  order.ID,
		"orderSn":  order.OrderSn,
		"memberId": order.MemberID,
		"status":   order.Status,
	}
	if err := s.eventPublisher.Publish(ctx, "order.refund.confirmed", event); err != nil {
		// 仅记录日志，不影响业务流程
		// s.logger.Error("Failed to publish order.refund.confirmed event", zap.Error(err))
	}

	return nil
}

// convertOrderToDTO 将订单领域模型转换为DTO
func convertOrderToDTO(order *model.Order) *dto.OrderResponse {
	if order == nil {
		return nil
	}

	// 转换订单项
	orderItems := make([]dto.OrderItemResponse, len(order.OrderItems))
	for i, item := range order.OrderItems {
		orderItems[i] = dto.OrderItemResponse{
			ID:                item.ID,
			ProductID:         item.ProductID,
			ProductSn:         item.ProductSn,
			ProductName:       item.ProductName,
			ProductPic:        item.ProductPic,
			ProductPrice:      item.ProductPrice.Amount,
			ProductQuantity:   item.ProductQuantity,
			ProductAttr:       item.ProductAttr,
			CouponAmount:      item.CouponAmount,
			PromotionAmount:   item.PromotionAmount,
			RealAmount:        item.RealAmount,
			GiftIntegration:   item.GiftIntegration,
			GiftGrowth:        item.GiftGrowth,
			ProductCategoryId: item.ProductCategoryId,
			CreatedAt:         item.CreatedAt,
		}
	}

	// 获取状态和支付方式的中文名称
	statusName, ok := dto.OrderStatusMap[string(order.Status)]
	if !ok {
		statusName = string(order.Status)
	}

	paymentTypeName, ok := dto.PaymentTypeMap[string(order.PayType)]
	if !ok {
		paymentTypeName = string(order.PayType)
	}

	sourceTypeName, ok := dto.OrderSourceTypeMap[order.SourceType]
	if !ok {
		sourceTypeName = "未知来源"
	}

	orderTypeName, ok := dto.OrderTypeMap[order.OrderType]
	if !ok {
		orderTypeName = "未知类型"
	}

	// 处理可能为零值的时间
	var paymentTime, deliveryTime, receiveTime, commentTime, modifyTime *time.Time
	if !order.PaymentTime.IsZero() {
		paymentTime = &order.PaymentTime
	}
	if !order.DeliveryTime.IsZero() {
		deliveryTime = &order.DeliveryTime
	}
	if !order.ReceiveTime.IsZero() {
		receiveTime = &order.ReceiveTime
	}
	if !order.CommentTime.IsZero() {
		commentTime = &order.CommentTime
	}
	if !order.ModifyTime.IsZero() {
		modifyTime = &order.ModifyTime
	}

	return &dto.OrderResponse{
		ID:                    order.ID,
		MemberID:              order.MemberID,
		OrderSn:               order.OrderSn,
		MemberUsername:        order.MemberUsername,
		TotalAmount:           order.TotalAmount.Amount,
		PayAmount:             order.PayAmount.Amount,
		FreightAmount:         order.FreightAmount.Amount,
		PromotionAmount:       order.PromotionAmount.Amount,
		IntegrationAmount:     order.IntegrationAmount.Amount,
		CouponAmount:          order.CouponAmount.Amount,
		DiscountAmount:        order.DiscountAmount.Amount,
		PayType:               string(order.PayType),
		SourceType:            order.SourceType,
		Status:                string(order.Status),
		OrderType:             order.OrderType,
		DeliveryCompany:       order.DeliveryCompany,
		DeliverySn:            order.DeliverySn,
		AutoConfirmDay:        order.AutoConfirmDay,
		Integration:           order.Integration,
		Growth:                order.Growth,
		PromotionInfo:         order.PromotionInfo,
		BillType:              order.BillType,
		BillHeader:            order.BillHeader,
		BillContent:           order.BillContent,
		BillReceiverPhone:     order.BillReceiverPhone,
		BillReceiverEmail:     order.BillReceiverEmail,
		ReceiverName:          order.ReceiverName,
		ReceiverPhone:         order.ReceiverPhone,
		ReceiverPostCode:      order.ReceiverPostCode,
		ReceiverProvince:      order.ReceiverProvince,
		ReceiverCity:          order.ReceiverCity,
		ReceiverDistrict:      order.ReceiverDistrict,
		ReceiverDetailAddress: order.ReceiverDetailAddress,
		Note:                  order.Note,
		ConfirmStatus:         order.ConfirmStatus,
		DeleteStatus:          order.DeleteStatus,
		UseIntegration:        order.UseIntegration,
		PaymentTime:           paymentTime,
		DeliveryTime:          deliveryTime,
		ReceiveTime:           receiveTime,
		CommentTime:           commentTime,
		ModifyTime:            modifyTime,
		OrderItems:            orderItems,
		CreatedAt:             order.CreatedAt,
		UpdatedAt:             order.UpdatedAt,
		StatusName:            statusName,
		PaymentTypeName:       paymentTypeName,
		SourceTypeName:        sourceTypeName,
		OrderTypeName:         orderTypeName,
	}
}

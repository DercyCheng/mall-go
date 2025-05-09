package domainservice

import (
	"context"
	"errors"
	"fmt"
	"time"

	"mall-go/services/order-service/domain/model"
	"mall-go/services/order-service/domain/repository"
	"mall-go/services/order-service/infrastructure/grpc"
	productpb "mall-go/services/product-service/proto"
	userpb "mall-go/services/user-service/proto"
)

// OrderDomainService 订单领域服务
// 处理跨聚合根或复杂的业务逻辑
type OrderDomainService struct {
	orderRepo  repository.OrderRepository
	grpcClient *grpc.ClientManager
}

// NewOrderDomainService 创建订单领域服务
func NewOrderDomainService(
	orderRepo repository.OrderRepository,
	grpcClient *grpc.ClientManager,
) *OrderDomainService {
	return &OrderDomainService{
		orderRepo:  orderRepo,
		grpcClient: grpcClient,
	}
}

// ConfirmOrder 确认订单
// 在创建订单时执行业务规则检查、库存锁定等操作
func (s *OrderDomainService) ConfirmOrder(ctx context.Context, order *model.Order) error {
	// 验证用户是否存在
	userClient, err := s.grpcClient.GetUserClient(ctx)
	if err != nil {
		return fmt.Errorf("获取用户服务客户端失败: %w", err)
	}

	userResp, err := userClient.GetUserInfo(ctx, &userpb.GetUserInfoRequest{
		UserId: order.MemberID,
	})
	if err != nil {
		return fmt.Errorf("验证用户信息失败: %w", err)
	}

	if !userResp.Success {
		return errors.New("用户不存在")
	}

	// 验证商品库存
	productClient, err := s.grpcClient.GetProductClient(ctx)
	if err != nil {
		return fmt.Errorf("获取产品服务客户端失败: %w", err)
	}

	for _, item := range order.OrderItems {
		// 获取产品库存
		stockResp, err := productClient.GetStock(ctx, &productpb.GetStockRequest{
			ProductId: item.ProductID,
		})
		if err != nil {
			return fmt.Errorf("获取产品 %s 库存失败: %w", item.ProductID, err)
		}

		if !stockResp.Success {
			return fmt.Errorf("获取产品 %s 库存信息失败: %s", item.ProductID, stockResp.Message)
		}

		// 检查库存是否充足
		if int(stockResp.Stock) < item.ProductQuantity {
			return fmt.Errorf("产品 %s 库存不足，剩余 %d，需要 %d",
				item.ProductName, stockResp.Stock, item.ProductQuantity)
		}

		// 锁定库存
		_, err = productClient.UpdateStock(ctx, &productpb.UpdateStockRequest{
			ProductId: item.ProductID,
			Stock:     stockResp.Stock - int32(item.ProductQuantity),
			LockStock: stockResp.LockStock + int32(item.ProductQuantity),
		})
		if err != nil {
			return fmt.Errorf("锁定产品 %s 库存失败: %w", item.ProductID, err)
		}
	}

	// 保存订单
	return s.orderRepo.Save(ctx, order)
}

// CancelOrder 取消订单
// 取消订单可能涉及到库存释放等跨服务操作
func (s *OrderDomainService) CancelOrder(ctx context.Context, orderID string, reason string) error {
	// 查询订单
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.New("订单不存在")
	}

	// 检查订单状态，只有待付款订单可以取消
	if order.Status != model.OrderStatusPending {
		return errors.New("只有待付款订单可以取消")
	}

	// 恢复库存
	productClient, err := s.grpcClient.GetProductClient(ctx)
	if err != nil {
		return fmt.Errorf("获取产品服务客户端失败: %w", err)
	}

	for _, item := range order.OrderItems {
		// 获取当前库存状态
		stockResp, err := productClient.GetStock(ctx, &productpb.GetStockRequest{
			ProductId: item.ProductID,
		})
		if err != nil {
			return fmt.Errorf("获取产品 %s 库存失败: %w", item.ProductID, err)
		}

		// 恢复库存，解锁锁定的库存
		_, err = productClient.UpdateStock(ctx, &productpb.UpdateStockRequest{
			ProductId: item.ProductID,
			Stock:     stockResp.Stock + int32(item.ProductQuantity),
			LockStock: stockResp.LockStock - int32(item.ProductQuantity),
		})
		if err != nil {
			return fmt.Errorf("恢复产品 %s 库存失败: %w", item.ProductID, err)
		}
	}

	// 取消订单
	err = order.Cancel(reason)
	if err != nil {
		return err
	}

	// 更新订单
	return s.orderRepo.Update(ctx, order)
}

// CloseTimeoutOrders 关闭超时未支付订单
func (s *OrderDomainService) CloseTimeoutOrders(ctx context.Context, timeout time.Duration) error {
	// 查询超时未支付订单
	// 这里需要仓储层支持根据创建时间和状态查询订单
	// orders, err := s.orderRepository.FindTimeoutOrders(ctx, timeout)
	// if err != nil {
	//     return err
	// }

	// 关闭超时订单
	// for _, order := range orders {
	//     if order.Status == model.OrderStatusPending {
	//         order.Close()
	//         s.orderRepository.Update(ctx, order)
	//
	//         // 释放库存
	//         // for _, item := range order.OrderItems {
	//         //     s.productClient.UnlockStock(ctx, item.ProductID, item.ProductQuantity)
	//         // }
	//     }
	// }

	return nil
}

// AutoConfirmReceivedOrders 自动确认收货
func (s *OrderDomainService) AutoConfirmReceivedOrders(ctx context.Context) error {
	// 查询需要自动确认收货的订单
	// 即发货时间 + 自动确认天数 < 当前时间 的订单
	// orders, err := s.orderRepository.FindAutoConfirmOrders(ctx)
	// if err != nil {
	//     return err
	// }

	// 自动确认收货
	// for _, order := range orders {
	//     if order.Status == model.OrderStatusShipping {
	//         order.Receive()
	//         s.orderRepository.Update(ctx, order)
	//     }
	// }

	return nil
}

// CalculateOrderAmount 计算订单金额
func (s *OrderDomainService) CalculateOrderAmount(ctx context.Context, order *model.Order) error {
	// 计算总金额
	var totalAmount float64
	for _, item := range order.OrderItems {
		totalAmount += item.ProductPrice.Amount * float64(item.ProductQuantity)
	}
	order.TotalAmount = model.Money{Amount: totalAmount, Currency: "CNY"}

	// 优惠金额
	var promotionAmount float64
	for _, item := range order.OrderItems {
		promotionAmount += item.PromotionAmount
	}
	order.PromotionAmount = model.Money{Amount: promotionAmount, Currency: "CNY"}

	// 优惠券金额
	var couponAmount float64
	for _, item := range order.OrderItems {
		couponAmount += item.CouponAmount
	}
	order.CouponAmount = model.Money{Amount: couponAmount, Currency: "CNY"}

	// 积分抵扣金额
	var integrationAmount float64
	if order.UseIntegration > 0 {
		// 假设1积分=0.01元
		integrationAmount = float64(order.UseIntegration) * 0.01
	}
	order.IntegrationAmount = model.Money{Amount: integrationAmount, Currency: "CNY"}

	// 计算实际支付金额
	payAmount := totalAmount - promotionAmount - couponAmount - integrationAmount + order.FreightAmount.Amount
	if payAmount < 0 {
		payAmount = 0
	}
	order.PayAmount = model.Money{Amount: payAmount, Currency: "CNY"}

	return nil
}

// ApplyPromotion 应用促销
func (s *OrderDomainService) ApplyPromotion(ctx context.Context, order *model.Order) error {
	// 此处应调用促销服务，获取适用的促销规则并应用
	// 示例逻辑
	// promotions, err := s.promotionClient.GetPromotions(ctx, order.MemberID, order.OrderItems)
	// if err != nil {
	//     return err
	// }

	// 应用促销规则
	// var promotionAmount float64 = 0
	// for _, promotion := range promotions {
	//     if promotion.Type == "discount" {
	//         promotionAmount += promotion.Amount
	//     }
	// }

	// 更新订单金额
	// order.PromotionAmount.Amount = promotionAmount
	// order.PayAmount.Amount = order.TotalAmount.Amount + order.FreightAmount.Amount - order.PromotionAmount.Amount - order.CouponAmount.Amount - order.IntegrationAmount.Amount

	return nil
}

// ApplyCoupon 应用优惠券
func (s *OrderDomainService) ApplyCoupon(ctx context.Context, order *model.Order, couponIDs []string) error {
	// 此处应调用优惠券服务，验证优惠券并计算优惠金额
	// 示例逻辑
	// coupons, err := s.couponClient.GetCoupons(ctx, order.MemberID, couponIDs)
	// if err != nil {
	//     return err
	// }

	// 计算优惠券优惠金额
	// var couponAmount float64 = 0
	// for _, coupon := range coupons {
	//     couponAmount += coupon.Amount
	// }

	// 更新订单金额
	// order.CouponAmount.Amount = couponAmount
	// order.PayAmount.Amount = order.TotalAmount.Amount + order.FreightAmount.Amount - order.PromotionAmount.Amount - order.CouponAmount.Amount - order.IntegrationAmount.Amount

	return nil
}

// ApplyIntegration 应用积分
func (s *OrderDomainService) ApplyIntegration(ctx context.Context, order *model.Order, useIntegration int) error {
	// 此处应调用会员服务，验证会员积分并计算抵扣金额
	// 示例逻辑
	// member, err := s.memberClient.GetMember(ctx, order.MemberID)
	// if err != nil {
	//     return err
	// }

	// 检查积分是否足够
	// if member.Integration < useIntegration {
	//     return errors.New("积分不足")
	// }

	// 计算积分抵扣金额(假设100积分=1元)
	// integrationAmount := float64(useIntegration) / 100.0

	// 更新订单金额
	// order.UseIntegration = useIntegration
	// order.IntegrationAmount.Amount = integrationAmount
	// order.PayAmount.Amount = order.TotalAmount.Amount + order.FreightAmount.Amount - order.PromotionAmount.Amount - order.CouponAmount.Amount - order.IntegrationAmount.Amount

	return nil
}

// CompletePaidOrder 完成已支付订单的后续处理
func (s *OrderDomainService) CompletePaidOrder(ctx context.Context, orderID string) error {
	// 查询订单
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.New("订单不存在")
	}

	// 确认产品库存变更
	productClient, err := s.grpcClient.GetProductClient(ctx)
	if err != nil {
		return fmt.Errorf("获取产品服务客户端失败: %w", err)
	}

	for _, item := range order.OrderItems {
		// 获取当前库存状态
		stockResp, err := productClient.GetStock(ctx, &productpb.GetStockRequest{
			ProductId: item.ProductID,
		})
		if err != nil {
			return fmt.Errorf("获取产品 %s 库存失败: %w", item.ProductID, err)
		}

		// 解锁库存并扣减已锁定的库存
		_, err = productClient.UpdateStock(ctx, &productpb.UpdateStockRequest{
			ProductId: item.ProductID,
			Stock:     stockResp.Stock,
			LockStock: stockResp.LockStock - int32(item.ProductQuantity),
		})
		if err != nil {
			return fmt.Errorf("解锁产品 %s 库存失败: %w", item.ProductID, err)
		}
	}

	// 更新订单状态
	if order.Status != model.OrderStatusPaid {
		order.Status = model.OrderStatusPaid
		order.PaymentTime = time.Now()
		err = s.orderRepo.Update(ctx, order)
		if err != nil {
			return fmt.Errorf("更新订单支付状态失败: %w", err)
		}
	}

	return nil
}

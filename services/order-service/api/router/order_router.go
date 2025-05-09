package router

import (
	"github.com/gin-gonic/gin"

	"mall-go/services/order-service/api/handler"
	"mall-go/services/order-service/api/middleware"
)

// OrderRouter 订单路由
type OrderRouter struct {
	orderHandler *handler.OrderHandler
	jwtSecret    string
}

// NewOrderRouter 创建订单路由
func NewOrderRouter(orderHandler *handler.OrderHandler, jwtSecret string) *OrderRouter {
	return &OrderRouter{
		orderHandler: orderHandler,
		jwtSecret:    jwtSecret,
	}
}

// Register 注册路由
func (r *OrderRouter) Register(e *gin.Engine) {
	// 基础路由组
	v1 := e.Group("/api/v1")

	// 健康检查
	v1.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "order-service",
		})
	})

	// 订单管理路由组 - 需要JWT认证
	orderAuth := v1.Group("/orders")
	orderAuth.Use(middleware.JWT(r.jwtSecret))
	{
		// 创建订单
		orderAuth.POST("", r.orderHandler.Create)
		// 获取用户订单列表
		orderAuth.GET("", r.orderHandler.ListMemberOrders)
		// 获取订单详情
		orderAuth.GET("/:id", r.orderHandler.Get)
		// 根据订单编号获取订单
		orderAuth.GET("/by-order-sn/:orderSn", r.orderHandler.GetByOrderSn)
		// 取消订单
		orderAuth.POST("/:id/cancel", r.orderHandler.Cancel)
		// 支付订单
		orderAuth.POST("/:id/pay", r.orderHandler.Pay)
		// 确认收货
		orderAuth.POST("/:id/receive", r.orderHandler.Receive)
		// 删除订单
		orderAuth.DELETE("/:id", r.orderHandler.Delete)
		// 申请退款
		orderAuth.POST("/:id/refund/apply", r.orderHandler.ApplyRefund)
		// 查看物流信息 - 需要实现此方法
		// orderAuth.GET("/:id/logistics", r.orderHandler.GetLogistics)
	}

	// 后台管理路由组 - 需要JWT认证和管理员权限
	adminAuth := v1.Group("/admin/orders")
	adminAuth.Use(middleware.JWT(r.jwtSecret))
	adminAuth.Use(middleware.Admin())
	{
		// 后台订单列表
		adminAuth.GET("", r.orderHandler.List)
		// 订单详情
		adminAuth.GET("/:id", r.orderHandler.Get)
		// 订单发货
		adminAuth.POST("/:id/ship", r.orderHandler.Ship)
		// 关闭订单 - 需要实现此方法
		// adminAuth.POST("/:id/close", r.orderHandler.CloseOrder)
		// 订单备注
		adminAuth.PUT("/:id/note", r.orderHandler.UpdateNote)
		// 修改收货人信息
		adminAuth.PUT("/:id/receiver", r.orderHandler.UpdateReceiverInfo)
		// 确认退款
		adminAuth.POST("/:id/refund/confirm", r.orderHandler.ConfirmRefund)
	}
}

package router

import (
	"github.com/gin-gonic/gin"

	"mall-go/services/order-service/api/handler"
	"mall-go/services/order-service/api/middleware"
)

// SetupRouter 设置HTTP路由
func SetupRouter(
	orderHandler *handler.OrderHandler,
) *gin.Engine {
	r := gin.Default()

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// API v1 路由组
	v1 := r.Group("/api/v1")

	// 使用认证中间件（实际项目中应实现适当的认证逻辑）
	authMiddleware := middleware.JWTAuth()

	// 订单相关路由
	orders := v1.Group("/orders")
	{
		// 公共API - 无需认证
		orders.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "Order service is healthy"})
		})

		// 需要认证的API
		authorized := orders.Group("/")
		authorized.Use(authMiddleware)
		{
			// 创建订单
			authorized.POST("", orderHandler.CreateOrder)

			// 支付订单
			authorized.POST("/pay", orderHandler.PayOrder)

			// 发货订单
			authorized.POST("/ship", orderHandler.ShipOrder)

			// 更新订单状态
			authorized.PUT("/status", orderHandler.UpdateOrderStatus)

			// 取消订单
			authorized.PUT("/:id/cancel", orderHandler.CancelOrder)

			// 获取订单详情
			authorized.GET("/:id", orderHandler.GetOrder)

			// 根据订单编号获取订单
			authorized.GET("/by-order-sn/:orderSN", orderHandler.GetOrderByOrderSN)

			// 获取用户订单列表
			authorized.GET("/user/:userID", orderHandler.GetUserOrders)

			// 获取当前用户的订单列表
			authorized.GET("/my-orders", orderHandler.GetUserOrders)

			// 搜索订单
			authorized.GET("/search", orderHandler.SearchOrders)

			// 获取订单统计数据
			authorized.GET("/statistics", orderHandler.GetOrderStatistics)

			// 删除订单
			authorized.DELETE("/:id", orderHandler.DeleteOrder)
		}
	}

	return r
}

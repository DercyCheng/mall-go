package router

import (
	"github.com/gin-gonic/gin"

	"mall-go/services/inventory-service/api/handler"
	"mall-go/services/inventory-service/api/middleware"
)

// SetupRouter 设置HTTP路由
func SetupRouter(
	inventoryHandler *handler.InventoryHandler,
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

	// 库存相关路由
	inventory := v1.Group("/inventory")
	{
		// 公共API - 无需认证
		inventory.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "Inventory service is healthy"})
		})

		// 允许公开查询的API
		inventory.GET("/product/:productId", inventoryHandler.GetInventoryByProductID)
		inventory.GET("/sku/:sku", inventoryHandler.GetInventoryBySKU)

		// 需要认证的API
		authorized := inventory.Group("/")
		authorized.Use(authMiddleware)
		{
			// 创建库存
			authorized.POST("", inventoryHandler.CreateInventory)

			// 更新库存信息
			authorized.PUT("/:id", inventoryHandler.UpdateInventory)

			// 调整库存数量
			authorized.POST("/adjust", inventoryHandler.AdjustInventory)

			// 入库操作
			authorized.POST("/inbound", inventoryHandler.InboundInventory)

			// 出库操作
			authorized.POST("/outbound", inventoryHandler.OutboundInventory)

			// 预留库存
			authorized.POST("/reserve", inventoryHandler.ReserveInventory)

			// 释放预留库存
			authorized.POST("/release", inventoryHandler.ReleaseInventory)

			// 锁定库存
			authorized.POST("/lock", inventoryHandler.LockInventory)

			// 解锁库存
			authorized.POST("/unlock", inventoryHandler.UnlockInventory)

			// 获取库存详情
			authorized.GET("/:id", inventoryHandler.GetInventory)

			// 获取低库存商品列表
			authorized.GET("/low-stock", inventoryHandler.GetLowStockInventories)

			// 搜索库存
			authorized.GET("/search", inventoryHandler.SearchInventories)

			// 获取库存操作记录
			authorized.GET("/operations/product/:productId", inventoryHandler.GetInventoryOperations)

			// 根据订单ID获取库存操作记录
			authorized.GET("/operations/order/:orderId", inventoryHandler.GetInventoryOperationsByOrderID)
		}
	}

	return r
}

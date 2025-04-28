package router

import (
	"github.com/gin-gonic/gin"

	"mall-go/services/product-service/api/handler"
	"mall-go/services/product-service/api/middleware"
)

// SetupRouter 设置API路由
func SetupRouter(r *gin.Engine, productHandler *handler.ProductHandler) {
	// 添加全局中间件
	r.Use(middleware.Cors())
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())

	// 健康检查接口
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})

	// API 版本分组
	v1 := r.Group("/api")
	{
		// 商品相关路由
		products := v1.Group("/products")
		{
			products.GET("", productHandler.List)            // 获取商品列表
			products.POST("", productHandler.Create)         // 创建商品
			products.GET("/search", productHandler.Search)   // 搜索商品
			products.GET("/:id", productHandler.Get)         // 获取商品详情
			products.PUT("/:id", productHandler.Update)      // 更新商品
			products.DELETE("/:id", productHandler.Delete)   // 删除商品
			products.PATCH("/:id/status", productHandler.UpdateStatus) // 更新商品状态
			products.POST("/:id/publish", productHandler.Publish)      // 发布商品
		}
	}
}
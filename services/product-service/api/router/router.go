package router

import (
	"github.com/gin-gonic/gin"

	"mall-go/services/product-service/api/handler"
	"mall-go/services/product-service/api/middleware"
)

// SetupRouter 设置路由
func SetupRouter(
	r *gin.Engine,
	productHandler *handler.ProductHandler,
) {
	// 中间件
	r.Use(middleware.Cors())
	r.Use(middleware.RequestLogger())
	r.Use(middleware.Recovery())

	// API v1 版本
	v1 := r.Group("/api/v1")
	{
		// 产品相关接口
		products := v1.Group("/products")
		{
			products.GET("", productHandler.List)                        // 获取产品列表
			products.POST("", productHandler.Create)                     // 创建产品
			products.GET("/:id", productHandler.Get)                     // 获取产品详情
			products.PUT("/:id", productHandler.Update)                  // 更新产品
			products.DELETE("/:id", productHandler.Delete)               // 删除产品
			products.PUT("/publish-status", productHandler.UpdatePublishStatus)     // 批量更新上架状态
			products.PUT("/new-status", productHandler.UpdateNewStatus)             // 批量更新新品状态
			products.PUT("/recommend-status", productHandler.UpdateRecommendStatus) // 批量更新推荐状态
			products.POST("/category/transfer", productHandler.TransferCategory)    // 转移商品分类
		}

		// 可以在此处添加更多产品相关路由，如品牌、分类等
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "UP",
		})
	})
}
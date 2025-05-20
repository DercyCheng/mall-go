package main

import (
	"mall-go/services/oss-service/api"
	"mall-go/services/oss-service/configs"

	"github.com/gin-gonic/gin"

	"go.uber.org/zap"
)

func main() {
	// 初始化配置
	cfg, err := configs.LoadConfig()
	if err != nil {
		panic(err)
	}

	// 初始化日志
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// 初始化Gin引擎
	router := gin.Default()

	// 注册路由
	api.RegisterRoutes(router, cfg, logger)

	// 启动服务
	logger.Info("Starting OSS service",
		zap.String("port", cfg.Server.Port))
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}

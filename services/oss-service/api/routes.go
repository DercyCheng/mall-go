package api

import (
	"mall-go/services/oss-service/configs"
	"mall-go/services/oss-service/infrastructure"
	"mall-go/services/oss-service/internal/repository"
	"mall-go/services/oss-service/internal/service"

	"github.com/gin-gonic/gin"

	"go.uber.org/zap"
)

func RegisterRoutes(router *gin.Engine, cfg *configs.Config, logger *zap.Logger) {
	// 初始化MongoDB
	mongoDB, err := infrastructure.NewMongoDB(&cfg.MongoDB, logger)
	if err != nil {
		logger.Fatal("Failed to initialize MongoDB", zap.Error(err))
	}
	defer mongoDB.Close()

	// 初始化repository
	fileRepo := repository.NewGridFSRepository(mongoDB.GridFS, logger)

	// 初始化service
	fileService := service.NewFileService(fileRepo, &cfg.Storage, logger)

	// 初始化handler
	fileHandler := NewFileHandler(fileService, logger)

	// 文件API路由组
	fileGroup := router.Group("/api/files")
	{
		fileGroup.POST("/upload", fileHandler.Upload)
		fileGroup.GET("/:id", fileHandler.Download)
		fileGroup.DELETE("/:id", fileHandler.Delete)
		fileGroup.GET("/:id/info", fileHandler.GetInfo)
	}

	logger.Info("API routes registered successfully")
}

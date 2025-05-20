package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mall-go/services/inventory-service/api/handler"
	"mall-go/services/inventory-service/api/router"
	"mall-go/services/inventory-service/application/service"
	"mall-go/services/inventory-service/domain/repository"
	domainService "mall-go/services/inventory-service/domain/service"
	"mall-go/services/inventory-service/infrastructure/config"
	"mall-go/services/inventory-service/infrastructure/persistence"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 初始化数据库连接
	db, err := persistence.NewDBConnection(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 初始化仓储
	var inventoryRepo repository.InventoryRepository
	inventoryRepo = persistence.NewInventoryRepository(db)

	// 初始化领域服务
	inventoryDomainService := domainService.NewInventoryDomainService(inventoryRepo)

	// 初始化应用服务
	inventoryService := service.NewInventoryService(inventoryDomainService, inventoryRepo)

	// 初始化处理器
	inventoryHandler := handler.NewInventoryHandler(inventoryService)

	// 设置路由
	r := router.SetupRouter(inventoryHandler)

	// 创建HTTP服务
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: r,
	}

	// 在单独的goroutine中启动服务
	go func() {
		log.Printf("Starting Inventory service on port %s\n", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}

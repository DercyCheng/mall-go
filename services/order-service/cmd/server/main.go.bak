package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	mysqldriver "gorm.io/driver/mysql"
	"gorm.io/gorm"

	ordergrpc "mall-go/services/order-service/api/grpc"
	"mall-go/services/order-service/api/handler"
	"mall-go/services/order-service/api/router"
	applicationservice "mall-go/services/order-service/application/service"
	domainservice "mall-go/services/order-service/domain/service"
	"mall-go/services/order-service/infrastructure/event"
	grpcclient "mall-go/services/order-service/infrastructure/grpc"
	mysqlrepo "mall-go/services/order-service/infrastructure/persistence/mysql"
	orderpb "mall-go/services/order-service/proto"
)

func main() {
	// 初始化日志
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	// 加载配置文件
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("./services/order-service/configs")
	if err := viper.ReadInConfig(); err != nil {
		logger.Fatal("Failed to read config file", zap.Error(err))
	}

	// 数据库配置
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("database.username"),
		viper.GetString("database.password"),
		viper.GetString("database.host"),
		viper.GetInt("database.port"),
		viper.GetString("database.database"),
	)

	// 初始化数据库连接
	db, err := gorm.Open(mysqldriver.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatal("Failed to connect database", zap.Error(err))
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatal("Failed to get db connection", zap.Error(err))
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 初始化事件发布器
	eventPublisher, err := event.NewRabbitMQEventPublisher(
		viper.GetString("rabbitmq.url"),
		viper.GetString("rabbitmq.exchange"),
		logger,
	)
	if err != nil {
		logger.Fatal("Failed to initialize event publisher", zap.Error(err))
	}
	defer eventPublisher.Close()

	// 初始化 gRPC 客户端管理器
	grpcClientManager := grpcclient.NewClientManager(logger)
	defer grpcClientManager.Close()

	// 初始化仓储层
	orderRepository := mysqlrepo.NewOrderRepository(db)

	// 初始化领域服务
	orderDomainService := domainservice.NewOrderDomainService(orderRepository, grpcClientManager)

	// 初始化应用服务
	orderAppService := applicationservice.NewOrderService(orderRepository, orderDomainService, eventPublisher)

	// 启动 gRPC 服务器
	go startGRPCServer(orderAppService, logger)

	// 初始化处理器
	orderHandler := handler.NewOrderHandler(orderAppService)

	// 初始化路由
	orderRouter := router.NewOrderRouter(orderHandler, viper.GetString("jwt.secret"))

	// 初始化Gin
	gin.SetMode(viper.GetString("server.mode"))
	engine := gin.Default()

	// 注册路由
	orderRouter.Register(engine)

	// 启动服务器
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("server.port")),
		Handler: engine,
	}

	// 优雅关闭
	go func() {
		logger.Info("Starting HTTP server", zap.String("address", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start HTTP server", zap.Error(err))
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down servers...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("HTTP Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}

// startGRPCServer 启动 gRPC 服务器
func startGRPCServer(orderService *applicationservice.OrderService, logger *zap.Logger) {
	// 获取 gRPC 配置
	grpcPort := viper.GetInt("grpc.port")
	if grpcPort == 0 {
		grpcPort = 50051 // 默认端口
	}

	// 创建监听器
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		logger.Fatal("Failed to listen for gRPC", zap.Error(err))
	}

	// 创建 gRPC 服务器
	grpcServer := grpc.NewServer()

	// 注册反射服务 - 支持像 grpcurl 这样的工具进行服务发现
	reflection.Register(grpcServer)

	// 创建 gRPC 服务处理器
	orderServer := ordergrpc.NewOrderServer(orderService)

	// 注册服务
	orderpb.RegisterOrderServiceServer(grpcServer, orderServer)

	// 启动 gRPC 服务器
	logger.Info("Starting gRPC server", zap.Int("port", grpcPort))
	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal("Failed to serve gRPC", zap.Error(err))
	}
}

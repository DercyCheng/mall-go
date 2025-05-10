package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"mall-go/services/inventory-service/api/handler"
	"mall-go/services/inventory-service/application/service"
	"mall-go/services/inventory-service/domain/repository"
	"mall-go/services/inventory-service/infrastructure/config"
	"mall-go/services/inventory-service/infrastructure/persistence/mysql"
	"mall-go/services/inventory-service/infrastructure/persistence/redis"
	"mall-go/services/inventory-service/proto/inventorypb"
)

// Package main serves as the entry point for the inventory service application
// It manages the lifecycle of the server, including configuration loading,
// dependency injection, service registration, and server lifecycle management.

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Set Gin mode based on config
	gin.SetMode(cfg.Server.Mode)

	// Initialize database connection
	db, err := mysql.NewMySQLConnection(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize Redis connection
	rdb, err := redis.NewRedisClient(cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer rdb.Close()

	// Initialize repositories
	inventoryRepo := mysql.NewInventoryRepository(db)
	inventoryHistoryRepo := mysql.NewInventoryHistoryRepository(db)
	warehouseRepo := mysql.NewWarehouseRepository(db)
	inventoryCache := redis.NewInventoryCache(rdb)

	// Initialize services
	inventoryService := service.NewInventoryServiceImpl(
		inventoryRepo,
		inventoryHistoryRepo,
		inventoryCache,
		warehouseRepo,
	)
	warehouseService := service.NewWarehouseServiceImpl(warehouseRepo)

	// Initialize gRPC handlers
	inventoryHandler := handler.NewInventoryHandler(inventoryService)
	warehouseHandler := handler.NewWarehouseHandler(warehouseService)

	// Start gRPC server
	grpcServer := startGRPCServer(cfg, inventoryHandler, warehouseHandler)
	defer grpcServer.GracefulStop()

	// Start HTTP server for health check and metrics
	httpServer := startHTTPServer(cfg)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		}
	}()

	// Wait for termination signal
	waitForTerminationSignal()

	log.Println("Inventory service shutting down gracefully...")
}

func startGRPCServer(cfg *config.Config, inventoryHandler *handler.InventoryHandler, warehouseHandler *handler.WarehouseHandler) *grpc.Server {
	// Create a listener on the specified port
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.GRPCPort))
	if err != nil {
		log.Fatalf("Failed to listen on port %d: %v", cfg.Server.GRPCPort, err)
	}

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Register services
	inventorypb.RegisterInventoryServiceServer(grpcServer, inventoryHandler)
	inventorypb.RegisterWarehouseServiceServer(grpcServer, warehouseHandler)

	// Register health check service
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("inventory-service", grpc_health_v1.HealthCheckResponse_SERVING)

	// Register reflection service on the server for service discovery
	reflection.Register(grpcServer)

	// Start the gRPC server
	go func() {
		log.Printf("gRPC server starting on port %d", cfg.Server.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC on port %d: %v", cfg.Server.GRPCPort, err)
		}
	}()

	return grpcServer
}

func startHTTPServer(cfg *config.Config) *http.Server {
	// Create a new Gin router
	router := gin.Default()

	// Register routes
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "UP",
			"service": "inventory-service",
		})
	})

	router.GET("/metrics", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})

	// Create HTTP server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.HTTPPort),
		Handler: router,
	}

	// Start HTTP server
	go func() {
		log.Printf("HTTP server starting on port %d", cfg.Server.HTTPPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	return server
}

func waitForTerminationSignal() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

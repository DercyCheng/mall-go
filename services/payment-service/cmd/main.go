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
	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"mall-go/services/payment-service/api/handler"
	"mall-go/services/payment-service/application/service"
	"mall-go/services/payment-service/domain/model"
	"mall-go/services/payment-service/domain/repository"
	"mall-go/services/payment-service/infrastructure/config"
	"mall-go/services/payment-service/infrastructure/persistence/mysql"
	redisCache "mall-go/services/payment-service/infrastructure/persistence/redis"
	"mall-go/services/payment-service/infrastructure/provider"
	"mall-go/services/payment-service/proto/paymentpb"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Setup logger
	setupLogger(cfg.Logging)

	// Initialize database
	db, err := setupDatabase(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize Redis
	rdb, err := setupRedis(cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Initialize repositories
	paymentRepo := mysql.NewPaymentRepository(db)
	refundRepo := mysql.NewPaymentRefundRepository(db)
	paymentCache := redisCache.NewPaymentCache(rdb)

	// Initialize payment providers
	paymentProviders := setupPaymentProviders(cfg.Payment)

	// Initialize services
	paymentService := service.NewPaymentServiceImpl(
		paymentRepo,
		paymentCache,
		refundRepo,
		paymentProviders,
	)
	refundService := service.NewRefundServiceImpl(
		paymentRepo,
		refundRepo,
		paymentProviders,
	)

	// Initialize gRPC handlers
	paymentHandler := handler.NewPaymentHandler(paymentService)
	refundHandler := handler.NewRefundHandler(refundService)

	// Start gRPC server
	go startGRPCServer(cfg.Server.GRPCPort, paymentHandler, refundHandler)

	// Start HTTP server
	go startHTTPServer(cfg.Server.HTTPPort, paymentHandler, refundHandler)

	// Wait for termination signal
	waitForTermination()
}

// setupLogger configures the logger
func setupLogger(cfg config.LoggingConfig) {
	// Configure logger based on configuration
	level := logger.Info
	switch cfg.Level {
	case "debug":
		level = logger.Info
	case "info":
		level = logger.Info
	case "warn":
		level = logger.Warn
	case "error":
		level = logger.Error
	}

	// Use gorm's logger for simplicity
	// In a real application, you might want to use a more advanced logger like zap or logrus
	gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: time.Second,
				LogLevel:      level,
				Colorful:      true,
			},
		),
	}
}

// setupDatabase initializes the database connection
func setupDatabase(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Set connection pool settings
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if cfg.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)
	}

	// Auto migrate the database schema
	err = db.AutoMigrate(
		&model.Payment{},
		&model.PaymentRefund{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// setupRedis initializes the Redis connection
func setupRedis(cfg config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}

// setupPaymentProviders initializes payment providers
func setupPaymentProviders(cfg config.PaymentConfig) map[string]service.PaymentProvider {
	providers := make(map[string]service.PaymentProvider)

	for _, providerCfg := range cfg.Providers {
		switch providerCfg.Name {
		case "alipay":
			providers["alipay"] = provider.NewAlipayProvider(provider.AlipayConfig{
				AppID:             providerCfg.AppID,
				MerchantPrivateKey: providerCfg.MerchantPrivateKey,
				AlipayPublicKey:   providerCfg.AlipayPublicKey,
				NotifyURL:         providerCfg.NotifyURL,
				ReturnURL:         providerCfg.ReturnURL,
				Sandbox:           providerCfg.Sandbox,
			})
		case "wechat":
			providers["wechat"] = provider.NewWeChatPayProvider(provider.WeChatPayConfig{
				AppID:      providerCfg.AppID,
				MerchantID: providerCfg.MerchantID,
				APIKey:     providerCfg.APIKey,
				NotifyURL:  providerCfg.NotifyURL,
				ReturnURL:  providerCfg.ReturnURL,
				Sandbox:    providerCfg.Sandbox,
			})
		}
	}

	return providers
}

// startGRPCServer starts the gRPC server
func startGRPCServer(port int, paymentHandler *handler.PaymentHandler, refundHandler *handler.RefundHandler) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen on port %d: %v", port, err)
	}

	grpcServer := grpc.NewServer()
	paymentpb.RegisterPaymentServiceServer(grpcServer, paymentHandler)
	paymentpb.RegisterRefundServiceServer(grpcServer, refundHandler)

	// Register reflection service for grpc_cli
	reflection.Register(grpcServer)

	log.Printf("gRPC server listening on port %d", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}

// startHTTPServer starts the HTTP server
func startHTTPServer(port int, paymentHandler *handler.PaymentHandler, refundHandler *handler.RefundHandler) {
	router := gin.Default()
	
	// TODO: Add HTTP endpoints for payment callbacks
	
	// Payment callback endpoint
	router.POST("/api/v1/payment/callback/:provider", func(c *gin.Context) {
		provider := c.Param("provider")
		
		// Parse form data or JSON data
		var params map[string]interface{}
		if c.Request.Header.Get("Content-Type") == "application/json" {
			if err := c.BindJSON(&params); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
				return
			}
		} else {
			if err := c.Request.ParseForm(); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form data"})
				return
			}
			
			params = make(map[string]interface{})
			for k, v := range c.Request.Form {
				if len(v) > 0 {
					params[k] = v[0]
				}
			}
		}
		
		// Create callback request
		callbackReq := &paymentpb.PaymentCallbackRequest{
			PaymentProvider: provider,
			Parameters:     make(map[string]string),
		}
		
		for k, v := range params {
			if s, ok := v.(string); ok {
				callbackReq.Parameters[k] = s
			} else {
				callbackReq.Parameters[k] = fmt.Sprintf("%v", v)
			}
		}
		
		// Process callback
		resp, err := paymentHandler.ProcessPaymentCallback(c.Request.Context(), callbackReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{"success": resp.Success, "message": resp.Message})
	})
	
	// Refund callback endpoint
	router.POST("/api/v1/refund/callback/:provider", func(c *gin.Context) {
		provider := c.Param("provider")
		
		// Parse form data or JSON data
		var params map[string]interface{}
		if c.Request.Header.Get("Content-Type") == "application/json" {
			if err := c.BindJSON(&params); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
				return
			}
		} else {
			if err := c.Request.ParseForm(); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form data"})
				return
			}
			
			params = make(map[string]interface{})
			for k, v := range c.Request.Form {
				if len(v) > 0 {
					params[k] = v[0]
				}
			}
		}
		
		// Create callback request
		callbackReq := &paymentpb.RefundCallbackRequest{
			PaymentProvider: provider,
			Parameters:     make(map[string]string),
		}
		
		for k, v := range params {
			if s, ok := v.(string); ok {
				callbackReq.Parameters[k] = s
			} else {
				callbackReq.Parameters[k] = fmt.Sprintf("%v", v)
			}
		}
		
		// Process callback
		resp, err := refundHandler.ProcessRefundCallback(c.Request.Context(), callbackReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{"success": resp.Success, "message": resp.Message})
	})

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	log.Printf("HTTP server listening on port %d", port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("failed to serve HTTP: %v", err)
	}
}

// waitForTermination waits for termination signals
func waitForTermination() {
	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
}

package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"mall-go/services/cart-service/api/handler"
	"mall-go/services/cart-service/application/service"
	"mall-go/services/cart-service/infrastructure/persistence/mysql"
	"mall-go/services/cart-service/infrastructure/persistence/redis"
	"mall-go/services/cart-service/infrastructure/config"
	"mall-go/services/cart-service/proto/cartpb"
)

func main() {
	// Initialize configuration
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	
	// Initialize MySQL database
	mysqlRepo, err := mysql.NewMySQLRepository(config.GetConfig().Database)
	if err != nil {
		log.Fatalf("Failed to initialize MySQL repository: %v", err)
	}

	// Initialize Redis client
	redisCache, err := redis.NewRedisCache(config.GetConfig().Redis)
	if err != nil {
		log.Fatalf("Failed to initialize Redis cache: %v", err)
	}
	
	// Initialize application services
	cartService := service.NewCartService(mysqlRepo, redisCache)
	
	// Initialize gRPC handlers
	cartHandler := handler.NewCartHandler(cartService)
	
	// Start gRPC server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor),
	)
	
	// Register services
	cartpb.RegisterCartServiceServer(grpcServer, cartHandler)
	reflection.Register(grpcServer)
	
	// Listen on configured port
	port := viper.GetString("server.grpc_port")
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}
	
	// Handle graceful shutdown
	go func() {
		log.Printf("Starting gRPC server on port %s", port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()
	
	// Wait for termination signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	
	log.Println("Shutting down server...")
	
	// Graceful shutdown
	grpcServer.GracefulStop()
	log.Println("Server exited properly")
}

func loggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	log.Printf("Method: %s, Duration: %s, Error: %v", info.FullMethod, time.Since(start), err)
	return resp, err
}

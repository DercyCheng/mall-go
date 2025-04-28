package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"

	"mall-go/pkg/registry"
	"mall-go/services/gateway-service/api/middleware"
	"mall-go/services/gateway-service/api/router"
	"mall-go/services/gateway-service/infrastructure/config"
)

var (
	configPath = flag.String("config", "configs/config.yml", "配置文件路径")
)

func main() {
	// 解析命令行参数
	flag.Parse()

	// 初始化配置
	absConfigPath, _ := filepath.Abs(*configPath)
	if err := config.LoadConfig(absConfigPath); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}
	log.Printf("成功加载配置: %s", absConfigPath)

	// 设置Gin模式
	if config.GlobalConfig.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 初始化Redis客户端(用于限流)
	var redisClient *redis.Client
	if config.GlobalConfig.RateLimit.Enabled && config.GlobalConfig.RateLimit.Type == "redis" {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", config.GlobalConfig.Redis.Host, config.GlobalConfig.Redis.Port),
			Password: config.GlobalConfig.Redis.Password,
			DB:       config.GlobalConfig.Redis.DB,
		})
		// 测试Redis连接
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if _, err := redisClient.Ping(ctx).Result(); err != nil {
			log.Printf("连接Redis失败: %v, 将使用内存限流", err)
			redisClient = nil
		} else {
			log.Printf("成功连接到Redis: %s:%d", config.GlobalConfig.Redis.Host, config.GlobalConfig.Redis.Port)
		}
	}

	// 初始化限流器
	middleware.InitRateLimiter(redisClient)

	// 创建Consul服务注册表
	consulRegistry, err := registry.NewConsulServiceRegistry(config.GlobalConfig.Consul.Address)
	if err != nil {
		log.Fatalf("初始化Consul客户端失败: %v", err)
	}

	// 创建Gin引擎
	r := gin.Default()

	// 创建API网关
	gateway := router.NewGateway(r, consulRegistry)
	gateway.SetupRoutes()

	// 服务注册
	var serviceRegistrationLifecycle *registry.ServiceRegistrationLifecycle
	if config.GlobalConfig.Consul.Address != "" {
		// 获取本机内网IP
		host, err := getLocalIP()
		if err != nil {
			log.Printf("获取本机IP失败: %v, 使用localhost代替", err)
			host = "localhost"
		}

		// 创建服务实例
		instanceID := fmt.Sprintf("%s-%s", config.GlobalConfig.Consul.ServiceName, uuid.New().String())
		serviceInstance := &registry.ServiceInstance{
			ID:      instanceID,
			Name:    config.GlobalConfig.Consul.ServiceName,
			Address: host,
			Port:    config.GlobalConfig.Server.Port,
			Tags:    config.GlobalConfig.Consul.Tags,
			Meta:    config.GlobalConfig.Consul.Meta,
		}

		// 创建服务注册生命周期管理
		serviceRegistrationLifecycle = registry.NewServiceRegistrationLifecycle(consulRegistry, serviceInstance)
		
		// 启动服务注册
		if err := serviceRegistrationLifecycle.Start(); err != nil {
			log.Printf("服务注册失败: %v", err)
		} else {
			log.Printf("服务已注册到Consul: %s", config.GlobalConfig.Consul.Address)
		}
	}

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.GlobalConfig.Server.Port),
		Handler:      r,
		ReadTimeout:  time.Duration(config.GlobalConfig.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(config.GlobalConfig.Server.WriteTimeout) * time.Second,
	}

	// 启动服务器
	go func() {
		log.Printf("API网关服务启动在端口: %d", config.GlobalConfig.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("启动服务失败: %v", err)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭服务...")

	// 从Consul注销服务
	if serviceRegistrationLifecycle != nil {
		if err := serviceRegistrationLifecycle.Stop(); err != nil {
			log.Printf("从Consul注销服务失败: %v", err)
		} else {
			log.Println("服务已从Consul注销")
		}
	}

	// 关闭HTTP服务器
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("服务器强制关闭: %v", err)
	}

	// 关闭Redis连接
	if redisClient != nil {
		if err := redisClient.Close(); err != nil {
			log.Printf("关闭Redis连接失败: %v", err)
		}
	}

	log.Println("服务已关闭")
}

// 获取本机内网IP
func getLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("没有找到有效的本地IP地址")
}
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
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"mall-go/pkg/auth"
	"mall-go/pkg/database"
	"mall-go/pkg/registry"
	"mall-go/services/product-service/api/handler"
	"mall-go/services/product-service/api/router"
	"mall-go/services/product-service/application/service"
	"mall-go/services/product-service/infrastructure/config"
	"mall-go/services/product-service/infrastructure/persistence/mysql"
)

var (
	configPath = flag.String("config", "configs/config.yml", "配置文件路径")
)

func main() {
	// 解析命令行参数
	flag.Parse()

	// 初始化配置
	if err := config.InitConfig(*configPath); err != nil {
		log.Fatalf("初始化配置失败: %v", err)
	}

	// 设置Gin模式
	gin.SetMode(config.GlobalConfig.Server.Mode)

	// 初始化JWT
	auth.InitJWTSecret(config.GlobalConfig.JWT.Secret)

	// 初始化数据库连接
	err := database.InitMySQL(
		config.GlobalConfig.Database.Username,
		config.GlobalConfig.Database.Password,
		config.GlobalConfig.Database.Host,
		fmt.Sprintf("%d", config.GlobalConfig.Database.Port),
		config.GlobalConfig.Database.Database,
	)
	if err != nil {
		log.Fatalf("初始化数据库连接失败: %v", err)
	}

	// 自动迁移数据库表结构
	if err := autoMigrateDB(); err != nil {
		log.Fatalf("自动迁移数据库表结构失败: %v", err)
	}

	// 初始化仓储
	productRepo := mysql.NewProductRepository()
	brandRepo := mysql.NewBrandRepository()
	categoryRepo := mysql.NewCategoryRepository()

	// 初始化服务
	productService := service.NewProductService(productRepo, brandRepo, categoryRepo)

	// 初始化处理器
	productHandler := handler.NewProductHandler(productService)

	// 初始化Gin引擎
	r := gin.Default()

	// 设置路由
	router.SetupRouter(r, productHandler)

	// 创建服务器
	srv := &http.Server{
		Addr:         config.GetServerAddress(),
		Handler:      r,
		ReadTimeout:  time.Duration(config.GlobalConfig.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(config.GlobalConfig.Server.WriteTimeout) * time.Second,
	}

	// 启动服务
	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("启动服务失败: %v", err)
		}
	}()

	// 服务注册
	serviceID := fmt.Sprintf("product-service-%s", uuid.New().String())
	log.Printf("服务ID: %s", serviceID)
	var serviceRegistrationLifecycle *registry.ServiceRegistrationLifecycle
	if config.GlobalConfig.Consul.Address != "" {
		// 创建Consul服务注册
		consulRegistry, err := registry.NewConsulServiceRegistry(config.GlobalConfig.Consul.Address)
		if err != nil {
			log.Printf("创建Consul客户端失败: %v, 将不会注册到服务发现中心", err)
		} else {
			// 获取本机内网IP
			host, err := getLocalIP()
			if err != nil {
				log.Printf("获取本机IP失败: %v, 使用localhost代替", err)
				host = "localhost"
			}

			// 创建服务实例
			serviceInstance := &registry.ServiceInstance{
				ID:      serviceID,
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
	}

	log.Printf("服务已启动，端口: %d", config.GlobalConfig.Server.Port)

	// 等待中断信号以优雅地关闭服务器（设置5秒的超时时间）
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭服务...")

	// 注销服务
	if serviceRegistrationLifecycle != nil {
		if err := serviceRegistrationLifecycle.Stop(); err != nil {
			log.Printf("从Consul注销服务失败: %v", err)
		} else {
			log.Println("服务已从Consul注销")
		}
	}

	// 设置5秒的超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("强制关闭服务: %v", err)
	}

	log.Println("服务已关闭")
}

// 自动迁移数据库表结构
func autoMigrateDB() error {
	// 获取数据库连接
	db := database.GetDB()

	// 自动迁移表结构
	return db.AutoMigrate(
		&mysql.ProductEntity{},
		&mysql.ProductAttributeEntity{},
		&mysql.ProductPromotionEntity{},
		&mysql.BrandEntity{},
		&mysql.CategoryEntity{},
	)
}

// 获取本机内网IP
func getLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if (err != nil) {
		return "", err
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("no valid local IP address found")
}